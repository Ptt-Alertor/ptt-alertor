resource "aws_ecs_cluster" "ptt-alertor" {
  name = "${var.ecs_cluster_name}"
}
