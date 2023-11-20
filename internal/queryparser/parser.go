package queryparser

type tKind int

const (
	tkIdentifier tKind = iota
	tkLiteralAny
	tkLiteralInt
	tkLiteralFloat
	tkLiteralString
	tkOperatorCompare
	tkOperatorLogical
)

func (t tKind) String() string {
	switch t {
	case tkIdentifier:
		return "identifier"
	case tkLiteralAny, tkLiteralInt, tkLiteralFloat, tkLiteralString:
		return "literal"
	case tkOperatorCompare:
		return "comparison operator"
	case tkOperatorLogical:
		return "logical operator"
	default:
		return "unknown"
	}
}

type token struct {
	pos   int
	kind  tKind
	value any
}

func Parse(s string) (Query, error) {
	tokens, err := tokenize(s)
	if err != nil {
		return Query{}, err
	}

	if err = checkSyntax(s, tokens); err != nil {
		return Query{}, err
	}

	return Query{
		tokens: tokens,
		args:   nil,
	}, nil
}
