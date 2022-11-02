output "role_name" {
  value       = join("", aws_iam_role.this.*.name)
  description = "The name of the IAM role created"
}

output "role_id" {
  value       = join("", aws_iam_role.this.*.unique_id)
  description = "The stable and unique string identifying the role"
}

output "role_arn" {
  value       = join("", aws_iam_role.this.*.arn)
  description = "The Amazon Resource Name (ARN) specifying the role"
}

output "iam_instance_profile" {
  value       = join("", aws_iam_instance_profile.this.*.arn)
  description = "The IAM Instance profile"
}