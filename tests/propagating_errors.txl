/// OUT = success
/// OUT = 10
/// ERR = Error from main function: "cannot half 1"
fn test?(x int) {
   if x > 10 {
       fail "Too large"
   }
   print("success")
}

fn half?(x int) -> int {
   if x == 1 {
      fail "cannot half 1"
   }
   return x/2
}

fn double?(x int) -> int {
   if x == 0 {
      fail "cannot double 0"
   }
   return x*2
}

fn main() {
   test(10)?

   a = 10
   print(a.half()?.double()?)

   b = half(1)?
}