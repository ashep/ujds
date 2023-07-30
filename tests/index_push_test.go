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

func TestIndexPush(tt *testing.T) {
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
	})

	tt.Run("OkWithSchema", func(t *testing.T) {
		ta := testapp.New(t)

		defer ta.Start(t)()
		defer ta.AssertNoLogErrors(t)

		cli := client.New("http://localhost:9000", "theAuthToken", &http.Client{})
		_, err := cli.I.PushIndex(context.Background(), connect.NewRequest(&v1.PushIndexRequest{
			Name:   "theIndexName",
			Schema: `{"foo":"bar"}`,
		}))

		assert.NoError(t, err)
	})
}
