# search-elasticsearch — documentation

  <img src=".github/assets/togo-mark.svg" alt="togo" height="64" />

## Overview

Package elasticsearch is an Elasticsearch/OpenSearch driver for togo search.
Both engines share the index/_doc and _search REST API, so one driver serves
both (registered as "elasticsearch" and "opensearch"). HTTP only — no SDK.

## Install

```bash
togo install togo-framework/search-elasticsearch
```

Set `SEARCH_DRIVER=elasticsearch`.

## Configuration

Environment variables read by this plugin (extracted from the source):

| Env var | Notes |
|---|---|
| `G` | _see provider docs_ |
| `SEARCH_PASSWORD` | _see provider docs_ |
| `SEARCH_URL` | _see provider docs_ |
| `SEARCH_USERNAME` | _see provider docs_ |

## Usage

```go
s := k.Search
s.Index(ctx, "posts", doc)
hits, _ := s.Search(ctx, "posts", "query")
```

## Links

- Marketplace: https://to-go.dev/marketplace
- Source: https://github.com/togo-framework/search-elasticsearch
- README: ../README.md
