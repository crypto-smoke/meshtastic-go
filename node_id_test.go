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

// TestRandomNodeID ensures that RandomNodeID generates a valid NodeID and that multiple calls generate different
// NodeIDs.
func TestRandomNodeID(t *testing.T) {
	nodeID1, err := RandomNodeID()
	if err != nil {
		t.Errorf("expected no error when generating the first node id, got %v", err)
	}
	t.Logf("nodeID1: %s", nodeID1)
	nodeID2, err := RandomNodeID()
	if err != nil {
		t.Errorf("expected no error when generating the second node id, got %v", err)
	}
	t.Logf("nodeID2: %s", nodeID2)
	if nodeID1 == nodeID2 {
		t.Errorf("expected random node ids to be different, got %s and %s", nodeID1, nodeID2)
	}
}
