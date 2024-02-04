package main

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go/transport"
	"github.com/crypto-smoke/meshtastic-go/transport/serial"
	"google.golang.org/protobuf/proto"
	"os"
	"os/signal"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.SetLevel(log.DebugLevel)

	ports := serial.GetPorts()
	serialConn, err := serial.Connect(ports[0])
	if err != nil {
		panic(err)
	}
	streamConn, err := transport.NewClientStreamConn(serialConn)
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := streamConn.Close(); err != nil {
			panic(err)
		}
	}()

	client := transport.NewClient(streamConn, false)
	client.Handle(new(pb.MeshPacket), func(msg proto.Message) {
		pkt := msg.(*pb.MeshPacket)
		log.Info("Received message from radio", "msg", pkt)
	})
	if client.Connect() != nil {
		panic("Failed to connect to the radio")
	}

	log.Info("Waiting for interrupt signal")
	<-ctx.Done()
}
