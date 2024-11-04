// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageInstanceProfilesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageInstanceProfilesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_instance_profiles.is_image_instance_profiles_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_instance_profiles.is_image_instance_profiles_instance", "instance_profiles.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageInstanceProfilesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_image_instance_profiles" "is_image_instance_profiles_instance" {
			identifier = "id"
		}
	`)
}
