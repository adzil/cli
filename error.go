// cli - An Extensible, POSIX Compatible Command-Line Argument Parser
// Copyright (c) 2017 Fadhli Dzil Ikram

package cli

import (
	"fmt"
)

type Error struct {
	Context error
}

func (e Error) Error() string {
	return e.Context.Error()
}

func NewError(format string, i ...interface{}) Error {
	return Error{Context: fmt.Errorf(format, i...)}
}

func NewRequireError(opt Option) Error {
	return NewError(
		"%s argument required in option %s--%s", opt.Config().InputType,
		getShortSyntax(opt, false), opt.Property().Name,
	)
}

func NewTypeError(opt Option) Error {
	return NewError(
		"Invalid argument type in option %s--%s (expect %s)", getShortSyntax(opt, false),
		opt.Property().Name, opt.Config().InputType,
	)
}
