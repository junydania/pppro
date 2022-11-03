
provider "aws" {
  region = "us-east-1"
}

locals {
  region = "us-east-1"
  name   = "ecs-pppro-app"

  user_data = <<-EOT
    #!/bin/bash
    cat <<'EOF' >> /etc/ecs/ecs.config
    ECS_CLUSTER=${local.name}
    ECS_LOGLEVEL=debug
    EOF
  EOT

  allow_all_ingress_rule = {
    key              = "allow_all_ingress"
    type             = "ingress"
    from_port        = 0
    to_port          = 0 # [sic] from and to port ignored when protocol is "-1", warning if not zero
    protocol         = "-1"
    description      = "Allow all ingress"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }
}

module "vpc" {
  source                           = "../../vpc"
  ipv4_primary_cidr_block          = "10.50.0.0/16"
  name                             = "vpc"
  assign_generated_ipv6_cidr_block = true
  context                          = module.this.context
}

module "subnets" {
  source               = "../../subnets"
  name                 = "kj-subnet"
  availability_zones   = var.availability_zones
  vpc_id               = module.vpc.vpc_id
  igw_id               = [module.vpc.igw_id]
  ipv4_enabled         = true
  ipv4_cidr_block      = [module.vpc.vpc_cidr_block]
  nat_gateway_enabled  = true
  nat_instance_enabled = false
  route_create_timeout = "5m"
  route_delete_timeout = "10m"
  context              = module.this.context
}

module "allow_all_sg" {
  source                     = "../../security-group"
  attributes                 = ["all", "tcp"]
  security_group_description = "Allow All traffic in"
  create_before_destroy      = true
  allow_all_egress           = true
  rules                      = [local.allow_all_ingress_rule]
  vpc_id                     = module.vpc.vpc_id
  context                    = module.this.context
}

data "aws_ssm_parameter" "ecs_optimized_ami" {
  name = "/aws/service/ecs/optimized-ami/amazon-linux-2/recommended"
}

module "autoscaling" {
  source = "../../autoscaling"
  for_each = {
    one = {
      instance_type = "t3.medium"
    }
    two = {
      instance_type = "t3.large"
    }
  }
  name = "${local.name}-${each.key}"

  image_id      = jsondecode(data.aws_ssm_parameter.ecs_optimized_ami.value)["image_id"]
  instance_type = each.value.instance_type

  security_groups                 = [module.allow_all_sg.id]
  user_data                       = base64encode(local.user_data)
  ignore_desired_capacity_changes = true

  create_iam_instance_profile = true
  iam_role_name               = local.name
  iam_role_description        = "ECS role for ${local.name}"
  iam_role_policies = {
    AmazonEC2ContainerServiceforEC2Role = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
    AmazonSSMManagedInstanceCore        = "arn:aws:iam::aws:policy/AmazonSSMManagedInstanceCore"
  }

  vpc_zone_identifier = module.subnets.private_subnet_ids
  health_check_type   = "EC2"
  min_size            = 1
  max_size            = 1
  desired_capacity    = 1

  autoscaling_group_tags = {
    AmazonECSManaged = true
  }
  protect_from_scale_in = true
}

resource "aws_cloudwatch_log_group" "this" {
  name              = "/aws/ecs/${local.name}"
  retention_in_days = 7
}

module "ecs" {
  source       = "../../ecs-cluster"
  cluster_name = local.name
  cluster_configuration = {
    execute_command_configuration = {
      logging = "OVERRIDE"
      log_configuration = {
        cloud_watch_log_group_name = aws_cloudwatch_log_group.this.name
      }
    }
  }

  default_capacity_provider_use_fargate = false

  # Capacity provider - autoscaling groups
  autoscaling_capacity_providers = {
    one = {
      auto_scaling_group_arn         = module.autoscaling["one"].autoscaling_group_arn
      managed_termination_protection = "ENABLED"

      managed_scaling = {
        maximum_scaling_step_size = 5
        minimum_scaling_step_size = 1
        status                    = "ENABLED"
        target_capacity           = 60
      }

      default_capacity_provider_strategy = {
        weight = 60
        base   = 20
      }
    }
    two = {
      auto_scaling_group_arn         = module.autoscaling["two"].autoscaling_group_arn
      managed_termination_protection = "ENABLED"

      managed_scaling = {
        maximum_scaling_step_size = 15
        minimum_scaling_step_size = 5
        status                    = "ENABLED"
        target_capacity           = 90
      }

      default_capacity_provider_strategy = {
        weight = 40
      }
    }
  }
}


module "container_definition" {
  source           = "../../ecs-task-definition"
  container_cpu    = 1024
  environment      = []
  essential        = true
  container_image  = "${module.ecr.pppro-build.repository_url}:latest"
  container_memory = 2048
  entrypoint       = []
  container_name   = "pppro-rest"
  log_configuration = {
    logDriver = "awslogs"
    options = {
      "awslogs-region" : "us-east-1",
      "awslogs-group" : "${aws_cloudwatch_log_group.this.name}",
      "awslogs-stream-prefix" : "ec2"
    }
    options = {
      awslogs-group         = aws_cloudwatch_log_group.this.name
      awslogs-region        = local.region
      awslogs-stream-prefix = "ec2"
    }
  }
  map_environment = {
    "BTC_RPCUSER" = "user"
    "BTC_TXINDEX" = "1"
  }

  port_mappings = [
    {
      containerPort = 8332
      hostPort      = 8332
      protocol      = "tcp"
    },
    {
      containerPort = 8333
      hostPort      = 8333
      protocol      = "tcp"
    }
  ]
  mount_points = [
    {
      containerPath = "/bitcoin/data"
      sourceVolume  = "pppro-efs"
      readOnly      = false
    }
  ]
  depends_on = [module.ecr]

}

module "ecs_service" {
  source                    = "../../ecs-service"
  container_definition_json = module.container_definition.json_map_encoded_list
  ecs_cluster_name          = module.ecs.cluster_name
  ecs_family_name           = "task-pppro-app"
  launch_type               = "EC2"
  enable_ecs_managed_tags   = true
  region                    = "us-east-1"
  encrypt_logs              = false
  kms_enabled               = false
  task_cpu                  = 1048
  task_memory               = 2048
  ecs_load_balancers        = []
  ecs_service_name          = "pppro-bitcoin"
  security_group_ids = [
    module.allow_all_sg.id
  ]
  scheduling_strategy = "DAEMON"
  subnet_ids          = module.subnets.private_subnet_ids
  service_role_arn    = "AWSServiceRoleForECS"
  volumes = [
    {
      host_path                   = ""
      name                        = "pppro-efs"
      docker_volume_configuration = []
      efs_volume_configuration = [
        {
          file_system_id          = module.efs.id
          root_directory          = "/"
          transit_encryption      = "ENABLED"
          transit_encryption_port = 2999
          authorization_config    = []
        }
      ]
    }
  ]
}

module "efs" {
  source = "../../efs"

  region                        = local.region
  vpc_id                        = module.vpc.vpc_id
  subnets                       = module.subnets.private_subnet_ids
  associated_security_group_ids = module.allow_all_sg.id

  access_points = {
    "data" = {
      posix_user = {
        gid            = "1001"
        uid            = "5000"
        secondary_gids = "1002,1003"
      }
      creation_info = {
        gid         = "1001"
        uid         = "5000"
        permissions = "0755"
      }
    }
  }
  transition_to_ia = ["AFTER_7_DAYS"]
}