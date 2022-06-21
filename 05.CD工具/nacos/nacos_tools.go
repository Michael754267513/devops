package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func main() {
	var  (

		server string
		config string
	)
	flag.StringVar(&server, "server", "", "服务器地址")
	flag.StringVar(&config, "config", "", "配置文件")
	flag.Parse()
	if server =="" || config =="" {
		fmt.Errorf("--server= --config= 其中参数为空，请重试")
		os.Exit(999)
	}
	sc := []constant.ServerConfig{
		*constant.NewServerConfig(server, 8848, constant.WithContextPath("/nacos")),
	}

	//create ClientConfig
	cc := *constant.NewClientConfig(
		constant.WithNamespaceId("online"),
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogDir("log"),
		constant.WithCacheDir("cache"),
		constant.WithLogLevel("debug"),
	)

	// create config client
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//publish config
	//config key=dataId+group+namespaceId

	_, err = client.PublishConfig(vo.ConfigParam{
		DataId:  config,
		Group:   "DEFAULT_GROUP",
		Content:  FileToString(config) ,
		Type: "properties",
	})
	if err != nil {
		fmt.Printf("PublishConfig err:%+v \n", err)
		fmt.Errorf("配置文件发布异常")
		os.Exit(9999)
	}
	time.Sleep(1 * time.Second)
}

func FileToString(filename string) string  {
	 f,err := ioutil.ReadFile(filename)
	 if err!=nil {}
	return string(f)
}
