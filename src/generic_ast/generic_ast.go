package generic_ast

import (
	context2 "context"
	"fmt"
	"reflect"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/styczynski/latte-compiler/src/parser/context"
)

type ComplexASTNode struct {
	BaseASTNode
	Tokens []lexer.Token
}

type BaseASTNode struct {
	Pos lexer.Position
	EndPos lexer.Position
}

func (ast *BaseASTNode) End() lexer.Position {
	return ast.EndPos
}

func (ast *BaseASTNode) Begin() lexer.Position {
	return ast.Pos
}

type NodeWithPosition interface {
	Begin() lexer.Position
	End() lexer.Position
}

type PrintableNode interface {
	Print(c *context.ParsingContext) string
}

type TraversableNode interface {
	NodeWithPosition
	GetChildren() []TraversableNode
	GetNode() interface{}
	Parent() TraversableNode
	OverrideParent(node TraversableNode)
}

func ReplaceExpressionRecursively(node TraversableNode, oldNode TraversableNode, newNode TraversableNode) TraversableNode {
	//fmt.Printf("REPLACE %p -> %p <%b>\n", oldNode, newNode, oldNode == newNode)
	if e, ok := node.(Expression); ok {
		var mapper ExpressionVisitor
		visitedNodes := map[Expression]struct{}{}
		mapper = func(parent Expression, e Expression, context VisitorContext) {
			//fmt.Printf("REPLACE %p\n", e)
			if _, ok := visitedNodes[e]; ok {
				return
			}
			visitedNodes[e] = struct{}{}
			if e.(TraversableNode) == oldNode {
				//if newExpr, ok := newNode.(Expression); ok {
				//	return newExpr.Map(parent, mapper, context)
				//}
				//fmt.Printf("-- REPLACE --\n")
				if reflect.TypeOf(e).Kind() == reflect.Ptr {
					// Pointer:
					reflect.ValueOf(e).Elem().Set(reflect.ValueOf(newNode).Elem())
					//reflect.New(reflect.ValueOf(newNode).Elem().Type()).Interface().(User)

				} else {
					// Not pointer:
					panic("Cannot replace node that is not a pointer")
				}
			}
			e.Visit(parent, mapper, context)
		}
		e.Visit(e, mapper, NewEmptyVisitorContext())
	}
	//fmt.Printf("Replacement done\n")
	return node
}

type NodeWithSyntaxValidation interface {
	Validate(c *context.ParsingContext) NodeError
}

type NormalNode interface {
	PrintableNode
	TraversableNode
}

type ConstFoldableNode interface {
	ConstFold() TraversableNode
}

type NodeWithFoldingValidation interface {
	ValidateConstFold() (error, TraversableNode)
}

type ConstExtractableNode interface {
	ExtractConst() (TraversableNode, bool)
}

type NormalNodeWithID interface {
	NormalNode
}

type NormalNodeSelection struct {
	node NormalNode
	id int
	describe func(src context.SelectionBlock, id int, mappingID func(block context.SelectionBlock) int) []string
}

func NewNormalNodeSelection(node NormalNode, id int, describe func(src context.SelectionBlock, id int, mappingID func(block context.SelectionBlock) int) []string) NormalNodeSelection {
	return NormalNodeSelection{
		node: node,
		describe: describe,
		id: id,
	}
}

func (sel NormalNodeSelection) GetID() int {
	return sel.id
}

func (sel NormalNodeSelection) GetNode() NormalNode {
	return sel.node
}

func (sel NormalNodeSelection) Describe(src context.SelectionBlock, id int, mappingID func(block context.SelectionBlock) int) []string {
	return sel.describe(src, id, mappingID)
}

func (sel NormalNodeSelection) Begin() (int, int) {
	n := sel.node.Begin()
	return n.Line, n.Column
}

func (sel NormalNodeSelection) End() (int, int) {
	n := sel.node.End()
	return n.Line, n.Column
}

type TraversableNodeToken struct {
	Token string
	BeginPos lexer.Position
	EndPos lexer.Position
	ParentNode TraversableNode
}

type TraversableNodeValue struct {
	Value interface{}
	Type string
	BeginPos lexer.Position
	EndPos lexer.Position
	ParentNode TraversableNode
}

func MakeTraversableNodeValue(parent TraversableNode, value interface{}, typeName string, begin lexer.Position, end lexer.Position) TraversableNode {
	return &TraversableNodeValue{
		Value: value,
		Type: typeName,
		BeginPos: begin,
		EndPos: end,
		ParentNode: parent,
	}
}

func MakeTraversableNodeToken(parent TraversableNode, value string, begin lexer.Position, end lexer.Position) TraversableNode {
	return &TraversableNodeToken{
		Token: value,
		BeginPos: begin,
		EndPos: end,
		ParentNode: parent,
	}
}

func (*TraversableNodeValue) GetChildren() []TraversableNode {
	return []TraversableNode{}
}

func (*TraversableNodeToken) GetChildren() []TraversableNode {
	return []TraversableNode{}
}

func (ast *TraversableNodeValue) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *TraversableNodeToken) OverrideParent(node TraversableNode) {
	ast.ParentNode = node
}

func (ast *TraversableNodeValue) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *TraversableNodeToken) Parent() TraversableNode {
	return ast.ParentNode
}

func (ast *TraversableNodeValue) GetNode() interface{} {
	return ast
}

func (ast *TraversableNodeToken) GetNode() interface{} {
	return ast
}

func (ast *TraversableNodeValue) End() lexer.Position {
	return ast.EndPos
}

func (ast *TraversableNodeToken) End() lexer.Position {
	return ast.EndPos
}

func (ast *TraversableNodeValue) Begin() lexer.Position {
	return ast.BeginPos
}

func (ast *TraversableNodeToken) Begin() lexer.Position {
	return ast.BeginPos
}

func (ast *TraversableNodeValue) Print(c *context.ParsingContext) string {
	return fmt.Sprintf("%v", ast.Value)
}

func (ast *TraversableNodeToken) Print(c *context.ParsingContext) string {
	return ast.Token
}

func TraverseAST(node TraversableNode, visitor func(ast TraversableNode)) {
	children := node.GetChildren()
	for _, child := range children {
		visitor(child)
		TraverseAST(child, visitor)
	}
}

type ExpressionVisitor = func (parent Expression, e Expression, context VisitorContext)
type ExpressionMapper = func (parent Expression, e Expression, context VisitorContext, backwards bool) Expression

type VisitorContext interface {
	context2.Context
	WithValue(key, val interface{}) VisitorContext
	Set(key, val interface{}) interface{}
	Get(key interface{}, defaultValue interface{}) interface{}
}

type VisitorContextImpl struct {
	context2.Context
}

func (c *VisitorContextImpl) Set(key, val interface{}) interface{} {
	v := c.Value(key)
	c.Context = context2.WithValue(c, key, val)
	return v
}

func (c *VisitorContextImpl) Get(key interface{}, defaultValue interface{}) interface{} {
	v := c.Context.Value(key)
	if v == nil {
		return defaultValue
	}
	return v
}

func (c *VisitorContextImpl) WithValue(key, val interface{}) VisitorContext {
	return &VisitorContextImpl{
		context2.WithValue(c, key, val),
	}
}

func NewEmptyVisitorContext() VisitorContext {
	return &VisitorContextImpl{
		context2.Background(),
	}
}

// An Expression is basically an AST node. In its simplest form, it's lambda calculus
type Expression interface {
	Body() Expression
	Map(parent Expression, mapper ExpressionMapper, context VisitorContext) Expression
	Visit(parent Expression, mapper ExpressionVisitor, context VisitorContext)
}

type Batch struct {
	Exp []Expression
}

type NodeError interface {
	ErrorName() string
	GetMessage() string
	GetCliMessage() string
	GetSource() Expression
}

func NewNodeError(errorName string, source Expression, message string, cliMessage string) *NodeErrorImpl {
	return &NodeErrorImpl{
		Message:    message,
		CliMessage: cliMessage,
		Source:     source,
		ErrorTypeName: errorName,
	}
}

type NodeErrorImpl struct {
	ErrorTypeName string
	Message string
	CliMessage string
	Source Expression
}

func (e *NodeErrorImpl) ErrorName() string {
	return e.ErrorTypeName
}

func (e *NodeErrorImpl) GetMessage() string {
	return e.Message
}

func (e *NodeErrorImpl) GetCliMessage() string {
	return e.CliMessage
}

func (e *NodeErrorImpl) GetSource() Expression {
	return e.Source
}
