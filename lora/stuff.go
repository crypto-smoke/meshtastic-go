// Package loraradio provides functionality to determine values for a LoRa radio link,
// including maximum data rate and link budget while accounting for receiver sensitivity.
package lora

import (
	"fmt"
	"math"
)

// TODO: needs cleanup. lots of gpt generated code

// MaxDataRate calculates the maximum data rate for a LoRa radio link based on the provided parameters.
func MaxDataRate(bandwidth, spreadingFactor, codeRate float64) float64 {
	// Assuming the LoRa modulation's data rate equation:
	// Data Rate = BW / (2^SF) * CR
	// where SF is the spreading factor, BW is the bandwidth in Hz, and CR is the code rate.
	return (bandwidth / math.Pow(2, spreadingFactor)) * codeRate * spreadingFactor
}

// Simplified link budget calculation per the Semtech calculator
func LinkBudget(rxSensitivity, transmitPower float64) float64 {
	return transmitPower - rxSensitivity
}

// SensitivityParams holds the parameters used for sensitivity calculations.
type SensitivityParams struct {
	Bandwidth       float64 // in Hz
	ImplementationL float64 // Implementation loss in dB, typically 1-3 dB
	SpreadingFactor int     // LoRa Spreading Factor
}

// SpreadingFactorData holds the data for each spreading factor.
type SpreadingFactorData struct {
	SF             int
	ChipsPerSymbol int
	DemodulatorSNR float64
}

var SpreadingFactors = []SpreadingFactorData{
	{SF: 5, ChipsPerSymbol: 32, DemodulatorSNR: -2.5},
	{SF: 6, ChipsPerSymbol: 64, DemodulatorSNR: -5},
	{SF: 7, ChipsPerSymbol: 128, DemodulatorSNR: -7.5},
	{SF: 8, ChipsPerSymbol: 256, DemodulatorSNR: -10},
	{SF: 9, ChipsPerSymbol: 512, DemodulatorSNR: -12.5},
	{SF: 10, ChipsPerSymbol: 1024, DemodulatorSNR: -15},
	{SF: 11, ChipsPerSymbol: 2048, DemodulatorSNR: -17.5},
	{SF: 12, ChipsPerSymbol: 4096, DemodulatorSNR: -20},
}

// CalculateSensitivity calculates the LoRa receiver sensitivity based on the provided parameters.
func CalculateSensitivity(params SensitivityParams, spreadingFactors []SpreadingFactorData) (float64, error) {
	if params.Bandwidth <= 0 {
		return 0, fmt.Errorf("bandwidth must be greater than 0")
	}

	// Find the SNR for the given spreading factor from the provided data.
	var snr float64
	found := false
	for _, sfData := range spreadingFactors {
		if sfData.SF == params.SpreadingFactor {
			snr = sfData.DemodulatorSNR
			found = true
			break
		}
	}
	if !found {
		return 0, fmt.Errorf("spreading factor data not found for SF=%d", params.SpreadingFactor)
	}

	// Thermal noise in dBm for the given bandwidth at room temperature (290K).
	thermalNoise := -174.0 // dBm/Hz

	// Noise figure in dBm for the given bandwidth.
	noiseFigure := thermalNoise + 10*math.Log10(params.Bandwidth)

	// Receiver sensitivity calculation in dBm.
	sensitivity := noiseFigure + snr + params.ImplementationL

	return sensitivity, nil
}
func SNR(spreadingFactor int) float64 {
	result := (float64(spreadingFactor) - 4) * -2.5
	return result * 100 / 100
}
