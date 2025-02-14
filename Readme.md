# GGIT - Git in go

1. [Git directory structure](docs/folder-structure.md)
2. [Structure and storing Blobs](docs/data-representation.md)



├── hello.txt
├── libs
│   ├── haha.txt
│   └── internal
│       └── internal.txt
└── world.txt
--------------------------------------------
entries := map[string]Object

entries["hello.txt"] = Entry("hello.txt")
entries["libs"] = Tree(
    NewEntry("haha.txt"),
)
entries["libs/internal"] = Tree(
    NewEntry("internal.txt")
)
entries["world.txt"] = Entry("world.txt")

func Build(r *Tree, entries map[string]Object) *Tree {
    for p, e := range entries { 
        if e.IsEntry(){
            r.AddEntry(e)
        }

        if e.IsTree(){
            r.AddEntry(e)
            r = Build(e, entries)
        }
    }
}

=> 

r := Tree(
    Entry("hello.txt"),
    Tree(
        Entry("libs/haha.txt"),
        Tree(
            Entry("libs/internal/internal.txt")      
        ),
    ),
    Entry("world.txt")
)

func Save(root *Tree) (string, error) {
    for _, e := range root.Entries() {
        if e.Entry() {
            oid := db.Save(e)
            e.Oid = oid
        }

        if e.Tree() {
            oid, _ := Save(e)
            e.Oid = oid
        }
    }
}

