//go:build functest

package tests

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestIndex_Push(main *testing.T) {
	main.Parallel()

	main.Run("InvalidAuthorization", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyIndexName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("InvalidIndexName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must match the regexp ^[a-zA-Z0-9.-]{1,255}$")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("Ok", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:  "theIndexName",
			Title: "theIndexTitle",
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "theIndexTitle", idx[0].Title.String)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Equal(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkEmptyTitle", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:  "theIndexName",
			Title: "",
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "", idx[0].Title.String)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Equal(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkPushTheSame", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:  "theIndexName",
			Title: "theIndexTitle",
		}))
		require.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:  "theIndexName",
			Title: "theIndexTitle",
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "theIndexTitle", idx[0].Title.String)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Greater(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkPushUpdate", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:  "theIndexName",
			Title: "theIndexTitle1",
		}))
		require.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:  "theIndexName",
			Title: "theIndexTitle2",
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "theIndexTitle2", idx[0].Title.String)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Greater(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoWarnsAndErrors()
	})
}
