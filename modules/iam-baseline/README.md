# iam-baseline

Provides an IAM baseline for the VPC a _Foundational_ module outputs.

See [the community module this module wraps](https://github.com/nozaq/terraform-aws-secure-baseline/tree/master/modules/iam-baseline) for more details.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:-----:|
| region | The region this module will run in | `string` | n/a | yes |
| tags | Tags to be opportunistically passed to created resources | `map` | `{}` | no |

## Outputs

| Name | Description |
|------|-------------|
| manager\_iam\_role | The IAM role used for the manager user. |
| master\_iam\_role | The IAM role used for the master user. |
| support\_iam\_role | The IAM role used for the support user. |

## Usage

```
module "iam_baseline" {
  source = "./modules/iam-baseline"
  region = var.region
}
```
