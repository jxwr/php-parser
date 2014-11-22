package ast

import (
	"bytes"
	"fmt"
	"strings"
)

// Node encapsulates every AST node.
type Node interface {
	String() string
}

// An Identifier is a raw string that can be used to identify
// a variable, function, class, constant, property, etc.
type Identifier struct {
	Parent Node
	Value  string
}

func (i Identifier) EvaluatesTo() Type {
	return String
}

func (i Identifier) String() string {
	return i.Value
}

type Variable struct {

	// Name is the identifier for the variable, which may be
	// a dynamic expression.
	Name Expression
	Type Type
}

func (v Variable) String() string {
	return "$" + v.Name.String()
}

type GlobalDeclaration struct {
	Identifiers []*Variable
}

func (g GlobalDeclaration) String() string {
	return "global"
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

// A statement is an executable piece of code. It may be as simple as
// a function call or a variable assignment. It also includes things like
// "if".
type Statement interface {
	Node
}

// EmptyStatement represents a statement that does nothing.
type EmptyStatement struct {
}

func (e EmptyStatement) String() string { return "" }

// An Expression is a snippet of code that evaluates to a single value when run
// and does not represent a program instruction.
type Expression interface {
	Node
	EvaluatesTo() Type
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

func (b BinaryExpression) String() string {
	return b.Operator
}

func (b BinaryExpression) EvaluatesTo() Type {
	return b.Type
}

type TernaryExpression struct {
	Condition, True, False Expression
	Type                   Type
}

func (t TernaryExpression) String() string {
	return "?:"
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

func (u UnaryExpression) String() string {
	if u.Preceding {
		return u.Operator + u.Operand.String()
	}
	return u.Operand.String() + u.Operator
}

func (u UnaryExpression) EvaluatesTo() Type {
	return Unknown
}

type ExpressionStmt struct {
	Expression
}

func (e ExpressionStmt) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
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

func (e EchoStmt) String() string {
	return "Echo"
}

// ReturnStmt represents a function return.
type ReturnStmt struct {
	Expression
}

func (r ReturnStmt) String() string {
	return fmt.Sprintf("return")
}

type BreakStmt struct {
	Expression
}

func (b BreakStmt) String() string {
	return "break"
}

type ContinueStmt struct {
	Expression
}

func (c ContinueStmt) String() string {
	return "continue"
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

func (i Include) String() string {
	return "include"
}

func (i Include) EvaluatesTo() Type {
	return AnyType
}

type ExitStmt struct {
	Expression Expression
}

func (e ExitStmt) String() string {
	return "exit"
}

type NewExpression struct {
	Class     Expression
	Arguments []Expression
}

func (n NewExpression) EvaluatesTo() Type {
	return Object
}

func (n NewExpression) String() string {
	return "new"
}

type AssignmentExpression struct {
	Assignee Assignable
	Value    Expression
	Operator string
}

func (a AssignmentExpression) EvaluatesTo() Type {
	return a.Value.EvaluatesTo()
}

func (a AssignmentExpression) String() string {
	return a.Operator
}

type Assignable interface {
	Node
	AssignableType() Type
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

func (f FunctionCallExpression) String() string {
	return fmt.Sprintf("%s()", f.FunctionName)
}

type Block struct {
	Statements []Statement
	Scope      Scope
}

func (b Block) String() string {
	return "{}"
}

type FunctionStmt struct {
	*FunctionDefinition
	Body *Block
}

func (f FunctionStmt) String() string {
	return fmt.Sprintf("Func: %s", f.Name)
}

type AnonymousFunction struct {
	ClosureVariables []FunctionArgument
	Arguments        []FunctionArgument
	Body             *Block
}

func (a AnonymousFunction) EvaluatesTo() Type {
	return Function
}

func (a AnonymousFunction) String() string {
	return "anonymous function"
}

type FunctionDefinition struct {
	Name      string
	Arguments []FunctionArgument
}

func (fd FunctionDefinition) String() string {
	return fmt.Sprintf("function %s( %s )", fd.Name, fd.Arguments)
}

type FunctionArgument struct {
	TypeHint string
	Default  Expression
	Variable *Variable
}

func (fa FunctionArgument) String() string {
	return fmt.Sprintf("Arg: %s", fa.TypeHint)
}

type Class struct {
	Name       string
	Extends    string
	Implements []string
	Methods    []Method
	Properties []Property
	Constants  []Constant
}

func (c Class) String() string {
	return fmt.Sprintf("class %s", c.Name)
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

func (i Interface) String() string {
	return fmt.Sprintf("interface %s extends %s", i.Name, strings.Join(i.Inherits, ", "))
}

type Property struct {
	Name           string
	Visibility     Visibility
	Type           Type
	Initialization Expression
}

func (p Property) String() string {
	return fmt.Sprintf("Prop: %s", p.Name)
}

func (p Property) AssignableType() Type {
	return p.Type
}

type PropertyExpression struct {
	Receiver Expression
	Name     Expression
	Type     Type
}

func (p PropertyExpression) String() string {
	return fmt.Sprintf("%s->%s", p.Receiver, p.Name)
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

func (c ClassExpression) String() string {
	return fmt.Sprintf("%s::", c.Receiver)
}

func (c ClassExpression) AssignableType() Type {
	return c.Type
}

type Method struct {
	*FunctionStmt
	Visibility Visibility
}

func (m Method) String() string {
	return m.Name
}

type MethodCallExpression struct {
	Receiver Expression
	*FunctionCallExpression
}

func (m MethodCallExpression) String() string {
	return fmt.Sprintf("%s->", m.Receiver)
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

func (i IfStmt) String() string {
	return "if"
}

type SwitchStmt struct {
	Expression  Expression
	Cases       []*SwitchCase
	DefaultCase *Block
}

func (s SwitchStmt) String() string {
	return "switch"
}

type SwitchCase struct {
	Expression Expression
	Block      Block
}

func (s SwitchCase) String() string {
	return "case"
}

type ForStmt struct {
	Initialization []Expression
	Termination    []Expression
	Iteration      []Expression
	LoopBlock      Statement
}

func (f ForStmt) String() string {
	return "for"
}

type WhileStmt struct {
	Termination Expression
	LoopBlock   Statement
}

func (w WhileStmt) String() string {
	return fmt.Sprintf("while")
}

type DoWhileStmt struct {
	Termination Expression
	LoopBlock   Statement
}

func (d DoWhileStmt) String() string {
	return fmt.Sprintf("do ... while")
}

type TryStmt struct {
	TryBlock     *Block
	FinallyBlock *Block
	CatchStmts   []*CatchStmt
}

func (t TryStmt) String() string {
	return "try"
}

type CatchStmt struct {
	CatchBlock *Block
	CatchType  string
	CatchVar   *Variable
}

func (c CatchStmt) String() string {
	return fmt.Sprintf("catch %s %s", c.CatchType, c.CatchVar)
}

type Literal struct {
	Type  Type
	Value string
}

func (l Literal) String() string {
	return fmt.Sprintf("Literal-%s: %s", l.Type, l.Value)
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

func (f ForeachStmt) String() string {
	return "foreach"
}

type ArrayExpression struct {
	ArrayType
	Pairs []ArrayPair
}

func (a ArrayExpression) String() string {
	return "array"
}

type ArrayPair struct {
	Key   Expression
	Value Expression
}

func (p ArrayPair) String() string {
	return fmt.Sprintf("%s => %s", p.Key, p.Value)
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

func (a ArrayLookupExpression) String() string {
	return fmt.Sprintf("%s[", a.Array)
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

func (a ArrayAppendExpression) String() string {
	return a.Array.String() + "[]"
}

type ShellCommand struct {
	Command string
}

func (s ShellCommand) String() string {
	return fmt.Sprintf("`%s`", s.Command)
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

func (l ListStatement) String() string {
	return fmt.Sprintf("list(%s)", l.Assignees)
}

type StaticVariableDeclaration struct {
	Declarations []Expression
}

func (s StaticVariableDeclaration) String() string {
	buf := bytes.NewBufferString("static ")
	for i, d := range s.Declarations {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(d.String())
	}
	return buf.String()
}

type DeclareBlock struct {
	Statements   *Block
	Declarations []string
}

func (d DeclareBlock) String() string {
	return "declare{}"
}
