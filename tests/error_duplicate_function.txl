/// ERR = error_duplicate_function.txl:6:8: function with name "hello" already exists in the same scope
fn hello() {
    print("Hello")
}

fn hello() {
    print("Duplicate")
}

fn main() {
    hello()
}