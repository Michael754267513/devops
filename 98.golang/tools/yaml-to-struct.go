/*
		Handpay ServiceMesh

           创建时间: 2020年11月25日15:55:24

	       少侠好武功,一起Giao起来
	  	 我说一Giao,你说Giao
		   一 Giao ？？？？

*/

package models

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type K8sNodeManeger struct {
	K8sLogs K8sLogs
}
type K8sLogs struct {
	EnableK8sLog    bool   `json:"enable_k8s_log"`    // 是否开启k8s日志归档
	ArchiveLogTimer string `json:"archive_log_timer"` // 日志归档时间 crontab 格式
	LogPath         string `json:"log_path"`          // 需要归档的日志路径
	ArchiveLogPath  string `json:"archive_log_path"`  // 归档后的日志路径
}

func (c *K8sNodeManeger) GetConfig(config string) (cf *K8sNodeManeger, err error) {
	yaml_file, err := ioutil.ReadFile(config)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(yaml_file, &cf)
	if err != nil {
		return
	}
	return
}
