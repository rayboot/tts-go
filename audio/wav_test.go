package audio_test

import (
	"bytes"
	"encoding/binary"
	"github.com/rayboot/tts-go/audio"
	"github.com/rayboot/tts-go/audio/errors"
	"github.com/rayboot/tts-go/audio/testutils"
	"github.com/stretchr/testify/require"
	"testing"
)

var apiErrorType *errors.APIError

// This test checks that if the NewWAV function is called, a correctly
// formatted struct is created with the default parameters
// as specified by the constants in the aduio package
func TestNewWAV(t *testing.T) {
	wav := audio.NewWAV()
	require.Equal(t, audio.DefaultNumChannels, wav.NumChannels)
	require.Equal(t, audio.DefaultSampleRate, wav.SampleRate)
	require.Equal(t, audio.DefaultAudioFormat, wav.AudioFormat)
	require.Equal(t, audio.DefaultBitsPerSample, wav.BitsPerSample)

	audioData := wav.AudioData()
	require.Equal(t, 0, len(audioData))
}

// If the WAVParams are specified, then the information should be
// correctly written into the WAV struct
func TestNewWAVFromParamsCustom(t *testing.T) {
	emptyAudio := make([]byte, 0)
	wav := audio.NewWAVFromParams(&audio.WAVParams{1, 10000, 8, emptyAudio})
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(10000), wav.SampleRate)
	require.Equal(t, uint16(8), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

func TestNewWAVFromParamsNotSpecified(t *testing.T) {
	// emptyAudio := make([]byte, 0)
	wav := audio.NewWAVFromParams(nil)
	require.Equal(t, audio.DefaultNumChannels, wav.NumChannels)
	require.Equal(t, audio.DefaultSampleRate, wav.SampleRate)
	require.Equal(t, audio.DefaultBitsPerSample, wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

// If some of the WAVParams are specified to be 0, then the default
// parameters specified in the wav.go file should be given
func TestNewWAVFromParamsWithZero(t *testing.T) {
	emptyAudio := make([]byte, 0)
	wav := audio.NewWAVFromParams(&audio.WAVParams{0, 0, 0, emptyAudio})
	require.Equal(t, audio.DefaultNumChannels, wav.NumChannels)
	require.Equal(t, audio.DefaultSampleRate, wav.SampleRate)
	require.Equal(t, audio.DefaultBitsPerSample, wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

// When given bytes of data, the information is properly parsed into a new
// WAV struct
func TestNewWAVFromData(t *testing.T) {
	emptyWAVFile := testutils.CreateEmptyWAVFile()

	wav, err := audio.NewWAVFromData(emptyWAVFile)
	require.Nil(t, err)
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(44100), wav.SampleRate)
	require.Equal(t, uint16(16), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

// When given bytes of data where the first couple of bytes is not part of the
// data file, then those parts should be regarded in reader the WAV header
func TestNewWAVFromDataWithNoise(t *testing.T) {
	buff := []byte{0x12, 0x34, 0x56, 0x78}
	emptyWAVFile := testutils.CreateEmptyWAVFile()

	wav, err := audio.NewWAVFromData(append(buff, emptyWAVFile...))
	require.Nil(t, err)
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(44100), wav.SampleRate)
	require.Equal(t, uint16(16), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

// A new reader with bytes of data should be properly parsed into a new
// WAV struct
func TestNewWAVFromReader(t *testing.T) {
	emptyWAVFile := testutils.CreateEmptyWAVFile()
	r := bytes.NewReader(emptyWAVFile)

	wav, err := audio.NewWAVFromReader(r)
	require.Nil(t, err)
	require.Equal(t, uint16(1), wav.NumChannels)
	require.Equal(t, uint32(44100), wav.SampleRate)
	require.Equal(t, uint16(16), wav.BitsPerSample)
	require.Equal(t, 0, len(wav.AudioData()))
}

// If AddAudioData function is passed a set of bytes, the wav
// structure should now contain the new data
func TestAddAudioData(t *testing.T) {
	emptyWAVFile := testutils.CreateEmptyWAVFile()

	audioData := []byte{0xff, 0xfd, 0x00, 0x00}

	wav, err := audio.NewWAVFromData(emptyWAVFile)
	wav.AddAudioData(audioData)

	require.Nil(t, err)
	require.Equal(t, 4, len(wav.AudioData()))
}

// If AddAudioData function is called on an empty set of bytes,
// the wav structure should remain the same
func TestAddAudioDataEmpty(t *testing.T) {
	emptyWAVFile := testutils.CreateEmptyWAVFile()

	audioData := make([]byte, 0)

	wav, err := audio.NewWAVFromData(emptyWAVFile)
	wav.AddAudioData(audioData)

	require.Nil(t, err)
	require.Equal(t, 0, len(wav.AudioData()))
}

func TestData(t *testing.T) {
	emptyWAVFile := testutils.CreateEmptyWAVFile()

	audioData := make([]byte, 4)
	binary.LittleEndian.PutUint32(audioData, 0x0000fdff)

	wav, err := audio.NewWAVFromData(emptyWAVFile)
	wav.AddAudioData(audioData)
	dataBytes := wav.Data()

	require.Nil(t, err)
	require.Equal(t, []byte{0xff, 0xfd, 0x0, 0x0}, dataBytes[44:48])
}
