package meshtastic

import (
	pbuf "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"encoding/binary"
	"fmt"
	"math"
)

type NodeID uint32

func (n NodeID) Uint32() uint32 {
	return uint32(n)
}

func (n NodeID) String() string {
	return fmt.Sprintf("!%08x", uint32(n))
}

// Bytes converts the NodeID to a byte slice
func (n NodeID) Bytes() []byte {
	bytes := make([]byte, 4) // uint32 is 4 bytes
	binary.BigEndian.PutUint32(bytes, n.Uint32())
	return bytes
}

// BroadcastNodeID is the special NodeID used when broadcasting a packet to a channel.
const BroadcastNodeID NodeID = math.MaxUint32

type Node struct {
	LongName      string
	ShortName     string
	ID            uint32
	HardwareModel pbuf.HardwareModel
}

// Not actually in use yet 😅
func (n *Node) EncryptPacket(pkt *pbuf.MeshPacket, channelName string, key []byte) *pbuf.MeshPacket {
	payload := pkt.GetPayloadVariant()
	_ = payload
	switch p := payload.(type) {
	case *pbuf.MeshPacket_Decoded:
		_ = p
		encrypted := pbuf.MeshPacket_Encrypted{
			Encrypted: nil,
		}
		_ = encrypted
	}
	return nil
}
