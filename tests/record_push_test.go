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

func TestRecord_Push(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "anInvalidAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{}))

		assert.EqualError(t, err, "invalid_argument: index get failed: invalid name: must match the regexp ^[a-zA-Z0-9_-]{1,64}$")
	})

	tt.Run("IndexNotFound", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{Index: "anIndex"}))

		assert.EqualError(t, err, "not_found: index is not found")
	})

	tt.Run("EmptyRecords", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", "{}")

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{Index: "theIndex"}))

		assert.EqualError(t, err, "invalid_argument: invalid records: must not be empty")
	})

	tt.Run("EmptyRecordID", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", "{}")

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id: "",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: invalid record (0) id: must not be empty")
	})

	tt.Run("EmptyRecordData", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", "{}")

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: "",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: invalid record (0) data: must not be empty")
	})

	tt.Run("InvalidDataJSON", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", "{}")

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: "{]",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: invalid record data (0): invalid json")
	})

	tt.Run("DataValidationFailed", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", `{"required": ["foo"]}`)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: "{}",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: invalid record data (0): (root): foo is required")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", `{"required": ["foo"]}`)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar"}`,
				},
			},
		}))

		require.NoError(t, err)

		rls := ta.DB().GetRecordLogs(t, "theIndex")
		require.Len(t, rls, 1)
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs := ta.DB().GetRecords(t, "theIndex")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 1, rcs[0].LogID)
		assert.Equal(t, []byte{0x27, 0x4d, 0x66, 0x42, 0x1b, 0x85, 0x49, 0x7e, 0x4c, 0x7d, 0xa9, 0xe7, 0xed, 0xa2, 0x37, 0xa6, 0x78, 0x65, 0x4a, 0x28, 0x50, 0x18, 0x54, 0xbc, 0x6c, 0x2d, 0x3c, 0xc1, 0xc5, 0x7a, 0x4f, 0xf2}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.NotZero(t, rcs[0].UpdatedAt)
	})

	tt.Run("OkUpdate", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", `{"required": ["foo"]}`)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar1"}`,
				},
			},
		}))
		require.NoError(t, err)

		_, err = cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar2"}`,
				},
			},
		}))
		require.NoError(t, err)

		rls := ta.DB().GetRecordLogs(t, "theIndex")
		require.Len(t, rls, 2)
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar1"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		assert.Equal(t, 2, rls[1].ID)
		assert.Equal(t, 1, rls[1].IndexID)
		assert.Equal(t, "theRecordID", rls[1].RecordID)
		assert.Equal(t, `{"foo": "bar2"}`, rls[1].Data)
		assert.Greater(t, rls[1].CreatedAt, rls[0].CreatedAt)

		rcs := ta.DB().GetRecords(t, "theIndex")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 2, rcs[0].LogID)
		assert.Equal(t, []byte{0x16, 0xfa, 0x32, 0xf6, 0x63, 0x2, 0x73, 0x62, 0x33, 0xec, 0xd1, 0xce, 0x21, 0xfb, 0x51, 0xbc, 0x9e, 0x40, 0x94, 0xce, 0x5e, 0x7e, 0x14, 0x74, 0xd1, 0xb5, 0x3a, 0xd6, 0x36, 0x3d, 0x3, 0x97}, rcs[0].Checksum)
		assert.Greater(t, rcs[0].UpdatedAt, rcs[0].CreatedAt)
	})

	tt.Run("OkUpdateWithSameData", func(t *testing.T) {
		ta := testapp.New(t)
		ta.DB().InsertIndex(t, "theIndex", `{"required": ["foo"]}`)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar"}`,
				},
			},
		}))
		require.NoError(t, err)

		_, err = cli.R.PushRecords(context.Background(), connect.NewRequest(&ujdsproto.PushRecordsRequest{
			Index: "theIndex",
			Records: []*ujdsproto.PushRecordsRequest_NewRecord{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar"}`,
				},
			},
		}))
		require.NoError(t, err)

		rls := ta.DB().GetRecordLogs(t, "theIndex")
		require.Len(t, rls, 1)
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs := ta.DB().GetRecords(t, "theIndex")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 1, rcs[0].LogID)
		assert.Equal(t, []byte{0x27, 0x4d, 0x66, 0x42, 0x1b, 0x85, 0x49, 0x7e, 0x4c, 0x7d, 0xa9, 0xe7, 0xed, 0xa2, 0x37, 0xa6, 0x78, 0x65, 0x4a, 0x28, 0x50, 0x18, 0x54, 0xbc, 0x6c, 0x2d, 0x3c, 0xc1, 0xc5, 0x7a, 0x4f, 0xf2}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.Equal(t, rcs[0].UpdatedAt, rcs[0].CreatedAt)
	})
}
