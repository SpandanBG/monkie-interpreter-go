package parser

import (
	"errors"
	"fmt"
	"strconv"

	"sudocoding.xyz/interpreter_in_go/src/lexer"
	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
	"sudocoding.xyz/interpreter_in_go/src/token"
)

// Operator precedence
type OpPrec int

const (
	_ OpPrec = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !x
	CALL        // function call
)

type (
	// prefixParserFn - function signature for parsing expressions of tokens with
	// prefix
	prefixParserFn func() ast.Expression

	// infixParserFn - function signature for parsing expressions of tokens with
	// infix. Takes the left expression as the arugment to the function
	infixParserFn func(leftExp ast.Expression) ast.Expression
)

// Parser - The parser of Monkie lang
type Parser struct {
	l             *lexer.Lexer_V2
	curToken      token.Token                        // Points to the currently pointing token
	peekToken     token.Token                        // Points to the next token
	errs          []error                            // List of errors that occured while parsing
	prefixParsers map[token.TokenType]prefixParserFn // map of prefix token parsers
	infixParsers  map[token.TokenType]infixParserFn  // map of infin token parsers
}

// New - Create a new parser with the lexer v2 and sets the curToken and the
// peekToken to the start of the program
func New(l *lexer.Lexer_V2) *Parser {
	p := &Parser{
		l:             l,
		prefixParsers: make(map[token.TokenType]prefixParserFn),
		infixParsers:  make(map[token.TokenType]infixParserFn),
	}

	p.registerPrefixParser(token.IDENT, p.parseIdentifier)
	p.registerPrefixParser(token.INT, p.parseIntegerLiteral)
	p.registerPrefixParser(token.BANG, p.parsePrefixExpression)
	p.registerPrefixParser(token.MINUS, p.parsePrefixExpression)

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

// registerPrefixParser - registers prefix token parsers
func (p *Parser) registerPrefixParser(tokenType token.TokenType, fn prefixParserFn) {
	p.prefixParsers[tokenType] = fn
}

// registerInfixParser - registers infix token parsers
func (p *Parser) registerInfixParser(tokenType token.TokenType, fn infixParserFn) {
	p.infixParsers[tokenType] = fn
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

// expectNextToken - tries to move to the next token if matches the expected
// token. otherwise returns an error
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
	default:
		if expStmt := p.parseExpressionStatement(); expStmt != nil {
			return expStmt
		}
		return nil
	}

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

	// TODO: We skipping the expression till we find the semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		if err := p.nextToken(); err != nil {
			fmt.Printf("Error looking for semicolon at end of let statement: %s\n", err.Error())
			return nil
		}
	}

	return stmt
}

// parseReturnStatement - parse a return statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: We skipping the expression till we find the semicolon

	for !p.curTokenIs(token.SEMICOLON) {
		if err := p.nextToken(); err != nil {
			fmt.Printf("Error looking for semicolon at end of let statement: %s\n", err.Error())
			return nil
		}
	}

	return stmt
}

// parseExpressionStatement - parse an expression statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	// This is to check if the next token is ; since ; will be optional for
	// expression statements
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

// parseExpression - parse an expression
func (p *Parser) parseExpression(opPrec OpPrec) ast.Expression {
	prefix := p.prefixParsers[p.curToken.Type]

	if prefix == nil {
		err := errors.New(
			fmt.Sprintf("No prefix parser function for %s found", string(p.curToken.Type)),
		)
		fmt.Println(err.Error())
		p.errs = append(p.errs, err)
		return nil
	}

	return prefix()
}

// parseIdentifier - parse an identifer expression
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral - parse an integer literal expression
func (p *Parser) parseIntegerLiteral() ast.Expression {
	intVal, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		pErr := errors.New(fmt.Sprintf("Error occured while parsing int literal %q with error %q", p.curToken.Literal, err.Error()))

		fmt.Println(pErr.Error())

		p.errs = append(p.errs, pErr)
		return nil
	}

	return &ast.IntegerLiteral{Token: p.curToken, Value: intVal}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	if err := p.nextToken(); err != nil {
		fmt.Println("Error occured while parsing next token of prefix expression: ", err.Error())
		return nil
	}

	expression.Right = p.parseExpression(PREFIX)
	return expression
}
