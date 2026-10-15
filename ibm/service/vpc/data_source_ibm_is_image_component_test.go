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

func TestAccIBMIsImageComponentDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageComponentDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "is_image_component_id"),
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
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "status_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_component.is_image_component_instance", "visibility"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageComponentDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_image_component" "is_image_component_instance" {
			id = "id"
		}
	`)
}


func TestDataSourceIBMIsImageComponentImageComponentAllowedUseToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["api_version"] = "2024-06-23"
		model["image"] = "operating_system.architecture == amd64"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageComponentAllowedUse)
	model.ApiVersion = core.StringPtr("2024-06-23")
	model.Image = core.StringPtr("operating_system.architecture == amd64")

	result, err := vpc.DataSourceIBMIsImageComponentImageComponentAllowedUseToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentImageComponentReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		model["name"] = "my-image-component"
		model["resource_type"] = "image_component"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ImageComponentReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	model.Name = core.StringPtr("my-image-component")
	model.ResourceType = core.StringPtr("image_component")

	result, err := vpc.DataSourceIBMIsImageComponentImageComponentReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsImageComponentDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentImageComponentResourceGroupToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
		model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"
		model["name"] = "my-resource-group"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageComponentResourceGroup)
	model.Href = core.StringPtr("https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345")
	model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")
	model.Name = core.StringPtr("my-resource-group")

	result, err := vpc.DataSourceIBMIsImageComponentImageComponentResourceGroupToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentImageComponentSoftwareToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["name"] = "cuda-toolkit"
		model["resource_type"] = "software_package"
		model["version"] = "12.4.0"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageComponentSoftware)
	model.Name = core.StringPtr("cuda-toolkit")
	model.ResourceType = core.StringPtr("software_package")
	model.Version = core.StringPtr("12.4.0")

	result, err := vpc.DataSourceIBMIsImageComponentImageComponentSoftwareToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentImageComponentSoftwarePackageToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["name"] = "cuda-toolkit"
		model["resource_type"] = "software_package"
		model["version"] = "12.4.0"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageComponentSoftwarePackage)
	model.Name = core.StringPtr("cuda-toolkit")
	model.ResourceType = core.StringPtr("software_package")
	model.Version = core.StringPtr("12.4.0")

	result, err := vpc.DataSourceIBMIsImageComponentImageComponentSoftwarePackageToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentImageComponentStatusReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "internal_error"
		model["message"] = "internal error (contact IBM support)"
		model["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageComponentStatusReason)
	model.Code = core.StringPtr("internal_error")
	model.Message = core.StringPtr("internal error (contact IBM support)")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons")

	result, err := vpc.DataSourceIBMIsImageComponentImageComponentStatusReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
