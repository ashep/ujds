package recordhandler_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"
	"github.com/ashep/ujds/internal/rpc/recordhandler"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/model"
	proto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
)

func TestRecordHandler_Push(tt *testing.T) {
	tt.Run("EmptyRecords", func(t *testing.T) {
		now := func() time.Time { return time.Unix(1234567890, 987654321) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir := &indexRepoMock{}
		defer ir.AssertExpectations(t)

		rr := &recordRepoMock{}
		defer rr.AssertExpectations(t)

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{}))

		assert.EqualError(t, err, "invalid_argument: empty records")
		assert.Empty(t, lb.String())
	})

	tt.Run("IndexRepoArgumentError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(1234567890, 987654321) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir := &indexRepoMock{}
		defer ir.AssertExpectations(t)
		ir.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, apperrors.InvalidArgError{
				Subj:   "theIndexRepoSubj",
				Reason: "theIndexRepoReason",
			})

		rr := &recordRepoMock{}
		defer rr.AssertExpectations(t)

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Records: []*proto.PushRequest_Record{{Index: "anIndex", Id: "anID", Data: "aData"}},
		}))

		assert.EqualError(t, err, "invalid_argument: invalid theIndexRepoSubj: theIndexRepoReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("IndexRepoNotFoundError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(1234567890, 987654321) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir := &indexRepoMock{}
		defer ir.AssertExpectations(t)
		ir.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, apperrors.NotFoundError{
				Subj: "theIndexRepoSubj",
			})

		rr := &recordRepoMock{}
		defer rr.AssertExpectations(t)

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Records: []*proto.PushRequest_Record{{Index: "anIndex", Id: "anID", Data: "aData"}},
		}))

		assert.EqualError(t, err, "not_found: theIndexRepoSubj is not found")
		assert.Empty(t, lb.String())
	})

	tt.Run("IndexRepoOtherError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(1234567890, 987654321) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir := &indexRepoMock{}
		defer ir.AssertExpectations(t)
		ir.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, errors.New("theIndexRepoOtherError"))

		rr := &recordRepoMock{}
		defer rr.AssertExpectations(t)

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Records: []*proto.PushRequest_Record{{Index: "anIndex", Id: "anID", Data: "aData"}},
		}))

		assert.EqualError(t, err, "internal: err_code: 1234567890987")
		assert.Equal(t, `{"level":"error","error":"theIndexRepoOtherError","proc":"","err_code":1234567890987,"message":"index repo get failed"}
`, lb.String())
	})

	tt.Run("RecordRepoInvalidArgError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(1234567890, 987654321) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir := &indexRepoMock{}
		defer ir.AssertExpectations(t)
		ir.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, nil)

		rr := &recordRepoMock{}
		defer rr.AssertExpectations(t)
		rr.On("Push", mock.Anything, mock.Anything).
			Return(apperrors.InvalidArgError{
				Subj:   "theErrorSubj",
				Reason: "theErrorReason",
			})

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{Records: []*proto.PushRequest_Record{
			{
				Index: "anIndex",
				Id:    "anID",
				Data:  "aData",
			},
		}}))

		assert.EqualError(t, err, "invalid_argument: invalid theErrorSubj: theErrorReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RecordRepoInternalError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(1234567890, 987654321) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir := &indexRepoMock{}
		defer ir.AssertExpectations(t)
		ir.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, nil)

		rr := &recordRepoMock{}
		defer rr.AssertExpectations(t)
		rr.On("Push", mock.Anything, mock.Anything).
			Return(errors.New("theRecordRepoError"))

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{Records: []*proto.PushRequest_Record{
			{
				Index: "anIndex",
				Id:    "anID",
				Data:  "aData",
			},
		}}))

		assert.EqualError(t, err, "internal: err_code: 1234567890987")
		assert.Equal(t, `{"level":"error","error":"theRecordRepoError","proc":"","err_code":1234567890987,"message":"record repo push failed"}
`, lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		now := func() time.Time { return time.Unix(1234567890, 987654321) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir := &indexRepoMock{}
		defer ir.AssertExpectations(t)
		ir.On("Get", mock.Anything, "theIndex").
			Return(model.Index{
				ID:     123,
				Schema: []byte("theIndexSchema"),
			}, nil)

		rr := &recordRepoMock{}
		defer rr.AssertExpectations(t)
		rr.On("Push", mock.Anything, []model.RecordUpdate{
			{
				ID:      "theRecordID",
				IndexID: 123,
				Schema:  []byte("theIndexSchema"),
				Data:    "theRecordData",
			},
		}).
			Return(nil)

		h := recordhandler.New(ir, rr, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{Records: []*proto.PushRequest_Record{
			{
				Index: "theIndex",
				Id:    "theRecordID",
				Data:  "theRecordData",
			},
		}}))

		require.NoError(t, err)
		assert.Empty(t, lb.String())
	})
}
