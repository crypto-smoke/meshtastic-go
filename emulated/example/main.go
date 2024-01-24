package main

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"encoding/base64"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/crypto-smoke/meshtastic-go/emulated"
	"github.com/crypto-smoke/meshtastic-go/mqtt"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	ctx := context.TODO()
	key, err := base64.URLEncoding.DecodeString("yGwxct335Pb8LdAt8Xk47_00bkbzEFQlwP45Id3HTjQ=")
	if err != nil {
		panic(err)
	}
	myNodeID := meshtastic.NodeID(3735928559)
	radio, err := emulated.NewRadio(emulated.Config{
		LongName:   "M7NOA - GLNG",
		ShortName:  "DOGS",
		NodeID:     meshtastic.NodeID(3735928559),
		MQTTClient: &mqtt.DefaultClient,
		Channels: &pb.ChannelSet{
			Settings: []*pb.ChannelSettings{
				{
					Name: "UKWide",
					Psk:  key,
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}
	go func() {
		// Forgive me, for I have sinned.
		time.Sleep(10 * time.Second)
		err := radio.ToRadio(ctx, &pb.ToRadio{
			PayloadVariant: &pb.ToRadio_Packet{
				Packet: &pb.MeshPacket{
					From: myNodeID.Uint32(),
					// This is hard coded to my node ID
					To: 2437877602,
					PayloadVariant: &pb.MeshPacket_Decoded{
						Decoded: &pb.Data{
							Portnum: pb.PortNum_TEXT_MESSAGE_APP,
							Payload: []byte("from main!!"),
						},
					},
				},
			},
		})
		if err != nil {
			log.Error("sending to radio", "error", err)
		}
	}()
	if err := radio.Run(ctx); err != nil {
		panic(err)
	}
}
