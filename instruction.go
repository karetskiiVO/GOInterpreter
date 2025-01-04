package main

import (
	"fmt"
	"reflect"
	"slices"

	"golang.org/x/exp/maps"
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
		err := argument.Execute(variables)
		if err != nil {
			return err
		}
	}
	args := slices.Clone(instr.program.stack[stacklen:])
	instr.program.stack = instr.program.stack[:stacklen]

	res, err := instr.program.functions[instr.functionID].Call(args...)
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

type BlockInstruction struct {
	program      *Program
	instructions []Instruction
}

func (instr *BlockInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)
	blockVariables := maps.Clone(variables)

	var err error = nil

	for _, instruction := range instr.instructions {
		err = instruction.Execute(blockVariables)
		
		if _, ok := err.(ReturnError); ok {
			break
		}
		if err != nil {
			return err
		}
	}

	for variable := range variables {
		variables[variable] = blockVariables[variable]
	}
	if stacklen > len(instr.program.stack) {
		return fmt.Errorf("wrong stack size")
	}
	instr.program.stack = instr.program.stack[:stacklen]

	return err
}

type IFInstruction struct {
	program   *Program
	statment  Instruction
	than      Instruction
	otherwise Instruction
}

func (instr *IFInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)
	err := instr.statment.Execute(variables)
	if err != nil {
		return err
	}

	if len(instr.program.stack) != stacklen+1 {
		return fmt.Errorf("wrong count of return values of statement")
	}

	statementValue, ok := instr.program.stack[len(instr.program.stack)-1].(bool)
	instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]
	if !ok {
		return fmt.Errorf("statement: %v(type: %v) is not bool", statementValue, reflect.TypeOf(statementValue))
	}

	if statementValue {
		err := instr.than.Execute(variables)
		if err != nil {
			return err
		}
	} else if instr.otherwise != nil {
		err := instr.otherwise.Execute(variables)
		if err != nil {
			return err
		}
	}

	return nil
}

type CompareInstruction struct {
	program     *Program
	lhv, rhv    Instruction
	compareType string
}

func (instr *CompareInstruction) Execute(variables map[string]any) error {
	stacklen := len(instr.program.stack)

	var err error
	err = instr.lhv.Execute(variables)
	if err != nil {
		return err
	}
	err = instr.rhv.Execute(variables)
	if err != nil {
		return err
	}

	if len(instr.program.stack) != stacklen+2 {
		return fmt.Errorf("wrong count of return values of statement")
	}

	rhv := instr.program.stack[len(instr.program.stack)-1]
	instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]
	lhv := instr.program.stack[len(instr.program.stack)-1]
	instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]

	res, err := CompareAny(lhv, rhv, instr.compareType)
	if err != nil {
		return err
	}

	instr.program.stack = append(instr.program.stack, res)

	return nil
}

type FORInstruction struct {
	program  *Program
	statment Instruction
	than     Instruction
}

func (instr *FORInstruction) Execute(variables map[string]any) error {
	statementValue := true
	for {
		if instr.statment != nil {
			ok := false
			stacklen := len(instr.program.stack)
			err := instr.statment.Execute(variables)
			if err != nil {
				return err
			}

			if len(instr.program.stack) != stacklen+1 {
				return fmt.Errorf("wrong count of return values of statement")
			}

			statementValue, ok = instr.program.stack[len(instr.program.stack)-1].(bool)
			instr.program.stack = instr.program.stack[:len(instr.program.stack)-1]
			if !ok {
				return fmt.Errorf("statement: %v(type: %v) is not bool", statementValue, reflect.TypeOf(statementValue))
			}
		}

		if !statementValue {
			break
		}

		err := instr.than.Execute(variables)
		if err != nil {

			if reflect.TypeOf(err) == reflect.TypeOf(BreakError{}) {
				break
			}

			return err
		}

	}

	return nil
}

type BreakInstruction struct{}

func (instr *BreakInstruction) Execute(variables map[string]any) error {
	return BreakError{}
}

type ReturnInstruction struct {
	program    *Program
	expression Instruction
}

func (instr *ReturnInstruction) Execute(variables map[string]any) error {
	hasExpression := (instr.expression != nil)
	_, hasReturnValue := variables["@result"]
	if hasExpression != hasReturnValue {
		return fmt.Errorf("has return value:%v expected:%v", hasExpression, hasReturnValue)
	}

	if hasReturnValue {
		stacklen := len(instr.program.stack)
		err := instr.expression.Execute(variables)
		if err != nil {
			return err
		}

		if len(instr.program.stack) != stacklen+1 {
			return fmt.Errorf("wrong count of return values of statement")
		}

		variables["@result"] = instr.program.stack[len(instr.program.stack)-1]
	}

	return ReturnError{}
}
