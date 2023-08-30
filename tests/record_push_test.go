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
	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestRecord_Push(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "anInvalidAuthToken", &http.Client{})
		_, err := cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{}))

		assert.EqualError(t, err, "invalid_argument: index get failed: invalid name: must match the regexp ^[a-zA-Z0-9_-]{1,64}$")
	})

	tt.Run("IndexNotFound", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{Index: "anIndex"}))

		assert.EqualError(t, err, "not_found: index is not found")
	})

	tt.Run("EmptyRecords", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{Index: "theIndex"}))

		assert.EqualError(t, err, "invalid_argument: invalid records: must not be empty")
	})

	tt.Run("EmptyRecordID", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
				{
					Id: "",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: invalid record (0) id: must not be empty")
	})

	tt.Run("EmptyRecordData", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
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

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
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

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex",
			Schema: `{"required": ["foo"]}`,
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
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

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex",
			Schema: `{"required": ["foo"]}`,
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
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
		assert.Equal(t, []byte{0xd7, 0xae, 0xdc, 0xee, 0x96, 0x9c, 0xe8, 0x55, 0x47, 0x83, 0xff, 0x37, 0x31, 0x3a, 0xcd, 0x6b, 0x7e, 0xcc, 0xa, 0xcf, 0x68, 0xcd, 0xfc, 0xdb, 0x86, 0x70, 0xd7, 0x65, 0xa6, 0x2, 0x9c, 0x0}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.NotZero(t, rcs[0].UpdatedAt)
	})

	tt.Run("OkUpdate", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex",
			Schema: `{"required": ["foo"]}`,
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar1"}`,
				},
			},
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
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
		assert.Equal(t, []byte{0x2e, 0x82, 0xd8, 0x32, 0xa1, 0x18, 0xff, 0x72, 0x77, 0x32, 0xf7, 0xb4, 0xec, 0x4f, 0x9c, 0xef, 0xf1, 0x16, 0x97, 0x5, 0xfc, 0xc7, 0xa0, 0xd1, 0xe9, 0x9f, 0xbb, 0x6a, 0x91, 0xca, 0x23, 0x72}, rcs[0].Checksum)
		assert.Greater(t, rcs[0].UpdatedAt, rcs[0].CreatedAt)
	})

	tt.Run("OkUpdateWithSameData", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex",
			Schema: `{"required": ["foo"]}`,
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar"}`,
				},
			},
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex",
			Records: []*recordproto.PushRequest_Record{
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
		assert.Equal(t, []byte{0xd7, 0xae, 0xdc, 0xee, 0x96, 0x9c, 0xe8, 0x55, 0x47, 0x83, 0xff, 0x37, 0x31, 0x3a, 0xcd, 0x6b, 0x7e, 0xcc, 0xa, 0xcf, 0x68, 0xcd, 0xfc, 0xdb, 0x86, 0x70, 0xd7, 0x65, 0xa6, 0x2, 0x9c, 0x0}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.Equal(t, rcs[0].UpdatedAt, rcs[0].CreatedAt)
	})

	tt.Run("OkDifferentIndicesWithSameData", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex1",
			Schema: `{"required": ["foo"]}`,
		}))
		require.NoError(t, err)

		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name:   "theIndex2",
			Schema: `{"required": ["foo"]}`,
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex1",
			Records: []*recordproto.PushRequest_Record{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar"}`,
				},
			},
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Index: "theIndex2",
			Records: []*recordproto.PushRequest_Record{
				{
					Id:   "theRecordID",
					Data: `{"foo":"bar"}`,
				},
			},
		}))
		require.NoError(t, err)

		rls := ta.DB().GetRecordLogs(t, "theIndex1")
		require.Len(t, rls, 1)
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs := ta.DB().GetRecords(t, "theIndex1")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 1, rcs[0].LogID)
		assert.Equal(t, []byte{0xd7, 0xae, 0xdc, 0xee, 0x96, 0x9c, 0xe8, 0x55, 0x47, 0x83, 0xff, 0x37, 0x31, 0x3a, 0xcd, 0x6b, 0x7e, 0xcc, 0xa, 0xcf, 0x68, 0xcd, 0xfc, 0xdb, 0x86, 0x70, 0xd7, 0x65, 0xa6, 0x2, 0x9c, 0x0}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.Equal(t, rcs[0].UpdatedAt, rcs[0].CreatedAt)

		rls = ta.DB().GetRecordLogs(t, "theIndex2")
		require.Len(t, rls, 1)
		assert.Equal(t, 2, rls[0].ID)
		assert.Equal(t, 2, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs = ta.DB().GetRecords(t, "theIndex2")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 2, rcs[0].IndexID)
		assert.Equal(t, 2, rcs[0].LogID)
		assert.Equal(t, []byte{0x35, 0x80, 0x91, 0x87, 0xc5, 0x85, 0xb6, 0x8a, 0xc8, 0xde, 0xab, 0xcd, 0x94, 0xff, 0x52, 0x1f, 0x50, 0x5e, 0xaa, 0x1c, 0x0, 0x70, 0xe4, 0x71, 0x1a, 0x82, 0x6e, 0xbe, 0x34, 0x25, 0x9e, 0xb5}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.Equal(t, rcs[0].UpdatedAt, rcs[0].CreatedAt)
	})
}
