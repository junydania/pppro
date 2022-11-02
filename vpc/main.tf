
locals {
  enabled                                   = module.this.enabled
  ipv6_egress_only_internet_gateway_enabled = local.enabled && var.ipv6_egress_only_internet_gateway_enabled
  ipv4_cidr_block_associations = local.enabled ? var.ipv4_additional_cidr_block_associations : {}

  ipv4_primary_cidr_block = var.ipv4_primary_cidr_block
  # assign_generated_ipv6_cidr_block was only briefly deprecated in favor of ipv6_enabled, so it retains
  # precedence. They both defaulted to `true` so we leave the default true.
  assign_generated_ipv6_cidr_block = var.assign_generated_ipv6_cidr_block
  dns_hostnames_enabled            = var.dns_hostnames_enabled
  dns_support_enabled              = var.dns_support_enabled
  classiclink_enabled              = var.classiclink_enabled
  classiclink_dns_support_enabled  = var.classiclink_dns_support_enabled
  default_security_group_deny_all  = var.default_security_group_deny_all
  internet_gateway_enabled         = var.internet_gateway_enabled
}

resource "aws_vpc" "default" {
  count               = local.enabled ? 1 : 0
  cidr_block          = local.ipv4_primary_cidr_block
  ipv4_ipam_pool_id   = try(var.ipv4_primary_cidr_block_association.ipv4_ipam_pool_id, null)
  ipv4_netmask_length = try(var.ipv4_primary_cidr_block_association.ipv4_netmask_length, null)

  ipv6_cidr_block     = try(var.ipv6_primary_cidr_block_association.ipv6_cidr_block, null)
  ipv6_ipam_pool_id   = try(var.ipv6_primary_cidr_block_association.ipv6_ipam_pool_id, null)
  ipv6_netmask_length = try(var.ipv6_primary_cidr_block_association.ipv6_netmask_length, null)

  instance_tenancy                 = var.instance_tenancy
  enable_dns_hostnames             = local.dns_hostnames_enabled
  enable_dns_support               = local.dns_support_enabled
  enable_classiclink               = local.classiclink_enabled
  enable_classiclink_dns_support   = local.classiclink_dns_support_enabled
  assign_generated_ipv6_cidr_block = local.assign_generated_ipv6_cidr_block
  tags                             = module.this.tags
}

# If `aws_default_security_group` is not defined, it will be created implicitly with access `0.0.0.0/0`
resource "aws_default_security_group" "default" {
  count = local.default_security_group_deny_all ? 1 : 0

  vpc_id = aws_vpc.default[0].id
  tags   = merge(module.this.tags, { Name = "Default Security Group" })
}

resource "aws_internet_gateway" "default" {
  count = local.internet_gateway_enabled ? 1 : 0

  vpc_id = aws_vpc.default[0].id
  tags   = module.this.tags
}

resource "aws_egress_only_internet_gateway" "default" {
  count = local.ipv6_egress_only_internet_gateway_enabled ? 1 : 0

  vpc_id = aws_vpc.default[0].id
  tags   = module.this.tags
}

resource "aws_vpc_ipv4_cidr_block_association" "default" {
  for_each = local.enabled ? local.ipv4_cidr_block_associations : {}

  cidr_block          = each.value.ipv4_cidr_block
  ipv4_ipam_pool_id   = each.value.ipv4_ipam_pool_id
  ipv4_netmask_length = each.value.ipv4_netmask_length

  vpc_id = aws_vpc.default[0].id

  dynamic "timeouts" {
    for_each = local.enabled && var.ipv4_cidr_block_association_timeouts != null ? [true] : []
    content {
      create = lookup(var.ipv4_cidr_block_association_timeouts, "create", null)
      delete = lookup(var.ipv4_cidr_block_association_timeouts, "delete", null)
    }
  }
}

resource "aws_vpc_ipv6_cidr_block_association" "default" {
  for_each = local.enabled ? var.ipv6_additional_cidr_block_associations : {}

  ipv6_cidr_block     = each.value.ipv6_cidr_block
  ipv6_ipam_pool_id   = each.value.ipv6_ipam_pool_id
  ipv6_netmask_length = each.value.ipv6_netmask_length

  vpc_id = aws_vpc.default[0].id

  dynamic "timeouts" {
    for_each = local.enabled && var.ipv6_cidr_block_association_timeouts != null ? [true] : []
    content {
      create = lookup(var.ipv6_cidr_block_association_timeouts, "create", null)
      delete = lookup(var.ipv6_cidr_block_association_timeouts, "delete", null)
    }
  }
}