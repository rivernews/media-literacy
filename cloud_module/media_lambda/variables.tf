// unique variables

variable description {
  type = string
  description = "Describe what the function is doing."
}

variable go_handler {
  type = string
  description = "The directory name of the go submodule."
}

variable debug {
  type = bool
  default = true
  description = "Set development mode to debug or not"
}

// upstream variables

variable "project_alias" {
  type        = string
  description = "Name prefix used for step function and related resources, including the domain name, so please only use [0-9a-z_-]"
}

variable environment_name {
  type = string
  default = ""
  description = "Empty string for Production, otherwise the environment name e.g. dev, stage, etc, make sure to use lowercase (s3 bucket only allows lower)"
}

variable "slack_post_webhook_url" {
  type = string
}

variable repo_dir {
  type = string
  description = "The absolute path of git repository path"
}

// downstream variables

variable attach_policy_json {
  type = bool
  default = false
}

variable policy_json {
  type = string
  default = null
}

variable attach_policy_statements {
  type = bool
  default = false
}

variable policy_statements {
  type = map
  default = {}
}

variable environment_variables {
  type = map
  default = {}
}