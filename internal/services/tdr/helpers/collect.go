package helpers

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/tdr"
)

func CollectDataType(d *schema.ResourceData) (tdr.DataType, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dt tdr.DataType
	if v, ok := d.GetOk("data_type"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			dt = tdr.DataType{
				System:  mVi["system"].(string),
				Code:    mVi["code"].(string),
			}
		}
	}
	return dt, diags
}

func CollectDeletionPolicy(d *schema.ResourceData) (tdr.DeletePolicy, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dp tdr.DeletePolicy
	if v, ok := d.GetOk("delete_policy"); ok {
		vL := v.(*schema.Set).List()
		for _, vi := range vL {
			mVi := vi.(map[string]interface{})
			dp = tdr.DeletePolicy{
				Duration:  mVi["duration"].(int),
				Unit:    mVi["unit"].(string),
			}
		}
	}
	return dp, diags
}