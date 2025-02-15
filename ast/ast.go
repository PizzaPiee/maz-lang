package ast

import (
	"bytes"
	"fmt"
	"maz-lang/token"
	"strings"
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
	return fmt.Sprintf("(%s%s)\n", pe.Prefix.Literal, pe.Value.String())
}

type InfixExpression struct {
	Left     Node
	Operator token.Token
	Right    Node
}

func (ie *InfixExpression) String() string {
	left := strings.TrimSpace(ie.Left.String())
	right := strings.TrimSpace(ie.Right.String())
	return fmt.Sprintf("(%s %s %s)\n", left, ie.Operator.Literal, right)
}

type IntegerLiteral struct {
	Value int64
}

func (il *IntegerLiteral) String() string { return fmt.Sprintf("%d\n", il.Value) }

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) String() string { return fmt.Sprintf("%v\n", bl.Value) }

type SyntaxError struct {
	Msg   string
	Token token.Token
}

func (se *SyntaxError) String() string {
	return fmt.Sprintf("\nSyntax error: %s\nError near: '%s'\n", se.Msg, se.Token.Literal)
}

func (se *SyntaxError) Error() string { return se.String() }

type LetStatement struct {
	Ident string
	Value Node
}

func (ls *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s;\n", ls.Ident, strings.Trim(ls.Value.String(), "\n"))
}

type Identifier struct {
	Name string
}

func (id *Identifier) String() string { return id.Name + "\n" }

type IfStatement struct {
	MainCondition  Node
	MainStatements []Node
	ElseIfs        []ElseIf
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

	if len(is.ElseStatements) > 0 {
		buffer.WriteString(" else {\n")
		for _, stmt := range is.ElseStatements {
			buffer.WriteString("\t" + stmt.String() + "\n")
		}
		buffer.WriteString("}\n")
	}

	return buffer.String()
}

type ElseIf struct {
	Condition  Node
	Statements []Node
}

func (ei *ElseIf) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf(" else if (%s) {\n", ei.Condition.String()))
	for _, stmt := range ei.Statements {
		buffer.WriteString("\t" + stmt.String() + "\n")
	}
	buffer.WriteString("}")

	return buffer.String()
}

type ReturnStatement struct {
	Expression Node
}

func (rs *ReturnStatement) String() string {
	return fmt.Sprintf("return %s;\n", strings.TrimSpace(rs.Expression.String()))
}

type Function struct {
	Name       string
	Parameters []Node
	Body       []Node
}

func (f *Function) String() string {
	var out bytes.Buffer

	out.WriteString("fn ")
	if f.Name != "" {
		out.WriteString(f.Name)
	}

	var params []string
	for _, p := range f.Parameters {
		params = append(params, strings.TrimSpace(p.String()))
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") ")

	out.WriteString("{\n")
	for _, stmt := range f.Body {
		out.WriteString("\t" + stmt.String())
	}
	out.WriteString("}\n")

	return out.String()
}
