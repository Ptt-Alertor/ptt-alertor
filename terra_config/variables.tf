#variable "aws_access_key" {
#  description = "The AWS access key."
#}

#variable "aws_secret_key" {
#  description = "The AWS secret key."
#}

variable "region" {
  description = "The AWS region to create resources in."
  default = "eu-west-1"
}

variable "availability_zones" {
  description = "The availability zones"
  default = "us-west-2a"
}

variable "autoscaling_group_name"      { default = "ecs-asg" }

variable "autoscaling_min_size"        { default = "1"}
variable "autoscaling_max_size"        { default = "10"}
variable "autoscaling_desired_size"    { default = "1"}

variable "launch_config_name"   { default = "ecs" }
variable "cloudwatch_log_group_name"   { default = "ptt-alertor"}

variable "ecs_cluster_name"            { default = "ptt-alertor-cluster" }

variable "ecs_service_name"            { default = "ptt-alertor-service" }
variable "ecs_service_desired_count"   { default = "1" }
variable "ecs_service_healthy_min"     { default = "0" }
variable "ecs_service_healthy_max"     { default = "200" }
variable "ecs_service_place_type"      { default = "binpack" }
variable "ecs_service_place_field"     { default = "cpu" }
variable "ecs_service_lb_contain_name" { default = "first" }
variable "ecs_service_lb_contain_port" { default = "9090" }

variable "s3_bucket_name"              { default = "ptt-alertor-bucket" }
variable "s3_tfstate_name"              { default = "ptt-alertor-terraform-state-file" }
variable "s3_log_name"              { default = "s3-liamlai-log" }

variable "ecr_repository_name" {
  description = "The username to use when connecting to the registry."
  default = "ptt-alertor-repo"
}

variable "registry_docker_image" {
   default = "896146012256.dkr.ecr.us-west-2.amazonaws.com/ptt-alertor-repo"
}

variable "registry_username" {
  description = "The username to use when connecting to the registry."
  default = "Registry"
}

/* ECS optimized AMIs per region */
variable "amis" {
  default = {
    us-west-2      = "ami-f173cc91"
#    ap-northeast-1 = "ami-8aa61c8a"
#    ap-southeast-2 = "ami-5ddc9f67"
#    eu-west-1      = "ami-2aaef35d"
#    us-east-1      = "ami-b540eade"
#    us-west-1      = "ami-5721df13"
#    us-west-2      = "ami-cb584dfb"
  }
}

variable "instance_type" {
  default = "t2.micro"
}

variable "key_name" {
  description = "The aws ssh key name."
  default = "ecs"
}

variable "key_file" {
  description = "The ssh public key for using with the cloud provider."
  default = "~/.ssh/id_rsa.pub"
}
