# queues

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
| [aws_sqs_queue.emails](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue) | resource |
| [aws_sqs_queue.emails_dlq](https://registry.terraform.io/providers/hashicorp/aws/latest/docs/resources/sqs_queue) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_prefix"></a> [prefix](#input\_prefix) | The string to prefix to the name of all resources created by this module | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_emails_arn"></a> [emails\_arn](#output\_emails\_arn) | ARN of the Emails queue |
| <a name="output_emails_name"></a> [emails\_name](#output\_emails\_name) | Name of the Emails queue |
| <a name="output_emails_url"></a> [emails\_url](#output\_emails\_url) | URL of the Emails queue |
<!-- END OF PRE-COMMIT-TERRAFORM DOCS HOOK -->
