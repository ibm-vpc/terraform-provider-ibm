// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"errors"
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIBMIsImageRecipeDeprecateBasic(t *testing.T) {
	name := fmt.Sprintf("tf-recipe-dep-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageRecipeDeprecateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageRecipeDeprecateConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMIsImageRecipeDeprecateExists("ibm_is_image_recipe.is_image_recipe_instance"),
					resource.TestCheckResourceAttr("ibm_is_image_recipe_deprecate.deprecate", "status", "deprecated"),
					resource.TestCheckResourceAttrSet("ibm_is_image_recipe_deprecate.deprecate", "deprecation_at"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageRecipeDeprecateDestroy(s *terraform.State) error {
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
			return fmt.Errorf("error checking for ImageRecipe (%s) destruction: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccCheckIBMIsImageRecipeDeprecateExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getImageRecipeOptions := &vpcv1.GetImageRecipeOptions{}
		getImageRecipeOptions.SetID(rs.Primary.ID)

		_, _, err = vpcClient.GetImageRecipe(getImageRecipeOptions)
		return err
	}
}

func testAccCheckIBMIsImageRecipeDeprecateConfig(name string) string {
	return fmt.Sprintf(`
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
			name             = "%s"
		}

		resource "ibm_is_image_recipe_deprecate" "deprecate" {
			image_recipe = ibm_is_image_recipe.is_image_recipe_instance.id
		}
	`, name)
}
