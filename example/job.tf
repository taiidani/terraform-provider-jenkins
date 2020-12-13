resource "jenkins_folder" "example" {
  name        = "folder-name"
  description = "A sample folder"
  template    = <<EOT
<com.cloudbees.hudson.plugins.folder.Folder plugin="cloudbees-folder@6.6">
    <actions/>
    <description>{{ .Description }}</description>
    <properties>
    <com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
        <inheritanceStrategy class="org.jenkinsci.plugins.matrixauth.inheritance.InheritParentStrategy"/>
        {{ range $value := .Permissions }}
        <permission>{{ $value }}</permission>
        {{ end }}
    </com.cloudbees.hudson.plugins.folder.properties.AuthorizationMatrixProperty>
    </properties>
    <icon class="com.cloudbees.hudson.plugins.folder.icons.StockFolderIcon"/>
</com.cloudbees.hudson.plugins.folder.Folder>
EOT

  lifecycle {
    ignore_changes = [template]
  }
}

resource "jenkins_job" "pipeline" {
  name     = "pipeline"
  folder   = jenkins_folder.example.id
  template = file("${path.module}/pipeline.xml")

  parameters = {
    description = "An example pipeline job"
  }
}

resource "jenkins_job" "freestyle" {
  name     = "freestyle"
  folder   = jenkins_folder.example.id
  template = file("${path.module}/freestyle.xml")

  parameters = {
    description = "An example freestyle job"
  }
}
