// cli - An Extensible, POSIX Compatible Command-Line Argument Parser
// Copyright (c) 2017 Fadhli Dzil Ikram

package cli

import (
	"fmt"
	"io"
	"os"
)

func printUsage(w io.Writer, ctx *Context) {
	var cmd *Command
	var err error
	// Get the current command stack
	if cmd, err = ctx.Stack.Get(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	// Get program name
	prog := ctx.Stack.serialize()
	// Print the usage header
	fmt.Fprintf(w, "\nUsage:  %s [OPTIONS]", prog)
	if len(cmd.Commands) > 0 {
		fmt.Fprintf(w, " COMMAND")
	} else {
		if len(cmd.Arguments) > 0 {
			fmt.Fprintf(w, " %s", cmd.Arguments)
		} else {
			fmt.Fprintf(w, " ARGUMENTS...")
		}
	}
	// Print the description header
	fmt.Fprintf(w, "\n\n%s\n", cmd.Description)
	// Print available options
	if len(cmd.Options) > 0 {
		fmt.Fprintf(w, "\nOptions:\n")
		syntax := make([]string, len(cmd.Options))
		syntaxLen := 0
		for i, opt := range cmd.Options {
			syntax[i] = fmt.Sprintf("%s--%s%s", getShortSyntax(opt, true), opt.Property().Name, getInputSyntax(opt))
			if len(syntax[i]) > syntaxLen {
				syntaxLen = len(syntax[i])
			}
		}
		format := fmt.Sprintf("  %%-%ds%%s\n", syntaxLen+3)
		for i, opt := range cmd.Options {
			fmt.Fprintf(w, format, syntax[i], opt.Property().Usage)
		}
	}
	// Print available commands
	if len(cmd.Commands) > 0 {
		fmt.Fprintf(w, "\nCommands:\n")
		syntaxLen := 0
		for _, subcmd := range cmd.Commands {
			if len(subcmd.Name) > syntaxLen {
				syntaxLen = len(subcmd.Name)
			}
		}
		format := fmt.Sprintf("  %%-%ds%%s\n", syntaxLen+3)
		for _, subcmd := range cmd.Commands {
			fmt.Fprintf(w, format, subcmd.Name, subcmd.Description)
		}
	}
	// Print remarks if available
	if len(cmd.Remarks) > 0 {
		fmt.Fprintf(w, "\n%s\n", cmd.Remarks)
	} else if len(cmd.Commands) > 0 {
		fmt.Fprintf(w, "\nUse '%s COMMAND --help' for more information about a command.\n", prog)
	}
}

func getInputSyntax(opt Option) string {
	if opt.Config().Input {
		return fmt.Sprintf(" %s", opt.Config().InputType)
	}
	return ""
}

func getShortSyntax(opt Option, space bool) string {
	// Write short syntax if available
	if opt.Property().Short != 0 {
		return fmt.Sprintf("-%c, ", opt.Property().Short)
	} else if space {
		return fmt.Sprintf("    ")
	}
	// Write empty string
	return ""
}
