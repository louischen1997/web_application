variable "aws_region" {
  type    = string
  default = "us-west-2"
}

variable "source_ami" {
  type    = string
  default = "ami-017fecd1353bcc96e" # Ubuntu 22.04 LTS
}

variable "ssh_username" {
  type    = string
  default = "ubuntu"
}

variable "subnet_id" {
  type    = string
  default = "subnet-0eef5723bff4bb36b"
}

# https://www.packer.io/plugins/builders/amazon/ebs
source "amazon-ebs" "my-ami" {
  region          = "${var.aws_region}"
  ami_name        = "csye6225_${formatdate("YYYY_MM_DD_hh_mm_ss", timestamp())}"
  ami_description = "AMI for CSYE 6225"

  ami_regions = [
    "us-west-2",
  ]
  
  ami_users = [
    "498018012136",
    "928544880062",
    "391108003343",
  ]

  aws_polling {
    delay_seconds = 120
    max_attempts  = 50
  }


  instance_type = "t2.micro"
  source_ami    = "${var.source_ami}"
  ssh_username  = "${var.ssh_username}"
  subnet_id     = "${var.subnet_id}"

  launch_block_device_mappings {
    delete_on_termination = true
    device_name           = "/dev/sda1"
    volume_size           = 8
    volume_type           = "gp2"
  }
}

build {
  sources = ["source.amazon-ebs.my-ami"]

  provisioner "file" {
  source = "./"
  destination = "/home/ubuntu/"
}

  provisioner "shell" {
    environment_vars = [
      "DEBIAN_FRONTEND=noninteractive",
      "CHECKPOINT_DISABLE=1"
    ]
    inline = [
      "sudo apt-get update",
      "sudo apt-get upgrade -y",
      "sudo apt-get install nginx -y",
      "sudo apt-get -y install golang-go",
      
      "sudo cp ./cloudwatch-config.json /opt/",

      //"sudo apt-get -y install mysql-server",
      /*
      "sudo mysql -u root -e \"ALTER USER 'root'@'localhost' IDENTIFIED WITH caching_sha2_password BY 'chenYTCfor6225';\"",
      "sudo mysql --user=root --password=chenYTCfor6225 -e \"create database Goapi;\"",
      "sudo mysql --user=root --password=chenYTCfor6225 -e \"CREATE USER 'newur'@'localhost' IDENTIFIED BY 'ytc6225forclass';\"",
      "sudo mysql --user=root --password=chenYTCfor6225 -e \"GRANT ALL PRIVILEGES ON *.* TO 'newur'@'localhost';\"",
      */

      "wget https://s3.us-west-2.amazonaws.com/amazoncloudwatch-agent-us-west-2/ubuntu/amd64/latest/amazon-cloudwatch-agent.deb",
      "sudo dpkg -i -E ./amazon-cloudwatch-agent.deb",

      "sudo apt-get -y install mysql-client-8.0",
      "sudo apt-get -y install systemd",
      "sudo cp ./app.service /etc/systemd/system/",
      "CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o systemd-test  main.go",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable app.service",
      "sudo systemctl start app.service",

      // "systemctl is-active mysql",
      // "sudo apt-get clean",
    ]
  }
  
}
