// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageInstanceProfilesDataSourceBasic(t *testing.T) {
	name := fmt.Sprintf("tfimg-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageInstanceProfilesDataSourceConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_instance_profiles.is_image_instance_profiles_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_instance_profiles.is_image_instance_profiles_instance", "instance_profiles.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_instance_profiles.is_image_instance_profiles_instance", "instance_profiles.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_instance_profiles.is_image_instance_profiles_instance", "instance_profiles.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_instance_profiles.is_image_instance_profiles_instance", "instance_profiles.0.resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageInstanceProfilesDataSourceConfigBasic(name string) string {
	return testAccCheckIBMISImageConfig(name) + fmt.Sprintf(` 
		data "ibm_is_image_instance_profiles" "is_image_instance_profiles_instance" {
			identifier = "ibm_is_image.isExampleImage.id"
		}
	`)
}
