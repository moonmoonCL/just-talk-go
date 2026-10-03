//go:build darwin

package hotkey

import (
	"testing"
	"time"
)

func newTestDarwinProvider(t *testing.T) *darwinProvider {
	t.Helper()
	return &darwinProvider{
		channels:             make(map[Combo]chan<- Event),
		tracker:              NewKeyStateTracker(),
		activeModifierCombos: make(map[Combo]bool),
	}
}

// keyDown/keyUp mirror handleCGEvent's routing: regular keys feed the
// tracker directly, modifier transitions arrive via flagsChanged.
func (p *darwinProvider) keyDown(t *testing.T, key KeyCode, mods Modifier, now time.Time) []Event {
	t.Helper()
	if key.IsModifier() {
		return p.processFlagsChanged(key, mods, now)
	}
	return p.tracker.KeyDown(key, now)
}

// The RightCtrl-only combo fires on right ctrl press/release and stays
// silent for left ctrl.
func TestDarwinRightCtrlOnlyCombo(t *testing.T) {
	p := newTestDarwinProvider(t)
	if _, err := p.RegisterWithOptions(Combo{Mods: ModRCtrl, Key: KeyNone}, RegisterOptions{}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()

	events := p.keyDown(t, KeyRCtrl, ModCtrl, now) // right ctrl down
	if len(events) != 1 || events[0].Type != KeyDown || events[0].Combo.Mods != ModRCtrl {
		t.Fatalf("right ctrl press: %+v", events)
	}

	events = p.processFlagsChanged(KeyRCtrl, 0, now) // right ctrl up
	if len(events) != 1 || events[0].Type != KeyUp {
		t.Fatalf("right ctrl release: %+v", events)
	}

	events = p.keyDown(t, KeyCtrl, ModCtrl, now) // left ctrl press
	if len(events) != 0 {
		t.Fatalf("left ctrl press must not fire RightCtrl combo: %+v", events)
	}
}

// A flagsChanged event for a modifier that is already believed held is a
// release even when the shared flag stays set (other side still down).
func TestDarwinBothSidesHeldReleaseDisambiguation(t *testing.T) {
	p := newTestDarwinProvider(t)
	if _, err := p.RegisterWithOptions(Combo{Mods: ModRCtrl, Key: KeyNone}, RegisterOptions{}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()

	p.processFlagsChanged(KeyRCtrl, ModCtrl, now)           // right down
	p.processFlagsChanged(KeyCtrl, ModCtrl, now)            // left down (flag unchanged)
	events := p.processFlagsChanged(KeyRCtrl, ModCtrl, now) // right up while left still holds flag
	if len(events) != 1 || events[0].Type != KeyUp {
		t.Fatalf("right ctrl release with left held: %+v", events)
	}
	if p.tracker.ActiveMods() != ModCtrl {
		t.Fatalf("left ctrl still held, activeMods = %v, want Ctrl", p.tracker.ActiveMods())
	}
}

// Side-agnostic combos must keep working: left ctrl still triggers a
// Ctrl-only combo, and standard combos fire from the right side too.
func TestDarwinSideAgnosticCombosUnchanged(t *testing.T) {
	p := newTestDarwinProvider(t)
	if _, err := p.RegisterWithOptions(Combo{Mods: ModCtrl, Key: KeyNone}, RegisterOptions{}); err != nil {
		t.Fatal(err)
	}
	now := time.Now()

	events := p.processFlagsChanged(KeyCtrl, ModCtrl, now)
	if len(events) != 1 || events[0].Type != KeyDown {
		t.Fatalf("left ctrl press on Ctrl-only combo: %+v", events)
	}

	events = p.processFlagsChanged(KeyCtrl, 0, now)
	if len(events) != 1 || events[0].Type != KeyUp {
		t.Fatalf("left ctrl release on Ctrl-only combo: %+v", events)
	}

	// Standard combo via right-side modifier.
	want := Combo{Mods: ModCtrl, Key: KeyF}
	if _, err := p.RegisterWithOptions(want, RegisterOptions{}); err != nil {
		t.Fatal(err)
	}
	p.processFlagsChanged(KeyRCtrl, ModCtrl, now) // right ctrl down
	events = p.keyDown(t, KeyF, ModCtrl, now)
	if len(events) == 0 || events[0].Combo != want || events[0].Type != KeyDown {
		t.Fatalf("Ctrl+F from right ctrl: %+v", events)
	}
}
