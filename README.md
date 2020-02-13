# terraform-aws-foundation-minimal

## Overview

This repository holds a module for AWS that provides a minimal [foundation](./docs/architecture.md) configuration. A public and private subnet are provided within available AZs for a specified region.

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:-----:|
| cidr\_block | CIDR block to configure public and private subnets with | `string` | `"10.0.0.0/16"` | no |
| environment | The environment this module will run in | `string` | n/a | yes |
| region | The region this module will run in | `string` | n/a | yes |
| subnet\_count | Number of `aws\_subnet` resources to distribute among AZs | `number` | `3` | no |

## Outputs

| Name | Description |
|------|-------------|
| default\_security\_group\_id | The default security group ID of the VPC created by this module |
| object | The VPC object created by this module |
| private\_cidr\_blocks | CIDR blocks of private subnets created by this module |
| private\_subnet\_ids | The private subnet IDs of the VPC created by this module |
| public\_cidr\_blocks | CIDR blocks of public subnets created by this module |
| public\_subnet\_ids | The public subnet IDs of the VPC created by this module |
| vpc\_id | The ID of the VPC created by this module |

## Usage

```hcl
module "foundation" {
  source = "..."

  environment = "staging"
  region      = "us-west-1"
}

resource "aws_instance" "public" {
  ami                         = data.aws_ami.debian.id
  instance_type               = "t2.nano"
  key_name                    = local.id
  subnet_id                   = element(module.foundation.public_subnet_ids, 0)
  vpc_security_group_ids      = [module.foundation.default_security_group_id]
  associate_public_ip_address = true
}

resource "aws_instance" "private" {
  ami                         = data.aws_ami.debian.id
  instance_type               = "t2.nano"
  key_name                    = local.id
  subnet_id                   = element(module.foundation.private_subnet_ids, 0)
  vpc_security_group_ids      = [module.foundation.default_security_group_id]
  associate_public_ip_address = false
}
```

Check out the [examples](../examples) for fully-working sample code that the [tests](../test) exercise. Testing patterns are documented [here](./docs/architecture.md#testing).

---

This repo has the following folder structure:

* root folder: The root folder contains a single, standalone, reusable, production-grade module.
* [modules](./modules): This folder may contain supporting modules to the root module.
* [examples](./examples): This folder shows examples of different ways to configure the root module and is typically exercised by tests.
* [test](./test): Automated tests for the modules and examples.

See the [official docs](https://www.terraform.io/docs/modules/index.html) for further details.

---

This repository was initialized with an Issue Template.
[See here](https://github.com/github/terraform-aws-foundation-minimal/issues/new/choose).
