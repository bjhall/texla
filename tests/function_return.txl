/// OUT = 100
/// OUT = 30

fn double(a int) -> int {
   return a * 2
}

fn mult(a int, b int) -> int {
   return a * b
}

fn main() {
    // Variable arguments
	a = 5
    b = 20
    result = mult(a, b)
	print(result)

	// Inline
    print(double(mult(5, 3)))
}