output "object" {
  description = "The VPC object created by this module"
  value       = module.foundation.object
}

output "vpc_id" {
  description = "The ID of the VPC created by this module"
  value       = module.foundation.vpc_id
}

output "public_subnet_ids" {
  description = "The public subnet IDs of the VPC created by this module"
  value       = module.foundation.public_subnet_ids
}

output "private_subnet_ids" {
  description = "The private subnet IDs of the VPC created by this module"
  value       = module.foundation.private_subnet_ids
}

output "public_cidr_blocks" {
  description = "CIDR blocks of public subnets created by this module"
  value       = module.foundation.public_cidr_blocks
}

output "private_cidr_blocks" {
  description = "CIDR blocks of private subnets created by this module"
  value       = module.foundation.private_cidr_blocks
}

output "default_security_group_id" {
  description = "The default security group ID of the VPC created by this module"
  value       = module.foundation.default_security_group_id
}

output "public_hostnames" {
  value = aws_instance.public.*.public_dns
}

output "private_hostnames" {
  value = aws_instance.private.*.private_dns
}
