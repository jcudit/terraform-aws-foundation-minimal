module "foundation" {
  source = "../../"

  environment = "staging"
  region      = "us-west-1"
}

# SSH Key

locals {
  id = "test-${var.environment}-${random_string.id.result}"

  ssh_public_key_path  = "${path.root}"
  public_key_filename  = "${local.ssh_public_key_path}/${local.id}.pub"
  private_key_filename = "${local.ssh_public_key_path}/${local.id}.priv"
}

resource "tls_private_key" "default" {
  algorithm = "RSA"
}

resource "aws_key_pair" "generated" {
  key_name   = local.id
  public_key = tls_private_key.default.public_key_openssh
}

resource "local_file" "private_key_pem" {
  content  = tls_private_key.default.private_key_pem
  filename = local.private_key_filename
}

resource "null_resource" "chmod" {
  triggers = {
    key_data = local_file.private_key_pem.content
  }

  provisioner "local-exec" {
    command = "chmod 600 ${local.private_key_filename}"
  }
}

resource "random_string" "id" {
  length  = 5
  upper   = false
  number  = false
  special = false
}

# AMI

data "aws_ami" "debian" {
  most_recent = true

  filter {
    name   = "name"
    values = ["debian-stretch-hvm-x86_64-gp2-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  # https://wiki.debian.org/Cloud/AmazonEC2Image/Stretch
  owners = ["379101102735"]
}

# Test instances

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
