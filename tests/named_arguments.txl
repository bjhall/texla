/// OUT = Name: Test1 | Double: 20 | Half: 1.5
/// OUT = Name: Test2 | Double: 4 | Half: 1.5
/// OUT = Name: Test3 | Double: 10 | Half: 2
/// OUT = Name: Test4 | Double: 10 | Half: 10

fn test(name str, double int, half float) {
   print("Name:", name, "| Double:", double*2, "| Half:", half/2)
}

fn main() {
   test("Test1", double=10, half=3)
   test(half="3", name="Test2", double=2)
   test(name="Test3", double=10/2, half=2+2)
   test("Test4", 5, 20)
}