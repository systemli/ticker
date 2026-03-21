# AGENTS.md

## Project Overview

Go REST API (Gin + GORM) for distributing short messages during time-sensitive events.
Bridges messages to Telegram, Mastodon, Bluesky, and Signal.

- **Language:** Go (see `go.mod` for version)
- **Web framework:** Gin
- **ORM:** GORM (SQLite, PostgreSQL, MySQL)
- **CLI:** Cobra

## Build / Test / Lint Commands

| Task | Command |
|------|---------|
| Build | `go build` |
| Test all | `go test ./...` |
| Test with coverage | `go test -coverprofile=coverage.txt -covermode=atomic ./...` |
| Single test suite | `go test -run TestTickerTestSuite ./internal/api/...` |
| Single subtest | `go test -run TestTickerTestSuite/TestGetTickers ./internal/api/...` |
| Single simple test | `go test -run TestContains ./internal/util/...` |
| Lint | `golangci-lint run --timeout 10m` |
| Format | `gofmt -w .` and `goimports -w .` |
| Tidy deps | `go mod tidy` |
| Generate mocks | `mockery` (configured in `.mockery.yml`) |

## Project Structure

All application code lives in `internal/`:

- `cmd/` - CLI entry points (Cobra)
- `internal/api/` - HTTP handlers (Gin), middleware, response types
- `internal/api/middleware/` - One subpackage per middleware
- `internal/api/response/` - Response DTOs and serializers
- `internal/bridge/` - Message bridging (Telegram, Mastodon, Bluesky, Signal)
- `internal/storage/` - Data models, GORM storage interface + implementation, migrations
- `internal/config/` - YAML + env var configuration
- `internal/cache/` - In-memory cache
- `internal/logger/` - Structured logging (logrus)
- `testdata/` - Test fixtures

## Code Style

### Imports

Two groups separated by a blank line:

1. Standard library
2. Everything else (third-party and internal mixed, alphabetically sorted)

Use import aliases only for naming collisions.

### Naming

- **Files:** `snake_case.go`, tests colocated as `snake_case_test.go`
- **Handlers:** HTTP verb prefix: `GetTickers`, `PostTicker`, `PutTicker`, `DeleteTicker`
- **Constructors:** `NewTicker()`, `NewSqlStorage()`, `NewCache()`
- **Short vars in narrow scope:** `c` for `*gin.Context`, `h` for handler, `s` for suite/storage, `err` for errors
- **Request types:** suffix `Param` (`TickerParam`, `MessageParam`)
- **Response types:** separate structs in `internal/api/response/`
- **Constants:** typed, PascalCase, grouped in `const` blocks

### Error Handling

- Early return on error (guard clause pattern)
- Lowercase error messages, no trailing punctuation
- Translate errors to structured API responses: `response.ErrorResponse(code, msg)`
- Log with context: `log.WithError(err).WithField("key", val).Error("description")`
- Don't log AND return the same error — choose one

### Formatting

- `gofmt` / `goimports` for all code
- Keep happy path left-aligned, return early to reduce nesting
- Favor clarity and simplicity over cleverness

## Testing Conventions

### Testify Suites (primary pattern for handler/integration tests)

```go
type TickerTestSuite struct {
    suite.Suite
    w     *httptest.ResponseRecorder
    ctx   *gin.Context
    store *storage.MockStorage
}

func (s *TickerTestSuite) SetupTest() { /* reset state per test */ }

func (s *TickerTestSuite) TestGetTickers() {
    s.Run("when not authorized", func() { /* ... */ })
    s.Run("happy path", func() { /* ... */ })
}

func TestTickerTestSuite(t *testing.T) {
    suite.Run(t, new(TickerTestSuite))
}
```

### Simple Tests (for models/utilities)

```go
func TestContains(t *testing.T) {
    assert.True(t, Contains([]int{1, 2}, 1))
    assert.False(t, Contains([]int{1, 2}, 3))
}
```

### Mocking

- **Mockery** generates mocks from interfaces (config: `.mockery.yml`)
- Pattern: `s.store.On("Method", mock.Anything).Return(val).Once()`
- Always call `s.store.AssertExpectations(s.T())` at end of subtests
- **gock** for HTTP mocking of external API calls
- Subtest names: human-readable scenarios (`"when storage returns error"`, `"happy path"`)

## HTTP Handler Patterns

Central `handler` struct holds all dependencies:

```go
type handler struct {
    config   config.Config
    storage  storage.Storage
    bridges  bridge.Bridges
    cache    *cache.Cache
    realtime *realtime.Engine
}
```

Handler methods follow this flow:

1. Extract entity from gin context (via middleware prefetch)
2. Validate/bind request body
3. Perform storage operation
4. Return `response.SuccessResponse(...)` or `response.ErrorResponse(...)`

## Logging

- Logrus with per-package logger: `var log = logger.GetWithPackage("api")`
- Structured fields: `.WithError(err)`, `.WithField("key", val)`
- Levels: `Error`, `Warn`, `Info`, `Fatal`

## Commits and PRs

Use **Gitmoji** in commit messages and PR titles:

- ✨ New feature
- 🐛 Bug fix
- ♻️ Refactor
- ✅ Add/update tests
- ⬆️ Upgrade dependencies
- 📝 Documentation
- 🧹 Chore/maintenance

## Linting

- Config: `.golangci.yml` (v2 format)
- Tests are excluded from linting
- Exclusion presets: comments, common-false-positives, legacy, std-error-handling
