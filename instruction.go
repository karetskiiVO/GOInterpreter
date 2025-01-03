package main

import "reflect"

type Instruction interface {
	Execute(variables map[string]any, context []any) error
}

type DefineVariableInstruction struct {
	Name string
	Type reflect.Type
}

func (instr *DefineVariableInstruction) Execute(variables map[string]any, context []any) error {
	println(instr.Name, instr.Type.String())
	println("implement")
	return nil
}
