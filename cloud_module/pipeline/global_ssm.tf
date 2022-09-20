data aws_ssm_parameter newssite_economy {
    name = "/app/media-literacy/newssites/ECONOMY"
}

locals {
    newssite_economy_tokens = split(",", data.aws_ssm_parameter.newssite_economy.value)
    newssite_economy_alias = local.newssite_economy_tokens[2]
}
