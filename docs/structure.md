## .git directory

1. .git/config - repository configuration
2. .git/description - name of the repository; it is used by gitweb
3. .git/HEAD - contains a reference to the current commit, either using commit id or symbolic reference to the current branch
4. .git/info - exclude files ( .git/info/exclude )
5. .git/hooks - git hooks
6. .git/objects - directory forms git's database, it's where  GIT stores all content ( source code and other assets  )
7. .git/refs - stores various kinds of pointers into .git/objects database. For example, files in
-.git/refs/heads store the latest commit on each local branch ,
-.git/refs/remotes store the latest commit on each remote branch
-.git/refs/tags stores tags.
8. .git/COMMIT_EDITMSG - file used to compose commit messages.
9. .git/index - cache storing information about each file in the current commit. updated whenever you add new files using `git add`
10. .git/logs - contains references, something that points to a commit, like HEAD or a branch name

Config example:
```
[core]
repositoryformatversion = 0 ( file format )
filemode = true ( whether file is executable etc )
bare = false ( not a bare repository, it's a repository where the user will edit the working copy of files and create commits )
logallrefupdates = true ( The reflog is enabled, meaning that all changes to files in .git/refs are logged in .git/logs )
```

On macOS there would be two more
```
ignorecase = true
precomposeunicode = true
```
