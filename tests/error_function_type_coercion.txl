/// ERR = Runtime error: string "hello" cannot be converted to integer
fn add(a int, b int) {
   print(a+b)
}

fn main() {
   add("hello", "world")
}