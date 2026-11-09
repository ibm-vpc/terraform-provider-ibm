// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsImageRecipeBasic(t *testing.T) {
	var conf vpcv1.ImageRecipe

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageRecipeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageRecipeConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageRecipeExists("ibm_is_image_recipe.is_image_recipe_instance", conf),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "id"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "crn"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "href"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "base_image"),
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "image_components.#", "1"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "status"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "version"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "visibility"),
				),
			},
		},
	})
}

func TestAccIBMIsImageRecipeAllArgs(t *testing.T) {
	var conf vpcv1.ImageRecipe
	name := fmt.Sprintf("tf-recipe-%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf-recipe-up-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageRecipeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageRecipeConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageRecipeExists("ibm_is_image_recipe.is_image_recipe_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "deprecation_at", "2029-06-01T00:00:00.000Z"),
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "obsolescence_at", "2029-12-01T00:00:00.000Z"),
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "image_components.#", "1"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "resource_group"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "crn"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe.is_image_recipe_instance", "href"),
				),
			},
			{
				Config: testAccCheckIBMIsImageRecipeConfig(nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "name", nameUpdate),
				),
			},
			{
				ResourceName:      "ibm_is_image_recipe.is_image_recipe_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsImageRecipeConfigBasic() string {
	return `
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
		}
	`
}

func testAccCheckIBMIsImageRecipeConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
			name             = "%s"
			deprecation_at   = "2029-06-01T00:00:00.000Z"
			obsolescence_at  = "2029-12-01T00:00:00.000Z"
			resource_group   = "fee82deba12e4c0fb69c3b09d1f12345"
		}
	`, name)
}

func testAccCheckIBMIsImageRecipeExists(n string, obj vpcv1.ImageRecipe) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getImageRecipeOptions := &vpcv1.GetImageRecipeOptions{}
		getImageRecipeOptions.SetID(rs.Primary.ID)

		imageRecipe, _, err := vpcClient.GetImageRecipe(getImageRecipeOptions)
		if err != nil {
			return err
		}

		obj = *imageRecipe
		return nil
	}
}

func testAccCheckIBMIsImageRecipeDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_image_recipe" {
			continue
		}

		getImageRecipeOptions := &vpcv1.GetImageRecipeOptions{}
		getImageRecipeOptions.SetID(rs.Primary.ID)

		_, response, err := vpcClient.GetImageRecipe(getImageRecipeOptions)
		if err == nil {
			return fmt.Errorf("ImageRecipe still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("error checking for ImageRecipe (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}
	return nil
}
