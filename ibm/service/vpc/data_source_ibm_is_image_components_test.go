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
	. "github.com/IBM-Cloud/terraform-provider-ibm/ibm/unittest"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageComponentsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageComponentsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_components.is_image_components_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_components.is_image_components_instance", "is_image_components_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_components.is_image_components_instance", "image_components.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageComponentsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_image_components" "is_image_components_instance" {
			id = "id"
		}
	`)
}


func TestDataSourceIBMIsImageComponentsImageComponentToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		imageComponentAllowedUseModel := make(map[string]interface{})
		imageComponentAllowedUseModel["api_version"] = "2024-06-23"
		imageComponentAllowedUseModel["image"] = "operating_system.architecture == amd64"

		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		imageComponentReferenceModel := make(map[string]interface{})
		imageComponentReferenceModel["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		imageComponentReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		imageComponentReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		imageComponentReferenceModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		imageComponentReferenceModel["name"] = "my-image-component"
		imageComponentReferenceModel["resource_type"] = "image_component"

		imageComponentResourceGroupModel := make(map[string]interface{})
		imageComponentResourceGroupModel["href"] = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
		imageComponentResourceGroupModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"
		imageComponentResourceGroupModel["name"] = "my-resource-group"

		imageComponentSoftwareModel := make(map[string]interface{})
		imageComponentSoftwareModel["name"] = "cuda-toolkit"
		imageComponentSoftwareModel["resource_type"] = "software_package"
		imageComponentSoftwareModel["version"] = "12.4.0"

		imageComponentStatusReasonModel := make(map[string]interface{})
		imageComponentStatusReasonModel["code"] = "internal_error"
		imageComponentStatusReasonModel["message"] = "internal error (contact IBM support)"
		imageComponentStatusReasonModel["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons"

		model := make(map[string]interface{})
		model["allowed_use"] = []map[string]interface{}{imageComponentAllowedUseModel}
		model["created_at"] = "2026-01-02T03:04:05.006Z"
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		model["dependencies"] = []map[string]interface{}{imageComponentReferenceModel}
		model["deprecation_at"] = "2026-01-02T03:04:05.006Z"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		model["name"] = "my-image-component"
		model["obsolescence_at"] = "2026-01-02T03:04:05.006Z"
		model["resource_group"] = []map[string]interface{}{imageComponentResourceGroupModel}
		model["resource_type"] = "image_component"
		model["software"] = []map[string]interface{}{imageComponentSoftwareModel}
		model["status"] = "available"
		model["status_reasons"] = []map[string]interface{}{imageComponentStatusReasonModel}
		model["visibility"] = "private"

		assert.Equal(t, result, model)
	}

	imageComponentAllowedUseModel := new(vpcv1.ImageComponentAllowedUse)
	imageComponentAllowedUseModel.ApiVersion = core.StringPtr("2024-06-23")
	imageComponentAllowedUseModel.Image = core.StringPtr("operating_system.architecture == amd64")

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	imageComponentReferenceModel := new(vpcv1.ImageComponentReference)
	imageComponentReferenceModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	imageComponentReferenceModel.Deleted = deletedModel
	imageComponentReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	imageComponentReferenceModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	imageComponentReferenceModel.Name = core.StringPtr("my-image-component")
	imageComponentReferenceModel.ResourceType = core.StringPtr("image_component")

	imageComponentResourceGroupModel := new(vpcv1.ImageComponentResourceGroup)
	imageComponentResourceGroupModel.Href = core.StringPtr("https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345")
	imageComponentResourceGroupModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")
	imageComponentResourceGroupModel.Name = core.StringPtr("my-resource-group")

	imageComponentSoftwareModel := new(vpcv1.ImageComponentSoftwarePackage)
	imageComponentSoftwareModel.Name = core.StringPtr("cuda-toolkit")
	imageComponentSoftwareModel.ResourceType = core.StringPtr("software_package")
	imageComponentSoftwareModel.Version = core.StringPtr("12.4.0")

	imageComponentStatusReasonModel := new(vpcv1.ImageComponentStatusReason)
	imageComponentStatusReasonModel.Code = core.StringPtr("internal_error")
	imageComponentStatusReasonModel.Message = core.StringPtr("internal error (contact IBM support)")
	imageComponentStatusReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons")

	model := new(vpcv1.ImageComponent)
	model.AllowedUse = imageComponentAllowedUseModel
	model.CreatedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	model.Dependencies = []vpcv1.ImageComponentReference{*imageComponentReferenceModel}
	model.DeprecationAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	model.Name = core.StringPtr("my-image-component")
	model.ObsolescenceAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.ResourceGroup = imageComponentResourceGroupModel
	model.ResourceType = core.StringPtr("image_component")
	model.Software = imageComponentSoftwareModel
	model.Status = core.StringPtr("available")
	model.StatusReasons = []vpcv1.ImageComponentStatusReason{*imageComponentStatusReasonModel}
	model.Visibility = core.StringPtr("private")

	result, err := vpc.DataSourceIBMIsImageComponentsImageComponentToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentsImageComponentAllowedUseToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["api_version"] = "2024-06-23"
		model["image"] = "operating_system.architecture == amd64"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageComponentAllowedUse)
	model.ApiVersion = core.StringPtr("2024-06-23")
	model.Image = core.StringPtr("operating_system.architecture == amd64")

	result, err := vpc.DataSourceIBMIsImageComponentsImageComponentAllowedUseToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentsImageComponentReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageComponentsImageComponentReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentsDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsImageComponentsDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentsImageComponentResourceGroupToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageComponentsImageComponentResourceGroupToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentsImageComponentSoftwareToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageComponentsImageComponentSoftwareToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentsImageComponentSoftwarePackageToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageComponentsImageComponentSoftwarePackageToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageComponentsImageComponentStatusReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageComponentsImageComponentStatusReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
