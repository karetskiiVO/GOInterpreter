package main

import (
	"fmt"
	"reflect"
)

type Function interface {
	Call(args ...any) ([]any, error)
	Name() string
}

type ReturnError struct{}

func (ReturnError) Error() string {
	panic("ReturnError must be handled it is not error")
}

type GenericFunction struct {
	name    string
	handler func(args ...any) error
}

func (gf GenericFunction) Call(args ...any) ([]any, error) {
	return []any{}, gf.handler(args...)
}

func (gf GenericFunction) Name() string {
	return gf.name
}

type InputVariable struct {
	Name string
	Type reflect.Type
}

type IntrpretatedFunction struct {
	inputVariables []InputVariable
	returnType     reflect.Type

	name         string
	instructions []Instruction
}

func NewIntrpretatedFunction(name string) *IntrpretatedFunction {
	return &IntrpretatedFunction{
		inputVariables: make([]InputVariable, 0),
		name:           name,
		instructions:   make([]Instruction, 0),
	}
}

func (f *IntrpretatedFunction) Name() string {
	return f.name
}

func (f *IntrpretatedFunction) Call(args ...any) ([]any, error) {
	if len(args) != len(f.inputVariables) {
		return nil, fmt.Errorf(
			"missmatch betweent count of arguments in function %v, given: %v expected: %v",
			f.name,
			len(args),
			len(f.inputVariables),
		)
	}

	variables := make(map[string]any)

	if f.returnType != nil {
		variables["@result"] = nil
	}

	for i, inputVariable := range f.inputVariables {
		if reflect.TypeOf(args[i]) != f.inputVariables[i].Type {
			return nil, fmt.Errorf(
				"missmatch type for argument %v, given: %v expected: %v",
				inputVariable.Name,
				reflect.TypeOf(args[i]).Name(),
				inputVariable.Type.Name(),
			)
		}

		variables[inputVariable.Name] = CloneAny(args[i])
	}

	for _, instruction := range f.instructions {
		err := instruction.Execute(variables)

		if err == nil {
			continue
		}
		if _, ok := err.(ReturnError); ok {
			break
		}

		return nil, err
	}

	var res []any
	if f.returnType != nil {
		if reflect.TypeOf(variables["@result"]) != f.returnType {
			return nil, fmt.Errorf(
				"wrong return type expected: %v, has: %v",
				f.returnType,
				reflect.TypeOf(variables["@result"]))
		}

		res = append(res, variables["@result"])
	}

	return res, nil
}

func (f *IntrpretatedFunction) RegisterArgument(argument InputVariable) error {
	for _, inputVariable := range f.inputVariables {
		if inputVariable.Name == argument.Name {
			return fmt.Errorf("variable %v double declared", argument.Name)
		}
	}

	f.inputVariables = append(f.inputVariables, argument)

	return nil
}
