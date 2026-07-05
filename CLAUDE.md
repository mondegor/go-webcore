# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`go-webcore` (GoWebCore) is a Go library (`github.com/mondegor/go-webcore`, Go 1.25) providing base building blocks for web services. It is consumed by other projects, not run as an application. Code comments and docs are written in **English** or **Russian**; match this convention when editing existing files.

The library's design philosophy: each `mr*` package defines **interfaces** (e.g. `HttpRouter`, `Localizer`), and third-party libraries are wired in through **adapter sub-packages**. This lets consuming projects swap implementations. When adding support for a new third-party library, follow the adapter pattern rather than importing it into the interface package.

This module is tightly coupled to the sibling module `github.com/mondegor/go-core`: several interface packages that used to live here (`mraccess`, `mrworker`, `mrrun`, `mridempotency`, the `Logger`/`mrtrace` types, the `mrstorage` interfaces) now live in **go-core**, and `go-webcore` keeps only the adapters that wire them to concrete libraries (Prometheus, chi, slog, etc.). When you can't find an interface here, it is almost certainly imported from `go-core`. There is a commented-out `replace` directive in `go.mod` for local development against a sibling checkout.

## Commands

Development tasks go through the `mrcmd` tool (see https://github.com/mondegor/mrcmd) wrapped by the `Makefile`:

- `make lint` — run gofumpt/goimports formatting + golangci-lint (config in `.golangci.yaml`, golangci-lint v2)
- `make test` — run all tests
- `make test-report` — tests with coverage report (`test-coverage-full.html`)
- `make generate` — run `go:generate`
- `make deps` — download dependencies
- `make check-and-fix` — generate + format + lint + test + plantuml (full pre-commit pass; defined in `Makefile.mk`)

To run a single test without mrcmd, standard Go works: `go test ./mrserver/request/parser/ -run TestName -v`.

There are currently no `//go:generate` directives or committed mocks in the tree; the `make generate` target stays for when they are reintroduced. Tests live in `_test` packages alongside the code they cover.

## Architecture

The repo is organized into independent `mr*` top-level packages:

- **mrserver** — HTTP serving, the largest package. Defines the central `HttpRouter` / `HttpController` / `HttpHandler` abstractions in `router.go`. Note `HttpHandlerFunc` returns an `error` (unlike `http.HandlerFunc`) so errors are handled centrally by middleware/adapters; `HttpHandler.Permission` is a string the access middleware uses to gate the handler. Request-statistics plumbing is in the root files (`RequestStat`/`RequestObserve` interfaces with Nop implementations, `RequestContainer`, `CacheableResponseWriter`). Sub-packages: adapters `mrchi` (go-chi), `mrjulienrouter` (httprouter); `middleware/` (access, access-token, idempotency, observer, recover, request-id — these wire in `go-core`'s `mraccess`/`mridempotency`); `mrresp/` (responses), `mrjson/`, `mrprometheus/` (metrics via `ObserveRequest`), `mrrscors/` (CORS), `httpserver/`, `request/` + `request/parser/` (request parsers), `observe/` (request/response body+reader+writer wrappers), `stat/` (request logger/metrics/tracer implementations).
- **mrcore** — cross-cutting interfaces: `Localizer` in `locale.go`; app initialization helpers in `initing/` (HTTP handler/controller/module init) and `mrinit/` (`prometheus.go`).
- **mrstorage** — storage adapters; `mrprometheus/` holds `DBCollector`, a Prometheus collector that exports DB connection-pool stats via `go-core`'s `mrstorage.DBStatProvider`.
- **mrlog** — logging adapters. Currently `slog/sentry/` (a Sentry handler for `slog`). The primary `Logger`/`mrtrace` types live in `go-core`.
- **mrview** — validation; `mrplayvalidator` adapter wraps go-playground/validator/v10.
- **mrclient** — outbound clients: `mail` (SMTP), `sentry`, `telegram`.
- **mrtests**, **mrdebug** — test helpers (`helpers/http_request.go`) and debug/no-op utilities.

`examples/` contains standalone runnable `main.go` demos (`validator`, `smtpmail`). `docs/` holds PlantUML diagrams (`make plantuml` regenerates SVGs). `grafana-dashboards/` ships ready-made dashboards for the Prometheus metrics.

## Conventions worth knowing

- `.golangci.yaml` is strict: `gochecknoglobals` and `gochecknoinits` are enabled (no global vars or `init()` funcs), `godot` requires comments to end with a period, error sentinels must be `Err`-prefixed and error types `Error`-suffixed (`errname`). Run `make lint` before considering work done.
- Sentinel/common errors are centralized rather than defined ad-hoc.
- Adapter packages named `mrprometheus` exist under multiple parents (`mrserver/mrprometheus`, `mrstorage/mrprometheus`) — they are distinct packages; check the import path.
- `.golangci.yaml.bak` and `.qwen/` are local cruft, not part of the project.
