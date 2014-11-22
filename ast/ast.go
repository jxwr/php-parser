package ast

/// Interfaces

type Node interface {
	Accept(Visitor)
}

type Expression interface {
	Node
	exprNode()
}

type Statement interface {
	Node
	stmtNode()
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

type Literal struct {
	Type  Type
	Value string
}

type ShellCommand struct {
	Command string
}

type Include struct {
	Expressions []Expression
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

type AnonymousFunction struct {
	ClosureVariables []FunctionArgument
	Arguments        []FunctionArgument
	Body             *Block
}

func (n Identifier) exprNode()             {}
func (n Variable) exprNode()               {}
func (n BinaryExpression) exprNode()       {}
func (n TernaryExpression) exprNode()      {}
func (n UnaryExpression) exprNode()        {}
func (n NewExpression) exprNode()          {}
func (n PropertyExpression) exprNode()     {}
func (n ClassExpression) exprNode()        {}
func (n AssignmentExpression) exprNode()   {}
func (n FunctionCallExpression) exprNode() {}
func (n ConstantExpression) exprNode()     {}
func (n ArrayExpression) exprNode()        {}
func (n ArrayLookupExpression) exprNode()  {}
func (n ArrayAppendExpression) exprNode()  {}
func (n ShellCommand) exprNode()           {}
func (n Literal) exprNode()                {}
func (n Include) exprNode()                {}
func (n AnonymousFunction) exprNode()      {}

func (n *Identifier) Accept(v Visitor)             { v.VisitIdentifier(n) }
func (n *Variable) Accept(v Visitor)               { v.VisitVariable(n) }
func (n *BinaryExpression) Accept(v Visitor)       { v.VisitBinaryExpression(n) }
func (n *TernaryExpression) Accept(v Visitor)      { v.VisitTernaryExpression(n) }
func (n *UnaryExpression) Accept(v Visitor)        { v.VisitUnaryExpression(n) }
func (n *NewExpression) Accept(v Visitor)          { v.VisitNewExpression(n) }
func (n *PropertyExpression) Accept(v Visitor)     { v.VisitPropertyExpression(n) }
func (n *ClassExpression) Accept(v Visitor)        { v.VisitClassExpression(n) }
func (n *AssignmentExpression) Accept(v Visitor)   { v.VisitAssignmentExpression(n) }
func (n *FunctionCallExpression) Accept(v Visitor) { v.VisitFunctionCallExpression(n) }
func (n *ConstantExpression) Accept(v Visitor)     { v.VisitConstantExpression(n) }
func (n *ArrayExpression) Accept(v Visitor)        { v.VisitArrayExpression(n) }
func (n *ArrayLookupExpression) Accept(v Visitor)  { v.VisitArrayLookupExpression(n) }
func (n *ArrayAppendExpression) Accept(v Visitor)  { v.VisitArrayAppendExpression(n) }
func (n *ShellCommand) Accept(v Visitor)           { v.VisitShellCommand(n) }
func (n *Literal) Accept(v Visitor)                { v.VisitLiteral(n) }
func (n *Include) Accept(v Visitor)                { v.VisitInclude(n) }
func (n *AnonymousFunction) Accept(v Visitor)      { v.VisitAnonymousFunction(n) }

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

type ForeachStmt struct {
	Source    Expression
	Key       *Variable
	Value     *Variable
	LoopBlock Statement
}

// list($a, $b, $c) = $my_array;
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

func (n GlobalDeclaration) stmtNode()         {}
func (n ExpressionStmt) stmtNode()            {}
func (n EmptyStatement) stmtNode()            {}
func (n EchoStmt) stmtNode()                  {}
func (n ReturnStmt) stmtNode()                {}
func (n BreakStmt) stmtNode()                 {}
func (n ContinueStmt) stmtNode()              {}
func (n ThrowStmt) stmtNode()                 {}
func (n IncludeStmt) stmtNode()               {}
func (n ExitStmt) stmtNode()                  {}
func (n FunctionCallStmt) stmtNode()          {}
func (n FunctionStmt) stmtNode()              {}
func (n FunctionDefinition) stmtNode()        {}
func (n Interface) stmtNode()                 {}
func (n DeclareBlock) stmtNode()              {}
func (n Class) stmtNode()                     {}
func (n Method) stmtNode()                    {}
func (n Block) stmtNode()                     {}
func (n IfStmt) stmtNode()                    {}
func (n SwitchStmt) stmtNode()                {}
func (n ForStmt) stmtNode()                   {}
func (n WhileStmt) stmtNode()                 {}
func (n DoWhileStmt) stmtNode()               {}
func (n TryStmt) stmtNode()                   {}
func (n CatchStmt) stmtNode()                 {}
func (n ForeachStmt) stmtNode()               {}
func (n ListStatement) stmtNode()             {}
func (n StaticVariableDeclaration) stmtNode() {}

func (n *GlobalDeclaration) Accept(v Visitor)  { v.VisitGlobalDeclaration(n) }
func (n *ExpressionStmt) Accept(v Visitor)     { v.VisitExpressionStmt(n) }
func (n *EmptyStatement) Accept(v Visitor)     { v.VisitEmptyStatement(n) }
func (n *EchoStmt) Accept(v Visitor)           { v.VisitEchoStmt(n) }
func (n *ReturnStmt) Accept(v Visitor)         { v.VisitReturnStmt(n) }
func (n *BreakStmt) Accept(v Visitor)          { v.VisitBreakStmt(n) }
func (n *ContinueStmt) Accept(v Visitor)       { v.VisitContinueStmt(n) }
func (n *ThrowStmt) Accept(v Visitor)          { v.VisitThrowStmt(n) }
func (n *IncludeStmt) Accept(v Visitor)        { v.VisitIncludeStmt(n) }
func (n *ExitStmt) Accept(v Visitor)           { v.VisitExitStmt(n) }
func (n *FunctionCallStmt) Accept(v Visitor)   { v.VisitFunctionCallStmt(n) }
func (n *FunctionStmt) Accept(v Visitor)       { v.VisitFunctionStmt(n) }
func (n *FunctionDefinition) Accept(v Visitor) { v.VisitFunctionDefinition(n) }
func (n *Interface) Accept(v Visitor)          { v.VisitInterface(n) }
func (n *DeclareBlock) Accept(v Visitor)       { v.VisitDeclareBlock(n) }
func (n *Class) Accept(v Visitor)              { v.VisitClass(n) }
func (n *Method) Accept(v Visitor)             { v.VisitMethod(n) }
func (n *Block) Accept(v Visitor)              { v.VisitBlock(n) }
func (n *IfStmt) Accept(v Visitor)             { v.VisitIfStmt(n) }
func (n *SwitchStmt) Accept(v Visitor)         { v.VisitSwitchStmt(n) }
func (n *ForStmt) Accept(v Visitor)            { v.VisitForStmt(n) }
func (n *WhileStmt) Accept(v Visitor)          { v.VisitWhileStmt(n) }
func (n *DoWhileStmt) Accept(v Visitor)        { v.VisitDoWhileStmt(n) }
func (n *TryStmt) Accept(v Visitor)            { v.VisitTryStmt(n) }
func (n *CatchStmt) Accept(v Visitor)          { v.VisitCatchStmt(n) }
func (n *ForeachStmt) Accept(v Visitor)        { v.VisitForeachStmt(n) }
func (n *ListStatement) Accept(v Visitor)      { v.VisitListStatement(n) }
func (n *StaticVariableDeclaration) Accept(v Visitor) {
	v.VisitStaticVariableDeclaration(n)
}
