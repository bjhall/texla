/// OUT = 4.5
/// OUT = 20.3
/// OUT = 4
/// OUT = 8

fn double(a int) -> int {
   return a * 2
}

fn mult(a int, b int) -> int {
   return a * b
}

fn half(a float) -> float {
   return a / 2
}

fn main() {
	// Inline coercion
	print(half(mult(3,3)))

	// Coercion of function return
	d = 4.3 + mult(4, 4)
	print(d)

	// Coercion of funciton arguments
	print(mult(2.5, 2.5))

	// Coercion bonanza
	print(half(double(mult(2.1, "4"))))
}