resource "aws_ecs_service" "ptt-alertor-service" {
  name = "${var.ecs_service_name}"
  cluster = "${aws_ecs_cluster.ptt-alertor.id}"
  task_definition = "${aws_ecs_task_definition.service.arn}"
  desired_count = "${var.ecs_service_desired_count}"
  iam_role = "${aws_iam_role.ptt-alertor_role.arn}"
  depends_on = ["aws_iam_role_policy.ptt-alertor_policy"]
  deployment_minimum_healthy_percent = "${var.ecs_service_healthy_min}"
  deployment_maximum_percent = "${var.ecs_service_healthy_max}"

  placement_strategy {
    type = "${var.ecs_service_place_type}"
    field = "${var.ecs_service_place_field}"
  }

  load_balancer {
    elb_name = "${aws_elb.ptt-alertor.name}"
    container_name = "${var.ecs_service_lb_contain_name}"
    container_port = "${var.ecs_service_lb_contain_port}"
  }

  #placement_constraints {
  #  type = "memberOf"
  #  expression = "attribute:ecs.availability-zone in [ us-west-2b ]"
  #}
}
