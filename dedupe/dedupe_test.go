package dedup

import (
	"testing"
	"testing/quick"
	"time"
)

const expiresAfter = time.Millisecond

func TestPacketDeduplicatorSeen(t *testing.T) {
	dedup := NewDeduplicator(expiresAfter)

	sender := uint32(1)
	packetID := uint32(1)

	// Test Seen with new packetID
	if dedup.Seen(sender, packetID) {
		t.Error("Expected the packet to not have been seen")
	}

	// Test Seen with same packetID again
	if !dedup.Seen(sender, packetID) {
		t.Error("Expected the packet to have been seen")
	}

	// Wait for expiration
	time.Sleep(expiresAfter + time.Millisecond)

	// Test Seen with same packetID after expiration
	if dedup.Seen(sender, packetID) {
		t.Error("Expected the packet to not have been seen after expiration")
	}
}

func TestDuplicatorProp(t *testing.T) {
	if err := quick.Check(func(s, p uint32) bool {
		dedup := NewDeduplicator(expiresAfter)
		// Test Seen with new packetID
		if dedup.Seen(s, p) {
			t.Error("Expected the packet to not have been seen")
			return false
		}

		// Test Seen with same packetID again
		if !dedup.Seen(s, p) {
			t.Error("Expected the packet to have been seen")
			return false
		}
		return true
	}, nil); err != nil {
		t.Error(err)
	}
}

func FuzzDup(f *testing.F) {
	f.Fuzz(func(t *testing.T, s uint32, p uint32) {
		dedup := NewDeduplicator(10*time.Second)
		// Test Seen with new packetID
		if dedup.Seen(s, p) {
			t.Errorf("Expected the packet %v/%v to not have been seen the first time", s, p)
		}

		// Test Seen with same packetID again
		if !dedup.Seen(s, p) {
			t.Errorf("Expected the packet %v/%v to have been seen", s, p)
		}

	})
}

func TestPacketDeduplicatorSeenData(t *testing.T) {
	dedup := NewDeduplicator(expiresAfter)

	data := []byte("test data")

	// Test SeenData with new data
	if dedup.SeenData(data) {
		t.Error("Expected the data to not have been seen")
	}

	// Test SeenData with same data again
	if !dedup.SeenData(data) {
		t.Error("Expected the data to have been seen")
	}

	// Wait for expiration
	time.Sleep(expiresAfter + time.Millisecond)

	// Test SeenData with same data after expiration
	if dedup.SeenData(data) {
		t.Error("Expected the data to not have been seen after expiration")
	}
}
