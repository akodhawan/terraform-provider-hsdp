package hsdp

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAIInferenceJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAIInferenceJobsRead,
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}

}

func dataSourceAIInferenceJobsRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	config := m.(*Config)
	endpoint := d.Get("endpoint").(string)
	client, err := config.getAIInferenceClientFromEndpoint(endpoint)
	if err != nil {
		return diag.FromErr(err)
	}

	jobs, _, err := client.Job.GetJobs(nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("ai_inference_compute_targets")

	var names []string
	var ids []string

	for _, job := range jobs {
		names = append(names, job.Name)
		ids = append(ids, job.ID)
	}
	_ = d.Set("names", names)
	_ = d.Set("ids", ids)

	return diags
}