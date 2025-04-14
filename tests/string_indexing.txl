/// OUT = C
/// OUT = G
/// OUT = B

fn double(x int) -> int {
   return x*2
}

fn main() {
   s = "ABCDEFG"

   a = s[2]
   print(a)

   print(s[double(2+1)])

   print(s["1"])
}