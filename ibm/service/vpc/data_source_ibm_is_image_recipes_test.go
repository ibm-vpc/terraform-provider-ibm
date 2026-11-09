// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageRecipesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageRecipesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.#"),
				),
			},
		},
	})
}

func TestAccIBMIsImageRecipesDataSourceAllArgs(t *testing.T) {
	imageRecipeName := fmt.Sprintf("tf-recipes-ds-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageRecipesDataSourceConfig(imageRecipeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "id"),
					resource.TestCheckResourceAttr("data.ibm_is_image_recipes.is_image_recipes_instance", "name", imageRecipeName),
					resource.TestCheckResourceAttr("data.ibm_is_image_recipes.is_image_recipes_instance", "visibility", "private"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.deprecation_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.name", imageRecipeName),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.obsolescence_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.updated_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipes.is_image_recipes_instance", "image_recipes.0.visibility"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageRecipesDataSourceConfigBasic() string {
	return `
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
		}

		data "ibm_is_image_recipes" "is_image_recipes_instance" {
		}
	`
}

func testAccCheckIBMIsImageRecipesDataSourceConfig(imageRecipeName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
			name             = "%s"
			deprecation_at   = "2029-06-01T00:00:00.000Z"
			obsolescence_at  = "2029-12-01T00:00:00.000Z"
			resource_group   = "fee82deba12e4c0fb69c3b09d1f12345"
		}

		data "ibm_is_image_recipes" "is_image_recipes_instance" {
			name       = ibm_is_image_recipe.is_image_recipe_instance.name
			visibility = "private"
		}
	`, imageRecipeName)
}
