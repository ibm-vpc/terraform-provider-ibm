// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageCompatibleComponentsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageCompatibleComponentsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_compatible_components.is_image_compatible_components_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_compatible_components.is_image_compatible_components_instance", "image_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_compatible_components.is_image_compatible_components_instance", "image_components.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageCompatibleComponentsDataSourceConfigBasic() string {
	return `
		data "ibm_is_images" "ds_images" {
			visibility = "public"
		}

		data "ibm_is_image_compatible_components" "is_image_compatible_components_instance" {
			image_id = data.ibm_is_images.ds_images.images.0.id
		}
	`
}
