# flag

Extensions to [github.com/spf13/pflag](https://pkg.go.dev/github.com/spf13/pflag):

- **Typed enum flags** — restrict a flag to a fixed set of values, with case-insensitive matching and helpful error messages.
- **Positional argument parsing** — declare expected positional arguments by name and parse them from the non-flag remainder of `os.Args`.

_Disclaimer: This is a personal project. The views, code, and opinions expressed here are my own and do not represent those of my current or past employers._

## Installation

```
go get bitbucket.org/glennhartmann/flag
```

## Usage

### Enum flags

```go
import flag "bitbucket.org/glennhartmann/flag/src/flag"

type Color int

const (
    Red Color = iota
    Green
    Blue
)

func (c Color) String() string {
    return [...]string{"red", "green", "blue"}[c]
}

var colors = []Color{Red, Green, Blue}

color := flag.Enum("color", colors, Red, "output color")
// or with a shorthand:
color := flag.EnumP("color", "c", colors, Red, "output color")

flag.Parse()
fmt.Println(*color) // e.g. "green"
```

Accepted values are matched case-insensitively against each option's `String()` representation. An invalid value produces an error listing the allowed options.

### Positional arguments

```go
name   := flag.PosString("name",   "", "the name to greet")
target := flag.PosString("target", "", "where to send the greeting")
// or a positional enum:
color  := flag.PosEnum("color", colors, Red, "output color")

flag.Parse() // parses flags, then positional args
fmt.Printf("Hello %s -> %s (in %s)\n", *name, *target, *color)
```

Positional arguments are matched in registration order. The number of non-flag arguments must match the number of registered positional arguments exactly.

### Using a FlagSet directly

All package-level functions have `FlagSet` method equivalents for use in multi-command programs:

```go
fs := flag.NewFlagSet()
color := flag.FlagSetEnum(fs, "color", colors, Red, "output color")
name  := fs.PosString("name", "", "the name")

fs.Parse()
```

## API overview

| Function | Description |
|---|---|
| `Enum[T](name, options, default, usage)` | Typed enum flag on `CommandLine` |
| `EnumP[T](name, short, options, default, usage)` | Same, with a shorthand character |
| `PosEnum[T](name, options, default, usage)` | Typed positional enum argument |
| `PosString(name, default, usage)` | Positional string argument |
| `Pos(val, name, usage)` | Positional argument from any `pflag.Value` |
| `Parse()` | Parse flags and positional arguments |
| `ParsePos()` | Parse only positional arguments (after pflag.Parse) |
| `ParsePosCustom(args)` | Parse positional arguments from a custom slice |

`FlagSet` method variants (`FlagSetEnum`, `FlagSetEnumP`, `FlagSetPosEnum`) and `Unsafe` variants (erased to `fmt.Stringer`) are also available.
