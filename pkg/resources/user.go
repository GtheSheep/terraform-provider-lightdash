package resources

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/lightdash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	roles = []string{
		"admin",
		"editor",
		"viewer",
	}
)

var userSchema = map[string]*schema.Schema{
	"email": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "User email address",
	},
	"first_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "First name",
	},
	"last_name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Last name",
	},
	"role": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "viewer",
		Description:  "Type of role for the user, one of viewer/ editor/ admin",
		ValidateFunc: validation.StringInSlice(roles, false),
	},
}

func ResourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Schema: userSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userID := d.Id()

	user, err := c.GetUser(userID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("first_name", user.FirstName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("last_name", user.LastName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("email", user.Email); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role", user.Role); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)

	var diags diag.Diagnostics

	firstName := d.Get("first_name").(string)
	lastName := d.Get("last_name").(string)
	email := d.Get("email").(string)
	role := d.Get("role").(string)

	user, err := c.CreateUser(email, firstName, lastName, role)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(*user.UserUUID))

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)
	userID := d.Id()

	if d.HasChange("role") {
		_, err = c.UpdateUser(userID, d.Get("role").(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceUserRead(ctx, d, m)
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)
	userID := d.Id()

	var diags diag.Diagnostics

	status, err = c.DeleteUser(userID)
	if (status != "ok") || (err != nil) {
		return diag.FromErr(err)
	}

	return diags
}
