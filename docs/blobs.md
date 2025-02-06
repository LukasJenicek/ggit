
## Database content

`git cat-file -p` - show object from its database
`git cat-file -p ${commit-hash}`

Example:
```
tree 863610802b71463d76f50a4aca56fffabf3330c5
parent 08c1413da7f6247de7320bf857784088e3743e45
author Lukas Jenicek <lukas.jenicek5@gmail.com> 1738684374 +0100
committer Lukas Jenicek <lukas.jenicek5@gmail.com> 1738684404 +0100

move docs to separate folder
```

Tree - ID of a tree, which represents your whole tree of files when this commit was made.

`git cat-file -p 863610802b71463d76f50a4aca56fffabf3330c5`

Shows:
```
100644 blob f1ca1419cbc72c555c80a490dfbe1485571a3950    .golangci.yml
100644 blob a51a19f57030705ee39b3a93829a6ffe34c841a9    Makefile
100644 blob fe3be0ce89c053d044d86ee5729eea4ff579c713    Readme.md
040000 tree 8b0ea11e2f310333c360ece26f5122d3e0348da1    cmd
040000 tree a52f3f39d964515c0265dbc0e2bef47fc7cfb760    docs
100644 blob 2d05ccc71b5ff29798a1023efd9a99aaa2c84028    go.mod
100644 blob 40b0112572cc9c6beeccfafc6b523647344f1397    go.sum
040000 tree 1fcb567218d3f78288f8ac7fd9c0efe061e67f1e    internal
040000 tree 4b2451128a88eb23539c41d422fe1565659d5057    testdata
```

Each entry in a tree is either another tree ( subdirectory ) or a blob ( regular file )

100644 - regular file, readable, but not executable

`git cat-file -p 40b0112572cc9c6beeccfafc6b523647344f1397`

Shows:
```
github.com/davecgh/go-spew v1.1.1 h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=
github.com/davecgh/go-spew v1.1.1/go.mod h1:J7Y8YcW2NihsgmVo/mv3lAwl/skON4iLHjSsI+c5H38=
github.com/pmezard/go-difflib v1.0.0 h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=
github.com/pmezard/go-difflib v1.0.0/go.mod h1:iKH77koFhYxTK1pcRnkKkqfTogsbg7gZNVY4sRDYZ/4=
github.com/stretchr/testify v1.9.0 h1:HtqpIVDClZ4nwg75+f6Lvsy/wHu+3BoSGCbBAcpTsTg=
github.com/stretchr/testify v1.9.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
github.com/stretchr/testify v1.10.0 h1:Xv5erBjTwe/5IxqUQTdXv5kgmIvbHo3QQyRwhJsOfJA=
github.com/stretchr/testify v1.10.0/go.mod h1:r2ic/lqez/lEtzL7wO/rwa5dbSLXVDPFyf8C91i36aY=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405 h1:yhCVgyC4o1eVCa2tZl7eS0r+SDo693bJlVdllGtEeKM=
gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405/go.mod h1:Co6ibVJAznAaIkqp8huTwlJQCZ016jof/cbN4VW5Yz0=
gopkg.in/yaml.v3 v3.0.1 h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=
gopkg.in/yaml.v3 v3.0.1/go.mod h1:K4uyk7z7BCEPqu6E+C64Yfv1cQ7kz7rIZviUmN+EgEM=
```

Commits point to trees, and trees point to blobs or trees.

## Blobs 
- Stores objects in separate folders to prevent from all objects being in the same directory. Which could be difficult to read for some operating systems.
- Saving disk space by compressing every object using DEFLATE compression algorithm. Implemented by widely used library called `zlib`. It's part of stdlib https://pkg.go.dev/compress/flate. [Wiki](https://en.wikipedia.org/wiki/Deflate)

```
├── objects
│├── 03
│   └── 060df5e568642c331494399bde381766d3b159
│├── 05
│   └── 9eada58c933d536fe575e5d435d5eaed82380f
```

Summary: 
1. Commit hash points to tree `git cat-file -p {commit-hash}`
2. Tree references either different trees ( folders ) or blobs ( files )
3. Files are compressed using deflate algorithm used in gzip, zlib and raw form. In case of git it's zlib what you looking for.
4. `make build` and run `cat .git/objects/f1/ca1419cbc72c555c80a490dfbe1485571a3950 | ./bin/inflate` and it should print blob content in raw form.