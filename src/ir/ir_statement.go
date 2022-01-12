package ir

import (
	"fmt"
	"strings"

	"github.com/alecthomas/participle/v2/lexer"

	"github.com/styczynski/latte-compiler/src/flow_analysis/cfg"
	"github.com/styczynski/latte-compiler/src/generic_ast"
	"github.com/styczynski/latte-compiler/src/parser/context"
	"github.com/styczynski/latte-compiler/src/parser/utils"
	"github.com/styczynski/latte-compiler/src/type_checker/hindley_milner"
)

type IRStatement struct {
	generic_ast.BaseASTNode
	Empty      *IREmpty      `@@`
	Exit       *IRExit       `| @@`
	If         *IRIf         `| @@`
	Phi        *IRPhi        `| @@`
	Expression *IRExpression `| @@`
	Copy       *IRCopy       `| @@`
	Const      *IRConst      `| @@`
	Call       *IRCall       `| @@`
	ParentNode generic_ast.TraversableNode

	VarIn   cfg.VariableSet
	VarOut  cfg.VariableSet
	Reached cfg.ReachedBlocks
	meta    IRMeta
}

func (ast *IRStatement) SetAllocationInfo(allocInfo IRAllocationMap) {
	ast.meta.Allocation = allocInfo
}

func (ast *IRStatement) GetAllocationInfo() IRAllocationMap {
	return ast.meta.Allocation
}

func (ast *IRStatement) SetFlowAnalysisProps(
	VarIn cfg.VariableSet,
	VarOut cfg.VariableSet,
	Reached cfg.ReachedBlocks,
) {
	ast.VarIn = VarIn
	ast.VarOut = VarOut
	ast.Reached = Reached
}

func WrapIREmpty() *IRStatement {
	return &IRStatement{
		Empty: &IREmpty{},
	}
}

func WrapIRPhi(ast *IRPhi) *IRStatement {
	ret := &IRStatement{
		Phi: ast,
	}
	ast.OverrideParent(ret)
	return ret
}

func WrapIRIf(ast *IRIf) *IRStatement {
	ret := &IRStatement{
		If: ast,
	}
	ast.OverrideParent(ret)
	return ret
}

func WrapIRExit(ast *IRExit) *IRStatement {
	ret := &IRStatement{
		Exit: ast,
	}
	ast.OverrideParent(ret)
	return ret
}

func WrapIRCall(ast *IRCall) *IRStatement {
	ret := &IRStatement{
		Call: ast,
	}
	ast.OverrideParent(ret)
	return ret
}

func WrapIRConst(ast *IRConst) *IRStatement {
	ret := &IRStatement{
		Const: ast,
	}
	ast.OverrideParent(ret)
	return ret
}

func WrapIRCopy(ast *IRCopy) *IRStatement {
	ret := &IRStatement{
		Copy: ast,
	}
	ast.OverrideParent(ret)
	return ret
}

func WrapIRExpression(ast *IRExpression) *IRStatement {
	ret := &IRStatement{
		Expression: ast,
	}
	ast.OverrideParent(ret)
	return ret
}

func (ast *IRStatement) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRStatement) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRStatement) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRStatement) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRStatement) GetNode() interface{} {
	return ast
}

func (ast *IRStatement) IsEmpty() bool {
	return (ast.Empty != nil || (!ast.IsExit() &&
		!ast.IsIf() &&
		!ast.IsPhi() &&
		!ast.IsCopy() &&
		!ast.IsConst() &&
		!ast.IsCall() &&
		!ast.IsExpression()))
}

func (ast *IRStatement) IsExit() bool {
	return ast.Exit != nil
}

func (ast *IRStatement) IsIf() bool {
	return ast.If != nil
}

func (ast *IRStatement) IsPhi() bool {
	return ast.Phi != nil
}

func (ast *IRStatement) IsCopy() bool {
	return ast.Copy != nil
}

func (ast *IRStatement) IsConst() bool {
	return ast.Const != nil
}

func (ast *IRStatement) IsCall() bool {
	return ast.Call != nil
}

func (ast *IRStatement) IsExpression() bool {
	return ast.Expression != nil
}

func (ast *IRStatement) GetChildren() []generic_ast.TraversableNode {
	if ast.IsEmpty() {
		return []generic_ast.TraversableNode{}
		//return []generic_ast.TraversableNode{ generic_ast.MakeTraversableNodeToken(ast, *ast.Empty, ast.Pos, ast.EndPos) }
	} else if ast.IsExit() {
		return []generic_ast.TraversableNode{ast.Exit}
	} else if ast.IsIf() {
		return []generic_ast.TraversableNode{ast.If}
	} else if ast.IsPhi() {
		return []generic_ast.TraversableNode{ast.Phi}
	} else if ast.IsCopy() {
		return []generic_ast.TraversableNode{ast.Copy}
	} else if ast.IsConst() {
		return []generic_ast.TraversableNode{ast.Const}
	} else if ast.IsCall() {
		return []generic_ast.TraversableNode{ast.Call}
	} else if ast.IsExpression() {
		return []generic_ast.TraversableNode{ast.Expression}
	}
	return []generic_ast.TraversableNode{}
}

func (ast *IRStatement) formatIRStatementInstruction(irstatement string, c *context.ParsingContext) string {
	if c.PrinterConfiguration.SkipStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = false
		return irstatement
	}
	return utils.PrintASTNode(c, ast, "%s%s %s %s %s", strings.Repeat("  ", c.BlockDepth), irstatement, ast.VarIn.String(), ast.VarOut.String(), ast.meta.String())
}

func (ast *IRStatement) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	fmt.Printf("RENAME IN STATMENT :DD %s\n", ast.VarOut)
	newVarIn := cfg.VariableSet{}
	for v, d := range ast.VarIn {
		newName := substDecl.Replace(substUsed.Replace(v))
		newVarIn[newName] = cfg.NewVariable(newName, d.Value())
	}
	newVarOut := cfg.VariableSet{}
	for v, d := range ast.VarOut {
		newName := substDecl.Replace(substUsed.Replace(v))
		newVarOut[newName] = cfg.NewVariable(newName, d.Value())
	}
	ast.VarIn = newVarIn
	ast.VarOut = newVarOut
	children := ast.GetChildren()
	if len(children) > 0 {
		children[0].(cfg.NodeWithVariableReplacement).RenameVariables(substUsed, substDecl)
	}
	fmt.Printf("RENAME IN STATMENT :C %s (with %v and %v)\n", ast.VarOut, substDecl, substUsed)
}

func (ast *IRStatement) Print(c *context.ParsingContext) string {
	ret := "UNKNOWN"
	propagateSkipIRStatementIdent := false
	if ast.IsEmpty() {
		ret = ";"
	} else if ast.IsExit() {
		ret = ast.Exit.Print(c)
	} else if ast.IsIf() {
		ret = ast.If.Print(c)
	} else if ast.IsPhi() {
		ret = ast.Phi.Print(c)
	} else if ast.IsCopy() {
		ret = ast.Copy.Print(c)
	} else if ast.IsConst() {
		ret = ast.Const.Print(c)
	} else if ast.IsCall() {
		ret = ast.Call.Print(c)
	} else if ast.IsExpression() {
		ret = utils.PrintASTNode(c, ast, "%s;", ast.Expression.Print(c))
	}
	if propagateSkipIRStatementIdent {
		c.PrinterConfiguration.SkipStatementIdent = true
		ret = ast.formatIRStatementInstruction(ret, c)
	} else {
		ret = ast.formatIRStatementInstruction(ret, c)
	}
	return ret
}

//////

func (ast *IRStatement) Body() generic_ast.Expression {
	if ast.IsEmpty() {
		return &IREmpty{}
		//return hindley_milner.Batch{Exp: []generic_ast.Expression{}}
	} else if ast.IsExit() {
		return ast.Exit
	} else if ast.IsExpression() {
		return ast.Expression
	} else if ast.IsCopy() {
		return ast.Copy
	} else if ast.IsConst() {
		return ast.Const
	} else if ast.IsCall() {
		return ast.Call
	}
	ast.BaseASTNode = generic_ast.BaseASTNode{}
	ast.ParentNode = nil
	//fmt.Printf("Failed for node %s\n", repr.String(ast))
	panic("Invalid IRStatement type")
}

func feedExpressionIntoIRStatement(node interface{}, base generic_ast.BaseASTNode) *IRStatement {
	if exitStmt, ok := node.(*IRExit); ok {
		return &IRStatement{
			BaseASTNode: base,
			Exit:        exitStmt,
		}
	} else if ifStmt, ok := node.(*IRIf); ok {
		return &IRStatement{
			BaseASTNode: base,
			If:          ifStmt,
		}
	} else if expr, ok := node.(*IRExpression); ok {
		return &IRStatement{
			BaseASTNode: base,
			Expression:  expr,
		}
	} else if expr, ok := node.(*IRCopy); ok {
		return &IRStatement{
			BaseASTNode: base,
			Copy:        expr,
		}
	} else if expr, ok := node.(*IRConst); ok {
		return &IRStatement{
			BaseASTNode: base,
			Const:       expr,
		}
	} else if expr, ok := node.(*IRCall); ok {
		return &IRStatement{
			BaseASTNode: base,
			Call:        expr,
		}
	}
	panic("Invalid irstatement type")
}

func (ast *IRStatement) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	mappedIRStatement := feedExpressionIntoIRStatement(mapper(ast, ast.Body(), context, false), ast.BaseASTNode)
	return mapper(parent, mappedIRStatement, context, true)
}

func (ast *IRStatement) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(ast, ast.Body(), context)
	mapper(parent, ast, context)
}

func (ast *IRStatement) ExpressionType() hindley_milner.ExpressionType { return hindley_milner.E_PROXY }

//

func (ast *IRStatement) BuildFlowGraph(builder cfg.CFGBuilder) {
	if ast.IsEmpty() {
		// Do nothing
	} else if ast.IsExit() {
		builder.BuildNode(ast.Exit)
	} else if ast.IsIf() {
		builder.BuildNode(ast.If)
	} else if ast.IsPhi() {
		builder.BuildNode(ast.Phi)
	} else if ast.IsExpression() {
		builder.BuildNode(ast.Expression)
	} else if ast.IsCopy() {
		builder.BuildNode(ast.Copy)
	} else if ast.IsConst() {
		builder.BuildNode(ast.Const)
	} else if ast.IsCall() {
		builder.BuildNode(ast.Call)
	}
}

func (ast *IRStatement) ResolveTypeOfVar(name string) IRType {
	if ast.IsEmpty() {
		return IR_UNKNOWN
	} else if ast.IsExit() {
		return IR_UNKNOWN
	} else if ast.IsIf() {
		return IR_UNKNOWN
	} else if ast.IsPhi() {
		if ast.Phi.TargetName == name {
			return ast.Phi.Type
		}
	} else if ast.IsExpression() {
		if ast.Expression.TargetName == name {
			return ast.Expression.Type
		}
	} else if ast.IsCopy() {
		if ast.Copy.TargetName == name {
			return ast.Copy.Type
		}
	} else if ast.IsConst() {
		if ast.Const.TargetName == name {
			return ast.Const.Type
		}
	} else if ast.IsCall() {
		if ast.Call.TargetName == name {
			return ast.Call.Type
		}
	}
	return IR_UNKNOWN
}

func (ast *IRStatement) LookupBlock(blockID int) *IRBlock {
	return ast.ParentNode.(*IRBlock).LookupBlock(blockID)
}
