package indexhandler_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/ashep/go-apperrors"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/rpc/indexhandler"
	"github.com/ashep/ujds/internal/validation"
	proto "github.com/ashep/ujds/sdk/proto/ujds/index/v1"
)

func TestIndexHandler_GetSchema(tt *testing.T) {
	tt.Run("InvalidName", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		nm := &nameValidatorMock{}
		defer nm.AssertExpectations(t)
		nm.On("Validate", "").Return(apperrors.InvalidArgError{Subj: "index name", Reason: "must not be empty"})

		h := indexhandler.New(nil, nil, nm, now, l)
		_, err := h.GetSchema(context.Background(), connect.NewRequest(&proto.GetSchemaRequest{Name: ""}))

		assert.EqualError(t, err, "invalid_argument: invalid index name: must not be empty")
		assert.Empty(t, lb.String())
	})

	tt.Run("Ok", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		nm := &nameValidatorMock{}
		defer nm.AssertExpectations(t)
		nm.On("Validate", "books").Return(nil)

		sm := &schemaMock{}
		defer sm.AssertExpectations(t)
		sm.On("SchemasFor", "books").Return([]validation.Schema{
			{Pattern: ".*", Schema: json.RawMessage(`{}`)},
			{Pattern: "books.*", Schema: json.RawMessage(`{"type":"object","required":["title"]}`)},
		})

		h := indexhandler.New(nil, sm, nm, now, l)
		res, err := h.GetSchema(context.Background(), connect.NewRequest(&proto.GetSchemaRequest{Name: "books"}))

		require.NoError(t, err)
		require.Len(t, res.Msg.Schemas, 2)

		assert.Equal(t, ".*", res.Msg.Schemas[0].Pattern)
		assert.Equal(t, `{}`, res.Msg.Schemas[0].Schema)

		assert.Equal(t, "books.*", res.Msg.Schemas[1].Pattern)
		assert.Equal(t, `{"type":"object","required":["title"]}`, res.Msg.Schemas[1].Schema)

		assert.Empty(t, lb.String())
	})

	tt.Run("NoMatchingSchemas", func(t *testing.T) {
		now := func() time.Time { return time.Unix(123456789, 0) }
		lb := &strings.Builder{}
		l := zerolog.New(lb)

		nm := &nameValidatorMock{}
		defer nm.AssertExpectations(t)
		nm.On("Validate", "books").Return(nil)

		sm := &schemaMock{}
		defer sm.AssertExpectations(t)
		sm.On("SchemasFor", "books").Return([]validation.Schema{})

		h := indexhandler.New(nil, sm, nm, now, l)
		res, err := h.GetSchema(context.Background(), connect.NewRequest(&proto.GetSchemaRequest{Name: "books"}))

		require.NoError(t, err)
		assert.Empty(t, res.Msg.Schemas)
		assert.Empty(t, lb.String())
	})
}
