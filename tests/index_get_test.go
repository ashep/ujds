//go:build functest

package tests

import (
	"context"
	"encoding/json"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestIndex_Get(main *testing.T) {
	main.Parallel()

	main.Run("InvalidAuthorization", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyIndexName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("InvalidIndexName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must match the regexp ^[a-zA-Z0-9.-]{1,255}$")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("IndexNotExists", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "not_found: index is not found")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("Ok", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		ta.DB().InsertIndex("theIndexName", "theIndexTitle")

		cli := ta.Client("")

		res, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Equal(t, "theIndexName", res.Msg.Name)
		assert.Equal(t, "theIndexTitle", res.Msg.Title)
		assert.NotZero(t, res.Msg.CreatedAt)
		assert.NotZero(t, res.Msg.UpdatedAt)
		assert.Empty(t, res.Msg.Schemas)

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkWithSchemas", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t, testapp.WithConfigOptionValidationIndex(
			"theIndex.*", json.RawMessage(`{"type":"object","required":["title"]}`)))

		ta.DB().InsertIndex("theIndexName", "theIndexTitle")

		cli := ta.Client("")

		res, err := cli.I.Get(context.Background(), connect.NewRequest(&indexproto.GetRequest{
			Name: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Equal(t, "theIndexName", res.Msg.Name)
		assert.Equal(t, []string{
			`{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","required":["title"]}`,
		}, res.Msg.Schemas)

		ta.AssertNoWarnsAndErrors()
	})
}
