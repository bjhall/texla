/// OUT = 1
/// OUT = 2
/// OUT = 3
/// OUT = 4
/// OUT = 2
/// OUT = 4
/// OUT = 6
/// OUT = 8
/// OUT = hello
/// OUT = world
/// OUT = 115
/// OUT = 116
/// OUT = 114

fn double(x int) -> int {
    return x*2
}

fn main() {
   list = [1,2,3,4]
   chars = ["1", "2", "3", "4"]

   for list -> i {
      print(i)
   }

   // Coerce control variable
   for chars -> char {
   	   print(double(char))
   }

   // Iterate slice literal
   for ["hello", "world"] -> word {
   	   print(word)
   }

   // Iterate string
   a = "str"
   for a -> c {
   	   print(c)
   }
}