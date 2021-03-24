# jenkins_job Data Source

Get the attributes of a job within Jenkins.

## Example Usage

```hcl
data "jenkins_job" "example" {
  name        = "job-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the job being read.
* `folder` - (Optional) The folder namespace containing this job.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The full canonical job path, E.G. `/job/job-name`.
* `template` - A Jenkins-compatible XML template to describe the job.
