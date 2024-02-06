package transport

import (
	"buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log/slog"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrTimeout = errors.New("timeout connecting to radio")
)

type HandlerFunc func(message proto.Message)

type Client struct {
	sc       *StreamConn
	handlers *HandlerRegistry
	log      *slog.Logger

	State State
}

type State struct {
	sync.RWMutex
	complete       bool
	configID       uint32
	nodeInfo       *meshtastic.MyNodeInfo
	deviceMetadata *meshtastic.DeviceMetadata
	nodes          []*meshtastic.NodeInfo
	channels       []*meshtastic.Channel
	configs        []*meshtastic.Config
	modules        []*meshtastic.ModuleConfig
}

func (s *State) Complete() bool {
	s.RLock()
	defer s.RUnlock()
	return s.complete
}

func (s *State) ConfigID() uint32 {
	s.RLock()
	defer s.RUnlock()
	return s.configID
}

func (s *State) NodeInfo() *meshtastic.MyNodeInfo {
	s.RLock()
	defer s.RUnlock()
	return s.nodeInfo
}

func (s *State) DeviceMetadata() *meshtastic.DeviceMetadata {
	s.RLock()
	defer s.RUnlock()
	return s.deviceMetadata
}

func (s *State) Nodes() []*meshtastic.NodeInfo {
	s.RLock()
	defer s.RUnlock()
	return s.nodes
}

func (s *State) Channels() []*meshtastic.Channel {
	s.RLock()
	defer s.RUnlock()
	return s.channels
}

func (s *State) Configs() []*meshtastic.Config {
	s.RLock()
	defer s.RUnlock()
	return s.configs
}

func (s *State) Modules() []*meshtastic.ModuleConfig {
	s.RLock()
	defer s.RUnlock()
	return s.modules
}

func (s *State) SetComplete(complete bool) {
	s.Lock()
	defer s.Unlock()
	s.complete = complete
}

func (s *State) SetConfigID(configID uint32) {
	s.Lock()
	defer s.Unlock()
	s.configID = configID
}

func (s *State) SetNodeInfo(nodeInfo *meshtastic.MyNodeInfo) {
	s.Lock()
	defer s.Unlock()
	s.nodeInfo = nodeInfo
}

func (s *State) SetDeviceMetadata(deviceMetadata *meshtastic.DeviceMetadata) {
	s.Lock()
	defer s.Unlock()
	s.deviceMetadata = deviceMetadata
}

func (s *State) AddNode(node *meshtastic.NodeInfo) {
	s.Lock()
	defer s.Unlock()
	s.nodes = append(s.nodes, node)
}

func (s *State) AddChannel(channel *meshtastic.Channel) {
	s.Lock()
	defer s.Unlock()
	s.channels = append(s.channels, channel)
}

func (s *State) AddConfig(config *meshtastic.Config) {
	s.Lock()
	defer s.Unlock()
	s.configs = append(s.configs, config)
}

func (s *State) AddModules(module *meshtastic.ModuleConfig) {
	s.Lock()
	defer s.Unlock()
	s.modules = append(s.modules, module)
}

func NewClient(sc *StreamConn, errorOnNoHandler bool) *Client {
	return &Client{
		log:      slog.Default().WithGroup("client"),
		sc:       sc,
		handlers: NewHandlerRegistry(errorOnNoHandler),
	}
}

// You have to send this first to get the radio into protobuf mode and have it accept and send packets via serial
func (c *Client) sendGetConfig() error {
	r := rand.Uint32()
	c.State.configID = r
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

func (c *Client) Connect(ctx context.Context) error {
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
				c.State.SetNodeInfo(msg.GetMyInfo())
				variant = c.State.nodeInfo
			case *meshtastic.FromRadio_Metadata:
				c.State.SetDeviceMetadata(msg.GetMetadata())
				variant = c.State.deviceMetadata
			case *meshtastic.FromRadio_NodeInfo:
				node := msg.GetNodeInfo()
				c.State.AddNode(node)
				variant = node
			case *meshtastic.FromRadio_Channel:
				channel := msg.GetChannel()
				c.State.AddChannel(channel)
				variant = channel
			case *meshtastic.FromRadio_Config:
				cfg := msg.GetConfig()
				c.State.AddConfig(cfg)
				variant = cfg
			case *meshtastic.FromRadio_ModuleConfig:
				cfg := msg.GetModuleConfig()
				c.State.AddModules(cfg)
				variant = cfg
			case *meshtastic.FromRadio_ConfigCompleteId:
				// logged here because it's not an actual proto.Message that we can call handlers on
				c.log.Debug("config complete")
				c.State.SetComplete(true)
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
				// logged here because it's not an actual proto.Message that we can call handlers on
				c.log.Debug("rebooted", "rebooted", msg.GetRebooted())

				continue
			case *meshtastic.FromRadio_XmodemPacket:
				variant = msg.GetXmodemPacket()
			case *meshtastic.FromRadio_Packet:
				variant = msg.GetPacket()
			default:
				c.log.Warn("unhandled protobuf from radio")
			}

			if !c.State.Complete() {
				continue
			}
			err = c.handlers.HandleMessage(variant)
			if err != nil {
				c.log.Error("error handling message", "err", err)
			}
		}
	}()
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ErrTimeout
		case <-ticker.C:
			done := c.State.complete
			if done {
				return nil
			}
		}
	}
}
