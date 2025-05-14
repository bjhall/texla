/// OUT = test	1	2	3
/// OUT = 
/// OUT = test2	3	6	9
/// OUT = 
/// OUT = test3	9	27	81
/// OUT = 
/// OUT = 0 HEADER: [test 1 2 3]
/// OUT = 1 ROW: [test2 3 6 9]
/// OUT = 2 ROW: [test3 9 27 81]

fn main() {
   read("tsv_test", chomp=false) -> l {
       print(l)
   }

   read("tsv_test", sep="\t") -> s, idx {
       if idx == 0 {
	       print(idx, "HEADER:", s)
	   } else {
	       print(idx, "ROW:", s)
	   }
   }
}