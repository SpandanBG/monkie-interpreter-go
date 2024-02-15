package ast

/*
 Top Down Operator Precedence (Pratt Parser)

 This is a recursive descent top-down parser.
*/

// ----------------------------------------------------------------------------

// Node - The AST is literally a tree with nodes. Each node is suppose to
// implement the `TokenLiteral` method which is to return the literal of the
// token. This will be used for debugging and testing purpose.
type Node interface {
	TokenLiteral() string
}

// Statement - The nodes in the AST that are statements are to implement this.
// the `statementNode` method is a dummy method to guide the Go compiler and
// possible causing it to throw errors when we use statement where an expression
// shouldn't been used.
type Statement interface {
	Node

	statementNode()
}

// Expression - The nodes in the AST that are expression to implement this.
// the `statementNode` method is a dummy method to guide the Go compiler and
// possible causing it to throw errors when we use where expression an statement
// shouldn't been used.
type Expression interface {
	Node

	expressionNode()
}

// Program - This is the root node of AST
type Program struct {
	Statements []Statement
}

// TokenLiteral - The root node's implementation of `TokenLiteral`.
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
