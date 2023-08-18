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
	ujdsproto "github.com/ashep/ujds/sdk/proto/ujds/v1"
)

func TestHandler_PushRecords(tt *testing.T) {
	tt.Parallel()

	tt.Run("IndexRepoInvalidArgumentError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, apperrors.InvalidArgError{
				Subj:   "theIndexRepoSubj",
				Reason: "theIndexRepoReason",
			}
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		require.EqualError(t, err, "invalid_argument: invalid theIndexRepoSubj: theIndexRepoReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("IndexRepoNotFoundError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, apperrors.NotFoundError{
				Subj: "theIndexRepoSubj",
			}
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		require.EqualError(t, err, "not_found: theIndexRepoSubj is not found")
		assert.Empty(t, lb.String())
	})

	tt.Run("IndexRepoOtherError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, errors.New("theIndexRepoError")
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		require.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theIndexRepoError","proc":"","err_code":123456789,"message":"index repo get failed"}`+"\n", lb.String())
	})

	tt.Run("RecordRepoInvalidArgError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, nil
		}

		rr.PushFunc = func(ctx context.Context, index model.Index, records []model.Record) error {
			return apperrors.InvalidArgError{
				Subj:   "theErrorSubj",
				Reason: "theErrorReason",
			}
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		require.EqualError(t, err, "invalid_argument: invalid theErrorSubj: theErrorReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RecordRepoOtherError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, nil
		}

		rr.PushFunc = func(ctx context.Context, index model.Index, records []model.Record) error {
			return errors.New("theRecordRepoError")
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		require.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRecordRepoError","proc":"","err_code":123456789,"message":"record repo push failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.GetFunc = func(ctx context.Context, name string) (model.Index, error) {
			return model.Index{}, nil
		}

		rr.PushFunc = func(ctx context.Context, index model.Index, records []model.Record) error {
			return nil
		}

		h := handler.New(ir, rr, now, l)
		_, err := h.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		require.NoError(t, err)
		assert.Empty(t, lb.String())
	})
}
