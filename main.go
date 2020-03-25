package main

import (
	"flag"
	"fmt"
	"github.com/dyike/log"
	"os"
	"time"
	"tts-go/tts"
)

var usageStr = `
Usage: tts-go [options]

资源的目录结构为
─res
   ├─msc_work_dir
   │  └─msc
   │      └─msc.cfg
   └─tts
      └─common.jet
      └─xiaofeng.jet
      └─xiaoyan.jet
讯飞参数中work_dir 设置为msc_work_dir的绝对路径
资源引用使用  ../../tts/***.jet

注： 如果appid过期，请更换appid 以及 jet文件！！

单次合成模式选项:
	-w <work_dir>				讯飞资源工作路径，按照上述目录结构设置为msc_work_dir的绝对路径
	-k <appid>					讯飞appid
    -t <text>                	待合成的文本
    -o <file>               	音频输出路径 
	-e <engine_type>			引擎  默认local
	-n <voice_name>				发音人  xiaoyan、xiaofeng
	-c <text_encoding>			编码  默认UTF8
	-j <tts_res_path>			合成资源路径   基于上述树形结构的相对路径  linux默认为   ../../tts/  windows平台默认为 ..\\..\\tts\\
	-r <sample_rate>			合成音频采样率  合成音频采样率，支持参数，16000，8000，默认为16000
	-s <speed>					语速  合成音频对应的语速，取值范围：[0,100]，数值越大语速越快。默认值：50
	-v <volume>					音量  合成音频的音量，取值范围：[0,100]，数值越大音量越大。默认值：50
	-p <pitch>					合成语调    通过此参数，设置合成返回音频的语调，值范围：[0，100]，默认：50
	-rd <rdn>					数字发音 合成音频数字发音，支持参数， 0 数值优先, 1 完全数值,2 完全字符串，3 字符串优先，默认值：0 
	-rc <rcn>					1 的中文发音   0(默认)：表示发音为yao   1：表示发音为yi
其他:
    -h                          查看帮助 
`

func main() {
	var text string
	var file string
	var help bool

	config := tts.TTSConf{}
	flag.StringVar(&config.WorkDir, "w", "", "msc_work_dir工作目录")
	flag.StringVar(&config.AppId, "k", "", "appid")
	flag.StringVar(&text, "t", "测试文本", "需要转换的文案")
	flag.StringVar(&file, "o", "test.wav", "转换后生成的文件名")
	flag.StringVar(&config.EngineType, "e", "local", "引擎")
	flag.StringVar(&config.VoiceName, "n", "xiaoyan", "发音人")
	flag.StringVar(&config.TextEncoding, "c", "UTF8", "编码")
	flag.StringVar(&config.TtsResPath, "j", fmt.Sprintf("..%s..%stts%s", string(os.PathSeparator), string(os.PathSeparator), string(os.PathSeparator)), "合成资源路径")
	flag.IntVar(&config.SampleRate, "r", 16000, "采样率")
	flag.IntVar(&config.Speed, "s", 50, "语速")
	flag.IntVar(&config.Volume, "v", 50, "音量")
	flag.IntVar(&config.Pitch, "p", 50, "合成语调")
	flag.IntVar(&config.Rdn, "rd", 2, "数字发音")
	flag.IntVar(&config.Rcn, "rc", 0, "1 的中文发音")
	flag.BoolVar(&help, "h", false, "帮帮帮助")
	flag.Parse()

	if help {
		fmt.Printf("%s\n", usageStr)
		return
	}

	log.Debug(config.WorkDir)

	if len(config.AppId) == 0 || len(config.WorkDir) == 0 || len(file) == 0 {
		log.Debug("appid / work_dir / file 没有输入内容")
		return
	}

	log.Debug("开始合成资源")
	t := time.Now()

	if err := tts.GetTTSWavFile(config.ToTTSParams(), config.ToTTSLoginParams(), text, file); err != nil {
		panic("转换失败")
		log.Error("%v", err)
		return
	}
	log.Debug("转换成功 耗时 %d", time.Since(t).Milliseconds())
}
