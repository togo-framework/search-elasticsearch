# search-elasticsearch

Elasticsearch / OpenSearch driver for togo [search](https://github.com/togo-framework/search)
(both share the index/_search API).

```bash
togo install togo-framework/search
togo install togo-framework/search-elasticsearch
```

Set `SEARCH_DRIVER=elasticsearch` (or `opensearch`) and `SEARCH_URL` (e.g. http://localhost:9200),
optional `SEARCH_USERNAME`/`SEARCH_PASSWORD`.
