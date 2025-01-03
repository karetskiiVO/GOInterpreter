package main

import (
	"fmt"

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

	callInstr := &FunctionCallInstruction{
		program:    l.program,
		functionID: functionID,
		arguments:  slices.Clone(l.instructionStack[len(l.instructionStack)-argumentsCnt : len(l.instructionStack)]),
	}

	l.instructionStack = append(l.instructionStack[:len(l.instructionStack)-argumentsCnt], callInstr)
}

func (l *GoCompilerListener) ExitStringUsing(ctx *parser.StringUsingContext) {
	str := ctx.GetText()
	l.instructionStack = append(l.instructionStack, &StringUsingInstruction{
		program: l.program,
		str:     str[1 : len(str)-1],
	})
}
