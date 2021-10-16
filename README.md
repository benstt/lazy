# Lazymaker
Simple CLI tool to easily create/manipulate programming files.

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
- ``lazy create myfile.rs`` will create a directory called ``rs_projects`` along with a file ``myfile.rs``, under `~/Documents/`.
- ``lazy create -t program.go`` will create the file the same way as explained above, and opens it in the terminal (vim, vi, nano)
- ``lazy compile foo.c`` will search for the file under `~/Documents/` and all subdirectories and compile it.
