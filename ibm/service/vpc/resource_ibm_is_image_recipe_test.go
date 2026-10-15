// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageRecipeBasic(t *testing.T) {
	var conf vpcv1.ImageRecipe

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageRecipeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageRecipeConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageRecipeExists("ibm_is_image_recipe.is_image_recipe_instance", conf),
				),
			},
		},
	})
}

func TestAccIBMIsImageRecipeAllArgs(t *testing.T) {
	var conf vpcv1.ImageRecipe
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageRecipeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageRecipeConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageRecipeExists("ibm_is_image_recipe.is_image_recipe_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsImageRecipeConfig(nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_image_recipe.is_image_recipe_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_image_recipe.is_image_recipe_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsImageRecipeConfigBasic() string {
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
	`, )
}

func testAccCheckIBMIsImageRecipeConfig(name string) string {
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

		// Try to find the key
		_, response, err := vpcClient.GetImageRecipe(getImageRecipeOptions)

		if err == nil {
			return fmt.Errorf("ImageRecipe still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("error checking for ImageRecipe (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsImageRecipeImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsImageRecipeImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.ResourceIBMIsImageRecipeDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeAccountReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "bb1b52262f7441a586f49068482f1e60"
		model["resource_type"] = "account"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.AccountReference)
	model.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	model.ResourceType = core.StringPtr("account")

	result, err := vpc.ResourceIBMIsImageRecipeAccountReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeRegionReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		model["name"] = "us-south"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.RegionReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	model.Name = core.StringPtr("us-south")

	result, err := vpc.ResourceIBMIsImageRecipeRegionReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeImageComponentReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsImageRecipeImageComponentReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeImageRecipeResourceGroupToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsImageRecipeImageRecipeResourceGroupToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeImageRecipeStatusReasonToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsImageRecipeImageRecipeStatusReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototype(t *testing.T) {
	checkResult := func(result vpcv1.ImageRecipeBaseImagePrototypeIntf) {
		model := new(vpcv1.ImageRecipeBaseImagePrototype)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByID) {
		model := new(vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByID)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByCRN) {
		model := new(vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByHref) {
		model := new(vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageComponentIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ImageComponentIdentityIntf) {
		model := new(vpcv1.ImageComponentIdentity)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageComponentIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByID(t *testing.T) {
	checkResult := func(result *vpcv1.ImageComponentIdentityImageComponentByID) {
		model := new(vpcv1.ImageComponentIdentityImageComponentByID)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.ImageComponentIdentityImageComponentByCRN) {
		model := new(vpcv1.ImageComponentIdentityImageComponentByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image-component:r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ImageComponentIdentityImageComponentByHref) {
		model := new(vpcv1.ImageComponentIdentityImageComponentByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_components/r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeResourceGroupPrototype(t *testing.T) {
	checkResult := func(result vpcv1.ImageRecipeResourceGroupPrototypeIntf) {
		model := new(vpcv1.ImageRecipeResourceGroupPrototype)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeResourceGroupPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeResourceGroupPrototypeResourceGroupIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ImageRecipeResourceGroupPrototypeResourceGroupIdentityByID) {
		model := new(vpcv1.ImageRecipeResourceGroupPrototypeResourceGroupIdentityByID)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeResourceGroupPrototypeResourceGroupIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatch(t *testing.T) {
	checkResult := func(result vpcv1.ImageRecipeBaseImagePatchIntf) {
		model := new(vpcv1.ImageRecipeBaseImagePatch)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatch(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ImageRecipeBaseImagePatchImageIdentityByID) {
		model := new(vpcv1.ImageRecipeBaseImagePatchImageIdentityByID)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.ImageRecipeBaseImagePatchImageIdentityByCRN) {
		model := new(vpcv1.ImageRecipeBaseImagePatchImageIdentityByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ImageRecipeBaseImagePatchImageIdentityByHref) {
		model := new(vpcv1.ImageRecipeBaseImagePatchImageIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}
