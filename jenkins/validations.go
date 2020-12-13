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
