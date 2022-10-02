# Based on
# https://github.com/nelg/terraform-aws-acmdemo/blob/master/ssl_cert.tf

locals {
  api_domain_name = "${local.project_name}.api.shaungc.com"
}

data "aws_route53_zone" "public" {
  name         = "shaungc.com"
  private_zone = false
}

# Resource `aws_acm_certificate`:
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/acm_certificate
resource "aws_acm_certificate" "api" {
  domain_name       = local.api_domain_name
  validation_method = "DNS"

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_route53_record" "dns_validation" {
  allow_overwrite = true
  name            = tolist(aws_acm_certificate.api.domain_validation_options)[0].resource_record_name
  records         = [tolist(aws_acm_certificate.api.domain_validation_options)[0].resource_record_value]
  type            = tolist(aws_acm_certificate.api.domain_validation_options)[0].resource_record_type
  zone_id         = data.aws_route53_zone.public.id
  ttl             = 60
}

# Resource `aws_acm_certificate_validation`
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/acm_certificate_validation
resource "aws_acm_certificate_validation" "api" {
  certificate_arn         = aws_acm_certificate.api.arn
  validation_record_fqdns = [aws_route53_record.dns_validation.fqdn]
}

# Based on
# https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/apigatewayv2_domain_name
resource "aws_route53_record" "api" {
  name    = local.api_domain_name
  type    = "A"
  zone_id = data.aws_route53_zone.public.zone_id

  alias {
    name                   = module.api.apigatewayv2_domain_name_configuration[0].target_domain_name
    zone_id                = module.api.apigatewayv2_domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}
