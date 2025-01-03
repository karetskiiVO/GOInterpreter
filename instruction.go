package main

import (
	"fmt"
	"reflect"
)

type Instruction interface {
	Execute(variables map[string]any) error
}

type DefineVariableInstruction struct {
	Name string
	Type reflect.Type
}

func (instr *DefineVariableInstruction) Execute(variables map[string]any) error {
	if _, ok := variables[instr.Name]; ok {
		return fmt.Errorf("varible %v has already defined", instr.Name)
	}

	variables[instr.Name] = NewVariable(instr.Type)

	println(reflect.TypeOf(variables[instr.Name]).String())
	return nil
}

type FunctionCallInstruction struct {
	program *Program

	functionID int
	arguments  []Instruction
}

func (instr *FunctionCallInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)

	for _, argument := range instr.arguments {
		argument.Execute(variables)
	}

	res, err := instr.program.functions[instr.functionID].Call(instr.program.stack[stacklen:]...)
	instr.program.stack = append(instr.program.stack, res...)
	return err
}

type StringUsingInstruction struct {
	program *Program
	str     string
}

func (instr StringUsingInstruction) Execute(variables map[string]any) error {
	instr.program.stack = append(instr.program.stack, instr.str)
	return nil
}
