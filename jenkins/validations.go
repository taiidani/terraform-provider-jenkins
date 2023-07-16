package jenkins

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func validateJobName(val interface{}, path cty.Path) diag.Diagnostics {
	if strings.Contains(val.(string), "/") {
		return diag.FromErr(fmt.Errorf("provided name includes path characters. Please use the 'folder' property if specifying a job within a subfolder"))
	}

	return diag.Diagnostics{}
}

func validateFolderName(val interface{}, path cty.Path) diag.Diagnostics {
	return diag.Diagnostics{}
}

// supportedCredentialScopes are the credential scope strings that Jenkins allows to be defined.
var supportedCredentialScopes = []string{"SYSTEM", "GLOBAL"}

// Deprecated: Use stringvalidator.OneOf against the `supportedCredentialScopes`.
func validateCredentialScope(val interface{}, path cty.Path) diag.Diagnostics {
	for _, supported := range supportedCredentialScopes {
		if val == supported {
			return diag.Diagnostics{}
		}
	}
	return diag.Errorf("Invalid scope: %s. Supported scopes are: %s", val, strings.Join(supportedCredentialScopes, ", "))
}
