package flag

import (
	"github.com/spf13/pflag"
)

// FlagSet holds registered positional arguments alongside the standard pflag
// flag set. Use NewFlagSet to create one, or use the package-level CommandLine
// instance for a single-command program.
type FlagSet struct {
	posArgs []posArg
}

type posArg struct {
	val         pflag.Value
	name, usage string
}

// CommandLine is the default FlagSet, used by the package-level functions.
var CommandLine = NewFlagSet()

// NewFlagSet creates an empty FlagSet with no positional arguments registered.
func NewFlagSet() *FlagSet {
	return &FlagSet{}
}

// Pos registers a positional argument on fs. Arguments are matched by
// position in the order they are registered.
func (fs *FlagSet) Pos(val pflag.Value, name, usage string) {
	fs.posArgs = append(fs.posArgs, posArg{val, name, usage})
}

// Pos registers a positional argument on CommandLine.
func Pos(val pflag.Value, name, usage string) {
	CommandLine.Pos(val, name, usage)
}

// Parse parses pflag flags from os.Args, then parses positional arguments
// from the remaining non-flag arguments.
func (fs *FlagSet) Parse() {
	pflag.Parse()
	fs.ParsePos()
}

// Parse parses flags and positional arguments on CommandLine.
func Parse() {
	CommandLine.Parse()
}

// ParsePos parses positional arguments from the non-flag arguments left over
// after pflag.Parse has already been called.
func (fs *FlagSet) ParsePos() {
	fs.ParsePosCustom(pflag.Args())
}

// ParsePos parses positional arguments on CommandLine.
func ParsePos() {
	CommandLine.ParsePos()
}

// ParsePosCustom parses positional arguments from the provided args slice.
// The slice must have exactly the same length as the number of registered
// positional arguments, or parsing will fail.
func (fs *FlagSet) ParsePosCustom(args []string) {
	if len(args) != len(fs.posArgs) {
		// TODO: have a lenient mode, allow for last arg to have indefinite size
		parsePosFail()
	}

	for i := range args {
		if err := fs.posArgs[i].val.Set(args[i]); err != nil {
			parsePosFail()
		}
	}
}

// ParsePosCustom parses positional arguments on CommandLine from the provided
// args slice.
func ParsePosCustom(args []string) {
	CommandLine.ParsePosCustom(args)
}

func parsePosFail() {
	// TODO: print usage and stuff - maybe try to make the cause of the failure obvious?
	panic("parsePosFail()")
	//os.Exit(1)
}
