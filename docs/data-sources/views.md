# jenkins_view Data Source

Get the attributes of a view within Jenkins.

## Example Usage

```hcl
data "jenkins_view" "example" {
  name        = "view-name"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the job being read.
* `folder` - (Optional) The folder namespace containing this job.


## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The full canonical job path, E.G. `/job/job-name`.
* `description` - The description of the view.
* `url` - The full url to the view.
