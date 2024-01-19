package lora

import "fmt"

// translated from https://sensing-labs.com/f-a-q/a-good-radio-level/

// Define signal quality and diagnostic notes.
type signalQuality string
type diagnosticNote string

const (
	Good signalQuality = "GOOD"
	Fair signalQuality = "FAIR"
	Bad  signalQuality = "BAD"
)

// getSignalQuality determines the signal quality based on RSSI and SNR.
func getSignalQuality(rssi, snr float64) signalQuality {
	// Define the boundaries for GOOD signal quality
	if snr >= -7 && rssi >= -115 {
		return Good
	}
	// Define the boundaries for FAIR signal quality
	if snr >= -15 && rssi >= -126 {
		return Fair
	}
	// If none of the above conditions are met, signal is BAD
	return Bad
}

// getDiagnosticNotes provides recommendations based on RSSI and SNR values.
func getDiagnosticNotes(rssi, snr float64) diagnosticNote {
	if rssi >= -115 && snr >= -7 {
		return "RF level is optimal to get a good reception reliability."
	} else if rssi >= -126 && snr >= -15 {
		return "RF level is not optimal but must be sufficient. Try to improve your device position if possible. You will have to monitor the stability of the RF level."
	} else {
		return "NOISY environment. Try to put device out of electromagnetic sources."
	}
}

func Demo() {
	// Example usage
	rssi := -120.0 // RSSI value
	snr := -10.0   // SNR value

	quality := getSignalQuality(rssi, snr)
	notes := getDiagnosticNotes(rssi, snr)

	fmt.Printf("The signal quality is %s.\n", quality)
	fmt.Printf("Diagnostic Notes: %s\n", notes)
}
