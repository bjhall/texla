/// ERR = Function "half" can return an error, but it is not handled

fn half?(x int) -> int {
   if x == 1 {
      fail "cannot half 1"
   }
   return x/2
}

fn main() {
   b = half(10)
}