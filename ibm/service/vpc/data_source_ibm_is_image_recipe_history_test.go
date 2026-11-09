// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageRecipeHistoryDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageRecipeHistoryDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "image_recipe_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "history.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "history.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "history.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "history.0.version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "history.0.status"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageRecipeHistoryDataSourceConfigBasic() string {
	return `
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
		}

		data "ibm_is_image_recipe_history" "is_image_recipe_history_instance" {
			image_recipe_id = ibm_is_image_recipe.is_image_recipe_instance.id
		}
	`
}
