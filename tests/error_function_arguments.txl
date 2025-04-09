/// ERR = Wrong number of arguments to add, expected 2, got 1
/// ERR = Value missing for argument "b" (Int) of function "add"

fn add(a int, b int) {
    print(a+b)
}


fn main() {
    add(10)
}