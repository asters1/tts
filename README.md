## 微软官网只有试听功能，本程序可以下载音频文件

## 如何使用?

1. 下载源代码使用

```
git clone https://github.com/asters1/tts

cd tts

go run main.go
```

2. 或者直接下载已经编译好的文件(需要配置tts.config文件)  
注:似乎交叉编译有问题，自己下载源码，自己编译吧！


3. tts.config(需要将此文件放置到tts同一目录下,并且新建一个名字为mp3的文件夹。否则会报错!!)


```
#此配置文件以#号作为注释
#
#
#
#
#语言默认为zh-CN
Language:zh-CN

#Name为发音员,默认为XiaoxiaoNeural

#发音员不止以下几个，懒得写了。自己找，网址为:
#https://azure.microsoft.com/zh-cn/services/cognitive-services/text-to-speech/#overview


#Name:XiaoxiaoNeural
#XiaoxiaoNeural     #晓晓
#XiaoyouNeural      #晓悠
Name:YunxiNeural    #云希
#-SC-YunxiNeural    #四川云希
#

#volume为音量,默认为100,区间[0，100]
volume:100

#rate为语速,默认为0,区间[-100,200]
rate:0

#pitch为音调,默认为0,区间[-50，50]
pitch:0

#path为保存音频路径，默认为当前目录下的mp3文件目录下
path:./mp3/
```







