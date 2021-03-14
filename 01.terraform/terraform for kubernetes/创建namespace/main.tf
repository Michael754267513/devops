provider "kubernetes" {
  config_path    = "~/.kube/config"
  #config_context = "my-context" 多环境使用
}

resource "kubernetes_namespace" "example" {
  metadata {
    name = "my-first-namespace"
  }
}


/*
provider "kubernetes" {
  host = "https://104.196.242.174"

  client_certificate     = "${file("~/.kube/client-cert.pem")}"
  client_key             = "${file("~/.kube/client-key.pem")}"
  cluster_ca_certificate = "${file("~/.kube/cluster-ca-cert.pem")}"
}
*/