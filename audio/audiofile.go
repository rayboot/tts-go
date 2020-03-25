// Package audio contains structs and functions to allow operating on audio data
package audio

import (
	"io"
	"io/ioutil"
	"os"
)

const (
	BufSize      = 1 << 10
	SilentThresh = 1 << 10
	SampleRate   = 16000
	NumChannels  = 1
)

// File represents an audio file. It wraps the raw WAV data and allows you
// to operate with it using high-level operations, such as padding, trimming,
// playback, writing to file, recording, etc.
type File struct {
	AudioData *WAV

	playing    bool
	shouldStop bool
}

// WriteToFile writes the audio data to a file.
func (f *File) WriteToFile(filename string) error {
	return ioutil.WriteFile(filename, f.AudioData.Data(), 0644)
}

// Pad adds silence to both the beginning and end of the audio data. Silence
// is specified in seconds.
func (f *File) Pad(seconds float64) {
	// calculate number of bytes needed to pad given amount of seconds
	bytes := float64(f.AudioData.NumChannels*f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * seconds
	padding := make([]byte, int(bytes))

	// copy the WAV parameters, set initial data to the left padding
	newWav := NewWAVFromParams(&WAVParams{
		NumChannels:   f.AudioData.NumChannels,
		SampleRate:    f.AudioData.SampleRate,
		BitsPerSample: f.AudioData.BitsPerSample,
		AudioData:     padding,
	})

	// add the original data and ther right padding
	newWav.AddAudioData(f.AudioData.AudioData())
	newWav.AddAudioData(padding)

	// set the audio data to the new wav
	f.AudioData = newWav
}

// PadLeft adds silence to the beginning of the audio data.
func (f *File) PadLeft(seconds float64) {
	// calculate number of bytes needed to pad given amount of seconds
	bytes := float64(f.AudioData.NumChannels*f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * seconds
	padding := make([]byte, int(bytes))

	// copy the WAV parameters, set initial data to the left padding
	newWav := NewWAVFromParams(&WAVParams{
		NumChannels:   f.AudioData.NumChannels,
		SampleRate:    f.AudioData.SampleRate,
		BitsPerSample: f.AudioData.BitsPerSample,
		AudioData:     padding,
	})

	// add the original data
	newWav.AddAudioData(f.AudioData.AudioData())

	// set the audio data to the new wav
	f.AudioData = newWav
}

// PadRight adds silence to the end of the audio data.
func (f *File) PadRight(seconds float64) {
	// calculate number of bytes needed to pad given amount of seconds
	bytes := float64(f.AudioData.NumChannels*f.AudioData.BitsPerSample/8) * float64(f.AudioData.SampleRate) * seconds
	padding := make([]byte, int(bytes))

	// add the padding to the right
	f.AudioData.AddAudioData(padding)
}

// TrimSilence trims silence from both ends of the audio data.
func (f *File) TrimSilence() {
	// TODO: calibrate constants
	f.AudioData.TrimSilent(0.03, 0.25)
}

// Stop the audio playback immediately.
func (f *File) Stop() {
	if f.playing {
		f.shouldStop = true
	}
}

// NewFileFromBytes creates a new audio.File from WAV data
func NewFileFromBytes(b []byte) (*File, error) {
	wav, err := NewWAVFromData(b)
	if err != nil {
		return nil, err
	}
	return &File{wav, false, false}, err
}

// NewFileFromReader creates a new audio.File from an io.Reader
func NewFileFromReader(r io.Reader) (*File, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return NewFileFromBytes(data)
}

// NewFileFromFile creates a new audio.File from an os.File
func NewFileFromFile(f *os.File) (*File, error) {
	return NewFileFromReader(f)
}

// NewFileFromFileName creates a new audio.File from the given filename
func NewFileFromFileName(f string) (*File, error) {
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return NewFileFromBytes(data)
}

// WAVData returns the wav data (header + audio data) contained in the
// audio file
func (f *File) WAVData() []byte {
	return f.AudioData.Data()
}
