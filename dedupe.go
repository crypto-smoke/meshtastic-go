package meshtastic

import (
	"encoding/hex"
	"fmt"
	"hash"
	"sync"
	"time"
)

// PacketDeduplicator is a structure that prevents processing of duplicate packets.
// It keeps a record of seen packets and the time they were last seen.
type PacketDeduplicator struct {
	hasher       hash.Hash            // hasher is used for generating packet identifiers.
	expiresAfter time.Duration        // expiresAfter defines the duration after which a seen packet record expires.
	sync.RWMutex                      // RWMutex is used to protect the seen map from concurrent access.
	seen         map[string]time.Time // seen maps a packet identifier to the last time it was seen.
}

// NewDeduplicator creates a new PacketDeduplicator with a given hasher and expiration duration for packet records.
// It starts a background goroutine to periodically clean up expired packet records.
func NewDeduplicator(hasher hash.Hash, expiresAfter time.Duration) *PacketDeduplicator {
	pd := PacketDeduplicator{
		seen:         make(map[string]time.Time),
		hasher:       hasher,
		expiresAfter: expiresAfter,
	}
	go func() {
		for {
			time.Sleep(expiresAfter)
			pd.Lock()
			for packet, timestamp := range pd.seen {
				if time.Since(timestamp) > pd.expiresAfter {
					delete(pd.seen, packet)
				}
			}
			pd.Unlock()
		}
	}()

	return &pd
}

// Seen checks whether a packet with the given sender and packetID has been seen before.
// If not, it records the packet as seen and returns false. Otherwise, it returns true.
func (p *PacketDeduplicator) Seen(sender, packetID uint32) bool {
	asString := fmt.Sprintf("%d-%d", sender, packetID)
	p.RLock()
	if _, exists := p.seen[asString]; !exists {
		p.RUnlock()
		p.Lock()
		defer p.Unlock()
		p.seen[asString] = time.Now()
		return false
	}
	p.RUnlock()
	return true
}

// SeenData checks whether the data has been seen before based on its hashed value.
// If not, it records the data as seen and returns false. Otherwise, it returns true.
func (p *PacketDeduplicator) SeenData(data []byte) bool {
	hashed := p.hasher.Sum(data)
	asHex := hex.EncodeToString(hashed)
	p.RLock()
	if _, exists := p.seen[asHex]; !exists {
		p.RUnlock()
		p.Lock()
		defer p.Unlock()
		p.seen[asHex] = time.Now()
		return false
	}
	p.RUnlock()

	return true
}
