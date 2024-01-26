package emulated

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/crypto-smoke/meshtastic-go/mqtt"
	"github.com/crypto-smoke/meshtastic-go/radio"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/proto"
	"sync"
	"time"
)

type Config struct {
	// Dependencies
	MQTTClient *mqtt.Client

	// Node configuration
	// NodeID is the ID of the node.
	NodeID meshtastic.NodeID
	// LongName is the long name of the node.
	LongName string
	// ShortName is the short name of the node.
	ShortName string
	// Channels is the set of channels the radio will listen and transmit on.
	// The first channel in the set is considered the primary channel and is used for broadcasting NodeInfo and Position
	Channels *pb.ChannelSet
	// BroadcastInterval is the interval at which the radio will broadcast NodeInfo and Position messages.
	// Defaults to 5 minutes.
	BroadcastInterval time.Duration
}

func (c *Config) validate() error {
	if c.MQTTClient == nil {
		return fmt.Errorf("MQTTClient is required")
	}
	if c.NodeID == 0 {
		return fmt.Errorf("NodeID is required")
	}
	if c.LongName == "" {
		// TODO: Generate from NodeID
		return fmt.Errorf("LongName is required")
	}
	if c.ShortName == "" {
		// TODO: Generate from NodeID
		return fmt.Errorf("ShortName is required")
	}
	if c.Channels == nil {
		return fmt.Errorf("Channels is required")
	}
	if len(c.Channels.Settings) == 0 {
		return fmt.Errorf("Channels.Settings should be non-empty")
	}
	if c.BroadcastInterval == 0 {
		c.BroadcastInterval = 5 * time.Minute
	}
	return nil
}

// Radio emulates a meshtastic Node, communicating with a meshtastic network via MQTT.
type Radio struct {
	cfg    Config
	mqtt   *mqtt.Client
	logger *log.Logger

	// TODO: rwmutex?? seperate mutexes??
	mu                   sync.Mutex
	fromRadioSubscribers map[chan<- *pb.FromRadio]struct{}
	nodeDB               map[uint32]*pb.User
	// packetID is incremented and included in each packet sent from the radio.
	// TODO: Eventually, we should offer an easy way of persisting this so that we can resume from where we left off.
	packetID uint32
}

func NewRadio(cfg Config) (*Radio, error) {
	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}
	return &Radio{
		cfg:                  cfg,
		logger:               log.With("radio", cfg.NodeID.String()),
		fromRadioSubscribers: map[chan<- *pb.FromRadio]struct{}{},
		mqtt:                 cfg.MQTTClient,
		nodeDB:               map[uint32]*pb.User{},
	}, nil
}

func (r *Radio) Run(ctx context.Context) error {
	if err := r.mqtt.Connect(); err != nil {
		return fmt.Errorf("connecting to mqtt: %w", err)
	}
	// TODO: Disconnect??

	// Subscribe to all configured channels
	for _, ch := range r.cfg.Channels.Settings {
		r.logger.Debug("subscribing to mqtt for channel", "channel", ch.Name)
		r.mqtt.Handle(ch.Name, r.handleMQTTMessage)
	}

	// TODO: Rethink concurrency. Do we want a goroutine servicing ToRadio and one servicing FromRadio?

	eg, egCtx := errgroup.WithContext(ctx)
	// Spin up goroutine to send NodeInfo every interval
	eg.Go(func() error {
		ticker := time.NewTicker(r.cfg.BroadcastInterval)
		defer ticker.Stop()
		for {
			if err := r.broadcastNodeInfo(ctx); err != nil {
				r.logger.Error("failed to broadcast node info", "err", err)
			}
			if err := r.broadcastPosition(ctx); err != nil {
				r.logger.Error("failed to broadcast position", "err", err)
			}
			select {
			case <-egCtx.Done():
				return nil
			case <-ticker.C:
			}
		}
	})

	return eg.Wait()
}

func (r *Radio) handleMQTTMessage(msg mqtt.Message) {
	// TODO: Determine how "github.com/eclipse/paho.mqtt.golang" handles concurrency. Do we need to dispatch here to
	// a goroutine which handles incoming messages to unblock this one?
	if err := r.tryHandleMQTTMessage(msg); err != nil {
		r.logger.Error("failed to handle incoming mqtt message", "err", err)
	}
}

func (r *Radio) tryHandleMQTTMessage(msg mqtt.Message) error {
	serviceEnvelope := &pb.ServiceEnvelope{}
	if err := proto.Unmarshal(msg.Payload, serviceEnvelope); err != nil {
		return fmt.Errorf("unmarshalling: %w", err)
	}
	meshPacket := serviceEnvelope.Packet

	// TODO: Attempt decryption first before dispatching to subscribers
	// TODO: This means we move this further below.
	if err := r.dispatchMessageToFromRadio(&pb.FromRadio{
		PayloadVariant: &pb.FromRadio_Packet{
			Packet: meshPacket,
		},
	}); err != nil {
		r.logger.Error("failed to dispatch message to FromRadio subscribers", "err", err)
	}

	// From now on, we only care about messages on the primary channel
	primaryName := r.cfg.Channels.Settings[0].Name
	primaryPSK := r.cfg.Channels.Settings[0].Psk
	if serviceEnvelope.ChannelId != primaryName {
		return nil
	}

	r.logger.Debug("received service envelope for primary channel", "serviceEnvelope", serviceEnvelope)
	// Check if we should try and decrypt the message
	var data *pb.Data
	switch payload := meshPacket.PayloadVariant.(type) {
	case *pb.MeshPacket_Decoded:
		data = payload.Decoded
	case *pb.MeshPacket_Encrypted:
		// TODO: Check if we have the key for this channel
		plaintext, err := radio.XOR(
			payload.Encrypted,
			primaryPSK,
			meshPacket.Id,
			meshPacket.From,
		)
		if err != nil {
			return fmt.Errorf("decrypting: %w", err)
		}
		data = &pb.Data{}
		if err := proto.Unmarshal(plaintext, data); err != nil {
			return fmt.Errorf("unmarshalling decrypted data: %w", err)
		}
	default:
		return fmt.Errorf("unknown payload variant %T", payload)
	}
	r.logger.Debug("received data for primary channel", "data", data)

	// For messages on the primary channel, we want to handle these and potentially update the nodeDB.
	switch data.Portnum {
	case pb.PortNum_NODEINFO_APP:
		user := &pb.User{}
		if err := proto.Unmarshal(data.Payload, user); err != nil {
			return fmt.Errorf("unmarshalling user: %w", err)
		}
		r.logger.Info("received NodeInfo", "user", user)
		// Update NodeDB
		r.mu.Lock()
		r.nodeDB[meshPacket.From] = user
		r.mu.Unlock()
	case pb.PortNum_TEXT_MESSAGE_APP:
		r.logger.Info("received TextMessage", "message", string(data.Payload))
	case pb.PortNum_ROUTING_APP:
		routingPayload := &pb.Routing{}
		if err := proto.Unmarshal(data.Payload, routingPayload); err != nil {
			return fmt.Errorf("unmarshalling routingPayload: %w", err)
		}
		r.logger.Info("received Routing", "routing", routingPayload)
	case pb.PortNum_POSITION_APP:
		positionPayload := &pb.Position{}
		if err := proto.Unmarshal(data.Payload, positionPayload); err != nil {
			return fmt.Errorf("unmarshalling positionPayload: %w", err)
		}
		r.logger.Info("received Position", "position", positionPayload)
	default:
		r.logger.Debug("received unhandled app payload", "data", data)
	}

	return nil
}

func (r *Radio) nextPacketID() uint32 {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.packetID++
	return r.packetID
}

func (r *Radio) sendPacket(ctx context.Context, packet *pb.MeshPacket) error {
	// TODO: Optimistically attempt to encrypt the packet here if we recognise the channel, encryption is enabled and
	// the payload is not currently encrypted.

	// sendPacket is responsible for setting the packet ID.
	r.packetID = r.nextPacketID()

	se := &pb.ServiceEnvelope{
		// TODO: Fetch channel to use based on packet.Channel rather than hardcoding to primary channel.
		ChannelId: r.cfg.Channels.Settings[0].Name,
		GatewayId: r.cfg.NodeID.String(),
		Packet:    packet,
	}
	bytes, err := proto.Marshal(se)
	if err != nil {
		return fmt.Errorf("marshalling service envelope: %w", err)
	}
	return r.mqtt.Publish(&mqtt.Message{
		Topic:   r.mqtt.GetFullTopicForChannel(r.cfg.Channels.Settings[0].Name) + "/" + r.cfg.NodeID.String(),
		Payload: bytes,
	})
}

func (r *Radio) broadcastNodeInfo(ctx context.Context) error {
	r.logger.Info("broadcasting NodeInfo")
	// TODO: Lots of stuff missing here. However, this is enough for it to show in the UI of another node listening to
	// the MQTT server.
	user := &pb.User{
		Id:        r.cfg.NodeID.String(),
		LongName:  r.cfg.LongName,
		ShortName: r.cfg.ShortName,
		HwModel:   pb.HardwareModel_PRIVATE_HW,
	}
	userBytes, err := proto.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshalling user: %w", err)
	}
	return r.sendPacket(ctx, &pb.MeshPacket{
		From: r.cfg.NodeID.Uint32(),
		To:   meshtastic.BroadcastNodeID.Uint32(),
		PayloadVariant: &pb.MeshPacket_Decoded{
			Decoded: &pb.Data{
				Portnum: pb.PortNum_NODEINFO_APP,
				Payload: userBytes,
			},
		},
	})
}

func (r *Radio) broadcastPosition(ctx context.Context) error {
	r.logger.Info("broadcasting Position")
	position := &pb.Position{
		// TODO: Make broadcasting positional optional and configurable.
		// Degrees divided by 1e-7 gives the value to use here.
		// Buckingham Palace :D
		LatitudeI:  515014760,
		LongitudeI: -1406340,
		Time:       uint32(time.Now().Unix()),
	}
	positionBytes, err := proto.Marshal(position)
	if err != nil {
		return fmt.Errorf("marshalling position: %w", err)
	}
	return r.sendPacket(ctx, &pb.MeshPacket{
		From: r.cfg.NodeID.Uint32(),
		To:   meshtastic.BroadcastNodeID.Uint32(),
		PayloadVariant: &pb.MeshPacket_Decoded{
			Decoded: &pb.Data{
				Portnum: pb.PortNum_POSITION_APP,
				Payload: positionBytes,
			},
		},
	})
}

// dispatchMessageToFromRadio sends a FromRadio message to all current subscribers to
// the FromRadio.
func (r *Radio) dispatchMessageToFromRadio(msg *pb.FromRadio) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for ch := range r.fromRadioSubscribers {
		// TODO: Make this way safer/resilient
		ch <- msg
	}
	return nil
}

func (r *Radio) FromRadio(ctx context.Context, ch chan<- *pb.FromRadio) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fromRadioSubscribers[ch] = struct{}{}
	// TODO: Unsubscribe from the channel when the context is cancelled??
	return nil
}

func (r *Radio) ToRadio(ctx context.Context, msg *pb.ToRadio) error {
	switch payload := msg.PayloadVariant.(type) {
	case *pb.ToRadio_Disconnect:
		r.logger.Info("received Disconnect from ToRadio")
	case *pb.ToRadio_Packet:
		r.logger.Info("received Packet from ToRadio")
		return r.sendPacket(ctx, payload.Packet)
	default:
		r.logger.Debug("unknown payload variant", "payload", payload)
	}
	return fmt.Errorf("not implemented")
}
