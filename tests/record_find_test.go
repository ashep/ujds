//go:build functest

package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/sdk/client"
	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestRecord_Find(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "anInvalidAuthToken", &http.Client{})
		_, err := cli.R.Find(context.Background(), connect.NewRequest(&recordproto.FindRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.Find(context.Background(), connect.NewRequest(&recordproto.FindRequest{}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
	})

	tt.Run("NoRecordsFound", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		res, err := cli.R.Find(context.Background(), connect.NewRequest(&recordproto.FindRequest{
			Index: "theIndex",
		}))

		assert.NoError(t, err)
		assert.Empty(t, res.Msg.Records)
		assert.Zero(t, res.Msg.Cursor)
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex1"}))
		require.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex2"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex1", Id: "theRecord1", Data: `{"foo1":"bar1"}`},
				{Index: "theIndex2", Id: "theRecord2", Data: `{"foo2":"bar2"}`},
				{Index: "theIndex1", Id: "theRecord3", Data: `{"foo3":"bar3"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.Find(context.Background(), connect.NewRequest(&recordproto.FindRequest{
			Index: "theIndex1",
		}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Records, 2)
		assert.Equal(t, uint64(0), res.Msg.Cursor)

		assert.Equal(t, "theRecord1", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(1), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex1", res.Msg.Records[0].Index)
		assert.NotZero(t, res.Msg.Records[0].CreatedAt)
		assert.Equal(t, res.Msg.Records[0].CreatedAt, res.Msg.Records[0].UpdatedAt)
		assert.Equal(t, res.Msg.Records[0].CreatedAt, res.Msg.Records[0].TouchedAt)
		assert.Equal(t, `{"foo1": "bar1"}`, res.Msg.Records[0].Data)

		assert.Equal(t, "theRecord3", res.Msg.Records[1].Id)
		assert.Equal(t, uint64(3), res.Msg.Records[1].Rev)
		assert.Equal(t, "theIndex1", res.Msg.Records[1].Index)
		assert.NotZero(t, res.Msg.Records[1].CreatedAt)
		assert.Equal(t, res.Msg.Records[1].CreatedAt, res.Msg.Records[1].UpdatedAt)
		assert.Equal(t, res.Msg.Records[1].CreatedAt, res.Msg.Records[1].TouchedAt)
		assert.Equal(t, `{"foo3": "bar3"}`, res.Msg.Records[1].Data)
	})

	tt.Run("OkWithSearch", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex1"}))
		require.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex2"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex1", Id: "theRecord1", Data: `{"foo1":"bar1"}`},
				{Index: "theIndex2", Id: "theRecord2", Data: `{"foo2":"bar2"}`},
				{Index: "theIndex1", Id: "theRecord3", Data: `{"foo3":"bar3"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.Find(context.Background(), connect.NewRequest(&recordproto.FindRequest{
			Index:  "theIndex1",
			Search: "foo3 = bar3",
		}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Records, 1)
		assert.Equal(t, uint64(0), res.Msg.Cursor)

		assert.Equal(t, "theRecord3", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(3), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex1", res.Msg.Records[0].Index)
		assert.NotZero(t, res.Msg.Records[0].CreatedAt)
		assert.Equal(t, res.Msg.Records[0].CreatedAt, res.Msg.Records[0].UpdatedAt)
		assert.Equal(t, res.Msg.Records[0].CreatedAt, res.Msg.Records[0].TouchedAt)
		assert.Equal(t, `{"foo3": "bar3"}`, res.Msg.Records[0].Data)
	})

	tt.Run("OkOffsetLimit", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
				Records: []*recordproto.PushRequest_Record{
					{Index: "theIndex", Id: fmt.Sprintf("theRecord%d", i), Data: fmt.Sprintf(`{"foo%d":"bar%d"}`, i, i)},
				},
			}))
			require.NoError(t, err)
		}

		cur := uint64(0)
		for i := 0; i < 10; i++ {
			res, err := cli.R.Find(context.Background(), connect.NewRequest(&recordproto.FindRequest{
				Index:  "theIndex",
				Limit:  1,
				Cursor: cur,
			}))

			require.NoError(t, err)
			require.Len(t, res.Msg.Records, 1)

			require.Equal(t, fmt.Sprintf("theRecord%d", i), res.Msg.Records[0].Id)
			assert.Equal(t, uint64(i+1), res.Msg.Records[0].Rev)
			assert.Equal(t, "theIndex", res.Msg.Records[0].Index)
			assert.NotZero(t, res.Msg.Records[0].CreatedAt)
			assert.Equal(t, res.Msg.Records[0].CreatedAt, res.Msg.Records[0].UpdatedAt)
			assert.Equal(t, res.Msg.Records[0].CreatedAt, res.Msg.Records[0].TouchedAt)
			assert.Equal(t, fmt.Sprintf(`{"foo%d": "bar%d"}`, i, i), res.Msg.Records[0].Data)

			cur = res.Msg.Cursor
		}
	})

	tt.Run("OkSince", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		for i := 0; i < 10; i++ {
			_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
				Records: []*recordproto.PushRequest_Record{
					{Index: "theIndex", Id: fmt.Sprintf("theRecord%d", i), Data: fmt.Sprintf(`{"foo%d":"bar%d"}`, i, i)},
				},
			}))
			require.NoError(t, err)
		}

		time.Sleep(time.Second * 2)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{Index: "theIndex", Id: "theRecord0", Data: `{"foo00":"bar00"}`},
				{Index: "theIndex", Id: "theRecord5", Data: `{"foo55":"bar55"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.Find(context.Background(), connect.NewRequest(&recordproto.FindRequest{
			Index: "theIndex",
			Since: time.Now().Add(-time.Second).Unix(),
		}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Records, 2)

		assert.Equal(t, "theRecord0", res.Msg.Records[0].Id)
		assert.Equal(t, uint64(11), res.Msg.Records[0].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[0].Index)
		assert.NotZero(t, res.Msg.Records[0].CreatedAt)
		assert.Greater(t, res.Msg.Records[0].UpdatedAt, res.Msg.Records[0].CreatedAt)
		assert.Equal(t, res.Msg.Records[0].TouchedAt, res.Msg.Records[0].UpdatedAt)
		assert.Equal(t, `{"foo00": "bar00"}`, res.Msg.Records[0].Data)

		assert.Equal(t, "theRecord5", res.Msg.Records[1].Id)
		assert.Equal(t, uint64(12), res.Msg.Records[1].Rev)
		assert.Equal(t, "theIndex", res.Msg.Records[1].Index)
		assert.NotZero(t, res.Msg.Records[1].CreatedAt)
		assert.Greater(t, res.Msg.Records[1].UpdatedAt, res.Msg.Records[1].CreatedAt)
		assert.Equal(t, res.Msg.Records[1].TouchedAt, res.Msg.Records[1].UpdatedAt)
		assert.Equal(t, `{"foo55": "bar55"}`, res.Msg.Records[1].Data)
	})
}
