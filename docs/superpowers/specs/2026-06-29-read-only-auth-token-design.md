# Read-Only Auth Token

**Date:** 2026-06-29
**Status:** Approved

## Summary

Add a second optional auth token (`auth_token_ro`) that grants access only to read-only RPC methods. The existing `auth_token` continues to grant full access. When a read-only token is used against a write method, the server returns `PermissionDenied`.

## Scope

- Single additional config value — no multi-token/role maps.
- Server side only; no SDK changes required.

## Read-Only vs Write Methods

| Service | Method | Classification |
|---------|--------|---------------|
| IndexService | Push | write |
| IndexService | Get | read-only |
| IndexService | List | read-only |
| IndexService | Clear | write |
| RecordService | Push | write |
| RecordService | Get | read-only |
| RecordService | Find | read-only |
| RecordService | History | read-only |

## Config Changes

`internal/app/config.go` — add to `Server`:

```go
AuthTokenRO string `json:"auth_token_ro" yaml:"auth_token_ro"`
```

Env var follows the same convention as `auth_token` → `SERVER_AUTH_TOKEN_RO`.

## Interceptor Changes

`internal/app/app.go` — replace `auth(token string)` with `auth(token, tokenRO string)`.

A package-level set holds the read-only procedure paths:

```
/ujds.index.v1.IndexService/Get
/ujds.index.v1.IndexService/List
/ujds.record.v1.RecordService/Get
/ujds.record.v1.RecordService/Find
/ujds.record.v1.RecordService/History
```

Decision logic per request (presented = `Authorization` header with `Bearer ` stripped):

1. `token == ""` → auth disabled; allow all (existing behavior unchanged).
2. presented matches `token` → allow (full access).
3. `tokenRO != ""` AND presented matches `tokenRO` AND procedure is in read-only set → allow.
4. `tokenRO != ""` AND presented matches `tokenRO` AND procedure is NOT in read-only set → `CodePermissionDenied`.
5. All other cases → `CodeUnauthenticated`.

## Error Codes

| Situation | gRPC code |
|-----------|-----------|
| Token absent or unknown | `Unauthenticated` |
| Read-only token on a write method | `PermissionDenied` |

## Edge Cases

- If `auth_token` is empty (auth disabled), `auth_token_ro` has no effect — server remains fully open.
- `auth_token_ro` may be omitted; the server behaves exactly as today.

## Tests

Unit tests for `auth()` in `internal/app/app_test.go` (new file) covering:

- Full token on a read method → allowed.
- Full token on a write method → allowed.
- Read-only token on a read method → allowed.
- Read-only token on a write method → `PermissionDenied`.
- Unknown token → `Unauthenticated`.
- Empty token (auth disabled) → allowed regardless of method.

## Documentation

Update `README.md` config reference to document `auth_token_ro` / `SERVER_AUTH_TOKEN_RO`.
