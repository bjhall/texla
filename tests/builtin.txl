/// OUT = abc123 6
/// OUT = [1 2 3 4 5 6] 6

fn sum(a float, b float) -> float {
   return a + b
}

fn main() {
   a = "a"
   a.append("bc")
   a.append(123)
   l = a.len()
   print(a, l)

   b = [1, 2, 3]
   append(b, 4)
   b.append(sum(2,3))
   b.append("6")
   print(b, b.len())
}