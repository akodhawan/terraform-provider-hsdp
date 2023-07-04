package tdr

import (
	"context"
	
	// "github.com/hashicorp/go-uuid"
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
		},
	}

}

func dataSourceTDRContractRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	organization_namespace := d.Get("organization_namespace").(string)

	endpoint := d.Get("tdr_endpoint").(string)

	dt ,_ := helpers.CollectDataType(d)

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

	tdrcontracts, resp, err := client.Contracts.GetContract(&contractOptions)
	// if err != nil {
	// 	return diag.FromErr(err)
	// } //Update to allow empty results
	d.SetId(organization_namespace+dataType)
	_ = d.Set("contracts", tdrcontracts)
	_ = d.Set("response", resp)
	_ = d.Set("err", err)

	return diags
}
