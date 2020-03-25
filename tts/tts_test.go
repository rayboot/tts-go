package tts

import (
	"fmt"
	"testing"
)

func TestTTS2(t *testing.T) {
	GetTTSWavFile("engine_type = local, voice_name = xiaoyan, text_encoding = UTF8, tts_res_path = fo|/home/project/go/src/tts-go/tts/bin/msc/res/tts/xiaoyan.jet;fo|/home/project/go/src/tts-go/tts/bin/msc/res/tts/common.jet, sample_rate= 16000, speed = 50, volume = 50, pitch = 50, rdn = 2",
		"appid = 5e0450b4, work_dir = .",
		"单次合成的文本",
		"test4.wav")
}

func TestTTS3(t *testing.T) {
	data, _ := TTSData("engine_type = local, voice_name = xiaoyan, text_encoding = UTF8, tts_res_path = fo|/home/project/go/src/tts-go/tts/bin/msc/res/tts/xiaoyan.jet;fo|/home/project/go/src/tts-go/tts/bin/msc/res/tts/common.jet, sample_rate= 16000, speed = 50, volume = 50, pitch = 50, rdn = 2",
		"appid = 5e0450b4, work_dir = .",
		"单次合成的文本")

	fmt.Print(data)
}

func TestTTSWindows(t *testing.T) {
	workDir := "C:\\Users\\admin\\go\\src\\tts-go\\tts\\res\\msc_work_dir\\"
	workDir2JetPath := "..\\..\\tts\\"
	params := fmt.Sprintf("engine_type = local, voice_name = xiaoyan, text_encoding = UTF8, tts_res_path = fo|%sxiaoyan.jet;fo|%scommon.jet, sample_rate= 16000, speed = 50, volume = 50, pitch = 50, rdn = 2",
		workDir2JetPath, workDir2JetPath)
	err := GetTTSWavFile(params,
		fmt.Sprintf("appid = 5e7b06e1, work_dir = %s", workDir),
		"单次合成的文本",
		"test4.wav")
	if err != nil {
		panic(err)
	}
}

func TestTTSLinux(t *testing.T) {
	workDir := "/home/project/go/src/tts-go/tts/res/msc_work_dir"
	workDir2JetPath := "../../tts/"
	params := fmt.Sprintf("engine_type = local, voice_name = xiaoyan, text_encoding = UTF8, tts_res_path = fo|%sxiaoyan.jet;fo|%scommon.jet, sample_rate= 16000, speed = 50, volume = 50, pitch = 50, rdn = 2",
		workDir2JetPath, workDir2JetPath)
	data, err := TTSData(params,
		fmt.Sprintf("appid = 5e7b06e1, work_dir = %s", workDir),
		"我们")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data))
	fmt.Println(data)
}
