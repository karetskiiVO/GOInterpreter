package main

import (
	"fmt"
	"strconv"

	"github.com/karetskiiVO/GOInterpreter/parser"
	"golang.org/x/exp/slices"
)

type GoCompilerListener struct {
	*parser.BaseGoListener

	instructionStack []Instruction
	program          *Program
	Errors           []error
}

func NewGoCompilerListener(program *Program) *GoCompilerListener {
	return &GoCompilerListener{
		program: program,
	}
}

func (l *GoCompilerListener) EnterFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	l.instructionStack = make([]Instruction, 0)
}
func (l *GoCompilerListener) ExitFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	function := l.program.functions[l.program.functionID[ctx.NAME().GetText()]]

	if intrpretedFunction, ok := function.(*IntrpretatedFunction); ok {
		intrpretedFunction.instructions = l.instructionStack
	}
}

func (l *GoCompilerListener) ExitVariableDefinition(ctx *parser.VariableDefinitionContext) {
	Type, err := ReflectType(ctx.Typename().GetText())

	if err != nil {
		l.Errors = append(l.Errors, err)
		return
	}

	l.instructionStack = append(l.instructionStack, &DefineVariableInstruction{
		Name: ctx.NAME().GetText(),
		Type: Type,
	})
}

func (l *GoCompilerListener) ExitCallExpression(ctx *parser.CallExpressionContext) {
	functionID, ok := l.program.functionID[ctx.NAME().GetText()]
	if !ok {
		l.Errors = append(l.Errors, fmt.Errorf("function '%v' undefined", ctx.NAME().GetText()))
		return
	}

	argumentsCnt := len(ctx.AllExpression())

	instruction := &FunctionCallInstruction{
		program:    l.program,
		functionID: functionID,
		arguments:  slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], instruction)
}

func (l *GoCompilerListener) ExitStringUsing(ctx *parser.StringUsingContext) {
	str := ctx.GetText()
	l.instructionStack = append(l.instructionStack, &StringUsingInstruction{
		program: l.program,
		str:     str[1 : len(str)-1],
	})
}

func (l *GoCompilerListener) ExitNumberUsing(ctx *parser.NumberUsingContext) {
	integer, err := strconv.Atoi(ctx.GetText())
	if err != nil {
		l.Errors = append(l.Errors, err)
		return
	}

	l.instructionStack = append(l.instructionStack, &IntUsingInstruction{
		program: l.program,
		integer: integer,
	})
}

func (l *GoCompilerListener) ExitVariableUsing(ctx *parser.VariableUsingContext) {
	l.instructionStack = append(l.instructionStack, &VariableUsingInstruction{
		program:      l.program,
		variableName: ctx.GetText(),
	})
}

func (l *GoCompilerListener) ExitBoolUsing(ctx *parser.BoolUsingContext) {
	l.instructionStack = append(l.instructionStack, &BoolUsingInstruction{
		program: l.program,
		boolVal: ctx.GetText() == "true",
	})
}

func (l *GoCompilerListener) ExitAssigment(ctx *parser.AssigmentContext) {
	instruction := l.instructionStack[len(l.instructionStack)-1]
	l.instructionStack[len(l.instructionStack)-1] = &AssigmentInstruction{
		program:     l.program,
		varName:     ctx.NAME().GetText(),
		instruction: instruction,
	}
}

func (l *GoCompilerListener) ExitExpressionAdd(ctx *parser.ExpressionAddContext) {
	argumentsCnt := len(ctx.AllExpressionSub())

	if argumentsCnt == 1 {
		return
	}

	instruction := &AddInstruction{
		program:      l.program,
		instructions: slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], instruction)
}

func (l *GoCompilerListener) ExitExpressionSub(ctx *parser.ExpressionSubContext) {
	argumentsCnt := len(ctx.AllExpressionMul())

	if argumentsCnt == 1 {
		return
	}

	instruction := &SubInstruction{
		program:      l.program,
		instructions: slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], instruction)
}

func (l *GoCompilerListener) ExitExpressionMul(ctx *parser.ExpressionMulContext) {
	argumentsCnt := len(ctx.AllExpressionDiv())

	if argumentsCnt == 1 {
		return
	}

	instruction := &MulInstruction{
		program:      l.program,
		instructions: slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], instruction)
}

func (l *GoCompilerListener) ExitExpressionDiv(ctx *parser.ExpressionDivContext) {
	argumentsCnt := len(ctx.AllExpressionLogic())

	if argumentsCnt == 1 {
		return
	}

	instruction := &DivInstruction{
		program:      l.program,
		instructions: slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], instruction)
}

func (l *GoCompilerListener) ExitExpressionLogic(ctx *parser.ExpressionLogicContext) {
	if ctx.ExpressionLogic() == nil {
		return
	}

	instruction := &NotInstruction{
		program:     l.program,
		instruction: l.instructionStack[len(l.instructionStack)-1],
	}
	l.instructionStack[len(l.instructionStack)-1] = instruction
}

func (l *GoCompilerListener) ExitExpressionLogicOr(ctx *parser.ExpressionLogicOrContext) {
	argumentsCnt := len(ctx.AllExpressionLogicAnd())

	if argumentsCnt == 1 {
		return
	}

	instruction := &DivInstruction{
		program:      l.program,
		instructions: slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], instruction)
}

func (l *GoCompilerListener) ExitExpressionLogicAnd(ctx *parser.ExpressionLogicAndContext) {
	argumentsCnt := len(ctx.AllCompareExpression())

	if argumentsCnt == 1 {
		return
	}

	instruction := &DivInstruction{
		program:      l.program,
		instructions: slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], instruction)
}

func (l *GoCompilerListener) ExitBlock(ctx *parser.BlockContext) {
	instructionCnt := len(ctx.AllLine())
	instructions := slices.Clone(l.instructionStack[len(l.instructionStack)-instructionCnt : len(l.instructionStack)])

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-instructionCnt], &BlockInstruction{
		program:      l.program,
		instructions: instructions,
	})
}

func (l *GoCompilerListener) ExitExpressionIF(ctx *parser.ExpressionIFContext) {
	res := &IFInstruction{}
	res.program = l.program

	if ctx.ExpressionELSE() != nil {
		res.otherwise = l.instructionStack[len(l.instructionStack)-1]
		l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]
	}

	res.than = l.instructionStack[len(l.instructionStack)-1]
	l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]

	res.statment = l.instructionStack[len(l.instructionStack)-1]
	l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]

	l.instructionStack = append(l.instructionStack, res)
}

func (l *GoCompilerListener) ExitCompareExpression(ctx *parser.CompareExpressionContext) {
	if len(ctx.AllSimpleExpresion()) == 1 {
		return
	}

	res := &CompareInstruction{}
	res.program = l.program
	res.compareType = ctx.COMPARETOKEN().GetText()

	res.rhv = l.instructionStack[len(l.instructionStack)-1]
	l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]

	res.lhv = l.instructionStack[len(l.instructionStack)-1]
	l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]

	l.instructionStack = append(l.instructionStack, res)
}

func (l *GoCompilerListener) ExitExpressionFOR(ctx *parser.ExpressionFORContext) {
	res := &FORInstruction{}
	res.program = l.program

	res.than = l.instructionStack[len(l.instructionStack)-1]
	l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]

	if ctx.Expression() != nil {
		res.statment = l.instructionStack[len(l.instructionStack)-1]
		l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]
	}

	l.instructionStack = append(l.instructionStack, res)
}

func (l *GoCompilerListener) ExitBreak(ctx *parser.BreakContext) {
	l.instructionStack = append(l.instructionStack, &BreakInstruction{})
}

func (l *GoCompilerListener) ExitFunctionReturn(ctx *parser.FunctionReturnContext) {
	res := &ReturnInstruction{
		program: l.program,
	}

	if ctx.Expression() != nil {
		res.expression = l.instructionStack[len(l.instructionStack)-1]
		l.instructionStack = l.instructionStack[:len(l.instructionStack)-1]
	}

	l.instructionStack = append(l.instructionStack, res)
}
