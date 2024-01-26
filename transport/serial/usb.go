package serial

import (
	"fmt"
	"github.com/charmbracelet/log"
	"go.bug.st/serial/enumerator"
)

type usbDevice struct {
	VID string
	PID string
}

var knownDevices = []usbDevice{
	{VID: "239A", PID: "8029"}, // rak4631_19003
}

func GetPorts() []string {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}
	var foundDevices []string
	if len(ports) == 0 {
		fmt.Println("No serial ports found!")
		return nil
	}
	for _, port := range ports {
		//fmt.Printf("Found port: %s\n", port.SettingName)
		if port.IsUSB {
			for _, device := range knownDevices {
				if device.VID != port.VID {
					continue
				}
				if device.PID != port.PID {
					continue
				}
				foundDevices = append(foundDevices, port.Name)
			}
		}
	}
	return foundDevices
}
