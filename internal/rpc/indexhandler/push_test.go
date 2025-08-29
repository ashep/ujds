package indexhandler_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"
	"github.com/ashep/ujds/internal/rpc/indexhandler"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestIndexHandler_Push(tt *testing.T) {
	tt.Run("RepoInvalidArgError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(apperrors.InvalidArgError{Subj: "theSubj", Reason: "theReason"})

		h := indexhandler.New(rm, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Schema: "{}",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid theSubj: theReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoInternalError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(errors.New("theRepoError"))

		h := indexhandler.New(rm, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Schema: "{}",
		}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoError","proc":"","err_code":123456789,"message":"index repo upsert failed"}`+"\n", lb.String())
	})

	tt.Run("IndexNotFound", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Upsert", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(apperrors.NotFoundError{Subj: "theNotFoundSubj"})

		h := indexhandler.New(rm, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Schema: "{}",
		}))

		assert.EqualError(t, err, "not_found: theNotFoundSubj is not found")
		assert.Empty(t, lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Upsert", mock.Anything, "theIndexName", "", `{"foo":"bar"}`).
			Return(nil)

		h := indexhandler.New(rm, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Title:  "",
			Schema: `{"foo":"bar"}`,
		}))

		assert.NoError(t, err)
		assert.Empty(t, lb.String())
	})

	tt.Run("OkWithTitle", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Upsert", mock.Anything, "theIndexName", "theIndexTitle", `{"foo":"bar"}`).
			Return(nil)

		h := indexhandler.New(rm, now, l)
		_, err := h.Push(context.Background(), connect.NewRequest(&proto.PushRequest{
			Name:   "theIndexName",
			Title:  "theIndexTitle",
			Schema: `{"foo":"bar"}`,
		}))

		assert.NoError(t, err)
		assert.Empty(t, lb.String())
	})
}
