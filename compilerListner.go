package main

import "github.com/karetskiiVO/GOInterpreter/parser"

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

// func (l *GoCompilerListener) Exo
