package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"sudocoding.xyz/interpreter_in_go/src/parser/ast"
)

type ObjectType string
type BuiltinFn func(args ...Object) Object

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	STRING_OBJ       ObjectType = "STRING"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION         ObjectType = "FUNCTION"
	BUILTIN_OBJ      ObjectType = "BUILTIN"
	ARRAY_OBJ        ObjectType = "ARRAY"
	HASH_OBJ         ObjectType = "HASH"
	QUOTE_OBJ        ObjectType = "QUOTE"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer - integer obj type
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Boolean - boolean obj type
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// String - string obj type
type String struct {
	Value string
}

func (b *String) Type() ObjectType {
	return STRING_OBJ
}

func (b *String) Inspect() string {
	return b.Value
}

// Null - null obj type
type Null struct{}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

func (n *Null) Inspect() string {
	return "null"
}

// Return Object - Wrap an object to identify if the object is returned and
// should skip further evaluation of statements in a block
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

// Error Object - err in the program
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s", e.Message)
}

// Function Object - functions
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FUNCTION
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

// Built-in Function - predefined function that Monkie lang provides
type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Type() ObjectType {
	return BUILTIN_OBJ
}

func (b *Builtin) Inspect() string {
	return "built-in function"
}

// Array - arrays
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) Inspect() string {
	var out bytes.Buffer

	elems := []string{}
	for _, e := range a.Elements {
		elems = append(elems, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elems, ", "))
	out.WriteString("]")

	return out.String()
}

// HashKey - hashed key for indexing
type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Hashable interface {
	Hash() HashKey
}

// Hash - hash data struct
type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HASH_OBJ
}

func (h *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

// Hashing Function for Objects

func (b *Boolean) Hash() HashKey {
	if b.Value {
		return HashKey{Type: b.Type(), Value: 1}
	}
	return HashKey{Type: b.Type(), Value: 0}
}

func (i *Integer) Hash() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) Hash() HashKey {
	h := fnv.New64()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Quote - Macro quote AST
type Quote struct {
	Node ast.Node
}

func (q *Quote) Type() ObjectType {
	return QUOTE_OBJ
}

func (q *Quote) Inspect() string {
	return fmt.Sprintf("QUOTE(%s)", q.Node.String())
}
