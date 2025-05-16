/// OUT = Slice has 2 elements

fn main() {
   string = "Test string"
   matches = string.capture("K")

   if matches {
       print("This should not be printed")
   }

   matches.append("Hello")
   matches.append("World")
   if matches {
       print("Slice has", len(matches), "elements")
   }
}