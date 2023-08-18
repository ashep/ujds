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

func TestRecord_Get(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "anInvalidAuthToken", &http.Client{})
		_, err := cli.R.GetRecord(context.Background(), connect.NewRequest(&ujdsproto.GetRecordRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.GetRecord(context.Background(), connect.NewRequest(&ujdsproto.GetRecordRequest{
			Index: "",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
	})

	tt.Run("IndexNotExists", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.GetRecord(context.Background(), connect.NewRequest(&ujdsproto.GetRecordRequest{
			Index: "theIndex",
		}))

		assert.EqualError(t, err, "not_found: record is not found")
	})

	tt.Run("RecordNotExists", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&ujdsproto.PushIndexRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.GetRecord(context.Background(), connect.NewRequest(&ujdsproto.GetRecordRequest{
			Index: "theIndex",
			Id:    "theRecord",
		}))

		assert.EqualError(t, err, "not_found: record is not found")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&ujdsproto.PushIndexRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{Id: "theRecord", Data: `{"foo":"bar"}`},
			},
		}))
		require.NoError(t, err)

		res, err := cli.R.GetRecord(context.Background(), connect.NewRequest(&ujdsproto.GetRecordRequest{
			Index: "theIndex",
			Id:    "theRecord",
		}))

		require.NoError(t, err)
		assert.Equal(t, "theRecord", res.Msg.Record.Id)
		assert.Equal(t, uint64(1), res.Msg.Record.Rev)
		assert.Equal(t, "theIndex", res.Msg.Record.Index)
		assert.NotEmpty(t, "theIndex", res.Msg.Record.CreatedAt)
		assert.Equal(t, res.Msg.Record.UpdatedAt, res.Msg.Record.CreatedAt)
		assert.Equal(t, `{"foo": "bar"}`, res.Msg.Record.Data)
	})
}
