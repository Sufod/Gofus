package test

import (
	"testing"
)

func TestClient(t *testing.T) {
	emulator := DofusServerEmulator{}
	go emulator.startClient()
	emulator.Start(t)
}

func TestWithDebug(t *testing.T) {
	emulator := DofusServerEmulator{}
	go emulator.Start(t)
	emulator.startClient()
}
