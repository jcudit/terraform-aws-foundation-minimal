output "master_iam_role" {
  description = "The IAM role used for the master user."
  value       = module.iam_baseline.master_iam_role
}

output "manager_iam_role" {
  description = "The IAM role used for the manager user."
  value       = module.iam_baseline.manager_iam_role
}

output "support_iam_role" {
  description = "The IAM role used for the support user."
  value       = module.iam_baseline.support_iam_role
}
