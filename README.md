# Redis Leaderboard

Go, Fiber, Redis, HTMX, and Templ leaderboard dashboard with a fixed control panel and automatic rank updates.

## Layout

- `internal/handler` for HTTP endpoints.
- `internal/redis` for Redis operations and seeding.
- `internal/model` for shared data types.
- `internal/view` for Templ components and rendering.

## Start

```bash
docker compose up -d
templ generate
go run .
```

## Environment

For local development use:

```bash
REDIS_ADDR=localhost:6379
SERVER_PORT=3000
```

When the Go app runs inside the Docker network, set:

```bash
REDIS_ADDR=redis:6379
```

## Verification

```bash
docker ps
redis-cli ping
redis-cli ZADD leaderboard 120 Alice
redis-cli ZADD leaderboard 90 Bob
redis-cli ZREVRANGE leaderboard 0 9 WITHSCORES
```

Expected output:

```text
Alice
120
Bob
90
```
