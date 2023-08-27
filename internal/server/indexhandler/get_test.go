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
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/server/indexhandler"

	"github.com/ashep/ujds/internal/model"
	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestHandler_Get(tt *testing.T) {
	tt.Parallel()

	tt.Run("RepoInvalidArgError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, apperrors.InvalidArgError{Subj: "theSubj", Reason: "theReason"}
		}

		h := indexhandler.New(ir, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid theSubj: theReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoNotFoundError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, apperrors.NotFoundError{Subj: "theSubj"}
		}

		h := indexhandler.New(ir, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "not_found: theSubj is not found")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoInternalError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, errors.New("theRepoError")
		}

		h := indexhandler.New(ir, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoError","proc":"","err_code":123456789,"message":"index repo get failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{
				ID:        123,
				Name:      "theIndexName",
				Schema:    []byte(`{"foo":"bar"}`),
				CreatedAt: time.Unix(123, 0),
				UpdatedAt: time.Unix(234, 0),
			}, nil
		}

		h := indexhandler.New(ir, now, l)
		res, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Empty(t, lb.String())

		assert.Equal(t, "theIndexName", res.Msg.Name)
		assert.Equal(t, uint64(time.Unix(123, 0).Unix()), res.Msg.CreatedAt)
		assert.Equal(t, uint64(time.Unix(234, 0).Unix()), res.Msg.UpdatedAt)
		assert.Equal(t, `{"foo":"bar"}`, res.Msg.Schema)
	})
}
