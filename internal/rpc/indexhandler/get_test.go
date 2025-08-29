package indexhandler_test

import (
	"context"
	"database/sql"
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
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/model"
	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestIndexHandler_Get(tt *testing.T) {
	tt.Run("RepoInvalidArgError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, apperrors.InvalidArgError{Subj: "theSubj", Reason: "theReason"})

		h := indexhandler.New(rm, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid theSubj: theReason")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoNotFoundError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, apperrors.NotFoundError{Subj: "theSubj"})

		h := indexhandler.New(rm, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "not_found: theSubj is not found")
		assert.Empty(t, lb.String())
	})

	tt.Run("RepoInternalError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Get", mock.Anything, mock.Anything).
			Return(model.Index{}, errors.New("theRepoError"))

		h := indexhandler.New(rm, now, l)
		_, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoError","proc":"","err_code":123456789,"message":"index repo get failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("Get", mock.Anything, "theIndexName").
			Return(model.Index{
				ID:        123,
				Name:      "theIndexName",
				Title:     sql.NullString{String: "theIndexTitle", Valid: true},
				Schema:    []byte(`{"foo":"bar"}`),
				CreatedAt: time.Unix(123, 0),
				UpdatedAt: time.Unix(234, 0),
			}, nil)

		h := indexhandler.New(rm, now, l)
		res, err := h.Get(context.Background(), connect.NewRequest(&proto.GetRequest{
			Name: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Empty(t, lb.String())

		assert.Equal(t, "theIndexName", res.Msg.Name)
		assert.Equal(t, "theIndexTitle", res.Msg.Title)
		assert.Equal(t, uint64(time.Unix(123, 0).Unix()), res.Msg.CreatedAt)
		assert.Equal(t, uint64(time.Unix(234, 0).Unix()), res.Msg.UpdatedAt)
		assert.Equal(t, `{"foo":"bar"}`, res.Msg.Schema)
	})
}
