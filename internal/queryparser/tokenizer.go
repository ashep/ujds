package queryparser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func tokenize(s string) ([]token, error) {
	var (
		tok token
		err error
		res = make([]token, 0)
	)

	nextTK := tkIdentifier
	for pos := 0; ; {
		if s[pos] == ' ' {
			pos++
			continue
		}

		switch nextTK {
		case tkIdentifier:
			tok, err = parseIdentifier(s, pos)
			nextTK = tkOperatorCompare
		case tkOperatorCompare:
			tok, err = parseOperator(s, pos)
			nextTK = tkLiteralAny
		case tkOperatorLogical:
			tok, err = parseOperator(s, pos)
			nextTK = tkIdentifier
		case tkLiteralAny:
			tok, err = parseLiteral(s, pos)
			nextTK = tkOperatorLogical
		default:
			return nil, fmt.Errorf("unexpected data at position %d; token kind %d", pos, nextTK)
		}

		pos = tok.pos

		if err != nil {
			return nil, fmt.Errorf("%w at position %d: %s[...]", err, pos, s[:pos])
		}

		res = append(res, tok)

		if pos == len(s) {
			break
		}
	}

	return res, nil
}

func parseIdentifier(s string, pos int) (token, error) {
	v := ""

	for stop := false; pos < len(s) && !stop; pos++ {
		switch {
		case s[pos] >= '0' && s[pos] <= '9', s[pos] >= 'A' && s[pos] <= 'Z', s[pos] >= 'a' && s[pos] <= 'z':
			v += string(s[pos])
		case s[pos] == '.', s[pos] == '_':
			v += string(s[pos])
		default:
			pos--
			stop = true
		}
	}

	if len(v) == 0 {
		return token{pos: pos}, errors.New("identifier expected")
	}

	return token{pos: pos, kind: tkIdentifier, value: v}, nil
}

func parseOperator(s string, pos int) (token, error) {
	v := ""

loop:
	for ; pos < len(s); pos++ {
		switch s[pos] {
		case '=', '<', '>', '!', '&', '|':
			v += string(s[pos])
		default:
			break loop
		}
	}

	switch v {
	case opEq, opEqEq, opNeq, opGt, opLt, opGte, opLte:
		return token{pos: pos, kind: tkOperatorCompare, value: v}, nil
	case opAnd, opOr:
		return token{pos: pos, kind: tkOperatorLogical, value: v}, nil
	default:
		return token{pos: pos}, errors.New("operator expected")
	}
}

func parseLiteral(s string, pos int) (token, error) {
	v := ""
	quotesCnt := 0
	nonNumCharsCnt := 0

	for stop := false; !stop && pos < len(s); pos++ {
		switch {
		case s[pos] == '"':
			quotesCnt++
			v += string(s[pos])
			if quotesCnt == 2 {
				stop = true
			}
		case s[pos] >= 'A' && s[pos] <= 'Z', s[pos] >= 'a' && s[pos] <= 'z', s[pos] == '_':
			v += string(s[pos])
			nonNumCharsCnt++
		case (s[pos] >= '0' && s[pos] <= '9') || s[pos] == '.':
			v += string(s[pos])
		default:
			if quotesCnt%2 != 0 {
				v += string(s[pos])
			} else {
				stop = true
			}
		}
	}

	if quotesCnt != 0 || nonNumCharsCnt > 0 {
		if !strings.HasPrefix(v, `"`) {
			v = `"` + v
		}
		if !strings.HasSuffix(v, `"`) {
			v = v + `"`
		}

		return token{pos: pos, kind: tkLiteralString, value: v}, nil
	}

	if strings.Contains(v, ".") {
		if _, err := strconv.ParseFloat(v, 64); err != nil {
			return token{pos: pos}, fmt.Errorf("parse float: %w", err)
		}

		return token{pos: pos, kind: tkLiteralFloat, value: v}, nil
	}

	if _, err := strconv.ParseInt(v, 10, 64); err != nil {
		return token{pos: pos}, fmt.Errorf("parse int: %w", err)
	}

	return token{pos: pos, kind: tkLiteralInt, value: v}, nil
}
