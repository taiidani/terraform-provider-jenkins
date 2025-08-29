resource "jenkins_credential_aws" "example" {
  name                  = "example-aws-credential"
  access_key            = "AKIAIOSFODNN7EXAMPLE"
  secret_key            = "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"
  iam_role_arn          = "arn:aws:iam::123456789012:role/MyIAMRoleName"
  iam_mfa_serial_number = "arn:aws:iam::123456789012:mfa/user-name"
}
