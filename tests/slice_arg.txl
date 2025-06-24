/// OUT = 15
/// OUT = [0 1 2 3 4 5 6 7 8 9 10]

fn sum(list []int) -> int {
   s = 0
   for list -> e {
      s = s + e
   }
   return s
}

fn double(length int) -> []int {
   out = [0]
   for 1..length -> i {
       out.append(i)
   }
   return out
}

fn main() {
   list = [1,2,3,4,5]
   a = sum(list)
   print(a)

   b = double(10)
   print(b)
}