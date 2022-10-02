locals {
  project_name = var.environment_name != "" ? "${var.project_alias}-${var.environment_name}" : "${var.project_alias}"
  environment = var.environment_name != "" ? var.environment_name : "prod"
}

locals {
  # amd64 is the x86 instruction set
  # arm is not (like M1), not supported by AWS lambda go runtime yet
  # https://stackoverflow.com/questions/26951940/how-do-i-make-go-get-to-build-against-x86-64-instead-of-i386
  go_build_flags = "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 "

  name = replace(var.go_handler, "_", "-")
}