# jenkins_job Resource

Manages a job within Jenkins.

## Example Usage

```hcl
resource jenkins_folder example {
  name = "folder-name"
}

resource jenkins_job example {
  name     = "example"
  folder   = jenkins_folder.example.id
  template = file("${path.module}/job.xml")

  parameters = {
    description = "An example job created from Terraform"
  }
}
```

And in `job.xml`:

```xml
<flow-definition plugin="workflow-job@2.25">
  <actions/>
  <description>{{ .Parameters.description }}</description>
  <keepDependencies>false</keepDependencies>
  <properties/>
  <definition class="org.jenkinsci.plugins.workflow.cps.CpsScmFlowDefinition" plugin="workflow-cps@2.59">
    <scm class="hudson.plugins.git.GitSCM" plugin="git@3.9.1">
      <configVersion>2</configVersion>
      <userRemoteConfigs>
        <hudson.plugins.git.UserRemoteConfig>
          <url>https://github.com/taiidani/terraform-provider-jenkins.git</url>
          <credentialsId>github</credentialsId>
        </hudson.plugins.git.UserRemoteConfig>
      </userRemoteConfigs>
      <branches>
        <hudson.plugins.git.BranchSpec>
          <name>*/master</name>
        </hudson.plugins.git.BranchSpec>
      </branches>
      <doGenerateSubmoduleConfigurations>false</doGenerateSubmoduleConfigurations>
      <submoduleCfg class="list"/>
      <extensions/>
    </scm>
    <scriptPath>Jenkinsfile</scriptPath>
    <lightweight>true</lightweight>
  </definition>
  <triggers/>
  <disabled>false</disabled>
</flow-definition>
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the job being created.
* `folder` - (Optional) The folder namespace to store the job in. If creating in a nested folder structure you may separate folder names with `/`, such as `parent/child`. This name cannot be changed once the folder has been created, and all parent folders must be created in advance.
* `parameters` - (Optional) A map of string values that are passed into the template for rendering.
* `template` - (Required) A Jenkins-compatible XML template to describe the job. You can retrieve an existing jobs' XML by appending `/config.xml` to its URL and viewing the source in your browser. The `template` property is rendered using a Golang template that takes the other resource arguments as variables. Do not include the XML prolog in the definition.

## Attribute Reference

All arguments above are exported.
