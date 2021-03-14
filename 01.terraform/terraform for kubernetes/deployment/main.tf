provider "kubernetes" {
  config_path    = "~/.kube/config"
  #config_context = "my-context" 多环境使用
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "${var.namespace}"
  }
}


resource "kubernetes_deployment" "example" {
  metadata {
    name = "terraform-example"
    namespace = "${var.namespace}"
    labels = {
      app = "nginx"
    }
  }

  spec {
    replicas = 3

    selector {
      match_labels = {
        app = "nginx"
      }
    }

    template {
      metadata {
        labels = {
          app = "nginx"
        }
      }

      spec {
        container {
          image = "nginx:1.7.8"
          name  = "example"

          resources {
            limits = {
              cpu    = "0.5"
              memory = "512Mi"
            }
            requests = {
              cpu    = "250m"
              memory = "50Mi"
            }
          }

          liveness_probe {
            http_get {
              path = "/"
              port = 80
            }

            initial_delay_seconds = 3
            period_seconds        = 3
          }
        }
      }
    }

  }

  // 定义依赖的资源
  depends_on = [
    kubernetes_namespace.example,
  ]
}