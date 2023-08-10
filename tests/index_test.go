//go:build functest

package tests

import (
	"context"
	"net/http"
	"testing"

	"github.com/bufbuild/connect-go"
	"github.com/stretchr/testify/assert"

	"github.com/ashep/ujds/sdk/client"
	v1 "github.com/ashep/ujds/sdk/proto/ujds/v1"
	"github.com/ashep/ujds/tests/testapp"
)

func TestIndex_Push(tt *testing.T) {
	tt.Run("InvalidAuthorization", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "anInvalidAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{}))

		assert.EqualError(t, err, "unauthenticated: not authorized")
	})

	tt.Run("EmptyIndexName", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "",
			Schema: "",
		}))

		assert.EqualError(t, err, "invalid_argument: name is empty")
	})

	tt.Run("InvalidSchema", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: "{]",
		}))

		assert.EqualError(t, err, "invalid_argument: invalid schema: invalid character ']' looking for beginning of object key string")
	})

	tt.Run("Ok", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: `{"foo":"bar"}`,
		}))
		assert.NoError(t, err)

		idx := ta.DB().GetIndices(t)
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, `{"foo": "bar"}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Equal(t, idx[0].UpdatedAt, idx[0].CreatedAt)
	})

	tt.Run("OkEmptySchema", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: "",
		}))
		assert.NoError(t, err)

		idx := ta.DB().GetIndices(t)
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, `{}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Equal(t, idx[0].UpdatedAt, idx[0].CreatedAt)
	})

	tt.Run("OkPushTheSame", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: `{"foo":"bar"}`,
		}))
		assert.NoError(t, err)

		_, err = cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: `{"foo":"bar"}`,
		}))
		assert.NoError(t, err)

		idx := ta.DB().GetIndices(t)
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, `{"foo": "bar"}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Greater(t, idx[0].UpdatedAt, idx[0].CreatedAt)
	})

	tt.Run("OkPushUpdate", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: `{"foo1":"bar1"}`,
		}))
		assert.NoError(t, err)

		_, err = cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: `{"foo2":"bar2"}`,
		}))
		assert.NoError(t, err)

		idx := ta.DB().GetIndices(t)
		assert.Len(t, idx, 1)
		assert.Equal(t, 1, idx[0].ID)
		assert.Equal(t, "theIndexName", idx[0].Name)
		assert.Equal(t, `{"foo2": "bar2"}`, idx[0].Schema)
		assert.NotZero(t, idx[0].CreatedAt)
		assert.Greater(t, idx[0].UpdatedAt, idx[0].CreatedAt)
	})
}
