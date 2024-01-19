package serial

import (
	"buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"encoding/binary"
	"fmt"
	"github.com/charmbracelet/log"
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

// Serial connection to a node
type Conn struct {
	serialPort string
	serialConn serial.Port
}

func NewConn(port string) *Conn {
	var c = Conn{serialPort: port}
	return &c
}

// You have to send this first to get the radio into protobuf mode and have it accept and send packets via serial
func (c *Conn) sendGetConfig() {
	r := rand.Uint32()

	//log.Info("want config id", r)
	msg := &meshtastic.ToRadio{
		PayloadVariant: &meshtastic.ToRadio_WantConfigId{
			WantConfigId: r,
		},
	}
	c.sendToRadio(msg)
}

// TODO: refactor this to lose the channels
func (c *Conn) Connect(ch chan *meshtastic.FromRadio, ch2 chan *meshtastic.ToRadio) error {
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
			c.sendToRadio(msg)
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
		continue

		// i think everthing below was just crap for initial debugging and testing. will cleanup (TODO)
		//fmt.Println("proto message parsed")

		switch msg2.PayloadVariant.(type) {
		case *meshtastic.FromRadio_Packet:
			//log.Debug("packet recvd")
			pkt := msg2.GetPacket()
			payload := pkt.GetPayloadVariant()
			switch payload.(type) {
			case *meshtastic.MeshPacket_Decoded:
				//log.Info("got decoded packet")
				data := pkt.GetDecoded()
				switch data.Portnum {
				case meshtastic.PortNum_TEXT_MESSAGE_APP:
					go func() {
						txtCh <- pkt
					}()
					//fmt.Println("payload:", string(data.Payload))
					//fmt.Println("to:", pkt.To)
				case meshtastic.PortNum_ADMIN_APP:
					var adminMsg meshtastic.AdminMessage
					err := proto.Unmarshal(data.Payload, &adminMsg)
					if err != nil {
						log.Fatal(err)
					}
					//	log.Println("admin msg", adminMsg)
					switch adminMsg.GetPayloadVariant().(type) {
					case *meshtastic.AdminMessage_GetDeviceMetadataResponse:
						//log.Info("metadata:", adminMsg.GetGetDeviceMetadataResponse())
					}
				}
				//log.Info("data", "packet", data)
			case *meshtastic.MeshPacket_Encrypted:
				log.Info("got encrypted packet")
			}
			//log.Info("got packet", pkt)
		case *meshtastic.FromRadio_Metadata:
			log.Info("got metadata", msg2.GetMetadata())

		case *meshtastic.FromRadio_ConfigCompleteId:
			log.Info("config complete", msg2.GetConfigCompleteId())
		case *meshtastic.FromRadio_MyInfo: //*generated.deleteme.FromRadio_MyInfo:
		//log.WithField("my_node", msg.GetMyInfo()).Debug("got my node info")
		//fmt.Println(msg2.GetMyInfo().MyNodeNum)
		/*
			case *message.FromRadio_Radio:
				log.WithField("radio", msg.GetRadio()).Debug("got radio config")
				m.radioConfig = msg.GetRadio()
			case *message.FromRadio_NodeInfo:
				log.WithField("node", msg.GetNodeInfo()).Debug("got node info")
				m.pub(TOPIC_NODE, msg.GetNodeInfo())
			case *message.FromRadio_ConfigCompleteId:
				// TODO: implement this
				log.WithField("node", msg.GetConfigCompleteId()).Debug("got config complete")

		*/
		case *meshtastic.FromRadio_MqttClientProxyMessage:
			mqttMessage := msg2.GetMqttClientProxyMessage()

			switch payload := mqttMessage.GetPayloadVariant().(type) {
			case *meshtastic.MqttClientProxyMessage_Data:
				//payload.Data
				fmt.Println("got data payload")
			case *meshtastic.MqttClientProxyMessage_Text:
				fmt.Println("payload text", payload.Text)
			}
			log.Info("got mqtt client proxy")

		default:
			//log.Info("fromRadio", msg2.GetPayloadVariant())
		}
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
func (c *Conn) sendToRadio(msg *meshtastic.ToRadio) error {
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
