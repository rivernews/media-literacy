provider "aws" {
  # please use env var to pass over credentials
  # https://registry.terraform.io/providers/hashicorp/aws/latest/docs
}

# Multi-provider
# https://www.terraform.io/docs/language/providers/configuration.html
provider "aws" {
  # This provider is solely used for ACM
  # Based on:
  # https://github.com/jareware/howto/blob/master/Using%20AWS%20ACM%20certificates%20with%20Terraform.md
  alias = "acm_provider"
  region = "us-west-2"
}
