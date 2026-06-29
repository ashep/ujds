# Read-Only Auth Token Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add an optional `auth_token_ro` config value that, when set, allows read-only RPC methods to be called using it while write methods still require the full `auth_token`.

**Architecture:** Extract auth decision logic from the `auth()` interceptor into a pure `checkAuth()` function (testable without mocking `connect.AnyRequest`). The interceptor delegates to `checkAuth()`, passing the presented token and procedure path. A package-level map holds the read-only procedure paths, populated with constants from the generated connect packages.

**Tech Stack:** Go, `connectrpc.com/connect`, `kelseyhightower/envconfig` (via `go-app/cfgloader`).

## Global Constraints

- No new dependencies.
- Keep all logic in `internal/app` package — no new packages.
- Follow existing patterns: bare `connect.NewError` calls, no wrapping.
- `auth_token == ""` means auth disabled; `auth_token_ro` has no effect in that case.

---

### Task 1: Config — add `AuthTokenRO` field

**Files:**
- Modify: `internal/app/config.go`

**Interfaces:**
- Produces: `Config.Server.AuthTokenRO string` (used in Task 2)

- [ ] **Step 1: Add `AuthTokenRO` to `Server` struct**

In `internal/app/config.go`, change the `Server` struct from:

```go
type Server struct {
	Addr      string `json:"addr" yaml:"addr"`
	AuthToken string `json:"auth_token" yaml:"auth_token"`
}
```

to:

```go
type Server struct {
	Addr        string `json:"addr" yaml:"addr"`
	AuthToken   string `json:"auth_token" yaml:"auth_token"`
	AuthTokenRO string `json:"auth_token_ro" yaml:"auth_token_ro"`
}
```

- [ ] **Step 2: Verify the project still compiles**

```bash
go build ./...
```

Expected: no output, exit 0.

- [ ] **Step 3: Commit**

```bash
git add internal/app/config.go
git commit -m "Add AuthTokenRO field to Server config"
```

---

### Task 2: Auth logic — `checkAuth()` function and updated interceptor

**Files:**
- Modify: `internal/app/app.go`
- Create: `internal/app/app_test.go`

**Interfaces:**
- Consumes: `Config.Server.AuthTokenRO string` from Task 1
- Produces: `checkAuth(token, tokenRO, presented, procedure string) error` (pure function, tested in this task)

- [ ] **Step 1: Write failing tests**

Create `internal/app/app_test.go`:

```go
package app

import (
	"testing"

	"connectrpc.com/connect"
)

func TestCheckAuth(t *testing.T) {
	const (
		full      = "full-token"
		ro        = "ro-token"
		readProc  = "/ujds.index.v1.IndexService/Get"
		writeProc = "/ujds.index.v1.IndexService/Push"
	)

	tests := []struct {
		name      string
		token     string
		tokenRO   string
		presented string
		procedure string
		wantCode  connect.Code
		wantErr   bool
	}{
		{
			name:      "auth disabled allows anything",
			token:     "",
			tokenRO:   "",
			presented: "",
			procedure: writeProc,
			wantErr:   false,
		},
		{
			name:      "full token on read method allowed",
			token:     full,
			tokenRO:   ro,
			presented: full,
			procedure: readProc,
			wantErr:   false,
		},
		{
			name:      "full token on write method allowed",
			token:     full,
			tokenRO:   ro,
			presented: full,
			procedure: writeProc,
			wantErr:   false,
		},
		{
			name:      "ro token on read method allowed",
			token:     full,
			tokenRO:   ro,
			presented: ro,
			procedure: readProc,
			wantErr:   false,
		},
		{
			name:      "ro token on write method returns PermissionDenied",
			token:     full,
			tokenRO:   ro,
			presented: ro,
			procedure: writeProc,
			wantErr:   true,
			wantCode:  connect.CodePermissionDenied,
		},
		{
			name:      "unknown token returns Unauthenticated",
			token:     full,
			tokenRO:   ro,
			presented: "wrong",
			procedure: readProc,
			wantErr:   true,
			wantCode:  connect.CodeUnauthenticated,
		},
		{
			name:      "empty presented token returns Unauthenticated",
			token:     full,
			tokenRO:   ro,
			presented: "",
			procedure: readProc,
			wantErr:   true,
			wantCode:  connect.CodeUnauthenticated,
		},
		{
			name:      "ro token not configured falls through to Unauthenticated",
			token:     full,
			tokenRO:   "",
			presented: ro,
			procedure: readProc,
			wantErr:   true,
			wantCode:  connect.CodeUnauthenticated,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := checkAuth(tc.token, tc.tokenRO, tc.presented, tc.procedure)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				cerr, ok := err.(*connect.Error)
				if !ok {
					t.Fatalf("expected *connect.Error, got %T", err)
				}
				if cerr.Code() != tc.wantCode {
					t.Fatalf("expected code %v, got %v", tc.wantCode, cerr.Code())
				}
			} else if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./internal/app/...
```

Expected: FAIL with `undefined: checkAuth`.

- [ ] **Step 3: Implement `readOnlyProcedures`, `checkAuth()`, and update `auth()`**

In `internal/app/app.go`:

Add imports (merge with existing):
```go
import (
    indexconnect "github.com/ashep/ujds/sdk/proto/ujds/index/v1/v1connect"
    recordconnect "github.com/ashep/ujds/sdk/proto/ujds/record/v1/v1connect"
)
```

Add after the existing imports, before `Run()`:

```go
var readOnlyProcedures = map[string]bool{
	indexconnect.IndexServiceGetProcedure:      true,
	indexconnect.IndexServiceListProcedure:     true,
	recordconnect.RecordServiceGetProcedure:    true,
	recordconnect.RecordServiceFindProcedure:   true,
	recordconnect.RecordServiceHistoryProcedure: true,
}

func checkAuth(token, tokenRO, presented, procedure string) error {
	if token == "" {
		return nil
	}
	if presented == token {
		return nil
	}
	if tokenRO != "" && presented == tokenRO {
		if readOnlyProcedures[procedure] {
			return nil
		}
		return connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}
	return connect.NewError(connect.CodeUnauthenticated, errors.New("not authorized"))
}
```

Replace the existing `auth()` function:

```go
func auth(token, tokenRO string) connect.UnaryInterceptorFunc {
	return func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			presented := strings.ReplaceAll(req.Header().Get("Authorization"), "Bearer ", "")
			if err := checkAuth(token, tokenRO, presented, req.Spec().Procedure); err != nil {
				return nil, err
			}
			return next(ctx, req)
		}
	}
}
```

In `Run()`, update the `icps` line from:

```go
icps := connect.WithInterceptors(auth(rt.Cfg.Server.AuthToken))
```

to:

```go
icps := connect.WithInterceptors(auth(rt.Cfg.Server.AuthToken, rt.Cfg.Server.AuthTokenRO))
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./internal/app/...
```

Expected: `ok  github.com/ashep/ujds/internal/app`

- [ ] **Step 5: Verify the full build still compiles**

```bash
go build ./...
```

Expected: no output, exit 0.

- [ ] **Step 6: Commit**

```bash
git add internal/app/app.go internal/app/app_test.go
git commit -m "Add read-only auth token support"
```

---

### Task 3: README documentation

**Files:**
- Modify: `README.md`

**Interfaces:**
- Consumes: nothing (doc-only task)

- [ ] **Step 1: Update the config file section**

In `README.md`, find the `server` config block:

```
- *optional* **object** `server`: server configuration.
    - *optional* **string** `address`: network address, default is `:9000`.
    - *optional* **string** `auth_token`: authorization token.
```

Replace with:

```
- *optional* **object** `server`: server configuration.
    - *optional* **string** `address`: network address, default is `:9000`.
    - *optional* **string** `auth_token`: authorization token (full access).
    - *optional* **string** `auth_token_ro`: read-only authorization token; grants access to read-only methods only (`IndexService/Get`, `IndexService/List`, `RecordService/Get`, `RecordService/Find`, `RecordService/History`).
```

- [ ] **Step 2: Update the env variables section**

Find:

```
- *optional* **string** `UJDS_SERVER_AUTHTOKEN`: server authorization token.
```

Replace with:

```
- *optional* **string** `UJDS_SERVER_AUTHTOKEN`: server authorization token (full access).
- *optional* **string** `UJDS_SERVER_AUTHTOKENRO`: read-only server authorization token.
```

- [ ] **Step 3: Update the Authorization section**

Find:

```
If the `server.auth_token` configuration parameter is specified, the server will expect an `Authorization: Bearer XXX`
HTTP header, where `XXX` must match the configured token value, otherwise the `403` HTTP status will be returned.
```

Replace with:

```
If the `server.auth_token` configuration parameter is specified, the server will expect an `Authorization: Bearer XXX`
HTTP header where `XXX` matches the configured token, otherwise `401` is returned. If `server.auth_token_ro` is also
configured, that token may be used instead for read-only methods (`IndexService/Get`, `IndexService/List`,
`RecordService/Get`, `RecordService/Find`, `RecordService/History`). Using the read-only token on a write method
returns `403`.
```

- [ ] **Step 4: Verify build and tests still pass**

```bash
go build ./... && go test ./internal/app/...
```

Expected: no errors.

- [ ] **Step 5: Commit**

```bash
git add README.md
git commit -m "Document auth_token_ro config and env var"
```
