# generoni

An opinionated, production-(almost)-ready template for Go web apps (a JSON API plus an embedded frontend) 'cos I got tired of copy and pasting. Still a WIP but should do 95% of the boilerplate for you :)

## Usage

1. Click **Use this template** to create your repo. A workflow renames the module path
   and replaces this README on first push, then deletes itself.

   > ❗
   > **Give it a moment to complete before cloning.** You'll know it's done when a new commit
   > titled "Rename template to <your_repo_name>" appears after "Initial commit".
2. Clone it. Optionally copy `.env.example` to `.env` to change the defaults (see `env.go`).
3. Run `make dev` to start the server with hot reload on `:8080`.
4. Try it:
   ```sh
   curl -i -X POST localhost:8080/api/v1/items \
     -H 'Content-Type: application/json' -d '{"title":"hi"}'
   curl localhost:8080/api/v1/items
   ```
   Open `http://localhost:8080` to see the placeholder frontend.
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
- The frontend in `web/dist` is embedded with `go:embed` and served by the same binary
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
web/                    frontend build output (web/dist), embedded into the binary
```

## Frontend

The server serves a web app from the same binary, so the browser talks to one origin (no CORS,
cookies just work). The API lives under `/api/v1`, and every other path serves `web/dist`, falling
back to `index.html` so client-side routes work.

- **Add your frontend:** put its source in `web/` (e.g. a Vite or React app) and make its build
  land in `web/dist/`, replacing the placeholder `index.html`. Vite does this by default. The Go
  build fails if `web/dist` is missing, so build the frontend first, commit `web/dist`, or add a
  Node stage to the `Dockerfile`.
- **Develop it:** run the frontend's dev server and proxy `/api` to `localhost:8080`.
- **Auth:** implement `middleware.Auth` (e.g. a session cookie) so the API checks who is calling.
- **CSP:** `middleware/security.go` sends `default-src 'self'`. Loosen it if you load scripts or
  images from elsewhere.

### Don't need a frontend? Remove it!

1. Delete `web/` and `internal/adapter/http/spa.go`.
2. Delete the `GET /api/` and `GET /` routes in `routes()` in `server.go`.
3. Remove `!web` from `.dockerignore`.
4. Optionally tighten the CSP in `middleware/security.go` to `default-src 'none'`.

## Deployment notes

- **Run one instance.** SQLite is a local file and the rate limiter is in memory.
- **Mount a volume at `/data`,** or the database is lost when the container is removed.
- **Back up the database.**  - the SQLite database it's just a file, so back it up however you like
- **Run behind exactly one trusted proxy.** The client IP is the right-most `X-Forwarded-For`
  value. Don't expose the app directly, recommended to use a reverse proxy like Caddy, or your
  cloud provider's load balancer.

## Licence

MIT
