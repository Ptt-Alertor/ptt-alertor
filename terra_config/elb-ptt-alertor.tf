# Create a new load balancer

resource "aws_elb" "ptt-alertor" {
  name = "ptt-alertor-terraform-elb"
  availability_zones = [ "us-west-2a" ,"us-west-2b" ]

  listener {
    instance_port = 80
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }
  listener {
    instance_port = 9090
    instance_protocol = "http"
    lb_port = 9090
    lb_protocol = "http"
  }

  #listener {
  #  instance_port = 22
  #  instance_protocol = "TCP"
  #  lb_port = 22
  #  lb_protocol = "TCP"
  #}
  health_check {
    healthy_threshold = 2
    unhealthy_threshold = 2
    timeout = 3
    target = "HTTP:80/"
    interval = 30
  }

  cross_zone_load_balancing = true
  idle_timeout = 400
  connection_draining = true
  connection_draining_timeout = 400
  security_groups = ["${aws_security_group.ecs.id}"]
  tags {
    Name = "ptt-alertor-terraform-elb"
  }
}
