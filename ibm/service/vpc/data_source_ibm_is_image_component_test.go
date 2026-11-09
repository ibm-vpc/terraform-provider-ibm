// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageComponentDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageComponentDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "image_component"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "allowed_use.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "dependencies.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "resource_group.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "software.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "visibility"),
				),
			},
		},
	})
}

func TestAccIBMIsImageComponentDataSourceByName(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageComponentDataSourceConfigByName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "visibility"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageComponentDataSourceConfigBasic() string {
	return `
		data "ibm_is_image_components" "components" {}

		data "ibm_is_image_component" "is_image_component_instance" {
			image_component = data.ibm_is_image_components.components.image_components.0.id
		}
	`
}

func testAccCheckIBMIsImageComponentDataSourceConfigByName() string {
	return `
		data "ibm_is_image_components" "components" {}

		data "ibm_is_image_component" "is_image_component_instance" {
			name = data.ibm_is_image_components.components.image_components.0.name
		}
	`
}
