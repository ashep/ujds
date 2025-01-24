package indexhandler_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/server/indexhandler"
	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestIndexHandler_Clear(tt *testing.T) {
	tt.Run("RepoInvalidArgError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Clear", mock.Anything, mock.Anything).
			Return(apperrors.InvalidArgError{Subj: "theSubj", Reason: "theReason"})

		h := indexhandler.New(rm, now, l)
		_, err := h.Clear(context.Background(), connect.NewRequest(&proto.ClearRequest{
			Name: "theIndexName",
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
		rm.On("Clear", mock.Anything, mock.Anything).
			Return(errors.New("theRepoError"))

		h := indexhandler.New(rm, now, l)
		_, err := h.Clear(context.Background(), connect.NewRequest(&proto.ClearRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoError","proc":"","err_code":123456789,"message":"index repo clear failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Clear", mock.Anything, mock.Anything).
			Return(nil)

		h := indexhandler.New(rm, now, l)
		_, err := h.Clear(context.Background(), connect.NewRequest(&proto.ClearRequest{
			Name: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Empty(t, lb.String())
	})
}
