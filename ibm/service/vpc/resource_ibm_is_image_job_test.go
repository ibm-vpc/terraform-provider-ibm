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

func TestAccIBMIsImageJobBasic(t *testing.T) {
	var conf vpcv1.ImageJob
	typeVar := "build"
	typeVarUpdate := "build"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageJobDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageJobConfigBasic(typeVar),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageJobExists("ibm_is_image_job.is_image_job_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVar),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsImageJobConfigBasic(typeVarUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVarUpdate),
				),
			},
		},
	})
}

func TestAccIBMIsImageJobAllArgs(t *testing.T) {
	var conf vpcv1.ImageJob
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	typeVar := "build"
	deleteImageJobOnImageDelete := "true"
	timeout := fmt.Sprintf("%d", acctest.RandIntRange(30, 480))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	typeVarUpdate := "build"
	deleteImageJobOnImageDeleteUpdate := "false"
	timeoutUpdate := fmt.Sprintf("%d", acctest.RandIntRange(30, 480))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageJobDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageJobConfig(name, typeVar, deleteImageJobOnImageDelete, timeout),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageJobExists("ibm_is_image_job.is_image_job_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVar),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "delete_image_job_on_image_delete", deleteImageJobOnImageDelete),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "timeout", timeout),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsImageJobConfig(nameUpdate, typeVarUpdate, deleteImageJobOnImageDeleteUpdate, timeoutUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVarUpdate),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "delete_image_job_on_image_delete", deleteImageJobOnImageDeleteUpdate),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "timeout", timeoutUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_image_job.is_image_job_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsImageJobConfigBasic(typeVar string) string {
	return fmt.Sprintf(`
		resource "ibm_is_image_job" "is_image_job_instance" {
			type = "%s"
		}
	`, typeVar)
}

func testAccCheckIBMIsImageJobConfig(name string, typeVar string, deleteImageJobOnImageDelete string, timeout string) string {
	return fmt.Sprintf(`

		resource "ibm_is_image_job" "is_image_job_instance" {
			name = "%s"
			type = "%s"
			delete_image_job_on_image_delete = %s
			image_recipe {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12"
				href = "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8/history/a3f8c2d1"
				id = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
				version = "a3f8c2d1"
			}
			image_template {
				name = "my-image"
				resource_group {
					id = "fee82deba12e4c0fb69c3b09d1f12345"
				}
			}
			timeout = %s
			worker {
				profile {
					href = "https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"
					name = "bx2-4x16"
				}
				subnet {
					crn = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
					href = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
					id = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				}
			}
		}
	`, name, typeVar, deleteImageJobOnImageDelete, timeout)
}

func testAccCheckIBMIsImageJobExists(n string, obj vpcv1.ImageJob) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getImageJobOptions := &vpcv1.GetImageJobOptions{}

		getImageJobOptions.SetID(rs.Primary.ID)

		imageJobIntf, _, err := vpcClient.GetImageJob(getImageJobOptions)
		if err != nil {
			return err
		}

		imageJob := imageJobIntf.(*vpcv1.ImageJob)
		obj = *imageJob
		return nil
	}
}

func testAccCheckIBMIsImageJobDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_image_job" {
			continue
		}

		getImageJobOptions := &vpcv1.GetImageJobOptions{}

		getImageJobOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := vpcClient.GetImageJob(getImageJobOptions)

		if err == nil {
			return fmt.Errorf("is_image_job still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("error checking for is_image_job (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsImageJobImageJobBuildImageRecipeToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8/history/a3f8c2d1"
		model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
		model["name"] = "my-image-recipe"
		model["resource_type"] = "image_recipe"
		model["version"] = "a3f8c2d1"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ImageJobBuildImageRecipe)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8/history/a3f8c2d1")
	model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")
	model.Name = core.StringPtr("my-image-recipe")
	model.ResourceType = core.StringPtr("image_recipe")
	model.Version = core.StringPtr("a3f8c2d1")

	result, err := vpc.ResourceIBMIsImageJobImageJobBuildImageRecipeToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.ResourceIBMIsImageJobDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobImageJobBuildImageTemplateToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		resourceGroupIdentityModel := make(map[string]interface{})
		resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

		model := make(map[string]interface{})
		model["name"] = "my-image"
		model["resource_group"] = []map[string]interface{}{resourceGroupIdentityModel}

		assert.Equal(t, result, model)
	}

	resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
	resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

	model := new(vpcv1.ImageJobBuildImageTemplate)
	model.Name = core.StringPtr("my-image")
	model.ResourceGroup = resourceGroupIdentityModel

	result, err := vpc.ResourceIBMIsImageJobImageJobBuildImageTemplateToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobResourceGroupIdentityToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ResourceGroupIdentity)
	model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

	result, err := vpc.ResourceIBMIsImageJobResourceGroupIdentityToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobResourceGroupIdentityResourceGroupIdentityByIDToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
	model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

	result, err := vpc.ResourceIBMIsImageJobResourceGroupIdentityResourceGroupIdentityByIDToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobImageJobBuildWorkerToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		instanceReferenceModel := make(map[string]interface{})
		instanceReferenceModel["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"
		instanceReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		instanceReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"
		instanceReferenceModel["id"] = "0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"
		instanceReferenceModel["name"] = "my-instance"

		imageJobBuildWorkerProfileModel := make(map[string]interface{})
		imageJobBuildWorkerProfileModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"
		imageJobBuildWorkerProfileModel["name"] = "bx2-4x16"
		imageJobBuildWorkerProfileModel["resource_type"] = "instance_profile"

		subnetReferenceModel := make(map[string]interface{})
		subnetReferenceModel["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		subnetReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		subnetReferenceModel["name"] = "my-subnet"
		subnetReferenceModel["resource_type"] = "subnet"

		model := make(map[string]interface{})
		model["instance"] = []map[string]interface{}{instanceReferenceModel}
		model["profile"] = []map[string]interface{}{imageJobBuildWorkerProfileModel}
		model["subnet"] = []map[string]interface{}{subnetReferenceModel}

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	instanceReferenceModel := new(vpcv1.InstanceReference)
	instanceReferenceModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0")
	instanceReferenceModel.Deleted = deletedModel
	instanceReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0")
	instanceReferenceModel.ID = core.StringPtr("0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0")
	instanceReferenceModel.Name = core.StringPtr("my-instance")

	imageJobBuildWorkerProfileModel := new(vpcv1.ImageJobBuildWorkerProfile)
	imageJobBuildWorkerProfileModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16")
	imageJobBuildWorkerProfileModel.Name = core.StringPtr("bx2-4x16")
	imageJobBuildWorkerProfileModel.ResourceType = core.StringPtr("instance_profile")

	subnetReferenceModel := new(vpcv1.SubnetReference)
	subnetReferenceModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.Deleted = deletedModel
	subnetReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	subnetReferenceModel.Name = core.StringPtr("my-subnet")
	subnetReferenceModel.ResourceType = core.StringPtr("subnet")

	model := new(vpcv1.ImageJobBuildWorker)
	model.Instance = instanceReferenceModel
	model.Profile = imageJobBuildWorkerProfileModel
	model.Subnet = subnetReferenceModel

	result, err := vpc.ResourceIBMIsImageJobImageJobBuildWorkerToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobInstanceReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"
		model["id"] = "0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0"
		model["name"] = "my-instance"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.InstanceReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::instance:0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0")
	model.ID = core.StringPtr("0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0")
	model.Name = core.StringPtr("my-instance")

	result, err := vpc.ResourceIBMIsImageJobInstanceReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobImageJobBuildWorkerProfileToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"
		model["name"] = "bx2-4x16"
		model["resource_type"] = "instance_profile"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageJobBuildWorkerProfile)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16")
	model.Name = core.StringPtr("bx2-4x16")
	model.ResourceType = core.StringPtr("instance_profile")

	result, err := vpc.ResourceIBMIsImageJobImageJobBuildWorkerProfileToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobSubnetReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		model["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
		model["name"] = "my-subnet"
		model["resource_type"] = "subnet"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.SubnetReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	model.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
	model.Name = core.StringPtr("my-subnet")
	model.ResourceType = core.StringPtr("subnet")

	result, err := vpc.ResourceIBMIsImageJobSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobImageJobStatusReasonsToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "internal_error"
		model["message"] = "A failure occurred"
		model["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-image-builder"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ImageJobStatusReasons)
	model.Code = core.StringPtr("internal_error")
	model.Message = core.StringPtr("A failure occurred")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-image-builder")

	result, err := vpc.ResourceIBMIsImageJobImageJobStatusReasonsToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsImageJobImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobImageRemoteContextImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsImageJobImageRemoteContextImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobAccountReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "bb1b52262f7441a586f49068482f1e60"
		model["resource_type"] = "account"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.AccountReference)
	model.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	model.ResourceType = core.StringPtr("account")

	result, err := vpc.ResourceIBMIsImageJobAccountReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobRegionReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		model["name"] = "us-south"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.RegionReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	model.Name = core.StringPtr("us-south")

	result, err := vpc.ResourceIBMIsImageJobRegionReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototype(t *testing.T) {
	checkResult := func(result vpcv1.ImageJobBuildImageRecipePrototypeIntf) {
		model := new(vpcv1.ImageJobBuildImageRecipePrototype)
		model.Version = core.StringPtr("a3f8c2d1")
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["version"] = "a3f8c2d1"
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
	model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByID(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByID) {
		model := new(vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByID)
		model.Version = core.StringPtr("a3f8c2d1")
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["version"] = "a3f8c2d1"
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByCRN) {
		model := new(vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByCRN)
		model.Version = core.StringPtr("a3f8c2d1")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["version"] = "a3f8c2d1"
	model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByHref) {
		model := new(vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByHref)
		model.Version = core.StringPtr("a3f8c2d1")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["version"] = "a3f8c2d1"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildImageTemplate(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildImageTemplate) {
		resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
		resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		model := new(vpcv1.ImageJobBuildImageTemplate)
		model.Name = core.StringPtr("my-image")
		model.ResourceGroup = resourceGroupIdentityModel

		assert.Equal(t, result, model)
	}

	resourceGroupIdentityModel := make(map[string]interface{})
	resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	model := make(map[string]interface{})
	model["name"] = "my-image"
	model["resource_group"] = []interface{}{resourceGroupIdentityModel}

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildImageTemplate(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToResourceGroupIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ResourceGroupIdentityIntf) {
		model := new(vpcv1.ResourceGroupIdentity)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMIsImageJobMapToResourceGroupIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToResourceGroupIdentityResourceGroupIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ResourceGroupIdentityResourceGroupIdentityByID) {
		model := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMIsImageJobMapToResourceGroupIdentityResourceGroupIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerPrototype(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildWorkerPrototype) {
		imageJobBuildWorkerProfilePrototypeModel := new(vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName)
		imageJobBuildWorkerProfilePrototypeModel.Name = core.StringPtr("bx2-4x16")

		imageJobBuildWorkerSubnetPrototypeModel := new(vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID)
		imageJobBuildWorkerSubnetPrototypeModel.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		model := new(vpcv1.ImageJobBuildWorkerPrototype)
		model.Profile = imageJobBuildWorkerProfilePrototypeModel
		model.Subnet = imageJobBuildWorkerSubnetPrototypeModel

		assert.Equal(t, result, model)
	}

	imageJobBuildWorkerProfilePrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerProfilePrototypeModel["name"] = "bx2-4x16"

	imageJobBuildWorkerSubnetPrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerSubnetPrototypeModel["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	model := make(map[string]interface{})
	model["profile"] = []interface{}{imageJobBuildWorkerProfilePrototypeModel}
	model["subnet"] = []interface{}{imageJobBuildWorkerSubnetPrototypeModel}

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototype(t *testing.T) {
	checkResult := func(result vpcv1.ImageJobBuildWorkerProfilePrototypeIntf) {
		model := new(vpcv1.ImageJobBuildWorkerProfilePrototype)
		model.Name = core.StringPtr("bx2-4x16")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "bx2-4x16"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName) {
		model := new(vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName)
		model.Name = core.StringPtr("bx2-4x16")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "bx2-4x16"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByHref) {
		model := new(vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototype(t *testing.T) {
	checkResult := func(result vpcv1.ImageJobBuildWorkerSubnetPrototypeIntf) {
		model := new(vpcv1.ImageJobBuildWorkerSubnetPrototype)
		model.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
	model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID) {
		model := new(vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID)
		model.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByCRN) {
		model := new(vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::subnet:0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByHref) {
		model := new(vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/subnets/0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobPrototype(t *testing.T) {
	checkResult := func(result vpcv1.ImageJobPrototypeIntf) {
		imageJobBuildImageRecipePrototypeModel := new(vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByID)
		imageJobBuildImageRecipePrototypeModel.Version = core.StringPtr("a3f8c2d1")
		imageJobBuildImageRecipePrototypeModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")

		resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
		resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		imageJobBuildImageTemplateModel := new(vpcv1.ImageJobBuildImageTemplate)
		imageJobBuildImageTemplateModel.Name = core.StringPtr("my-image")
		imageJobBuildImageTemplateModel.ResourceGroup = resourceGroupIdentityModel

		imageJobBuildWorkerProfilePrototypeModel := new(vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName)
		imageJobBuildWorkerProfilePrototypeModel.Name = core.StringPtr("bx2-4x16")

		imageJobBuildWorkerSubnetPrototypeModel := new(vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID)
		imageJobBuildWorkerSubnetPrototypeModel.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		imageJobBuildWorkerPrototypeModel := new(vpcv1.ImageJobBuildWorkerPrototype)
		imageJobBuildWorkerPrototypeModel.Profile = imageJobBuildWorkerProfilePrototypeModel
		imageJobBuildWorkerPrototypeModel.Subnet = imageJobBuildWorkerSubnetPrototypeModel

		model := new(vpcv1.ImageJobPrototype)
		model.Name = core.StringPtr("my-image-job")
		model.Type = core.StringPtr("build")
		model.DeleteImageJobOnImageDelete = core.BoolPtr(true)
		model.ImageRecipe = imageJobBuildImageRecipePrototypeModel
		model.ImageTemplate = imageJobBuildImageTemplateModel
		model.Timeout = core.Int64Ptr(int64(120))
		model.Worker = imageJobBuildWorkerPrototypeModel

		assert.Equal(t, result, model)
	}

	imageJobBuildImageRecipePrototypeModel := make(map[string]interface{})
	imageJobBuildImageRecipePrototypeModel["version"] = "a3f8c2d1"
	imageJobBuildImageRecipePrototypeModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"

	resourceGroupIdentityModel := make(map[string]interface{})
	resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	imageJobBuildImageTemplateModel := make(map[string]interface{})
	imageJobBuildImageTemplateModel["name"] = "my-image"
	imageJobBuildImageTemplateModel["resource_group"] = []interface{}{resourceGroupIdentityModel}

	imageJobBuildWorkerProfilePrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerProfilePrototypeModel["name"] = "bx2-4x16"

	imageJobBuildWorkerSubnetPrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerSubnetPrototypeModel["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	imageJobBuildWorkerPrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerPrototypeModel["profile"] = []interface{}{imageJobBuildWorkerProfilePrototypeModel}
	imageJobBuildWorkerPrototypeModel["subnet"] = []interface{}{imageJobBuildWorkerSubnetPrototypeModel}

	model := make(map[string]interface{})
	model["name"] = "my-image-job"
	model["type"] = "build"
	model["delete_image_job_on_image_delete"] = true
	model["image_recipe"] = []interface{}{imageJobBuildImageRecipePrototypeModel}
	model["image_template"] = []interface{}{imageJobBuildImageTemplateModel}
	model["timeout"] = int(120)
	model["worker"] = []interface{}{imageJobBuildWorkerPrototypeModel}

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsImageJobMapToImageJobPrototypeImageJobBuildPrototype(t *testing.T) {
	checkResult := func(result *vpcv1.ImageJobPrototypeImageJobBuildPrototype) {
		imageJobBuildImageRecipePrototypeModel := new(vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByID)
		imageJobBuildImageRecipePrototypeModel.Version = core.StringPtr("a3f8c2d1")
		imageJobBuildImageRecipePrototypeModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")

		resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
		resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		imageJobBuildImageTemplateModel := new(vpcv1.ImageJobBuildImageTemplate)
		imageJobBuildImageTemplateModel.Name = core.StringPtr("my-image")
		imageJobBuildImageTemplateModel.ResourceGroup = resourceGroupIdentityModel

		imageJobBuildWorkerProfilePrototypeModel := new(vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName)
		imageJobBuildWorkerProfilePrototypeModel.Name = core.StringPtr("bx2-4x16")

		imageJobBuildWorkerSubnetPrototypeModel := new(vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID)
		imageJobBuildWorkerSubnetPrototypeModel.ID = core.StringPtr("0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		imageJobBuildWorkerPrototypeModel := new(vpcv1.ImageJobBuildWorkerPrototype)
		imageJobBuildWorkerPrototypeModel.Profile = imageJobBuildWorkerProfilePrototypeModel
		imageJobBuildWorkerPrototypeModel.Subnet = imageJobBuildWorkerSubnetPrototypeModel

		model := new(vpcv1.ImageJobPrototypeImageJobBuildPrototype)
		model.Name = core.StringPtr("my-image-job")
		model.DeleteImageJobOnImageDelete = core.BoolPtr(true)
		model.ImageRecipe = imageJobBuildImageRecipePrototypeModel
		model.ImageTemplate = imageJobBuildImageTemplateModel
		model.Timeout = core.Int64Ptr(int64(120))
		model.Type = core.StringPtr("build")
		model.Worker = imageJobBuildWorkerPrototypeModel

		assert.Equal(t, result, model)
	}

	imageJobBuildImageRecipePrototypeModel := make(map[string]interface{})
	imageJobBuildImageRecipePrototypeModel["version"] = "a3f8c2d1"
	imageJobBuildImageRecipePrototypeModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"

	resourceGroupIdentityModel := make(map[string]interface{})
	resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	imageJobBuildImageTemplateModel := make(map[string]interface{})
	imageJobBuildImageTemplateModel["name"] = "my-image"
	imageJobBuildImageTemplateModel["resource_group"] = []interface{}{resourceGroupIdentityModel}

	imageJobBuildWorkerProfilePrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerProfilePrototypeModel["name"] = "bx2-4x16"

	imageJobBuildWorkerSubnetPrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerSubnetPrototypeModel["id"] = "0717-7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	imageJobBuildWorkerPrototypeModel := make(map[string]interface{})
	imageJobBuildWorkerPrototypeModel["profile"] = []interface{}{imageJobBuildWorkerProfilePrototypeModel}
	imageJobBuildWorkerPrototypeModel["subnet"] = []interface{}{imageJobBuildWorkerSubnetPrototypeModel}

	model := make(map[string]interface{})
	model["name"] = "my-image-job"
	model["delete_image_job_on_image_delete"] = true
	model["image_recipe"] = []interface{}{imageJobBuildImageRecipePrototypeModel}
	model["image_template"] = []interface{}{imageJobBuildImageTemplateModel}
	model["timeout"] = int(120)
	model["type"] = "build"
	model["worker"] = []interface{}{imageJobBuildWorkerPrototypeModel}

	result, err := vpc.ResourceIBMIsImageJobMapToImageJobPrototypeImageJobBuildPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}
