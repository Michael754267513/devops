/*
		Handpay ServiceMesh

           创建时间: 2020年11月25日15:55:24

	       少侠好武功,一起Giao起来
	  	 我说一Giao,你说Giao
		   一 Giao ？？？？

*/

package global

import (
	"os"

	"github.com/sirupsen/logrus"
)

var (
	Log *logrus.Logger
)

func init() {
	Log = logrus.New()
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	//Log.SetFormatter(&logrus.JSONFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	Log.SetOutput(os.Stdout)
	//设置最低loglevel
	Log.SetLevel(logrus.InfoLevel)
	//Log.SetReportCaller(true)

	file, err := os.OpenFile("logs/k8s-node-manager.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}
}
