package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jxwr/php-parser/ast"
	php "github.com/jxwr/php-parser/parser"
)

func main() {
	bytes, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	parser := php.NewParser(string(bytes))
	nodes, _ := parser.Parse()
	for _, node := range nodes {
		switch v := node.(type) {
		case ast.EchoStmt:
			s := v.Expressions[0].(ast.Literal)
			fmt.Println(strings.TrimSpace(s.Value))
		case ast.ExpressionStmt:
			call := v.Expression.(*ast.FunctionCallExpression)
			fmt.Println(call)
		case *ast.IfStmt:
			fmt.Println("if ", v.Condition, " {")
			fmt.Printf("%#v", v.TrueBranch)
			fmt.Println("} else {")
			fmt.Printf("%#v", v.FalseBranch)
			fmt.Println("}")
		default:
			fmt.Printf("%#v\n", node)
		}
	}
}
