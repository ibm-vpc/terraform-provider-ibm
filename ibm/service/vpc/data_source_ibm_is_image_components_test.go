// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageComponentsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageComponentsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_components.is_image_components_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_components.is_image_components_instance", "image_components.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageComponentsDataSourceConfigBasic() string {
	return `
		data "ibm_is_image_components" "is_image_components_instance" {
		}
	`
}
