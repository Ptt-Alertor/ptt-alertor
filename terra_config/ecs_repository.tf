resource "aws_ecr_repository" "ptt-alertor-repo" {
  name = "${var.ecr_repository_name}"
}
