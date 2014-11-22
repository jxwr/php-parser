package parser

import (
	"github.com/jxwr/php-parser/ast"
	"github.com/jxwr/php-parser/token"
)

func (p *Parser) parseIf() *ast.IfStmt {
	p.expect(token.OpenParen)
	n := &ast.IfStmt{}
	n.Condition = p.parseNextExpression()
	p.expect(token.CloseParen)

	p.next()
	n.TrueBranch = p.parseControlBlock(token.EndIf, token.ElseIf, token.Else)
	n.FalseBranch = &ast.Block{}

	blockStyle := false
	switch p.current.Typ {
	case token.ElseIf, token.Else, token.EndIf:
	default:
		p.next()
		blockStyle = true
	}

	switch p.current.Typ {
	case token.ElseIf:
		n.FalseBranch = p.parseIf()
	case token.Else:
		p.next()
		if p.current.Typ == token.If {
			n.FalseBranch = p.parseIf()
		} else {
			n.FalseBranch = p.parseControlBlock(token.EndIf)
		}
	default:
		if blockStyle {
			p.backup()
		}
	}

	return n
}

func (p *Parser) parseWhile() ast.Statement {
	p.expect(token.OpenParen)
	term := p.parseNextExpression()
	p.expect(token.CloseParen)
	p.next()
	block := p.parseControlBlock(token.EndWhile)
	return &ast.WhileStmt{
		Termination: term,
		LoopBlock:   block,
	}
}

func (p *Parser) parseForeach() ast.Statement {
	stmt := &ast.ForeachStmt{}
	p.expect(token.OpenParen)
	stmt.Source = p.parseNextExpression()
	p.expect(token.AsOperator)
	if p.peek().Typ == token.AmpersandOperator {
		p.expect(token.AmpersandOperator)
	}
	p.expect(token.VariableOperator)
	p.next()
	first := ast.NewVariable(p.current.Val)
	if p.peek().Typ == token.ArrayKeyOperator {
		stmt.Key = first
		p.expect(token.ArrayKeyOperator)
		if p.peek().Typ == token.AmpersandOperator {
			p.expect(token.AmpersandOperator)
		}
		p.expect(token.VariableOperator)
		p.next()
		stmt.Value = ast.NewVariable(p.current.Val)
	} else {
		stmt.Value = first
	}
	p.expect(token.CloseParen)
	p.next()
	stmt.LoopBlock = p.parseControlBlock(token.EndForeach)
	return stmt
}

func (p *Parser) parseControlBlock(end ...token.Token) ast.Statement {
	// try to parse this in bash style, but it requires an end token
	if len(end) > 0 && p.current.Typ == token.TernaryOperator2 {
		return p.parseStatementsUntil(end...)
	}
	return p.parseStmt()
}

func (p *Parser) parseFor() ast.Statement {
	stmt := &ast.ForStmt{}
	p.expect(token.OpenParen)
	stmt.Initialization = p.parseExpressionsUntil(token.Comma, token.StatementEnd)
	stmt.Termination = p.parseExpressionsUntil(token.Comma, token.StatementEnd)
	stmt.Iteration = p.parseExpressionsUntil(token.Comma, token.CloseParen)
	p.expectCurrent(token.CloseParen)
	p.next()
	stmt.LoopBlock = p.parseControlBlock(token.EndFor)
	return stmt
}

func (p *Parser) parseDo() ast.Statement {
	block := p.parseBlock()
	p.expect(token.While)
	p.expect(token.OpenParen)
	term := p.parseNextExpression()
	p.expect(token.CloseParen)
	p.expectStmtEnd()
	return &ast.DoWhileStmt{
		Termination: term,
		LoopBlock:   block,
	}
}

func (p *Parser) parseSwitch() ast.Statement {
	stmt := ast.SwitchStmt{}
	p.expect(token.OpenParen)
	stmt.Expression = p.parseExpression()
	p.expectCurrent(token.CloseParen)
	p.expect(token.BlockBegin, token.TernaryOperator2)
	p.next()
	for {
		switch p.current.Typ {
		case token.Case:
			expr := p.parseNextExpression()
			p.expect(token.TernaryOperator2, token.StatementEnd)
			p.next()
			stmt.Cases = append(stmt.Cases, &ast.SwitchCase{
				Expression: expr,
				Block:      *(p.parseSwitchBlock()),
			})
		case token.Default:
			p.expect(token.TernaryOperator2, token.StatementEnd)
			p.next()
			stmt.DefaultCase = p.parseSwitchBlock()
		case token.BlockEnd, token.EndSwitch:
			return &stmt
		default:
			p.errorf("Unexpected token. in switch statement:", p.current)
			return nil
		}
	}
}

func (p *Parser) parseSwitchBlock() *ast.Block {
	needBlockEnd := false
	if p.current.Typ == token.BlockBegin {
		needBlockEnd = true
		p.next()
	}
	block := &ast.Block{
		Statements: make([]ast.Statement, 0),
	}
stmtLoop:
	for {
		switch p.current.Typ {
		case token.BlockEnd:
			if needBlockEnd {
				needBlockEnd = false
				p.next()
			}
			fallthrough
		case token.Case, token.Default, token.EndSwitch:
			break stmtLoop
		default:
			stmt := p.parseStmt()
			if stmt == nil {
				p.errorf("Invalid statement in switch block", p.current)
				break stmtLoop
			}
			block.Statements = append(block.Statements, stmt)
			p.next()
		}
	}
	if needBlockEnd {
		p.errorf("switch case needs block end")
	}
	return block
}

func (p *Parser) parseDeclareBlock() *ast.DeclareBlock {
	declare := &ast.DeclareBlock{Declarations: make([]string, 0)}

	p.expectCurrent(token.Declare)
	p.expect(token.OpenParen)

	declare.Declarations = append(declare.Declarations, p.parseDeclareElement())

	p.next()
	for p.current.Typ == token.Comma {
		declare.Declarations = append(declare.Declarations, p.parseDeclareElement())
		p.next()
	}

	p.expectCurrent(token.CloseParen)

	if p.peek().Typ == token.BlockBegin {
		declare.Statements = p.parseBlock()
	} else {
		p.expect(token.StatementEnd)
	}
	return declare
}

func (p *Parser) parseDeclareElement() string {
	element := ""
	p.expect(token.Identifier)
	element += p.current.Val

	p.expect(token.AssignmentOperator)
	element += p.current.Val

	p.parseNextExpression()
	element += p.current.Val
	return element
}
