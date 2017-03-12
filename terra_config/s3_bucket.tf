resource "aws_s3_bucket" "ptt-alertor" {
    bucket = "${var.s3_bucket_name}"
    acl = "public-read-write"

    tags {
        Name = "My bucket"
        Environment = "Dev"
    }
}
