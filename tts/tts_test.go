package tts

import (
	"fmt"
	"testing"
	"time"
)

func TestTTSWindows(t *testing.T) {
	tipTime := time.Now()
	workDir := "C:\\Users\\admin\\go\\src\\tts-go\\tts\\res\\msc_work_dir\\"
	workDir2JetPath := "..\\..\\tts\\"
	params := fmt.Sprintf("engine_type = local, voice_name = xiaoyan, text_encoding = UTF8, tts_res_path = fo|%sxiaoyan.jet;fo|%scommon.jet, sample_rate= 16000, speed = 50, volume = 50, pitch = 50, rdn = 2",
		workDir2JetPath, workDir2JetPath)

	Login(fmt.Sprintf("appid = 5e7b06e1, work_dir = %s", workDir))
	err := GetTTSWavFile(params,
		"单次合成的文本",
		"test4.wav")
	if err != nil {
		panic(err)
	}
	fmt.Println("app elapsed:", time.Since(tipTime).Milliseconds())
}

func TestTTSLinux(t *testing.T) {
	tipTime := time.Now()
	workDir := "/home/project/go/src/tts-go/tts/res/msc_work_dir"
	workDir2JetPath := "../../tts/"
	params := fmt.Sprintf("engine_type = local, voice_name = xiaoyan, text_encoding = UTF8, tts_res_path = fo|%sxiaoyan.jet;fo|%scommon.jet, sample_rate= 16000, speed = 50, volume = 50, pitch = 50, rdn = 2",
		workDir2JetPath, workDir2JetPath)
	Login(fmt.Sprintf("appid = 5e7b06e1, work_dir = %s", workDir))
	data, err := TTSData(params,
		"我们")
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data))
	fmt.Println(data)
	fmt.Println("app elapsed:", time.Since(tipTime).Milliseconds())
}
