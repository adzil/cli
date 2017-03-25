// cli - An Extensible, POSIX Compatible Command-Line Argument Parser
// Copyright (c) 2017 Fadhli Dzil Ikram

package cli

import (
	"fmt"
	"os"
	"strings"
)

var ErrorHandler = HandlerFunc(func(ctx *Context) error {
	// Print error information
	fmt.Fprintln(os.Stderr, ctx.Error.Error())
	// Print usage information
	printUsage(os.Stderr, ctx)
	// Exit with error code
	os.Exit(1)
	// There is no error returned to parent
	return nil
})

var OptionParser = OptionHandlerFunc(func(ctx *Context, opt Option) error {
	var arg string
	var err error
	// Get argument without actually popping it
	if arg, err = ctx.Arguments.Get(); err != nil {
		return err
	}
	// Check for option flag
	if len(arg) < 2 || arg[0] != '-' {
		return ErrEndOfOption
	}
	// Check for option terminator
	if arg == "--" {
		ctx.Arguments.Pop()
		return ErrEndOfOption
	}
	// Check for short option support
	if opt.Property().Short != 0 && byte(opt.Property().Short) == arg[1] {
		// Check if short option is input or not
		if opt.Config().Input {
			// Check for value position
			if len(arg) > 2 {
				ctx.Arguments.SetFirst(arg[2:])
			} else {
				ctx.Arguments.Pop()
			}
		} else if len(arg) > 2 {
			// Is not an input, check length and do proper option decomposition
			// Create new argument array
			args := make([]string, 0, len(ctx.Arguments)+len(arg)-3)
			// Decomposition
			for i := 2; i < len(arg); i++ {
				args = append(args, "-"+arg[i:i+1])
			}
			args = append(args, ctx.Arguments[1:]...)
			// Set new argument list
			ctx.Arguments.Set(args)
		} else {
			// Pop from argument list
			ctx.Arguments.Pop()
		}
		// Run the option processor
		return opt.Parse(ctx)
	}
	// Filter out non extended prefix
	if arg[1] != '-' {
		return ErrNextOption
	}
	// Trim arg string and do exact compare
	arg = arg[2:]
	name := opt.Property().Name
	if arg == name {
		ctx.Arguments.Pop()
		return opt.Parse(ctx)
	} else if strings.HasPrefix(arg, name) && len(arg) > len(name) && arg[len(name)] == '=' {
		ctx.Arguments.SetFirst(arg[len(name)+1:])
		return opt.Parse(ctx)
	}
	// Done, switch to next option
	return ErrNextOption
})

type OptionHandler interface {
	Parse(ctx *Context, opt Option) error
}

type OptionHandlerFunc func(ctx *Context, opt Option) error

func (h OptionHandlerFunc) Parse(ctx *Context, opt Option) error {
	return h(ctx, opt)
}

type Handler interface {
	Run(ctx *Context) error
}

type HandlerFunc func(ctx *Context) error

func (h HandlerFunc) Run(ctx *Context) error {
	return h(ctx)
}
