/// OUT = true false
/// OUT = false
/// OUT = a is true
/// OUT = a is true
/// OUT = b is not true
/// OUT = b is false
/// OUT = b is not true
fn main() {
   a = true
   b = false
   print(a, b)
   print(a == b)

   if a == true {
       print("a is true")
   }

   if a {
       print("a is true")
   }

   if b != true {
   	   print("b is not true")
   }

   if b == false {
   	   print("b is false")
   }

   if !b {
       print("b is not true")
   }

}