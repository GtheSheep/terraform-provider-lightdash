package resources_test

import (
	"os"
	"testing"

	"github.com/gthesheep/terraform-provider-lightdash/pkg/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func providers() map[string]*schema.Provider {
	p := provider.Provider()
	return map[string]*schema.Provider{
		"lightdash": p,
	}
}

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = provider.Provider()
	testAccProviders = map[string]*schema.Provider{
		"lightdash": testAccProvider,
	}
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("LIGHTDASH_URL"); v == "" {
		t.Fatal("LIGHTDASH_URL must be set for acceptance tests")
	}
	if v := os.Getenv("LIGHTDASH_USERNAME"); v == "" {
		t.Fatal("LIGHTDASH_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("LIGHTDASH_PASSWORD"); v == "" {
		t.Fatal("LIGHTDASH_PASSWORD must be set for acceptance tests")
	}
}
