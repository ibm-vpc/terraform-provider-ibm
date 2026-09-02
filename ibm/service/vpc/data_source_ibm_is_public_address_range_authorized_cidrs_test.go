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

func TestAccIBMIsPublicAddressRangeAuthorizedCidrsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeAuthorizedCidrsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidrs.is_public_address_range_authorized_cidrs_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidrs.is_public_address_range_authorized_cidrs_instance", "authorized_cidrs.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeAuthorizedCidrsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_authorized_cidrs" "is_public_address_range_authorized_cidrs_instance" {
			allocation_profile_family = "provider"
			availability_mode = "zonal"
		}
	`)
}
