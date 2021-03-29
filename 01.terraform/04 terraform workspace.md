# Terraform 多环境部署
- 同一个项目(创建一个nginx的pod) 
- 同一个kubernetes 
- 不同namespaces
## 文档和代码地址
1. 工作目录 
    https://github.com/Michael754267513/devops/tree/main/01.terraform/terraform%20for%20kubernetes/multi_namespace_project
2. 文档路径：
    https://github.com/Michael754267513/devops/blob/main/01.terraform/04%20terraform%20workspace.md


## 核心实现
- terraform workspace 命令来切换不同的化境
### 创建一个新环境步骤
1. 编写tf文件，定义好需要的资源
2. 新建环境(terraform workspace new <project_namespace>)
3. 运行命令新建(terraform apply ....)
4. 验证资源是否创建成功(kubectl get ...)
### 删除一个环境步骤
1. 切换到环境下面(terraform workspace select <env>)
2. 删除新建的资源(terraform  destroy   -auto-approve&& kubectl get ...)
3. 切换到其他环境(terraform workspace select <env>)
4. 删除环境(terraform workspace delete <env>)
5. 验证资源是否删除(terraform workspace list)

## demo应用 nginx
- main.tf
``` 
provider "kubernetes" {
  config_path    = "~/.kube/config"
}



resource "kubernetes_namespace" "example" {
  metadata {
    name = var.namespace
  }
}

resource "kubernetes_deployment" "example" {
  depends_on = [kubernetes_namespace.example]
  metadata {
    name = "terraform-example"
    namespace = var.namespace
    labels = {
      test = "MyExampleApp"
    }
  }

  spec {
    replicas = var.replicas

    selector {
      match_labels = {
        test = "MyExampleApp"
      }
    }

    template {
      metadata {
        labels = {
          test = "MyExampleApp"
        }
      }

      spec {
        container {
          image = "nginx:1.7.8"
          name  = "example"

          resources {
            limits = {
              cpu    = "50m"
              memory = "10Mi"
            }
            requests = {
              cpu    = "20m"
              memory = "10Mi"
            }
          }

          liveness_probe {
            http_get {
              path = "/"
              port = 80

              http_header {
                name  = "X-Custom-Header"
                value = "Awesome"
              }
            }

            initial_delay_seconds = 3
            period_seconds        = 3
          }
        }
      }
    }
  }
}


data "kubernetes_all_namespaces" "allns" {}

output "all-ns" {
  value = data.kubernetes_all_namespaces.allns.namespaces
}


```
- ariables.tf
``` 
variable "replicas" {
  default = 2
}

variable "env" {
  default = "example"
}

variable "namespace" {
  default = "example"
}
```

### 运行命令
- 查看 新建工作目录

```text
[root@allinone ns]# terraform workspace list  // 查看环境列表
  default

[root@allinone ns]# terraform workspace new prj1 // 新建工作目录
Created and switched to workspace "prj1"!

You're now on a new, empty workspace. Workspaces isolate their state,
so if you run "terraform plan" Terraform will not see any existing state
for this configuration.
[root@allinone ns]# terraform workspace new prj2  // 新建工作目录
Created and switched to workspace "prj2"!

You're now on a new, empty workspace. Workspaces isolate their state,
so if you run "terraform plan" Terraform will not see any existing state
for this configuration.
[root@allinone ns]# 
```
- 查看当前workspace
```text
[root@allinone ns]# terraform workspace show 
prj2
```
#### 运行prj2 
```text
[root@allinone ns]# terraform  apply -var "env=prj2" -var "namespace=prj2"   -auto-approve // prj2 运行tf，传入参数
kubernetes_namespace.example: Creating...
kubernetes_namespace.example: Creation complete after 0s [id=prj2]
kubernetes_deployment.example: Creating...
kubernetes_deployment.example: Creation complete after 7s [id=prj2/terraform-example]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

all-ns = tolist([
  "cert-manager",
  "default",
  "ingress-nginx",
  "kube-node-lease",
  "kube-public",
  "kube-system",
  "vela-system",
]) 
[root@allinone ns]# kubectl get all -n prj2   // 查看运行结果
NAME                                     READY   STATUS    RESTARTS   AGE
pod/terraform-example-559dc4dd76-7z262   1/1     Running   0          18s
pod/terraform-example-559dc4dd76-skdzs   1/1     Running   0          18s

NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/terraform-example   2/2     2            2           18s

NAME                                           DESIRED   CURRENT   READY   AGE
replicaset.apps/terraform-example-559dc4dd76   2         2         2       18s
[root@allinone ns]# 

```
#### 运行prj1
- 切换workspace
```text
[root@allinone ns]# terraform workspace list  // 查看当前工作目录
  default
  prj1
* prj2

[root@allinone ns]# terraform workspace select prj1 // 切换工作目录
Switched to workspace "prj1".
[root@allinone ns]# terraform workspace list 
  default
* prj1
  prj2

[root@allinone ns]# 
```
- 运行命令
```text
[root@allinone ns]# terraform  apply -var "env=prj1" -var "namespace=prj1"   -auto-approve
kubernetes_namespace.example: Creating...
kubernetes_namespace.example: Creation complete after 0s [id=prj1]
kubernetes_deployment.example: Creating...
kubernetes_deployment.example: Creation complete after 7s [id=prj1/terraform-example]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

all-ns = tolist([
  "cert-manager",
  "default",
  "ingress-nginx",
  "kube-node-lease",
  "kube-public",
  "kube-system",
  "prj2",
  "vela-system",
])
[root@allinone ns]# kubectl get all -n prj1 
NAME                                     READY   STATUS    RESTARTS   AGE
pod/terraform-example-559dc4dd76-545tf   1/1     Running   0          23s
pod/terraform-example-559dc4dd76-nzwjt   1/1     Running   0          23s

NAME                                READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/terraform-example   2/2     2            2           23s

NAME                                           DESIRED   CURRENT   READY   AGE

```

### 删除工作目录
- 切换到需要删除的workspace，执行destory
```text
[root@allinone ns]# terraform  destroy   -auto-approve   
kubernetes_deployment.example: Destroying... [id=prj1/terraform-example]
kubernetes_deployment.example: Destruction complete after 0s
kubernetes_namespace.example: Destroying... [id=prj1]
kubernetes_namespace.example: Still destroying... [id=prj1, 10s elapsed]
kubernetes_namespace.example: Still destroying... [id=prj1, 20s elapsed]
kubernetes_namespace.example: Destruction complete after 22s

Destroy complete! Resources: 2 destroyed.
```
-  删除工作目录
```text
[root@allinone ns]# terraform workspace select default
Switched to workspace "default".
[root@allinone ns]# 
[root@allinone ns]# terraform workspace delete prj1
Deleted workspace "prj1"!
```

------------------------------------
**关注不迷路，不定期分享devops相关技术文章**

![](../public/img/11.png) ![](../public/img/12.jpg)