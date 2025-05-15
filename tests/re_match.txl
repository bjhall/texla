/// OUT = 0 '^This' matches!
/// OUT = 1 'This$' doesn't match!
/// OUT = 2 'l.ng\ssent' matches!
/// OUT = 3 'T\w{3} is' matches!
/// OUT = 4 '  ' doesn't match!

fn main() {

   haystack = "This is a fairly long sentence"
   patterns = ["^This", "This$", "l.ng\\ssent", "T\\w{3} is", "  "]

   for patterns -> pattern, idx {
      if haystack.match(pattern) {
   	      print(idx, "'"+pattern+"'", "matches!")
   	  } else {
          print(idx, "'"+pattern+"'", "doesn't match!")
      }
   }
}