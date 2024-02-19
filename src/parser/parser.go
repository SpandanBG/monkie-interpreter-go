package parser

import (
	"errors"
	"fmt"

	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
	"sudocoding.xyz/interpreter_in_go/src/token"
)

// Parser - The parser of Monkie lang
type Parser struct {
	l         *lexer.Lexer_V2
	curToken  token.Token // Points to the currently pointing token
	peekToken token.Token // Points to the next token
	errs      []error     // List of errors that occured while parsing
}

// New - Create a new parser with the lexer v2 and sets the curToken and the
// peekToken to the start of the program
func New(l *lexer.Lexer_V2) *Parser {
	p := &Parser{l: l}

	p.nextToken()
	p.nextToken()

	return p
}

// ParseProgram - Start creating the AST of the input souce code
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.curTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.nextToken()
	}

	return program
}

// Errors - Returns the list of errors that was found by the parser
func (p *Parser) Errors() []error {
	return p.errs
}

// nextToken - updates the curToken and peekToken to their next respective values
func (p *Parser) nextToken() error {
	if p.curTokenIs(token.EOF) {
		err := errors.New("Attempted to read next token at end of file")
		p.errs = append(p.errs, err)
		return err
	}

	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken_V2()
	return nil
}

func (p *Parser) expectNextToken(expectedType token.TokenType) error {
	if !p.peekTokenIs(expectedType) {
		err := errors.New(
			fmt.Sprintf("Next token expected %+v. Got %+v.", expectedType, p.peekToken.Type),
		)
		p.errs = append(p.errs, err)
		return err
	}
	p.nextToken()
	return nil
}

// curTokenIs - Returns true if the curTokken in the parser matches the expected
func (p *Parser) curTokenIs(expected token.TokenType) bool {
	return p.curToken.Type == expected
}

// peekTokenIs - Returns true if the peekToken in the parser matches the expected
func (p *Parser) peekTokenIs(expected token.TokenType) bool {
	return p.peekToken.Type == expected
}

// parseStatement - checks the list of statement tokens parses accordingly
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		// We check for nil here, cause in Go, the nil interface will have the type
		// thus the ast.Statement will not be considered nil even if the data is empty
		// cause it will have the type.
		if letStmt := p.parseLetStatement(); letStmt != nil {
			return letStmt
		}
		return nil
	case token.RETURN:
		if returnStmt := p.parseReturnStatement(); returnStmt != nil {
			return returnStmt
		}
		return nil
	}

	return nil
}

// parseExpression - parse an expression
func (p *Parser) parseExpression() ast.Expression {
	// TODO: We skipping the expression till we find the semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		if err := p.nextToken(); err != nil {
			fmt.Printf("Error looking for semicolon at end of let statement: %s\n", err.Error())
			return nil
		}
	}

	return nil
}

// parseLetStatement - parses a let statement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if err := p.expectNextToken(token.IDENT); err != nil {
		fmt.Printf("Error while parsing let statement: %s\n", err.Error())
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	if err := p.expectNextToken(token.ASSIGN); err != nil {
		fmt.Printf("Error while parsing let statement: %s\n", err.Error())
		return nil
	}

	stmt.Value = p.parseExpression()
	return stmt
}

// parseReturnStatement - parse a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()
	stmt.ReturnValue = p.parseExpression()
	return stmt
}