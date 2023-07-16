resource "jenkins_credential_secret_file" "global" {
  name   = "global-secret-file"
  filename = "secret-file.txt"
  // This can also be read directy from file like this:
  // filebase64("${path.module}/hello.txt")
  secretbytes = base64encode("My secret file content.")
}

resource "jenkins_credential_secret_file" "folder" {
  name   = "folder-secret-file"
  folder = jenkins_folder.example.id
  filename = "secret-file.txt"
  // This can also be read directy from file like this:
  // filebase64("${path.module}/hello.txt")
  secretbytes = base64encode("My secret file content.")
}
