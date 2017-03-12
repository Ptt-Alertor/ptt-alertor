resource "aws_instance" "example" {
  ami           = "ami-f173cc91"
  instance_type = "t2.micro"
}
