package test

import (
	"testing"
)

func TestClient(t *testing.T) {
	emulator := DofusServerEmulator{}
	emulator.Start(t)
}
