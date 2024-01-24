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
	NodeID     meshtastic.NodeID
	LongName   string
	ShortName  string
	Channels   *pb.ChannelSet
	MQTTClient *mqtt.Client
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
}

func NewRadio(cfg Config) (*Radio, error) {
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

	for _, ch := range r.cfg.Channels.Settings {
		r.logger.Debug("subscribing to mqtt for channel", "channel", ch.Name)
		r.mqtt.Handle(ch.Name, r.handleMQTTMessage)
	}

	// TODO: Rethink concurrency. Do we want a goroutine servicing ToRadio and one servicing FromRadio?

	eg, egCtx := errgroup.WithContext(ctx)
	// Spin up goroutine to send NodeInfo every interval
	eg.Go(func() error {
		// TODO: Make interval configurable
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()
		for {
			if err := r.sendNodeInfo(ctx); err != nil {
				r.logger.Error("failed to send node info", "err", err)
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

	// TODO: Dispatch to FromRadio subscribers
	// TODO: Do we need to attempt to decrypt before dispatching to FromRadio subscribers?
	r.dispatchMessageToFromRadio(&pb.FromRadio{
		PayloadVariant: &pb.FromRadio_Packet{
			Packet: meshPacket,
		},
	})

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
		r.logger.Info("updated nodeDB", "node_db", r.nodeDB)
		r.mu.Unlock()
	case pb.PortNum_TEXT_MESSAGE_APP:
		r.logger.Info("received TextMessage", "message", string(data.Payload))
	case pb.PortNum_ROUTING_APP:
		routingPayload := &pb.Routing{}
		if err := proto.Unmarshal(data.Payload, routingPayload); err != nil {
			return fmt.Errorf("unmarshalling routingPayload: %w", err)
		}
		r.logger.Info("received Routing", "routing", routingPayload)
	default:
		r.logger.Debug("received unhandled app payload", "data", data)
	}

	return nil
}

func (r *Radio) sendPacket(ctx context.Context, packet *pb.MeshPacket) error {
	// TODO: This does not necessarily send on the primary channel.
	// TODO: Do we try to encrypt the packet here?
	// TODO: Packet ID should be incremented for each packet sent.
	se := &pb.ServiceEnvelope{
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

func (r *Radio) sendNodeInfo(ctx context.Context) error {
	// TODO: Lots of stuff missing here. We should use encryption if possible and we should send this every configured
	// interval. However, this is enough for it to show in the UI of another node listening to the MQTT servr.
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
		// TODO: Calculate the correct To address to use here.
		To: 3806663900,
		PayloadVariant: &pb.MeshPacket_Decoded{
			Decoded: &pb.Data{
				Portnum: pb.PortNum_NODEINFO_APP,
				Payload: userBytes,
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
		// TODO: Do we need to decorate the outgoing packet.
		return r.sendPacket(ctx, payload.Packet)
	default:
		r.logger.Debug("unknown payload variant", "payload", payload)
	}
	return fmt.Errorf("not implemented")
}
