package ast

import (
	"fmt"
	"reflect"
	"testing"
)

func Test_Modify(t *testing.T) {
	one := func() Expression { return &IntegerLiteral{Value: 1} }
	two := func() Expression { return &IntegerLiteral{Value: 2} }

	turnOneIntoTwo := func(node Node) Node {
		integer, ok := node.(*IntegerLiteral)
		if !ok {
			return node
		}

		if integer.Value != 1 {
			return node
		}

		integer.Value = 2
		return integer
	}

	for _, test := range []struct {
		input    Node
		expected Node
	}{
		{one(), two()},
		{
			&Program{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
			&Program{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
		},
		{
			&InfixExpression{Left: one(), Operator: "+", Right: one()},
			&InfixExpression{Left: two(), Operator: "+", Right: two()},
		},
		{
			&PrefixExpression{Operator: "-", Right: one()},
			&PrefixExpression{Operator: "-", Right: two()},
		},
		{
			&IndexExpression{Left: one(), Index: one()},
			&IndexExpression{Left: two(), Index: two()},
		},
		{
			&IfExpression{
				Condition:   one(),
				Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
				Alternative: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
			},
			&IfExpression{
				Condition:   two(),
				Consequence: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
				Alternative: &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
			},
		},
		{
			&ReturnStatement{ReturnValue: one()},
			&ReturnStatement{ReturnValue: two()},
		},
		{
			&LetStatement{Value: one()},
			&LetStatement{Value: two()},
		},
		{
			&Assignment{Value: one()},
			&Assignment{Value: two()},
		},
		{
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body:       &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: one()}}},
			},
			&FunctionLiteral{
				Parameters: []*Identifier{},
				Body:       &BlockStatement{Statements: []Statement{&ExpressionStatement{Expression: two()}}},
			},
		},
		{
			&ArrayLiteral{Elements: []Expression{one(), one()}},
			&ArrayLiteral{Elements: []Expression{two(), two()}},
		},
	} {
		t.Run(fmt.Sprintf("Test modify for %v", test.input), func(t *testing.T) {
			modified := Modify(test.input, turnOneIntoTwo)

			eq(t, true, reflect.DeepEqual(modified, test.expected), fmt.Sprintf("Expected = %#v\nGot = %#v", test.expected, modified))
		})
	}

	hashLiteral := &HashLiteral{
		Pairs: map[Expression]Expression{
			one(): one(),
			one(): one(),
		},
	}

	Modify(hashLiteral, turnOneIntoTwo)

	for key, value := range hashLiteral.Pairs {
		key, _ := key.(*IntegerLiteral)
		eq(t, int64(2), key.Value)

		value, _ := value.(*IntegerLiteral)
		eq(t, int64(2), value.Value)
	}
}
