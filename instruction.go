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

func (instr *StringUsingInstruction) Execute(variables map[string]any) error {
	instr.program.stack = append(instr.program.stack, CloneAny(instr.str))
	return nil
}

type IntUsingInstruction struct {
	program *Program
	integer int
}

func (instr *IntUsingInstruction) Execute(variables map[string]any) error {
	instr.program.stack = append(instr.program.stack, instr.integer)
	return nil
}

type VariableUsingInstruction struct {
	program      *Program
	variableName string
}

func (instr *VariableUsingInstruction) Execute(variables map[string]any) error {
	val, ok := variables[instr.variableName]
	if !ok {
		return fmt.Errorf("variable %v not declarated", instr.variableName)
	}

	instr.program.stack = append(instr.program.stack, CloneAny(val))
	return nil
}

type BoolUsingInstruction struct {
	program *Program
	boolVal bool
}

func (instr *BoolUsingInstruction) Execute(variables map[string]any) error {
	instr.program.stack = append(instr.program.stack, instr.boolVal)
	return nil
}

type AssigmentInstruction struct {
	program     *Program
	varName     string
	instruction Instruction
}

func (instr *AssigmentInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)
	err := instr.instruction.Execute(variables)
	if err != nil {
		return err
	}

	if stacklen == len(instr.program.stack) {
		return fmt.Errorf("right hand value has no return statment")
	}
	val, ok := variables[instr.varName]
	if !ok {
		return fmt.Errorf("variable %v undefined", instr.varName)
	}
	if reflect.TypeOf(val) != reflect.TypeOf(instr.program.stack[stacklen]) {
		return fmt.Errorf(
			"mismatcn types expected: %v, actual: %v",
			reflect.TypeOf(val),
			reflect.TypeOf(instr.program.stack[stacklen]))
	}

	variables[instr.varName] = CloneAny(instr.program.stack[stacklen])
	instr.program.stack = instr.program.stack[:stacklen]
	return nil
}

type AddInstruction struct {
	program      *Program
	instructions []Instruction
}

func (instr *AddInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)

	for idx := range instr.instructions {
		instruction := instr.instructions[len(instr.instructions)-1-idx]

		err := instruction.Execute(variables)
		if err != nil {
			return err
		}
	}

	if len(instr.program.stack)-stacklen != len(instr.instructions) {
		return fmt.Errorf(
			"missmatch between return values expected: %v actual: %v",
			len(instr.instructions),
			len(instr.program.stack)-stacklen,
		)
	}

	for idx := 0; idx < len(instr.instructions)-1; idx++ {
		val1 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]
		val2 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]

		sum, err := AddAny(val1, val2)
		if err != nil {
			return err
		}
		instr.program.stack = append(instr.program.stack, sum)
	}

	return nil
}

type MulInstruction struct {
	program      *Program
	instructions []Instruction
}

func (instr *MulInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)

	for idx := range instr.instructions {
		instruction := instr.instructions[len(instr.instructions)-1-idx]

		err := instruction.Execute(variables)
		if err != nil {
			return err
		}
	}

	if len(instr.program.stack)-stacklen != len(instr.instructions) {
		return fmt.Errorf(
			"missmatch between return values expected: %v actual: %v",
			len(instr.instructions),
			len(instr.program.stack)-stacklen,
		)
	}

	for idx := 0; idx < len(instr.instructions)-1; idx++ {
		val1 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]
		val2 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]

		sum, err := MulAny(val1, val2)
		if err != nil {
			return err
		}
		instr.program.stack = append(instr.program.stack, sum)
	}

	return nil
}

type SubInstruction struct {
	program      *Program
	instructions []Instruction
}

func (instr *SubInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)

	for idx := range instr.instructions {
		instruction := instr.instructions[len(instr.instructions)-1-idx]

		err := instruction.Execute(variables)
		if err != nil {
			return err
		}
	}

	if len(instr.program.stack)-stacklen != len(instr.instructions) {
		return fmt.Errorf(
			"missmatch between return values expected: %v actual: %v",
			len(instr.instructions),
			len(instr.program.stack)-stacklen,
		)
	}

	for idx := 0; idx < len(instr.instructions)-1; idx++ {
		val1 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]
		val2 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]

		sum, err := SubAny(val1, val2)
		if err != nil {
			return err
		}
		instr.program.stack = append(instr.program.stack, sum)
	}

	return nil
}

type DivInstruction struct {
	program      *Program
	instructions []Instruction
}

func (instr *DivInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)

	for idx := range instr.instructions {
		instruction := instr.instructions[len(instr.instructions)-1-idx]

		err := instruction.Execute(variables)
		if err != nil {
			return err
		}
	}

	if len(instr.program.stack)-stacklen != len(instr.instructions) {
		return fmt.Errorf(
			"missmatch between return values expected: %v actual: %v",
			len(instr.instructions),
			len(instr.program.stack)-stacklen,
		)
	}

	for idx := 0; idx < len(instr.instructions)-1; idx++ {
		val1 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]
		val2 := instr.program.stack[len(instr.program.stack)-1]
		instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]

		sum, err := DivAny(val1, val2)
		if err != nil {
			return err
		}
		instr.program.stack = append(instr.program.stack, sum)
	}

	return nil
}
