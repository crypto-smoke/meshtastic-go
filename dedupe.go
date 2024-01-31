package meshtastic

import (
	"encoding/hex"
	"fmt"
	"hash"
	"sync"
	"time"
)

type PacketDeduplicator struct {
	hasher       hash.Hash
	expiresAfter time.Duration
	sync.RWMutex
	seen map[string]time.Time
}

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
