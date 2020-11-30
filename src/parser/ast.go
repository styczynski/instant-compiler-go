package parser

import (
	"fmt"
	"strings"
)

type LatteProgram struct {
	Definitions []*TopDef `@@*`
}

func (ast *LatteProgram) Print(c *ParsingContext) string {
	defs := []string{}
	for _, def := range ast.Definitions {
		defs = append(defs, def.Print(c))
	}
	return fmt.Sprintf("%s", strings.Join(defs, ";\n"))
}

type FnDef struct {
	ReturnType Type `@@`
	Name string `@Ident`
	Arg []*Arg `"(" (@@ ( "," @@ )*)? ")"`
	Body Block `@@`
}

func (ast *FnDef) Print(c *ParsingContext) string {
	argsList := []string{}
	for _, arg := range ast.Arg {
		argsList = append(argsList, arg.Print(c))
	}

	return fmt.Sprintf("%s %s(%s) %s",
		ast.ReturnType.Print(c),
		ast.Name,
		strings.Join(argsList, ", "),
		ast.Body.Print(c))
}

type Block struct {
	Statements []*Statement `"{" @@* "}"`
}

func (ast *Block) Print(c *ParsingContext) string {
	statementsList := []string{}
	for _, statement := range ast.Statements {
		statementsList = append(statementsList, statement.Print(c))
	}
	return fmt.Sprintf("{\n%s\n}", strings.Join(statementsList, "\n"))
}

type Arg struct {
	ArgumentType Type `@@`
	Name string `@Ident`
}

func (ast *Arg) Print(c *ParsingContext) string {
	return fmt.Sprintf("%s %s", ast.ArgumentType.Print(c), ast.Name)
}

type Type struct {
	Name string `(@ "int" | "void")`
}

func (ast *Type) Print(c *ParsingContext) string {
	return fmt.Sprintf("%s", ast.Name)
}

type TopDef struct {
	Function *FnDef `@@`
}

func (ast *TopDef) Print(c *ParsingContext) string {
	return fmt.Sprintf("%s", ast.Function.Print(c))
}

type Return struct {
	Expression *Expression `"return" (@@)? ";"`
}

func (ast *Return) Print(c *ParsingContext) string {
	return fmt.Sprintf("return %s;", ast.Expression.Print(c))
}

type Statement struct {
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

func (ast *Statement) Print(c *ParsingContext) string {
	if ast.IsEmpty() {
		return ";"
	} else if ast.IsBlockStatement() {
		return ast.BlockStatement.Print(c)
	} else if ast.IsDeclaration() {
		return ast.Declaration.Print(c)
	} else if ast.IsAssignment() {
		return ast.Assignment.Print(c)
	} else if ast.IsUnaryStatement() {
		return ast.UnaryStatement.Print(c)
	} else if ast.IsReturn() {
		return ast.Return.Print(c)
	} else if ast.IsIf() {
		return ast.If.Print(c)
	} else if ast.IsWhile() {
		return ast.While.Print(c)
	} else if ast.IsExpression() {
		return fmt.Sprintf("%s;", ast.Expression.Print(c))
	}
	return "UNKNOWN"
}

type Assignment struct {
	TargetName string `@Ident`
	Value *Expression `"=" @@ ";"`
}

func (ast *Assignment) Print(c *ParsingContext) string {
	return fmt.Sprintf("%s = %s;", ast.TargetName, ast.Value.Print(c))
}

type UnaryStatement struct {
	TargetName *string `@Ident`
	Operation string `@( "+" "+" | "-" "-" ) ";"`
}

func (ast *UnaryStatement) Print(c *ParsingContext) string {
	return fmt.Sprintf("%s%s;", *ast.TargetName, ast.Operation)
}

type If struct {
	Condition *Expression `"if" "(" @@ ")"`
	Then *Statement `@@`
	Else *Statement `( "else" @@ )?`
}

func (ast *If) HasElseBlock() bool {
	return ast.Else != nil
}

func (ast *If) Print(c *ParsingContext) string {
	if ast.HasElseBlock(){
		return fmt.Sprintf("if (%s) { %s } else { %s }", ast.Condition.Print(c), ast.Then.Print(c), ast.Else.Print(c))
	}
	return fmt.Sprintf("if (%s) { %s }", ast.Condition.Print(c), ast.Then.Print(c))
}

type While struct {
	Condition *Expression `"while" "(" @@ ")"`
	Do *Statement `@@`
}

func (ast *While) Print(c *ParsingContext) string {
	return fmt.Sprintf("while (%s) { %s }", ast.Condition.Print(c), ast.Do.Print(c))
}

type Declaration struct {
	DeclarationType Type `@@`
	Items []*DeclarationItem `( @@ ( "," @@ )* ) ";"`
}

func (ast *Declaration) Print(c *ParsingContext) string {
	declarationItemsList := []string{}
	for _, item := range ast.Items {
		declarationItemsList = append(declarationItemsList, item.Print(c))
	}
	return fmt.Sprintf("%s %s", ast.DeclarationType.Print(c), strings.Join(declarationItemsList, ", "))
}

type DeclarationItem struct {
	Name string `@Ident`
	Initializer *Expression `( "=" @@ )?`
}

func (ast *DeclarationItem) HasInitializer() bool {
	return ast.Initializer != nil
}

func (ast *DeclarationItem) Print(c *ParsingContext) string {
	if ast.HasInitializer() {
		return fmt.Sprintf("%s = %s", ast.Name, ast.Initializer.Print(c))
	}
	return fmt.Sprintf("%s", ast.Name)
}