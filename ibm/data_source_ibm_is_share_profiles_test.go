// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmIsShareProfilesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareProfilesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profiles.is_share_profiles", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profiles.is_share_profiles", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profiles.is_share_profiles", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profiles.is_share_profiles", "next.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profiles.is_share_profiles", "profiles.#"),
					resource.TestCheckResourceAttr("data.ibm_is_share_profiles.is_share_profiles", "family", "tiered"),
					resource.TestCheckResourceAttr("data.ibm_is_share_profiles.is_share_profiles", "href", "https://us-south.iaas.cloud.ibm.com/v1/share/profiles/tier-3iops"),
					resource.TestCheckResourceAttr("data.ibm_is_share_profiles.is_share_profiles", "name", "tier-3iops"),
					resource.TestCheckResourceAttr("data.ibm_is_share_profiles.is_share_profiles", "resource_type", "share_profile"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profiles.is_share_profiles", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIbmIsShareProfilesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_share_profiles" "is_share_profiles" {
		}
	`)
}
