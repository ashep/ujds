package indexhandler_test

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/ujds/internal/indexrepo"
	"github.com/ashep/ujds/internal/rpc/indexhandler"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestIndexHandler_List(tt *testing.T) {
	tt.Run("RepoError", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("List", mock.Anything).
			Return([]indexrepo.Index(nil), errors.New("theRepoListError"))

		h := indexhandler.New(rm, now, l)
		_, err := h.List(context.Background(), connect.NewRequest(&proto.ListRequest{}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoListError","proc":"","err_code":123456789,"message":"index repo list failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("List", mock.Anything).
			Return([]indexrepo.Index{
				{
					ID:        123,
					Name:      "theIndex1",
					Title:     sql.NullString{String: "theTitle1", Valid: true},
					CreatedAt: time.Unix(234, 0),
					UpdatedAt: time.Unix(345, 0),
				},
				{
					ID:        321,
					Name:      "theIndex2",
					Title:     sql.NullString{String: "theTitle2", Valid: true},
					CreatedAt: time.Unix(432, 0),
					UpdatedAt: time.Unix(543, 0),
				},
			}, nil)

		h := indexhandler.New(rm, now, l)
		res, err := h.List(context.Background(), connect.NewRequest(&proto.ListRequest{}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Indices, 2)
		assert.Empty(t, lb.String())

		assert.Equal(t, "theIndex1", res.Msg.Indices[0].Name)
		assert.Equal(t, "theTitle1", res.Msg.Indices[0].Title)

		assert.Equal(t, "theIndex2", res.Msg.Indices[1].Name)
		assert.Equal(t, "theTitle2", res.Msg.Indices[1].Title)
	})

	tt.Run("OkWithFilter", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		rm := &repoMock{}
		defer rm.AssertExpectations(t)
		rm.On("List", mock.Anything).
			Return([]indexrepo.Index{
				{
					ID:        123,
					Name:      "theIndex1Foo",
					Title:     sql.NullString{String: "theTitle1", Valid: true},
					CreatedAt: time.Unix(234, 0),
					UpdatedAt: time.Unix(345, 0),
				},
				{
					ID:        321,
					Name:      "theIndex2Bar",
					Title:     sql.NullString{String: "theTitle2", Valid: true},
					CreatedAt: time.Unix(432, 0),
					UpdatedAt: time.Unix(543, 0),
				},
			}, nil)

		h := indexhandler.New(rm, now, l)
		res, err := h.List(context.Background(), connect.NewRequest(&proto.ListRequest{
			Filter: &proto.ListRequestFilter{
				Names: []string{"theIndex2*"},
			},
		}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Indices, 1)
		assert.Empty(t, lb.String())

		assert.Equal(t, "theIndex2Bar", res.Msg.Indices[0].Name)
		assert.Equal(t, "theTitle2", res.Msg.Indices[0].Title)
	})
}
