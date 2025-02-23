package object

import (
	"fmt"
	"maz-lang/ast"
)

type ObjectType string

const (
	INTEGER_OBJ = "INT"
	BOOLEAN_OBJ = "BOOL"
	NULL_OBJ    = "NULL"
	ERROR_OBJ   = "ERROR"
	FUNCDEF_OBJ = "FUNCDEF"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }

type Error struct {
	Value error
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return e.Value.Error() }

type FunctionDef struct {
	Fn ast.FunctionDefinition
}

func (f *FunctionDef) Type() ObjectType { return FUNCDEF_OBJ }
func (f *FunctionDef) Inspect() string  { return fmt.Sprintf("<fn %s>", f.Fn.Name) }
