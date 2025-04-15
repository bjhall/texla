/// OUT = 20
/// OUT = 15

fn sub(x int, y int) -> int {
   return x-y
}

fn mult(x int, fact int = 2, plus int = 0) -> int {
	return x * fact + plus
}

fn one() -> int {
   return 1
}

fn main() {
    a = 10
	b = a.mult(plus = 10, fact = 1).sub(10).mult()
   	print(b)

	print(one().mult(10, 5))
}