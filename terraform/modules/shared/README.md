# shared

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

No modules.

## Resources

| Name | Type |
|------|------|
| [aws_ssm_parameter.base_domain](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/data-sources/ssm_parameter) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_environment"></a> [environment](#input\_environment) | Environment name (stage, prod) | `string` | n/a | yes |
| <a name="input_project"></a> [project](#input\_project) | Project name | `string` | `"gestione-caselo"` | no |
| <a name="input_subdomain"></a> [subdomain](#input\_subdomain) | The subdomain from which the application will be served | `string` | `""` | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_domain"></a> [domain](#output\_domain) | The full domain of the application |
| <a name="output_github_repo"></a> [github\_repo](#output\_github\_repo) | Path of the Github project |
| <a name="output_prefix"></a> [prefix](#output\_prefix) | The string to prefix to the name of all resources created by this module |
| <a name="output_tags"></a> [tags](#output\_tags) | Common tags to apply to all resources |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
