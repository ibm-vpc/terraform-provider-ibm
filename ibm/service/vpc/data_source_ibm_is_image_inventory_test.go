// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.117.0-7f07c563-20260915-094553
*/

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageInventoryDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageInventoryDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_inventory.is_image_inventory_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_inventory.is_image_inventory_instance", "is_image_inventory_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_inventory.is_image_inventory_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_inventory.is_image_inventory_instance", "packages.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_inventory.is_image_inventory_instance", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageInventoryDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_image_inventory" "is_image_inventory_instance" {
			id = "id"
		}
	`)
}


func TestDataSourceIBMIsImageInventoryImageInventoryPackageToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["name"] = "cuda-toolkit-12-4"
		model["version"] = "12.4.0-1"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageInventoryPackage)
	model.Name = core.StringPtr("cuda-toolkit-12-4")
	model.Version = core.StringPtr("12.4.0-1")

	result, err := vpc.DataSourceIBMIsImageInventoryImageInventoryPackageToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
