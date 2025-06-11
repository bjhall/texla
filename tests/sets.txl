/// OUT = map[Bonnie:{} Clyde:{} Johnny:{}]
/// OUT = Clyde is in the set

fn main() {
   names = set("Bonnie", "Clyde")
   names.add("Johnny")
   names.add("Clyde")
   print(names)
   if names.has("Clyde") {
       print("Clyde is in the set")
   }
   if names.has("Fred") {
       print("Fred is in the set")
   }
}