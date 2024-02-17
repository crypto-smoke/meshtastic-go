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
	// Some arbitrary time.  Ideally the generator would provide one, but the go quickcheck generators aren't very good
	now := time.Now()
	if err := quick.Check(func(s, p uint32) bool {
		dedup := NewDeduplicator(expiresAfter)
		// Test Seen with new packetID
		if dedup.seenAt(now, s, p) {
			t.Error("Expected the packet to not have been seen")
			return false
		}

		// Test Seen with same packetID again
		if !dedup.seenAt(now, s, p) {
			t.Error("Expected the packet to have been seen")
			return false
		}

		// Even up to expiry
		if !dedup.seenAt(now.Add(expiresAfter), s, p) {
			t.Error("Expected the packet to have been seen all the way to expiry")
			return false
		}

		// But not in The Future
		if dedup.seenAt(now.Add(expiresAfter+time.Nanosecond), s, p) {
			t.Error("Expected the packet to not have been seen after expiry")
			return false
		}

		return true
	}, nil); err != nil {
		t.Error(err)
	}
}

func FuzzDup(f *testing.F) {
	now := time.Now()
	f.Fuzz(func(t *testing.T, s1 uint32, s2 uint32, p1 uint32, p2 uint32) {
		dedup := NewDeduplicator(expiresAfter)
		// Test Seen with new packetID
		if dedup.seenAt(now, s1, p1) {
			t.Error("Expected the packet to not have been seen")
		}

		// Test Seen with same packetID again
		if !dedup.seenAt(now, s1, p1) {
			t.Error("Expected the packet to have been seen")
		}

		// Even up to expiry
		if !dedup.seenAt(now.Add(expiresAfter), s1, p1) {
			t.Error("Expected the packet to have been seen all the way to expiry")
		}

		// But not in The Future
		if dedup.seenAt(now.Add(expiresAfter+time.Nanosecond), s1, p1) {
			t.Error("Expected the packet to not have been seen after expiry")
		}

		//
		// This remaining bits of this property are not valid if s1 == s2 and p1 == p2
		//
		if s1 == s2 && p1 == p2 {
			return
		}

		if dedup.seenAt(now, s2, p2) {
			t.Error("Expected a different packet to not have been seen")
		}

		// Test Seen with same packetID again
		if !dedup.seenAt(now, s2, p2) {
			t.Error("Expected a different packet to have been seen")
		}

		// Even up to expiry
		if !dedup.seenAt(now.Add(expiresAfter), s2, p2) {
			t.Error("Expected a different packet to have been seen all the way to expiry")
		}

		// But not in The Future
		if dedup.seenAt(now.Add(expiresAfter+time.Nanosecond), s2, p2) {
			t.Error("Expected a different packet to not have been seen after expiry")
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
