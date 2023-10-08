resource "random_pet" "folder_name" {}

resource "jenkins_folder" "example" {
  name        = random_pet.folder_name.id
  description = "A sample folder"
}
