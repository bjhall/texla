/// ERR = error_nongenerator_with_arrow.txl:3:15: cannot put "->" after non-generator function "len"
fn main() {
   len("aa") -> i {
     print("A")
   }
}