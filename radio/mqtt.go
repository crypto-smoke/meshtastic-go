package radio

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

// some debugging junk that needs to be deleted
type FakeRadio struct {
	ID uint32
}

func NewFakeRadio() (*FakeRadio, error) {
	opts := mqtt.NewClientOptions().AddBroker("tcp://mqtt.meshtastic.org:1883").SetClientID("poopypants").SetUsername("meshdev").SetPassword("large4cats")
	opts.SetKeepAlive(2 * time.Second)
	//opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	ugh := c.Subscribe("msh/2/c/#", 0, func(client mqtt.Client, message mqtt.Message) {

	})
	_ = ugh

	return nil, nil
}
