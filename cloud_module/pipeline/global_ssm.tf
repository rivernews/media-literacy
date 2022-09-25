data aws_ssm_parameter newssite_economy {
    name = "/app/media-literacy/newssites/ECONOMY"
}

data aws_ssm_parameter media_table {
  name  = "/app/media-literacy/table"
}

locals {
    newssite_economy_tokens = split(",", data.aws_ssm_parameter.newssite_economy.value)
    newssite_economy_alias = local.newssite_economy_tokens[2]

    _media_table_tokens = split(",", data.aws_ssm_parameter.media_table)
    media_table_arn = local._media_table_tokens[0]
    media_table_id = local._media_table_tokens[1]
}
