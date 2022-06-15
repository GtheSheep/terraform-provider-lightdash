package resources

import (
	"context"

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
		Description: "User email address to send invite link to",
	},
	"role": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "viewer",
		Description:  "Type of role for the user, one of viewer/ editor/ admin",
		ValidateFunc: validation.StringInSlice(roles, false),
	},
	"invite_code": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Original invite code for the user",
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

	email := d.Get("email").(string)
	role := d.Get("role").(string)

	inviteLink, err := c.CreateInviteLink(email)
	if err != nil {
		return diag.FromErr(err)
	}

	newUser := lightdash.User{
		Email:      email,
		InviteCode: &inviteLink.InviteCode,
		UserUUID:   inviteLink.UserUUID,
	}
	if err != nil {
		return diag.FromErr(err)
	}

	user, err := c.UpdateUser(newUser.UserUUID, role)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(user.UserUUID)

	if err := d.Set("invite_code", &inviteLink.InviteCode); err != nil {
		return diag.FromErr(err)
	}

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)
	userID := d.Id()

	if d.HasChange("role") {
		_, err := c.UpdateUser(userID, d.Get("role").(string))
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

	inviteCode := d.Get("invite_code").(string)
	status, err := c.DeleteInviteLink(inviteCode)
	if (status != "ok") || (err != nil) {
		status, err := c.DeleteUser(userID)
		if (status != "ok") || (err != nil) {
			return diag.FromErr(err)
		}
	}

	return diags
}
