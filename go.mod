module github.com/ecastellanosr/rssagg

go 1.23.2

replace github.com/ecastellanosr/rssagg/internal/config v0.0.0 => ./internal/config

require (
	github.com/ecastellanosr/rssagg/internal/config v0.0.0
	github.com/google/uuid v1.6.0
	github.com/lib/pq v1.10.9
)
