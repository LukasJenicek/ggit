## Representation

### Tree
Save latest commit to clipboard: `git log --pretty=format:"%H" -n 1 | xclip -selection clipboard`
Show cat output `git cat-file -p ${commit hash}`
```
Output:
tree 6fe6d378a7602f20a1520f8f4036249c98f9c92e
parent a6cbc3988f024f6a4eae00527f3dd93301048f2b
author Lukas Jenicek <lukas.jenicek5@gmail.com> 1738839458 +0100
committer Lukas Jenicek <lukas.jenicek5@gmail.com> 1738839458 +0100
```
Print content of tree: `git cat-file -p 6fe6d378a7602f20a1520f8f4036249c98f9c92e`
```
100644 blob be22bdeea290f39e10c63f4f199a156ba1d9728a    .gitignore
100644 blob e8a66e2beca2c379b99da3d5e407a68feb610bf8    .golangci.yml
100644 blob ba15e464ca9aa8b09025af14ebf8ef3f6d97b643    Makefile
100644 blob 66f164baf0ad030e4427289233a75c7006d5a38d    Readme.md
040000 tree 12741077f4898ab15349e660e45b0fcb5fe4718f    cmd
040000 tree 4a555e59d48269376fc5d3f63f9e23cf316131c2    docs
100644 blob 2d05ccc71b5ff29798a1023efd9a99aaa2c84028    go.mod
100644 blob 40b0112572cc9c6beeccfafc6b523647344f1397    go.sum
040000 tree 1fcb567218d3f78288f8ac7fd9c0efe061e67f1e    internal
040000 tree 4b2451128a88eb23539c41d422fe1565659d5057    testdata
```

Print content using hexdump
`cat .git/objects/6f/e6d378a7602f20a1520f8f4036249c98f9c92e | ./bin/inflate | hexdump -C`
```
00000000  74 72 65 65 20 33 35 31  00 31 30 30 36 34 34 20  |tree 351.100644 |
00000010  2e 67 69 74 69 67 6e 6f  72 65 00 be 22 bd ee a2  |.gitignore.."...|
00000020  90 f3 9e 10 c6 3f 4f 19  9a 15 6b a1 d9 72 8a 31  |.....?O...k..r.1|
00000030  30 30 36 34 34 20 2e 67  6f 6c 61 6e 67 63 69 2e  |00644 .golangci.|
00000040  79 6d 6c 00 e8 a6 6e 2b  ec a2 c3 79 b9 9d a3 d5  |yml...n+...y....|
00000050  e4 07 a6 8f eb 61 0b f8  31 30 30 36 34 34 20 4d  |.....a..100644 M|
00000060  61 6b 65 66 69 6c 65 00  ba 15 e4 64 ca 9a a8 b0  |akefile....d....|
00000070  90 25 af 14 eb f8 ef 3f  6d 97 b6 43 31 30 30 36  |.%.....?m..C1006|
00000080  34 34 20 52 65 61 64 6d  65 2e 6d 64 00 66 f1 64  |44 Readme.md.f.d|
00000090  ba f0 ad 03 0e 44 27 28  92 33 a7 5c 70 06 d5 a3  |.....D'(.3.\p...|
000000a0  8d 34 30 30 30 30 20 63  6d 64 00 12 74 10 77 f4  |.40000 cmd..t.w.|
000000b0  89 8a b1 53 49 e6 60 e4  5b 0f cb 5f e4 71 8f 34  |...SI.`.[.._.q.4|
000000c0  30 30 30 30 20 64 6f 63  73 00 4a 55 5e 59 d4 82  |0000 docs.JU^Y..|
000000d0  69 37 6f c5 d3 f6 3f 9e  23 cf 31 61 31 c2 31 30  |i7o...?.#.1a1.10|
000000e0  30 36 34 34 20 67 6f 2e  6d 6f 64 00 2d 05 cc c7  |0644 go.mod.-...|
000000f0  1b 5f f2 97 98 a1 02 3e  fd 9a 99 aa a2 c8 40 28  |._.....>......@(|
00000100  31 30 30 36 34 34 20 67  6f 2e 73 75 6d 00 40 b0  |100644 go.sum.@.|
00000110  11 25 72 cc 9c 6b ee cc  fa fc 6b 52 36 47 34 4f  |.%r..k....kR6G4O|
00000120  13 97 34 30 30 30 30 20  69 6e 74 65 72 6e 61 6c  |..40000 internal|
00000130  00 1f cb 56 72 18 d3 f7  82 88 f8 ac 7f d9 c0 ef  |...Vr...........|
00000140  e0 61 e6 7f 1e 34 30 30  30 30 20 74 65 73 74 64  |.a...40000 testd|
00000150  61 74 61 00 4b 24 51 12  8a 88 eb 23 53 9c 41 d4  |ata.K$Q....#S.A.|
00000160  22 fe 15 65 65 9d 50 57                           |"..ee.PW|
00000168
```

1. Line: `tree 351.100644 ` => `{object type} {byte length}{null byte}{100644}`
`tree` = `74 72 65 65`
` `    = `20`
`351`  = `33 35 31`
`null byte`    = `00`
`100644` = `31 30 30 36 34 34`
` ` = `20`

2 - 4 Line: `.gitignore.."........?O...k..r.100644 ` => `{filename}{null byte}{objectId}{100644}`

- objectId is 20 bytes

### Blob
Example: `.gitignore`

```
00000000  62 6c 6f 62 20 31 32 00  62 69 6e 2f 0a 76 65 6e  |blob 12.bin/.ven|
00000010  64 6f 72 2f                                       |dor/|
00000014
```

`{object type} {size in bytes}{null byte}{content}`

### Commit
`git cat-file -p 01d4ac8be89add3577cde77a3fb7156cf14fea6c`

```
tree 6fe6d378a7602f20a1520f8f4036249c98f9c92e
parent a6cbc3988f024f6a4eae00527f3dd93301048f2b
author Lukas Jenicek <lukas.jenicek5@gmail.com> 1738839458 +0100
committer Lukas Jenicek <lukas.jenicek5@gmail.com> 1738839458 +0100

util for reading compressed files in git
```

tree - represents the state of your files at that point in the history
author - name and email address of the person who created the commit
commiter - same as author, can only differ if you amend somebody else commit or cherry-pick

followed by blank line and then the commit message

### Computing object id
- It's using SHA-1 hash
- SHA1 always generates 160 bits hash size
- Check this example where you can see how object id is generated: `go run cmd/utils/object-id/main.go`

Lib: https://manpages.ubuntu.com/manpages/trusty/man1/zlib-flate.1.html
Using zlib-flate package: `cat .git/objects/f1/ca1419cbc72c555c80a490dfbe1485571a3950 | zlib-flate -uncompress`

Output:
`blob 3821linters` and rest of the file

Cmd: `cat .git/objects/f1/ca1419cbc72c555c80a490dfbe1485571a3950 | zlib-flate -uncompress | hexdump -C`
Hexdump shows numeric values of all those bytes written in hexadecimal

Output: `00000000  62 6c 6f 62 20 33 38 32  31 00 6c 69 6e 74 65 72  |blob 3821.linter|`

The column on the right, between the | characters, displays the corresponding ASCII character
that each byte represents. If the byte does not represent a printable character in ASCII20 then
hexdump prints a .

So, Git stores blobs by prepending them with the word blob, a space, the length of the blob,
and a null byte, and then compressing the result using zlib.

00 is a null bytes which is typically used to separate bits of information in a binary format.