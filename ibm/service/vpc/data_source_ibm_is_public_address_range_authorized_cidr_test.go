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

func TestAccIBMIsPublicAddressRangeAuthorizedCIDRDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "authorized_cidr_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "allocation.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "availability_mode"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "cidr"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "network_prefix_length"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_authorized_cidr.is_public_address_range_authorized_cidr_instance", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeAuthorizedCIDRDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_authorized_cidr" "is_public_address_range_authorized_cidr_instance" {
			authorized_cidr_id = "r134-ad0758f0-887d-48fd-a93a-bd1d72474585"
		}
	`)
}
