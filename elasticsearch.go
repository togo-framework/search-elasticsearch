// Package elasticsearch is an Elasticsearch/OpenSearch driver for togo search.
// Both engines share the index/_doc and _search REST API, so one driver serves
// both (registered as "elasticsearch" and "opensearch"). HTTP only — no SDK.
package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/togo-framework/search"
	"github.com/togo-framework/togo"
)

func init() {
	factory := func(k *togo.Kernel) (search.Searcher, error) {
		base := os.Getenv("SEARCH_URL")
		if base == "" {
			return nil, errors.New("search-elasticsearch: SEARCH_URL not set")
		}
		return &client{
			base: strings.TrimRight(base, "/"),
			user: os.Getenv("SEARCH_USERNAME"),
			pass: os.Getenv("SEARCH_PASSWORD"),
			http: &http.Client{Timeout: 15 * time.Second},
		}, nil
	}
	search.RegisterDriver("elasticsearch", factory)
	search.RegisterDriver("opensearch", factory)
}

type client struct {
	base, user, pass string
	http             *http.Client
}

func (c *client) do(ctx context.Context, method, path string, body any) (*http.Response, error) {
	var r io.Reader
	if body != nil {
		buf, _ := json.Marshal(body)
		r = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.base+path, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.user != "" {
		req.SetBasicAuth(c.user, c.pass)
	}
	return c.http.Do(req)
}

func (c *client) Index(ctx context.Context, index, id string, doc map[string]any) error {
	resp, err := c.do(ctx, http.MethodPut, fmt.Sprintf("/%s/_doc/%s", index, id), doc)
	if err != nil {
		return err
	}
	return drain(resp)
}

func (c *client) Delete(ctx context.Context, index, id string) error {
	resp, err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/%s/_doc/%s", index, id), nil)
	if err != nil {
		return err
	}
	return drain(resp)
}

func (c *client) Search(ctx context.Context, index, query string, limit int) ([]search.Hit, error) {
	if limit <= 0 {
		limit = 20
	}
	body := map[string]any{
		"size":  limit,
		"query": map[string]any{"multi_match": map[string]any{"query": query}},
	}
	resp, err := c.do(ctx, http.MethodPost, "/"+index+"/_search", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("search-elasticsearch: status %d: %s", resp.StatusCode, string(b))
	}
	var out struct {
		Hits struct {
			Hits []struct {
				ID     string         `json:"_id"`
				Score  float64        `json:"_score"`
				Source map[string]any `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	hits := make([]search.Hit, 0, len(out.Hits.Hits))
	for _, h := range out.Hits.Hits {
		hits = append(hits, search.Hit{ID: h.ID, Score: h.Score, Doc: h.Source})
	}
	return hits, nil
}

func drain(resp *http.Response) error {
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("search-elasticsearch: status %d: %s", resp.StatusCode, string(b))
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	return nil
}
