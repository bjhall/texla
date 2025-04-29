/// ERR = Function "half" is not fallible, do not put ? after the call to it

fn half(x int) -> int {
   return x/2
}

fn main() {
   b = half(10)?
}