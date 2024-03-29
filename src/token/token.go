package token

// TokenType - Defines the type of token. Here we are using string so that we
// can have full flexibility for different kinds of token. This will also help
// us to debug the string out.
// However, a more efficient way for token types would be to use integers or
// byte.
type TokenType string

// Token - We hold the token type and its corresponding literal
type Token struct {
	Type    TokenType
	Literal string
}

// The tokens of our langauge
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF               = "EOF"

	// Identifiers + Literals
	IDENT TokenType = "IDENT" // add, foobar, x, y, ...
	INT             = "INT"   // 12345
	STR             = "STR"   // string

	// Operators
	ASSIGN   TokenType = "="
	PLUS               = "+"
	MINUS              = "-"
	BANG               = "!"
	ASTERISK           = "*"
	SLASH              = "/"

	// Equality
	GT     TokenType = ">"
	LT               = "<"
	EQ               = "=="
	NOT_EQ           = "!="
	GTE              = ">="
	LTE              = "<="

	// Delimiters
	COMMA     TokenType = ","
	SEMICOLON           = ";"
	COLON               = ":"
	LPAREN              = "("
	RPAREN              = ")"
	LBRACE              = "{"
	RBRACE              = "}"
	LBRACKET            = "["
	RBRACKET            = "]"

	// Keywords
	FUNCTION TokenType = "FUNCTION"
	LET                = "LET"
	TRUE               = "TRUE"
	FALSE              = "FALSE"
	IF                 = "IF"
	ELSE               = "ELSE"
	RETURN             = "RETURN"
	MACRO              = "MACRO"

	// String Tokens
	DOUBLE_QUOTES TokenType = "\""
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"macro":  MACRO,
}

// LookupIdent - Checks the keywords map. If the keyword is mapped to a token type
// returns the token type otherwise returns token type as `IDENT`
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
