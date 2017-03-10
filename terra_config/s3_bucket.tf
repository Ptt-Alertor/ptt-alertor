resource "aws_s3_bucket" "ptt-alertor" {
    bucket = "ptt-alertor-bucket"
    acl = "public-read-write"

    tags {
        Name = "My bucket"
        Environment = "Dev"
    }
}
