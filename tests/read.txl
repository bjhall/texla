/// OUT = monkey tail 6 11
/// OUT = tiger tail 5 10
/// OUT = elephant tail 8 13
/// OUT = police tail 6 11
/// OUT = DONE

fn concat(s1 str, s2 str) -> str {
   return s1+" "+s2
}

fn main() {
   read("test_file") -> row {
       a = concat(row, "tail")
       print(a, len(row), a.len())
   }

   print("DONE")
}