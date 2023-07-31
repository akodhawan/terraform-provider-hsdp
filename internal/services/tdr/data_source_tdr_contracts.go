package tdr

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/tdr"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/tdr/helpers"
)

func DataSourceTDRContract() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTDRContractRead,
		Schema: map[string]*schema.Schema{
			"tdr_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"organization_namespace": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"data_type": {
				Type:     schema.TypeSet,
				Elem:     dataTypeSchema(),
				MaxItems: 1,
				MinItems: 0,
				Optional: true,
				ForceNew: true,
			},
			"_count": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"send_notifications": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"delete_policy": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"json_schema": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
		},
	}

}

func dataSourceTDRContractRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	organization_namespace := d.Get("organization_namespace").(string)

	endpoint := d.Get("tdr_endpoint").(string)

	dt, _ := helpers.CollectDataType(d)

	var dataType string

	if dt != (tdr.DataType{}) {
		dataType = dt.System + "|" + dt.Code
	}

	count := d.Get("_count").(int)

	client, err := c.GetTDRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	contractOptions := tdr.GetContractOptions{
		Organization: &organization_namespace,
		DataType:     &dataType,
		Count:        &count,
	}

	tdrcontracts, _, err := client.Contracts.GetContract(&contractOptions)

	if err != nil {
		return diag.FromErr(err)
	} //Update to allow empty results

	d.SetId(organization_namespace + dataType)

	for _, r := range tdrcontracts {
		schema, _ := json.Marshal(r.Schema)
		deletePolicy, _ := json.Marshal(r.DeletePolicy)
		_ = d.Set("organization_namespace", r.Organization)
		_ = d.Set("delete_policy", string(deletePolicy))
		_ = d.Set("send_notifications", r.SendNotifications)
		_ = d.Set("json_schema", string(schema))
	}

	return diags
}
