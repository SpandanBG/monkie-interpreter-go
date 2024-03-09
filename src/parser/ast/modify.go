package ast

type ModifyFunc func(Node) Node

func Modify(node Node, modifier ModifyFunc) Node {
	switch node := node.(type) {
	case *Program:
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}
	case *ExpressionStatement:
		node.Expression, _ = Modify(node.Expression, modifier).(Expression)
	case *InfixExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *PrefixExpression:
		node.Right, _ = Modify(node.Right, modifier).(Expression)
	case *IndexExpression:
		node.Left, _ = Modify(node.Left, modifier).(Expression)
		node.Index, _ = Modify(node.Index, modifier).(Expression)
	case *BlockStatement:
		for i, stmt := range node.Statements {
			node.Statements[i], _ = Modify(stmt, modifier).(Statement)
		}
	case *IfExpression:
		node.Condition, _ = Modify(node.Condition, modifier).(Expression)
		node.Consequence, _ = Modify(node.Consequence, modifier).(*BlockStatement)
		node.Alternative, _ = Modify(node.Alternative, modifier).(*BlockStatement)
	case *ReturnStatement:
		node.ReturnValue, _ = Modify(node.ReturnValue, modifier).(Expression)
	case *LetStatement:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *Assignment:
		node.Value, _ = Modify(node.Value, modifier).(Expression)
	case *FunctionLiteral:
		node.Body, _ = Modify(node.Body, modifier).(*BlockStatement)
	case *ArrayLiteral:
		for i, elm := range node.Elements {
			node.Elements[i], _ = Modify(elm, modifier).(Expression)
		}
	case *HashLiteral:
		newPairs := make(map[Expression]Expression)
		for key, value := range node.Pairs {
			newKey, _ := Modify(key, modifier).(Expression)
			newValue, _ := Modify(value, modifier).(Expression)
			newPairs[newKey] = newValue
		}
		node.Pairs = newPairs
	}

	return modifier(node)
}
