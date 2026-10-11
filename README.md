# generoni

An opinionated, production-ready template for Go web apps 'cos I got tired of copy and pasting.

Click **Use this template** to create your own repo. The rename script handles, well, the renaming.

## Stack

- `net/http` standard library
- `html/template` for server rendering
- htmx + Alpine.js (CSP build)
- Tailwind CSS via the standalone CLI
- SQLite with sqlc and golang-migrate
- air for hot reload
- Multi-stage Docker build to a distroless image

## Layout

Ports and adapters, following Ben Johnson's [Standard Package Layout][spl].
Entries marked `TODO` are planned, not written yet.

```
sqlc.yaml                          sqlc config                            TODO
cmd/server/
  main.go                          wiring: config, logger, server, routes
internal/
  config/config.go                 env to typed struct, fails fast
  logging/logging.go               slog setup
  middleware/                      recover, requestid, requestlog,
                                   security, ratelimit, cache, nodirlist
  handlers/
    health.go  pages.go            handlers with no dependencies
    item.go                        handlers that take the service        TODO
  item/                            THE DOMAIN. imports nothing of ours   TODO
    item.go                        Item, Title (validated), ErrNotFound
    service.go                     Service: Create, List, Get
    repo.go                        type Repo interface, the port
    service_test.go                unit tests against a fake Repo
  storage/
    sqlite/                        the driven adapter                    TODO
      sqlite.go                    driver, pragmas, migrate on boot
      item.go                      implements item.Repo
      mapper.go                    db.Item to item.Item and back
      item_test.go                 against a temp .db file
      queries.sql                  hand written SQL, sqlc INPUT
      migrations/                  golang-migrate, embedded
        0001_init.up.sql
        0001_init.down.sql
      db/                          sqlc OUTPUT. generated, committed,
        db.go                      never edited by hand
        models.go
        queries.sql.go
web/
  web.go                           go:embed of templates and static
  templates/                       base, home, 404, 500
  static/                          style.css
data/                              gitignored, holds the .db file        TODO
```

### Dependency direction

```
handlers --> item <-- storage/sqlite
              ^
       imports nothing
```

`item` never imports `database/sql`, sqlc or `net/http`. Both arrows point
inward, so the domain does not move when the database or the transport does.

### Deviations from Standard Package Layout

Two, both deliberate:

1. **Domain types live in `internal/item`, not the root package.** The root
   package would be named after the repository, and a repository called
   `my-app` gives `package my-app`, which is not a valid Go identifier. The
   rename script makes the canonical form impossible.
2. **Adapters live under `internal/`, not at the root.** Nobody should import
   this app's adapters from outside it.

[spl]: https://www.gobeyond.dev/standard-package-layout/

## Licence

MIT
