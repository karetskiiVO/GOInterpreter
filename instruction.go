package main

import (
	"fmt"
	"reflect"
)

type Instruction interface {
	Execute(variables map[string]any, context []any) error
}

type DefineVariableInstruction struct {
	Name string
	Type reflect.Type
}

func (instr *DefineVariableInstruction) Execute(variables map[string]any, context []any) error {
	if _, ok := variables[instr.Name]; ok {
		return fmt.Errorf("varible %v has already defined", instr.Name)
	} 

	variables[instr.Name] = reflect.New(instr.Type)

	return nil
}

