// DaVinci Script

package parser

import (
	"fmt"
	. "github.com/DavinciScript/Davi/lexer"
	"strings"
)

type Program struct {
	Statements Block
}

func (p *Program) String() string {
	return p.Statements.String()
}

type Block []Statement

func (b Block) String() string {
	lines := []string{}
	for _, s := range b {
		lines = append(lines, fmt.Sprintf("%s", s))
	}
	return strings.Join(lines, "\n")
}

type Statement interface {
	Position() Position
	statementNode()
}

type Assign struct {
	pos    Position
	Target Expression
	Value  Expression
}

func (s *Assign) statementNode()     {}
func (s *Assign) Position() Position { return s.pos }

func (s *Assign) String() string {
	return fmt.Sprintf("%s = %s", s.Target, s.Value)
}

type OuterAssign struct {
	pos   Position
	Name  string
	Value Expression
}

func (s *OuterAssign) statementNode()     {}
func (s *OuterAssign) Position() Position { return s.pos }

func (s *OuterAssign) String() string {
	return fmt.Sprintf("outer %s = %s", s.Name, s.Value)
}

type If struct {
	pos       Position
	Condition Expression
	Body      Block
	Else      Block
}

func (s *If) statementNode()     {}
func (s *If) Position() Position { return s.pos }

func indent(s string) string {
	input := strings.Split(s, "\n")
	output := []string{}
	for _, line := range input {
		output = append(output, "    "+line)
	}
	return strings.Join(output, "\n")
}

func (s *If) String() string {
	str := fmt.Sprintf("if %s {\n%s\n}", s.Condition, indent(s.Body.String()))
	if len(s.Else) > 0 {
		str += fmt.Sprintf(" else {\n%s\n}", indent(s.Else.String()))
	}
	return str
}

type While struct {
	pos       Position
	Condition Expression
	Body      Block
}

func (s *While) statementNode()     {}
func (s *While) Position() Position { return s.pos }

func (s *While) String() string {
	return fmt.Sprintf("while %s {\n%s\n}", s.Condition, indent(s.Body.String()))
}

type For struct {
	pos      Position
	Name     string
	Iterable Expression
	Body     Block
}

func (s *For) statementNode()     {}
func (s *For) Position() Position { return s.pos }

func (s *For) String() string {
	return fmt.Sprintf("for %s in %s {\n%s\n}", s.Name, s.Iterable, indent(s.Body.String()))
}

type Return struct {
	pos    Position
	Result Expression
}

func (s *Return) statementNode()     {}
func (s *Return) Position() Position { return s.pos }

func (s *Return) String() string {
	return fmt.Sprintf("return %s", s.Result)
}

type ExpressionStatement struct {
	pos        Position
	Expression Expression
}

func (s *ExpressionStatement) statementNode()     {}
func (s *ExpressionStatement) Position() Position { return s.pos }

func (s *ExpressionStatement) String() string {
	return fmt.Sprintf("%s", s.Expression)
}

type FunctionDefinition struct {
	pos        Position
	Name       string
	Parameters []string
	Ellipsis   bool
	Body       Block
}

func (s *FunctionDefinition) statementNode()     {}
func (s *FunctionDefinition) Position() Position { return s.pos }

func (s *FunctionDefinition) String() string {
	ellipsisStr := ""
	if s.Ellipsis {
		ellipsisStr = "..."
	}
	bodyStr := ""
	if len(s.Body) != 0 {
		bodyStr = "\n" + indent(s.Body.String()) + "\n"
	}
	return fmt.Sprintf("function %s(%s%s) {%s}",
		s.Name, strings.Join(s.Parameters, ", "), ellipsisStr, bodyStr)
}

type Expression interface {
	Position() Position
	expressionNode()
}

type Binary struct {
	pos      Position
	Left     Expression
	Operator Token
	Right    Expression
}

func (e *Binary) expressionNode()    {}
func (e *Binary) Position() Position { return e.pos }

func (e *Binary) String() string {
	return fmt.Sprintf("(%s %s %s)", e.Left, e.Operator, e.Right)
}

type Unary struct {
	pos      Position
	Operator Token
	Operand  Expression
}

func (e *Unary) expressionNode()    {}
func (e *Unary) Position() Position { return e.pos }

func (e *Unary) String() string {
	space := ""
	if e.Operator == NOT {
		space = " "
	}
	return fmt.Sprintf("(%s%s%s)", e.Operator, space, e.Operand)
}

type Call struct {
	pos       Position
	Function  Expression
	Arguments []Expression
	Ellipsis  bool
}

func (e *Call) expressionNode()    {}
func (e *Call) Position() Position { return e.pos }

func (e *Call) String() string {
	args := []string{}
	for _, arg := range e.Arguments {
		args = append(args, fmt.Sprintf("%s", arg))
	}
	ellipsisStr := ""
	if e.Ellipsis {
		ellipsisStr = "..."
	}
	return fmt.Sprintf("%s(%s%s)", e.Function, strings.Join(args, ", "), ellipsisStr)
}

type Literal struct {
	pos   Position
	Value interface{}
}

func (e *Literal) expressionNode()    {}
func (e *Literal) Position() Position { return e.pos }

func (e *Literal) String() string {
	if e.Value == nil {
		return "nil"
	}
	if s, ok := e.Value.(string); ok {
		return fmt.Sprintf("%q", s)
	}
	return fmt.Sprintf("%v", e.Value)
}

type List struct {
	pos    Position
	Values []Expression
}

func (e *List) expressionNode()    {}
func (e *List) Position() Position { return e.pos }

func (e *List) String() string {
	values := []string{}
	for _, value := range e.Values {
		values = append(values, fmt.Sprintf("%s", value))
	}
	return fmt.Sprintf("[%s]", strings.Join(values, ", "))
}

type MapItem struct {
	Key   Expression
	Value Expression
}

type Map struct {
	pos   Position
	Items []MapItem
}

func (e *Map) expressionNode()    {}
func (e *Map) Position() Position { return e.pos }

func (e *Map) String() string {
	items := []string{}
	for _, item := range e.Items {
		items = append(items, fmt.Sprintf("%s: %s", item.Key, item.Value))
	}
	return fmt.Sprintf("{%s}", strings.Join(items, ", "))
}

type ClassDefinition struct {
	pos       Position // The position in the source code where the class is defined
	ClassName string   // The name of the class
	Parent    *string  // Optional parent class (for inheritance)
	Body      []Statement
}

func (e *ClassDefinition) statementNode()     {}
func (e *ClassDefinition) Position() Position { return e.pos }
func (e *ClassDefinition) String() string {
	bodyStr := ""
	if len(e.Body) != 0 {
		//bodyStr = "\n" + indent(e.Body.String()) + "\n"
	}
	print("class %s {%s}", e.ClassName, bodyStr)
	return fmt.Sprintf("class %s {%s}", e.ClassName, bodyStr)
}

type NewExpression struct {
	pos       Position
	ClassName string
	Arguments []string
}

func (e *NewExpression) expressionNode()    {}
func (e *NewExpression) Position() Position { return e.pos }
func (e *NewExpression) String() string {
	return fmt.Sprintf("new %s(%s)", e.ClassName, strings.Join(e.Arguments, ", "))
}

// PropertyAccess represents accessing a property of an object, e.g., `$object->property`
type PropertyAccess struct {
	pos      Position   // Position of the `->` operator in the source code
	Object   Expression // The object whose property is being accessed
	Property string     // The name of the property being accessed
}

func (e *PropertyAccess) expressionNode()    {}
func (e *PropertyAccess) Position() Position { return e.pos }

// MethodCall represents a method call on an object, e.g., `$object->method(arg1, arg2)`
type MethodCall struct {
	pos       Position     // Position of the `->` operator in the source code
	Object    Expression   // The object on which the method is being called
	Method    string       // The name of the method being called
	Arguments []Expression // Arguments passed to the method
}

func (e *MethodCall) expressionNode()    {}
func (e *MethodCall) Position() Position { return e.pos }

type FunctionExpression struct {
	pos        Position
	Parameters []string
	Ellipsis   bool
	Body       Block
}

func (e *FunctionExpression) expressionNode()    {}
func (e *FunctionExpression) Position() Position { return e.pos }

func (e *FunctionExpression) String() string {
	ellipsisStr := ""
	if e.Ellipsis {
		ellipsisStr = "..."
	}
	bodyStr := ""
	if len(e.Body) != 0 {
		bodyStr = "\n" + indent(e.Body.String()) + "\n"
	}
	return fmt.Sprintf("function(%s%s) {%s}", strings.Join(e.Parameters, ", "), ellipsisStr, bodyStr)
}

type Subscript struct {
	pos       Position
	Container Expression
	Subscript Expression
}

func (e *Subscript) expressionNode()    {}
func (e *Subscript) Position() Position { return e.pos }

func (e *Subscript) String() string {
	return fmt.Sprintf("%s[%s]", e.Container, e.Subscript)
}

type Variable struct {
	pos  Position
	Name string
}

func (e *Variable) expressionNode()    {}
func (e *Variable) Position() Position { return e.pos }

func (e *Variable) String() string {
	return e.Name
}

type SemiTag struct {
	pos Position
}

func (e *SemiTag) expressionNode()    {}
func (e *SemiTag) Position() Position { return e.pos }
