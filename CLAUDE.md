# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`go-webcore` (GoWebCore) is a Go library (`github.com/mondegor/go-webcore`, Go 1.25) providing base building blocks for web services. It is consumed by other projects, not run as an application. Code comments and docs are written in **English** or **Russian**; match this convention when editing existing files.

The library's design philosophy: each `mr*` package defines **interfaces** (e.g. `HttpRouter`, `Localizer`, `RightsSource`), and third-party libraries are wired in through **adapter sub-packages**. This lets consuming projects swap implementations. When adding support for a new third-party library, follow the adapter pattern rather than importing it into the interface package.

## Commands

Development tasks go through the `mrcmd` tool (see https://github.com/mondegor/mrcmd) wrapped by the `Makefile`:

- `make lint` — run gofumpt/goimports formatting + golangci-lint (config in `.golangci.yaml`, golangci-lint v2)
- `make test` — run all tests
- `make test-report` — tests with coverage report (`test-coverage-full.html`)
- `make generate` — run `go:generate` (regenerates gomock mocks)
- `make deps` — download dependencies
- `make check-and-fix` — generate + format + lint + test + plantuml (full pre-commit pass; defined in `Makefile.mk`)

To run a single test without mrcmd, standard Go works: `go test ./mrworker/process/schedule/ -run TestName -v`.

Mocks are generated with `golang/mock`'s `mockgen` via `//go:generate` directives in `*_test.go` files, output to local `mock/` sub-packages (e.g. `mrworker/process/schedule/mock/`). Regenerate with `make generate`, not by hand.

## Architecture

The repo is organized into independent `mr*` top-level packages:

- **mrcore** — cross-cutting interfaces like `Localizer`; app initialization helpers in `initing/` and `mrinit/`.
- **mrserver** — HTTP serving. Defines the central `HttpRouter` / `HttpController` / `HttpHandler` abstractions in `router.go`. Note `HttpHandlerFunc` returns an `error` (unlike `http.HandlerFunc`) so errors are handled centrally by middleware/adapters. Adapters: `mrchi` (go-chi), `mrjulienrouter` (httprouter); plus `middleware/`, `mrresp/` (responses), `mrjson/`, `mrprometheus/` (metrics via `ObserveRequest`), `mrrscors/` (CORS), `httpserver/`, `request/` (request parsers).
- **mrworker** — concurrent background processing. Defines `Task`/`Job` (scheduled work) and `MessageConsumer`/`MessageHandler`/`MessageBatchHandler` (queue processing, PULL and PUSH models). `process/` holds the runnable services: `schedule` (TaskScheduler), `consume`, `collect`, `signal`, `onstartup`. `period/strategy/` defines task scheduling periods.
- **mrrun** — process lifecycle. `AppRunner` runs a group of `Process` services concurrently with start-order dependencies (`Add`, `AddFirstProcess`, `AddNextProcess`) and graceful shutdown; health/readiness probes in `health_probe.go`/`prepare_probes.go`.
- **mraccess** — role-based access control: `RightsSource`, `RightsChecker`, roles, permissions, privileges. Used by `mrserver/middleware` to gate handlers via `HttpHandler.Permission`.
- **mrlog** — logging; `slog` adapter (with `sentry/`). Note: the primary `Logger`/`mrtrace` types are pulled from the sibling module `github.com/mondegor/go-sysmess`.
- **mrview** — validation; `mrplayvalidator` adapter wraps go-playground/validator/v10.
- **mrclient** — outbound clients: `mail`, `sentry`, `telegram`.
- **mridempotency** — idempotency `Provider`/`Responser` interfaces with `nopprovider`/`nopresponser` no-op implementations.
- **mrtests**, **mrdebug** — test helpers and debug/no-op utilities.

`examples/` contains standalone runnable `main.go` demos (status, validator, shutdown, smtpmail). `docs/` holds PlantUML diagrams (`make plantuml` regenerates SVGs). `grafana-dashboards/` ships ready-made dashboards for the Prometheus metrics.

## Conventions worth knowing

- `.golangci.yaml` is strict: `gochecknoglobals` and `gochecknoinits` are enabled (no global vars or `init()` funcs), `godot` requires comments to end with a period, error sentinels must be `Err`-prefixed and error types `Error`-suffixed (`errname`). Run `make lint` before considering work done.
- Generics are used throughout `mrworker` (`MessageConsumer[T]`, etc.).
- Sentinel/common errors are centralized rather than defined ad-hoc.
- `.golangci.yaml.bak` and `.qwen/` are local cruft, not part of the project.
