# testing terraform install

# will need to setup 'aws cli' on machine running this script
# and then add key buy running 'aws configure'
provider "aws" {
    profile = "default"
    region = var.region
}

# Ubuntu 16.04 LTS AMI = "ami-2757f631"
# Ubuntu 16.10 AMI = "ami-b374d5a5"
resource "aws_instance" "example" {
    ami = "ami-b374d5a5"
    instance_type = "t2.micro"

    provisioner "local-exec"{
        command = "echo ${aws_instance.example.public_ip} > ip_address.txt"
    }
}
