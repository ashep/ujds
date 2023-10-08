package recordhandler_test

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

	"github.com/ashep/ujds/internal/server/recordhandler"

	"github.com/ashep/ujds/internal/model"
	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func TestRecordHandler_Get(tt *testing.T) {
	tt.Run("RecordRepoInvalidArgumentError", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.GetFunc = func(ctx context.Context, index, id string) (model.Record, error) {
			return model.Record{}, apperrors.InvalidArgError{
				Subj:   "theRecordRepoSubj",
				Reason: "theRecordRepoReason",
			}
		}

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{}))

		assert.EqualError(t, err, "invalid_argument: invalid theRecordRepoSubj: theRecordRepoReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RecordRepoNotFoundError", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.GetFunc = func(ctx context.Context, index, id string) (model.Record, error) {
			return model.Record{}, apperrors.NotFoundError{
				Subj: "theRecordRepoSubj",
			}
		}

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{}))

		assert.EqualError(t, err, "not_found: theRecordRepoSubj is not found")
		assert.Empty(t, lb.String())
	})

	tt.Run("RecordRepoInternalError", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.GetFunc = func(ctx context.Context, index, id string) (model.Record, error) {
			return model.Record{}, errors.New("theRecordRepoError")
		}

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRecordRepoError","proc":"","err_code":123456789,"message":"record repo push failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.GetFunc = func(ctx context.Context, index, id string) (model.Record, error) {
			assert.Equal(t, "theIndexName", index)
			assert.Equal(t, "theRecordID", id)

			return model.Record{
				ID:        id,
				IndexID:   123,
				Rev:       234,
				Data:      "theData",
				CreatedAt: time.Unix(345, 0),
				UpdatedAt: time.Unix(456, 0),
			}, nil
		}

		h := recordhandler.New(ir, rr, now, l)
		res, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Index: "theIndexName",
			Id:    "theRecordID",
		}))

		require.NoError(t, err)
		assert.Equal(t, "theRecordID", res.Msg.Record.Id)
		assert.Equal(t, "theIndexName", res.Msg.Record.Index)
		assert.Equal(t, uint64(234), res.Msg.Record.Rev)
		assert.Equal(t, "theData", res.Msg.Record.Data)
		assert.Equal(t, time.Unix(345, 0).Unix(), res.Msg.Record.CreatedAt)
		assert.Equal(t, time.Unix(456, 0).Unix(), res.Msg.Record.UpdatedAt)
		assert.Empty(t, lb.String())
	})
}
