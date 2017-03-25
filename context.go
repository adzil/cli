// cli - An Extensible, POSIX Compatible Command-Line Argument Parser
// Copyright (c) 2017 Fadhli Dzil Ikram

package cli

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

var ErrEmptyContextStack = errors.New("CLI: Context stack is empty")

var ErrEmptyContextArgument = errors.New("CLI: Context argument slice empty")

var ErrInvalidStruct = errors.New("CLI: Receiver struct should be a pointer")

const ContextStackGrowSize = 4

type ContextStack []*Command

func (s *ContextStack) Push(cmd *Command) {
	// Check if stack is empty or full
	if len(*s) == cap(*s) {
		// Create new larger slice and copy all of its contents
		tmp := make(ContextStack, len(*s), cap(*s)+ContextStackGrowSize)
		copy(tmp, *s)
		*s = tmp
	}
	// Add command to stack
	*s = append(*s, cmd)
}

func (s ContextStack) Get() (*Command, error) {
	// Check if context stack is empty
	if len(s) == 0 {
		return nil, ErrEmptyContextStack
	}
	// Return the most current command stack
	return s[len(s)-1], nil
}

func (s ContextStack) serialize() string {
	var buf bytes.Buffer
	// Print all command name to Buffer
	for _, stack := range s {
		buf.WriteString(stack.Name)
		buf.WriteByte(' ')
	}
	// Check if buffer actually contains serialized data
	if buf.Len() > 0 {
		return buf.String()[:buf.Len()-1]
	}
	// Otherwise, return empty string
	return ""
}

type ContextArguments []string

func (a ContextArguments) Get() (string, error) {
	// Check if context arguments is empty
	if len(a) == 0 {
		return "", ErrEmptyContextArgument
	}
	// Return the first argument list entry
	return a[0], nil
}

func (a ContextArguments) SetFirst(arg string) error {
	// Check if argument list is empty
	if len(a) == 0 {
		return ErrEmptyContextArgument
	}
	// Set first argument
	a[0] = arg
	return nil
}

func (a *ContextArguments) Set(args []string) {
	// Replace argument list
	*a = args
}

func (a *ContextArguments) Pop() error {
	if len(*a) == 0 {
		return ErrEmptyContextArgument
	}
	// Slice the argument list
	(*a) = (*a)[1:]
	return nil
}

type Context struct {
	Application *Application
	Arguments   ContextArguments
	Error       error
	Options     map[string]interface{}
	Stack       ContextStack
}

func (c Context) Map(v interface{}) error {
	// Find reflection value of receiver struct
	recv := reflect.ValueOf(v).Elem()
	// Iterate on every field in struct
	for i := 0; i < recv.NumField(); i++ {
		// Get destination field info
		dst := recv.Field(i)
		tag := recv.Type().Field(i).Tag.Get("cli")
		// Check if destination field tag is exists on options
		if opt, ok := c.Options[tag]; ok {
			// If destination is a pointer, initialize its value and change the
			// destination to the actual element
			if dst.Kind() == reflect.Ptr {
				dst.Set(reflect.New(dst.Type().Elem()))
				dst = dst.Elem()
			}
			// Get source reflection value
			src := reflect.ValueOf(opt)
			// Check for destination and source type mismatch
			if src.Kind() != dst.Kind() {
				return fmt.Errorf(
					"CLI: Cannot assign mismatched type option '%s' (%s) to field '%s' (%s)",
					tag, src.Kind().String(), recv.Type().Field(i).Name, dst.Kind().String(),
				)
			}
			// Set struct value
			dst.Set(src)
		}
	}
	return nil
}

func NewContext(app *Application, args []string) *Context {
	return &Context{
		Application: app,
		Arguments:   ContextArguments(args),
		Options:     make(map[string]interface{}),
	}
}
