/// ERR = PARSE ERROR: Function with name "hello" already exists in the same scope
fn hello() {
    print("Hello")
}

fn hello() {
    print("Duplicate")
}

fn main() {
    hello()
}