package queryparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Query struct {
	tokens []token
	args   []any
}

func (q Query) String(fieldName string, firstArgIndex int) string {
	argCnt := firstArgIndex
	res := ""

	for i, tok := range q.tokens {
		switch tok.kind {
		case tkIdentifier:
			res += q.formatIdentifier(fieldName, tok, q.tokens[i+2])
		case tkOperatorCompare, tkOperatorLogical:
			res += q.formatOperator(tok)
		case tkLiteralInt, tkLiteralFloat:
			res += "$" + strconv.Itoa(argCnt)
			argCnt++
		case tkLiteralString:
			// res += `'"$` + strconv.Itoa(argCnt) + `"'`
			res += `'"' || $` + strconv.Itoa(argCnt) + ` || '"'`
			argCnt++
		}

		res += " "
	}

	return strings.TrimSpace(res)
}

func (q Query) Args() []any {
	res := make([]any, 0)

	for _, tok := range q.tokens {
		switch tok.kind {
		case tkLiteralInt:
			v, _ := strconv.Atoi(tok.value.(string))
			res = append(res, v)
		case tkLiteralFloat:
			v, _ := strconv.ParseFloat(tok.value.(string), 64)
			res = append(res, v)
		case tkLiteralString:
			res = append(res, tok.value.(string))
		}
	}

	return res
}

func (q Query) formatIdentifier(fName string, idf, arg token) string {
	idfParts := strings.Split(idf.value.(string), ".")
	for i := range idfParts {
		idfParts[i] = "'" + idfParts[i] + "'"
	}

	res := "(" + fName + "->" + strings.Join(idfParts, "->") + ")"

	switch arg.kind {
	case tkLiteralInt:
		res += "::int"
	case tkLiteralFloat:
		res += "::float"
	case tkLiteralString:
		res += "::text"
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
		return tok.value.(string)
	}
}

func checkSyntax(input string, tokens []token) error {
	exprIsComplete := false
	expectedTKinds := []tKind{tkIdentifier}

	for _, tok := range tokens {
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
			return fmt.Errorf("%s exepected: %s", strings.Join(expStr, ", "), input[:tok.pos])
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
