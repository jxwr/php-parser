package ast

/// Interfaces

type Node interface {
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Assignable interface {
	Node
}

// AnyType is a bitmask of all the valid types.
const AnyType = String | Integer | Float | Boolean | Null | Resource | Array | Object

/// Expression

type Identifier struct {
	Parent Node
	Value  string
}

type Variable struct {
	// Name is the identifier for the variable, which may be
	// a dynamic expression.
	Name Expression
	Type Type
}

type BinaryExpression struct {
	Antecedent Expression
	Subsequent Expression
	Type       Type
	Operator   string
}

type TernaryExpression struct {
	Condition, True, False Expression
	Type                   Type
}

type UnaryExpression struct {
	Operand   Expression
	Operator  string
	Preceding bool
}

type NewExpression struct {
	Class     Expression
	Arguments []Expression
}

type AssignmentExpression struct {
	Assignee Assignable
	Value    Expression
	Operator string
}

type FunctionCallExpression struct {
	FunctionName Expression
	Arguments    []Expression
}

type ConstantExpression struct {
	*Variable
}

type ArrayExpression struct {
	ArrayType
	Pairs []ArrayPair
}

type ArrayPair struct {
	Key   Expression
	Value Expression
}

type ArrayLookupExpression struct {
	Array Expression
	Index Expression
}

type ArrayAppendExpression struct {
	Array Expression
}

/// Statements

type GlobalDeclaration struct {
	Identifiers []*Variable
}

type EmptyStatement struct {
}

type ExpressionStmt struct {
	Expression
}

type EchoStmt struct {
	Expressions []Expression
}

type ReturnStmt struct {
	Expression
}

type BreakStmt struct {
	Expression
}

type ContinueStmt struct {
	Expression
}

type ThrowStmt struct {
	Expression
}

type IncludeStmt struct {
	Include
}

type Include struct {
	Expressions []Expression
}

type ExitStmt struct {
	Expression Expression
}

type FunctionCallStmt struct {
	FunctionCallExpression
}

type Block struct {
	Statements []Statement
	Scope      Scope
}

type FunctionStmt struct {
	*FunctionDefinition
	Body *Block
}

type AnonymousFunction struct {
	ClosureVariables []FunctionArgument
	Arguments        []FunctionArgument
	Body             *Block
}

type FunctionDefinition struct {
	Name      string
	Arguments []FunctionArgument
}

type FunctionArgument struct {
	TypeHint string
	Default  Expression
	Variable *Variable
}

type Class struct {
	Name       string
	Extends    string
	Implements []string
	Methods    []Method
	Properties []Property
	Constants  []Constant
}

type Constant struct {
	*Variable
	Value interface{}
}

type Interface struct {
	Name      string
	Inherits  []string
	Methods   []Method
	Constants []Constant
}

type Property struct {
	Name           string
	Visibility     Visibility
	Type           Type
	Initialization Expression
}

type PropertyExpression struct {
	Receiver Expression
	Name     Expression
	Type     Type
}

type ClassExpression struct {
	Receiver   Expression
	Expression Expression
	Type       Type
}

type Method struct {
	*FunctionStmt
	Visibility Visibility
}

type MethodCallExpression struct {
	Receiver Expression
	*FunctionCallExpression
}

type Visibility int

const (
	Private Visibility = iota
	Protected
	Public
)

type IfStmt struct {
	Condition   Expression
	TrueBranch  Statement
	FalseBranch Statement
}

type SwitchStmt struct {
	Expression  Expression
	Cases       []*SwitchCase
	DefaultCase *Block
}

type SwitchCase struct {
	Expression Expression
	Block      Block
}

type ForStmt struct {
	Initialization []Expression
	Termination    []Expression
	Iteration      []Expression
	LoopBlock      Statement
}

type WhileStmt struct {
	Termination Expression
	LoopBlock   Statement
}

type DoWhileStmt struct {
	Termination Expression
	LoopBlock   Statement
}

type TryStmt struct {
	TryBlock     *Block
	FinallyBlock *Block
	CatchStmts   []*CatchStmt
}

type CatchStmt struct {
	CatchBlock *Block
	CatchType  string
	CatchVar   *Variable
}

type Literal struct {
	Type  Type
	Value string
}

type ForeachStmt struct {
	Source    Expression
	Key       *Variable
	Value     *Variable
	LoopBlock Statement
}

type ShellCommand struct {
	Command string
}

type ListStatement struct {
	Assignees []Assignable
	Value     Expression
	Operator  string
}

type StaticVariableDeclaration struct {
	Declarations []Expression
}

type DeclareBlock struct {
	Statements   *Block
	Declarations []string
}
