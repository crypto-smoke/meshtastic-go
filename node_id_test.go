package meshtastic

import (
	"bytes"
	"testing"
)

const testNodeID = 3735928559

func TestNodeID_Uint32(t *testing.T) {
	nodeID := NodeID(testNodeID)
	got := nodeID.Uint32()
	if got != testNodeID {
		t.Errorf("expected %v, got %v", testNodeID, got)
	}
}

func TestNodeID_Bytes(t *testing.T) {
	nodeID := NodeID(testNodeID)
	want := []byte{0xde, 0xad, 0xbe, 0xef}
	got := nodeID.Bytes()
	if !bytes.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestNodeID_String(t *testing.T) {
	nodeID := NodeID(testNodeID)
	want := "!deadbeef"
	got := nodeID.String()
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestNodeID_DefaultShortName(t *testing.T) {
	nodeID := NodeID(testNodeID)
	want := "beef"
	got := nodeID.DefaultShortName()
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}

func TestNodeID_DefaultLongName(t *testing.T) {
	nodeID := NodeID(testNodeID)
	want := "Meshtastic beef"
	got := nodeID.DefaultLongName()
	if want != got {
		t.Errorf("expected %v, got %v", want, got)
	}
}
