/// ERR = Runtime error: string "hello" cannot be converted to integer
fn sum(a int, b int) {
   print(a+b)
}

fn main() {
   sum("hello", "world")
}