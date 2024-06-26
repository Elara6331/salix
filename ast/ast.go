package ast

import "fmt"

func PosError(n Node, format string, v ...any) error {
	return fmt.Errorf(n.Pos().String()+": "+format, v...)
}

type Node interface {
	Pos() Position
}

type Position struct {
	Name string
	Line int
	Col  int
}

func (p Position) String() string {
	return fmt.Sprintf("%s: line %d, col %d", p.Name, p.Line, p.Col)
}

type Nil struct {
	Position Position
}

func (n Nil) Pos() Position {
	return n.Position
}

type Tag struct {
	Name     Ident
	Params   []Node
	HasBody  bool
	Position Position
}

func (t Tag) Pos() Position {
	return t.Position
}

type ExprTag struct {
	Value       Node
	IgnoreError bool
	Position    Position
}

func (et ExprTag) Pos() Position {
	return et.Position
}

type EndTag struct {
	Name     Ident
	Position Position
}

func (et EndTag) Pos() Position {
	return et.Position
}

type Text struct {
	Data     []byte
	Position Position
}

func (t Text) Pos() Position {
	return t.Position
}

type Value struct {
	Node
	Not bool
}

type Map struct {
	Map      map[Node]Node
	Position Position
}

func (m Map) Pos() Position {
	return m.Position
}

type Array struct {
	Array    []Node
	Position Position
}

func (a Array) Pos() Position {
	return a.Position
}

type Expr struct {
	First    Node
	Operator Operator
	Rest     []Expr
	Position Position
}

func (e Expr) Pos() Position {
	return e.Position
}

type Assignment struct {
	Name     Ident
	Value    Node
	Position Position
}

func (a Assignment) Pos() Position {
	return a.Position
}

type FuncCall struct {
	Name     Ident
	Params   []Node
	Position Position
}

func (fc FuncCall) Pos() Position {
	return fc.Position
}

type MethodCall struct {
	Value    Node
	Name     Ident
	Params   []Node
	Position Position
}

func (mc MethodCall) Pos() Position {
	return mc.Position
}

type FieldAccess struct {
	Value    Node
	Name     Ident
	Position Position
}

func (fa FieldAccess) Pos() Position {
	return fa.Position
}

type Index struct {
	Value    Node
	Index    Node
	Position Position
}

func (i Index) Pos() Position {
	return i.Position
}

type Ident struct {
	Value    string
	Position Position
}

func (id Ident) Pos() Position {
	return id.Position
}

type String struct {
	Value    string
	Position Position
}

func (s String) Pos() Position {
	return s.Position
}

type Float struct {
	Value    float64
	Position Position
}

func (f Float) Pos() Position {
	return f.Position
}

type Integer struct {
	Value    int64
	Position Position
}

func (i Integer) Pos() Position {
	return i.Position
}

type Bool struct {
	Value    bool
	Position Position
}

func (b Bool) Pos() Position {
	return b.Position
}

type Operator struct {
	Value    string
	Position Position
}

func (op Operator) Pos() Position {
	return op.Position
}

type Ternary struct {
	Condition Node
	IfTrue    Node
	Else      Node
	Position  Position
}

func (t Ternary) Pos() Position {
	return t.Position
}

type VariableOr struct {
	Variable Ident
	Or       Node
	Position Position
}

func (vo VariableOr) Pos() Position {
	return vo.Position
}
