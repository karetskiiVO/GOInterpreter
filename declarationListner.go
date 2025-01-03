package main

import (
	"github.com/karetskiiVO/GOInterpreter/parser"
)

type GoDeclarationListener struct {
	*parser.BaseGoListener

	program *Program
	Errors  []error
}

func NewGoDeclarationListener(program *Program) *GoDeclarationListener {
	return &GoDeclarationListener{
		program: program,
	}
}

func (l *GoDeclarationListener) EnterFunctionDefinition(ctx *parser.FunctionDefinitionContext) {
	res := NewIntrpretatedFunction(ctx.NAME().GetText())

	if ctx.Arguments() != nil {
		for i := range ctx.Arguments().AllNAME() {
			varName := ctx.Arguments().AllNAME()[i].GetText()
			varType := ctx.Arguments().AllTypename()[i].NAME().GetText()

			inputVariable := InputVariable{}

			inputVariable.Name = varName

			var err error
			inputVariable.Type, err = ReflectType(varType)
			if err != nil {
				l.Errors = append(l.Errors, err)
			}

			err = res.RegisterArgument(inputVariable)
			if err != nil {
				l.Errors = append(l.Errors, err)
			}
		}
	}

	if ctx.Typename() != nil {
		var err error
		res.returnType, err = ReflectType(ctx.Typename().GetText())

		if err != nil {
			l.Errors = append(l.Errors, err)
		}
	}

	err := l.program.RegisterFunction(res)
	if err != nil {
		l.Errors = append(l.Errors, err)
	}
}
