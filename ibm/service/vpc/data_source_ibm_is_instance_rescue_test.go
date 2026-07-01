// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 99-SNAPSHOT-c913c18c-20260623-001551
 */

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMVPCInstanceRescueDataSourceBasic(t *testing.T) {
	instanceRescueInstanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMVPCInstanceRescueDataSourceConfigBasic(instanceRescueInstanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_vpc_instance_rescue.vpc_instance_rescue_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_vpc_instance_rescue.vpc_instance_rescue_instance", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_vpc_instance_rescue.vpc_instance_rescue_instance", "image.#"),
					resource.TestCheckResourceAttrSet("data.ibm_vpc_instance_rescue.vpc_instance_rescue_instance", "keys.#"),
					resource.TestCheckResourceAttrSet("data.ibm_vpc_instance_rescue.vpc_instance_rescue_instance", "rescue_volume_attachment.#"),
				),
			},
		},
	})
}

func testAccCheckIBMVPCInstanceRescueDataSourceConfigBasic(instanceRescueInstanceID string) string {
	return fmt.Sprintf(`
		data "ibm_is_image" "rescue_image" {
			name = "ibm-ubuntu-22-04-minimal-amd64-1"
		}

		resource "ibm_is_instance_rescue" "vpc_instance_rescue_instance" {
			instance_id = "%s"
			
			image {
				id = data.ibm_is_image.rescue_image.id
			}
			
			rescue_volume_attachment {
				name = "rescue-volume-attachment"
				delete_volume_on_instance_delete = true
				volume {
					name = "rescue-volume"
					profile = "general-purpose"
				}
			}
		}

		data "ibm_is_instance_rescue" "vpc_instance_rescue_instance" {
			instance_id = ibm_is_instance_rescue.vpc_instance_rescue_instance.instance_id
		}
	`, instanceRescueInstanceID)
}

func TestDataSourceIBMISInstanceRescueImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMISInstanceRescueImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMISInstanceRescueDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueAccountReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "bb1b52262f7441a586f49068482f1e60"
		model["resource_type"] = "account"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.AccountReference)
	model.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	model.ResourceType = core.StringPtr("account")

	result, err := vpc.DataSourceIBMISInstanceRescueAccountReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueRegionReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		model["name"] = "us-south"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.RegionReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	model.Name = core.StringPtr("us-south")

	result, err := vpc.DataSourceIBMISInstanceRescueRegionReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueKeyReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["fingerprint"] = "SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		model["id"] = "r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		model["name"] = "my-key-1"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.KeyReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	model.Deleted = deletedModel
	model.Fingerprint = core.StringPtr("SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	model.ID = core.StringPtr("r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	model.Name = core.StringPtr("my-key-1")

	result, err := vpc.DataSourceIBMISInstanceRescueKeyReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueInstanceRescuePasswordToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		instanceRescueEncryptionKeyModel := make(map[string]interface{})
		instanceRescueEncryptionKeyModel["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		instanceRescueEncryptionKeyModel["deleted"] = []map[string]interface{}{deletedModel}
		instanceRescueEncryptionKeyModel["fingerprint"] = "SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY"
		instanceRescueEncryptionKeyModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		instanceRescueEncryptionKeyModel["id"] = "r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		instanceRescueEncryptionKeyModel["name"] = "my-key-1"

		model := make(map[string]interface{})
		model["encrypted_password"] = "qQ+/YEApnl1ZtEgIrfprzb065307thTkzlnLqL5ICpesdbBN03dyCQ=="
		model["encryption_key"] = []map[string]interface{}{instanceRescueEncryptionKeyModel}

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	instanceRescueEncryptionKeyModel := new(vpcv1.InstanceRescueEncryptionKey)
	instanceRescueEncryptionKeyModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	instanceRescueEncryptionKeyModel.Deleted = deletedModel
	instanceRescueEncryptionKeyModel.Fingerprint = core.StringPtr("SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY")
	instanceRescueEncryptionKeyModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	instanceRescueEncryptionKeyModel.ID = core.StringPtr("r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	instanceRescueEncryptionKeyModel.Name = core.StringPtr("my-key-1")

	model := new(vpcv1.InstanceRescuePassword)
	model.EncryptedPassword = createMockByteArray("qQ+/YEApnl1ZtEgIrfprzb065307thTkzlnLqL5ICpesdbBN03dyCQ==")
	model.EncryptionKey = instanceRescueEncryptionKeyModel

	result, err := vpc.DataSourceIBMISInstanceRescueInstanceRescuePasswordToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["fingerprint"] = "SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		model["id"] = "r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
		model["name"] = "my-key-1"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.InstanceRescueEncryptionKey)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	model.Deleted = deletedModel
	model.Fingerprint = core.StringPtr("SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	model.ID = core.StringPtr("r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
	model.Name = core.StringPtr("my-key-1")

	result, err := vpc.DataSourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		volumeAttachmentDeviceModel := make(map[string]interface{})
		volumeAttachmentDeviceModel["id"] = "0717-80b3e36e-41f4-40e9-bd56-beae81792a68-679qb"

		volumeReferenceVolumeAttachmentContextModel := make(map[string]interface{})
		volumeReferenceVolumeAttachmentContextModel["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"
		volumeReferenceVolumeAttachmentContextModel["deleted"] = []map[string]interface{}{deletedModel}
		volumeReferenceVolumeAttachmentContextModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/volumes/r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"
		volumeReferenceVolumeAttachmentContextModel["id"] = "r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"
		volumeReferenceVolumeAttachmentContextModel["name"] = "my-volume"
		volumeReferenceVolumeAttachmentContextModel["resource_type"] = "volume"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["device"] = []map[string]interface{}{volumeAttachmentDeviceModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/volume_attachments/0717-82cbf856-9cbb-45fb-b62f-d7bcef32399a"
		model["id"] = "0717-82cbf856-9cbb-45fb-b62f-d7bcef32399a"
		model["name"] = "my-volume-attachment"
		model["volume"] = []map[string]interface{}{volumeReferenceVolumeAttachmentContextModel}

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	volumeAttachmentDeviceModel := new(vpcv1.VolumeAttachmentDevice)
	volumeAttachmentDeviceModel.ID = core.StringPtr("0717-80b3e36e-41f4-40e9-bd56-beae81792a68-679qb")

	volumeReferenceVolumeAttachmentContextModel := new(vpcv1.VolumeReferenceVolumeAttachmentContext)
	volumeReferenceVolumeAttachmentContextModel.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5")
	volumeReferenceVolumeAttachmentContextModel.Deleted = deletedModel
	volumeReferenceVolumeAttachmentContextModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volumes/r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5")
	volumeReferenceVolumeAttachmentContextModel.ID = core.StringPtr("r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5")
	volumeReferenceVolumeAttachmentContextModel.Name = core.StringPtr("my-volume")
	volumeReferenceVolumeAttachmentContextModel.ResourceType = core.StringPtr("volume")

	model := new(vpcv1.VolumeAttachmentReferenceInstanceContext)
	model.Deleted = deletedModel
	model.Device = volumeAttachmentDeviceModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/volume_attachments/0717-82cbf856-9cbb-45fb-b62f-d7bcef32399a")
	model.ID = core.StringPtr("0717-82cbf856-9cbb-45fb-b62f-d7bcef32399a")
	model.Name = core.StringPtr("my-volume-attachment")
	model.Volume = volumeReferenceVolumeAttachmentContextModel

	result, err := vpc.DataSourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "0717-80b3e36e-41f4-40e9-bd56-beae81792a68-679qb"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VolumeAttachmentDevice)
	model.ID = core.StringPtr("0717-80b3e36e-41f4-40e9-bd56-beae81792a68-679qb")

	result, err := vpc.DataSourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/volumes/r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"
		model["id"] = "r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5"
		model["name"] = "my-volume"
		model["resource_type"] = "volume"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.VolumeReferenceVolumeAttachmentContext)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south-1:a/aa2432b1fa4d4ace891e9b80fc104e34::volume:r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volumes/r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5")
	model.ID = core.StringPtr("r006-1a6b7274-678d-4dfb-8981-c71dd9d4daa5")
	model.Name = core.StringPtr("my-volume")
	model.ResourceType = core.StringPtr("volume")

	result, err := vpc.DataSourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

// Helper function to create a pointer to a byte array
func createMockByteArray(s string) *[]byte {
	b := []byte(s)
	return &b
}
