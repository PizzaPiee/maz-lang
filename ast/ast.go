package ast

import (
	"bytes"
	"fmt"
	"maz-lang/token"
)

type Node interface {
	String() string
}

type Program struct {
	Statements []Node
}

func (p *Program) String() string {
	var buffer bytes.Buffer

	for _, stmt := range p.Statements {
		buffer.WriteString(stmt.String())
	}

	return buffer.String()
}

type PrefixExpression struct {
	Prefix token.Token
	Value  Node
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Prefix.Literal, pe.Value.String())
}

type InfixExpression struct {
	Left     Node
	Operator token.Token
	Right    Node
}

func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left.String(), ie.Operator.Literal, ie.Right.String())
}

type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) String() string { return fmt.Sprintf("%d", il.Value) }

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) String() string { return fmt.Sprintf("%v", bl.Value) }

type SyntaxError struct {
	Msg   string
	Token token.Token
}

func (se *SyntaxError) String() string {
	return fmt.Sprintf("Syntax error: %s\nError near: %s", se.Msg, se.Token.Literal)
}

func (se *SyntaxError) Error() string { return se.String() }

type LetStatement struct {
	Ident string
	Value Node
}

func (ls *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s;", ls.Ident, ls.Value.String())
}

type Identifier struct {
	Name string
}

func (id *Identifier) String() string { return id.Name }

type IfStatement struct {
	MainCondition  Node
	MainStatements []Node
	ElseIfs        []ElseIf
	ElseCondition  Node
	ElseStatements []Node
}

func (is *IfStatement) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("if (%s) {\n", is.MainCondition.String()))
	for _, stmt := range is.MainStatements {
		buffer.WriteString("\t" + stmt.String() + "\n")
	}
	buffer.WriteString("}")

	for _, elseIf := range is.ElseIfs {
		buffer.WriteString(elseIf.String())
	}

	if is.ElseCondition != nil {
		buffer.WriteString(" else {")
		for _, stmt := range is.ElseStatements {
			buffer.WriteString("\t" + stmt.String() + "\n")
		}
	}

	return buffer.String()
}

type ElseIf struct {
	Condition  string
	Statements []Node
}

func (ei *ElseIf) String() string {
	var buffer bytes.Buffer

	buffer.WriteString("else if (%s) {")
	for _, stmt := range ei.Statements {
		buffer.WriteString("\t" + stmt.String() + "\n")
	}
	buffer.WriteString("}")

	return buffer.String()
}
