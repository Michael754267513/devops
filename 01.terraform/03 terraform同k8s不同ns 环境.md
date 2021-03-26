# Terraform 多环境部署
- 同一个项目(创建一个nginx的pod) 
- 同一个kubernetes 
- 不同namespaces

## 核心实现
- terraform env 命令来切换不通的化境

## demo应用 nginx
- main.tf 文件
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

- variables.tf 变量文件
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
- 新建两个环境 env1/env2 两个环境
``` 
terraform env new env1
terraform env new env2
```
- 查看环境列表
``` 
[root@allinone ns]# terraform env list 
 # terraform  env 不建议使用，建议用workspace 

  default
  env1
* env2
```
- 切换环境
``` 
[root@allinone ns]# terraform env select env1 
Warning: the "terraform env" family of commands is deprecated.

"Workspace" is now the preferred term for what earlier Terraform versions
called "environment", to reduce ambiguity caused by the latter term colliding
with other concepts.

The "terraform workspace" commands should be used instead. "terraform env"
will be removed in a future Terraform version.

Switched to workspace "env1".
```
- 运行env1 
``` 
[root@allinone ns]# terraform  apply -var "env=env1" -var "namespace=env1"   -auto-approve
kubernetes_namespace.example: Creating...
kubernetes_namespace.example: Creation complete after 0s [id=env1]
kubernetes_deployment.example: Creating...
kubernetes_deployment.example: Creation complete after 8s [id=env1/terraform-example]

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

```
- 切换环境运行env2
``` 
[root@allinone ns]# terraform env select env2
Warning: the "terraform env" family of commands is deprecated.

"Workspace" is now the preferred term for what earlier Terraform versions
called "environment", to reduce ambiguity caused by the latter term colliding
with other concepts.

The "terraform workspace" commands should be used instead. "terraform env"
will be removed in a future Terraform version.

Switched to workspace "env2".
[root@allinone ns]# terraform  apply -var "env=env2" -var "namespace=env2"   -auto-approve
kubernetes_namespace.example: Creating...
kubernetes_namespace.example: Creation complete after 0s [id=env2]
kubernetes_deployment.example: Creating...
kubernetes_deployment.example: Creation complete after 8s [id=env2/terraform-example]

Apply complete! Resources: 2 added, 0 changed, 0 destroyed.

Outputs:

all-ns = tolist([
  "cert-manager",
  "default",
  "env1",
  "ingress-nginx",
  "kube-node-lease",
  "kube-public",
  "kube-system",
  "vela-system",
])

```
- 结果验证
```
[root@allinone ns]# kubectl get ns | grep env
env1              Active   5m35s
env2              Active   3m45s
[root@allinone ns]# kubectl get pod -n env1 
NAME                                 READY   STATUS    RESTARTS   AGE
terraform-example-559dc4dd76-8wvxp   1/1     Running   0          6m11s
terraform-example-559dc4dd76-zlb9j   1/1     Running   0          6m11s
[root@allinone ns]# kubectl get pod -n env2 
NAME                                 READY   STATUS    RESTARTS   AGE
terraform-example-559dc4dd76-pjdr5   1/1     Running   0          4m24s
terraform-example-559dc4dd76-thsjl   1/1     Running   0          4m24s
```
- 注 demo 没有添加依赖,最好先定义下服务间的依赖 在deployment里面定义
``` 
  // 定义依赖的资源
  depends_on = [
    kubernetes_namespace.example,
  ]
```

