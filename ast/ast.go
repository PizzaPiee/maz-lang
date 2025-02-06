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
	Value  token.Token
}

func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Prefix.Literal, pe.Value.Literal)
}
