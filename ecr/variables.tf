variable "scan_images_on_push" {
  type        = bool
  description = "Indicates whether images are scanned after being pushed to the repository (true) or not (false)"
  default     = true
}

variable "untagged_description" {
  type        = string
  default     = "Remove untagged images"
  description = "Set to any description for untagged rules"
}

variable "untagged_tag_status" {
  type        = string
  default     = "untagged"
  description = "Set to any tag status for untagged rules"
}

variable "untagged_count_type" {
  type        = string
  default     = "imageCountMoreThan"
  description = "Set to any count type for untagged rules"
}

variable "untagged_count_number" {
  type        = number
  default     = 1
  description = "Set to any count number for untagged rules"
}

variable "untagged_count_unit" {
  type        = string
  default     = null
  description = "Set to `days` count unit on if countType is set to sinceImagePushed for untagged rules"
}

variable "remove_description" {
  type        = string
  default     = "Rotate images when it reaches max count images stored"
  description = "Set to any description for remove rules"
}

variable "remove_tag_status" {
  type        = string
  default     = "any"
  description = "Set to any tag status for remove rules"
}

variable "remove_count_type" {
  type        = string
  default     = "imageCountMoreThan"
  description = "Set to any count type for remove rules"
}

variable "protected_description" {
  type        = string
  default     = "Protects images tagged"
  description = "Set to any description for protected rules"
}

variable "protected_tag_status" {
  type        = string
  default     = "tagged"
  description = "Set to any tag status for protected rules"
}

variable "protected_count_type" {
  type        = string
  default     = "imageCountMoreThan"
  description = "Set to any count type for protected rules"
}

variable "protected_count_number" {
  type        = number
  default     = 999999
  description = "Set to any count number for protected rules"
}

variable "max_image_count" {
  type        = number
  description = "How many Docker Image versions AWS ECR will store"
  default     = 500
}

variable "image_names" {
  type        = list(string)
  default     = []
  description = "List of Docker local image names, used as repository names for AWS ECR "
}

variable "image_tag_mutability" {
  type        = string
  default     = "IMMUTABLE"
  description = "The tag mutability setting for the repository. Must be one of: `MUTABLE` or `IMMUTABLE`"
}

variable "enable_lifecycle_policy" {
  type        = bool
  description = "Set to false to prevent the module from adding any lifecycle policies to any repositories"
  default     = true
}

variable "enable_iam_policy" {
  type        = bool
  description = "Set to false to prevent the module from adding any iam policies to any repositories"
  default     = true
}

variable "protected_tags" {
  type        = set(string)
  description = "Name of image tags prefixes that should not be destroyed. Useful if you tag images with names like `dev`, `staging`, and `prod`"
  default     = []
}

variable "encryption_configuration" {
  type = object({
    encryption_type = string
    kms_key         = any
  })
  description = "ECR encryption configuration"
  default     = null
}

variable "policy" {
  type        = string
  default     = ""
  sensitive   = true
  description = "A valid policy JSON document. For more information about building AWS IAM policy documents with Terraform."
}
