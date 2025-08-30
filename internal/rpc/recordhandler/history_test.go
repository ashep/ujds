package recordhandler_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"
	"github.com/ashep/ujds/internal/recordrepo"
	"github.com/ashep/ujds/internal/rpc/recordhandler"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func TestRecordHandler_History(tt *testing.T) {
	//nolint:dupl // this is test
	tt.Run("RecordRepoInvalidArgumentError", func(t *testing.T) {
		ir := &indexRepoMock{}

		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr := &recordRepoMock{}
		rr.On("History", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]recordrepo.Record(nil), uint64(0), apperrors.InvalidArgError{
				Subj:   "theRecordRepoSubj",
				Reason: "theRecordRepoReason",
			})

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.History(context.Background(), connect.NewRequest(&proto.HistoryRequest{}))

		assert.EqualError(t, err, "invalid_argument: invalid theRecordRepoSubj: theRecordRepoReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RecordRepoInternalError", func(t *testing.T) {
		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr := &recordRepoMock{}
		rr.On("History", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]recordrepo.Record(nil), uint64(0), errors.New("theRecordRepoInternalError"))

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.History(context.Background(), connect.NewRequest(&proto.HistoryRequest{}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRecordRepoInternalError","proc":"","err_code":123456789,"message":"record repo history failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr := &recordRepoMock{}
		rr.On("History", mock.Anything, "theIndexName", "theRecordID", time.Unix(55, 0), uint64(77), uint32(66)).
			Return([]recordrepo.Record{
				{
					ID:        "theRecord1",
					IndexID:   11,
					Rev:       111,
					Data:      "theData1",
					CreatedAt: time.Unix(1122, 0),
					UpdatedAt: time.Unix(1133, 0),
				},
				{
					ID:        "theRecord2",
					IndexID:   22,
					Rev:       222,
					Data:      "theData2",
					CreatedAt: time.Unix(2222, 0),
					UpdatedAt: time.Unix(2233, 0),
				},
			}, uint64(78), nil)

		h := recordhandler.New(ir, rr, now, l)
		res, err := h.History(context.Background(), connect.NewRequest(&proto.HistoryRequest{
			Index:  "theIndexName",
			Id:     "theRecordID",
			Since:  55,
			Limit:  66,
			Cursor: 77,
		}))

		require.NoError(t, err)
		assert.Empty(t, lb.String())
		require.Len(t, res.Msg.Records, 2)
		assert.Equal(t, uint64(78), res.Msg.Cursor)

		assert.Equal(t, "theRecord1", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(111), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndexName", res.Msg.Records[0].Index)
		assert.Equal(t, "theData1", res.Msg.Records[0].Data)
		assert.Equal(t, int64(1122), res.Msg.Records[0].CreatedAt)
		assert.Equal(t, int64(0), res.Msg.Records[0].UpdatedAt)
		assert.Equal(t, int64(0), res.Msg.Records[0].TouchedAt)

		assert.Equal(t, "theRecord2", res.Msg.Records[1].Id)
		assert.Equal(t, uint64(222), res.Msg.Records[1].Rev)
		assert.Equal(t, "theIndexName", res.Msg.Records[1].Index)
		assert.Equal(t, "theData2", res.Msg.Records[1].Data)
		assert.Equal(t, int64(2222), res.Msg.Records[1].CreatedAt)
		assert.Equal(t, int64(0), res.Msg.Records[1].UpdatedAt)
		assert.Equal(t, int64(0), res.Msg.Records[1].TouchedAt)
	})
}
