/// OUT = str is not empty
/// OUT = str2 is empty
/// OUT = int1 is not 0
/// OUT = int2 is 0
fn main() {
   str = "string"
   if str {
       print("str is not empty")
   }

   str2 = ""
   if !str2 {
   	  print("str2 is empty")
   }

   int1 = 3
   if int1 {
   	  print("int1 is not 0")
   }

   int2 = 0
   if !int2 {
       print("int2 is 0")
   }
}