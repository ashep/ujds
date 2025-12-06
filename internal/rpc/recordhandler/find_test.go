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

func TestRecordHandler_GetAll(tt *testing.T) {
	//nolint:dupl // this is a test
	tt.Run("RecordRepoInvalidArgumentError", func(t *testing.T) {
		ir := &indexRepoMock{}

		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr := &recordRepoMock{}
		rr.On("Find", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return([]recordrepo.Record(nil), uint64(0), apperrors.InvalidArgError{
				Subj:   "theRecordRepoSubj",
				Reason: "theRecordRepoReason",
			})

		idxNameValidator := &stringValidatorMock{}
		recIDValidator := &stringValidatorMock{}
		recDataValidator := &stringValidatorMock{}

		h := recordhandler.New(ir, rr, idxNameValidator, recIDValidator, recDataValidator, now, l)
		_, err := h.Find(context.Background(), connect.NewRequest(&proto.FindRequest{}))

		assert.EqualError(t, err, "invalid_argument: invalid theRecordRepoSubj: theRecordRepoReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RecordRepoInternalError", func(t *testing.T) {
		ir := &indexRepoMock{}
		rr := &recordRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr.On("Find", mock.Anything, mock.Anything).
			Return([]recordrepo.Record(nil), uint64(0), errors.New("theRecordRepoError"))

		idxNameValidator := &stringValidatorMock{}
		recIDValidator := &stringValidatorMock{}
		recDataValidator := &stringValidatorMock{}

		h := recordhandler.New(ir, rr, idxNameValidator, recIDValidator, recDataValidator, now, l)
		_, err := h.Find(context.Background(), connect.NewRequest(&proto.FindRequest{}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRecordRepoError","proc":"","err_code":123456789,"message":"record repo find failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rr := &recordRepoMock{}
		rr.On("Find", mock.Anything, mock.Anything).
			Return([]recordrepo.Record{
				{
					ID:        "theRecordID1",
					IndexID:   11,
					Rev:       123,
					Data:      `{"foo1":"bar1"}`,
					CreatedAt: time.Unix(111, 0),
					UpdatedAt: time.Unix(112, 0),
					TouchedAt: time.Unix(113, 0),
				},
				{
					ID:        "theRecordID2",
					IndexID:   22,
					Rev:       234,
					Data:      `{"foo2":"bar2"}`,
					CreatedAt: time.Unix(221, 0),
					UpdatedAt: time.Unix(222, 0),
					TouchedAt: time.Unix(223, 0),
				},
			}, uint64(345), nil)

		idxNameValidator := &stringValidatorMock{}
		recIDValidator := &stringValidatorMock{}
		recDataValidator := &stringValidatorMock{}

		h := recordhandler.New(ir, rr, idxNameValidator, recIDValidator, recDataValidator, now, l)
		res, err := h.Find(context.Background(), connect.NewRequest(&proto.FindRequest{
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
		assert.Equal(t, time.Unix(112, 0).Unix(), res.Msg.Records[0].UpdatedAt)
		assert.Equal(t, time.Unix(113, 0).Unix(), res.Msg.Records[0].TouchedAt)

		assert.Equal(t, "theRecordID2", res.Msg.Records[1].Id)
		assert.Equal(t, "theIndexName", res.Msg.Records[1].Index)
		assert.Equal(t, uint64(234), res.Msg.Records[1].Rev)
		assert.Equal(t, `{"foo2":"bar2"}`, res.Msg.Records[1].Data)
		assert.Equal(t, time.Unix(221, 0).Unix(), res.Msg.Records[1].CreatedAt)
		assert.Equal(t, time.Unix(222, 0).Unix(), res.Msg.Records[1].UpdatedAt)
		assert.Equal(t, time.Unix(223, 0).Unix(), res.Msg.Records[1].TouchedAt)
	})
}
