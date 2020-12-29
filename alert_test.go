package beeep

import (
	"testing"
)

func TestAlert(t *testing.T) {
	var action [][]string
	err := Alert("Alert title", "Message body", "assets/warning.png", action)
	if err != nil {
		t.Error(err)
	}
}
