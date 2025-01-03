package main

import (
	"fmt"
	"os"

	"github.com/antlr4-go/antlr/v4"
	"github.com/jessevdk/go-flags"
	"github.com/karetskiiVO/GOInterpreter/parser"
)

func main() {
	var options struct {
		Args struct {
			SourceFileName string
		} `positional-args:"yes" required:"1"`
	}

	flagsParser := flags.NewParser(&options, flags.Default&(^flags.PrintErrors))
	_, err := flagsParser.Parse()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	input, err := antlr.NewFileStream(options.Args.SourceFileName)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	lexer := parser.NewGoLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	parser := parser.NewGoParser(stream)

	tree := parser.Program() // Начинаем с корневого узла

	program := NewProgram()

	declarationListner := NewGoDeclarationListener(program)
	antlr.ParseTreeWalkerDefault.Walk(declarationListner, tree)
	if len(declarationListner.Errors) != 0 {
		for _, err := range declarationListner.Errors {
			fmt.Println(err)
		}

		os.Exit(1)
	}

	compileListner := NewGoCompilerListener(program)
	antlr.ParseTreeWalkerDefault.Walk(compileListner, tree)
	if len(compileListner.Errors) != 0 {
		for _, err := range compileListner.Errors {
			fmt.Println(err)
		}

		os.Exit(1)
	}

	err = program.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
