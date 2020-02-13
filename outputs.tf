output "object" {
  description = "The VPC object created by this module"
  value       = aws_vpc.main
}

output "vpc_id" {
  description = "The ID of the VPC created by this module"
  value       = aws_vpc.main.id
}

output "public_subnet_ids" {
  description = "The public subnet IDs of the VPC created by this module"
  value       = aws_subnet.public.*.id
}

output "private_subnet_ids" {
  description = "The private subnet IDs of the VPC created by this module"
  value       = aws_subnet.private.*.id
}

output "public_cidr_blocks" {
  description = "CIDR blocks of public subnets created by this module"
  value       = aws_subnet.public.*.cidr_block
}

output "private_cidr_blocks" {
  description = "CIDR blocks of private subnets created by this module"
  value       = aws_subnet.private.*.cidr_block
}

output "default_security_group_id" {
  description = "The default security group ID of the VPC created by this module"
  value       = aws_default_security_group.default.id
}
