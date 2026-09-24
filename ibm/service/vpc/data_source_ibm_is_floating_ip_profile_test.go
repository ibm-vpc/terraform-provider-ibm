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

func TestAccIBMIsFloatingIPProfileDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsFloatingIPProfileDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profile.is_floating_ip_profile_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profile.is_floating_ip_profile_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profile.is_floating_ip_profile_instance", "family"),
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profile.is_floating_ip_profile_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profile.is_floating_ip_profile_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_floating_ip_profile.is_floating_ip_profile_instance", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsFloatingIPProfileDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_floating_ip_profile" "is_floating_ip_profile_instance" {
			name = "user-ipv4"
		}
	`)
}
