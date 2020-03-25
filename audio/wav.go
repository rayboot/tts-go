package audio

import (
	"encoding/binary"
	"io"
	"io/ioutil"
	"math"
	"tts-go/audio/errors"
)

// WAV-related constants.
const (
	// DefaultNumChannels is 1 (mono audio)
	DefaultNumChannels uint16 = 1
	// DefaultSampleRate is 16KHz
	DefaultSampleRate uint32 = 16000
	// DefaultAudioFormat is 1 (raw, uncompressed PCM waveforms)
	DefaultAudioFormat uint16 = 1
	// DefaultBitsPerSample is 16 (2 bytes per sample).
	DefaultBitsPerSample uint16 = 16
)

// WAV represents a PCM audio file in the WAV container format. It keeps
// a high-level description of the parameters of the file, along with the
// raw audio bytes, until it needs to be written to a file, stream, or array.
// It is based on the WAV formatting as specified: http://soundfile.sapp.org/doc/WaveFormat/
type WAV struct {
	// NumChannels is the number of channels the WAV file has. 1 = mono,
	// 2 = stereo, etc. This affects the block align and also number of
	// bytes per sample: (BitsPerSample / 8) * NumChannels.
	NumChannels uint16
	// SampleRate is the number of samples taken per second.
	SampleRate uint32
	// AudioFormat is the type of Audio that is encoded in the WAV file.
	// In most scenarios, this will be 1 (1 = raw, uncompressed PCM audio),
	// since WAV doesn't support compression.
	AudioFormat uint16
	// BitsPerSample is the width of each sample. 16 bits means each sample
	// is two bytes.
	BitsPerSample uint16

	// audioData is the raw audio data stored in the WAV file.
	audioData []byte
}

// WAVParams are a set of parameters used to create a WAV file. Its fields
// correspond directly to the WAV file.
type WAVParams struct {
	NumChannels   uint16
	SampleRate    uint32
	BitsPerSample uint16
	AudioData     []byte
}

// NewWAV returns a new, empty WAV file using the default parameters.
func NewWAV() *WAV {
	// create a new default WAV
	return &WAV{
		NumChannels:   DefaultNumChannels,
		SampleRate:    DefaultSampleRate,
		AudioFormat:   DefaultAudioFormat,
		BitsPerSample: DefaultBitsPerSample,
		audioData:     make([]byte, 0),
	}
}

// NewWAVFromParams returns a new WAV file from the passed in parameters
// If any of the parameters are 0, it will be given the default value.
func NewWAVFromParams(params *WAVParams) *WAV {
	if params == nil {
		return NewWAV()
	}
	if params.NumChannels == 0 {
		params.NumChannels = DefaultNumChannels
	}
	if params.SampleRate == 0 {
		params.SampleRate = DefaultSampleRate
	}
	if params.BitsPerSample == 0 {
		params.BitsPerSample = DefaultBitsPerSample
	}
	if params.AudioData == nil || len(params.AudioData) == 0 {
		params.AudioData = make([]byte, 0)
	}
	return &WAV{
		NumChannels:   params.NumChannels,
		SampleRate:    params.SampleRate,
		AudioFormat:   DefaultAudioFormat,
		BitsPerSample: params.BitsPerSample,
		audioData:     params.AudioData,
	}
}

// NewWAVFromData creates a WAV format struct from the given data buffer
// The buffer is broken up into its respective information and that
func NewWAVFromData(data []byte) (*WAV, error) {
	// find the end of ChunkID denoted by RIFF
	// This marks the beginning of the WAV file
	i := 4
	for i < len(data) && data[i-4] != 'R' || data[i-3] != 'I' || data[i-2] != 'F' || data[i-1] != 'F' {
		i++
	}

	dataLen := len(data) - i
	if dataLen <= 0 {
		return nil, errors.NewFromErrorCodeInfo(errors.WAVCorruptFile, "The letters `RIFF` should exist from bytes 0 to 3 in big endian form from the start of the header to indicate that it is a RIFF header.")
	}

	// hOff is the header offset. Even though the header length is actually
	// 44, we find where the data begins by looking for the letters
	// "RIFF" which is from bytes 0 to 3. The variable i at this point
	// pointing to right past the "RIFF" letters
	hOff := i - 4

	if (len(data) - hOff - 44) < 0 {
		return nil, errors.NewFromErrorCode(errors.WAVCorruptFile)
	}

	// Verifies that "WAVE" letters exist in big endian form
	if data[hOff+8] != 'W' || data[hOff+9] != 'A' || data[hOff+10] != 'V' || data[hOff+11] != 'E' {
		return nil, errors.NewFromErrorCodeInfo(errors.WAVCorruptFile, "The letters `WAVE` should exist from bytes 8 to 11 in big endian form from the start of the header to indicate that it is a WAVE format file.")
	}

	// Verifies that "fmt " letters exist in big endian form
	if data[hOff+12] != 'f' || data[hOff+13] != 'm' || data[hOff+14] != 't' || data[hOff+15] != ' ' {
		return nil, errors.NewFromErrorCodeInfo(errors.WAVCorruptFile, "The letters `fmt ` should exist from bytes 12 to 15 in big endian form from the start of the header to indicate the subchunk 1 ID")
	}

	// Verifies that the "data" letters exist in big endian form
	if data[hOff+36] != 'd' || data[hOff+37] != 'a' || data[hOff+38] != 't' || data[hOff+39] != 'a' {
		return nil, errors.NewFromErrorCodeInfo(errors.WAVCorruptFile, "The letters `data` should exist from bytes 36 to 39 in big endian form from the start of the header to indicate the subchunk 2 ID.")
	}

	numChannels := binary.LittleEndian.Uint16(data[hOff+22 : hOff+24])
	sampleRate := binary.LittleEndian.Uint32(data[hOff+24 : hOff+28])
	bitsPerSample := binary.LittleEndian.Uint16(data[hOff+34 : hOff+36])

	// The actual sound data begins at byte 44 from the beginning of the header
	audioData := data[hOff+44:]

	return &WAV{
		NumChannels:   numChannels,
		SampleRate:    sampleRate,
		AudioFormat:   DefaultAudioFormat,
		BitsPerSample: bitsPerSample,
		audioData:     audioData,
	}, nil
}

// NewWAVFromReader takes in a reader and creates a new WAV file.
func NewWAVFromReader(reader io.Reader) (*WAV, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return NewWAVFromData(b)
}

// TrimSilent is called on a WAV struct to trim the silent portions from the
// ends of the file while leaving a certain amount of padding. The padding input is
// specified in seconds. The threshold input is a decimal (between 0 and 1) and is
// relative to the maximum amplitude of the waveform
func (w *WAV) TrimSilent(threshold float64, padding float64) {
	// sample size in bytes
	sampleSize := int(w.NumChannels * w.BitsPerSample / 8)
	// number of bytes to examine in each step
	step := 1024

	// get max amplitude
	maxPossibleAmp := math.Exp2(float64(w.BitsPerSample)) / 2.0
	// silenceThresh is a percentage of the maximum wave height
	silenceThresh := threshold * maxPossibleAmp

	// Trimming the beginning
	N1 := 0
	for N1 < len(w.audioData) {
		sampleRMS := rms(sampleSize, w.audioData[N1:N1+(sampleSize*step)])
		if sampleRMS > silenceThresh {
			break
		}
		N1 += sampleSize * step
	}

	// Trimming the end
	N2 := len(w.audioData)
	for N2 >= 0 {
		sampleRMS := rms(sampleSize, w.audioData[N2-(sampleSize*step):N2])
		if sampleRMS > silenceThresh {
			break
		}
		N2 -= sampleSize * step
	}

	paddingSamples := int(padding * float64(w.SampleRate) * float64(sampleSize))
	w.audioData = w.audioData[N1-paddingSamples : N2+paddingSamples]
}

// AddAudioData adds the passed-in audio bytes to the WAV struct
func (w *WAV) AddAudioData(d []byte) {
	// add audio data to existing data
	if d != nil && len(d) > 0 {
		w.audioData = append(w.audioData, d...)
	}
}

// AudioData returns the raw audio data
func (w *WAV) AudioData() []byte {
	return w.audioData
}

// Data creates the header and data based on the WAV struct and returns
// a fully formatted WAV file
func (w *WAV) Data() []byte {
	// find first data index
	dataLen := len(w.audioData)
	headerLen := 44
	chunkSize := (dataLen + headerLen - 8)

	// first create the header, then append the rest of the file
	wav := make([]byte, headerLen)

	// RIFF header
	wav[0] = 'R'
	wav[1] = 'I'
	wav[2] = 'F'
	wav[3] = 'F'

	// chunk size
	binary.LittleEndian.PutUint32(wav[4:8], uint32(chunkSize))

	// Format (WAVE)
	wav[8] = 'W'
	wav[9] = 'A'
	wav[10] = 'V'
	wav[11] = 'E'
	// Metadata subchunk ID ("fmt ")
	wav[12] = 'f'
	wav[13] = 'm'
	wav[14] = 't'
	wav[15] = ' '

	// Metadata subchunk size (16)
	binary.LittleEndian.PutUint32(wav[16:20], 16)
	// Audio format (PCM = 1)
	binary.LittleEndian.PutUint16(wav[20:22], w.AudioFormat)
	// Num Channels (Mono = 1)
	binary.LittleEndian.PutUint16(wav[22:24], w.NumChannels)
	// Sample Rate (16000 Hz)
	binary.LittleEndian.PutUint32(wav[24:28], w.SampleRate)

	// Byte Rate = SampleRate * NumChannels * BitsPerSample/8 = 32000
	byteRate := w.SampleRate * uint32(w.NumChannels) * uint32(w.BitsPerSample) / 8
	binary.LittleEndian.PutUint32(wav[28:32], byteRate)

	// Block Align = NumChannels * BitsPerSample/8
	blockAlign := w.NumChannels * w.BitsPerSample / 8
	binary.LittleEndian.PutUint16(wav[32:34], blockAlign)

	// Bits per sample = 16
	binary.LittleEndian.PutUint16(wav[34:36], w.BitsPerSample)

	// Data subchunk ID ("data")
	wav[36] = 'd'
	wav[37] = 'a'
	wav[38] = 't'
	wav[39] = 'a'

	// Data length
	binary.LittleEndian.PutUint32(wav[40:44], uint32(dataLen))

	return append(wav, w.audioData[0:]...)
}
