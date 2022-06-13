package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/lightdash"
	"github.com/gthesheep/terraform-provider-lightdash/pkg/resources"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIGHTDASH_URL", nil),
				Description: "URL for your Lighdash instance",
			},
			"username": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIGHTDASH_USERNAME", nil),
				Description: "Username for your Lighdash account",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("LIGHTDASH_PASSWORD", nil),
				Description: "Password for your Lighdash account",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"lightdash_user": resources.ResourceUser(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	url := d.Get("url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)

	var diags diag.Diagnostics

	if (url != "") && (username != "") && (password != "") {
		c, err := lightdash.NewClient(&url, &username, &password)

		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to login to Lightdash",
				Detail:   err.Error(),
			})
			return nil, diags
		}

		return c, diags
	}

	c, err := lightdash.NewClient(nil, nil)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Lightdash client",
			Detail:   err.Error(),
		})
		return nil, diags
	}

	return c, diags
}
