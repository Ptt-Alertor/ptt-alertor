resource "aws_ecs_service" "ptt-alertor-service" {
  name = "ptt-alertor-service"
  cluster = "${aws_ecs_cluster.ptt-alertor.id}"
  task_definition = "${aws_ecs_task_definition.service.arn}"
  desired_count = 1
  iam_role = "${aws_iam_role.ptt-alertor_role.arn}"
  depends_on = ["aws_iam_role_policy.ptt-alertor_policy"]
  deployment_minimum_healthy_percent = 0
  deployment_maximum_percent = 200

  placement_strategy {
    type = "binpack"
    field = "cpu"
  }

  load_balancer {
    elb_name = "${aws_elb.ptt-alertor.name}"
    container_name = "first"
    container_port = 9090
  }

  #placement_constraints {
  #  type = "memberOf"
  #  expression = "attribute:ecs.availability-zone in [ us-west-2b ]"
  #}
}
