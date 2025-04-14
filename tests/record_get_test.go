//go:build functest

package tests

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestRecord_Get(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("anInvalidAuthToken")
		_, err := cli.R.Get(context.Background(), connect.NewRequest(&recordproto.GetRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")
		_, err := cli.R.Get(context.Background(), connect.NewRequest(&recordproto.GetRequest{
			Index: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
	})

	tt.Run("EmptyRecordId", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")
		_, err := cli.R.Get(context.Background(), connect.NewRequest(&recordproto.GetRequest{
			Index: "theIndex",
			Id:    "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid record id: must not be empty")
	})

	tt.Run("IndexNotExists", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")
		_, err := cli.R.Get(context.Background(), connect.NewRequest(&recordproto.GetRequest{
			Index: "theIndex",
			Id:    "theRecord",
		}))

		assert.EqualError(t, err, "not_found: record is not found")
	})

	tt.Run("RecordNotExists", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Get(context.Background(), connect.NewRequest(&recordproto.GetRequest{
			Index: "theIndex",
			Id:    "theRecord",
		}))

		assert.EqualError(t, err, "not_found: record is not found")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start().Stop()
		defer ta.AssertNoLogErrors()

		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex", Id: "theRecord", Data: `{"foo":"bar"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.Get(context.Background(), connect.NewRequest(&recordproto.GetRequest{
			Index: "theIndex",
			Id:    "theRecord",
		}))

		require.NoError(t, err)
		assert.Equal(t, "theRecord", res.Msg.Record.Id)
		assert.Equal(t, uint64(1), res.Msg.Record.Rev)
		assert.Equal(t, "theIndex", res.Msg.Record.Index)
		assert.NotEmpty(t, "theIndex", res.Msg.Record.CreatedAt)
		assert.Equal(t, res.Msg.Record.UpdatedAt, res.Msg.Record.CreatedAt)
		assert.Equal(t, res.Msg.Record.TouchedAt, res.Msg.Record.CreatedAt)
		assert.Equal(t, `{"foo": "bar"}`, res.Msg.Record.Data)
	})
}
