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

func TestAccIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocations.is_public_address_range_authorized_cidr_allocations_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocations.is_public_address_range_authorized_cidr_allocations_instance", "authorized_cidr_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocations.is_public_address_range_authorized_cidr_allocations_instance", "allocations.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_authorized_cidr_allocations" "is_public_address_range_authorized_cidr_allocations_instance" {
			authorized_cidr_id = "r134-7be42030-e392-43b0-9ae8-2a8f2798c6f1"
			allocations_resource_type = "public_address_range"
		}
	`)
}
