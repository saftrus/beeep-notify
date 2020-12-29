package beeep

import (
	"testing"
)

func TestNotify(t *testing.T) {
	var action [][]string
	err := Notify("Notify title", "Message body", "assets/information.png", action)
	if err != nil {
		t.Error(err)
	}
}
