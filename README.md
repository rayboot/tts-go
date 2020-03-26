## tts-go

文字转语音工具，基于讯飞语音离线sdk实现。

support windows and linux

### 环境依赖

* go version >= 1.11, go mod

### 使用前准备

- 讯飞语音appid
- 讯飞语音sdk linux平台或者windows平台
- linux平台 将 libmsc.so 文件移动到 LD_LIBRARY_PATH 目录中
- windows 64位 平台 将 msc_x64.dll 移动到 system32目录中


### 使用方法

```
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
```

### 使用示例

```

tts-go.exe -w C:\\Users\\admin\\go\\src\\tts-go\\res\\msc_work_dir\\ -k your_appid -t 测试一下

tts-go -w /path/to/msc_work_dir -k your_appid -t 测试一下

```

### thanks for



[auroraapi](https://github.com/auroraapi/aurora-go)

[wanliu](https://github.com/wanliu/xf)