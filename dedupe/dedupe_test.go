package dedupe_test

import (
	"crypto/md5"
	"github.com/crypto-smoke/meshtastic-go/dedupe"
	"testing"
	"time"
)

func TestPacketDeduplicatorSeen(t *testing.T) {
	hasher := md5.New()
	expiresAfter := 100 * time.Millisecond
	dedup := dedupe.NewDeduplicator(hasher, expiresAfter)

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
	time.Sleep(expiresAfter + 100*time.Millisecond)

	// Test Seen with same packetID after expiration
	if dedup.Seen(sender, packetID) {
		t.Error("Expected the packet to not have been seen after expiration")
	}
}

func TestPacketDeduplicatorSeenData(t *testing.T) {
	hasher := md5.New()
	expiresAfter := 100 * time.Millisecond
	dedup := dedupe.NewDeduplicator(hasher, expiresAfter)

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
	time.Sleep(expiresAfter + 100*time.Millisecond)

	// Test SeenData with same data after expiration
	if dedup.SeenData(data) {
		t.Error("Expected the data to not have been seen after expiration")
	}
}
