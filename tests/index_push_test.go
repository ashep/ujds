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

func TestIndex_Push(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
		ta.AssertNoLogErrors()
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		ta.AssertNoLogErrors()
	})

	tt.Run("InvalidIndexName", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must match the regexp ^[a-zA-Z0-9.-]{1,255}$")
		ta.AssertNoLogErrors()
	})

	tt.Run("InvalidSchema", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Schema: "{]",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid schema: invalid character ']' looking for beginning of object key string")
		ta.AssertNoLogErrors()
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Title:  "theIndexTitle",
			Schema: `{"foo":"bar"}`,
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "theIndexTitle", idx[0].Title.String)
		assert.Equal(t, `{"foo": "bar"}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Equal(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoLogErrors()
	})

	tt.Run("OkEmptyTitle", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Schema: "{}",
			Title:  "",
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, `{}`, idx[0].Schema)
		assert.Equal(t, "", idx[0].Title.String)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Equal(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoLogErrors()
	})

	tt.Run("OkEmptySchema", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Title:  "theIndexTitle",
			Schema: "",
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "theIndexTitle", idx[0].Title.String)
		assert.Equal(t, `{}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Equal(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoLogErrors()
	})

	tt.Run("OkPushTheSame", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Title:  "theIndexTitle",
			Schema: `{"foo":"bar"}`,
		}))
		require.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Title:  "theIndexTitle",
			Schema: `{"foo":"bar"}`,
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "theIndexTitle", idx[0].Title.String)
		assert.Equal(t, `{"foo": "bar"}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Greater(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoLogErrors()
	})

	tt.Run("OkPushUpdate", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Title:  "theIndexTitle1",
			Schema: `{"foo1":"bar1"}`,
		}))
		require.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndexName",
			Title:  "theIndexTitle2",
			Schema: `{"foo2":"bar2"}`,
		}))
		require.NoError(t, err)

		idx := ta.DB().GetIndices()
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, "theIndexTitle2", idx[0].Title.String)
		assert.Equal(t, `{"foo2": "bar2"}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Greater(t, idx[0].UpdatedAt, idx[0].CreatedAt)

		ta.AssertNoLogErrors()
	})
}
