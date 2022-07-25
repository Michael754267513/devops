# 基于yaml文件的 agent 指定工作目录


```azure

pipeline {
  agent {
    kubernetes {
      customWorkspace 'some/other/path'
      defaultContainer 'maven'
      yamlFile 'KubernetesPod.yaml'
    }
  }

  stages {
    stage('Run maven') {
      steps {
        sh 'mvn -version'
        sh "echo Workspace dir is ${pwd()}"
      }
    }
  }
}

```