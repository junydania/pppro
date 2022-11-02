
variable "region" {
  type        = string
  description = "AWS region"
  default     = "us-west-2"
}

variable "availability_zones" {
  type        = list(string)
  description = "List of Availability Zones where subnets will be created"
}