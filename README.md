# generoni

An opinionated, production-ready template for Go REST APIs 'cos I got tired of copy and pasting.

Click **Use this template** to create your own repo. A workflow renames the module path and this
title on the first push, then deletes itself.

## Stack

- `net/http` standard library, JSON in and out
- SQLite with sqlc and golang-migrate
- air for hot reload
- Multi-stage Docker build to a distroless image

## Layout

Hexagonal (ports and adapters). `core` holds the business logic, `adapter`
holds everything that talks to the outside world, and `container` wires them.

```
cmd/server/main.go                     config, logger, container, server, shutdown
internal/
  core/
    domain/                            entities, their rules and expected errors
      item.go                          Item, NewItem
      pagination.go                    Pagination, NewPagination
      errors.go                        Kind, Error, NewError, Err* values
    port/                              interfaces the core needs and offers
      item.go                          ItemService, ItemRepository
    service/                           business logic, one package per entity
      item/service.go
  adapter/
    config/environment/env.go          env to typed struct, fails fast
    container/container.go             builds DB, repositories, services, handlers
    logger/logger.go                   slog setup, JSON only
    http/
      server.go                        http.Server, routes, health
      respond/respond.go               Decode, JSON, Error: {"error":{"code"}},
                                       HTTP errors, Kind to status
      middleware/                      recover, requestid, requestlog,
                                       security, ratelimit, auth (stub), timeout
      handlers/                        one package per entity
        item/
          handler.go                   ItemHandler: Get, List, Create
          dto.go                       request and response shapes
          validation.go                path and query parsing, pagination
    storage/
      sqlite/
        db.go                          driver, pragmas, Open, Migrate
        item_repository.go             implements port.ItemRepository
        mapper.go                      queries.Item to domain.Item
        sqlc.yaml                      sqlc config
        query.sql                      hand written SQL, sqlc INPUT
        migrations/                    golang-migrate, embedded
          000001_init.up.sql
          000001_init.down.sql
        queries/                       sqlc OUTPUT. generated, committed,
          db.go                        never edited by hand
          models.go
          query.sql.go
data/                                  gitignored, holds the .db file
```

### Dependency direction

```
adapter/http/handlers --> core/port <-- core/service --> core/port <-- adapter/storage
                              |
                         core/domain
```

`core` never imports an adapter. Handlers depend on `port.ItemService`, the
service depends on `port.ItemRepository`, and only `container` knows which
concrete types fill them.

### Adding an entity

1. `core/domain/<entity>.go` with its constructor and rules, and its errors in `core/domain/errors.go`
2. `core/port/<entity>.go` with the service and repository interfaces
3. `core/service/<entity>/service.go`
4. A migration, queries in `query.sql`, then `sqlc generate -f internal/adapter/storage/sqlite/sqlc.yaml`
5. `adapter/storage/sqlite/<entity>_repository.go` and its mapper
6. `adapter/http/handlers/<entity>/` with `handler.go`, `dto.go`, `validation.go`
7. Wire it in `container.go` and add the routes in `server.go`

## Deployment

### Run one instance

SQLite is a local file and the rate limiter lives in memory, so run a single instance.
Several instances can't safely share one SQLite file, and each would count rate limits on
its own. To scale out, move to a shared database (e.g. Postgres) and a shared rate limiter
(e.g. Redis).

### Persist the database

Mount a volume at `/data`, for example `docker run -p 8080:8080 -v app-data:/data <image>`.
Without it the SQLite file lives inside the container and is lost when the container is
removed.

### Back up the database

The SQLite file is the only copy of your data. Take volume or disk snapshots with your
host, or stream it continuously to S3-compatible storage with Litestream.

### Run behind one trusted proxy

The client IP is the right-most `X-Forwarded-For` value, so the app assumes exactly one
trusted proxy in front of it that appends the address it saw. Never expose the app
directly. A client could then forge the header and bypass the rate limit.

## Licence

MIT
