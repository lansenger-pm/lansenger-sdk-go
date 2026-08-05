# AGENTS.md — lansenger-sdk-go

Go SDK for the Lansenger Smart Bot API. Published via git tags on
`github.com/lansenger-pm/lansenger-sdk-go` (consumed through pkg.go modules).

## How to run

- Tests: `go test ./...`
- Build CLI: `go build ./cmd/lansenger`
- Publish: push a git tag `vx.y.z` (e.g. `git tag -a v0.10.2 -m "..." && git push origin v0.10.2`)

## Tech stack

Go (stdlib only for the SDK; cobra for the CLI in `cmd/lansenger`). No build step —
Go source is consumed directly via the module path.

## Layout

- `*.go` — SDK source (`client.go`, `config.go`, `auth.go`, `chats.go`, `version.go`)
- `cmd/lansenger/` — CLI
- `*_test.go` — go test suite

## Release rules — CRITICAL

### Version numbers live in MULTIPLE places — update ALL of them together

Before tagging a release, every one of these must hold the same version:

| File | Symbol |
|------|--------|
| `version.go` | `const Version = "x.y.z"` |
| git tag | `vx.y.z` (annotated, pushed to origin) |

`cmd/lansenger/main.go` reads `lansenger.Version` for `--version`, so it stays in sync
automatically — but the tag must point at the commit that has the updated `version.go`.
If you tag before bumping `version.go`, delete and recreate the tag on the correct commit.

### NEVER tag/publish without a full green test run

`go test ./...` MUST pass (0 failures) before `git tag`. No exceptions — not "the
failure is just a version string", not "I'll fix it in the next release". A red test
run means the release is not ready.

### CI-driven publishing

Releases are published exclusively by pushing a git tag (`vx.y.z`). The
`Release` GitHub Actions workflow verifies `version.go` matches the tag and
runs `go test ./...`; pkg.go.dev indexes the tag automatically. There is no
manual upload step for Go modules — the tag itself IS the release.

### Pass-through (external token) mode

- `NewClientWithConfig(cfg)` accepts a `Config` with `AppToken` set — used by the CLI
  and for external-token scenarios.
- `NewClientWithToken(appToken, userToken)` is the convenience constructor for
  pass-through mode (reads `LANSENGER_API_GATEWAY_URL` from env).
- `NewClient(appID, appSecret)` is the standard-mode constructor (SDK auto-refreshes
  appToken). `appID`/`appSecret` may be empty when `AppToken` is set on the config.
- `TokenManager.GetToken` returns `config.AppToken` directly in external mode (no refresh).
- `UserTokenManager` honors an externally-provided `config.UserToken` directly instead
  of loading from the credential store.

## Current status

v0.11.0. High-risk write gate (exit 10 + `--yes`/`--dry-run`) added to CLI; `NewClientWithToken`; `UserTokenManager` honors external userToken.
