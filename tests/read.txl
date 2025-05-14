/// OUT = monkey tail 6 11
/// OUT = tiger tail 5 10
/// OUT = elephant tail 8 13
/// OUT = police tail 6 11
/// OUT = 0 monkey
/// OUT = 1 tiger
/// OUT = 2 elephant
/// OUT = 3 police

fn concat(s1 str, s2 str) -> str {
   return s1+" "+s2
}

fn main() {
   read("test_file") -> row {
       a = concat(row, "tail")
       print(a, len(row), a.len())
   }

   read("test_file") -> row, idx {
       print(idx,row)
   }

}