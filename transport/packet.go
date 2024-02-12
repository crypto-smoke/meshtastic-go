package transport

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crypto-smoke/meshtastic-go/radio"
	"google.golang.org/protobuf/proto"
	"log/slog"
)

// xorHash computes a simple XOR hash of the provided byte slice.
func xorHash(p []byte) uint8 {
	var code uint8
	for _, b := range p {
		code ^= b
	}
	return code
}

// GenerateHash returns the hash for a given channel by XORing the channel name and PSK.
func ChannelHash(channelName string, channelKey []byte) (uint32, error) {
	if len(channelKey) == 0 {
		return 0, fmt.Errorf("channel key cannot be empty")
	}

	h := xorHash([]byte(channelName))
	h ^= xorHash(channelKey)

	return uint32(h), nil
}
func tojson(data *pb.Data) ([]byte, error) {
	str, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		slog.Error("failed marshalling protobuf to json", "err", err)
		return nil, err
	}
	return str, nil
}
func (c *Client) GetData(pkt *pb.MeshPacket) (*pb.Data, error) {
	switch payload := pkt.GetPayloadVariant().(type) {
	case *pb.MeshPacket_Decoded:
		// json marshal
		return payload.Decoded, nil
	case *pb.MeshPacket_Encrypted:
		channels := c.State.Channels()
		for _, ch := range channels {
			if ch.Settings.Name == "" || ch.Settings.Psk == nil {
				continue
			}
			hash, err := ChannelHash(ch.Settings.Name, ch.Settings.Psk)
			if err != nil {
				return nil, err
			}
			if hash == pkt.Channel {
				data, err := c.TryDecode(pkt, ch.Settings.Psk)
				if err != nil {
					return nil, err
				}
				return data, nil
			}
		}
	}
	return nil, errors.New("unknown packet type")
}
func (c *Client) ToJson(pkt *pb.MeshPacket) ([]byte, error) {
	switch payload := pkt.GetPayloadVariant().(type) {
	case *pb.MeshPacket_Decoded:
		// json marshal
		return tojson(payload.Decoded)
	case *pb.MeshPacket_Encrypted:
		channels := c.State.Channels()
		for _, ch := range channels {
			if ch.Settings.Name == "" || ch.Settings.Psk == nil {
				continue
			}
			hash, err := ChannelHash(ch.Settings.Name, ch.Settings.Psk)
			if err != nil {
				return nil, err
			}
			if hash == pkt.Channel {
				data, err := c.TryDecode(pkt, ch.Settings.Psk)
				if err != nil {
					return nil, err
				}
				return tojson(data)
			}
		}
	}

	/*
		decoded, err := c.TryDecode(pkt, radio.DefaultKey)
		if err != nil {
			//log.Error("failed decoding packet", "err", err, "topic", m.Topic, "payload", hex.EncodeToString(m.Payload))
			if !errors.Is(err, radio.ErrDecrypt) {
				//	log.Error("failed decoding packet", "err", err, "topic", m.Topic, "payload", hex.EncodeToString(m.Payload))
				data, jsonErr := json.Marshal(packet)
				_ = data
				if jsonErr != nil {
					//		log.Error("failed marshalling packet for debug output", "jsonErr", jsonErr)
				}
				//		log.Info("packet", "json", string(data))
			}
		}

	*/
	return nil, nil
}

var ErrDecrypt = errors.New("unable to decrypt payload")

func (c *Client) TryDecode(packet *pb.MeshPacket, key []byte) (*pb.Data, error) {
	switch packet.GetPayloadVariant().(type) {
	case *pb.MeshPacket_Decoded:
		return packet.GetDecoded(), nil
	case *pb.MeshPacket_Encrypted:

		decrypted, err := radio.XOR(packet.GetEncrypted(), key, packet.Id, packet.From)
		if err != nil {
			//log.Error("failed decrypting packet", "error", err)
			return nil, ErrDecrypt
		}

		var meshPacket pb.Data
		err = proto.Unmarshal(decrypted, &meshPacket)
		if err != nil {
			return nil, ErrDecrypt
		}
		return &meshPacket, nil
	default:
		return nil, errors.New("unknown payload variant")
	}
}
