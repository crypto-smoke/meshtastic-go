package transport

import (
	"buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"fmt"
	"github.com/charmbracelet/log"
	"google.golang.org/protobuf/proto"
	"math/rand"
)

type HandlerFunc func(message proto.Message)

type Client struct {
	sc       *StreamConn
	handlers *HandlerRegistry
	log      *log.Logger

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

func NewClient(sc *StreamConn, errorOnNoHandler bool) *Client {
	return &Client{
		log:      log.WithPrefix("client"),
		sc:       sc,
		handlers: NewHandlerRegistry(errorOnNoHandler),
	}
}

// You have to send this first to get the radio into protobuf mode and have it accept and send packets via serial
func (c *Client) sendGetConfig() error {
	r := rand.Uint32()
	c.config.configID = r
	msg := &meshtastic.ToRadio{
		PayloadVariant: &meshtastic.ToRadio_WantConfigId{
			WantConfigId: r,
		},
	}
	c.log.Debug("sending want config", "id", r)
	if err := c.sc.Write(msg); err != nil {
		return fmt.Errorf("writing want config command: %w", err)
	}
	c.log.Debug("sent want config")
	return nil
}

func (c *Client) Handle(kind proto.Message, handler MessageHandler) {
	c.handlers.RegisterHandler(kind, handler)
}

func (c *Client) SendToRadio(msg *meshtastic.ToRadio) error {
	return c.sc.Write(msg)
}

func (c *Client) Connect() error {
	if err := c.sendGetConfig(); err != nil {
		return fmt.Errorf("requesting config: %w", err)
	}
	go func() {
		for {
			msg := &meshtastic.FromRadio{}
			err := c.sc.Read(msg)
			if err != nil {
				c.log.Error("error reading from radio", "err", err)
				continue
			}
			c.log.Debug("received message from radio", "msg", msg)
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
	return nil
}
