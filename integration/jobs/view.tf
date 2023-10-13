resource "jenkins_view" "example" {
  name   = "example"
  folder = jenkins_folder.example.id
}

data "jenkins_view" "example" {
  depends_on = [jenkins_view.example]
  name       = "example"
  folder     = jenkins_folder.example.id
}

output "view" {
  value = data.jenkins_view.example
}
