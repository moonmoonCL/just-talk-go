package config

import (
	"testing"

	"github.com/c/just-talk-go/hotkey"
)

func TestParseHotkeySideModifiers(t *testing.T) {
	cases := []struct {
		in   string
		want hotkey.Combo
	}{
		{"RightCtrl", hotkey.Combo{Mods: hotkey.ModRCtrl, Key: hotkey.KeyNone}},
		{"Right Ctrl", hotkey.Combo{Mods: hotkey.ModRCtrl, Key: hotkey.KeyNone}},
		{"rctrl", hotkey.Combo{Mods: hotkey.ModRCtrl, Key: hotkey.KeyNone}},
		{"RightCmd", hotkey.Combo{Mods: hotkey.ModRSuper, Key: hotkey.KeyNone}},
		{"RightOption", hotkey.Combo{Mods: hotkey.ModRAlt, Key: hotkey.KeyNone}},
		{"right option", hotkey.Combo{Mods: hotkey.ModRAlt, Key: hotkey.KeyNone}},
		{"RightShift+F9", hotkey.Combo{Mods: hotkey.ModRShift, Key: hotkey.KeyF9}},
		{"Ctrl+Option", hotkey.Combo{Mods: hotkey.ModCtrl | hotkey.ModAlt, Key: hotkey.KeyNone}},
	}
	for _, c := range cases {
		got, err := ParseHotkey(c.in)
		if err != nil {
			t.Errorf("ParseHotkey(%q): %v", c.in, err)
			continue
		}
		if got != c.want {
			t.Errorf("ParseHotkey(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestSideModifierStringRoundTrip(t *testing.T) {
	combo := hotkey.Combo{Mods: hotkey.ModRCtrl, Key: hotkey.KeyNone}
	if got := combo.String(); got != "RightCtrl" {
		t.Fatalf("combo.String() = %q, want %q", got, "RightCtrl")
	}
	again, err := ParseHotkey(combo.String())
	if err != nil || again != combo {
		t.Fatalf("round trip of %q: %v, %v", combo.String(), again, err)
	}
}
