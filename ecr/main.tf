locals {
  image_names = length(var.image_names)
}
resource "aws_ecr_repository" "name" {
  for_each             = toset(var.image_names)
  name                 = each.value
  image_tag_mutability = var.image_tag_mutability

  dynamic "encryption_configuration" {
    for_each = var.encryption_configuration == null ? [] : [var.encryption_configuration]
    content {
      encryption_type = encryption_configuration.value.encryption_type
      kms_key         = encryption_configuration.value.kms_key
    }
  }

  image_scanning_configuration {
    scan_on_push = var.scan_images_on_push
  }
}

locals {
  untagged_image_rule = [{
    rulePriority = length(var.protected_tags) + 1
    description  = var.untagged_description
    selection = {
      tagStatus   = var.untagged_tag_status
      countType   = var.untagged_count_type
      countNumber = var.untagged_count_number
      countUnit   = var.untagged_count_unit
    }
    action = {
      type = "expire"
    }
  }]

  remove_old_image_rule = [{
    rulePriority = length(var.protected_tags) + 2
    description  = var.remove_description
    selection = {
      tagStatus   = var.remove_tag_status
      countType   = var.remove_count_type
      countNumber = var.max_image_count
    }
    action = {
      type = "expire"
    }
  }]

  protected_tag_rules = [
    for index, tagPrefix in zipmap(range(length(var.protected_tags)), tolist(var.protected_tags)) :
    {
      rulePriority = tonumber(index) + 1
      description  = var.protected_description
      selection = {
        tagStatus     = var.protected_tag_status
        tagPrefixList = [tagPrefix]
        countType     = var.protected_count_type
        countNumber   = var.protected_count_number
      }
      action = {
        type = "expire"
      }
    }
  ]
}

resource "aws_ecr_lifecycle_policy" "name" {
  for_each   = toset(var.enable_lifecycle_policy ? var.image_names : [])
  repository = aws_ecr_repository.name[each.value].name

  policy = jsonencode({
    rules = concat(local.protected_tag_rules, local.untagged_image_rule, local.remove_old_image_rule)
  })
}

resource "aws_ecr_repository_policy" "name" {
  for_each   = toset(var.enable_iam_policy ? var.image_names : [])
  repository = aws_ecr_repository.name[each.value].name
  policy     = var.policy
}