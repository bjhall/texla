/// OUT = a-b-c
/// OUT = 0:1:2:3:4:5:6:7:8:9:10
/// OUT = 0.1 < 0.2 < 0.3 < 0.4

fn main() {
   a = ["a", "b", "c"]
   print(a.join("-"))

   b = 0..10
   print(b.join(":"))

   c = [0.1, 0.2, 0.3, 0.4]
   print(c.join(" < "))
}