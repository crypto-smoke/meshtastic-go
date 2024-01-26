package main

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/crypto-smoke/meshtastic-go/emulated"
	"github.com/crypto-smoke/meshtastic-go/mqtt"
	"github.com/crypto-smoke/meshtastic-go/radio"
	"time"
)

func main() {
	ctx := context.Background()
	log.SetLevel(log.DebugLevel)

	key, err := radio.ParseKey("yGwxct335Pb8LdAt8Xk47_00bkbzEFQlwP45Id3HTjQ=")
	if err != nil {
		panic(err)
	}
	myNodeID := meshtastic.NodeID(3735928559)
	r, err := emulated.NewRadio(emulated.Config{
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
		err := r.ToRadio(ctx, &pb.ToRadio{
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
	if err := r.Run(ctx); err != nil {
		panic(err)
	}
}
