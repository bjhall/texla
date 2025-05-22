/// OUT = Both are > 5
/// OUT = Any above 10? true

fn main() {
   a = 10
   b = 11
   if a > 5 && b > 5 {
       print("Both are > 5")
   }

   l = a > 10 || b > 10
   print("Any above 10?", l)
}