# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`go-webcore` (GoWebCore) is a Go library (`github.com/mondegor/go-webcore`, Go 1.26) providing base building blocks for web services. It is consumed by other projects, not run as an application. Code comments and docs are written in **English** or **Russian**; match this convention when editing existing files.

The library's design philosophy: each `mr*` package defines **interfaces** (e.g. `HttpRouter`, `Localizer`), and third-party libraries are wired in through **adapter sub-packages**. This lets consuming projects swap implementations. When adding support for a new third-party library, follow the adapter pattern rather than importing it into the interface package.

This module is tightly coupled to the sibling module `github.com/mondegor/go-core`: several interface packages that used to live here (`mraccess`, `mrprocess`, `mrrun`, `mridempotency`, the `Logger`/`mrtrace` types, the `mrstorage` interfaces) now live in **go-core**, and `go-webcore` keeps only the adapters that wire them to concrete libraries (Prometheus, chi, slog, etc.). When you can't find an interface here, it is almost certainly imported from `go-core`. Frequently used go-core packages: `errors`, `mrlog`, `mrtrace`, `mrtype` (+ `mrtype/parse`), `mrlocale`, `mraccess`, `util/*`. There is a commented-out `replace` directive in `go.mod` for local development against a sibling checkout.

## Commands

Development tasks go through the `mrcmd` tool (see https://github.com/mondegor/mrcmd) wrapped by the `Makefile`:

- `make lint` — run gofumpt/goimports/gci formatting + golangci-lint (config in `.golangci.yaml`, golangci-lint v2); `mrcmd go-dev lint` runs only golangci-lint
- `make test` — run all tests
- `make test-report` — tests with coverage report (`test-coverage-full.html`)
- `make generate` — run `go:generate`
- `make deps` — download dependencies
- `make deps-upgrade` — upgrade dependencies (`go get -u ./...` + `go mod tidy`)
- `make plantuml` — render `.puml` diagrams under `docs/` to images
- `make check-and-fix` — generate + format + lint + test + plantuml (full pre-commit pass; defined in `Makefile.mk`)

Without `mrcmd` installed:
- All tests: `go test ./...`
- Single package: `go test ./mrserver/request/parser/`
- Single test: `go test -run TestName -v ./mrserver/request/parser/`

There are currently no `//go:generate` directives or committed mocks in the tree; the `make generate` target stays for when they are reintroduced. Tests live in `_test` packages alongside the code they cover.

## Architecture

The repo is organized into independent `mr*` top-level packages:

- **mrserver** — HTTP serving, the largest package. Defines the central `HttpRouter` / `HttpController` / `HttpHandler` abstractions in `router.go`. Note `HttpHandlerFunc` returns an `error` (unlike `http.HandlerFunc`) so errors are handled centrally by middleware/adapters; `HttpHandler.Permission` is a string the access middleware uses to gate the handler. Request-statistics plumbing is in the root files (`RequestStat`/`RequestObserve` interfaces with Nop implementations, `RequestContainer`, `CacheableResponseWriter`). `response.go` defines the response interfaces (`ResponseEncoder`, `ResponseSender`, `FileResponseSender`, `ErrorResponseSender`, `ErrorStatusMapper` + `NewHttpErrorStatusMapper`); header names are in `header_keys.go`. Sub-packages: adapters `mrchi` (go-chi), `mrjulienrouter` (httprouter); `middleware/` (access, access-token, idempotency, observer, recover, request-id — these wire in `go-core`'s `mraccess`/`mridempotency`); `mrresp/` (responses), `mrjson/`, `mrprometheus/` (metrics via `ObserveRequest`), `mrrscors/` (CORS), `httpserver/`, `request/` + `request/parser/` (request parsers) + `request/validate/` (aggregate parser facades `Parser`/`ListParser`/`ContextParser` with their `Request*Parser` interfaces), `observe/` (request/response body+reader+writer wrappers), `stat/` (request logger/metrics/tracer implementations).
- **mrcore** — cross-cutting interfaces: `Localizer` in `locale.go`; app initialization helpers in `initing/` (HTTP handler/controller/module init) and `mrinit/` (`prometheus.go`).
- **mrstorage** — storage adapters; `mrprometheus/` holds `DBCollector`, a Prometheus collector that exports DB connection-pool stats via `go-core`'s `mrstorage.DBStatProvider`.
- **mrlog** — logging adapters. Currently `slog/sentry/` (a Sentry handler for `slog`). The primary `Logger`/`mrtrace` types live in `go-core`.
- **mrview** — validation; `mrplayvalidator` adapter wraps go-playground/validator/v10.
- **mrclient** — outbound clients: `MailSender`/`MessengerSender` interfaces in `client.go`; implementations `mail` (SMTP), `telegram`; plus `sentry`.
- **mrtests** — test helpers (`helpers/http_request.go`).
- **mrdebug** — debug utilities: `MultipartForm`/`MultipartFileHeader` (log multipart contents) and `PrepareNopServiceWithTimeoutToStart` (a nop service for testing process start/stop, e.g. with an app runner).

`examples/` contains standalone runnable `main.go` demos (`validator`, `smtpmail`). `mrserver/README.md` (Russian) is the authoritative description of the HTTP infrastructure (request lifecycle, access-check table, subsystems). `docs/` holds its C4 diagrams, laid out like go-core's: PlantUML sources (`docs/packages/c4/mrserver_*.puml`, `docs/diagrams/c4/hld.puml`, one shared `.iuml` per C4 element in `docs/components/c4/`, catalogued in `_list.puml`) and SVGs in `docs/resources/` referenced from `mrserver/README.md`; `make plantuml` regenerates them (it needs a TTY and an existing `docs/resources/`). `grafana-dashboards/` ships ready-made dashboards for the Prometheus metrics.

## Conventions worth knowing

- `.golangci.yaml` is strict: `gochecknoglobals` and `gochecknoinits` are enabled (no global vars or `init()` funcs), `godot` requires comments to end with a period, error sentinels must be `Err`-prefixed and error types `Error`-suffixed (`errname`). Run `make lint` before considering work done.
- Sentinel errors are built with go-core `errors` factories (e.g. `errors.NewInternalProto(...)`, see `mrclient/mail/message.go`), not `errors.New`.
- Adapter packages named `mrprometheus` exist under multiple parents (`mrserver/mrprometheus`, `mrstorage/mrprometheus`) — they are distinct packages; check the import path.
- `.env`, `Makefile.mk`, `tmp/`, `.cache/` are local and gitignored (`.env.dist` is the committed template for `mrcmd`).
- Project skills live in `.claude/skills/`: `go-style-guide` (code style), `check-go-style` (style compliance check), `audit-go-package` (deep package audit).
