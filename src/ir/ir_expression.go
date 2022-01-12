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

type IROperator string
type IROperatorKind string

const (
	IR_OP_KIND_NUMERIC IROperatorKind = "NUMERIC"
	IR_OP_KIND_LOGIC   IROperatorKind = "LOGIC"
	IR_OP_KIND_ANY     IROperatorKind = "ANY"
)

const (
	IR_OP_ADD           IROperator = "Add"
	IR_OP_SUB           IROperator = "Sub"
	IR_OP_MUL           IROperator = "Mul"
	IR_OP_DIV           IROperator = "Div"
	IR_OP_SELF_ADD      IROperator = "SelfAdd"
	IR_OP_SELF_SUB      IROperator = "SelfSub"
	IR_OP_SELF_MUL      IROperator = "SelfMul"
	IR_OP_SELF_DIV      IROperator = "SelfDiv"
	IR_OP_EQ            IROperator = "Eq"
	IR_OP_NOT_EQ        IROperator = "NotEq"
	IR_OP_LT            IROperator = "Lt"
	IR_OP_GT            IROperator = "Gt"
	IR_OP_LTEQ          IROperator = "LtEq"
	IR_OP_GTEQ          IROperator = "GtEq"
	IR_OP_AND           IROperator = "And"
	IR_OP_OR            IROperator = "Or"
	IR_OP_SELF_AND      IROperator = "SelfAnd"
	IR_OP_SELF_OR       IROperator = "SelfOr"
	IR_OP_UNARY_NEG     IROperator = "NumberNeg"
	IR_OP_UNARY_LOG_NEG IROperator = "LogicalNeg"
)

type IROperatorSpecs struct {
	ArgsCount          int
	MappedName         string
	TranslationEnabled bool
	Kind               IROperatorKind
}

var OPERATORS_SPECS = map[IROperator]IROperatorSpecs{
	IR_OP_ADD: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          2,
		MappedName:         "+",
		TranslationEnabled: true,
	},
	IR_OP_SUB: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          2,
		MappedName:         "-",
		TranslationEnabled: true,
	},
	IR_OP_MUL: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          2,
		MappedName:         "*",
		TranslationEnabled: true,
	},
	IR_OP_DIV: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          2,
		MappedName:         "/",
		TranslationEnabled: true,
	},
	IR_OP_SELF_ADD: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          1,
		MappedName:         "+=",
		TranslationEnabled: false,
	},
	IR_OP_SELF_SUB: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          1,
		MappedName:         "-=",
		TranslationEnabled: false,
	},
	IR_OP_SELF_MUL: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          1,
		MappedName:         "*=",
		TranslationEnabled: false,
	},
	IR_OP_SELF_DIV: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          1,
		MappedName:         "/=",
		TranslationEnabled: false,
	},
	IR_OP_EQ: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         "==",
		TranslationEnabled: true,
	},
	IR_OP_NOT_EQ: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         "!=",
		TranslationEnabled: true,
	},
	IR_OP_LT: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         "<",
		TranslationEnabled: true,
	},
	IR_OP_GT: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         ">",
		TranslationEnabled: true,
	},
	IR_OP_LTEQ: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         "<=",
		TranslationEnabled: true,
	},
	IR_OP_GTEQ: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         ">=",
		TranslationEnabled: true,
	},
	IR_OP_AND: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         "&&",
		TranslationEnabled: true,
	},
	IR_OP_OR: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          2,
		MappedName:         "||",
		TranslationEnabled: true,
	},
	IR_OP_SELF_AND: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          1,
		MappedName:         "&=",
		TranslationEnabled: false,
	},
	IR_OP_SELF_OR: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          1,
		MappedName:         "|=",
		TranslationEnabled: false,
	},
	IR_OP_UNARY_NEG: IROperatorSpecs{
		Kind:               IR_OP_KIND_NUMERIC,
		ArgsCount:          1,
		MappedName:         "-",
		TranslationEnabled: true,
	},
	IR_OP_UNARY_LOG_NEG: IROperatorSpecs{
		Kind:               IR_OP_KIND_LOGIC,
		ArgsCount:          1,
		MappedName:         "!",
		TranslationEnabled: true,
	},
}

func CreateIROperator(name string, argsCount int, kind IROperatorKind) IROperator {
	for v, spec := range OPERATORS_SPECS {
		if spec.ArgsCount == argsCount && spec.TranslationEnabled && (kind == IR_OP_KIND_ANY || kind == spec.Kind) {
			return v
		}
	}
	panic(fmt.Sprintf("Failed to translate operator: [%s] with kind %s with argsCount == %d", name, kind, argsCount))
}

type IRExpression struct {
	generic_ast.BaseASTNode
	Type           IRType
	TargetName     string
	Operation      IROperator
	ArgumentsTypes []IRType
	Arguments      []string
	ParentNode     generic_ast.TraversableNode
}

func (ast *IRExpression) Parent() generic_ast.TraversableNode {
	return ast.ParentNode
}

func (ast *IRExpression) OverrideParent(node generic_ast.TraversableNode) {
	ast.ParentNode = node
}

func (ast *IRExpression) Begin() lexer.Position {
	return ast.Pos
}

func (ast *IRExpression) End() lexer.Position {
	return ast.EndPos
}

func (ast *IRExpression) GetNode() interface{} {
	return ast
}

func (ast *IRExpression) GetChildren() []generic_ast.TraversableNode {
	return []generic_ast.TraversableNode{
		generic_ast.MakeTraversableNodeToken(ast, ast.TargetName, ast.Pos, ast.EndPos),
		generic_ast.MakeTraversableNodeToken(ast, string(ast.Operation), ast.Pos, ast.EndPos),
	}
}

func (ast *IRExpression) Print(c *context.ParsingContext) string {
	args := []string{}
	argsTypes := []string{}
	for i, arg := range ast.Arguments {
		args = append(args, arg)
		argsTypes = append(argsTypes, string(ast.ArgumentsTypes[i]))
	}
	return utils.PrintASTNode(c, ast, "%s %s = (%s) %s(%s)", ast.Type, ast.TargetName, strings.Join(argsTypes, ","), ast.Operation, strings.Join(args, ","))
}

//

func (ast *IRExpression) Map(parent generic_ast.Expression, mapper generic_ast.ExpressionMapper, context generic_ast.VisitorContext) generic_ast.Expression {
	return mapper(parent, &IRExpression{
		BaseASTNode: ast.BaseASTNode,
		Operation:   ast.Operation,
		Arguments:   ast.Arguments,
		TargetName:  ast.TargetName,
		ParentNode:  parent.(generic_ast.TraversableNode),
	}, context, true)
}

func (ast *IRExpression) Visit(parent generic_ast.Expression, mapper generic_ast.ExpressionVisitor, context generic_ast.VisitorContext) {
	mapper(parent, ast, context)
}

func (ast *IRExpression) Body() generic_ast.Expression {
	return hindley_milner.Batch{}
}

func (ast *IRExpression) GetAssignedVariables(wantMembers bool, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRExpression) GetUsedVariables(vars cfg.VariableSet, visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	for _, arg := range ast.Arguments {
		vars.Add(cfg.NewVariable(arg, nil))
	}
	return vars
}

func (ast *IRExpression) GetDeclaredVariables(visitedMap map[generic_ast.TraversableNode]struct{}) cfg.VariableSet {
	return cfg.NewVariableSet(cfg.NewVariable(ast.TargetName, nil))
}

func (ast *IRExpression) RenameVariables(substUsed, substDecl cfg.VariableSubstitution) {
	for i, arg := range ast.Arguments {
		ast.Arguments[i] = substUsed.Replace(arg)
	}
	ast.TargetName = substDecl.Replace(ast.TargetName)
}
