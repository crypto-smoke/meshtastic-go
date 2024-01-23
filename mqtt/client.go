package mqtt

import (
	"errors"
	"github.com/charmbracelet/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"sync"
	"time"
)

type Client struct {
	server    string
	username  string
	password  string
	topicRoot string
	clientID  string
	client    mqtt.Client
	sync.RWMutex
	channelHandlers map[string][]HandlerFunc
}

type HandlerFunc func(message Message)

var DefaultClient = Client{
	server:    "tcp://mqtt.meshtastic.org:1883",
	username:  "meshdev",
	password:  "large4cats",
	topicRoot: "msh/2", //TODO: this will need to change

	channelHandlers: make(map[string][]HandlerFunc),
}

func NewClient(url, username, password, rootTopic string) *Client {
	return &Client{
		server:          url,
		username:        username,
		password:        password,
		topicRoot:       rootTopic,
		channelHandlers: make(map[string][]HandlerFunc),
	}
}

func (c *Client) TopicRoot() string {
	return c.topicRoot
}

func (c *Client) Connect() error {
	var alphabet = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	c.clientID = randomString(23, alphabet)

	mqtt.DEBUG = log.StandardLog(log.StandardLogOptions{ForceLevel: log.DebugLevel})
	mqtt.ERROR = log.StandardLog(log.StandardLogOptions{ForceLevel: log.ErrorLevel})
	opts := mqtt.NewClientOptions().
		AddBroker(c.server).
		SetUsername(c.username).
		SetOrderMatters(false).
		SetPassword(c.password).
		SetClientID(c.clientID).
		SetCleanSession(false)
	opts.SetKeepAlive(30 * time.Second)
	opts.SetResumeSubs(true)
	//opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(5 * time.Second)
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(1 * time.Minute)
	opts.SetConnectionLostHandler(func(c mqtt.Client, err error) {
		log.Error("mqtt connection lost", "err", err)
	})
	opts.SetReconnectingHandler(func(c mqtt.Client, options *mqtt.ClientOptions) {
		log.Info("mqtt reconnecting")
	})
	opts.SetOnConnectHandler(func(client mqtt.Client) {
		log.Info("connected to", "server", c.server)
	})
	c.client = mqtt.NewClient(opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// MQTT Message
type Message struct {
	Topic    string
	Payload  []byte
	Retained bool
}

// Publish a message to the broker
func (c *Client) Publish(m *Message) error {
	tok := c.client.Publish(m.Topic, 0, m.Retained, m.Payload)
	if !tok.WaitTimeout(10 * time.Second) {
		tok.Wait()
		return errors.New("timeout on mqtt publish")
	}
	if tok.Error() != nil {
		return tok.Error()
	}
	return nil
}

// Register a handler for messages on the specified channel
func (c *Client) Handle(channel string, h HandlerFunc) {
	c.Lock()
	defer c.Unlock()
	topic := c.GetFullTopicForChannel(channel)
	c.channelHandlers[channel] = append(c.channelHandlers[channel], h)
	c.client.Subscribe(topic+"/+", 0, c.handleBrokerMessage)
}
func (c *Client) GetFullTopicForChannel(channel string) string {
	return c.topicRoot + "/c/" + channel
}
func (c *Client) GetChannelFromTopic(topic string) string {
	trimmed := strings.TrimPrefix(topic, c.topicRoot+"/c/")
	sepIndex := strings.Index(trimmed, "/")
	if sepIndex > 0 {
		return trimmed[:sepIndex]
	}
	return trimmed
}
func (c *Client) handleBrokerMessage(client mqtt.Client, message mqtt.Message) {
	msg := Message{
		Topic:    message.Topic(),
		Payload:  message.Payload(),
		Retained: message.Retained(),
	}
	c.RLock()
	defer c.RUnlock()
	channel := c.GetChannelFromTopic(msg.Topic)
	chans := c.channelHandlers[channel]
	if len(chans) == 0 {
		log.Error("no handlers found", "topic", channel)
	}
	for _, ch := range chans {
		go ch(msg)
	}
}
