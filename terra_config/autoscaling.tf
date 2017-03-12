resource "aws_autoscaling_group" "ecs" {
  name                 = "${var.autoscaling_group_name}"
  availability_zones   = ["${split(",", var.availability_zones)}"]
  launch_configuration = "${aws_launch_configuration.ecs.name}"
  min_size             = "${var.autoscaling_min_size}"
  max_size             = "${var.autoscaling_max_size}"
  desired_capacity     = "${var.autoscaling_desired_size}"
}
/* SSH key pair */
resource "aws_key_pair" "ecs" {
  key_name   = "${var.key_name}"
  public_key = "${file(var.key_file)}"
}

/**
 * Launch configuration used by autoscaling group
 */
resource "aws_launch_configuration" "ecs" {
  name                 = "${var.launch_configuration_name}"
  #image_id             = "${lookup(var.amis, var.region)}"
  image_id             = "ami-022b9262"
  instance_type        = "${var.instance_type}"
  key_name             = "${aws_key_pair.ecs.key_name}"
  iam_instance_profile = "${aws_iam_instance_profile.ecs.id}"
  security_groups      = ["${aws_security_group.ecs.id}"]
  iam_instance_profile = "${aws_iam_instance_profile.ecs.name}"
  user_data            = "#!/bin/bash\necho ECS_CLUSTER=${aws_ecs_cluster.ptt-alertor.name} > /etc/ecs/ecs.config\nmkdir /etc/ecs/config\nyum install -y aws-cli\naws s3 cp s3://ptt-alertor-bucket/config/redis.json /etc/ecs/config/\naws s3 cp s3://ptt-alertor-bucket/config/mailgun.json /etc/ecs/config/\n "
}
