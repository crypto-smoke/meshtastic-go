package main

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go/transport"
	"github.com/crypto-smoke/meshtastic-go/transport/serial"
	"google.golang.org/protobuf/proto"
	"time"
)

func main() {
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

	client := transport.NewClient(streamConn, false)
	client.Handle(new(pb.FromRadio), func(msg proto.Message) {
		pkt := msg.(*pb.FromRadio)
		log.Info("Received message from radio", "msg", pkt)
	})
	if client.Connect() != nil {
		panic("Failed to connect to the radio")
	}

	time.Sleep(50 * time.Second)
}
