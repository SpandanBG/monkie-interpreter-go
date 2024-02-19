package ast

import (
	"bytes"

	"sudocoding.xyz/interpreter_in_go/src/token"
)

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) statementNode() {}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral())
	out.WriteString(" ")

	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	return out.String()
}
