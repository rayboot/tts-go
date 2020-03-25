package tts

import "C"
import (
	"fmt"
	"os"
	"tts-go/audio"
)

type TTSConf struct {
	WorkDir      string `json:"work_dir"`
	AppId        string `json:"appid"`
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

func (ttsConf TTSConf) ToTTSLoginParams() string {
	return fmt.Sprintf("appid = %s, work_dir = %s",
		ttsConf.AppId,
		ttsConf.WorkDir)
}

// 生成wav文件
func GetTTSWavFile(ttsParmas, loginParams, speedTxt, desPath string) error {
	audioData, err := TTSData(ttsParmas, loginParams, speedTxt)
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

// 获取语音的二进制数据
func TTSData(ttsParmas, loginParams, speedTxt string) ([]byte, error) {
	if speedTxt == "" {
		return audio.NewWAV().Data(), nil
	}

	if err := MSPLogin(loginParams); err != nil {
		return nil, err
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
	fmt.Print("正在合成 ...\n")
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
