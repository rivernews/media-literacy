terraform {
  required_version = ">= 1.0.2" # minimum version that supports M1

  required_providers {
    # please use env var to pass over credentials
    # https://registry.terraform.io/providers/hashicorp/aws/latest/docs
    aws = ">= 4.9.0" # minimum version based on `terraform init` error message
  }
}
