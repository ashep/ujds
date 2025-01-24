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

func TestIndex_List(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.I.List(context.Background(), connect.NewRequest(&indexproto.ListRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		ta.DB().InsertIndex("theIndexName1", "theIndexTitle1", `{}`)
		ta.DB().InsertIndex("theIndexName2", "theIndexTitle2", `{}`)

		cli := ta.Client("")
		res, err := cli.I.List(context.Background(), connect.NewRequest(&indexproto.ListRequest{}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Indices, 2)

		assert.Equal(t, "theIndexName1", res.Msg.Indices[0].Name)
		assert.Equal(t, "theIndexTitle1", res.Msg.Indices[0].Title)

		assert.Equal(t, "theIndexName2", res.Msg.Indices[1].Name)
		assert.Equal(t, "theIndexTitle2", res.Msg.Indices[1].Title)
	})

	tt.Run("OkWithFilter", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		ta.DB().InsertIndex("theIndexName1Foo", "theIndexTitle1", `{}`)
		ta.DB().InsertIndex("theIndexName2Bar", "theIndexTitle2", `{}`)

		cli := ta.Client("")
		res, err := cli.I.List(context.Background(), connect.NewRequest(&indexproto.ListRequest{
			Filter: &indexproto.ListRequestFilter{
				Names: []string{"theIndexName2*"},
			},
		}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Indices, 1)

		assert.Equal(t, "theIndexName2Bar", res.Msg.Indices[0].Name)
		assert.Equal(t, "theIndexTitle2", res.Msg.Indices[0].Title)
	})
}
