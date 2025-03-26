package main

import (
	"fmt"
	"bytes"
	"os"
	"os/exec"
	"github.com/bjhall/texla/parser"
)

func main() {

	code, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "I/O ERROR: Could not read input file %q", os.Args[1])
		os.Exit(1)
	}
	DEBUG := false

	tokens, err := parser.Tokenize(string(code)+"\n")
	if err != nil {
		fmt.Fprintf(os.Stderr, "LEXER ERROR: %s\n", err)
		os.Exit(1)
	}

	if DEBUG {
		fmt.Println("TOKENS:")
		for idx, token := range tokens {
			fmt.Println(idx, token)
		}
	}

	ast, err := parser.Parse(tokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "PARSE ERROR: %s\n", err)
		os.Exit(1)
	}

	if DEBUG {
		fmt.Println("\nAST:")
		ast.Print(0)
		fmt.Println()
	}

	typed_ast, err := parser.CheckTypes(ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	transpiledCode, err := parser.GenerateCode(typed_ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if DEBUG {
		fmt.Println("GENERATED CODE:")
		fmt.Println(transpiledCode)
	}

	// Write go code to a file
	err = os.WriteFile("a.go", []byte(transpiledCode), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command("go", "build", "-o", "a",  "a.go")
    cmd.Stdout = &outbuf
    cmd.Stderr = &errbuf

	err = cmd.Run()
    stdout := outbuf.String()
    stderr := errbuf.String()

	if err != nil {
		fmt.Println("GO COMPILATION ERROR. THIS SHOULD NEVER HAPPEN!")
		fmt.Println(stdout)
		fmt.Println(stderr)
		os.Exit(1)
	}

	cmd = exec.Command("./a")
    cmd.Stdout = &outbuf
    cmd.Stderr = &errbuf
	err = cmd.Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", errbuf.String())
	}

	if DEBUG {
		fmt.Println("OUTPUT:")
	}

	fmt.Print(string(outbuf.String()))
}

//for line in read("test.tsv", header=true, sep="\t") {
//
//}
