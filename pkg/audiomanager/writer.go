package audiomanager

import (
	"math"
)

// IsSilentAudioDataEx determines if the audio data is silent based on enhanced RMS and optional zero-crossing rate.
// Parameters:
// - data: slice of int16 audio samples.
// - rmsThresholdDb: threshold in decibels (e.g., -50.0) below which audio is considered silent.
// - useZCR: whether to use zero-crossing rate as an additional criterion.
// - zcrThreshold: threshold for zero-crossing rate (0.0 to 1.0), only used if useZCR is true.
func IsSilentAudioDataEx(data []int16, rmsThresholdDb float64, useZCR bool, zcrThreshold float64) bool {
	if len(data) == 0 {
		return true // 빈 데이터는 무음으로 처리
	}

	var sum int64
	var sumSquares float64
	var zeroCrossings int

	// Initialize previous sample for ZCR calculation
	prevSample := float64(data[0])
	sum += int64(data[0])
	sumSquares += prevSample * prevSample

	for i := 1; i < len(data); i++ {
		sample := float64(data[i])
		sum += int64(sample)
		sumSquares += sample * sample

		// Zero-crossing detection
		if (prevSample < 0 && sample >= 0) || (prevSample > 0 && sample <= 0) {
			zeroCrossings++
		}
		prevSample = sample
	}

	// Calculate mean to remove DC offset
	mean := float64(sum) / float64(len(data))

	// Calculate RMS without DC offset
	rms := math.Sqrt(sumSquares/float64(len(data)) - mean*mean)
	if math.IsNaN(rms) || rms < 0 {
		return true // Invalid RMS is treated as silent
	}

	// Convert RMS to decibels relative to full scale (dBFS)
	// Assuming full scale is 32768 for int16
	if rms == 0 {
		return true // 완전 무음
	}
	rmsDb := 20 * math.Log10(rms/32768.0)

	// Determine silence based on RMS threshold
	isSilent := rmsDb < rmsThresholdDb

	// Optionally, include zero-crossing rate as an additional criterion
	if useZCR {
		zcr := float64(zeroCrossings) / float64(len(data))
		isSilent = isSilent && (zcr < zcrThreshold)
	}

	return isSilent
}
