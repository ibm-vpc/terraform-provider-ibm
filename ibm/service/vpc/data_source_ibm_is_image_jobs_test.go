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
	. "github.com/IBM-Cloud/terraform-provider-ibm/ibm/unittest"
	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsImageJobsDataSourceBasic(t *testing.T) {
	imageJobType := "build"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageJobsDataSourceConfigBasic(imageJobType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.#"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.type", imageJobType),
				),
			},
		},
	})
}

func TestAccIBMIsImageJobsDataSourceAllArgs(t *testing.T) {
	imageJobName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	imageJobType := "build"
	imageJobDeleteImageJobOnImageDelete := "true"
	imageJobTimeout := fmt.Sprintf("%d", acctest.RandIntRange(30, 480))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsImageJobsDataSourceConfig(imageJobName, imageJobType, imageJobDeleteImageJobOnImageDelete, imageJobTimeout),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.completed_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.name", imageJobName),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.started_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.status"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.type", imageJobType),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.delete_image_job_on_image_delete", imageJobDeleteImageJobOnImageDelete),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.timeout", imageJobTimeout),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageJobsDataSourceConfigBasic(imageJobType string) string {
	return fmt.Sprintf(`
		resource "ibm_is_image_job" "is_image_job_instance" {
			type = "%s"
		}

		data "ibm_is_image_jobs" "is_image_jobs_instance" {
			name = ibm_is_image_job.is_image_job_instance.name
			status = ibm_is_image_job.is_image_job_instance.status
		}
	`, imageJobType)
}

func testAccCheckIBMIsImageJobsDataSourceConfig(imageJobName string, imageJobType string, imageJobDeleteImageJobOnImageDelete string, imageJobTimeout string) string {
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

		data "ibm_is_image_jobs" "is_image_jobs_instance" {
			name = ibm_is_image_job.is_image_job_instance.name
			status = ibm_is_image_job.is_image_job_instance.status
		}
	`, imageJobName, imageJobType, imageJobDeleteImageJobOnImageDelete, imageJobTimeout)
}

func TestDataSourceIBMIsImageJobsImageJobToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		imageJobStatusReasonsModel := make(map[string]interface{})
		imageJobStatusReasonsModel["code"] = "internal_error"
		imageJobStatusReasonsModel["message"] = "A failure occurred"
		imageJobStatusReasonsModel["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-image-builder"

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

		imageJobBuildImageRecipeModel := make(map[string]interface{})
		imageJobBuildImageRecipeModel["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12"
		imageJobBuildImageRecipeModel["deleted"] = []map[string]interface{}{deletedModel}
		imageJobBuildImageRecipeModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8/history/a3f8c2d1"
		imageJobBuildImageRecipeModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
		imageJobBuildImageRecipeModel["name"] = "my-image-recipe"
		imageJobBuildImageRecipeModel["resource_type"] = "image_recipe"
		imageJobBuildImageRecipeModel["version"] = "a3f8c2d1"

		resourceGroupIdentityModel := make(map[string]interface{})
		resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

		imageJobBuildImageTemplateModel := make(map[string]interface{})
		imageJobBuildImageTemplateModel["name"] = "my-image"
		imageJobBuildImageTemplateModel["resource_group"] = []map[string]interface{}{resourceGroupIdentityModel}

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

		imageJobBuildWorkerModel := make(map[string]interface{})
		imageJobBuildWorkerModel["instance"] = []map[string]interface{}{instanceReferenceModel}
		imageJobBuildWorkerModel["profile"] = []map[string]interface{}{imageJobBuildWorkerProfileModel}
		imageJobBuildWorkerModel["subnet"] = []map[string]interface{}{subnetReferenceModel}

		model := make(map[string]interface{})
		model["completed_at"] = "2026-01-02T03:04:05.006Z"
		model["created_at"] = "2026-01-02T03:04:05.006Z"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image/jobs/r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf"
		model["id"] = "r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf"
		model["name"] = "my-image-job"
		model["resource_type"] = "image_job"
		model["started_at"] = "2026-01-02T03:04:05.006Z"
		model["status"] = "running"
		model["status_reasons"] = []map[string]interface{}{imageJobStatusReasonsModel}
		model["type"] = "build"
		model["delete_image_job_on_image_delete"] = true
		model["image"] = []map[string]interface{}{imageReferenceModel}
		model["image_recipe"] = []map[string]interface{}{imageJobBuildImageRecipeModel}
		model["image_template"] = []map[string]interface{}{imageJobBuildImageTemplateModel}
		model["timeout"] = int(120)
		model["worker"] = []map[string]interface{}{imageJobBuildWorkerModel}

		assert.Equal(t, result, model)
	}

	imageJobStatusReasonsModel := new(vpcv1.ImageJobStatusReasons)
	imageJobStatusReasonsModel.Code = core.StringPtr("internal_error")
	imageJobStatusReasonsModel.Message = core.StringPtr("A failure occurred")
	imageJobStatusReasonsModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-image-builder")

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

	imageJobBuildImageRecipeModel := new(vpcv1.ImageJobBuildImageRecipe)
	imageJobBuildImageRecipeModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12")
	imageJobBuildImageRecipeModel.Deleted = deletedModel
	imageJobBuildImageRecipeModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8/history/a3f8c2d1")
	imageJobBuildImageRecipeModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")
	imageJobBuildImageRecipeModel.Name = core.StringPtr("my-image-recipe")
	imageJobBuildImageRecipeModel.ResourceType = core.StringPtr("image_recipe")
	imageJobBuildImageRecipeModel.Version = core.StringPtr("a3f8c2d1")

	resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
	resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

	imageJobBuildImageTemplateModel := new(vpcv1.ImageJobBuildImageTemplate)
	imageJobBuildImageTemplateModel.Name = core.StringPtr("my-image")
	imageJobBuildImageTemplateModel.ResourceGroup = resourceGroupIdentityModel

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

	imageJobBuildWorkerModel := new(vpcv1.ImageJobBuildWorker)
	imageJobBuildWorkerModel.Instance = instanceReferenceModel
	imageJobBuildWorkerModel.Profile = imageJobBuildWorkerProfileModel
	imageJobBuildWorkerModel.Subnet = subnetReferenceModel

	model := new(vpcv1.ImageJob)
	model.CompletedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.CreatedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image/jobs/r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf")
	model.ID = core.StringPtr("r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf")
	model.Name = core.StringPtr("my-image-job")
	model.ResourceType = core.StringPtr("image_job")
	model.StartedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.Status = core.StringPtr("running")
	model.StatusReasons = []vpcv1.ImageJobStatusReasons{*imageJobStatusReasonsModel}
	model.Type = core.StringPtr("build")
	model.DeleteImageJobOnImageDelete = core.BoolPtr(true)
	model.Image = imageReferenceModel
	model.ImageRecipe = imageJobBuildImageRecipeModel
	model.ImageTemplate = imageJobBuildImageTemplateModel
	model.Timeout = core.Int64Ptr(int64(120))
	model.Worker = imageJobBuildWorkerModel

	result, err := vpc.DataSourceIBMIsImageJobsImageJobToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageJobStatusReasonsToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsImageJobStatusReasonsToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsImageJobsDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageRemoteContextImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsImageRemoteContextImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsAccountReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "bb1b52262f7441a586f49068482f1e60"
		model["resource_type"] = "account"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.AccountReference)
	model.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	model.ResourceType = core.StringPtr("account")

	result, err := vpc.DataSourceIBMIsImageJobsAccountReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsRegionReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		model["name"] = "us-south"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.RegionReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	model.Name = core.StringPtr("us-south")

	result, err := vpc.DataSourceIBMIsImageJobsRegionReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageJobBuildImageRecipeToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsImageJobBuildImageRecipeToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageJobBuildImageTemplateToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsImageJobBuildImageTemplateToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsResourceGroupIdentityToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ResourceGroupIdentity)
	model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

	result, err := vpc.DataSourceIBMIsImageJobsResourceGroupIdentityToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsResourceGroupIdentityResourceGroupIdentityByIDToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
	model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

	result, err := vpc.DataSourceIBMIsImageJobsResourceGroupIdentityResourceGroupIdentityByIDToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageJobBuildWorkerToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsImageJobBuildWorkerToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsInstanceReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsInstanceReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageJobBuildWorkerProfileToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsImageJobBuildWorkerProfileToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsSubnetReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsImageJobsSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsImageJobsImageJobBuildToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		imageJobStatusReasonsModel := make(map[string]interface{})
		imageJobStatusReasonsModel["code"] = "internal_error"
		imageJobStatusReasonsModel["message"] = "A failure occurred"
		imageJobStatusReasonsModel["more_info"] = "https://cloud.ibm.com/docs/vpc?topic=vpc-image-builder"

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

		imageJobBuildImageRecipeModel := make(map[string]interface{})
		imageJobBuildImageRecipeModel["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12"
		imageJobBuildImageRecipeModel["deleted"] = []map[string]interface{}{deletedModel}
		imageJobBuildImageRecipeModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8/history/a3f8c2d1"
		imageJobBuildImageRecipeModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8"
		imageJobBuildImageRecipeModel["name"] = "my-image-recipe"
		imageJobBuildImageRecipeModel["resource_type"] = "image_recipe"
		imageJobBuildImageRecipeModel["version"] = "a3f8c2d1"

		resourceGroupIdentityModel := make(map[string]interface{})
		resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

		imageJobBuildImageTemplateModel := make(map[string]interface{})
		imageJobBuildImageTemplateModel["name"] = "my-image"
		imageJobBuildImageTemplateModel["resource_group"] = []map[string]interface{}{resourceGroupIdentityModel}

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

		imageJobBuildWorkerModel := make(map[string]interface{})
		imageJobBuildWorkerModel["instance"] = []map[string]interface{}{instanceReferenceModel}
		imageJobBuildWorkerModel["profile"] = []map[string]interface{}{imageJobBuildWorkerProfileModel}
		imageJobBuildWorkerModel["subnet"] = []map[string]interface{}{subnetReferenceModel}

		model := make(map[string]interface{})
		model["completed_at"] = "2026-01-02T03:04:05.006Z"
		model["created_at"] = "2026-01-02T03:04:05.006Z"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/image/jobs/r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf"
		model["id"] = "r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf"
		model["name"] = "my-image-job"
		model["resource_type"] = "image_job"
		model["started_at"] = "2026-01-02T03:04:05.006Z"
		model["status"] = "running"
		model["status_reasons"] = []map[string]interface{}{imageJobStatusReasonsModel}
		model["delete_image_job_on_image_delete"] = true
		model["image"] = []map[string]interface{}{imageReferenceModel}
		model["image_recipe"] = []map[string]interface{}{imageJobBuildImageRecipeModel}
		model["image_template"] = []map[string]interface{}{imageJobBuildImageTemplateModel}
		model["timeout"] = int(120)
		model["type"] = "build"
		model["worker"] = []map[string]interface{}{imageJobBuildWorkerModel}

		assert.Equal(t, result, model)
	}

	imageJobStatusReasonsModel := new(vpcv1.ImageJobStatusReasons)
	imageJobStatusReasonsModel.Code = core.StringPtr("internal_error")
	imageJobStatusReasonsModel.Message = core.StringPtr("A failure occurred")
	imageJobStatusReasonsModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/docs/vpc?topic=vpc-image-builder")

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

	imageJobBuildImageRecipeModel := new(vpcv1.ImageJobBuildImageRecipe)
	imageJobBuildImageRecipeModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::image-recipe:my-rhel9-cuda12")
	imageJobBuildImageRecipeModel.Deleted = deletedModel
	imageJobBuildImageRecipeModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image_recipes/r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8/history/a3f8c2d1")
	imageJobBuildImageRecipeModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7f1acb8")
	imageJobBuildImageRecipeModel.Name = core.StringPtr("my-image-recipe")
	imageJobBuildImageRecipeModel.ResourceType = core.StringPtr("image_recipe")
	imageJobBuildImageRecipeModel.Version = core.StringPtr("a3f8c2d1")

	resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityResourceGroupIdentityByID)
	resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

	imageJobBuildImageTemplateModel := new(vpcv1.ImageJobBuildImageTemplate)
	imageJobBuildImageTemplateModel.Name = core.StringPtr("my-image")
	imageJobBuildImageTemplateModel.ResourceGroup = resourceGroupIdentityModel

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

	imageJobBuildWorkerModel := new(vpcv1.ImageJobBuildWorker)
	imageJobBuildWorkerModel.Instance = instanceReferenceModel
	imageJobBuildWorkerModel.Profile = imageJobBuildWorkerProfileModel
	imageJobBuildWorkerModel.Subnet = subnetReferenceModel

	model := new(vpcv1.ImageJobBuild)
	model.CompletedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.CreatedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/image/jobs/r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf")
	model.ID = core.StringPtr("r006-3e5acf73-9349-4ef8-a87b-74c7dafb52cf")
	model.Name = core.StringPtr("my-image-job")
	model.ResourceType = core.StringPtr("image_job")
	model.StartedAt = CreateMockDateTime("2026-01-02T03:04:05.006Z")
	model.Status = core.StringPtr("running")
	model.StatusReasons = []vpcv1.ImageJobStatusReasons{*imageJobStatusReasonsModel}
	model.DeleteImageJobOnImageDelete = core.BoolPtr(true)
	model.Image = imageReferenceModel
	model.ImageRecipe = imageJobBuildImageRecipeModel
	model.ImageTemplate = imageJobBuildImageTemplateModel
	model.Timeout = core.Int64Ptr(int64(120))
	model.Type = core.StringPtr("build")
	model.Worker = imageJobBuildWorkerModel

	result, err := vpc.DataSourceIBMIsImageJobsImageJobBuildToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
