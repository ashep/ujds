package indexhandler_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/model"
	"github.com/ashep/ujds/internal/server/indexhandler"
	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestIndexHandler_List(tt *testing.T) {
	tt.Parallel()

	tt.Run("RepoError", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.ListFunc = func(ctx context.Context) ([]model.Index, error) {
			return nil, errors.New("theRepoListError")
		}

		h := indexhandler.New(ir, now, l)
		_, err := h.List(context.Background(), connect.NewRequest(&proto.ListRequest{}))

		assert.EqualError(t, err, "internal: err_code: 123456789")
		assert.Equal(t, `{"level":"error","error":"theRepoListError","proc":"","err_code":123456789,"message":"index repo list failed"}`+"\n", lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		t.Parallel()

		ir := &indexRepoMock{}
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		ir.ListFunc = func(ctx context.Context) ([]model.Index, error) {
			return []model.Index{
				{
					ID:        123,
					Name:      "theIndex1",
					Schema:    []byte("theSchema1"),
					CreatedAt: time.Unix(234, 0),
					UpdatedAt: time.Unix(345, 0),
				},
				{
					ID:        321,
					Name:      "theIndex2",
					Schema:    []byte("theSchema2"),
					CreatedAt: time.Unix(432, 0),
					UpdatedAt: time.Unix(543, 0),
				},
			}, nil
		}

		h := indexhandler.New(ir, now, l)
		res, err := h.List(context.Background(), connect.NewRequest(&proto.ListRequest{}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Indices, 2)
		assert.Empty(t, lb.String())

		assert.Equal(t, "theIndex1", res.Msg.Indices[0].Name)
		assert.Equal(t, "theIndex2", res.Msg.Indices[1].Name)
	})
}
