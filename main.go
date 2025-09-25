package main

import (
	"bytes"
	"strings"
	"flag"
	"fmt"
	"github.com/bjhall/texla/parser"
	"os"
	"os/exec"
)

func main() {

	debugFlag := flag.Bool("debug", false, "Print debug information")
	flag.Parse()

	DEBUG := *debugFlag

	inFiles := flag.Args()
	code, err := os.ReadFile(inFiles[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "I/O ERROR: Could not read input file %q", os.Args[1])
		os.Exit(1)
	}

	tokens, err := parser.Tokenize(string(code) + "\n", 0)
	if err != nil {
		fmt.Fprintf(os.Stderr, "LEXER ERROR: %s\n", err)
		os.Exit(1)
	}

	if DEBUG {
		fmt.Fprint(os.Stderr, "CODE:", string(code[:])+"\n")

		fmt.Println("TOKENS:")
		for idx, token := range tokens {
			fmt.Println(idx, token)
		}
	}

	ast, err := parser.Parse(tokens, inFiles)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
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

	if DEBUG {
		fmt.Println("\nTYPE CHECKED AST:")
		ast.Print(0)
		fmt.Println()
	}

	transpiledCode, err := parser.GenerateCode(typed_ast)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	if DEBUG {
		fmt.Println("GENERATED CODE:")
		lines := strings.Split(transpiledCode, "\n")
		for i := 0; i < len(lines); i++ {
			fmt.Printf("%d: %s\n", i+1, lines[i])
		}
	}

	// Write go code to a file
	err = os.WriteFile("/tmp/a.go", []byte(transpiledCode), 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer os.Remove("/tmp/a.go")

	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command("go", "build", "-o", "/tmp/a", "/tmp/a.go")
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

	cmd = exec.Command("/tmp/a")
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
