//go:build functest

package tests

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"

	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestIndex_Clear(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")
		_, err := cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{
			Name: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
	})

	tt.Run("InvalidIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")
		_, err := cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{
			Name: "the n@me",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must match the regexp ^[a-zA-Z0-9.-]{1,255}$")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")
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
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex1", Id: "foo", Data: "{}"},
				{Index: "theIndex1", Id: "bar", Data: "{}"},
				{Index: "theIndex1", Id: "baz", Data: "{}"},
				{Index: "theIndex2", Id: "foo", Data: "{}"},
				{Index: "theIndex2", Id: "bar", Data: "{}"},
				{Index: "theIndex2", Id: "baz", Data: "{}"},
			},
		}))
		assert.NoError(t, err)

		rcs := ta.DB().GetRecords("theIndex1")
		assert.Len(t, rcs, 3)

		rcs = ta.DB().GetRecords("theIndex2")
		assert.Len(t, rcs, 3)

		_, err = cli.I.Clear(context.Background(), connect.NewRequest(&indexproto.ClearRequest{
			Name: "theIndex1",
		}))
		assert.NoError(t, err)

		rcs = ta.DB().GetRecords("theIndex1")
		assert.Len(t, rcs, 0)

		rcs = ta.DB().GetRecords("theIndex2")
		assert.Len(t, rcs, 3)
	})
}
