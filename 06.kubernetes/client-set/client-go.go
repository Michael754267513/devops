
//
// 初始化k8s客户端
func InitClient() (clientset *kubernetes.Clientset, err error) {
	var (
		restConf *rest.Config
	)

	if restConf, err = GetRestConf(); err != nil {
		return
	}

	// 生成clientset配置
	if clientset, err = kubernetes.NewForConfig(restConf); err != nil {
		goto END
	}
END:
	return
}

// 获取k8s restful client配置
func GetRestConf() (restConf *rest.Config, err error) {
	var (
		kubeconfig []byte
	)
	// 读kubeconfig文件
	if kubeconfig, err = ioutil.ReadFile("./admim.config"); err == nil {
		// 生成rest client配置
		if restConf, err = clientcmd.RESTConfigFromKubeConfig(kubeconfig); err != nil {
			goto END
		}
	} else {
		// 从容器SA里面获取配置文件
		restConf, err = rest.InClusterConfig()
		if err != nil {
			goto END
		}
	}

END:
	return
}

func Logger(err error) {
	global.Log.Errorf("获取k8s 配置失败")
	global.Log.Error(err)
}
