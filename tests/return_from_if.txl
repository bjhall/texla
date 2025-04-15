/// OUT = true
/// OUT = false

fn test(x int) -> bool {
   if x > 10 {
   	  return true
   }
   return false
}

fn main() {
   print(test(11))
   print(test(9))
}