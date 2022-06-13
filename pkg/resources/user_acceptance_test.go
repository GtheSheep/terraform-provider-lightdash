package resources_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/lightdash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLightdashUserResource(t *testing.T) {

	firstName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	lastName := strings.ToUpper(acctest.RandStringFromCharSet(10, acctest.CharSetAlpha))
	email := "gthesheep@gmail.com"
	role := "editor"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLightdashUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLightdashUserResourceBasicConfig(firstName, lastName, email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLightdashUserExists("lightdash_user.test_user"),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "first_name", firstName),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "last_name", lastName),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "email", email),
				),
			},
			// MODIFY
			{
				Config: testAccLightdashUserResourceFullConfig(firstName, lastName, email, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLightdashUserExists("lightdash_user.test_user"),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "first_name", firstName),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "last_name", lastName),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "email", email),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "role", role),
				),
			},
			// IMPORT
			{
				ResourceName:            "lightdash_user.test_user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
			},
		},
	})
}

func testAccLightdashUserResourceBasicConfig(firstName, lastName, email string) string {
	return fmt.Sprintf(`
resource "lightdash_user" "test_user" {
    first_name = "%s"
    last_name = "%s"
    email = "%s"
}
`, firstName, lastName, email)
}

func testAccLightdashUserResourceFullConfig(firstName, lastName, email, role string) string {
	return fmt.Sprintf(`
resource "lightdash_user" "test_user" {
    first_name = "%s"
    last_name = "%s"
    email = "%s"
    role = "%s"
}
`, firstName, lastName, email, role)
}

func testAccCheckLightdashUserExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}
		apiClient := testAccProvider.Meta().(*lightdash.Client)
		_, err := apiClient.GetUser(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching item with resource %s. %s", resource, err)
		}
		return nil
	}
}

func testAccCheckLightdashUserDestroy(s *terraform.State) error {
	apiClient := testAccProvider.Meta().(*lightdash.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "lightdash_user" {
			continue
		}
		_, err := apiClient.GetUser(rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("User still exists")
		}
		notFoundErr := "not found"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}
