package resources_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/lightdash"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// The destroy is currently testing the invite links in the current state as the user hasn't accepted
func TestAccLightdashUserResource(t *testing.T) {

	email := "gthesheep@gmail.com"
	role := "editor"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLightdashUserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLightdashUserResourceBasicConfig(email),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLightdashUserExists("lightdash_user.test_user"),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "email", email),
				),
			},
			// MODIFY
			{
				Config: testAccLightdashUserResourceFullConfig(email, role),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLightdashUserExists("lightdash_user.test_user"),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "email", email),
					resource.TestCheckResourceAttr("lightdash_user.test_user", "role", role),
				),
			},
			// IMPORT
			{
				ResourceName:            "lightdash_user.test_user",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"invite_code"},
			},
		},
	})
}

func testAccLightdashUserResourceBasicConfig(email string) string {
	return fmt.Sprintf(`
resource "lightdash_user" "test_user" {
    email = "%s"
}
`, email)
}

func testAccLightdashUserResourceFullConfig(email, role string) string {
	return fmt.Sprintf(`
resource "lightdash_user" "test_user" {
    email = "%s"
    role = "%s"
}
`, email, role)
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

		if rs.Primary.Attributes["invite_code"] != "" {
			_, err := apiClient.GetInviteLink(rs.Primary.Attributes["invite_code"])
			if err == nil {
				return fmt.Errorf("Invite link still exists")
			}
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
