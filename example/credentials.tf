resource jenkins_credential_username global {
  name     = "global-username2"
  username = "foo"
  # password = "barsoom"
}

resource jenkins_credential_username folder {
  name     = "folder-username2"
  folder   = jenkins_folder.example.name
  username = "folder-foo"
  password = "barsoom"
}
