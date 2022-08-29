package data_sources

import (
	"context"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/lightdash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var organizationSchema = map[string]*schema.Schema{
	"organization_uuid": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "UUID of the organization",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Given name for organizatiom",
	},
}

func DatasourceOrganization() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceOrganizationRead,
		Schema:      organizationSchema,
	}
}

func datasourceOrganizationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)

	var diags diag.Diagnostics

	organization, err := c.GetOrganization()
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("organization_uuid", organization.UUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", organization.Name); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(organization.UUID)

	return diags
}
