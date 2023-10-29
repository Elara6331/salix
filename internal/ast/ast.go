package ast

import "fmt"

func PosError(n Node, filename, format string, v ...any) error {
	return fmt.Errorf(filename+":"+n.Pos().String()+": "+format, v...)
}

type Node interface {
	Pos() Position
}

type Position struct {
	Line int
	Col  int
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Col)
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
	Value    Node
	Position Position
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

type Expr struct {
	Segment  Node
	Logical  Logical
	Rest     []Expr
	Position Position
}

func (e Expr) Pos() Position {
	return e.Position
}

type ExprSegment struct {
	Value    Node
	Operator Operator
	Rest     []ExprSegment
	Position Position
}

func (es ExprSegment) Pos() Position {
	return es.Position
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

type Op interface {
	Op() string
}

type Operator struct {
	Value    string
	Position Position
}

func (op Operator) Pos() Position {
	return op.Position
}

func (op Operator) Op() string {
	return op.Value
}

type Logical struct {
	Value    string
	Position Position
}

func (l Logical) Pos() Position {
	return l.Position
}

func (l Logical) Op() string {
	return l.Value
}
