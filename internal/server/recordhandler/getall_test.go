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

	"github.com/ashep/ujds/internal/model"
	"github.com/ashep/ujds/internal/server/recordhandler"
	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func TestRecordHandler_GetAll(tt *testing.T) {
	tt.Run("RecordRepoInvalidArgumentError", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.GetAllFunc = func(ctx context.Context, index string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error) {
			return nil, 0, apperrors.InvalidArgError{
				Subj:   "theRecordRepoSubj",
				Reason: "theRecordRepoReason",
			}
		}

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.GetAll(context.Background(), connect.NewRequest(&proto.GetAllRequest{}))

		assert.EqualError(t, err, "invalid_argument: invalid theRecordRepoSubj: theRecordRepoReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RecordRepoInternalError", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.GetAllFunc = func(ctx context.Context, index string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error) {
			return nil, 0, errors.New("theRecordRepoError")
		}

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.GetAll(context.Background(), connect.NewRequest(&proto.GetAllRequest{}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRecordRepoError","proc":"","err_code":123456789,"message":"record repo get all failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.GetAllFunc = func(ctx context.Context, index string, since time.Time, cursor uint64, limit uint32) ([]model.Record, uint64, error) {
			assert.Equal(t, "theIndexName", index)
			assert.Equal(t, time.Unix(0, 0), since)
			assert.Equal(t, uint64(0), cursor)
			assert.Equal(t, uint32(500), limit)

			return []model.Record{
				{
					ID:        "theRecordID1",
					IndexID:   11,
					Rev:       123,
					Data:      `{"foo1":"bar1"}`,
					CreatedAt: time.Unix(111, 0),
					UpdatedAt: time.Unix(222, 0),
				},
				{
					ID:        "theRecordID2",
					IndexID:   22,
					Rev:       234,
					Data:      `{"foo2":"bar2"}`,
					CreatedAt: time.Unix(333, 0),
					UpdatedAt: time.Unix(444, 0),
				},
			}, 345, nil
		}

		h := recordhandler.New(ir, rr, now, l)
		res, err := h.GetAll(context.Background(), connect.NewRequest(&proto.GetAllRequest{
			Index: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Empty(t, lb.String())

		require.Len(t, res.Msg.Records, 2)
		assert.Equal(t, uint64(345), res.Msg.Cursor)

		assert.Equal(t, "theRecordID1", res.Msg.Records[0].Id)
		assert.Equal(t, "theIndexName", res.Msg.Records[0].Index)
		assert.Equal(t, uint64(123), res.Msg.Records[0].Rev)
		assert.Equal(t, `{"foo1":"bar1"}`, res.Msg.Records[0].Data)
		assert.Equal(t, time.Unix(111, 0).Unix(), res.Msg.Records[0].CreatedAt)
		assert.Equal(t, time.Unix(222, 0).Unix(), res.Msg.Records[0].UpdatedAt)

		assert.Equal(t, "theRecordID2", res.Msg.Records[1].Id)
		assert.Equal(t, "theIndexName", res.Msg.Records[1].Index)
		assert.Equal(t, uint64(234), res.Msg.Records[1].Rev)
		assert.Equal(t, `{"foo2":"bar2"}`, res.Msg.Records[1].Data)
		assert.Equal(t, time.Unix(333, 0).Unix(), res.Msg.Records[1].CreatedAt)
		assert.Equal(t, time.Unix(444, 0).Unix(), res.Msg.Records[1].UpdatedAt)
	})
}
