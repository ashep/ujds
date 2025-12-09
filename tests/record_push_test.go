//go:build functest

package tests

import (
	"context"
	"encoding/json"
	"testing"

	"connectrpc.com/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	indexproto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
	recordproto "github.com/ashep/ujds/sdk/proto/ujds/record/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestRecord_Push(main *testing.T) {
	main.Parallel()

	main.Run("InvalidAuthorization", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("anInvalidAuthToken")

		_, err := cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyRecords", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{}))

		assert.EqualError(t, err, "invalid_argument: empty records")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyIndexName", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("IndexNotFound", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "anUnknownIndex",
				},
			},
		}))

		assert.EqualError(t, err, "not_found: index is not found")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyRecordID", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: record 0, id=: validation failed: invalid record id: must not be empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("EmptyRecordData", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  "",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: record 0, id=theRecordID: validation failed: invalid json: empty")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("MalformedDataJSON", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{Name: "theIndex"}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  "{]",
				},
			},
		}))

		assert.EqualError(t, err, "invalid_argument: record 0, id=theRecordID: validation failed: invalid json schema or data: invalid character ']' looking for beginning of object key string")
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("SchemaValidationFailed", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t, testapp.WithConfigOptionValidationIndex("theIndex", json.RawMessage(`{"required": ["foo"]}`)))
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "theIndex",
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  "{}",
				},
			},
		}))

		assert.EqualError(t, err, `invalid_argument: record 0, id=theRecordID: validation failed: invalid json: (root): foo is required`)
		ta.AssertNoWarnsAndErrors()
	})

	main.Run("Ok", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t, testapp.WithConfigOptionValidationIndex("theIndex", json.RawMessage(`{"required": ["foo"]}`)))
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "theIndex",
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  `{"foo":"bar"}`,
				},
			},
		}))

		require.NoError(t, err)

		rls := ta.DB().GetRecordLogs("theIndex")
		require.Len(t, rls, 1)
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs := ta.DB().GetRecords("theIndex")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 1, rcs[0].LogID)
		assert.Equal(t, `{"foo": "bar"}`, rcs[0].Data)
		assert.Equal(t, []byte{0xd7, 0xae, 0xdc, 0xee, 0x96, 0x9c, 0xe8, 0x55, 0x47, 0x83, 0xff, 0x37, 0x31, 0x3a, 0xcd, 0x6b, 0x7e, 0xcc, 0xa, 0xcf, 0x68, 0xcd, 0xfc, 0xdb, 0x86, 0x70, 0xd7, 0x65, 0xa6, 0x2, 0x9c, 0x0}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.Equal(t, rcs[0].CreatedAt, rcs[0].UpdatedAt) // the record has no updates
		assert.Equal(t, rcs[0].CreatedAt, rcs[0].TouchedAt) // the record has no touches

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkUpdate", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t, testapp.WithConfigOptionValidationIndex("theIndex", json.RawMessage(`{"required": ["foo"]}`)))
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "theIndex",
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  `{"foo":"bar1"}`,
				},
			},
		}))
		require.NoError(t, err)

		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  `{"foo":"bar2"}`,
				},
			},
		}))
		require.NoError(t, err)

		// We must have 2 log records: the first one for the insert, the second one for the update
		rls := ta.DB().GetRecordLogs("theIndex")
		require.Len(t, rls, 2)

		// First log record
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar1"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		// Second log record
		assert.Equal(t, 2, rls[1].ID)
		assert.Equal(t, 1, rls[1].IndexID)
		assert.Equal(t, "theRecordID", rls[1].RecordID)
		assert.Equal(t, `{"foo": "bar2"}`, rls[1].Data)
		assert.Greater(t, rls[1].CreatedAt, rls[0].CreatedAt)

		// Actual record state
		rcs := ta.DB().GetRecords("theIndex")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 2, rcs[0].LogID)
		assert.Equal(t, `{"foo": "bar2"}`, rcs[0].Data)
		assert.Equal(t, []byte{0x2e, 0x82, 0xd8, 0x32, 0xa1, 0x18, 0xff, 0x72, 0x77, 0x32, 0xf7, 0xb4, 0xec, 0x4f, 0x9c, 0xef, 0xf1, 0x16, 0x97, 0x5, 0xfc, 0xc7, 0xa0, 0xd1, 0xe9, 0x9f, 0xbb, 0x6a, 0x91, 0xca, 0x23, 0x72}, rcs[0].Checksum)
		assert.Greater(t, rcs[0].UpdatedAt, rcs[0].CreatedAt) // the record was updated after creation
		assert.Equal(t, rcs[0].UpdatedAt, rcs[0].TouchedAt)   // and was not touched after the update

		ta.AssertNoWarnsAndErrors()
	})

	main.Run("OkUpdateWithSameData", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t, testapp.WithConfigOptionValidationIndex("theIndex", json.RawMessage(`{"required": ["foo"]}`)))
		cli := ta.Client("")

		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "theIndex",
		}))
		require.NoError(t, err)

		// Insert record
		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  `{"foo":"bar"}`,
				},
			},
		}))
		require.NoError(t, err)

		// Update it
		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex",
					Id:    "theRecordID",
					Data:  `{"foo":"bar"}`, // the same as on first push
				},
			},
		}))
		require.NoError(t, err)

		rls := ta.DB().GetRecordLogs("theIndex")
		require.Len(t, rls, 1) // only one log record despite two pushes
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo": "bar"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs := ta.DB().GetRecords("theIndex")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 1, rcs[0].LogID)
		assert.Equal(t, []byte{0xd7, 0xae, 0xdc, 0xee, 0x96, 0x9c, 0xe8, 0x55, 0x47, 0x83, 0xff, 0x37, 0x31, 0x3a, 0xcd, 0x6b, 0x7e, 0xcc, 0xa, 0xcf, 0x68, 0xcd, 0xfc, 0xdb, 0x86, 0x70, 0xd7, 0x65, 0xa6, 0x2, 0x9c, 0x0}, rcs[0].Checksum)
		assert.NotZero(t, rcs[0].CreatedAt)
		assert.Equal(t, rcs[0].CreatedAt, rcs[0].UpdatedAt)   // no data updated after second push
		assert.Greater(t, rcs[0].TouchedAt, rcs[0].UpdatedAt) // but the record was touched

		ta.AssertNoWarnsAndErrors()
	})

	// Check that records with same IDs but in different indices are not interfering
	main.Run("OkDifferentIndicesWithSameData", func(t *testing.T) {
		t.Parallel()
		ta := testapp.New(t)
		cli := ta.Client("")

		// Create first index
		_, err := cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "theIndex1",
		}))
		require.NoError(t, err)

		// Create second index
		_, err = cli.I.Push(context.Background(), connect.NewRequest(&indexproto.PushRequest{
			Name: "theIndex2",
		}))
		require.NoError(t, err)

		// Push a record to the first index
		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex1",
					Id:    "theRecordID",
					Data:  `{"foo1":"bar1"}`,
				},
			},
		}))
		require.NoError(t, err)

		// Push a record to the second index; note the same record ID as in the first index
		_, err = cli.R.Push(context.Background(), connect.NewRequest(&recordproto.PushRequest{
			Records: []*recordproto.PushRequest_Record{
				{
					Index: "theIndex2",
					Id:    "theRecordID",
					Data:  `{"foo2":"bar2"}`,
				},
			},
		}))
		require.NoError(t, err)

		rls := ta.DB().GetRecordLogs("theIndex1")
		require.Len(t, rls, 1)
		assert.Equal(t, 1, rls[0].ID)
		assert.Equal(t, 1, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo1": "bar1"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs := ta.DB().GetRecords("theIndex1")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 1, rcs[0].IndexID)
		assert.Equal(t, 1, rcs[0].LogID)
		assert.Equal(t, `{"foo1": "bar1"}`, rcs[0].Data)
		assert.Equal(t, []byte{0x0, 0xdb, 0x9a, 0xb1, 0xa1, 0xff, 0x7d, 0x9d, 0x7a, 0x93, 0xf1, 0xf7, 0x4f, 0xb9, 0x20, 0x7, 0x34, 0x5f, 0xa3, 0x85, 0x5c, 0xd6, 0x98, 0xcc, 0x9e, 0x35, 0x5d, 0x43, 0x93, 0x4e, 0x64, 0x90}, rcs[0].Checksum)

		rls = ta.DB().GetRecordLogs("theIndex2")
		require.Len(t, rls, 1)
		assert.Equal(t, 2, rls[0].ID)
		assert.Equal(t, 2, rls[0].IndexID)
		assert.Equal(t, "theRecordID", rls[0].RecordID)
		assert.Equal(t, `{"foo2": "bar2"}`, rls[0].Data)
		assert.NotZero(t, rls[0].CreatedAt)

		rcs = ta.DB().GetRecords("theIndex2")
		require.Len(t, rcs, 1)
		assert.Equal(t, "theRecordID", rcs[0].ID)
		assert.Equal(t, 2, rcs[0].IndexID)
		assert.Equal(t, 2, rcs[0].LogID)
		assert.Equal(t, `{"foo2": "bar2"}`, rcs[0].Data)
		assert.Equal(t, []byte{0x3c, 0x79, 0x3c, 0xf4, 0xdd, 0x1c, 0x1f, 0x6d, 0x37, 0x11, 0xe3, 0x1c, 0xaf, 0xcf, 0x74, 0xe0, 0xcc, 0x9f, 0x7b, 0xcb, 0x1d, 0x1f, 0x3b, 0x58, 0xb2, 0x72, 0xe3, 0x60, 0x6e, 0x61, 0xbe, 0x8d}, rcs[0].Checksum)

		ta.AssertNoWarnsAndErrors()
	})
}
