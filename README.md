# go-railway-template

A modern Go REST API starter for [Railway](https://railway.com) — current Go, zero third-party dependencies, a ~10 MB static image, and the classic Go-on-PaaS failure modes already engineered out.

[![Deploy on Railway](https://railway.com/button.svg)](https://railway.app/new?github_url=https://github.com/lNamelessl/go-railway-template)

## What's inside

| Piece | Choice | Why |
|---|---|---|
| Go | **1.27** (pinned in `go.mod` and the Dockerfile) | Current stable release |
| Router | **Standard library** `net/http` (Go 1.22+ pattern routing) | `POST /items`, `GET /items/{id}` — the same patterns the popular `httprouter` package popularized, now built in. Zero dependencies to pin, audit, or go stale. |
| Runtime image | **`scratch`** | `CGO_ENABLED=0` static binary + CA certs. No OS, no package manager, no CVE surface from base-image libraries. |
| Store | In-memory (`sync.RWMutex`) | Runs with zero config; swap for a database in your own service. Restarts reset state — intentional for a starter. |

### The failure modes, engineered out

1. **`$PORT` binding** — Railway injects `PORT`; the server reads it (default `8080` for local dev) and logs the listen address at startup. Never hardcode a port.
2. **`/healthz`** — a root-level, redirect-free JSON healthcheck, wired to Railway via `healthcheckPath` in `railway.json`.
3. **Graceful shutdown** — SIGTERM/SIGINT trigger a 10-second drain (`server.Shutdown`), so in-flight requests finish instead of getting RST'd on redeploys.
4. **Static builds** — `CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w"` in the Dockerfile. The binary runs on an empty `scratch` image as a non-root user.

## Endpoints

```bash
# Healthcheck
curl https://<your-domain>/healthz
# {"status":"ok"}

# Create an item
curl -X POST https://<your-domain>/items \
  -H 'Content-Type: application/json' \
  -d '{"name":"hello"}'
# {"id":"<uuid-like hex>","name":"hello"}

# List items
curl https://<your-domain>/items
# {"items":[{"id":"...","name":"hello"}]}

# Get one item
curl https://<your-domain>/items/<id>
```

## Deploying to Railway

Click the deploy button above, or from scratch:

1. Push this repo to GitHub.
2. Create a Railway project from the repo — Railway auto-detects the root `Dockerfile`.
3. That's it. `PORT`, healthchecking, and restart policy come from `railway.json`.

### The PORT note

Railway assigns `PORT` as an environment variable and routes your public domain to it. This template binds to `$PORT` automatically. If you're porting your own service, make sure it does the same — hardcoding `:3000` is the #1 reason Go deploys go "success" but never answer on the public domain.

### Build flags

- `CGO_ENABLED=0` — pure-Go static binary; required for `scratch`, and it removes glibc/musl mismatches between build and runtime.
- `-trimpath` — removes local filesystem paths from the binary.
- `-ldflags="-s -w"` — strips the symbol table and DWARF debug info (smaller image).

If your service needs cgo (e.g. SQLite drivers), keep `CGO_ENABLED=0` and use a pure-Go driver, or switch the runtime stage to `alpine`.

## Swap-in guide: your existing Go service

1. **Handlers** — register your routes in `registerRoutes` (`api.go`). Stdlib patterns: `"GET /things/{id}"`, `r.PathValue("id")`.
2. **Keep `main.go`'s shape** — `$PORT` read, `http.Server` with timeouts, and the `signal.NotifyContext` graceful shutdown. Copy those ~40 lines into your `main` if you keep your own tree.
3. **Keep `/healthz`** (or rename it and update `healthcheckPath` in `railway.json`). Railway healthchecks must hit a real, redirect-free route.
4. **Dependencies** — add yours to `go.mod` and commit `go.sum`; the Dockerfile's `go mod download` layer picks them up. This template starts with zero.
5. **Store** — replace `store.go` with your persistence layer; nothing else touches it.

## Local development

```bash
go run .            # binds :8080
curl localhost:8080/healthz
```

Docker:

```bash
docker build -t go-railway-template .
docker run -p 8080:8080 go-railway-template
```
