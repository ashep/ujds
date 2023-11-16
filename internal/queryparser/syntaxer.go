package queryparser

import (
	"errors"
	"fmt"
	"strings"
)

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
