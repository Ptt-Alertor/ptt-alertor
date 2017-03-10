resource "aws_iam_role_policy_attachment" "ptt2ptt-policy" {
    role = "${aws_iam_role.ptt-alertor_role.name}"
    policy_arn = "${aws_iam_policy.ptt-alertor-policy.arn}"
}
resource "aws_iam_role_policy_attachment" "ptt2cloudwatch-policy" {
    role = "${aws_iam_role.ptt-alertor_role.name}"
    policy_arn ="${aws_iam_policy.ECS-CloudWatchLogs.arn}"
}
resource "aws_iam_role_policy_attachment" "ptt2ecs-policy" {
    role = "${aws_iam_role.ptt-alertor_role.name}"
    policy_arn ="arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceRole"
}
