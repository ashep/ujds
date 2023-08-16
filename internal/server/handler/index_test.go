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

	"github.com/ashep/ujds/internal/model"
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
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoOtherError", func(t *testing.T) {
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

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.UpsertFunc = func(ctx context.Context, name string, schema string) error {
			assert.Equal(t, "theIndexName", name)
			assert.Equal(t, `{"foo":"bar"}`, schema)
			return nil
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: `{"foo":"bar"}`,
		}))

		assert.NoError(t, err)
	})
}

func TestHandler_GetIndex(tt *testing.T) {
	tt.Parallel()

	tt.Run("RepoInvalidArgError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, apperrors.InvalidArgError{Subj: "theSubj", Reason: "theReason"}
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.GetIndex(context.Background(), connect.NewRequest(&v1.GetIndexRequest{
			Name: "theIndexName",
		}))

		require.EqualError(t, err, "invalid_argument: invalid theSubj: theReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoNotFoundError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, apperrors.NotFoundError{Subj: "theSubj"}
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.GetIndex(context.Background(), connect.NewRequest(&v1.GetIndexRequest{
			Name: "theIndexName",
		}))

		require.EqualError(t, err, "not_found: theSubj is not found")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoOtherError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, errors.New("theRepoError")
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.GetIndex(context.Background(), connect.NewRequest(&v1.GetIndexRequest{
			Name: "theIndexName",
		}))

		require.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoError","proc":"","err_code":123456789,"message":"index repo get failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
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

		h := handler.New(ir, rr, now, l)
		res, err := h.GetIndex(context.Background(), connect.NewRequest(&v1.GetIndexRequest{
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
