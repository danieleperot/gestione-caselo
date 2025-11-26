# cloudwatch

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
| [aws_cloudwatch_log_group.emails](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |
| [aws_cloudwatch_log_group.eventform](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/cloudwatch_log_group) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_environment"></a> [environment](#input\_environment) | The specific environment that we are controlling for the project, like 'stage' or 'prod' | `string` | n/a | yes |
| <a name="input_project"></a> [project](#input\_project) | The name of the project! | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | AWS region where resources are being deployed | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
