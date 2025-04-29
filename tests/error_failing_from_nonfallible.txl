/// ERR = Cannot use `fail` in non-fallible function

fn test() {
   fail "Not allowed"
}

fn main() {
   test()
}