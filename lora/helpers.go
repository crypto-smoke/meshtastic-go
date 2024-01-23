// Package lora provides utilities to assess the signal quality of LoRa (Long Range) communication
// based on RSSI (Received Signal Strength Indicator) and SNR (Signal-to-Noise Ratio) values.
package lora

import "fmt"

// SignalQuality defines a type for representing the quality of a signal.
type SignalQuality string

const (
	// SignalQualityGood indicates a good signal quality.
	SignalQualityGood SignalQuality = "GOOD"
	// SignalQualityFair indicates a fair signal quality.
	SignalQualityFair SignalQuality = "FAIR"
	// SignalQualityBad indicates a bad signal quality.
	SignalQualityBad SignalQuality = "BAD"
)

// GetSignalQuality determines the signal quality based on RSSI and SNR.
// A GOOD signal is determined by SNR >= -7 and RSSI >= -115.
// A FAIR signal is determined by SNR >= -15 and RSSI >= -126.
// Any signal that does not meet the GOOD or FAIR criteria is considered BAD.
func GetSignalQuality(rssi, snr float64) SignalQuality {
	// Define the boundaries for GOOD signal quality
	if snr >= -7 && rssi >= -115 {
		return SignalQualityGood
	}
	// Define the boundaries for FAIR signal quality
	if snr >= -15 && rssi >= -126 {
		return SignalQualityFair
	}
	// If none of the above conditions are met, signal is BAD
	return SignalQualityBad
}

// GetDiagnosticNotes provides recommendations based on RSSI and SNR values.
// It returns different notes depending on the quality of the RF level.
func GetDiagnosticNotes(rssi, snr float64) string {
	if rssi >= -115 && snr >= -7 {
		return "RF level is optimal to get a good reception reliability."
	} else if rssi >= -126 && snr >= -15 {
		return "RF level is not optimal but must be sufficient. Try to improve your device position if possible. You will have to monitor the stability of the RF level."
	} else {
		return "NOISY environment. Try to put device out of electromagnetic sources."
	}
}

// ExampleGetSignalQuality demonstrates the usage of the GetSignalQuality function.
func ExampleGetSignalQuality() {
	rssi := -120.0 // RSSI value
	snr := -10.0   // SNR value

	quality := GetSignalQuality(rssi, snr)
	fmt.Printf("The signal quality is %s.\n", quality)
	// Output: The signal quality is FAIR.
}

// ExampleGetDiagnosticNotes demonstrates the usage of the GetDiagnosticNotes function.
func ExampleGetDiagnosticNotes() {
	rssi := -130.0 // RSSI value
	snr := -20.0   // SNR value

	notes := GetDiagnosticNotes(rssi, snr)
	fmt.Printf("Diagnostic Notes: %s\n", notes)
	// Output: Diagnostic Notes: NOISY environment. Try to put device out of electromagnetic sources.
}
