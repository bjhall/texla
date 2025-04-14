/// OUT = 9
/// OUT = 16
/// OUT = 49
/// OUT = 4
/// OUT = (3*2)^2 == 36


fn double(x int) -> int {
   return x*2
}

fn main() {
   square = [0, 1, 4, 9, 16, 25, 36, 49, 64, 81, 100]
   a = square[3]
   print(a)

   strIdx = "4"
   b = square[strIdx]
   print(b)

   floatIdx = 7.3
   print(square[floatIdx])

   print(square[1+1])

   if square[double(2+1)] == 36 {
      print("(3*2)^2 == 36")
   }
}