package lora

import "testing"

// TestGetSignalQuality tests the GetSignalQuality function with different RSSI and SNR values.
func TestGetSignalQuality(t *testing.T) {
	tests := []struct {
		rssi     float64
		snr      float64
		expected SignalQuality
	}{
		{-115, -7, SignalQualityGood},
		{-120, -10, SignalQualityFair},
		{-130, -20, SignalQualityBad},
	}

	for _, test := range tests {
		actual := GetSignalQuality(test.rssi, test.snr)
		if actual != test.expected {
			t.Errorf("GetSignalQuality(%v, %v) = %v; want %v", test.rssi, test.snr, actual, test.expected)
		}
	}
}

// TestGetDiagnosticNotes tests the GetDiagnosticNotes function with different RSSI and SNR values.
func TestGetDiagnosticNotes(t *testing.T) {
	tests := []struct {
		rssi     float64
		snr      float64
		expected string
	}{
		{-115, -7, "RF level is optimal to get a good reception reliability."},
		{-126, -15, "RF level is not optimal but must be sufficient. Try to improve your device position if possible. You will have to monitor the stability of the RF level."},
		{-130, -20, "NOISY environment. Try to put device out of electromagnetic sources."},
	}

	for _, test := range tests {
		actual := GetDiagnosticNotes(test.rssi, test.snr)
		if actual != test.expected {
			t.Errorf("GetDiagnosticNotes(%v, %v) = %v; want %v", test.rssi, test.snr, actual, test.expected)
		}
	}
}
