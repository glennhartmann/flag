package flag

import (
	"fmt"
	"testing"
)

// color is a simple fmt.Stringer used as an enum type in tests.
type color int

const (
	red color = iota
	green
	blue
)

func (c color) String() string {
	switch c {
	case red:
		return "red"
	case green:
		return "green"
	case blue:
		return "blue"
	}
	return "unknown"
}

var colorOptions = []color{red, green, blue}

// --- enumValue tests ---

func TestEnumValue_SetValid(t *testing.T) {
	ev := enumInternal[color]("c", colorOptions, red, "")
	if err := ev.Set("green"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev.v != green {
		t.Errorf("got %v, want green", ev.v)
	}
}

func TestEnumValue_SetCaseInsensitive(t *testing.T) {
	ev := enumInternal[color]("c", colorOptions, red, "")
	if err := ev.Set("BLUE"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ev.v != blue {
		t.Errorf("got %v, want blue", ev.v)
	}
}

func TestEnumValue_SetInvalid(t *testing.T) {
	ev := enumInternal[color]("c", colorOptions, red, "")
	if err := ev.Set("yellow"); err == nil {
		t.Fatal("expected error for invalid value, got nil")
	}
}

func TestEnumValue_StringReturnsDefault(t *testing.T) {
	ev := enumInternal[color]("c", colorOptions, green, "")
	if got := ev.String(); got != "green" {
		t.Errorf("got %q, want %q", got, "green")
	}
}

func TestEnumValue_StringAfterSet(t *testing.T) {
	ev := enumInternal[color]("c", colorOptions, red, "")
	_ = ev.Set("blue")
	if got := ev.String(); got != "blue" {
		t.Errorf("got %q, want %q", got, "blue")
	}
}

func TestEnumValue_TypeContainsOptions(t *testing.T) {
	ev := enumInternal[color]("c", colorOptions, red, "")
	typ := ev.Type()
	// Type() should mention all option values.
	for _, opt := range colorOptions {
		s := fmt.Sprintf("%q", opt)
		if !contains(typ, opt.String()) {
			t.Errorf("Type() = %q does not mention option %s", typ, s)
		}
	}
}

// --- str tests ---

func TestStr_SetAndString(t *testing.T) {
	s := &str{}
	if err := s.Set("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := s.String(); got != "hello" {
		t.Errorf("got %q, want %q", got, "hello")
	}
}

func TestStr_Type(t *testing.T) {
	s := &str{}
	if got := s.Type(); got != "string" {
		t.Errorf("got %q, want %q", got, "string")
	}
}

// --- FlagSet positional parsing tests ---

func TestParsePosCustom_CorrectArgs(t *testing.T) {
	fs := NewFlagSet()
	a := fs.PosString("first", "", "first arg")
	b := fs.PosString("second", "", "second arg")

	fs.ParsePosCustom([]string{"foo", "bar"})

	if *a != "foo" {
		t.Errorf("first = %q, want %q", *a, "foo")
	}
	if *b != "bar" {
		t.Errorf("second = %q, want %q", *b, "bar")
	}
}

func TestParsePosCustom_TooFewArgsPanics(t *testing.T) {
	fs := NewFlagSet()
	fs.PosString("first", "", "first arg")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for too few args, got none")
		}
	}()
	fs.ParsePosCustom([]string{})
}

func TestParsePosCustom_TooManyArgsPanics(t *testing.T) {
	fs := NewFlagSet()
	fs.PosString("first", "", "first arg")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for too many args, got none")
		}
	}()
	fs.ParsePosCustom([]string{"a", "b"})
}

func TestParsePosCustom_NoArgs(t *testing.T) {
	// A FlagSet with no positional args registered should accept an empty slice.
	fs := NewFlagSet()
	fs.ParsePosCustom([]string{}) // should not panic
}

// --- FlagSetPosEnum tests ---

func TestFlagSetPosEnum_Valid(t *testing.T) {
	fs := NewFlagSet()
	got := FlagSetPosEnum[color](fs, "color", colorOptions, red, "pick a color")

	fs.ParsePosCustom([]string{"blue"})

	if *got != blue {
		t.Errorf("got %v, want blue", *got)
	}
}

func TestFlagSetPosEnum_Invalid(t *testing.T) {
	fs := NewFlagSet()
	FlagSetPosEnum[color](fs, "color", colorOptions, red, "pick a color")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for invalid enum value, got none")
		}
	}()
	fs.ParsePosCustom([]string{"yellow"})
}

func TestFlagSetPosEnum_DefaultNotOverwritten(t *testing.T) {
	fs := NewFlagSet()
	// Register but never parse — value should remain the default.
	got := FlagSetPosEnum[color](fs, "color", colorOptions, green, "pick a color")
	if *got != green {
		t.Errorf("got %v, want green (default)", *got)
	}
}

// --- helpers ---

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		len(s) > 0 && len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
