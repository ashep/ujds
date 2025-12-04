//go:build functest

package tests

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestRecord_History(main *testing.T) {
	main.Parallel()

	main.Run("InvalidAuthorization", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("anInvalidAuthToken")

		_, err := cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("NoRecords", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		res, err := cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index: "theIndex",
			Id:    "theRecord",
		}))

		require.NoError(t, err)
		assert.Empty(t, res.Msg.Records)
		assert.Zero(t, res.Msg.Cursor)
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyIndexName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index: "",
			Id:    "theRecord",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyRecordID", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index: "theIndex",
			Id:    "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid record id: must not be empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("Ok", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex", Id: "theRecord", Data: `{"foo1":"bar1"}`},
				{Index: "theIndex", Id: "theRecord", Data: `{"foo2":"bar2"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index: "theIndex",
			Id:    "theRecord",
		}))
		require.NoError(t, err)

		require.Len(t, res.Msg.Records, 2)
		require.Zero(t, res.Msg.Cursor)

		assert.Equal(t, "theRecord", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(2), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[0].Index)
		assert.Equal(t, `{"foo2": "bar2"}`, res.Msg.Records[0].Data)

		assert.Equal(t, "theRecord", res.Msg.Records[1].Id)
		assert.Equal(t, uint64(1), res.Msg.Records[1].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[1].Index)
		assert.Equal(t, `{"foo1": "bar1"}`, res.Msg.Records[1].Data)

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkPaginated", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex", Id: "theRecord", Data: `{"foo1":"bar1"}`},
				{Index: "theIndex", Id: "theRecord", Data: `{"foo2":"bar2"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index: "theIndex",
			Id:    "theRecord",
			Limit: 1,
		}))
		require.NoError(t, err)

		require.Len(t, res.Msg.Records, 1)
		require.Equal(t, uint64(2), res.Msg.Cursor)
		assert.Equal(t, "theRecord", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(2), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[0].Index)
		assert.Equal(t, `{"foo2": "bar2"}`, res.Msg.Records[0].Data)

		res, err = cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index:  "theIndex",
			Id:     "theRecord",
			Limit:  1,
			Cursor: 2,
		}))
		require.NoError(t, err)

		require.Len(t, res.Msg.Records, 1)
		require.Zero(t, res.Msg.Cursor)
		assert.Equal(t, "theRecord", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(1), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[0].Index)
		assert.Equal(t, `{"foo1": "bar1"}`, res.Msg.Records[0].Data)

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkSince", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex", Id: "theRecord", Data: `{"foo1":"bar1"}`},
			},
		}))
		require.NoError(t, err)

		time.Sleep(time.Second)
		since := time.Now()

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex", Id: "theRecord", Data: `{"foo2":"bar2"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index: "theIndex",
			Id:    "theRecord",
			Limit: 1,
		}))
		require.NoError(t, err)

		require.Len(t, res.Msg.Records, 1)
		require.Equal(t, uint64(2), res.Msg.Cursor)
		assert.Equal(t, "theRecord", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(2), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[0].Index)
		assert.Equal(t, `{"foo2": "bar2"}`, res.Msg.Records[0].Data)

		res, err = cli.R.History(context.Background(), connect.NewRequest(&recordproto.HistoryRequest{
			Index: "theIndex",
			Id:    "theRecord",
			Since: since.Unix(),
		}))
		require.NoError(t, err)

		require.Len(t, res.Msg.Records, 1)
		require.Zero(t, res.Msg.Cursor)
		assert.Equal(t, "theRecord", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(2), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[0].Index)
		assert.Equal(t, `{"foo2": "bar2"}`, res.Msg.Records[0].Data)

		ta.AssertNoWarnsAndErrors()
	})
}
