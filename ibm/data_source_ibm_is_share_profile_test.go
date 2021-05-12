// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmIsShareProfileDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareProfileDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profile.is_share_profile", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profile.is_share_profile", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profile.is_share_profile", "family"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profile.is_share_profile", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_profile.is_share_profile", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIbmIsShareProfileDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_share_profile" "is_share_profile" {
			name = "%s"
		}
	`, shareProfileName)
}
