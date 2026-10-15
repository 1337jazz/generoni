# TODO: production-readiness

## Must fix

- [x] Add starter CI: `go vet`, `go test -race`, `docker build` (`.github/workflows/ci.yml`).
- [x] Remove or gate `GET /panic` in `internal/adapter/http/server.go`.
- [x] Fix `clientip.FromRequest`: take the right-most `X-Forwarded-For` value (one proxy).
- [x] Document the SQLite limit: one writer, so run one instance. Rate limiter is per instance too.

## Should have

- [x] Slim `go.mod`: run air via `go run ...@version` in the Makefile, not a `tool` directive.
- [ ] Extend the Makefile: `test`, `lint`, `vet`, `sqlc`, `docker`.
- [ ] Validate config on load: `RATE_LIMIT > 0`, `RATE_LIMIT_WINDOW > 0`,
      `REQUEST_TIMEOUT < WRITE_TIMEOUT`.
- [ ] Return JSON for mux 404 and 405 (currently plain text).
- [ ] Add an auth middleware stub and decide on CORS.
- [ ] Add metrics (`/metrics`) and/or OpenTelemetry. Honour an inbound request ID
      header (e.g. `X-Request-Id`).
- [ ] Add a `/v1` route prefix and an OpenAPI spec.
- [ ] Plan DB backups (Litestream or volume snapshots). Note it in the README.
- [ ] README: mount a volume at `/data` (`-v app-data:/data`). Without it data is lost.
- [ ] README: assume exactly one trusted proxy, never expose the app directly.

## Smaller

- [ ] Require the `check` and `docker` CI jobs via branch protection on `main`.
- [ ] Optional CI extras: `govulncheck`, `golangci-lint`.
- [ ] `Timeout` middleware only sets a ctx deadline. Add a test or comment.
- [ ] Add `/readyz` (DB ping) only if deploying behind a load balancer or orchestrator.
- [ ] On SIGTERM, wait briefly before `Shutdown` so a proxy can drain (if behind one).
- [ ] Pin Dockerfile base images. Distroless has no curl, so skip `HEALTHCHECK`.
- [ ] Consider an HSTS header (skip if the proxy sets it).
- [ ] Run `make dev` to confirm the air banner is gone (`startup_banner = ""`).
