package hotkey

import (
	"testing"
	"time"
)

// Right-side modifier keys must feed the side-agnostic modifier mask so
// standard combos (e.g. Ctrl+F) trigger from either side.
func TestSideModifierDrivesStandardCombo(t *testing.T) {
	tr := NewKeyStateTracker()
	ch := make(chan Event, 8)
	combo := Combo{Mods: ModCtrl, Key: KeyF}
	tr.Watch(combo, ch)

	now := time.Now()
	tr.KeyDown(KeyRCtrl, now)
	if got := tr.ActiveMods(); got != ModCtrl {
		t.Fatalf("after RightCtrl down, activeMods = %v, want Ctrl", got)
	}

	events := tr.KeyDown(KeyF, now)
	if len(events) != 1 || events[0].Combo != combo || events[0].Type != KeyDown {
		t.Fatalf("Ctrl+F did not fire from right-side ctrl: %+v", events)
	}
}

func TestSideModifierReleaseKeepsOtherSideHeld(t *testing.T) {
	tr := NewKeyStateTracker()
	now := time.Now()

	tr.KeyDown(KeyRCtrl, now)
	tr.KeyDown(KeyCtrl, now)
	tr.KeyUp(KeyRCtrl, now)
	if got := tr.ActiveMods(); got != ModCtrl {
		t.Fatalf("after right released while left held, activeMods = %v, want Ctrl", got)
	}

	tr.KeyUp(KeyCtrl, now)
	if got := tr.ActiveMods(); got != ModNone {
		t.Fatalf("after both released, activeMods = %v, want None", got)
	}
}

func TestLeftModifierReleaseKeepsRightSideHeld(t *testing.T) {
	tr := NewKeyStateTracker()
	now := time.Now()

	tr.KeyDown(KeyCtrl, now)
	tr.KeyDown(KeyRCtrl, now)
	tr.KeyUp(KeyCtrl, now)
	if got := tr.ActiveMods(); got != ModCtrl {
		t.Fatalf("after left released while right held, activeMods = %v, want Ctrl", got)
	}
}
