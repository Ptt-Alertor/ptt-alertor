resource "aws_s3_bucket" "ptt-alertor" {
    bucket = "${var.s3_bucket_name}"
    acl = "public-read-write"
    tags {
        Name = "My bucket"
        Environment = "Dev"
    }
    logging {
        target_bucket = "${var.s3_log_name}"
    }
}

resource "aws_s3_bucket" "terraform_state_file" {
    bucket = "${var.s3_tfstate_name}"
    acl = "public-read-write"
    tags {
        Name = "My bucket"
        Environment = "Dev"
    }
}

resource "aws_s3_bucket" "bucket" {
    bucket = "${var.s3_log_name}"
    acl = "public-read-write"
    tags {
        Name = "My bucket"
        Environment = "Dev"
    }
}
