/**
 * Provides internal access to container ports
 */
resource "aws_security_group" "ecs" {
  name = "ecs-sg"
  description = "Container Instance Allowed Ports"

  ingress {
    from_port = 1
    to_port   = 65535
    protocol  = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port = 0
    to_port = 0
    protocol = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags {
    Name = "ecs-sg"
  }
}
