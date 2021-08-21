terraform {
  backend "s3" {
    bucket = "iriversland-cloud"
    key    = "terraform/media-literacy-dev.remote-terraform.tfstate"
  }
}
