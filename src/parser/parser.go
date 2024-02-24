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
	DIVIDE      // /
	PREFIX      // -X or !x
	CALL        // function call
)

var precedences = map[token.TokenType]OpPrec{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.LTE:      LESSGREATER,
	token.GTE:      LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    DIVIDE,
}

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
	p.registerPrefixParser(token.TRUE, p.parseBoolean)
	p.registerPrefixParser(token.FALSE, p.parseBoolean)
	p.registerPrefixParser(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefixParser(token.IF, p.parseIfExpression)
	p.registerPrefixParser(token.FUNCTION, p.parseFunctionLiteral)

	p.registerInfixParser(token.PLUS, p.parseInfixExpression)
	p.registerInfixParser(token.MINUS, p.parseInfixExpression)
	p.registerInfixParser(token.SLASH, p.parseInfixExpression)
	p.registerInfixParser(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixParser(token.EQ, p.parseInfixExpression)
	p.registerInfixParser(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfixParser(token.LT, p.parseInfixExpression)
	p.registerInfixParser(token.LTE, p.parseInfixExpression)
	p.registerInfixParser(token.GT, p.parseInfixExpression)
	p.registerInfixParser(token.GTE, p.parseInfixExpression)

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

// peekPrecedence - returns the precedence value for the peekToken
func (p *Parser) peekPrecedence() OpPrec {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence - returns the precedence value for the curToken
func (p *Parser) curPrecedence() OpPrec {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && opPrec < p.peekPrecedence() {
		infix := p.infixParsers[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
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

// parsePrefixExpression - parse a prefix expression
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

// parseInfixExpression - parse an infix expression
func (p *Parser) parseInfixExpression(leftExp ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Left:     leftExp,
		Operator: p.curToken.Literal,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseBoolean - parse a boolean
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curTokenIs(token.TRUE),
	}
}

// parseGroupedExpression - parses a grouped (anything within LPAREN and RPAREN)
// expression
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()
	exp := p.parseExpression(LOWEST)
	if err := p.expectNextToken(token.RPAREN); err != nil {
		fmt.Println("Error occured while parsing grouped expression: ", err.Error())
		return nil
	}
	return exp
}

// parseIfExpression - parse an if expression
func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.curToken}

	if err := p.expectNextToken(token.LPAREN); err != nil {
		fmt.Println("Expected missing ( in if condition: ", err.Error())
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if err := p.expectNextToken(token.RPAREN); err != nil {
		fmt.Println("Expected missing ) in if condition: ", err.Error())
		return nil
	}

	if err := p.expectNextToken(token.LBRACE); err != nil {
		fmt.Println("Expected missing { in if body: ", err.Error())
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if err := p.expectNextToken(token.LBRACE); err != nil {
			fmt.Println("Expected missing { in else body: ", err.Error())
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

// parseBlockStatement - parses statements withing curly braces
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if stmt := p.parseStatement(); stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseFunctionLiteral - parse a function
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	if err := p.expectNextToken(token.LPAREN); err != nil {
		fmt.Println("Expected ( for function params is missing: ", err.Error())
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if err := p.expectNextToken(token.LBRACE); err != nil {
		fmt.Println("Expected { for function body is missing: ", err.Error())
		return nil
	}

	lit.Body = p.parseBlockStatement()
	return lit
}

// parseFunctionParameters - parse function parameters
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()
	identifiers = append(identifiers, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		identifiers = append(identifiers, &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal})
	}

	if err := p.expectNextToken(token.RPAREN); err != nil {
		fmt.Println("Expected ) for function params is missing: ", err.Error())
		return nil
	}

	return identifiers
}
