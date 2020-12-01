package parser

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
)

func makeBlockFromStatement(statement *Statement) *Block {
	if statement.IsBlockStatement() {
		return statement.BlockStatement
	}
	return &Block{
		Statements: []*Statement{ statement },
	}
}

func makeBlockFromExpression(expression *Expression) *Block {
	return makeBlockFromStatement(&Statement{
		Expression: expression,
	})
}

type LatteProgram struct {
	BaseASTNode
	Definitions []*TopDef `@@*`
}

func (ast *LatteProgram) Begin() lexer.Position {
	return ast.Pos
}

func (ast *LatteProgram) End() lexer.Position {
	return ast.EndPos
}

func (ast *LatteProgram) GetNode() interface{} {
	return ast
}

func (ast *LatteProgram) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Definitions))
	for _, child := range ast.Definitions {
		nodes = append(nodes, child)
	}
	return nodes
}

func (ast *LatteProgram) Print(c *ParsingContext) string {
	defs := []string{}
	for _, def := range ast.Definitions {
		defs = append(defs, def.Print(c))
	}
	return printNode(c, ast, "%s\n", strings.Join(defs, "\n\n"))
}

type FnDef struct {
	BaseASTNode
	ReturnType Type `@@`
	Name string `@Ident`
	Arg []*Arg `"(" (@@ ( "," @@ )*)? ")"`
	Body Block `@@`
}

func (ast *FnDef) Begin() lexer.Position {
	return ast.Pos
}

func (ast *FnDef) End() lexer.Position {
	return ast.EndPos
}

func (ast *FnDef) GetNode() interface{} {
	return ast
}

func (ast *FnDef) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Arg) + 3)
	nodes = append(nodes, &ast.ReturnType)
	nodes = append(nodes, MakeTraversableNodeToken(ast.Name, ast.Pos, ast.EndPos))

	for _, child := range ast.Arg {
		nodes = append(nodes, child)
	}
	nodes = append(nodes, &ast.Body)

	return nodes
}

func (ast *FnDef) Print(c *ParsingContext) string {
	argsList := []string{}
	for _, arg := range ast.Arg {
		argsList = append(argsList, arg.Print(c))
	}

	return printNode(c, ast, "%s %s(%s) %s",
		ast.ReturnType.Print(c),
		ast.Name,
		strings.Join(argsList, ", "),
		ast.Body.Print(c))
}

type Block struct {
	BaseASTNode
	Statements []*Statement `"{" @@* "}"`
}

func (ast *Block) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Block) End() lexer.Position {
	return ast.EndPos
}

func (ast *Block) GetNode() interface{} {
	return ast
}

func (ast *Block) Print(c *ParsingContext) string {
	statementsList := []string{}
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
	}
	c.BlockPush()
	for _, statement := range ast.Statements {
		statementsList = append(statementsList, statement.Print(c))
	}
	c.BlockPop()
	return printNode(c, ast, "{\n%s\n%s}", strings.Join(statementsList, "\n"), strings.Repeat("  ", c.BlockDepth))
}

func (ast *Block) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Statements))
	for _, child := range ast.Statements {
		nodes = append(nodes, child)
	}
	return nodes
}

type Arg struct {
	BaseASTNode
	ArgumentType Type `@@`
	Name string `@Ident`
}

func (ast *Arg) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Arg) End() lexer.Position {
	return ast.EndPos
}

func (ast *Arg) GetNode() interface{} {
	return ast
}

func (ast *Arg) GetChildren() []TraversableNode {
	return []TraversableNode{
		&ast.ArgumentType,
		MakeTraversableNodeToken(ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *Arg) Print(c *ParsingContext) string {
	return printNode(c, ast, "%s %s", ast.ArgumentType.Print(c), ast.Name)
}

type Type struct {
	BaseASTNode
	Name string `(@ "int" | "void")`
}

func (ast *Type) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Type) End() lexer.Position {
	return ast.EndPos
}

func (ast *Type) GetNode() interface{} {
	return ast
}

func (ast *Type) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast.Name, ast.Pos, ast.EndPos),
	}
}

func (ast *Type) Print(c *ParsingContext) string {
	return printNode(c, ast, "%s", ast.Name)
}

type TopDef struct {
	BaseASTNode
	Function *FnDef `@@`
}

func (ast *TopDef) Begin() lexer.Position {
	return ast.Pos
}

func (ast *TopDef) End() lexer.Position {
	return ast.EndPos
}

func (ast *TopDef) GetNode() interface{} {
	return ast
}

func (ast *TopDef) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Function,
	}
}

func (ast *TopDef) Print(c *ParsingContext) string {
	return printNode(c, ast, "%s", ast.Function.Print(c))
}

type Return struct {
	BaseASTNode
	Expression *Expression `"return" (@@)? ";"`
}

func (ast *Return) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Return) End() lexer.Position {
	return ast.EndPos
}

func (ast *Return) GetNode() interface{} {
	return ast
}

func (ast *Return) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Expression,
	}
}

func (ast *Return) Print(c *ParsingContext) string {
	return printNode(c, ast, "return %s;", ast.Expression.Print(c))
}

type Statement struct {
	BaseASTNode
	Empty *string `";"`
	BlockStatement *Block `| @@`
	Declaration *Declaration `| @@`
	Assignment *Assignment `| @@`
	UnaryStatement *UnaryStatement `| @@`
	Return *Return `| @@`
	If *If `| @@`
	While *While `| @@`
	Expression *Expression `| @@ ";"`
}

func (ast *Statement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Statement) End() lexer.Position {
	return ast.EndPos
}

func (ast *Statement) GetNode() interface{} {
	return ast
}

func (ast *Statement) IsEmpty() bool {
	return ast.Empty != nil
}

func (ast *Statement) IsBlockStatement() bool {
	return ast.BlockStatement != nil
}

func (ast *Statement) IsDeclaration() bool {
	return ast.Declaration != nil
}

func (ast *Statement) IsAssignment() bool {
	return ast.Assignment != nil
}

func (ast *Statement) IsUnaryStatement() bool {
	return ast.UnaryStatement != nil
}

func (ast *Statement) IsReturn() bool {
	return ast.Return != nil
}

func (ast *Statement) IsIf() bool {
	return ast.If != nil
}

func (ast *Statement) IsWhile() bool {
	return ast.While != nil
}

func (ast *Statement) IsExpression() bool {
	return ast.Expression != nil
}

func (ast *Statement) GetChildren() []TraversableNode {
	if ast.IsEmpty() {
		return []TraversableNode{ MakeTraversableNodeToken(*ast.Empty, ast.Pos, ast.EndPos) }
	} else if ast.IsBlockStatement() {
		return []TraversableNode{ ast.BlockStatement }
	} else if ast.IsDeclaration() {
		return []TraversableNode{ ast.Declaration }
	} else if ast.IsAssignment() {
		return []TraversableNode{ ast.Assignment }
	} else if ast.IsUnaryStatement() {
		return []TraversableNode{ ast.UnaryStatement }
	} else if ast.IsReturn() {
		return []TraversableNode{ ast.Return }
	} else if ast.IsIf() {
		return []TraversableNode{ ast.If }
	} else if ast.IsWhile() {
		return []TraversableNode{ ast.While }
	} else if ast.IsExpression() {
		return []TraversableNode{ ast.Expression }
	}
	return []TraversableNode{}
}

func (ast *Statement) formatStatementInstruction(statement string, c *ParsingContext) string {
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
		return statement
	}
	return printNode(c, ast, "%s%s", strings.Repeat("  ", c.BlockDepth), statement)
}

func (ast *Statement) Print(c *ParsingContext) string {
	ret := "UNKNOWN"
	propagateSkipStatementIdent := false
	if ast.IsEmpty() {
		ret = ";"
	} else if ast.IsBlockStatement() {
		if c.PrinterConfiguration.SkipStatementIdent {
			propagateSkipStatementIdent = true
		}
		ret = ast.BlockStatement.Print(c)
	} else if ast.IsDeclaration() {
		ret = ast.Declaration.Print(c)
	} else if ast.IsAssignment() {
		ret = ast.Assignment.Print(c)
	} else if ast.IsUnaryStatement() {
		ret = ast.UnaryStatement.Print(c)
	} else if ast.IsReturn() {
		ret = ast.Return.Print(c)
	} else if ast.IsIf() {
		ret =  ast.If.Print(c)
	} else if ast.IsWhile() {
		ret = ast.While.Print(c)
	} else if ast.IsExpression() {
		ret = printNode(c, ast, "%s;", ast.Expression.Print(c))
	}
	if propagateSkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = true
		ret = ast.formatStatementInstruction(ret, c)
	} else {
		ret = ast.formatStatementInstruction(ret, c)
	}
	return ret
}

type Assignment struct {
	BaseASTNode
	TargetName string `@Ident`
	Value *Expression `"=" @@ ";"`
}

func (ast *Assignment) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Assignment) End() lexer.Position {
	return ast.EndPos
}

func (ast *Assignment) GetNode() interface{} {
	return ast
}

func (ast *Assignment) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast.TargetName, ast.Pos, ast.EndPos),
		ast.Value,
	}
}

func (ast *Assignment) Print(c *ParsingContext) string {
	return printNode(c, ast, "%s = %s;", ast.TargetName, ast.Value.Print(c))
}

type UnaryStatement struct {
	BaseASTNode
	TargetName *string `@Ident`
	Operation string `@( "+" "+" | "-" "-" ) ";"`
}

func (ast *UnaryStatement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *UnaryStatement) End() lexer.Position {
	return ast.EndPos
}

func (ast *UnaryStatement) GetNode() interface{} {
	return ast
}

func (ast *UnaryStatement) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(*ast.TargetName, ast.Pos, ast.EndPos),
		MakeTraversableNodeToken(ast.Operation, ast.Pos, ast.EndPos),
	}
}

func (ast *UnaryStatement) Print(c *ParsingContext) string {
	return printNode(c, ast, "%s%s;", *ast.TargetName, ast.Operation)
}

type If struct {
	BaseASTNode
	Condition *Expression `"if" "(" @@ ")"`
	Then *Statement `@@`
	Else *Statement `( "else" @@ )?`
}

func (ast *If) Begin() lexer.Position {
	return ast.Pos
}

func (ast *If) End() lexer.Position {
	return ast.EndPos
}

func (ast *If) GetNode() interface{} {
	return ast
}

func (ast *If) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Condition,
		ast.Then,
		ast.Else,
	}
}

func (ast *If) HasElseBlock() bool {
	return ast.Else != nil
}

func (ast *If) Print(c *ParsingContext) string {
	if ast.HasElseBlock(){
		return printNode(c, ast, "if (%s) %s else %s", ast.Condition.Print(c), makeBlockFromStatement(ast.Then).Print(c), makeBlockFromStatement(ast.Else).Print(c))
	}
	return printNode(c, ast, "if (%s) %s", ast.Condition.Print(c), ast.Then.Print(c))
}

type While struct {
	BaseASTNode
	Condition *Expression `"while" "(" @@ ")"`
	Do *Statement `@@`
}

func (ast *While) Begin() lexer.Position {
	return ast.Pos
}

func (ast *While) End() lexer.Position {
	return ast.EndPos
}

func (ast *While) GetNode() interface{} {
	return ast
}

func (ast *While) Print(c *ParsingContext) string {
	c.PrinterConfiguration.SkipStatementIdent = true
	body := ast.Do.Print(c)
	return printNode(c, ast, "while (%s) %s", ast.Condition.Print(c), body)
}

func (ast *While) GetChildren() []TraversableNode {
	return []TraversableNode{
		ast.Condition,
		ast.Do,
	}
}

type Declaration struct {
	BaseASTNode
	DeclarationType Type `@@`
	Items []*DeclarationItem `( @@ ( "," @@ )* ) ";"`
}

func (ast *Declaration) Begin() lexer.Position {
	return ast.Pos
}

func (ast *Declaration) End() lexer.Position {
	return ast.EndPos
}

func (ast *Declaration) GetNode() interface{} {
	return ast
}

func (ast *Declaration) GetChildren() []TraversableNode {
	nodes := make([]TraversableNode, len(ast.Items)+1)
	nodes = append(nodes, &ast.DeclarationType)
	for _, child := range ast.Items {
		nodes = append(nodes, child)
	}
	return nodes
}

func (ast *Declaration) Print(c *ParsingContext) string {
	declarationItemsList := []string{}
	for _, item := range ast.Items {
		declarationItemsList = append(declarationItemsList, item.Print(c))
	}
	return printNode(c, ast, "%s %s", ast.DeclarationType.Print(c), strings.Join(declarationItemsList, ", "))
}

type DeclarationItem struct {
	BaseASTNode
	Name string `@Ident`
	Initializer *Expression `( "=" @@ )?`
}

func (ast *DeclarationItem) Begin() lexer.Position {
	return ast.Pos
}

func (ast *DeclarationItem) End() lexer.Position {
	return ast.EndPos
}

func (ast *DeclarationItem) GetNode() interface{} {
	return ast
}

func (ast *DeclarationItem) GetChildren() []TraversableNode {
	return []TraversableNode{
		MakeTraversableNodeToken(ast.Name, ast.Pos, ast.EndPos),
		ast.Initializer,
	}
}

func (ast *DeclarationItem) HasInitializer() bool {
	return ast.Initializer != nil
}

func (ast *DeclarationItem) Print(c *ParsingContext) string {
	if ast.HasInitializer() {
		return printNode(c, ast, "%s = %s", ast.Name, ast.Initializer.Print(c))
	}
	return printNode(c, ast, "%s", ast.Name)
}