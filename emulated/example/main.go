package main

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/crypto-smoke/meshtastic-go/emulated"
	"github.com/crypto-smoke/meshtastic-go/mqtt"
	"github.com/crypto-smoke/meshtastic-go/radio"
	"golang.org/x/sync/errgroup"
	"time"
)

func main() {
	// TODO: Flesh this example out and make it configurable
	ctx := context.Background()
	log.SetLevel(log.DebugLevel)

	myNodeID := meshtastic.NodeID(3735928559)
	r, err := emulated.NewRadio(emulated.Config{
		LongName:   "EXAMPLE",
		ShortName:  "EMPL",
		NodeID:     myNodeID,
		MQTTClient: &mqtt.DefaultClient,
		Channels: &pb.ChannelSet{
			Settings: []*pb.ChannelSettings{
				{
					Name: "LongFast",
					Psk:  radio.DefaultKey,
				},
			},
		},
		BroadcastNodeInfoInterval: 5 * time.Minute,

		BroadcastPositionInterval: 5 * time.Minute,
		// Hardcoded to the position of Buckingham Palace.
		PositionLatitudeI:  515014760,
		PositionLongitudeI: -1406340,
		PositionAltitude:   2,
	})
	if err != nil {
		panic(err)
	}

	eg, egCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := r.Run(egCtx); err != nil {
			return fmt.Errorf("running radio: %w", err)
		}
		return nil
	})

	eg.Go(func() error {
		// Forgive me, for I have sinned.
		// TODO: We need a way of knowing the radio has come up and is ready that's better than waiting ten seconds.
		select {
		case <-egCtx.Done():
			return nil
		case <-time.After(10 * time.Second):
		}

		err := r.ToRadio(egCtx, &pb.ToRadio{
			PayloadVariant: &pb.ToRadio_Packet{
				Packet: &pb.MeshPacket{
					From: myNodeID.Uint32(),
					// This is hard coded to Noah's node ID
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
			return fmt.Errorf("sending to radio: %w", err)
		}

		return nil
	})

	eg.Go(func() error {
		ch := make(chan *pb.FromRadio)
		defer close(ch)
		err := r.FromRadio(egCtx, ch)
		if err != nil {
			return fmt.Errorf("setting up FromRadio subscriber: %w", err)
		}

		for {
			select {
			case <-egCtx.Done():
				return nil
			case fromRadio := <-ch:
				log.Info("FromRadio!!", "packet", fromRadio)
			}
		}
	})

	if err := eg.Wait(); err != nil {
		panic(err)
	}
}
