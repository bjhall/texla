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
	case "regexMatch":
		return `
func ___regexMatch(haystack string, regexpStr string) bool {
    matched, err := regexp.MatchString(regexpStr, haystack)
    if err != nil {
        return false
    }
    return matched
}
`

	case "regexCapture":
		return `
func ___regexCapture(haystack string, regexpStr string) []string {
    regex, err := regexp.Compile(regexpStr)
    if err != nil {
        return []string{}
    }
    matches := regex.FindStringSubmatch(haystack)
    if len(matches) > 0 {
        return matches[1:]
    }
    return matches
}
`

	case "regexFind":
		return `
func ___regexFind(haystack string, regexpStr string) []string {
    regex, err := regexp.Compile(regexpStr)
    if err != nil {
        return []string{}
    }
    matches := regex.FindAllString(haystack, -1)
    return matches
}
`

	case "slurpFile":
		return `
func ___slurpFile(path string) string {
    b, err := os.ReadFile(path)
    if err != nil {
        panic("Cannot read file")
    }
    return string(b)
}
`
	case "setContains":
		return `
func ___setContains[T comparable](haystack map[T]struct{}, needle T) bool {
    _, found := haystack[needle]
    return found
}
`
	case "setUnion":
		return `
func ___setUnion[T comparable](set1 map[T]struct{}, set2 map[T]struct{}) map[T]struct{} {
    union := map[T]struct{}{}
    for v := range set1 {
        union[v] = struct{}{}
    }
    for v := range set2 {
        union[v] = struct{}{}
    }
    return union
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
	case "createRange", "setContains", "setUnion":
		return []string{}
	case "joinIntSlice", "joinFloatSlice":
		return []string{"strings", "strconv"}
	case "handleNonPropagatableError", "slurpFile":
		return []string{"os"}
	case "regexMatch", "regexCapture", "regexFind":
		return []string{"regexp"}
	default:
		panic("Unknown prelude")
	}
}
