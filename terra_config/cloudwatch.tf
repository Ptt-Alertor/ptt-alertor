resource "aws_cloudwatch_log_group" "ptt-alertor" {
  name = "ptt-alertor"

  tags {
    Environment = "beta"
    Application = "ptt-alertor"
  }
}
