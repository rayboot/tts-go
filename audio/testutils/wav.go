package testutils

import (
	"encoding/binary"
)

func CreateEmptyWAVFile() []byte {
	emptyWAVFile := make([]byte, 44)
	// "RIFF" marker
	binary.BigEndian.PutUint32(emptyWAVFile[0:4], 0x52494646)
	// file size
	binary.LittleEndian.PutUint32(emptyWAVFile[4:8], 0)
	// "WAVE" type
	binary.BigEndian.PutUint32(emptyWAVFile[8:12], 0x57415645)
	// "fmt" section
	binary.BigEndian.PutUint32(emptyWAVFile[12:16], 0x666d7420)
	// length of fmt section
	binary.LittleEndian.PutUint32(emptyWAVFile[16:20], 16)
	// audio format
	binary.LittleEndian.PutUint16(emptyWAVFile[20:22], 1)
	// num channels
	binary.LittleEndian.PutUint16(emptyWAVFile[22:24], 1)
	// sample rate
	binary.LittleEndian.PutUint32(emptyWAVFile[24:28], 44100)
	// byte rate ((Sample Rate * Bit Size * Channels) / 8)
	binary.LittleEndian.PutUint32(emptyWAVFile[28:32], 44100)
	// block align ((bit size * channels) / 8)
	binary.LittleEndian.PutUint16(emptyWAVFile[32:34], 1)
	// bits per sample
	binary.LittleEndian.PutUint16(emptyWAVFile[34:36], 16)
	// "data" marker
	binary.BigEndian.PutUint32(emptyWAVFile[36:40], 0x64617461)
	// data size
	binary.LittleEndian.PutUint32(emptyWAVFile[40:44], 0)

	return emptyWAVFile
}
