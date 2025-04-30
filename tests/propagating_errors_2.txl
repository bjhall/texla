/// ERR = Error from main function: "Cannot divide by zero"
/// OUT = 9

fn divide?(numerator int, denominator int) -> int {
   if denominator == 0 {
      fail "Cannot divide by zero"
   }
   return numerator/denominator
}

fn calc?(x int, y int) -> int {
   return divide(x-1, y-1)?
}

fn do?() {
   a = calc(10,2)?
   print(a)
   b = calc(10,1)?
   print(b)
}

fn main() {
   do()?
}