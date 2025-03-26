package main

import (
	"fmt"
	"os"
	"strings"
	"bufio"
)

func main() {
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("Could not open file `enums`")
		os.Exit(1)
	}
	defer file.Close()
	fmt.Println("package parser\n")
	scanner := bufio.NewScanner(file)
	firstMember := true
	enumType := ""
	functionStr := ""
	typeName := ""
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		if line[0] == '-' {
			if enumType != "" {
				fmt.Println(")\n\n")
				fmt.Println(functionStr+"\n\tdefault: return \"???\"\n\t}\n}\n\n")
				functionStr = ""
			}
			h := strings.Split(line[1:]," ")
			if len(h) != 2 {
				fmt.Println(h)
				fmt.Println("Header needs name and type")
				os.Exit(1)
			}
			firstMember = true
			enumType = h[0]
			fmt.Printf("type %s %s\n", h[0], h[1])
			fmt.Println("const (")
			typeName = h[0]
			functionStr += fmt.Sprintf("func (s %s) String() string {\n\tswitch s {\n", h[0])
		} else {
			name := strings.TrimSuffix(line, "\n")
			niceName := strings.ReplaceAll(name, typeName, "")
			fmt.Printf("\t%s", name)
			functionStr += fmt.Sprintf("\tcase %s: return \"%s\"\n", name, niceName)
			if firstMember {
				fmt.Printf(" %s = iota", enumType)
				firstMember = false
			}
			fmt.Println()
		}
	}
	fmt.Println(")")
	fmt.Println(functionStr+"\n\tdefault: return \"???\"\n\t}\n}\n")
}
