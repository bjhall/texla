/// OUT = 10
/// OUT = 7
/// OUT = 2
/// OUT = 9
/// OUT = 16
/// OUT = 17
/// OUT = 17.799999999999997

fn addi(a int, b int) {
   print(a+b)
}

fn addf(a float, b float) {
   print(a+b)
}

fn main() {

   // Literal float to int coercion
   addi(3.3, 6+1)

   // Float variables to int
   apa = 3.9
   b = 4.1
   addi(apa, b)

   // Literal int to float coercion
   addf(1, 1)

   // Int variables to float
   i = 3
   j = 6
   addf(i, j)

   // String literals to int
   addi("7", "9")

   // String variables to int
   k = "8"
   l = "9"
   addi(k, l)

   // String variables to float
   m = "9.1"
   n = "8.7"
   addf(m, n)

   
}