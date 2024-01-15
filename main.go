package main

import (
	"context"
	"fmt"
	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/openapi"
	"github.com/tencent-connect/botgo/token"
	"github.com/tencent-connect/botgo/websocket"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"runtime"
	"time"
)

var databaseURL = "kali:kali@tcp(127.0.0.1:3306)/boot"

type Processor struct {
	method openapi.OpenAPI
}

var processor Processor

type Config struct {
	Appid    string `yaml:"appid"`
	Token    string `yaml:"token"`
	Database string `yaml:"database"`
}

func main() {
	var ConfigData Config
	ConfigData, _ = loadConfig("/config.yaml")
	databaseURL = ConfigData.Database
	ctx := context.Background()
	// 加载 appid 和 token
	botToken := token.New(token.TypeBot)
	if err := botToken.LoadFromConfig(getConfigPath("/config.yaml")); err != nil {
		log.Fatalln(err)
	}

	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(botToken).WithTimeout(3 * time.Second)
	wsInfo, err := api.WS(ctx, nil, "")
	if err != nil {
		log.Fatalln(err)
	}
	processor = Processor{method: api}

	intent := websocket.RegisterHandlers(
		// at 机器人事件，目前是在这个事件处理中有逻辑，会回消息，其他的回调处理都只把数据打印出来，不做任何处理
		ATMessageEventHandler(),
	)
	if err = botgo.NewSessionManager().Start(wsInfo, botToken, &intent); err != nil {
		log.Fatalln(err)
	}

}

func getConfigPath(name string) string {
	_, _, _, ok := runtime.Caller(1)
	if ok {
		return fmt.Sprintf("%s/%s", "/", name)
	}
	return ""
}
func loadConfig(filename string) (Config, error) {
	var config Config

	// 读取文件内容
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		return config, err
	}

	// 解析YAML内容到配置结构体
	err = yaml.Unmarshal(fileContent, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}
