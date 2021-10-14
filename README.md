# Lazymaker
CLI tool to easily create/manipulate programming files.

## Requirements
[Go](https://golang.org/dl/)

## Setup
Add $GOPATH to $PATH, if you haven't already:
```bash
 export PATH=$PATH:$(go env GOPATH)/bin
```

Then simply run
```
go get github.com/grbenjamin/lazy
```

## Use

Run `lazy help` to see all available commands.
Examples:
- ``lazy create myfile.cpp`` will create a directory called ``cpp_projects`` along with a file ``myfile.cpp``, on the current root directory.
- ``lazy create -t program.go`` will create the file the same way as explained above, and opens it in the terminal (vim, vi, nano)
- ``lazy create -o foo.rs`` will create the file and open it with the OS preferred application.
