/// OUT = olleh

fn backwards(s str) -> str {
   r = ""
   for [4,3,2,1,0] -> idx {
       r = r + s[idx]
   }
   return r
}

fn main() {
   print("hello".backwards())
}