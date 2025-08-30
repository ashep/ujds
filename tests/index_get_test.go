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

func TestIndex_Get(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
		ta.AssertNoLogErrors()
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		ta.AssertNoLogErrors()
	})

	tt.Run("InvalidIndexName", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must match the regexp ^[a-zA-Z0-9.-]{1,255}$")
		ta.AssertNoLogErrors()
	})

	tt.Run("IndexNotExists", func(t *testing.T) {
		ta := testapp.New(t).Start()

		cli := ta.Client("")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "not_found: index is not found")
		ta.AssertNoLogErrors()
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t).Start()

		ta.DB().InsertIndex("theIndexName", "theIndexTitle", `{"foo":"bar"}`)

		cli := ta.Client("")

		res, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Equal(t, "theIndexName", res.Msg.Name)
		assert.Equal(t, "theIndexTitle", res.Msg.Title)
		assert.Equal(t, `{"foo": "bar"}`, res.Msg.Schema)
		assert.NotZero(t, res.Msg.CreatedAt)
		assert.NotZero(t, res.Msg.UpdatedAt)

		ta.AssertNoLogErrors()
	})
}
