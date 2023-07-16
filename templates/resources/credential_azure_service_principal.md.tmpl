# resource_jenkins_credential_azure_service_principal Resource

Manages an Azure Service Principal credential within Jenkins. This credential may then be referenced within jobs that are created.

~> The "client_secret" property may leave plain-text secret id in your state file. If using the property to manage the secret id in Terraform, ensure that your state file is properly secured and encrypted at rest.

~> The Jenkins installation that uses this resource is expected to have the [Azure Credentials Plugin](https://plugins.jenkins.io/azure-credentials/) installed in their system.

## Example Usage

```hcl
resource jenkins_credential_azure_service_principal foo {
    name = "example-secret"
    subscription_id = "01234567-89ab-cdef-0123-456789abcdef"
    client_id = "abcdef01-2345-6789-0123-456789abcdef"
    client_secret = "super-secret"
    tenant = "01234567-89ab-cdef-abcd-456789abcdef"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the credentials being created. This maps to the ID property within Jenkins, and cannot be changed once set.
* `domain` - (Optional) The domain store to place the credentials into. If not set will default to the global credentials store.
* `folder` - (Optional) The folder namespace to store the credentials in. If not set will default to global Jenkins credentials.
* `scope` - (Optional) The visibility of the credentials to Jenkins agents. This must be set to either "GLOBAL" or "SYSTEM". If not set will default to "GLOBAL".
* `description` - (Optional) A human readable description of the credentials being stored.
* `subscription_id` - (Required) The Azure subscription id mapped to the Azure Service Principal.
* `client_id` - (Required) The client id (application id) of the Azure Service Principal.
* `client_secret` - (Optional) The client secret of the Azure Service Principal. Cannot be used with `certificate_id`. Has to be specified, if `certificate_id` is not specified.
* `certificate_id` - (Optional) The certificate reference of the Azure Service Principal, pointing to a Jenkins certificate credential. Cannot be used with `client_secret`. Has to be specified, if `client_secret` is not specified.
* `tenant` - (Required) The Azure Tenant ID of the Azure Service Principal.
* `azure_environment_name` - (Optional) The Azure Cloud enviroment name. Allowed values are "Azure", "Azure China", "Azure Germany", "Azure US Government".
* `service_management_url` - (Optional) Override the Azure management endpoint URL for the selected Azure environment.
* `authentication_endpoint` - (Optional) Override the Azure Active Directory endpoint for the selected Azure environment.
* `resource_manager_endpoint` - (Optional) Override the Azure resource manager endpoint URL for the selected Azure environment.
* `graph_endpoint` - (Optional) Override the Azure graph endpoint URL for the selected Azure environment.

## Attribute Reference

All arguments above are exported.
