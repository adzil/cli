// cli - An Extensible, POSIX Compatible Command-Line Argument Parser
// Copyright (c) 2017 Fadhli Dzil Ikram

package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const ListOptionGrowSize = 8

var ErrEndOfOption = errors.New("CLI: End of option")
var ErrNextOption = errors.New("CLI: Next option")

var HelpOption = HandlerOption{
	OptionProperty: OptionProperty{
		Name:  "help",
		Short: 'h',
		Usage: "Usage help",
	},
	Handler: HandlerFunc(func(ctx *Context) error {
		printUsage(os.Stdout, ctx)
		os.Exit(0)
		return nil
	}),
}

var VersionOption = HandlerOption{
	OptionProperty: OptionProperty{
		Name:  "version",
		Short: 'v',
		Usage: "Show program version",
	},
	Handler: HandlerFunc(func(ctx *Context) error {
		fmt.Println(ctx.Application.Version)
		os.Exit(0)
		return nil
	}),
}

type OptionProperty struct {
	Name  string
	Short rune
	Usage string
}

type OptionConfig struct {
	Input     bool
	InputType string
}

type Option interface {
	Property() OptionProperty
	Config() OptionConfig
	Parse(ctx *Context) error
}

type OptionInit interface {
	Init(ctx *Context)
}

type BoolOption struct {
	OptionProperty
}

func (o BoolOption) Property() OptionProperty {
	return o.OptionProperty
}

func (o BoolOption) Config() OptionConfig {
	return OptionConfig{Input: false}
}

func (o BoolOption) Parse(ctx *Context) error {
	// Set bool option in context
	ctx.Options[o.Name] = true
	// Return with no error
	return nil
}

func OptBool(name string, short rune, usage string) BoolOption {
	return BoolOption{
		OptionProperty: OptionProperty{
			Name:  name,
			Short: short,
			Usage: usage,
		},
	}
}

type TBoolOption struct {
	OptionProperty
}

func (o TBoolOption) Property() OptionProperty {
	return o.OptionProperty
}

func (o TBoolOption) Config() OptionConfig {
	return OptionConfig{Input: false}
}

func (o TBoolOption) Init(ctx *Context) {
	// Set default value in context
	ctx.Options[o.Name] = true
}

func (o TBoolOption) Parse(ctx *Context) error {
	// Set bool option in context
	ctx.Options[o.Name] = false
	// Return with no error
	return nil
}

func OptTBool(name string, short rune, usage string) TBoolOption {
	return TBoolOption{
		OptionProperty: OptionProperty{
			Name:  name,
			Short: short,
			Usage: usage,
		},
	}
}

type HandlerOption struct {
	OptionProperty
	OptionConfig
	Handler Handler
}

func (o HandlerOption) Property() OptionProperty {
	return o.OptionProperty
}

func (o HandlerOption) Config() OptionConfig {
	return o.OptionConfig
}

func (o HandlerOption) Parse(ctx *Context) error {
	return o.Handler.Run(ctx)
}

type StringOption struct {
	OptionProperty
}

func (o StringOption) Property() OptionProperty {
	return o.OptionProperty
}

func (o StringOption) Config() OptionConfig {
	return OptionConfig{Input: true, InputType: "string"}
}

func (o StringOption) Parse(ctx *Context) error {
	var err error
	var val string
	// Get value from arguments
	if val, err = ctx.Arguments.Get(); err != nil {
		return NewRequireError(o)
	}
	ctx.Arguments.Pop()
	// Store value in context
	ctx.Options[o.Property().Name] = val
	return nil
}

func OptString(name string, short rune, usage string) StringOption {
	return StringOption{
		OptionProperty: OptionProperty{
			Name:  name,
			Short: short,
			Usage: usage,
		},
	}
}

type ListOption struct {
	OptionProperty
}

func (o ListOption) Property() OptionProperty {
	return o.OptionProperty
}

func (o ListOption) Config() OptionConfig {
	return OptionConfig{Input: true, InputType: "list"}
}

func (o ListOption) Parse(ctx *Context) error {
	var err error
	var val string
	// Get value from arguments
	if val, err = ctx.Arguments.Get(); err != nil {
		return NewRequireError(o)
	}
	ctx.Arguments.Pop()
	// Get list from context
	var list []string
	if src, ok := ctx.Options[o.Property().Name].([]string); ok {
		list = src
	}
	// Resize list if too small
	if len(list) == cap(list) {
		newlist := make([]string, len(list), cap(list)+ListOptionGrowSize)
		copy(newlist, list)
		list = newlist
	}
	// Append newly added entity
	list = append(list, val)
	// Update the string array in context
	ctx.Options[o.Property().Name] = list
	return nil
}

func OptList(name string, short rune, usage string) ListOption {
	return ListOption{
		OptionProperty: OptionProperty{
			Name:  name,
			Short: short,
			Usage: usage,
		},
	}
}

type IntOption struct {
	OptionProperty
}

func (o IntOption) Property() OptionProperty {
	return o.OptionProperty
}

func (o IntOption) Config() OptionConfig {
	return OptionConfig{Input: true, InputType: "int"}
}

func (o IntOption) Parse(ctx *Context) error {
	var err error
	var val string
	// Get value from arguments
	if val, err = ctx.Arguments.Get(); err != nil {
		return NewRequireError(o)
	}
	ctx.Arguments.Pop()
	// Try to parse the value
	var dst int
	if dst, err = strconv.Atoi(val); err != nil {
		return NewTypeError(o)
	}
	// Store value in context
	ctx.Options[o.Property().Name] = dst
	return nil
}

func OptInt(name string, short rune, usage string) IntOption {
	return IntOption{
		OptionProperty: OptionProperty{
			Name:  name,
			Short: short,
			Usage: usage,
		},
	}
}
