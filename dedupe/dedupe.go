package dedup

import (
	"hash/maphash"
	"sync"
	"time"
)

// PacketDeduplicator is a structure that prevents processing of duplicate packets.
// It keeps a record of seen packets and the time they were last seen.
type PacketDeduplicator struct {
	expiresAfter time.Duration // expiresAfter defines the duration after which a seen packet record expires.
	seed         maphash.Seed
	sync.Mutex                 // protects the seen map from concurrent access.
	seen         map[uint64]time.Time // seen maps a packet identifier to the last time it was seen.
}

// NewDeduplicator creates a new PacketDeduplicator with a given hasher and expiration duration for packet records.
// It starts a background goroutine to periodically clean up expired packet records.
func NewDeduplicator(expiresAfter time.Duration) *PacketDeduplicator {
	pd := &PacketDeduplicator{
		seen:         map[uint64]time.Time{},
		seed:         maphash.MakeSeed(),
		expiresAfter: expiresAfter,
	}
	go func() {
		for range time.NewTicker(time.Second * 10).C {
			pd.Lock()
			for packet, timestamp := range pd.seen {
				if time.Since(timestamp) > pd.expiresAfter {
					delete(pd.seen, packet)
				}
			}
			pd.Unlock()
		}
	}()

	return pd
}

func (p *PacketDeduplicator) isSeen(k uint64) bool {
	p.Lock()
	defer p.Unlock()
	d, exists := p.seen[k]
	if !exists {
		p.seen[k] = time.Now()
		return false
	}
	if time.Since(d) > p.expiresAfter {
		delete(p.seen, k)
		return false
	}
	return true
}

// Seen checks whether a packet with the given sender and packetID has been seen before.
// If not, it records the packet as seen and returns false. Otherwise, it returns true.
func (p *PacketDeduplicator) Seen(sender, packetID uint32) bool {
	return p.isSeen((uint64(sender) << 32) | uint64(packetID))
}

// SeenData checks whether the data has been seen before based on its hashed value.
// If not, it records the data as seen and returns false. Otherwise, it returns true.
func (p *PacketDeduplicator) SeenData(data []byte) bool {
	return p.isSeen(maphash.Bytes(p.seed, data))
}
