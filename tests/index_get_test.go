//go:build functest

package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/sdk/client"
	ujdsproto "github.com/ashep/ujds/sdk/proto/ujds/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestIndex_Get(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "anInvalidAuthToken", &http.Client{})
		_, err := cli.I.GetIndex(context.Background(), connect.NewRequest(&ujdsproto.GetIndexRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.GetIndex(context.Background(), connect.NewRequest(&ujdsproto.GetIndexRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid name: must match the regexp ^[a-zA-Z0-9_-]{1,64}$")
	})

	tt.Run("InvalidIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.GetIndex(context.Background(), connect.NewRequest(&ujdsproto.GetIndexRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid name: must match the regexp ^[a-zA-Z0-9_-]{1,64}$")
	})

	tt.Run("IndexNotExists", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.GetIndex(context.Background(), connect.NewRequest(&ujdsproto.GetIndexRequest{
			Name: "theIndexName",
		}))

		assert.EqualError(t, err, "not_found: index is not found")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		ta.DB().InsertIndex(t, "theIndexName", `{"foo":"bar"}`)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		res, err := cli.I.GetIndex(context.Background(), connect.NewRequest(&ujdsproto.GetIndexRequest{
			Name: "theIndexName",
		}))

		require.NoError(t, err)
		assert.Equal(t, "theIndexName", res.Msg.Name)
		assert.Equal(t, `{"foo": "bar"}`, res.Msg.Schema)
		assert.NotZero(t, res.Msg.CreatedAt)
		assert.NotZero(t, res.Msg.UpdatedAt)
	})
}
