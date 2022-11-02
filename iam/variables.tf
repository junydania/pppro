variable "role_name" {
  description = "IAM role name"
  type        = string
  default     = "default"
}

variable "policy_name" {
  description = "IAM policy name"
  type        = string
  default     = "default"
}

variable "policy_description" {
  description = "Description of the IAM policy."
  type        = string
  default     = "Managed by Terraform"
}

variable "role_policy_arns" {
  description = "List of policy ARNs to use for default role"
  type        = list(string)
  default     = ["arn:aws:iam::aws:policy/ReadOnlyAccess"]
}

variable "custom_policy_documents" {
  description = "List of managed policy ARNs to use for default role"
  type        = list(string)
  default     = []
}

variable "inline_policy_documents" {
  description = "List of inline policy ARNs to use for default role"
  type        = list(string)
  default     = []
}

variable "max_session_duration" {
  type        = number
  default     = 3600
  description = "The maximum session duration (in seconds) for the role. Can have a value from 1 hour to 12 hours"
}

variable "force_detach_policies" {
  description = "Whether policies should be detached from this role when destroying"
  type        = bool
  default     = false
}

variable "permissions_boundary_arn" {
  description = "Permissions boundary ARN to use for admin role"
  type        = string
  default     = ""
}

variable "role_description" {
  type        = string
  description = "The description of the IAM role that is visible in the IAM role"
  default     = "Managed by Terraform"
}

variable "tags" {
  type    = map(string)
  default = null
}

variable "users" {
    type = list(string)
    description = "List of users"
    defaulr = []
}
