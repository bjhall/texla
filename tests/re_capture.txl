/// OUT = 0 is
/// OUT = 1 a

fn main() {

   text = "this isis a string not a list"

   matches = text.capture("(is) (a)")
   for matches -> match, idx {
       print(idx, match)
   }

}