// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsImagePackagesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImagePackagesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_packages.is_image_packages_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_packages.is_image_packages_instance", "image_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_packages.is_image_packages_instance", "packages.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImagePackagesDataSourceConfigBasic() string {
	return `
		data "ibm_is_images" "ds_images" {
			visibility = "private"
		}

		data "ibm_is_image_packages" "is_image_packages_instance" {
			image_id = data.ibm_is_images.ds_images.images.0.id
		}
	`
}

func TestDataSourceIBMIsImagePackagesImagePackageToMap(t *testing.T) {
	model := new(vpcv1.ImagePackage)
	model.Name = core.StringPtr("cuda-toolkit-12-4")
	model.Version = core.StringPtr("12.4.0-1")

	result, err := vpc.DataSourceIBMIsImagePackagesImagePackageToMap(model)
	assert.Nil(t, err)

	expected := map[string]interface{}{
		"name":    "cuda-toolkit-12-4",
		"version": "12.4.0-1",
	}
	assert.Equal(t, expected, result)
}
