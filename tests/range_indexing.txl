/// OUT = [1 2]
/// OUT = [0 1]
/// OUT = [3 4 5 6]
/// OUT = lo w
/// OUT = worl

fn double(x int) -> int {
    return x*2
}

fn main() {
   elements = ["0", "1", "2", "3", "4", "5", "6"]
   print(elements[1..3])
   print(elements[..2])
   print(elements[3..])

   text = "Hello world"
   print(text[4-1..7])
   print(text[double(3)..double(5)])
}