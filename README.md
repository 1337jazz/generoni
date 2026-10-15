# generoni

An opinionated, production-ready (single instance) template for Go REST APIs 'cos I got tired of copy and pasting. Still a WIP but should do 95% of the boilerplate for you :)

## Usage

1. Click **Use this template** to create your repo. A workflow renames the module path and
   this title on the first push, then deletes itself.
2. Clone it. Optionally copy `.env.example` to `.env` to change the defaults (see `env.go`).
3. Run `make dev` to start the server with hot reload on `:8080`.
4. Try it:
   ```sh
   curl -i -X POST localhost:8080/v1/items \
     -H 'Content-Type: application/json' -d '{"title":"hi"}'
   curl localhost:8080/v1/items
   ```
5. Build and run in Docker:
   ```sh
   make docker
   docker run -p 8080:8080 -v app-data:/data server
   ```

Other `make` targets: `test`, `vet`, `sqlc`, `migrate-new name=<name>`, `migrate-up`,
`migrate-down`. They all do what you'd expect.

To add your own resource, follow `item` through each layer below.

## Stack

- `net/http` standard library, with `alice` for middleware chains and `httprate` for rate limiting
- SQLite (`modernc.org/sqlite`), `sqlc` for queries, `golang-migrate` for migrations
- `caarlos0/env` for config, `slog` for JSON logs
- `air` for hot reload, a multi-stage Docker build to a distroless image

## Structure

Hexagonal (ports and adapters). `core` never imports `adapter`, and `container` wires them.

```
cmd/server/             entry point: config, logger, container, server, shutdown
internal/
  core/
    domain/             entities, rules and errors
    port/               interfaces the core needs and offers
    service/            business logic
  adapter/
    http/               server, routes, middleware, handlers
    storage/sqlite/     database, migrations, sqlc queries
    config/             env vars to a typed struct
    logger/             slog setup
    container/          builds and connects everything
```

## Deployment notes

- **Run one instance.** SQLite is a local file and the rate limiter is in memory.
- **Mount a volume at `/data`,** or the database is lost when the container is removed.
- **Back up the database.** It is the only copy of your data.
- **Run behind exactly one trusted proxy.** The client IP is the right-most `X-Forwarded-For`
  value. Don't expose the app directly.

## Licence

MIT
