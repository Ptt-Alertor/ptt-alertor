resource "aws_cloudwatch_log_group" "ptt-alertor" {
  name = "${var.cloudwatch_log_group_name}"

  tags {
    Environment = "beta"
    Application = "ptt-alertor"
  }
}
