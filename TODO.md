# TODO: production-readiness

## Must fix

- [ ] Add tests: handlers (`httptest`), SQLite repository, middleware.
- [ ] Add CI: `go vet`, `go test -race`, `govulncheck`, `golangci-lint`, `docker build`.
- [x] Remove or gate `GET /panic` in `internal/adapter/http/server.go`.
- [ ] Add `/readyz` that pings the DB. Point the ALB health check at it.
- [ ] Fix `clientip.FromRequest`: it returns the raw `X-Forwarded-For` chain.
      Take one entry (left-most or right-most, depending on trusted proxies).
- [ ] Document the SQLite limit: one writer, so one ECS task. Rate limiter is per instance too.

## Should have

- [ ] Slim `go.mod`: the `air` tool directive pulls in hugo, sass and protobuf.
      Move tools to a separate module.
- [ ] Extend the Makefile: `test`, `lint`, `vet`, `sqlc`, `docker`.
- [ ] Validate config on load: `RATE_LIMIT > 0`, `RATE_LIMIT_WINDOW > 0`,
      `REQUEST_TIMEOUT < WRITE_TIMEOUT`.
- [ ] Return JSON for mux 404 and 405 (currently plain text).
- [ ] Add an auth middleware stub and decide on CORS.
- [ ] Add metrics (`/metrics`) and/or OpenTelemetry. Honour inbound `X-Request-Id`
      or `X-Amzn-Trace-Id`.
- [ ] Add a `/v1` route prefix and an OpenAPI spec.
- [ ] Plan DB backups (Litestream or EFS snapshots). Note it in the README.

## Smaller

- [ ] `Timeout` middleware only sets a ctx deadline. Add a test or comment.
- [ ] On SIGTERM, fail readiness and wait briefly before `Shutdown`.
- [ ] Pin Dockerfile base images. Distroless has no curl, so skip `HEALTHCHECK`.
- [ ] Consider an HSTS header (skip if the proxy sets it).
- [ ] Run `make dev` to confirm the air banner is gone (`startup_banner = ""`).
