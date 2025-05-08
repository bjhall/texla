/// OUT = PART: this
/// OUT = PART: is
/// OUT = PART: a
/// OUT = PART: string
/// OUT = Integer part: 1 Fractional part: 75
/// OUT = [1 3]
/// OUT = hello world how are you

fn div(x int, y int) -> float {
   return x / y
}

fn main() {

   a = "this is a string"
   parts = a.split(" ")
   for parts -> part {
   	   print("PART:", part)
   }

   p = div(7, 4).split(".")
   print("Integer part:", p[0], "Fractional part:", p[1])

   l = 123.split("2")
   print(l)

   o = "hello_world_how_are_you"
   print(o.split("_").join(" "))
}