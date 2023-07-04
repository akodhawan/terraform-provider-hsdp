package tdr

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/philips-software/go-hsdp-api/tdr"
	"github.com/philips-software/terraform-provider-hsdp/internal/services/tdr/helpers"
	"github.com/philips-software/terraform-provider-hsdp/internal/config"
)

func ResourceTDRContract() *schema.Resource {
	return &schema.Resource{
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CreateContext: resourceTDRContractCreate,
		ReadContext:   resourceTDRContractRead,
		DeleteContext: resourceTDRContractDelete,

		Schema: map[string]*schema.Schema{
			"tdr_endpoint": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"organization": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"data_type": {
				Type:     schema.TypeSet,
				Elem:     dataTypeSchema(),
				MaxItems: 1,
				MinItems: 0,
				Required: true,
				ForceNew: true,
			},
			"send_notifications": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"delete_policy": {
				Type:     schema.TypeSet,
				Elem:     deletePolicySchema(),
				MaxItems: 1,
				ForceNew: true,
				Required: true,
			},
			"json_schema": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func dataTypeSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"system": {
				Type:     schema.TypeString,
				Required: true,
			},
			"code": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func deletePolicySchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"duration": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"unit": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"DAY", "MONTH", "YEAR",
				}, false),
			},
		},
	}
}

func resourceTDRContractCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	endpoint := d.Get("tdr_endpoint").(string)

	client, err := c.GetTDRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	tdrNamespaceOrg := d.Get("organization").(string)
	dataType ,_ := helpers.CollectDataType(d)
	dtId := dataType.System+"|"+dataType.Code
	sendNotifications := d.Get("send_notifications").(bool)

	deletePolicy , _ := helpers.CollectDeletionPolicy(d)
	schema := d.Get("json_schema").(string)

	tdrContract := tdr.Contract{
		Organization:      tdrNamespaceOrg,
		DataType:          dataType,
		SendNotifications: sendNotifications,
		DeletePolicy:      deletePolicy,
		Schema:            json.RawMessage(schema),
	}
	_, resp, err := client.Contracts.CreateContract(tdrContract)
	if err != nil {
		if resp == nil {
			return diag.FromErr(err)
		}
		if resp.StatusCode() != http.StatusConflict {
			return diag.FromErr(err)
		}
		// Search for existing
		contracts, _, err2 := client.Contracts.GetContract(nil)
		if err2 != nil {
			return diag.FromErr(fmt.Errorf("on match attempt during Create conflict: %w", err))
		}
		for _, tdrContract := range contracts {
			if dtId == tdrContract.ID {
				d.SetId(tdrContract.ID)
				return resourceTDRContractRead(ctx, d, m)
			}
		}
		return diag.FromErr(err)
	}
	d.SetId(dtId)
	return resourceTDRContractRead(ctx, d, m)
}

func resourceTDRContractRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*config.Config)

	var diags diag.Diagnostics

	organization_namespace := d.Get("organization").(string)

	endpoint := d.Get("tdr_endpoint").(string)

	dt ,_ := helpers.CollectDataType(d)
	dataType := dt.System + "|" + dt.Code

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

	contract, _, err := client.Contracts.GetContract(&contractOptions)
	if err != nil {
		return diag.FromErr(err)
	}
	_ = d.Set("contract", contract)

	d.SetId(dataType)

	return diags
}


func resourceTDRContractDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*config.Config)

	endpoint := d.Get("tdr_endpoint").(string)

	client, err := c.GetTDRClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}
	defer client.Close()

	d.SetId("") // This is by design currently
	return diags
}