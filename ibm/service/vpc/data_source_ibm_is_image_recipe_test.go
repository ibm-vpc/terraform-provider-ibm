// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.117.0-7f07c563-20260915-094553
*/

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageRecipeDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageRecipeDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "is_image_recipe_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "base_image.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "image_components.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "resource_group.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "status_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "updated_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "visibility"),
				),
			},
		},
	})
}

func TestAccIBMIsImageRecipeDataSourceAllArgs(t *testing.T) {
	imageRecipeName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageRecipeDataSourceConfig(imageRecipeName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "is_image_recipe_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "base_image.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "deprecation_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "image_components.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "image_components.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "image_components.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "image_components.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_image_recipe.is_image_recipe_instance", "image_components.0.name", imageRecipeName),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "image_components.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "obsolescence_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "resource_group.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "status_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "status_reasons.0.code"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "status_reasons.0.message"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "status_reasons.0.more_info"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "superseded_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "updated_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_recipe.is_image_recipe_instance", "visibility"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageRecipeDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image {
				crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
				href = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
				id = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			}
			image_components {
				crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
				id = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
				name = "my-image-component"
				resource_type = "image_component"
			}
		}

		data "ibm_is_image_recipe" "is_image_recipe_instance" {
			is_image_recipe_id = "is_image_recipe_id"
		}
	`)
}

func testAccCheckIBMIsImageRecipeDataSourceConfig(imageRecipeName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_image_recipe" "is_image_recipe_instance" {
			base_image {
				crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
				href = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
				id = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			}
			deprecation_at = "2004-10-28T04:39:00.000Z"
			image_components {
				crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
				id = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
				name = "my-image-component"
				resource_type = "image_component"
			}
			name = "%s"
			obsolescence_at = "2004-10-28T04:39:00.000Z"
			resource_group {
				id = "fee82deba12e4c0fb69c3b09d1f12345"
			}
		}

		data "ibm_is_image_recipe" "is_image_recipe_instance" {
			is_image_recipe_id = "is_image_recipe_id"
		}
	`, imageRecipeName)
}

func TestDataSourceIBMIsImageRecipeImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageRecipeImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsImageRecipeDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeAccountReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "bb1b52262f7441a586f49068482f1e60"
		model["resource_type"] = "account"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.AccountReference)
	model.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	model.ResourceType = core.StringPtr("account")

	result, err := vpc.DataSourceIBMIsImageRecipeAccountReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeRegionReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		model["name"] = "us-south"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.RegionReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	model.Name = core.StringPtr("us-south")

	result, err := vpc.DataSourceIBMIsImageRecipeRegionReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeImageComponentReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageRecipeImageComponentReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeImageRecipeResourceGroupToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageRecipeImageRecipeResourceGroupToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageRecipeImageRecipeStatusReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageRecipeImageRecipeStatusReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
