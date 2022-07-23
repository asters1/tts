package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/asters1/tools"
	"github.com/gorilla/websocket"
)

func GetToken() string {
	res := tools.RequestClient("https://azure.microsoft.com/en-gb/services/cognitive-services/text-to-speech/", "get", "", "")
	token := tools.Re(res, `token: \"(.*?)\"`)[1]
	return token

}
func GetISOTime() string {
	T := time.Now().String()
	return T[:23][:10] + "T" + T[:23][11:] + "Z"

}
func GetLocalTime() string {
	BJ, _ := time.LoadLocation("Asia/Shanghai")
	T := time.Now().In(BJ).String()
	return T[:23][:10] + "_" + T[:19][11:]

}

func main() {
	config, _ := tools.GetConfig("./tts.config")
	fmt.Println("获取uuid...")
	uuid := tools.GetUUID()
	fmt.Println("获取token...")

	token := GetToken()
	WssUrl := `wss://eastus.tts.speech.microsoft.com/cognitiveservices/websocket/v1?Authorization=` + token + `&X-ConnectionId=` + uuid
	dl := websocket.Dialer{
		EnableCompression: true,
	}

	fmt.Println("创建websocket连接...")
	conn, _, err := dl.Dial(WssUrl, tools.GetHeader(
		`Accept-Encoding:gzip
		User-Agent:Mozilla/5.0 (Linux; Android 7.1.2; M2012K11AC Build/N6F26Q; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/81.0.4044.117 Mobile Safari/537.36
		Origin:https://azure.microsoft.com`,
	))
	if err != nil {
		fmt.Println("websocket连接失败")
	}
	defer conn.Close()

	m1 := "Path: speech.config\r\nX-RequestId: " + uuid + "\r\nX-Timestamp: " + GetISOTime() + "\r\nContent-Type: application/json\r\n\r\n{\"context\":{\"system\":{\"name\":\"SpeechSDK\",\"version\":\"1.19.0\",\"build\":\"JavaScript\",\"lang\":\"JavaScript\",\"os\":{\"platform\":\"Browser/Linux x86_64\",\"name\":\"Mozilla/5.0 (X11; Linux x86_64; rv:78.0) Gecko/20100101 Firefox/78.0\",\"version\":\"5.0 (X11)\"}}}}"
	m2 := "Path: synthesis.context\r\nX-RequestId: " + uuid + "\r\nX-Timestamp: " + GetISOTime() + "\r\nContent-Type: application/json\r\n\r\n{\"synthesis\":{\"audio\":{\"metadataOptions\":{\"sentenceBoundaryEnabled\":false,\"wordBoundaryEnabled\":false},\"outputFormat\":\"audio-24khz-160kbitrate-mono-mp3\"}}}"
	conn.WriteMessage(websocket.TextMessage, []byte(m1))
	conn.WriteMessage(websocket.TextMessage, []byte(m2))
	for {
		fmt.Println("请输入文本内容:(退出请输入q,然后回车)")
		readers := bufio.NewReader(os.Stdin)

		text, _, _ := readers.ReadLine()
		if string(text) == "q" {
			break
		}
		SSML := `<speak xmlns="http://www.w3.org/2001/10/synthesis" xmlns:mstts="http://www.w3.org/2001/mstts" xmlns:emo="http://www.w3.org/2009/10/emotionml" version="1.0" xml:lang="en-US">
        <voice name="` + config["Language"] + `-` + config["Name"] + `">
            <mstts:express-as style="general" >
                <prosody rate="` + config["rate"] + `%" volume="` + config["volume"] + `" pitch="` + config["pitch"] + `%">` + string(text) + `</prosody>
            </mstts:express-as>
        </voice>
    </speak>`
		m3 := "Path: ssml\r\nX-RequestId: " + tools.GetUUID() + "\r\nX-Timestamp: " + GetISOTime() + "\r\nContent-Type: application/ssml+xml\r\n\r\n" + SSML
		conn.WriteMessage(websocket.TextMessage, []byte(m3))

		var Adata []byte
		fmt.Println("正在下载文件...")
		for {
			Num, message, err := conn.ReadMessage()
			time.Sleep(time.Second)
			if err != nil {
				fmt.Println(err)
				break
			}
			if Num == 2 {
				index := strings.Index(string(message), "Path:audio")

				data := []byte(string(message)[index+12:])
				Adata = append(Adata, data...)
			} else if Num == 1 && string(message)[len(string(message))-14:len(string(message))-6] == "turn.end" {
				fmt.Println("已完成")
				break
			}

		}
		Adata = Adata[:len(Adata)-2400]
		fmt.Println("文件存放路径为:" + config["path"] + config["Language"] + "-" + config["Name"] + "-" + GetLocalTime() + ".mp3")
		ioutil.WriteFile(config["path"]+config["Language"]+"-"+config["Name"]+"-"+GetLocalTime()+".mp3", Adata, 0666)
	}
}
