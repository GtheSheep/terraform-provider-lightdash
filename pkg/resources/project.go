package resources

import (
	"context"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/lightdash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	projectTypes = []string{
		"DEFAULT",
		"DEVELOPMENT",
	}
	wareHouseTypes = []string{
		"snowflake",
	}
	dbtConnectionTypes = []string{
		"github",
	}
)

var projectSchema = map[string]*schema.Schema{
	"organization_uuid": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "UUID of the organization to create the project in",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Project name",
	},
	"type": &schema.Schema{
		Type:         schema.TypeString,
		Required:     true,
		Description:  "Type of project to create, either DEFAULT or DEVELOPMENT",
		ValidateFunc: validation.StringInSlice(projectTypes, false),
	},
	"dbt_connection_type": &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "github",
		Description:  "dbt project connection type, currently only support 'github', which is the default",
		ValidateFunc: validation.StringInSlice(dbtConnectionTypes, false),
	},
	"dbt_connection_repository": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Repository name in <org>/<repo> format",
	},
	"dbt_connection_branch": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "main",
		Description: "Branch to use, default 'main'",
	},
	"dbt_connection_project_sub_path": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "/",
		Description: "Sub path to find the project in the repo, default '/'",
	},
	"dbt_connection_host_domain": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "github.com",
		Description: "Host domain of the repo, default 'github.com'",
	},
	"warehouse_connection_type": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "snowflake",
		Description: "Type of warehouse to connect to, currently only 'snowflake', as a default",
	},
	"warehouse_connection_account": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Account identifier, including region/ cloud path",
	},
	"warehouse_connection_role": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Role to connect to the warehouse with",
	},
	"warehouse_connection_database": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Database to connect to",
	},
	"warehouse_connection_schema": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "PUBLIC",
		Description: "Schema to connect to, default 'PUBLIC'",
	},
	"warehouse_connection_client_session_keep_alive": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Client session keep alive param, default `false`",
	},
	"warehouse_connection_warehouse": &schema.Schema{
		Type:        schema.TypeString,
		Required:    true,
		Description: "Warehouse to use",
	},
	"warehouse_connection_threads": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     1,
		Description: "Number of threads to use, default `1`",
	},
}

func ResourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		ReadContext:   resourceProjectRead,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,

		Schema: projectSchema,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceProjectRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)

	var diags diag.Diagnostics

	projectID := d.Id()

	project, err := c.GetProject(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", project.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("organization_uuid", project.OrganisationUUID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", project.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dbt_connection_type", project.DbtConnection.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dbt_connection_repository", project.DbtConnection.Repository); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dbt_connection_branch", project.DbtConnection.Branch); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dbt_connection_project_sub_path", project.DbtConnection.ProjectSubPath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("dbt_connection_host_domain", project.DbtConnection.HostDomain); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_type", project.WarehouseConnection.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_account", project.WarehouseConnection.Account); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_role", project.WarehouseConnection.Role); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_database", project.WarehouseConnection.Database); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_schema", project.WarehouseConnection.Schema); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_client_session_keep_alive", project.WarehouseConnection.ClientSessionKeepAlive); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_warehouse", project.WarehouseConnection.Warehouse); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse_connection_threads", project.WarehouseConnection.Threads); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)

	var diags diag.Diagnostics

	name := d.Get("name").(string)
	organizationUUID := d.Get("organization_uuid").(string)
	projectType := d.Get("type").(string)
	dbtConnectionType := d.Get("dbt_connection_type").(string)
	dbtConnectionRepository := d.Get("dbt_connection_repository").(string)
	dbtConnectionBranch := d.Get("dbt_connection_branch").(string)
	dbtConnectionProjectSubPath := d.Get("dbt_connection_project_sub_path").(string)
	dbtConnectionHostDomain := d.Get("dbt_connection_host_domain").(string)
	warehouseConnectionType := d.Get("warehouse_connection_type").(string)
	warehouseConnectionAccount := d.Get("warehouse_connection_account").(string)
	warehouseConnectionRole := d.Get("warehouse_connection_role").(string)
	warehouseConnectionDatabase := d.Get("warehouse_connection_database").(string)
	warehouseConnectionSchema := d.Get("warehouse_connection_schema").(string)
	warehouseConnectionClientSessionKeepAlive := d.Get("warehouse_connection_client_session_keep_alive").(bool)
	warehouseConnectionWarehouse := d.Get("warehouse_connection_warehouse").(string)
	warehouseConnectionThreads := d.Get("warehouse_connection_threads").(int)

	dbtConnection := lightdash.DbtConnection{
		Type:           dbtConnectionType,
		Repository:     dbtConnectionRepository,
		Branch:         dbtConnectionBranch,
		ProjectSubPath: dbtConnectionProjectSubPath,
		HostDomain:     dbtConnectionHostDomain,
	}
	warehouseConnection := lightdash.WarehouseConnection{
		Type:                   warehouseConnectionType,
		Account:                warehouseConnectionAccount,
		Role:                   warehouseConnectionRole,
		Database:               warehouseConnectionDatabase,
		Warehouse:              warehouseConnectionWarehouse,
		Schema:                 warehouseConnectionSchema,
		ClientSessionKeepAlive: warehouseConnectionClientSessionKeepAlive,
		Threads:                warehouseConnectionThreads,
	}

	project, err := c.CreateProject(organizationUUID, name, projectType, dbtConnection, warehouseConnection)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(project.ProjectUUID)

	resourceProjectRead(ctx, d, m)

	return diags
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// TODO: Implement Updates

	return resourceUserRead(ctx, d, m)
}

func resourceProjectDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*lightdash.Client)

	projectID := d.Id()

	var diags diag.Diagnostics

	status, err := c.DeleteProject(projectID)
	if (status != "ok") || (err != nil) {
		return diag.FromErr(err)
	}

	return diags
}
