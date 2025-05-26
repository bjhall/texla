/// OUT = 1 1
/// OUT = 2 * 10 = 20
/// OUT = Name: test.person
/// OUT = Domain: email
/// OUT = TLD: com
/// OUT = 1234560123 [1 1]

fn double(x int) -> int {
   return x * 2
}


fn main() {
   b = a = 1
   print(a, b)

   value = 10
   if a = double(value) < 30 {
       print("2 *", value, "=", a)
   }

   string = "test.person@email.com"

   if matches = string.capture("^(.+)@(\\w+)\\.(\\w+)$") {
      print("Name:", matches[0])
      print("Domain:", matches[1])
      print("TLD:", matches[2])
   }

   if matches = string.capture("HELLO") {
      print("This should not be printed")
   }

    numbers = "1234560123"
    if ones = numbers.find("1") && ones.len() == 2 {
        print(numbers, ones)
    }
}