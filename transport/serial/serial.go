package serial

import (
	"buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"encoding/binary"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go/transport"
	"go.bug.st/serial"
	"google.golang.org/protobuf/proto"
	"io"
	"math/rand"
	"time"
)

const (
	WAIT_AFTER_WAKE = 100 * time.Millisecond
	START1          = 0x94
	START2          = 0xc3
	PACKET_MTU      = 512
	PORT_SPEED      = 115200 //921600
)

type HandlerFunc func(message proto.Message)

// Serial connection to a node
type Conn struct {
	serialPort string
	serialConn serial.Port
	handlers   *transport.HandlerRegistry

	config struct {
		complete bool
		configID uint32
		*meshtastic.MyNodeInfo
		*meshtastic.DeviceMetadata
		nodes    []*meshtastic.NodeInfo
		channels []*meshtastic.Channel
		config   []*meshtastic.Config
		modules  []*meshtastic.ModuleConfig
	}
}

func NewConn(port string, errorOnNoHandler bool) *Conn {
	var c = Conn{serialPort: port,
		handlers: transport.NewHandlerRegistry(errorOnNoHandler)}
	return &c
}

// You have to send this first to get the radio into protobuf mode and have it accept and send packets via serial
func (c *Conn) sendGetConfig() {
	r := rand.Uint32()
	c.config.configID = r
	//log.Info("want config id", r)
	msg := &meshtastic.ToRadio{
		PayloadVariant: &meshtastic.ToRadio_WantConfigId{
			WantConfigId: r,
		},
	}
	c.SendToRadio(msg)
}
func (c *Conn) Handle(kind proto.Message, handler transport.MessageHandler) {
	c.handlers.RegisterHandler(kind, handler)
}

func (c *Conn) Connect() error {
	mode := &serial.Mode{
		BaudRate: PORT_SPEED,
	}
	port, err := serial.Open(c.serialPort, mode)
	if err != nil {
		return err
	}
	c.serialConn = port
	ch := make(chan *meshtastic.FromRadio)
	go c.decodeProtos(false, ch)
	go func() {
		for {
			msg := <-ch
			var variant proto.Message
			switch msg.GetPayloadVariant().(type) {
			// These pbufs all get sent upon initial connection to the node
			case *meshtastic.FromRadio_MyInfo:
				c.config.MyNodeInfo = msg.GetMyInfo()
				variant = c.config.MyNodeInfo
			case *meshtastic.FromRadio_Metadata:
				c.config.DeviceMetadata = msg.GetMetadata()
				variant = c.config.DeviceMetadata
			case *meshtastic.FromRadio_NodeInfo:
				node := msg.GetNodeInfo()
				c.config.nodes = append(c.config.nodes, node)
				variant = node
			case *meshtastic.FromRadio_Channel:
				channel := msg.GetChannel()
				c.config.channels = append(c.config.channels, channel)
				variant = channel
			case *meshtastic.FromRadio_Config:
				cfg := msg.GetConfig()
				c.config.config = append(c.config.config, cfg)
				variant = cfg
			case *meshtastic.FromRadio_ModuleConfig:
				cfg := msg.GetModuleConfig()
				c.config.modules = append(c.config.modules, cfg)
				variant = cfg
			case *meshtastic.FromRadio_ConfigCompleteId:
				// done getting config info
				//fmt.Println("config complete")
				c.config.complete = true
				/*
					out, err := json.MarshalIndent(c.config, "", "  ")
					if err != nil {
						log.Error("failed marshalling", "err", err)
						continue
					}
					fmt.Println(string(out))
					out, err = json.MarshalIndent(c.config.config, "", "  ")
					if err != nil {
						log.Error("failed marshalling", "err", err)
						continue
					}
					fmt.Println(string(out))

				*/
				continue
				// below are packets not part of initial connection

			case *meshtastic.FromRadio_LogRecord:
				variant = msg.GetLogRecord()
			case *meshtastic.FromRadio_MqttClientProxyMessage:
				variant = msg.GetMqttClientProxyMessage()
			case *meshtastic.FromRadio_QueueStatus:
				variant = msg.GetQueueStatus()
			case *meshtastic.FromRadio_Rebooted:
				// true if radio just rebooted
				fmt.Print("rebooted", msg.GetRebooted())
				continue
			case *meshtastic.FromRadio_XmodemPacket:
				variant = msg.GetXmodemPacket()

			case *meshtastic.FromRadio_Packet:
				variant = msg.GetPacket()
			default:
				log.Error("unhandled protobuf from radio")
			}
			if !c.config.complete {
				continue
			}
			err = c.handlers.HandleMessage(variant)
			if err != nil {
				log.Error("error handling message", "err", err)
			}
		}
	}()

	c.sendGetConfig()
	return nil
}
func (c *Conn) ConnectOld(ch chan *meshtastic.FromRadio, ch2 chan *meshtastic.ToRadio) error {
	mode := &serial.Mode{
		BaudRate: PORT_SPEED,
	}
	port, err := serial.Open(c.serialPort, mode)
	if err != nil {
		return err
	}
	c.serialConn = port

	go c.decodeProtos(false, ch)
	go func() {
		for {
			msg := <-ch2
			c.SendToRadio(msg)
		}
	}()
	c.sendGetConfig()
	return nil
}

var txtCh = make(chan *meshtastic.MeshPacket)

func (c *Conn) decodeProtos(printDebug bool, ch chan *meshtastic.FromRadio) {
	for {
		data, err := readUntilProtobuf(c.serialConn, printDebug)
		if err != nil {
			log.Info("error:", err)
			continue
		}
		//log.Info("read from serial and got proto")
		var msg2 meshtastic.FromRadio
		err = proto.Unmarshal(data, &msg2)
		if err != nil {
			log.Fatal(err)
		}
		ch <- &msg2
	}
}
func readUntilProtobuf(reader io.Reader, printDebug bool) ([]byte, error) {
	buf := make([]byte, 4)
	for {
		// Read the first byte, looking for START1.
		_, err := io.ReadFull(reader, buf[:1])
		if err != nil {
			return nil, err
		}

		// Check for START1.
		if buf[0] != 0x94 {
			if printDebug {
				fmt.Print(string(buf[0]))
			}
			continue
		}

		// Read the second byte, looking for START2.
		_, err = io.ReadFull(reader, buf[1:2])
		if err != nil {
			return nil, err
		}

		// Check for START2.
		if buf[1] != 0xc3 {
			continue
		}

		// The next two bytes should be the length of the protobuf message.
		_, err = io.ReadFull(reader, buf[2:])
		if err != nil {
			return nil, err
		}

		length := int(binary.BigEndian.Uint16(buf[2:]))
		if length > PACKET_MTU {
			//packet corrupt, start over
			continue
		}
		//fmt.Println("got packet from node with length", length)
		data := make([]byte, length)

		// Read the protobuf data.
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
}

func (c *Conn) flushPort() error {
	flush := make([]byte, 32, 32)
	for j := 0; j < len(flush); j++ {
		flush[j] = START2
	}
	_, err := c.serialConn.Write(flush)
	if err != nil {
		return err
	}
	return nil
}
func (c *Conn) SendToRadio(msg *meshtastic.ToRadio) error {
	err := c.flushPort()
	if err != nil {
		return err
	}
	//fmt.Printf("Sent %v bytes\n", n)
	data, err := proto.Marshal(msg)
	if err != nil {
		panic(err)
	}
	time.Sleep(WAIT_AFTER_WAKE)

	datalen := len(data)
	header := []byte{START1, START2, byte(datalen >> 8), byte(datalen)}
	data = append(header, data...)
	_, err = c.serialConn.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Sent %v bytes\n", n)
	return nil
}
