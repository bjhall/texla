/// ERR = Value missing for argument "double" (int) of function "test"
fn test(name str, double int, half float) {
   print("Name:", name, "| Double:", double*2, "| Half:", half/2)
}

fn main() {
   test(4, name="10", half=3)
}