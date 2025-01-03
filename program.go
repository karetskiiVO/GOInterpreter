package main

import "fmt"

type Program struct {
	functions  []Function
	functionID map[string]int
}

func NewProgram() *Program {
	res := &Program{
		functions:  make([]Function, 0),
		functionID: map[string]int{},
	}

	res.RegisterFunction(GenericFunction{
		name: "print",
		handler: func(args ...any) error {
			print(args)
			return nil
		},
	})
	res.RegisterFunction(GenericFunction{
		name: "println",
		handler: func(args ...any) error {
			print(args)
			return nil
		},
	})
	res.RegisterFunction(GenericFunction{
		name: "panic",
		handler: func(args ...any) error {
			if len(args) != 1 {
				return fmt.Errorf("the \"panic\" function has an incorrect number of arguments")
			}

			return fmt.Errorf("%v", args[0])
		},
	})

	return res
}

func (prog *Program) RegisterFunction(function Function) error {
	if _, ok := prog.functionID[function.Name()]; ok {
		return fmt.Errorf("function %v already defined", function.Name())
	}

	prog.functionID[function.Name()] = len(prog.functions)
	prog.functions = append(prog.functions, function)

	return nil
}

func (prog *Program) Execute() error {
	id, ok := prog.functionID["main"]
	if !ok {
		return fmt.Errorf("there is no 'main'")
	}
	res, err := prog.functions[id].Call()
	if len(res) != 0 {
		return fmt.Errorf("'main' can't have return value")
	}
	
	return err
}
