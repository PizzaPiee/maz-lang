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
