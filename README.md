Command-Line Parser Library
============================

An extensible, POSIX compatible command-line argument parser with no external
dependencies for Golang.

## Work in Progress and TODOs

This documentation is currently in progress. Incomplete API docs is expected.

Add documentation to the Golang code.

## Usage

Install the package with

```
go get github.com/adzil/cli
```

And write a couple line of codes.

```go
package main

import (
	"fmt"
	"os"

	"github.com/adzil/cli"
)

func main() {
	// Setup CLI app
	app := cli.NewApp("myapp", "My new app written in Golang", "USERNAME", "0.1.0")
	app.SetOptions(cli.HelpOption, cli.VersionOption)
	app.SetHandler(func(ctx *cli.Context) error {
		var user string
		var err error
		// Get username argument from user
		if user, err = ctx.Arguments.Get(); err != nil {
			return cli.NewError("Required USERNAME argument in program.")
		}
		// Print back input username to user
		fmt.Println("Hello,", user, "!")
		return nil
	})
	// Run the application
	app.Run(os.Args)
}
```

Then, build and run your code.

```
$ go build
$ ./myapp world
Hello, world !
```

You can ask the CLI for help by using `--help` option.

```
$ ./myapp --help

Usage:  myapp [OPTIONS] USERNAME

My new app written in Golang

Options:
  -h, --help      Usage help
  -v, --version   Show program version
```

## Add New CLI Options

(documentation in progress)

## Add New CLI Commands

(documentation in progress)

## Bugs and Feature Requests

Please open new issue at GitHub.

## License

This code is MIT licensed. Copyright (c) 2017 Fadhli Dzil Ikram.
