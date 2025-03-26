package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"path/filepath"
)

func main() {
	sourceFiles, err := filepath.Glob("*.txl")
	if err != nil {
		fmt.Println("Not able go get test files...")
		os.Exit(1)
	}

	passCount := 0
	for i, sourceFile := range sourceFiles {
		fmt.Printf("TEST %d: %s: ", i, sourceFile)
		if test(sourceFile) {
			fmt.Println("PASS")
			passCount += 1
		}
	}

	fmt.Printf("\nSUMMARY: passed %d/%d tests\n", passCount, len(sourceFiles))
	if passCount < len(sourceFiles) {
		os.Exit(1)
	}
	os.Exit(0)
}

func test(sourceFile string) bool {

	// Collect expected STDOUT from the header
	expectedOut, expectedErr, err := collectExpectations(sourceFile)
	if err != nil {
		fmt.Printf("Error collecting expectations: %v\n", err)
		os.Exit(1)
	}

	// Run the script and collect the observed output
	observedOut, observedErr, err := runProgram("texla", sourceFile)
	if err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}


	outPass, errPass := compareOutputs(expectedOut, expectedErr, observedOut, observedErr)
	if outPass && errPass {
		return true
	}

	fmt.Println("FAIL")

	if !outPass {
		fmt.Println("        Expected output:")
		for i, output := range expectedOut {
			fmt.Printf("           %d: %s\n", i+1, output)
		}
		fmt.Println("\n        Observed output:")
		for i, output := range observedOut {
			fmt.Printf("           %d: %s\n", i+1, output)
		}
	}

	if !errPass {
		fmt.Println("        Expected error:")
		for i, output := range expectedErr {
			fmt.Printf("           %d: %s\n", i+1, output)
		}
		fmt.Println("\n        Observed error:")
		for i, output := range observedErr {
			fmt.Printf("           %d: %s\n", i+1, output)
		}
	}
	return false
}

func collectExpectations(filename string) ([]string, []string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	var expectedOutputs []string
	var expectedErrors []string
	outRegex := regexp.MustCompile(`^/// OUT\s*=\s*(.*)$`)
	errRegex := regexp.MustCompile(`^/// ERR\s*=\s*(.*)$`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		outMatches := outRegex.FindStringSubmatch(line)
		if len(outMatches) > 1 {
			expectedOutputs = append(expectedOutputs, outMatches[1])
		}
		errMatches := errRegex.FindStringSubmatch(line)
		if len(errMatches) > 1 {
			expectedErrors = append(expectedErrors, errMatches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	return expectedOutputs, expectedErrors, nil
}

func runProgram(program, inputFile string) ([]string, []string, error) {
	cmd := exec.Command(program, inputFile)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	_ = cmd.Run()
	/*if err != nil {
		fmt.Println(err)
		return []string{}, []string{}, err
	}*/

	errStr := strings.TrimSuffix(stderr.String(), "\n")
	outStr := strings.TrimSuffix(stdout.String(), "\n")

	outLines := []string{}
	errLines := []string{}
	if outStr != "" {
		outLines = strings.Split(outStr, "\n")
	}
	if errStr != "" {
		errLines = strings.Split(errStr, "\n")
	}
	return outLines, errLines, nil
}

func compareOutputs(expectedOut []string, expectedErr []string, observedOut []string, observedErr []string) (bool, bool) {

	outPass := true
	errPass := true

	if len(expectedOut) != len(observedOut) {
		outPass = false
	} else {
		for i, _ := range expectedOut {
			if expectedOut[i] != observedOut[i] {
				outPass = false
			}
		}
	}


	if len(expectedErr) != len(observedErr) {
		errPass = false
	} else {
		for i, _ := range expectedErr {
			if expectedErr[i] != observedErr[i] {
				errPass = false
			}
		}
	}

	return outPass, errPass
}
