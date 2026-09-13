# Go REST API Starter — Go on Railway, done right

[![Deploy on Railway](https://railway.com/button.svg)](https://railway.com/deploy/go-railway-template)

A modern Go REST API starter: current Go, zero third-party dependencies, a ~10 MB static image, and the classic Go-on-PaaS failure modes already engineered out — `$PORT` binding, a real healthcheck, and graceful shutdown on SIGTERM.

**What ships in the box:**

- **Go 1.27** pinned in `go.mod` and the multi-stage Dockerfile (`golang:1.27-alpine` builder → `scratch` runtime)
- **Standard library routing** (Go 1.22+ enhanced `net/http`) — `POST /items`, `GET /items/{id}` — no router dependency to go stale
- **`CGO_ENABLED=0` static binary** with `-trimpath -ldflags="-s -w"`, running as non-root on an empty `scratch` image: **~10 MB total**, no base-image CVE surface
- **`/healthz`** JSON healthcheck wired via `healthcheckPath` in the repo's `railway.json`
- **Graceful shutdown** — SIGTERM triggers a 10-second connection drain, so redeploys don't RST in-flight requests
- **Binds to `$PORT`** (Railway-injected) and logs the listen address at startup

**Endpoints:** `GET /healthz` · `POST /items` · `GET /items` · `GET /items/{id}` — all JSON. The in-memory store is a placeholder for your own persistence layer.

# Deploy and Host

Deploy with the one-click button above. Railway provisions a single service named `api`, builds it from the repo's root Dockerfile, injects `PORT`, waits for `/healthz` to return 200, and assigns a public domain — no configuration, no prompts. The first deploy typically goes green in under a minute; the image is ~10 MB.

## About Hosting

You host one Go service serving a small JSON REST API. The API is stateless by design: state lives in an in-memory store that resets on restart, so the service can be redeployed or scaled freely. Swapping in a database is your one integration point — replace `store.go`, add the database service, and reference `DATABASE_URL` in your code. Nothing else in the template touches storage.

The service runs as a non-root user on a `scratch` (empty) runtime image. There is no shell, no package manager, and no base OS library surface — just your static binary and CA certificates.

## Why Deploy

Writing a Go service is easy; getting one to *stay* green on a PaaS is where the classic failures live. This template fixes them up front:

- **Port binding** — the server reads `$PORT` instead of hardcoding one. This is the #1 reason Go deploys build fine but never answer on the public domain.
- **Healthcheck** — Railway probes `/healthz` before marking the deploy live, so broken releases never receive traffic.
- **Graceful shutdown** — SIGTERM drains connections for up to 10 seconds during redeploys.
- **Static builds** — `CGO_ENABLED=0` eliminates glibc/musl mismatches and enables the `scratch` runtime.

Zero third-party Go dependencies also means zero `go.sum` surprises, zero audit surface, and builds that keep working.

## Common Use Cases

- Starting a new Go microservice or REST API on Railway without boilerplate
- A reference for wiring an existing Go service to Railway (PORT handling, healthcheck, graceful shutdown, Dockerfile) — the README has a swap-in guide
- A teaching example of production-grade Go server setup: timeouts, signal handling, static builds
- A tiny, cheap always-on API (e.g. a webhook receiver or a personal API endpoint) — a service this small runs for a few dollars a month

## Dependencies for

None. The template has zero third-party Go dependencies and zero external services — no databases, no caches, no queues. Everything compiles from the standard library into one static binary.

### Deployment Dependencies

- A Railway account (the one-click button handles project creation)
- No deploy-form variables: the only environment variable is `PORT`, injected by Railway automatically
- No database, no storage volume, no plugins
