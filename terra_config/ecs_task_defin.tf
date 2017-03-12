resource "aws_ecs_task_definition" "service" {
  family = "service"
  container_definitions = "${file("task-definitions/service.json")}"

  volume {
    name = "config"
    host_path = "/etc/ecs/config"
  }

  placement_constraints {
    type = "memberOf"
    expression = "attribute:ecs.availability-zone in [us-west-2a, us-west-2b]"
  }
}

