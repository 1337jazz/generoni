# generoni

An opinionated, production-ready template for Go REST APIs 'cos I got tired of copy and pasting.

Click **Use this template** to create your own repo. The rename script handles, well, the renaming.

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
                                       security, ratelimit, timeout
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

## Licence

MIT
