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
		"databricks",
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
		Sensitive:   true,
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
		Type:         schema.TypeString,
		Optional:     true,
		Default:      "snowflake",
		Description:  "Type of warehouse to connect to, must be one of 'snowflake' or 'databricks', 'snowflake' is the default",
		ValidateFunc: validation.StringInSlice(wareHouseTypes, false),
	},
	"databricks_connection_server_host_name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Databricks - Server host name for connection",
	},
	"databricks_connection_http_path": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Databricks - HTTP path for connection",
	},
	"databricks_connection_personal_access_token": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Databricks - Personal access token for connection",
	},
	"databricks_connection_catalog": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Databricks - Catalog name for connection",
	},
	"warehouse_connection_account": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Snowflake - Account identifier, including region/ cloud path",
	},
	"warehouse_connection_role": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Snowflake - Role to connect to the warehouse with",
	},
	"warehouse_connection_database": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Snowflake - Database to connect to",
	},
	"warehouse_connection_schema": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "PUBLIC",
		Description: "Snowflake - Schema to connect to, default 'PUBLIC'",
	},
	"warehouse_connection_client_session_keep_alive": &schema.Schema{
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "Snowflake - Client session keep alive param, default `false`",
	},
	"warehouse_connection_warehouse": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Snowflake - Warehouse to use",
	},
	"warehouse_connection_threads": &schema.Schema{
		Type:        schema.TypeInt,
		Optional:    true,
		Default:     1,
		Description: "Snowflake - Number of threads to use, default `1`",
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
	if err := d.Set("databricks_connection_server_host_name", project.WarehouseConnection.ServerHostName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("databricks_connection_http_path", project.WarehouseConnection.HTTPPath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("databricks_connection_catalog", project.WarehouseConnection.Catalog); err != nil {
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
	databricksConnectionServerHostName := d.Get("databricks_connection_server_host_name").(string)
	databricksConnectionHttpPath := d.Get("databricks_connection_http_path").(string)
	databricksConnectionPersonalAccessToken := d.Get("databricks_connection_personal_access_token").(string)
	databricksConnectionCatalog := d.Get("databricks_connection_catalog").(string)

	dbtConnection := lightdash.DbtConnection{
		Type:           dbtConnectionType,
		Repository:     dbtConnectionRepository,
		Branch:         dbtConnectionBranch,
		ProjectSubPath: dbtConnectionProjectSubPath,
		HostDomain:     dbtConnectionHostDomain,
	}
	warehouseConnection := lightdash.WarehouseConnection{
		Type: warehouseConnectionType,
	}

	if warehouseConnection.Type == "snowflake" {
		warehouseConnection.Account = warehouseConnectionAccount
		warehouseConnection.Role = warehouseConnectionRole
		warehouseConnection.Database = warehouseConnectionDatabase
		warehouseConnection.Warehouse = warehouseConnectionWarehouse
		warehouseConnection.Schema = warehouseConnectionSchema
		warehouseConnection.ClientSessionKeepAlive = warehouseConnectionClientSessionKeepAlive
		warehouseConnection.Threads = warehouseConnectionThreads
	}
	if warehouseConnection.Type == "databricks" {
		warehouseConnection.ServerHostName = databricksConnectionServerHostName
		warehouseConnection.HTTPPath = databricksConnectionHttpPath
		warehouseConnection.PersonalAccessToken = databricksConnectionPersonalAccessToken
		warehouseConnection.Catalog = databricksConnectionCatalog
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

	return resourceProjectRead(ctx, d, m)
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
