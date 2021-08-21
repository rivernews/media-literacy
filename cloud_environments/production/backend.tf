terraform {
  backend "s3" {
    bucket = "iriversland-cloud"
    key    = "terraform/media-literacy-prod.remote-terraform.tfstate"
  }
}
