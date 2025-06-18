/// OUT = map[Bonnie:{} Bruce:{} Clyde:{} Fred:{} Johnny:{}]
/// OUT = Clyde is in the set

fn main() {
   names1 = set("Bonnie", "Fred")
   names2 = set("Clyde", "Bruce")
   names = union(names1, names2)
   names.add("Johnny")
   names.add("Clyde")
   print(names)
   if names.has("Clyde") {
       print("Clyde is in the set")
   }
   names.del("Fred")
   if names.has("Fred") {
       print("Fred is in the set")
   }
}