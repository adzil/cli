// cli - An Extensible, POSIX Compatible Command-Line Argument Parser
// Copyright (c) 2017 Fadhli Dzil Ikram

package cli

import (
	"os"
)

type Command struct {
	Options      []Option
	Commands     []Command
	Name         string
	Description  string
	Arguments    string
	Remarks      string
	Handler      Handler
	ErrorHandler Handler
}

func Cmd(name string, description string) Command {
	return Command{
		Name:        name,
		Description: description,
	}
}

func Subcmd(name string, description string, arguments string) Command {
	return Command{
		Name:        name,
		Description: description,
		Arguments:   arguments,
	}
}

func (c *Command) SetCommands(cmd ...Command) {
	c.Commands = cmd
}

func (c *Command) SetOptions(opt ...Option) {
	c.Options = opt
}

func (c *Command) SetHandler(fn HandlerFunc) {
	c.Handler = fn
}

func (c *Command) exec(ctx *Context) error {
	var err error

	// Add current command to context stack
	ctx.Stack.Push(c)
	// Run option initialisation function if available
	for _, opt := range c.Options {
		if initOpt, ok := opt.(OptionInit); ok {
			initOpt.Init(ctx)
		}
	}
	// Process option until end of option flag no argument left
	for len(c.Options) > 0 && err != ErrEndOfOption && len(ctx.Arguments) > 0 {
		// Iterate over option list
		for _, opt := range c.Options {
			if err = ctx.Application.OptionHandler.Parse(ctx, opt); err == nil || err == ErrEndOfOption {
				break
			} else if err != ErrNextOption {
				return err
			}
		}
		// Checks if option not found
		if err == ErrNextOption {
			var opt string
			// Get unknown options
			if opt, err = ctx.Arguments.Get(); err != nil {
				return err
			}
			ctx.Arguments.Pop()
			// Return syntax error
			return NewError("Unknown option '%s'", opt)
		}
	}
	// Run command handler if defined
	if c.Handler != nil {
		if err = c.Handler.Run(ctx); err != nil {
			return err
		}
	}
	// Execute subcommand finder if defined
	if len(c.Commands) > 0 {
		var cmd string
		if cmd, err = ctx.Arguments.Get(); err == ErrEmptyContextArgument {
			printUsage(os.Stderr, ctx)
			os.Exit(0)
			return nil
		} else if err != nil {
			return err
		}
		ctx.Arguments.Pop()
		// Iterate over subcommand list
		for _, command := range c.Commands {
			if command.Name == cmd {
				// Run subcommand Exec()
				return command.Exec(ctx)
			}
		}
		// Tell user that we don't know the command
		return NewError("Unknown command '%s'", cmd)
	}
	// Return execution without error
	return nil
}

func (c *Command) Exec(ctx *Context) error {
	var err error
	// Run the internal executor and get the error result
	if err = c.exec(ctx); err != nil {
		// Check if error handler is exists
		if c.ErrorHandler != nil {
			// Write error to application Context
			ctx.Error = err
			// Run the error handler
			return c.ErrorHandler.Run(ctx)
		}
	}
	// No error handler present, pass error to parent
	return err
}

type Application struct {
	Command
	Version       string
	OptionHandler OptionHandler
}

func NewApp(name string, description string, arguments string, version string) Application {
	return Application{
		Command: Command{
			Name:        name,
			Description: description,
			Arguments:   arguments,
		},
		Version: version,
	}
}

func NewCmdApp(name string, description string, version string) Application {
	return Application{
		Command: Command{
			Name:        name,
			Description: description,
		},
		Version: version,
	}
}

func (a *Application) Run(osArgs []string) error {
	// Initialize context
	ctx := NewContext(a, osArgs[1:])
	// Set the default option handler if not set
	if a.OptionHandler == nil {
		a.OptionHandler = OptionParser
	}
	// Set the default error handler if not set
	if a.ErrorHandler == nil {
		a.ErrorHandler = ErrorHandler
	}
	// Run the top level Exec()
	return a.Exec(ctx)
}
