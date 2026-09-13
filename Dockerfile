# ---- Build stage: pinned current Go on Alpine ----
FROM golang:1.27.1-alpine AS build

WORKDIR /src

# Cached dependency download layer (no-op for this template — zero third-party
# deps — but keeps the standard layering when you add requirements to go.mod).
COPY go.mod ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 produces a fully static binary that runs on scratch.
# -trimpath and -ldflags="-s -w" strip build paths and debug info.
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /bin/server .

# ---- Runtime stage: scratch (empty image) + static binary + CA certs ----
FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /bin/server /server

# Non-root numeric UID/GID (scratch has no /etc/passwd to name a user).
USER 65532:65532

EXPOSE 8080
ENTRYPOINT ["/server"]
