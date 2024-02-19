package ast

import "sudocoding.xyz/interpreter_in_go/src/token"

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression == nil {
		return ""
	}
	return es.Expression.String()
}
