/// OUT = 100
/// OUT = 101
/// OUT = 102
/// OUT = 103
/// OUT = 1
/// OUT = 2
/// OUT = 3
/// OUT = 4
/// OUT = 1 2
/// OUT = 1 3
/// OUT = 2 3
/// OUT = 2 4
/// OUT = 3 4
/// OUT = 3 5

fn main() {

   for 100..103 -> a {
   	   print(a)
   }

   a = 2
   for a/2..a*2 -> i {
       print(i)
   }

   for 1..3 -> i {
       for i+1..i+2 -> j {
	       print(i, j)
	   }
   }
}