locals {
  ecr_names = {
    "pppro-hello" = {
      image_names           = ["pppro-hello"]
      scan_images_on_push   = false
      untagged_description  = "Delete untagged images for operator"
      untagged_count_type   = "sinceImagePushed"
      untagged_count_number = 15
      untagged_count_unit   = "days"
      image_tag_mutability  = "MUTABLE"
      policy                = data.aws_iam_policy_document.ecr.json
    },
  }
}

data "aws_iam_policy_document" "ecr" {
  version = "2012-10-17"
  statement {
    sid    = "pppro-operations"
    effect = "Allow"
    principals {
      type = "AWS"
      identifiers = [
        "arn:aws:iam::719802944938:root"
      ]
    }
    actions = [
      "ecr:BatchCheckLayerAvailability",
      "ecr:BatchGetImage",
      "ecr:CompleteLayerUpload",
      "ecr:GetDownloadUrlForLayer",
      "ecr:InitiateLayerUpload",
      "ecr:PutImage",
      "ecr:UploadLayerPart"
    ]
  }
}

module "ecr" {
  for_each              = local.ecr_names
  source                = "../../ecr"
  image_names           = each.value["image_names"]
  scan_images_on_push   = each.value["scan_images_on_push"]
  untagged_description  = each.value["untagged_description"]
  untagged_count_type   = each.value["untagged_count_type"]
  untagged_count_number = each.value["untagged_count_number"]
  untagged_count_unit   = each.value["untagged_count_unit"]
  image_tag_mutability  = each.value["image_tag_mutability"]
  policy                = each.value["policy"]
}