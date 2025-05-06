/// OUT = The division failed: Division by zero
/// OUT = 1
/// OUT = 2
/// OUT = 3
/// OUT = Result: 0
/// OUT = 2

fn div?(x int, y int) -> float {
   if y == 0 {
       fail("Division by zero")
   }
   return x / y
}


fn main() {
   result = div(10, 0) ? {
       print("The division failed:", err)
	   for 1..3 -> i {
	       print(i)
	   }

   }

   print("Result:", result)

   result = div(10, 5)?
   print(result)
}