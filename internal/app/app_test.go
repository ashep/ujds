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
