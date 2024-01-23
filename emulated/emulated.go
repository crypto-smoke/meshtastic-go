package emulated

import (
	pb "buf.build/gen/go/meshtastic/protobufs/protocolbuffers/go/meshtastic"
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/crypto-smoke/meshtastic-go"
	"github.com/crypto-smoke/meshtastic-go/mqtt"
	"github.com/crypto-smoke/meshtastic-go/radio"
	"google.golang.org/protobuf/proto"
	"sync"
)

type Config struct {
	ID         string
	Channels   *pb.ChannelSet
	MQTTClient *mqtt.Client
}

// Radio emulates a meshtastic Node, communicating with a meshtastic network via MQTT.
type Radio struct {
	cfg    Config
	mqtt   *mqtt.Client
	logger *log.Logger

	mu                   sync.Mutex
	fromRadioSubscribers map[chan<- *pb.FromRadio]struct{}
	nodeDB               map[uint32]*pb.User
}

func NewRadio(cfg Config) (*Radio, error) {
	return &Radio{
		cfg:    cfg,
		logger: log.With("radio", cfg.ID),
		mqtt:   cfg.MQTTClient,
		nodeDB: map[uint32]*pb.User{},
	}, nil
}

func (r *Radio) Run(ctx context.Context) error {
	if err := r.mqtt.Connect(); err != nil {
		return fmt.Errorf("connecting to mqtt: %w", err)
	}

	for _, ch := range r.cfg.Channels.Settings {
		r.logger.Debug("subscribing to mqtt for channel", "channel", ch.Name)
		r.mqtt.Handle(ch.Name, r.handleMQTTMessage)
	}

	if err := r.sendNodeInfo(ctx); err != nil {
		r.logger.Error("failed to send node info", "err", err)
	}

	<-ctx.Done()
	return nil
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
		r.logger.Info("updated nodeDB", "node_count", len(r.nodeDB))
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

func (r *Radio) sendNodeInfo(ctx context.Context) error {
	// TODO: Lots of stuff missing here. We should use encryption if possible and we should send this every configured
	// interval. However, this is enough for it to show in the UI of another node listening to the MQTT servr.
	nodeID := meshtastic.NodeID(6969420)

	user := &pb.User{
		Id:        nodeID.String(),
		LongName:  "M7NOA - GPHR",
		ShortName: "GPHR",
		HwModel:   pb.HardwareModel_PRIVATE_HW,
	}
	userBytes, err := proto.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshalling user: %w", err)
	}
	se := &pb.ServiceEnvelope{
		ChannelId: r.cfg.Channels.Settings[0].Name,
		GatewayId: nodeID.String(),
		Packet: &pb.MeshPacket{
			From: nodeID.Uint32(),
			To:   3806663900,
			PayloadVariant: &pb.MeshPacket_Decoded{
				Decoded: &pb.Data{
					Portnum: pb.PortNum_NODEINFO_APP,
					Payload: userBytes,
				},
			},
		},
	}
	bytes, err := proto.Marshal(se)
	if err != nil {
		return fmt.Errorf("marshalling se: %w", err)
	}
	return r.mqtt.Publish(&mqtt.Message{
		Topic:   r.mqtt.GetFullTopicForChannel(r.cfg.Channels.Settings[0].Name) + "/" + nodeID.String(),
		Payload: bytes,
	})
}

// dispatchMessageToFromRadio sends a FromRadio message to all current subscribers to
// the FromRadio.
func (r *Radio) dispatchMessageToFromRadio(msg *pb.FromRadio) error {
	return fmt.Errorf("not implemented")
}

func (r *Radio) FromRadio(ctx context.Context, ch chan<- *pb.FromRadio) error {
	return fmt.Errorf("not implemented")
}

func (r *Radio) ToRadio(ctx context.Context, msg *pb.ToRadio) error {
	return fmt.Errorf("not implemented")
}
