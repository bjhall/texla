/// OUT = Hello
/// OUT = Hi !!!
fn say(s str, scream bool) {
   if scream {
       print(s, "!!!")
   }
   if !scream {
       print(s)
   }
}

fn main() {
   say("Hello", false)
   say("Hi", true)
}