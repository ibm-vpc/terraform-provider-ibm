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

func TestAccIBMIsPublicAddressRangeAuthorizedCIDRAllocationDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRAllocationDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocation.is_public_address_range_authorized_cidr_allocation_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocation.is_public_address_range_authorized_cidr_allocation_instance", "authorized_cidr_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocation.is_public_address_range_authorized_cidr_allocation_instance", "authorized_cidr_allocation_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocation.is_public_address_range_authorized_cidr_allocation_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocation.is_public_address_range_authorized_cidr_allocation_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocation.is_public_address_range_authorized_cidr_allocation_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr_allocation.is_public_address_range_authorized_cidr_allocation_instance", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRAllocationDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_authorized_cidr_allocation" "is_public_address_range_authorized_cidr_allocation_instance" {
			authorized_cidr_id = "r134-7be42030-e392-43b0-9ae8-2a8f2798c6f1"
			authorized_cidr_allocation_id = "r134-703cae15-af08-4438-927b-4588d42c108c"
		}
	`)
}
