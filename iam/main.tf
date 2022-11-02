locals {
  naming_convention = "btc-prod"
  users             = toset(var.users)
}

resource "aws_iam_group" "groups" {
  name = "${local.naming_convention}-group"
}

//Create IAM users
resource "aws_iam_user" "users" {
  for_each      = local.users
  name          = "${local.naming_convention}-user"
  force_destroy = var.force_destroy
  tags          = var.tags
}

resource "aws_iam_group_membership" "team" {
  name = "${local.naming_convention}-group_membership"
  users = var.users
  group = aws_iam_group.groups.name
}


data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "AWS"
      identifiers = [aws_iam_group.groups.arn]
    }
    
    actions = ["sts:AssumeRole"]
  }
}


############################################
# Create IAM role                          #
############################################
resource "aws_iam_role" "this" {
  name                  = "${local.naming_convention}-role"
  assume_role_policy    = data.aws_iam_policy_document.assume_role.json
  description           = var.role_description
  max_session_duration  = var.max_session_duration
  force_detach_policies = var.force_detach_policies
  permissions_boundary  = var.permissions_boundary_arn
}