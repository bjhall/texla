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

	case "createRange":
		return `
func ___createRange(from int, to int) []int {
    a := make([]int, to-from+1)
    for i := range a {
        a[i] = from + i
    }
    return a
}
`
	case "joinIntSlice":
		return `
func ___joinIntSlice(slice []int, sep string) string {
    parts := make([]string, len(slice))
    for i, part := range slice {
        parts[i] = strconv.Itoa(part)
    }
    return strings.Join(parts, sep)
}
`

	case "joinFloatSlice":
		return `
func ___joinFloatSlice(slice []float64, sep string) string {
    parts := make([]string, len(slice))
    for i, part := range slice {
        parts[i] = strconv.FormatFloat(part, 'f', -1, 64)
    }
    return strings.Join(parts, sep)
}`

	case "handleNonPropagatableError":
		return `
func ___handleNonPropagatableError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error from main function: %q\n", err.Error());
        os.Exit(1)
    }
}
`

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
	case "createRange":
		return []string{}
	case "joinIntSlice", "joinFloatSlice":
		return []string{"strings", "strconv"}
	case "handleNonPropagatableError":
		return []string{"os"}
	default:
		panic("Unknown prelude")
	}
}
