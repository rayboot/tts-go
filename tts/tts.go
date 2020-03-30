package tts

import "C"
import (
	"fmt"
	"github.com/rayboot/tts-go/audio"
	"os"
)

type TTSLoginConf struct {
	WorkDir string `json:"work_dir"`
	AppId   string `json:"appid"`
}

type TTSConf struct {
	EngineType   string `json:"engine_type"`
	VoiceName    string `json:"voice_name"`
	TextEncoding string `json:"text_encoding"`
	TtsResPath   string `json:"tts_res_path"`
	SampleRate   int    `json:"sample_rate"`
	Speed        int    `json:"speed"`
	Volume       int    `json:"volume"`
	Pitch        int    `json:"pitch"`
	Rdn          int    `json:"rdn"`
	Rcn          int    `json:"rcn"`
}

func (ttsConf TTSConf) ToTTSParams() string {
	return fmt.Sprintf("engine_type = %s, voice_name = %s, text_encoding = %s, tts_res_path = %s, sample_rate= %d, speed = %d, volume = %d, pitch = %d, rdn = %d, rcn = %d",
		ttsConf.EngineType,
		ttsConf.VoiceName,
		ttsConf.TextEncoding,
		fmt.Sprintf("fo|%s%s.jet;fo|%scommon.jet", ttsConf.TtsResPath, ttsConf.VoiceName, ttsConf.TtsResPath),
		ttsConf.SampleRate,
		ttsConf.Speed,
		ttsConf.Volume,
		ttsConf.Pitch,
		ttsConf.Rdn,
		ttsConf.Rcn)
}

func (ttsLoginConf TTSLoginConf) ToTTSLoginParams() string {
	return fmt.Sprintf("appid = %s, work_dir = %s",
		ttsLoginConf.AppId,
		ttsLoginConf.WorkDir)
}

// 生成wav文件
func GetTTSWavFile(ttsParmas, speedTxt, desPath string) error {
	audioData, err := TTSData(ttsParmas, speedTxt)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(desPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(audioData)
	if err != nil {
		return err
	}
	return nil
}

func Login(loginParams string) error {
	return MSPLogin(loginParams)
}

func Logout() error {
	return MSPLogout()
}

// 获取语音的二进制数据
func TTSData(ttsParmas, speedTxt string) ([]byte, error) {
	if speedTxt == "" {
		return audio.NewWAV().Data(), nil
	}
	/* 开始合成 */
	sessionID, err := QTTSSessionBegin(ttsParmas)
	if err != nil {
		return nil, err
	}
	defer QTTSSessionEnd(sessionID, "Normal")
	if err = QTTSTextPut(sessionID, speedTxt, ""); err != nil {
		return nil, err
	}
	audioData := audio.NewWAV()
	for true {
		data, synthStatus, err := QTTSAudioGet(sessionID)
		if err != nil {
			return nil, err
		}
		if data != nil {
			audioData.AddAudioData(data)
		}
		if synthStatus == 2 {
			break
		}
	}
	return audioData.Data(), nil
}
