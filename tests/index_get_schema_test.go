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

// pairs maps the schema response to {pattern: schema-string} pairs.
func pairs(schemas []*indexproto.GetSchemaResponse_Schema) map[string]string {
	res := make(map[string]string, len(schemas))
	for _, s := range schemas {
		res[s.Pattern] = s.Schema
	}
	return res
}

// patterns returns the schema patterns in response order.
func patterns(schemas []*indexproto.GetSchemaResponse_Schema) []string {
	res := make([]string, 0, len(schemas))
	for _, s := range schemas {
		res = append(res, s.Pattern)
	}
	return res
}

func TestIndex_GetSchema(main *testing.T) {
	main.Parallel()

	main.Run("InvalidAuthorization", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.I.GetSchema(context.Background(), connect.NewRequest(&indexproto.GetSchemaRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.GetSchema(context.Background(), connect.NewRequest(&indexproto.GetSchemaRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("InvalidName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		_, err := cli.I.GetSchema(context.Background(), connect.NewRequest(&indexproto.GetSchemaRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must match the regexp ^[a-zA-Z0-9.-]{1,255}$")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("NoConfiguredSchemasReturnsEmpty", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)

		cli := ta.Client("")
		res, err := cli.I.GetSchema(context.Background(), connect.NewRequest(&indexproto.GetSchemaRequest{
			Name: "books",
		}))

		require.NoError(t, err)
		assert.Empty(t, res.Msg.Schemas)
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("MatchingNameReturnsMatchWithoutCatchAll", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t, testapp.WithConfigOptionValidationIndex(
			"books.*", json.RawMessage(`{"type":"object","required":["title"]}`)))

		cli := ta.Client("")
		res, err := cli.I.GetSchema(context.Background(), connect.NewRequest(&indexproto.GetSchemaRequest{
			Name: "books",
		}))

		require.NoError(t, err)
		assert.Equal(t, []string{"books.*"}, patterns(res.Msg.Schemas))
		assert.Equal(t, map[string]string{
			"books.*": `{"$schema":"http://json-schema.org/draft-07/schema#","type":"object","required":["title"]}`,
		}, pairs(res.Msg.Schemas))
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("NonMatchingNameReturnsEmpty", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t, testapp.WithConfigOptionValidationIndex(
			"books.*", json.RawMessage(`{"type":"object"}`)))

		cli := ta.Client("")
		res, err := cli.I.GetSchema(context.Background(), connect.NewRequest(&indexproto.GetSchemaRequest{
			Name: "movies",
		}))

		require.NoError(t, err)
		assert.Empty(t, res.Msg.Schemas)
		ta.AssertNoWarnsAndErrors()
	})
}
