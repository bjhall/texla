/// OUT = 20
/// OUT = 30
/// OUT = 40
/// OUT = 50
/// OUT = Hi John
/// OUT = Hello Jean!!!
/// OUT = Hej Paul!!!
/// OUT = Bonjour Mary
/// OUT = Hi Kate!!!

fn mult(x int, times int = 2) -> int {
   return x * times
}

fn hello(name str, phrase str = "Hi", scream bool = false) {
   if !scream {
      print(phrase, name)
   }
   if scream {
      print(phrase, name+"!!!")
   }
}

fn main() {
	print(mult(10))
	print(mult(10, 3))
	print(mult(10, times=4))
	print(mult(times=5, x=10))

	hello("John")
	hello("Jean", "Hello", true)
	hello(scream = true, phrase = "Hej", name = "Paul")
	hello("Mary", phrase = "Bonjour")
	hello(scream = true, name = "Kate")
}