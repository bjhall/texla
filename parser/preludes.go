package parser

func preludeCode(name string) string {
	switch name {
	case "stringToFloat":
		return `
func stringToFloat(s string) float64 {
    f, err := strconv.ParseFloat(s, 64)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Runtime error: string %q cannot be converted to float", s)
        os.Exit(99)
    }
    return f
}
`
	case "stringToInt":
		return `
func stringToInt(s string) int {
    i, err := strconv.Atoi(s)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Runtime error: string %q cannot be converted to integer", s)
        os.Exit(99)
    }
    return i
}
`
	case "intToString":
		return ""
	default:
		panic("Unknown prelude")

	}
}

func preludeImports(name string) []string {
	switch name {
	case "stringToInt", "stringToFloat":
		return []string{"fmt", "strconv", "os"}
	case "intToString":
		return []string{"strconv"}
	default:
		panic("Unknown prelude")
	}
}
