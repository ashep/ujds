package queryparser_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ashep/ujds/internal/queryparser"
)

// select * from record_log where data->>'brand'='Brooks Brothers' and (data->'srp')::int>150;

func TestParse_SyntaxErrorIdentifierExpected(tt *testing.T) {
	tt.Run("1", func(t *testing.T) {
		_, err := queryparser.Parse(`"`)
		assert.EqualError(t, err, "identifier expected at position 0: [...]")
	})

	tt.Run("2", func(t *testing.T) {
		_, err := queryparser.Parse(`=`)
		assert.EqualError(t, err, "identifier expected at position 0: [...]")
	})

	tt.Run("3", func(t *testing.T) {
		_, err := queryparser.Parse(`&&`)
		assert.EqualError(t, err, "identifier expected at position 0: [...]")
	})

	tt.Run("4", func(t *testing.T) {
		_, err := queryparser.Parse(`a = 1 && &&`)
		assert.EqualError(t, err, "identifier expected at position 9: a = 1 && [...]")
	})
}

func TestParse_Basic(tt *testing.T) {
	tt.Run("SingleIdentifier", func(t *testing.T) {
		_, err := queryparser.Parse(`foo`)
		assert.EqualError(t, err, "incomplete expression")
	})

	tt.Run("SingleIdentifierAndOperator", func(t *testing.T) {
		_, err := queryparser.Parse(`foo=`)
		assert.EqualError(t, err, "incomplete expression")
	})

	tt.Run("SingleIdentifierAndOperatorAndIdentifier", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123`, q.String("data"))
	})

	tt.Run("OperatorLeftWhitespace", func(t *testing.T) {
		q, err := queryparser.Parse(`foo =123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123`, q.String("data"))
	})

	tt.Run("OperatorRightWhitespace", func(t *testing.T) {
		q, err := queryparser.Parse(`foo= 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123`, q.String("data"))
	})

	tt.Run("OperatorBothWhitespace", func(t *testing.T) {
		q, err := queryparser.Parse(`foo = 123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123`, q.String("data"))
	})

	tt.Run("LogicalOperator", func(t *testing.T) {
		_, err := queryparser.Parse(`foo = 123 &&`)
		assert.EqualError(t, err, "incomplete expression")
	})
}

func TestParse_Literal(tt *testing.T) {
	tt.Run("String", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=bar`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = '"bar"'`, q.String("data"))
	})

	tt.Run("QuotedString", func(t *testing.T) {
		q, err := queryparser.Parse(`foo="bar"`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = '"bar"'`, q.String("data"))
	})

	tt.Run("Int", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123`, q.String("data"))
	})

	tt.Run("Float", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=123.45`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123.45`, q.String("data"))
	})
}

func TestParse_OperatorEq(tt *testing.T) {
	tt.Run("Eq", func(t *testing.T) {
		q, err := queryparser.Parse(`foo=123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123`, q.String("data"))
	})

	tt.Run("EqEq", func(t *testing.T) {
		q, err := queryparser.Parse(`foo==123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') = 123`, q.String("data"))
	})
}

func TestParse_OperatorNEq(tt *testing.T) {
	tt.Run("NEq", func(t *testing.T) {
		q, err := queryparser.Parse(`foo!=123`)
		require.NoError(t, err)
		assert.Equal(t, `(data->'foo') != 123`, q.String("data"))
	})
}
