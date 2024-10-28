package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/lightdash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLightdashProjectResource(t *testing.T) {

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	nameDatabricks := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLightdashProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLightdashProjectResourceBasicConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLightdashProjectExists("lightdash_project.test_project"),
					resource.TestCheckResourceAttr("lightdash_project.test_project", "name", name),
				),
			},
			// MODIFY
			// IMPORT
			{
				ResourceName:            "lightdash_project.test_project",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLightdashProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLightdashProjectResourceBasicDatabricksConfig(nameDatabricks),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLightdashProjectExists("lightdash_project.test_databricks_project"),
					resource.TestCheckResourceAttr("lightdash_project.test_databricks_project", "name", nameDatabricks),
				),
			},
			// MODIFY
			// IMPORT
			{
				ResourceName:            "lightdash_project.test_databricks_project",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccLightdashProjectResourceBasicConfig(name string) string {
	return fmt.Sprintf(`
data "lightdash_organization" "test_org" {
}

resource "lightdash_project" "test_project" {
    name = "%s"
    organization_uuid = data.lightdash_organization.test_org.organization_uuid
    type = "DEFAULT"
    dbt_connection_repository = "gthesheep/terraform-provider-dbt-cloud"
		warehouse_connection_type = "snowflake"
    warehouse_connection_account = "abc-123.eu-west-1"
    warehouse_connection_role = "ACCOUNTADMIN"
    warehouse_connection_database = "DB"
    warehouse_connection_warehouse = "TEST_WH"
}
`, name)
}

func testAccLightdashProjectResourceBasicDatabricksConfig(name string) string {
	return fmt.Sprintf(`
data "lightdash_organization" "test_org" {
}

resource "lightdash_project" "test_databricks_project" {
    name = "%s"
    organization_uuid = data.lightdash_organization.test_org.organization_uuid
    type = "DEFAULT"
    dbt_connection_repository = "gthesheep/terraform-provider-dbt-cloud"
		warehouse_connection_type = "databricks"
    databricks_connection_server_host_name = "help-im-on-databricks.com"
    databricks_connection_http_path = "moo/baa"
    databricks_connection_personal_access_token = "abcdefg123"
    databricks_connection_catalog = "PROD"
}
`, name)
}

func testAccCheckLightdashProjectExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*lightdash.Client)
		_, err := apiClient.GetProject(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLightdashProjectDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*lightdash.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lightdash_project" {
			continue
		}

		_, err := apiClient.GetProject(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Project still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
