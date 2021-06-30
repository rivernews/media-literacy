terraform {
  backend "s3" {
    bucket = "iriversland-cloud"
    key    = "terraform/media-literacy-workload.remote-terraform.tfstate"
  }
}
