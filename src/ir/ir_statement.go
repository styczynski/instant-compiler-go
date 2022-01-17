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
	MacroCall  *IRMacroCall  `| @@`
	ParentNode generic_ast.TraversableNode

	VarIn   cfg.VariableSet
	VarOut  cfg.VariableSet
	Reached cfg.ReachedBlocks
	meta    IRMeta

	Comment string
}

func (ast *IRStatement) SetTargetAllocationConstraints(targetName string, cons IRAllocationConstraints) *IRStatement {
	targetCons := IRAllocationContraintsMap{}
	targetCons[targetName] = cons
	return ast.SetTargetAllocationConstraintsMap(targetCons)
}

func (ast *IRStatement) SetTargetAllocationConstraintsMap(targetCons IRAllocationContraintsMap) *IRStatement {
	if len(targetCons) > 1 {
		panic(fmt.Sprintf("Invalid allocation constraints: %s", targetCons.String()))
	}
	ast.meta.AllocationContraints = targetCons
	return ast
}

func (ast *IRStatement) GetAllocationTargetContraints() (string, IRAllocationConstraints) {
	for k, v := range ast.meta.AllocationContraints {
		return k, v
	}
	return "", IRAllocationConstraints{}
}

func (ast *IRStatement) SetComment(format string, values ...interface{}) *IRStatement {
	ast.Comment = fmt.Sprintf(format, values...)
	return ast
}

func (ast *IRStatement) CopyDataForAllocationShadow(stmt *IRStatement) *IRStatement {
	ast.VarIn = stmt.VarIn
	ast.VarOut = stmt.VarOut
	ast.BaseASTNode = ast.BaseASTNode
	ast.ParentNode = ast.ParentNode
	ast.Reached = ast.Reached
	return ast
}

func (ast *IRStatement) CopyDataForAllocation(stmt *IRStatement) *IRStatement {
	ast.meta.TargetAllocation = stmt.meta.TargetAllocation
	ast.meta.AllocationContraints = stmt.meta.AllocationContraints
	ast.meta.ContextAllocation = nil //stmt.meta.ContextAllocation
	ast.VarIn = stmt.VarIn
	ast.VarOut = stmt.VarOut
	ast.BaseASTNode = ast.BaseASTNode
	ast.ParentNode = ast.ParentNode
	ast.Reached = ast.Reached
	return ast
}

func (ast *IRStatement) SetAllocationInfo(targetAlloc IRAllocationMap, contextAlloc IRAllocationMap) *IRStatement {
	if len(targetAlloc) > 1 {
		panic(fmt.Sprintf("Invalid allocation info: %s", targetAlloc.String()))
	}
	ast.meta.TargetAllocation = targetAlloc
	ast.meta.ContextAllocation = contextAlloc
	return ast
}

func (ast *IRStatement) TryToGetAllocationTarget() (string, IRAllocation, bool) {
	for k, v := range ast.meta.TargetAllocation {
		return k, v, true
	}
	return "", nil, false
}

func (ast *IRStatement) GetAllocationTarget() (string, IRAllocation) {
	name, alloc, ok := ast.TryToGetAllocationTarget()
	if !ok {
		panic("Missing allocation info")
	}
	return name, alloc
}

func (ast *IRStatement) GetAllocationContext() IRAllocationMap {
	return ast.meta.ContextAllocation
}

func (ast *IRStatement) SetFlowAnalysisProps(
	VarIn cfg.VariableSet,
	VarOut cfg.VariableSet,
	Reached cfg.ReachedBlocks,
) *IRStatement {
	ast.VarIn = VarIn
	ast.VarOut = VarOut
	ast.Reached = Reached
	return ast
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

func WrapIRMacroCall(ast *IRMacroCall) *IRStatement {
	ret := &IRStatement{
		MacroCall: ast,
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
		!ast.IsMacroCall() &&
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

func (ast *IRStatement) IsMacroCall() bool {
	return ast.MacroCall != nil
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
	} else if ast.IsMacroCall() {
		return []generic_ast.TraversableNode{ast.MacroCall}
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
	return utils.PrintASTNode(c, ast, "%s%s %s %s %s (begin %s)", strings.Repeat("  ", c.BlockDepth), irstatement, ast.VarIn.String(), ast.VarOut.String(), ast.meta.String(), ast.BaseASTNode.Begin())
	//return utils.PrintASTNode(c, ast, "%s%s %s", strings.Repeat("  ", c.BlockDepth), irstatement, ast.meta.AllocationContraints)
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
	} else if ast.IsMacroCall() {
		ret = ast.MacroCall.Print(c)
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

	if len(ast.Comment) > 0 {
		ret = fmt.Sprintf("%s ; %s", ret, ast.Comment)
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
	} else if ast.IsMacroCall() {
		return ast.MacroCall
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
	} else if expr, ok := node.(*IRMacroCall); ok {
		return &IRStatement{
			BaseASTNode: base,
			MacroCall:   expr,
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
	} else if ast.IsMacroCall() {
		builder.BuildNode(ast.MacroCall)
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
	} else if ast.IsMacroCall() {
		if ast.MacroCall.Var == name {
			return ast.MacroCall.Type
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
