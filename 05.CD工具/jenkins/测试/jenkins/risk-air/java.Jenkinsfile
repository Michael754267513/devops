#!groovy
@Library(value="handpay-library@master", changelog=true) _

// def gitlabWebhook = new org.devops.gitlab()
def formatMsg = new org.devops.tools()
def gitlabWebhook= new org.devops.gitlab()
def tomail= new org.devops.toemail()

String branchName = "${env.branch}"

branchName = branch - "refs/heads/"

pipeline {
    agent any
    
    environment {
        // harbor_url = 'https://harbor.devops.hpay' 
        image  =  "harbor.devops.hpay"+"/"+"${system}"+"/"+"${projectName}"+":"+"${commitSha}"
    }
    options {
        timeout(time: 1, unit: 'HOURS') 
        timestamps()
    }
 
    
    stages {
        stage('获取代码') {
            steps { 

                checkout([$class: 'GitSCM', branches: [[name: "${branchName}"]], extensions: [], userRemoteConfigs: [[credentialsId: 'gitlab-user', url: "${gitlabUrl}"]]]) 
               
            }
        }
        stage('编译') {
            steps { 
                
                    sh "/opt/apache-maven-3.8.6/bin/mvn clean install"
                
               
            }
        }
        stage('打包镜像') {
            steps { 
                 
                    sh "docker build -t  harbor.devops.hpay/${system}/${projectName}:${commitSha} ."
                
               
            }
        }
        stage('推送镜像仓库') {
            steps { 
                
                  sh "docker push  harbor.devops.hpay/${system}/${projectName}:${commitSha}"
               
            }
        } 
        stage('部署') {
            steps { 
                script {
                    formatMsg.PrintMes("项目：${projectName} 分支：${branchName} 部署完成 ","green")
                }
               
            }
        } 
     
    }
    post {
        always {
            echo 'always'
            // deleteDir() /* clean up our workspace */
        }
        success {
            echo 'I succeeeded!'
            script {
                 formatMsg.PrintMes('I succeeeded!',"red")
                 tomail.Email("succeeeded","${tomailsend}")
                 //gitlabWebhook.ChangeCommitStatus("${project_id}","${commitSha}","success")
            }
                
            
            
        }
        unstable {
            echo 'I am unstable :/'
            script { 
                 tomail.Email("unstable","${tomailsend}")
                 //gitlabWebhook.ChangeCommitStatus("${project_id}","${commitSha}","success")
            }
                
        }
        failure {
            echo 'I failed :('
            script { 
                 tomail.Email("failure","${tomailsend}")
                 //gitlabWebhook.ChangeCommitStatus("${project_id}","${commitSha}","success")
            }
            
        }
        changed {
            echo 'Things were different before...'
        }
        aborted {
            echo '喔霍取消了'
            script { 
                 tomail.Email("喔霍取消了","${tomailsend}")
                 //gitlabWebhook.ChangeCommitStatus("${project_id}","${commitSha}","success")
            }
            
        }
    }
}
