# route53

<!-- BEGINNING OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
## Requirements

| Name | Version |
|------|---------|
| <a name="requirement_terraform"></a> [terraform](#requirement\_terraform) | ~> 1.11 |
| <a name="requirement_aws"></a> [aws](#requirement\_aws) | ~> 6.0 |

## Providers

| Name | Version |
|------|---------|
| <a name="provider_aws"></a> [aws](#provider\_aws) | 6.22.1 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_shared"></a> [shared](#module\_shared) | ../../modules/shared | n/a |

## Resources

| Name | Type |
|------|------|
| [aws_route53_zone.main](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/route53_zone) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_environment"></a> [environment](#input\_environment) | The specific environment that we are controlling for the project, like 'stage' or 'prod' | `string` | n/a | yes |
| <a name="input_project"></a> [project](#input\_project) | The name of the project! | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | AWS region where resources are being deployed | `string` | n/a | yes |
| <a name="input_subdomain"></a> [subdomain](#input\_subdomain) | The subdomain from which the application will be served | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_hosted_zone_id"></a> [hosted\_zone\_id](#output\_hosted\_zone\_id) | ID of the Route53 hosted zone created for the domain |
| <a name="output_ns_records"></a> [ns\_records](#output\_ns\_records) | List of nameserver records that must be configured on the parent domain to point the subdomain of the app to the Route53 hosted zone |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
