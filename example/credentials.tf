resource jenkins_credential_username global {
  username = "foo"
  password = "barsoom"
}

resource jenkins_credential_username folder {
  folder   = jenkins_folder.example.name
  username = "folder-foo"
  password = "barsoom"
}
