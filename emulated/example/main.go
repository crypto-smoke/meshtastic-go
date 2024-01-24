package main

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"encoding/base64"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/crypto-smoke/meshtastic-go/emulated"
	"github.com/crypto-smoke/meshtastic-go/mqtt"
)

func main() {
	log.SetLevel(log.DebugLevel)
	key, err := base64.URLEncoding.DecodeString("yGwxct335Pb8LdAt8Xk47_00bkbzEFQlwP45Id3HTjQ=")
	if err != nil {
		panic(err)
	}
	radio, err := emulated.NewRadio(emulated.Config{
		LongName:   "M7NOA - GPHR",
		ShortName:  "GPHR",
		NodeID:     meshtastic.NodeID(6969420),
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
	if err := radio.Run(context.TODO()); err != nil {
		panic(err)
	}
}
