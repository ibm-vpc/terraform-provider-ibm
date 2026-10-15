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

func TestAccIBMIsImageRecipeHistoryDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageRecipeHistoryDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "is_image_recipe_history_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe_history.is_image_recipe_history_instance", "history.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageRecipeHistoryDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_image_recipe_history" "is_image_recipe_history_instance" {
			id = "id"
		}
	`)
}


func TestDataSourceIBMIsImageRecipeHistoryImageRecipeHistoryToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		accountReferenceModel := make(map[string]interface{})
		accountReferenceModel["id"] = "bb1b52262f7441a586f49068482f1e60"
		accountReferenceModel["resource_type"] = "account"

		regionReferenceModel := make(map[string]interface{})
		regionReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		regionReferenceModel["name"] = "us-south"

		imageRemoteContextImageReferenceModel := make(map[string]interface{})
		imageRemoteContextImageReferenceModel["account"] = []map[string]interface{}{accountReferenceModel}
		imageRemoteContextImageReferenceModel["region"] = []map[string]interface{}{regionReferenceModel}

		imageReferenceModel := make(map[string]interface{})
		imageReferenceModel["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		imageReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		imageReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		imageReferenceModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		imageReferenceModel["name"] = "my-image"
		imageReferenceModel["remote"] = []map[string]interface{}{imageRemoteContextImageReferenceModel}
		imageReferenceModel["resource_type"] = "image"

		imageComponentReferenceModel := make(map[string]interface{})
		imageComponentReferenceModel["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		imageComponentReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		imageComponentReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		imageComponentReferenceModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
		imageComponentReferenceModel["name"] = "my-image-component"
		imageComponentReferenceModel["resource_type"] = "image_component"

		imageRecipeResourceGroupModel := make(map[string]interface{})
		imageRecipeResourceGroupModel["href"] = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
		imageRecipeResourceGroupModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"
		imageRecipeResourceGroupModel["name"] = "my-resource-group"

		imageRecipeStatusReasonModel := make(map[string]interface{})
		imageRecipeStatusReasonModel["code"] = "internal_error"
		imageRecipeStatusReasonModel["message"] = "internal error (contact IBM support)"
		imageRecipeStatusReasonModel["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons"

		model := make(map[string]interface{})
		model["base_image"] = []map[string]interface{}{imageReferenceModel}
		model["created_at"] = "2026-01-02T03:04:05.006Z"
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12"
		model["deprecation_at"] = "2026-01-02T03:04:05.006Z"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
		model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
		model["image_components"] = []map[string]interface{}{imageComponentReferenceModel}
		model["name"] = "my-image-recipe"
		model["obsolescence_at"] = "2026-01-02T03:04:05.006Z"
		model["resource_group"] = []map[string]interface{}{imageRecipeResourceGroupModel}
		model["resource_type"] = "image_recipe"
		model["status"] = "available"
		model["status_reasons"] = []map[string]interface{}{imageRecipeStatusReasonModel}
		model["superseded_at"] = "2026-01-02T03:04:05.006Z"
		model["updated_at"] = "2026-01-02T03:04:05.006Z"
		model["version"] = "a3f8c2d1"
		model["visibility"] = "private"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	accountReferenceModel := new(vpcv1.AccountReference)
	accountReferenceModel.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	accountReferenceModel.ResourceType = core.StringPtr("account")

	regionReferenceModel := new(vpcv1.RegionReference)
	regionReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	regionReferenceModel.Name = core.StringPtr("us-south")

	imageRemoteContextImageReferenceModel := new(vpcv1.ImageRemoteContextImageReference)
	imageRemoteContextImageReferenceModel.Account = accountReferenceModel
	imageRemoteContextImageReferenceModel.Region = regionReferenceModel

	imageReferenceModel := new(vpcv1.ImageReference)
	imageReferenceModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
	imageReferenceModel.Deleted = deletedModel
	imageReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
	imageReferenceModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
	imageReferenceModel.Name = core.StringPtr("my-image")
	imageReferenceModel.Remote = imageRemoteContextImageReferenceModel
	imageReferenceModel.ResourceType = core.StringPtr("image")

	imageComponentReferenceModel := new(vpcv1.ImageComponentReference)
	imageComponentReferenceModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	imageComponentReferenceModel.Deleted = deletedModel
	imageComponentReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	imageComponentReferenceModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
	imageComponentReferenceModel.Name = core.StringPtr("my-image-component")
	imageComponentReferenceModel.ResourceType = core.StringPtr("image_component")

	imageRecipeResourceGroupModel := new(vpcv1.ImageRecipeResourceGroup)
	imageRecipeResourceGroupModel.Href = core.StringPtr("https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345")
	imageRecipeResourceGroupModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")
	imageRecipeResourceGroupModel.Name = core.StringPtr("my-resource-group")

	imageRecipeStatusReasonModel := new(vpcv1.ImageRecipeStatusReason)
	imageRecipeStatusReasonModel.Code = core.StringPtr("internal_error")
	imageRecipeStatusReasonModel.Message = core.StringPtr("internal error (contact IBM support)")
	imageRecipeStatusReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons")

	model := new(vpcv1.ImageRecipeHistory)
	model.BaseImage = imageReferenceModel
	model.CreatedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12")
	model.DeprecationAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")
	model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")
	model.ImageComponents = []vpcv1.ImageComponentReference{*imageComponentReferenceModel}
	model.Name = core.StringPtr("my-image-recipe")
	model.ObsolescenceAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.ResourceGroup = imageRecipeResourceGroupModel
	model.ResourceType = core.StringPtr("image_recipe")
	model.Status = core.StringPtr("available")
	model.StatusReasons = []vpcv1.ImageRecipeStatusReason{*imageRecipeStatusReasonModel}
	model.SupersededAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.UpdatedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.Version = core.StringPtr("a3f8c2d1")
	model.Visibility = core.StringPtr("private")

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryImageRecipeHistoryToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryImageReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		accountReferenceModel := make(map[string]interface{})
		accountReferenceModel["id"] = "bb1b52262f7441a586f49068482f1e60"
		accountReferenceModel["resource_type"] = "account"

		regionReferenceModel := make(map[string]interface{})
		regionReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		regionReferenceModel["name"] = "us-south"

		imageRemoteContextImageReferenceModel := make(map[string]interface{})
		imageRemoteContextImageReferenceModel["account"] = []map[string]interface{}{accountReferenceModel}
		imageRemoteContextImageReferenceModel["region"] = []map[string]interface{}{regionReferenceModel}

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
		model["name"] = "my-image"
		model["remote"] = []map[string]interface{}{imageRemoteContextImageReferenceModel}
		model["resource_type"] = "image"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	accountReferenceModel := new(vpcv1.AccountReference)
	accountReferenceModel.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	accountReferenceModel.ResourceType = core.StringPtr("account")

	regionReferenceModel := new(vpcv1.RegionReference)
	regionReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	regionReferenceModel.Name = core.StringPtr("us-south")

	imageRemoteContextImageReferenceModel := new(vpcv1.ImageRemoteContextImageReference)
	imageRemoteContextImageReferenceModel.Account = accountReferenceModel
	imageRemoteContextImageReferenceModel.Region = regionReferenceModel

	model := new(vpcv1.ImageReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
	model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
	model.Name = core.StringPtr("my-image")
	model.Remote = imageRemoteContextImageReferenceModel
	model.ResourceType = core.StringPtr("image")

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryImageRemoteContextImageReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		accountReferenceModel := make(map[string]interface{})
		accountReferenceModel["id"] = "bb1b52262f7441a586f49068482f1e60"
		accountReferenceModel["resource_type"] = "account"

		regionReferenceModel := make(map[string]interface{})
		regionReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		regionReferenceModel["name"] = "us-south"

		model := make(map[string]interface{})
		model["account"] = []map[string]interface{}{accountReferenceModel}
		model["region"] = []map[string]interface{}{regionReferenceModel}

		assert.Equal(t, result, model)
	}

	accountReferenceModel := new(vpcv1.AccountReference)
	accountReferenceModel.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	accountReferenceModel.ResourceType = core.StringPtr("account")

	regionReferenceModel := new(vpcv1.RegionReference)
	regionReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	regionReferenceModel.Name = core.StringPtr("us-south")

	model := new(vpcv1.ImageRemoteContextImageReference)
	model.Account = accountReferenceModel
	model.Region = regionReferenceModel

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryImageRemoteContextImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryAccountReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "bb1b52262f7441a586f49068482f1e60"
		model["resource_type"] = "account"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.AccountReference)
	model.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	model.ResourceType = core.StringPtr("account")

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryAccountReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryRegionReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		model["name"] = "us-south"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.RegionReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	model.Name = core.StringPtr("us-south")

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryRegionReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryImageComponentReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryImageComponentReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryImageRecipeResourceGroupToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
		model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"
		model["name"] = "my-resource-group"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageRecipeResourceGroup)
	model.Href = core.StringPtr("https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345")
	model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")
	model.Name = core.StringPtr("my-resource-group")

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryImageRecipeResourceGroupToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeHistoryImageRecipeStatusReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "internal_error"
		model["message"] = "internal error (contact IBM support)"
		model["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageRecipeStatusReason)
	model.Code = core.StringPtr("internal_error")
	model.Message = core.StringPtr("internal error (contact IBM support)")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-image-recipe-status-reasons")

	result, err := vpc.DataSourceIBMIsImageRecipeHistoryImageRecipeStatusReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
