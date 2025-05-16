/// OUT = a
/// OUT = c
/// OUT = d
fn main() {

   list = ["a", "b", "c", "d", "e"]
   for list -> letter {
        if letter == "b" { continue }
		if letter == "e" { break }
		print(letter)
   }
}