# Go REST API Starter for Railway — Implementation Plan

## Research results (Step 1, done)

- **Current stable Go: 1.27.1** (per go.dev/dl; previous line 1.26.8). Pin the builder to `golang:1.27.1-alpine`.
- **Router: standard library `net/http` (Go 1.22+ enhanced routing).** Justification, which will go in the README: the incumbent uses `julienschmidt/httprouter`, which has been effectively unmaintained for years; Go 1.22+ stdlib routing now supports the same method+path patterns (`POST /items`, `GET /items/{id}`) with equivalent trie-based performance. Going stdlib means `go.mod` has **zero third-party dependencies** — the strongest possible answer to "pinned deps": nothing to pin, audit, or break. Swap-in cost for httprouter users is trivial (patterns are nearly identical) and documented.
- **Failure modes engineered out** (all four from the brief):
  1. `PORT` from env with `8080` default, logged at startup
  2. `/healthz` JSON handler, wired as Railway's `healthcheckPath`
  3. Graceful shutdown: `signal.NotifyContext` on SIGTERM/SIGINT → `server.Shutdown(10s timeout)`
  4. `CGO_ENABLED=0` static build with `-trimpath -ldflags="-s -w"`
- Runtime image: **`scratch`** — static binary + `ca-certificates` copied from the builder, non-root numeric `USER`. This maximizes the "tiny image" selling point; fallback to `alpine` if anything runtime-related misbehaves on Railway.

## Scaffold (Step 2)

Repo: public GitHub **`go-railway-template`** via `github-repo-setup`, scaffolded in the (empty) working directory `C:\Users\HomePC\Documents\code\go`, excluding `.mimosa/` via `.gitignore`.

```
Dockerfile          # multi-stage: golang:1.27.1-alpine build → scratch runtime
railway.json        # DOCKERFILE builder, healthcheckPath /healthz, ON_FAILURE restart policy
go.mod              # go 1.27, zero requires
main.go             # PORT, logging, http.Server (read/write timeouts), graceful shutdown
api.go              # handlers: healthz + items CRUD, JSON encode/decode helpers
store.go            # in-memory store (sync.RWMutex), items with UUID ids
.gitignore / .dockerignore
README.md
```

**API surface (3 JSON endpoints + healthz):**
- `GET /healthz` → `{"status":"ok"}` (no redirects, root-level — per the Railway trailing-slash gotcha)
- `POST /items` `{"name":"..."}` → 201 + created item
- `GET /items` → `{"items":[...]}` (in-memory; documented as a starter store to replace)
- `GET /items/{id}` → item or 404 JSON error

**Dockerfile:** pinned builder, `go mod download` layer kept for swap-in users, `CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w"`, scratch runtime with ca-certificates and `USER 65532:65532`. Note: no container-level `HEALTHCHECK` instruction (no curl/wget in scratch) — Railway's platform healthcheck via `railway.json` is the mechanism.

**README:** what it is, deploy button, endpoint examples (curl), swap-in guide for an existing Go service (drop your handlers into `api.go`, or replace `main.go`'s mux — keep PORT/healthz/shutdown), PORT note (Railway injects it, bind to `$PORT`), build flags explained, image-size badge/number.

## Deploy & debug (Step 3) — `railway-deploy` skill

1. `railway whoami` (workspace check) → `railway init` → `railway up` in background, poll build logs.
2. `railway domain` → curl `/healthz` (expect 200), then POST/GET `/items` round-trip.
3. Confirm newest deployment is SUCCESS; trigger a redeploy and confirm a clean restart (in-memory state resets — documented behavior).
4. Record final **image size** (from build output / `railway ssh` binary size) for the README and listing.

## Publish & verify (Step 4) — `railway-template-publish` + `railway-template-listing` skills

1. `railway service source connect --repo <owner>/go-railway-template` (repo-sourced, not image-sourced), let the triggered build settle.
2. `railway templates create --json` → audit `serializedConfig`: repo source, public domain captured, **zero required variables** (zero-prompt design — there are no secrets or literals to strip; the app has no env vars besides Railway-injected `PORT`).
3. Write `TEMPLATE_OVERVIEW.md` with the deploy button and the five exact required headings (`# Deploy and Host`, `## About Hosting`, `## Why Deploy`, `## Common Use Cases`, `## Dependencies for` + `### Deployment Dependencies`).
4. `railway templates publish <id> --category Starters --description "..." --readme-file TEMPLATE_OVERVIEW.md` (description ≤75 chars, e.g. "Go REST API starter for Railway — static binary, healthcheck, one click"), then swap the repo README button to the published `railway.com/deploy/<code>` URL and re-publish the listing.
5. **Verify 2× from the marketplace:** two fresh `railway init` + `railway deploy -t <code>` projects (zero `-v` flags expected — zero-prompt), each checked for: deployment SUCCESS, `/healthz` 200, POST+GET `/items` round-trip. Acceptance = 2/2 green.

## Notes

- The Mimosa SSRF constraint (validate outbound URLs, reject loopback/private hosts) is satisfied by construction: the API makes no outbound requests.
- Step 5 (re-scan in a month) is out of scope for this session — I'll note it in the final summary.

## Acceptance mapping

- [ ] Published template live; 2/2 fresh one-click deploys succeed → Step 4.5
- [ ] PORT binding + `/healthz` + graceful shutdown → `main.go`, verified live
- [ ] Static minimal image; deps pinned (zero third-party); listing complete → Dockerfile + go.mod + TEMPLATE_OVERVIEW.md