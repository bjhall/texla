/// OUT = 1 : monkey
/// OUT = 2 : tiger
/// OUT = 3 : elephant
/// OUT = 4 : police

fn main() {
   code = slurp("test_file").split("\n")
   for code -> line, idx {
       print(idx+1, ":", line)
   }
}