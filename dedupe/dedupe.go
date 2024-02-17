package dedup

import (
	"hash/maphash"
	"math/rand"
	"sync"
	"time"
)

// PacketDeduplicator is a structure that prevents processing of duplicate packets.
// It keeps a record of seen packets and the time they were last seen.
type PacketDeduplicator struct {
	expiresAfter time.Duration // expiresAfter defines the duration after which a seen packet record expires.
	seed         maphash.Seed
	sync.Mutex                        // protects the seen map from concurrent access.
	seen         map[uint64]time.Time // seen maps a packet identifier to the last time it was seen.
}

// NewDeduplicator creates a new PacketDeduplicator with a given hasher and expiration duration for packet records.
// It starts a background goroutine to periodically clean up expired packet records.
func NewDeduplicator(expiresAfter time.Duration) *PacketDeduplicator {
	return &PacketDeduplicator{
		seen:         map[uint64]time.Time{},
		seed:         maphash.MakeSeed(),
		expiresAfter: expiresAfter,
	}
}

func (p *PacketDeduplicator) isSeen(now time.Time, k uint64) bool {
	// Checking for dedup purges ~5% of the time
	if rand.Intn(100) < 5 {
		p.deleteUntil(now.Add(-p.expiresAfter))
	}

	p.Lock()
	defer p.Unlock()
	d, exists := p.seen[k]
	if !exists {
		p.seen[k] = now
		return false
	}
	if now.Sub(d) > p.expiresAfter {
		delete(p.seen, k)
		return false
	}
	return true
}

// Seen checks whether a packet with the given sender and packetID has been seen before.
// If not, it records the packet as seen and returns false. Otherwise, it returns true.
func (p *PacketDeduplicator) Seen(sender, packetID uint32) bool {
	return p.seenAt(time.Now(), sender, packetID)
}

// SeenData checks whether the data has been seen before based on its hashed value.
// If not, it records the data as seen and returns false. Otherwise, it returns true.
func (p *PacketDeduplicator) SeenData(data []byte) bool {
	return p.seenDataAt(time.Now(), data)
}

//
// These are used internally and are test hooks allowing us to avoid the clock.
//

func (p *PacketDeduplicator) deleteUntil(t time.Time) {
	p.Lock()
	defer p.Unlock()
	for packet, timestamp := range p.seen {
		if timestamp.Before(t) {
			delete(p.seen, packet)
		}
	}
}

// Seendata with an explicit time
func (p *PacketDeduplicator) seenDataAt(now time.Time, data []byte) bool {
	return p.isSeen(now, maphash.Bytes(p.seed, data))
}

// Seen with an explicit time
func (p *PacketDeduplicator) seenAt(now time.Time, sender uint32, packetID uint32) bool {
	return p.isSeen(now, (uint64(sender)<<32)|uint64(packetID))
}
