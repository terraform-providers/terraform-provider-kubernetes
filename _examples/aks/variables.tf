variable "location" {
  type = string
  default = "westus2"
}

resource "random_id" "cluster_name" {
  byte_length = 5
}

locals {
  cluster_name                = "tf-k8s-${random_id.cluster_name.hex}"
  cluster_credentials_updated = timestamp()
}
