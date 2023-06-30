package tdr

import (
	"context"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/philips-software/go-hsdp-api/tdr"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func DataSourceTDRContract() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTDRContractRead,
		Schema: map[string]*schema.Schema{
			"tdr_endpoint": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"organization_namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data_type": {
				Type:     schema.TypeSet,
				Elem:     dataTypeSchema(),
				MaxItems: 1,
				Optional: true,
			},
			"_count": {
				Type:     schema.TypeInt,
				Default:  100,
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

	dtSystem := d.Get("dataType.system").(string)
	dtCode := d.Get("dataType.code").(string)
	dataType := dtSystem + "|" + dtCode

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

	contracts, _, err := client.Contracts.GetContract(&contractOptions)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("contracts", contracts)
	// _ = d.Set("total", (*bundleResponse).Total)

	result, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(result)

	return diags
}
