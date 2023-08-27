package indexhandler_test

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

	"github.com/ashep/ujds/internal/server/indexhandler"
	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestHandler_Push(tt *testing.T) {
	tt.Parallel()

	tt.Run("RepoInvalidArgError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.UpsertFunc = func(ctx context.Context, name string, schema string) error {
			return apperrors.InvalidArgError{Subj: "theSubj", Reason: "theReason"}
		}

		h := indexhandler.New(ir, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Schema: "{}",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid theSubj: theReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoInternalError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.UpsertFunc = func(ctx context.Context, name string, schema string) error {
			return errors.New("theRepoError")
		}

		h := indexhandler.New(ir, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Schema: "{}",
		}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoError","proc":"","err_code":123456789,"message":"index repo upsert failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.UpsertFunc = func(ctx context.Context, name string, schema string) error {
			assert.Equal(t, "theIndexName", name)
			assert.Equal(t, `{"foo":"bar"}`, schema)
			return nil
		}

		h := indexhandler.New(ir, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Schema: `{"foo":"bar"}`,
		}))

		assert.NoError(t, err)
	})
}
