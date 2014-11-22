package ast

// An Identifier is a raw string that can be used to identify
// a variable, function, class, constant, property, etc.
type Identifier struct {
	Parent Node
	Value  string
}

func (i Identifier) EvaluatesTo() Type {
	return String
}

type Variable struct {

	// Name is the identifier for the variable, which may be
	// a dynamic expression.
	Name Expression
	Type Type
}

type GlobalDeclaration struct {
	Identifiers []*Variable
}

func (i Variable) AssignableType() Type {
	return i.Type
}

// EvaluatesTo returns the known type of the variable.
func (i Variable) EvaluatesTo() Type {
	return i.Type
}

// NewVariable intializes a variable node with its name being a simple
// identifier and its type set to AnyType. The name argument should not
// include the $ operator.
func NewVariable(name string) *Variable {
	return &Variable{Name: Identifier{Value: name}, Type: AnyType}
}

// EmptyStatement represents a statement that does nothing.
type EmptyStatement struct {
}

// AnyType is a bitmask of all the valid types.
const AnyType = String | Integer | Float | Boolean | Null | Resource | Array | Object

// OperatorExpression is an expression that applies an operator to one, two, or three
// operands. The operator determines how many operands it should contain.
type BinaryExpression struct {
	Antecedent Expression
	Subsequent Expression
	Type       Type
	Operator   string
}

func (b BinaryExpression) EvaluatesTo() Type {
	return b.Type
}

type TernaryExpression struct {
	Condition, True, False Expression
	Type                   Type
}

func (t TernaryExpression) EvaluatesTo() Type {
	return t.Type
}

// UnaryExpression is an expression that applies an operator to only one operand. The
// operator may precede or follow the operand.
type UnaryExpression struct {
	Operand   Expression
	Operator  string
	Preceding bool
}

func (u UnaryExpression) EvaluatesTo() Type {
	return Unknown
}

type ExpressionStmt struct {
	Expression
}

// Echo returns a new echo statement.
func Echo(exprs ...Expression) EchoStmt {
	return EchoStmt{Expressions: exprs}
}

// Echo represents an echo statement. It may be either a literal statement
// or it may be from data outside PHP-mode, such as "here" in: <? not here ?> here <? not here ?>
type EchoStmt struct {
	Expressions []Expression
}

// ReturnStmt represents a function return.
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

func (i Include) EvaluatesTo() Type {
	return AnyType
}

type ExitStmt struct {
	Expression Expression
}

type NewExpression struct {
	Class     Expression
	Arguments []Expression
}

func (n NewExpression) EvaluatesTo() Type {
	return Object
}

type AssignmentExpression struct {
	Assignee Assignable
	Value    Expression
	Operator string
}

func (a AssignmentExpression) EvaluatesTo() Type {
	return a.Value.EvaluatesTo()
}

type FunctionCallStmt struct {
	FunctionCallExpression
}

type FunctionCallExpression struct {
	FunctionName Expression
	Arguments    []Expression
}

func (f FunctionCallExpression) EvaluatesTo() Type {
	return String | Integer | Float | Boolean | Null | Resource | Array | Object
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

func (a AnonymousFunction) EvaluatesTo() Type {
	return Function
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

type ConstantExpression struct {
	*Variable
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

func (p Property) AssignableType() Type {
	return p.Type
}

type PropertyExpression struct {
	Receiver Expression
	Name     Expression
	Type     Type
}

func (p PropertyExpression) AssignableType() Type {
	return p.Type
}

func (p PropertyExpression) EvaluatesTo() Type {
	return AnyType
}

type ClassExpression struct {
	Receiver   Expression
	Expression Expression
	Type       Type
}

func NewClassExpression(r string, e Expression) *ClassExpression {
	return &ClassExpression{
		Receiver:   Identifier{Value: r},
		Expression: e,
	}
}

func (c ClassExpression) EvaluatesTo() Type {
	return AnyType
}

func (c ClassExpression) AssignableType() Type {
	return c.Type
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

func (l Literal) EvaluatesTo() Type {
	return l.Type
}

type ForeachStmt struct {
	Source    Expression
	Key       *Variable
	Value     *Variable
	LoopBlock Statement
}

type ArrayExpression struct {
	ArrayType
	Pairs []ArrayPair
}

type ArrayPair struct {
	Key   Expression
	Value Expression
}

func (a ArrayExpression) EvaluatesTo() Type {
	return Array
}

func (a ArrayExpression) AssignableType() Type {
	return AnyType
}

type ArrayLookupExpression struct {
	Array Expression
	Index Expression
}

func (a ArrayLookupExpression) EvaluatesTo() Type {
	return AnyType
}

func (a ArrayLookupExpression) AssignableType() Type {
	return AnyType
}

type ArrayAppendExpression struct {
	Array Expression
}

func (a ArrayAppendExpression) EvaluatesTo() Type {
	return AnyType
}

func (a ArrayAppendExpression) AssignableType() Type {
	return AnyType
}

type ShellCommand struct {
	Command string
}

func (s ShellCommand) EvaluatesTo() Type {
	return String
}

type ListStatement struct {
	Assignees []Assignable
	Value     Expression
	Operator  string
}

func (l ListStatement) EvaluatesTo() Type {
	return Array
}

type StaticVariableDeclaration struct {
	Declarations []Expression
}

type DeclareBlock struct {
	Statements   *Block
	Declarations []string
}
