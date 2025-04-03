/// OUT = a is smaller than b
/// OUT = a is different from b
/// OUT = 9.5 + 1 is greater than 10
/// OUT = 9.5 is greater or equal to 5.5
/// OUT = 10.5 is equal to '10.5'
/// OUT = 'one' is not 'two'
/// OUT = 'one' is 'one'

fn main() {
   a = 10
   b = 20
   if a < b {
   	  print("a is smaller than b")
   }

   if a != b {
      print("a is different from b")
   }

   c = 9.5
   if c + 1 > a {
      print(c, "+ 1 is greater than 10")
   }

   if c >= 5.5 {
      print("9.5 is greater or equal to 5.5")
   }

   d = 9.5
   if d+1 == "10.5" {
      print("10.5 is equal to '10.5'")
   }

   str1 = "one"
   str2 = "two"
   if str1 != str2 {
      print("'one' is not 'two'")
   }

   if str1 == "one" {
      print("'one' is 'one'")
   }
}