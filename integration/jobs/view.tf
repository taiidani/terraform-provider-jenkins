resource "jenkins_view" "example" {
  name = "example"
  assigned_projects = [
    jenkins_folder.example.name,
  ]
}

data "jenkins_view" "example" {
  depends_on = [jenkins_view.example]
  name       = "example"
}

output "view" {
  value = data.jenkins_view.example
}
