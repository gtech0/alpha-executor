package service

import (
	"alpha-executor/model"
	"fmt"
	"io"
)

func GenerateAST(reader io.Reader) {
	lexer := model.NewLexer(reader)
	output := lexer.Lex()
	//parser := NewParser(output)

	for _, token := range output {
		fmt.Printf("%d\t%s\n", token.Type, token.Value)
	}
}
