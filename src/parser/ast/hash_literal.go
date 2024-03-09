package ast

import (
	"bytes"
	"fmt"
	"strings"

	"sudocoding.xyz/interpreter_in_go/src/token"
)

// { <expression> : <expression>, ... }
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	list := []string{}
	for key, value := range hl.Pairs {
		list = append(list, fmt.Sprintf("%s : %s", key.String(), value.String()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(list, ", "))
	out.WriteString("}")

	return out.String()
}
