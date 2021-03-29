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


