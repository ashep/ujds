package queryparser

import (
	"strings"
)

type Query []token

func (q Query) String(fieldName string) string {
	res := ""

	for _, tok := range q {
		switch tok.kind {
		case tkIdentifier:
			res += q.formatIdentifier(fieldName, tok)
		case tkOperatorCompare, tkOperatorLogical:
			res += q.formatOperator(tok)
		case tkLiteralInt, tkLiteralFloat, tkLiteralString:
			res += q.formatLiteral(tok)
		}

		res += " "
	}

	return strings.TrimSpace(res)
}

func (q Query) formatIdentifier(fName string, idf token) string {
	idfParts := strings.Split(idf.value, ".")
	for i := range idfParts {
		idfParts[i] = "'" + idfParts[i] + "'"
	}

	res := "(" + fName + "->" + strings.Join(idfParts, "->") + ")"

	// switch arg.kind {
	// case tkLiteralInt:
	// 	res += "::int"
	// case tkLiteralFloat:
	// 	res += "::float"
	// }

	return res
}

func (q Query) formatOperator(tok token) string {
	switch tok.value {
	case opEq, opEqEq:
		return "="
	case opAnd:
		return "AND"
	case opOr:
		return "OR"
	default:
		return tok.value
	}
}

func (q Query) formatLiteral(tok token) string {
	switch tok.kind {
	case tkLiteralInt, tkLiteralFloat:
		return tok.value
	default:
		return `'` + tok.value + `'`
	}
}
