# Lazymaker
Simple CLI tool to easily create/manipulate programming files.

## Requirements
[Go 1.17+](https://golang.org/dl/)

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
### Available commands
You can use three different commands:
- ``lazy create``: creates a file.
- ``lazy compile``: compiles a file.
- ``lazy run``: runs a file.
Run `lazy help` for help of any command.

Examples:
- ``lazy create -t myfile.rs`` will create a directory called ``rs_projects`` along with a file ``myfile.rs`` under `~/Documents/`, and open it with the terminal.
- ``lazy compile program.c otherfile.cpp`` will search for each file under the ``~/Documents/`` directory and all subdirectories and compile them.
- ``lazy run foo.c`` will search for the file under `~/Documents/` and all subdirectories and run it. If the file is not compiled yet, it will do so and then run it.
