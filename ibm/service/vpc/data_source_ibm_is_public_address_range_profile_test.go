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

func TestAccIBMIsPublicAddressRangeProfileDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeProfileDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.is_public_address_range_profile_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.is_public_address_range_profile_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.is_public_address_range_profile_instance", "family"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.is_public_address_range_profile_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.is_public_address_range_profile_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.is_public_address_range_profile_instance", "resource_type"),
				),
			},
		},
	})
}

// TestAccIBMIsPublicAddressRangeProfileDataSourceIPv6 verifies the provider-ipv6 profile
// introduced for IPv6 public address ranges. Its targetable_resource_types must list
// virtual_network_interface (not vpc).
func TestAccIBMIsPublicAddressRangeProfileDataSourceIPv6(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsPublicAddressRangeProfileDataSourceConfigIPv6(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.ipv6_profile", "id"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "name", "ibm-ipv6"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "family", "provider"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "ip_version", "ipv6"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profile.ipv6_profile", "href"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "resource_type", "public_address_range_profile"),
					// targetable_resource_types — provider-ipv6 targets a VNI
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "targetable_resource_types.#", "1"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "targetable_resource_types.0.type", "enum"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "targetable_resource_types.0.values.#", "1"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profile.ipv6_profile", "targetable_resource_types.0.values.0", "virtual_network_interface"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeProfileDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_profile" "is_public_address_range_profile_instance" {
			name = "public-address-range-provider-ipv4"
		}
	`)
}

func testAccCheckIBMIsPublicAddressRangeProfileDataSourceConfigIPv6() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_profile" "ipv6_profile" {
			name = "ibm-ipv6"
		}
	`)
}
