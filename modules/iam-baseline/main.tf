provider "aws" {
  region = var.region
}

data "aws_caller_identity" "current" {
}

resource "aws_iam_user" "admin" {
  name = "admin"
}

module "iam_baseline" {

  source  = "nozaq/secure-baseline/aws//modules/iam-baseline"
  version = "0.16.0"

  aws_account_id                 = data.aws_caller_identity.current.account_id
  support_iam_role_principal_arn = aws_iam_user.admin.arn

  tags = var.tags

}
