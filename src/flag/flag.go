// Package flag implements new github.com/spf13/pflag-compatible flag types.
package flag

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// TODO: this file name doesn't really make sense anymore

func (fs *FlagSet) EnumUnsafe(name string, options []fmt.Stringer, defaultVal fmt.Stringer, usage string) *fmt.Stringer {
	return FlagSetEnum[fmt.Stringer](fs, name, options, defaultVal, usage)
}

func FlagSetEnum[T fmt.Stringer](fs *FlagSet, name string, options []T, defaultVal T, usage string) *T {
	v := enumInternal[T](name, options, defaultVal, usage)
	pflag.Var(v, name, usage)
	return &v.v
}

func Enum[T fmt.Stringer](name string, options []T, defaultVal T, usage string) *T {
	return FlagSetEnum[T](CommandLine, name, options, defaultVal, usage)
}

func (fs *FlagSet) EnumPUnsafe(name, shorthand string, options []fmt.Stringer, defaultVal fmt.Stringer, usage string) *fmt.Stringer {
	return FlagSetEnumP[fmt.Stringer](fs, name, shorthand, options, defaultVal, usage)
}

func FlagSetEnumP[T fmt.Stringer](fs *FlagSet, name, shorthand string, options []T, defaultVal T, usage string) *T {
	v := enumInternal[T](name, options, defaultVal, usage)
	pflag.VarP(v, name, shorthand, usage)
	return &v.v
}

func EnumP[T fmt.Stringer](name, shorthand string, options []T, defaultVal T, usage string) *T {
	return FlagSetEnumP[T](CommandLine, name, shorthand, options, defaultVal, usage)
}

func (fs *FlagSet) PosEnumUnsafe(name string, options []fmt.Stringer, defaultVal fmt.Stringer, usage string) *fmt.Stringer {
	return FlagSetPosEnum[fmt.Stringer](fs, name, options, defaultVal, usage)
}

func FlagSetPosEnum[T fmt.Stringer](fs *FlagSet, name string, options []T, defaultVal T, usage string) *T {
	v := enumInternal[T](name, options, defaultVal, usage)
	fs.Pos(v, name, usage)
	return &v.v
}

func PosEnum[T fmt.Stringer](name string, options []T, defaultVal T, usage string) *T {
	return FlagSetPosEnum[T](CommandLine, name, options, defaultVal, usage)
}

func enumInternal[T fmt.Stringer](name string, options []T, defaultVal T, usage string) *enumValue[T] {
	s := make(map[string]T, len(options))
	for _, opt := range options {
		s[strings.ToLower(opt.String())] = opt
	}
	return &enumValue[T]{defaultVal, s, options}
}

// implements https://pkg.go.dev/github.com/spf13/pflag#Value
type enumValue[T fmt.Stringer] struct {
	v T
	s map[string]T
	o []T
}

func (ev *enumValue[T]) String() string {
	return ev.v.String()
}

func (ev *enumValue[T]) Set(s string) error {
	v, ok := ev.s[strings.ToLower(s)]
	if !ok {
		return fmt.Errorf("got %q, allowed values are %q", s, ev.o)
	}
	ev.v = v
	return nil
}

func (ev *enumValue[T]) Type() string {
	return fmt.Sprintf("Enum: %q", ev.o)
}

// TODO: fill out and move to different file
func (fs *FlagSet) PosString(name string, defaultVal string, usage string) *string {
	ss := &str{usage}
	fs.Pos(ss, name, usage)
	return &ss.v
}

// implements https://pkg.go.dev/github.com/spf13/pflag#Value
type str struct{ v string }

func (ss *str) String() string {
	return ss.v
}

func (ss *str) Set(s string) error {
	ss.v = s
	return nil
}

func (ev *str) Type() string {
	return "string"
}
