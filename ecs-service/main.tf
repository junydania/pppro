locals {
  enabled = var.create_ecs_service
  task_role_arn           = try(var.task_role_arn[0], tostring(var.task_role_arn), "")
  create_task_role        = local.enabled && length(var.task_role_arn) == 0
  task_exec_role_arn      = try(var.task_exec_role_arn[0], tostring(var.task_exec_role_arn), "")
  create_exec_role        = local.enabled && length(var.task_exec_role_arn) == 0
  enable_ecs_service_role = var.network_mode != "awsvpc" && length(var.ecs_load_balancers) <= 1
}

resource "aws_ecs_service" "ignore_changes_desired_count" {
  count                              = local.enabled ? 1 : 0
  name                               = var.ecs_service_name
  task_definition                    = coalesce(var.task_definition, "${join("", aws_ecs_task_definition.default.*.family)}:${join("", aws_ecs_task_definition.default.*.revision)}")
  desired_count                      = var.desired_count
  deployment_maximum_percent         = var.deployment_maximum_percent
  deployment_minimum_healthy_percent = var.deployment_minimum_healthy_percent
  health_check_grace_period_seconds  = var.health_check_grace_period_seconds
  launch_type                        = length(var.capacity_provider_strategies) > 0 ? null : var.launch_type
  platform_version                   = var.launch_type == "FARGATE" ? var.platform_version : null
  scheduling_strategy                = var.launch_type == "FARGATE" ? "REPLICA" : var.scheduling_strategy
  enable_ecs_managed_tags            = var.enable_ecs_managed_tags
  iam_role                           = local.enable_ecs_service_role ? coalesce(var.service_role_arn, join("", aws_iam_role.ecs_service.*.arn)) : null
  wait_for_steady_state              = var.wait_for_steady_state
  force_new_deployment               = var.force_new_deployment
  enable_execute_command             = var.exec_enabled

  dynamic "capacity_provider_strategy" {
    for_each = var.capacity_provider_strategies
    content {
      capacity_provider = capacity_provider_strategy.value.capacity_provider
      weight            = capacity_provider_strategy.value.weight
      base              = lookup(capacity_provider_strategy.value, "base", null)
    }
  }

  dynamic "service_registries" {
    for_each = var.service_registries
    content {
      registry_arn   = service_registries.value.registry_arn
      port           = lookup(service_registries.value, "port", null)
      container_name = lookup(service_registries.value, "container_name", null)
      container_port = lookup(service_registries.value, "container_port", null)
    }
  }

  dynamic "ordered_placement_strategy" {
    for_each = var.ordered_placement_strategy
    content {
      type  = ordered_placement_strategy.value.type
      field = lookup(ordered_placement_strategy.value, "field", null)
    }
  }

  dynamic "placement_constraints" {
    for_each = var.service_placement_constraints
    content {
      type       = placement_constraints.value.type
      expression = lookup(placement_constraints.value, "expression", null)
    }
  }

  dynamic "load_balancer" {
    for_each = var.ecs_load_balancers
    content {
      container_name   = load_balancer.value.container_name
      container_port   = load_balancer.value.container_port
      elb_name         = lookup(load_balancer.value, "elb_name", null)
      target_group_arn = lookup(load_balancer.value, "target_group_arn", null)
    }
  }

  cluster        = var.ecs_cluster_name
  propagate_tags = var.propagate_tags
  tags           = var.tags

  deployment_controller {
    type = var.deployment_controller_type
  }

  dynamic "network_configuration" {
    for_each = var.network_mode == "awsvpc" ? ["true"] : []
    content {
      security_groups  = var.security_group_ids
      subnets          = var.subnet_ids
      assign_public_ip = var.assign_public_ip
    }
  }

  deployment_circuit_breaker {
    enable   = var.circuit_breaker_deployment_enabled
    rollback = var.circuit_breaker_rollback_enabled
  }

  lifecycle {
    ignore_changes = [desired_count]
  }
}