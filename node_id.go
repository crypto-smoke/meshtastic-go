package meshtastic

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
)

// NodeID holds the node identifier. This is a uint32 value which uniquely identifies a node within a mesh.
type NodeID uint32

const (
	// BroadcastNodeID is the special NodeID used when broadcasting a packet to a channel.
	BroadcastNodeID NodeID = math.MaxUint32
	// ReservedNodeIDThreshold is the threshold at which NodeIDs are considered reserved. Random NodeIDs should not
	// be generated below this threshold.
	// Source: https://github.com/meshtastic/firmware/blob/d1ea58975755e146457a8345065e4ca357555275/src/mesh/NodeDB.cpp#L461
	reservedNodeIDThreshold NodeID = 4
)

// Uint32 returns the underlying uint32 value of the NodeID.
func (n NodeID) Uint32() uint32 {
	return uint32(n)
}

// String converts the NodeID to a hex formatted string.
// This is typically how NodeIDs are displayed in Meshtastic UIs.
func (n NodeID) String() string {
	return fmt.Sprintf("!%08x", uint32(n))
}

// Bytes converts the NodeID to a byte slice
func (n NodeID) Bytes() []byte {
	bytes := make([]byte, 4) // uint32 is 4 bytes
	binary.BigEndian.PutUint32(bytes, n.Uint32())
	return bytes
}

// DefaultLongName returns the default long node name based on the NodeID.
// Source: https://github.com/meshtastic/firmware/blob/d1ea58975755e146457a8345065e4ca357555275/src/mesh/NodeDB.cpp#L382
func (n NodeID) DefaultLongName() string {
	return "Meshtastic " + n.DefaultShortName()
}

// DefaultShortName returns the default short node name based on the NodeID.
// Last two bytes of the NodeID represented in hex.
// Source: https://github.com/meshtastic/firmware/blob/d1ea58975755e146457a8345065e4ca357555275/src/mesh/NodeDB.cpp#L382
func (n NodeID) DefaultShortName() string {
	bytes := make([]byte, 4) // uint32 is 4 bytes
	binary.BigEndian.PutUint32(bytes, n.Uint32())
	return fmt.Sprintf("%04x", bytes[2:])
}

// RandomNodeID returns a randomised NodeID.
// It's recommended to call this the first time a node is started and persist the result.
//
// Hardware meshtastic nodes first try a NodeID of the last four bytes of the BLE MAC address. If that ID is already in
// use or invalid, a random NodeID is generated.
// Source: https://github.com/meshtastic/firmware/blob/d1ea58975755e146457a8345065e4ca357555275/src/mesh/NodeDB.cpp#L466
func RandomNodeID() (NodeID, error) {
	randBytes := make([]byte, 4)
	// TODO: Ensure generated value is not below ReservedNodeIDThreshold and not BroadcastNodeID
	_, err := rand.Read(randBytes)
	if err != nil {
		return NodeID(0), fmt.Errorf("reading entropy: %w", err)
	}
	return NodeID(binary.BigEndian.Uint32(randBytes)), nil
}
