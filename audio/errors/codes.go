package errors

// ErrorCode is a code for an error
type ErrorCode string

// Define the various error codes possible
const (
	SpeechNilAudio                 = "SpeechNilAudio"
	WAVCorruptFile                 = "WAVCorruptFile"
	AudioFileOutputStreamNotOpened = "AudioFileOutputStreamNotOpened"
	AudioFileNotWritableStream     = "AudioFileNotWritableStream"
)

// errorMessages converts an error code to its corresponding message
var errorMessages = map[ErrorCode]string{
	SpeechNilAudio:                 "The audio file was nil. In order to convert a Speech object to Text, it must have a valid audio file. Usually, this means you created a Speech object that wasn't created using one of the Listen methods.",
	WAVCorruptFile:                 "The WAV file was corrupted and did not have a correctly formatted RIFF header. Check the file to make sure it was not corrupted or incomplete.",
	AudioFileOutputStreamNotOpened: "PortAudio encountered an error in opening the audio stream, which is usually due to an error in connecting to the input and/or output device.",
	AudioFileNotWritableStream:     "The data could not be written into the stream. You may have attempted to write to a callback stream, tried to write to an input-only stream, created a buffer with incorrect parameters, or did not open the stream at all.",
}
