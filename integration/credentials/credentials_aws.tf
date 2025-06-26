resource "jenkins_credential_aws" "global" {
  name          = "global-aws-cred"
  access_key    = "foo"
  secret_key    = "bar"
}

data "jenkins_credential_aws" "global" {
  depends_on = [jenkins_credential_aws.global]
  name       = "global-aws-cred"
}

output "aws_cred" {
  value     = data.jenkins_credential_aws.global
  sensitive = true
}

resource "jenkins_credential_aws" "folder" {
  name         = "folder-aws-cred"
  folder       = jenkins_folder.example.id
  access_key   = "foo"
  secret_key   = "bar"
}

resource "jenkins_credential_aws" "global_iam" {
  name                  = "global-aws-iam"
  access_key            = "foo-global-iam"
  secret_key            = "bar"
  iam_role_arn          = "my-role-arn"
  iam_mfa_serial_number = "my-mfa-serial-number"
}

resource "jenkins_credential_aws" "folder_iam" {
  name                  = "folder-aws-iam"
  folder                = jenkins_folder.example.id
  access_key            = "foo-folder-iam"
  secret_key            = "bar"
  iam_role_arn          = "my-role-arn"
  iam_mfa_serial_number = "my-mfa-serial-number"
}

data "jenkins_credential_aws" "folder_iam" {
  depends_on = [jenkins_credential_aws.folder_iam]
  folder     = jenkins_folder.example.id
  name       = "folder-aws-iam"
}

output "aws_cred_folder" {
  value     = data.jenkins_credential_aws.folder_iam
  sensitive = true
  
}