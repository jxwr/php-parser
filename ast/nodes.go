package ast

// NewVariable intializes a variable node with its name being a simple
// identifier and its type set to AnyType. The name argument should not
// include the $ operator.
func NewVariable(name string) *Variable {
	return &Variable{Name: &Identifier{Value: name}, Type: AnyType}
}

// Echo returns a new echo statement.
func Echo(exprs ...Expression) *EchoStmt {
	return &EchoStmt{Expressions: exprs}
}

func NewClassExpression(r string, e Expression) *ClassExpression {
	return &ClassExpression{
		Receiver:   &Identifier{Value: r},
		Expression: e,
	}
}
