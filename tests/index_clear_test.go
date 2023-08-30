//go:build functest

package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"

	"github.com/ashep/ujds/sdk/client"
	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestIndex_Clear(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "anInvalidAuthToken", &http.Client{})
		_, err := cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid name: must match the regexp ^[a-zA-Z0-9_-]{1,64}$")
	})

	tt.Run("InvalidIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid name: must match the regexp ^[a-zA-Z0-9_-]{1,64}$")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex1",
			Schema: "",
		}))
		assert.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex2",
			Schema: "",
		}))
		assert.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex1",
			Records: []*recordproto.PushRequest_Record{
				{Id: "foo", Data: "{}"},
				{Id: "bar", Data: "{}"},
				{Id: "baz", Data: "{}"},
			},
		}))
		assert.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex2",
			Records: []*recordproto.PushRequest_Record{
				{Id: "foo", Data: "{}"},
				{Id: "bar", Data: "{}"},
				{Id: "baz", Data: "{}"},
			},
		}))
		assert.NoError(t, err)

		rcs := ta.DB().GetRecords(t, "theIndex1")
		assert.Len(t, rcs, 3)

		rcs = ta.DB().GetRecords(t, "theIndex2")
		assert.Len(t, rcs, 3)

		_, err = cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{
			Name: "theIndex1",
		}))
		assert.NoError(t, err)

		rcs = ta.DB().GetRecords(t, "theIndex1")
		assert.Len(t, rcs, 0)

		rcs = ta.DB().GetRecords(t, "theIndex2")
		assert.Len(t, rcs, 3)
	})
}
