package parser

import (
	"fmt"
	"strings"
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
	Definitions []*TopDef `@@*`
}

func (ast *LatteProgram) Print(c *ParsingContext) string {
	defs := []string{}
	for _, def := range ast.Definitions {
		defs = append(defs, def.Print(c))
	}
	return fmt.Sprintf("%s\n", strings.Join(defs, "\n\n"))
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
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
	}
	c.BlockPush()
	for _, statement := range ast.Statements {
		statementsList = append(statementsList, statement.Print(c))
	}
	c.BlockPop()
	return fmt.Sprintf("{\n%s\n%s}", strings.Join(statementsList, "\n"), strings.Repeat("  ", c.BlockDepth))
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

func (ast *Statement) formatStatementInstruction(statement string, c *ParsingContext) string {
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
		return statement
	}
	return fmt.Sprintf("%s%s", strings.Repeat("  ", c.BlockDepth), statement)
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
		ret = fmt.Sprintf("%s;", ast.Expression.Print(c))
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
		return fmt.Sprintf("if (%s) %s else %s", ast.Condition.Print(c), makeBlockFromStatement(ast.Then).Print(c), makeBlockFromStatement(ast.Else).Print(c))
	}
	return fmt.Sprintf("if (%s) %s", ast.Condition.Print(c), ast.Then.Print(c))
}

type While struct {
	Condition *Expression `"while" "(" @@ ")"`
	Do *Statement `@@`
}

func (ast *While) Print(c *ParsingContext) string {
	c.PrinterConfiguration.SkipStatementIdent = true
	body := ast.Do.Print(c)
	return fmt.Sprintf("while (%s) %s", ast.Condition.Print(c), body)
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