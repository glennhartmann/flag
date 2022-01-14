package flag

import (
	"github.com/spf13/pflag"
)

type FlagSet struct {
	posArgs []posArg
}

type posArg struct {
	val         pflag.Value
	name, usage string
}

var CommandLine = NewFlagSet()

func NewFlagSet() *FlagSet {
	return &FlagSet{}
}

func (fs *FlagSet) Pos(val pflag.Value, name, usage string) {
	fs.posArgs = append(fs.posArgs, posArg{val, name, usage})
}

func Pos(val pflag.Value, name, usage string) {
	CommandLine.Pos(val, name, usage)
}

func (fs *FlagSet) Parse() {
	pflag.Parse()
	fs.ParsePos()
}

func Parse() {
	CommandLine.Parse()
}

func (fs *FlagSet) ParsePos() {
	fs.ParsePosCustom(pflag.Args())
}

func ParsePos() {
	CommandLine.ParsePos()
}

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

func ParsePosCustom(args []string) {
	CommandLine.ParsePosCustom(args)
}

func parsePosFail() {
	// TODO: print usage and stuff - maybe try to make the cause of the failure obvious?
	panic("parsePosFail()")
	//os.Exit(1)
}
