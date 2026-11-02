// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.111.0-1bfb72c2-20260206-185521
 */

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsFloatingIPProfilesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsFloatingIPProfilesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profiles.is_floating_ip_profiles_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profiles.is_floating_ip_profiles_instance", "profiles.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsFloatingIPProfilesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_floating_ip_profiles" "is_floating_ip_profiles_instance" {
		}
	`)
}
