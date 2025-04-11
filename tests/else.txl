/// OUT = 10 is less than 20
/// OUT = 25 is greater than 20
/// OUT = Exactly 20!

fn check(val int) {
   if val > 20 {
      print(val, "is greater than 20")
   } else if val < 20 {
      print(val, "is less than 20")
   } else {
      print("Exactly 20!")
   }
}

fn main() {
   check(10)
   check(25)
   check(20)
}