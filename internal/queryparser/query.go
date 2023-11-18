package queryparser

import (
	"errors"
	"fmt"
	"strings"
)

type Query []token

func (q Query) String(fieldName string) string {
	res := ""

	for i, tok := range q {
		switch tok.kind {
		case tkIdentifier:
			res += q.formatIdentifier(fieldName, tok, q[i+2])
		case tkOperatorCompare, tkOperatorLogical:
			res += q.formatOperator(tok)
		case tkLiteralInt, tkLiteralFloat, tkLiteralString:
			res += q.formatLiteral(tok)
		}

		res += " "
	}

	return strings.TrimSpace(res)
}

func (q Query) formatIdentifier(fName string, idf, arg token) string {
	idfParts := strings.Split(idf.value, ".")
	for i := range idfParts {
		idfParts[i] = "'" + idfParts[i] + "'"
	}

	res := "(" + fName + "->" + strings.Join(idfParts, "->") + ")"

	switch arg.kind {
	case tkLiteralInt:
		res += "::int"
	case tkLiteralFloat:
		res += "::float"
	}

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

func checkSyntax(input string, query Query) error {
	exprIsComplete := false
	expectedTKinds := []tKind{tkIdentifier}

	for _, tok := range query {
		found := false
		for _, tk := range expectedTKinds {
			if tok.kind == tk {
				found = true
				break
			}
		}

		if !found {
			var expStr []string
			for _, tk := range expectedTKinds {
				expStr = append(expStr, tk.String())
			}
			return fmt.Errorf("%s exepected: %s[...]", strings.Join(expStr, ", "), input[:tok.pos])
		}

		switch tok.kind {
		case tkIdentifier:
			expectedTKinds = []tKind{tkOperatorCompare}
			exprIsComplete = false
		case tkOperatorCompare:
			expectedTKinds = []tKind{tkLiteralInt, tkLiteralFloat, tkLiteralString}
			exprIsComplete = false
		case tkOperatorLogical:
			expectedTKinds = []tKind{tkIdentifier}
			exprIsComplete = false
		case tkLiteralInt, tkLiteralFloat, tkLiteralString:
			expectedTKinds = []tKind{tkOperatorLogical}
			exprIsComplete = true
		default:
			return fmt.Errorf("unexpected token %s at position %d", tok.kind, tok.pos)
		}
	}

	if !exprIsComplete {
		return errors.New("incomplete expression")
	}

	return nil
}
