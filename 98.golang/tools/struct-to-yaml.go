/*
		Handpay ServiceMesh

           创建时间: 2020年11月25日15:55:24

	       少侠好武功,一起Giao起来
	  	 我说一Giao,你说Giao
		   一 Giao ？？？？

*/

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.handpay.com.cn/it/k8s-node-manager/models"
	"gopkg.in/yaml.v2"
)

func main() {
	config := models.K8sNodeManeger{K8sLogs: models.K8sLogs{
		EnableK8sLog:    true,
		ArchiveLogTimer: "0 0 04 * * *",
		LogPath:         "/opt/logs",
		ArchiveLogPath:  "/opt/archive",
	}}
	b, err := yaml.Marshal(config)
	if err != nil {
		fmt.Println(err)
	}
	_, err = os.OpenFile("config.yaml", os.O_CREATE|os.O_WRONLY, 0644)
	err = ioutil.WriteFile("config.yaml", b, 0644)

}
