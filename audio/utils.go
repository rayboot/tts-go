package audio

import (
	"encoding/binary"
	"math"
)

// rms calculates the root-mean-square of a sequence of audio data. For
// now, it assumes that the data is in 16-bit mono samples. Thus, the passed
// value for `sampleSize` MUST be 2. This will change once we figure out
// how to read a variable size of data during runtime.
func rms(sampleSize int, audioData []byte) float64 {
	sum := 0.0
	for i := 0; i < len(audioData)-1; i += sampleSize {
		// had to hard code to Uint16 or else it tries to 8 bytes for Uint64
		// note that sampleSize MUST be 2 for this to work
		val := binary.LittleEndian.Uint16(audioData[i:(i + sampleSize)])
		sum += float64(val * val)
	}

	return math.Sqrt(sum / (float64(len(audioData) / sampleSize)))
}

// isSilent determines whether an audio slice is silent or not
func isSilent(audio []int16) bool {
	var max int16 = audio[0]
	for _, value := range audio {
		if value > max {
			max = value
		}
	}
	return max < SilentThresh
}

// recordResponse is a response emitted on the channel returns from the
// record function.
type recordResponse struct {
	// Data contains the raw data converted from the n-bit sample that
	// is read from PortAudio
	Data []byte
	// Samples contains the raw samples read from PortAudio (do not use --
	// this is overwritten in every message sent on the channel. Its use
	// is purely internal)
	Samples []int16
	// Error if an error occurred
	Error error
}
