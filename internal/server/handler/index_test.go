package handler_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/ashep/go-apperrors"
	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/server/handler"
	v1 "github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func TestHandler_PushIndex(tt *testing.T) {
	tt.Parallel()

	tt.Run("RepoInvalidArgError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.UpsertFunc = func(ctx context.Context, name string, schema string) error {
			return apperrors.InvalidArgError{Subj: "theSubj", Reason: "theReason"}
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: "{}",
		}))

		require.EqualError(t, err, "invalid_argument: invalid theSubj: theReason")
		assert.Equal(t, "", lb.String())
	})

	tt.Run("RepoError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.UpsertFunc = func(ctx context.Context, name string, schema string) error {
			return errors.New("theRepoError")
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: "{}",
		}))

		require.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoError","proc":"","err_code":123456789,"message":"index repo upsert failed"}`+"\n", lb.String())
	})
}
