// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMVPCInstanceRescueBasic(t *testing.T) {
	var conf vpcv1.InstanceRescue
	instanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMVPCInstanceRescueDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMVPCInstanceRescueConfigBasic(instanceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMVPCInstanceRescueExists("ibm_vpc_instance_rescue.vpc_instance_rescue_instance", conf),
					resource.TestCheckResourceAttr("ibm_vpc_instance_rescue.vpc_instance_rescue_instance", "instance_id", instanceID),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_vpc_instance_rescue.vpc_instance_rescue_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMVPCInstanceRescueConfigBasic(instanceID string) string {
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
	`, instanceID)
}

func testAccCheckIBMVPCInstanceRescueExists(n string, obj vpcv1.InstanceRescue) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getInstanceRescueOptions := &vpcv1.GetInstanceRescueOptions{}

		// The ID is just the instance ID, no need to split
		getInstanceRescueOptions.SetInstanceID(rs.Primary.ID)

		instanceRescue, _, err := vpcClient.GetInstanceRescue(getInstanceRescueOptions)
		if err != nil {
			return err
		}

		obj = *instanceRescue
		return nil
	}
}

func testAccCheckIBMVPCInstanceRescueDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_instance_rescue" {
			continue
		}

		getInstanceRescueOptions := &vpcv1.GetInstanceRescueOptions{}

		// The ID is just the instance ID, no need to split
		getInstanceRescueOptions.SetInstanceID(rs.Primary.ID)

		// Try to find the rescue
		_, response, err := vpcClient.GetInstanceRescue(getInstanceRescueOptions)

		if err == nil {
			return fmt.Errorf("InstanceRescue still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for InstanceRescue (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMISInstanceRescueImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMISInstanceRescueImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.ResourceIBMISInstanceRescueDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueAccountReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "bb1b52262f7441a586f49068482f1e60"
		model["resource_type"] = "account"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.AccountReference)
	model.ID = core.StringPtr("bb1b52262f7441a586f49068482f1e60")
	model.ResourceType = core.StringPtr("account")

	result, err := vpc.ResourceIBMISInstanceRescueAccountReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueRegionReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south"
		model["name"] = "us-south"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.RegionReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south")
	model.Name = core.StringPtr("us-south")

	result, err := vpc.ResourceIBMISInstanceRescueRegionReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueKeyReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMISInstanceRescueKeyReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["id"] = "0717-80b3e36e-41f4-40e9-bd56-beae81792a68-679qb"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VolumeAttachmentDevice)
	model.ID = core.StringPtr("0717-80b3e36e-41f4-40e9-bd56-beae81792a68-679qb")

	result, err := vpc.ResourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueInstanceRescuePasswordToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMISInstanceRescueInstanceRescuePasswordToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToKeyIdentity(t *testing.T) {
	checkResult := func(result vpcv1.KeyIdentityIntf) {
		model := new(vpcv1.KeyIdentity)
		model.ID = core.StringPtr("r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45")
		model.Fingerprint = core.StringPtr("SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45"
	model["fingerprint"] = "SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY"

	result, err := vpc.ResourceIBMISInstanceRescueMapToKeyIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToKeyIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.KeyIdentityByID) {
		model := new(vpcv1.KeyIdentityByID)
		model.ID = core.StringPtr("r006-82679077-ac3b-4c10-be16-63e9c21f0f45")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-82679077-ac3b-4c10-be16-63e9c21f0f45"

	result, err := vpc.ResourceIBMISInstanceRescueMapToKeyIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToKeyIdentityByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.KeyIdentityByCRN) {
		model := new(vpcv1.KeyIdentityByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::key:r006-82679077-ac3b-4c10-be16-63e9c21f0f45"

	result, err := vpc.ResourceIBMISInstanceRescueMapToKeyIdentityByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToKeyIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.KeyIdentityByHref) {
		model := new(vpcv1.KeyIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/keys/r006-82679077-ac3b-4c10-be16-63e9c21f0f45"

	result, err := vpc.ResourceIBMISInstanceRescueMapToKeyIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToKeyIdentityByFingerprint(t *testing.T) {
	checkResult := func(result *vpcv1.KeyIdentityByFingerprint) {
		model := new(vpcv1.KeyIdentityByFingerprint)
		model.Fingerprint = core.StringPtr("SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["fingerprint"] = "SHA256:yxavE4CIOL2NlsqcurRO3xGjkP6m/0mp8ugojH5yxlY"

	result, err := vpc.ResourceIBMISInstanceRescueMapToKeyIdentityByFingerprint(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToImageIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ImageIdentityIntf) {
		model := new(vpcv1.ImageIdentity)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMISInstanceRescueMapToImageIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToImageIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ImageIdentityByID) {
		model := new(vpcv1.ImageIdentityByID)
		model.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMISInstanceRescueMapToImageIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToImageIdentityByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.ImageIdentityByCRN) {
		model := new(vpcv1.ImageIdentityByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::image:r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMISInstanceRescueMapToImageIdentityByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToImageIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ImageIdentityByHref) {
		model := new(vpcv1.ImageIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/images/r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	result, err := vpc.ResourceIBMISInstanceRescueMapToImageIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToInstanceRescueVolumeAttachmentPrototype(t *testing.T) {
	checkResult := func(result *vpcv1.InstanceRescueVolumeAttachmentPrototype) {
		volumeAllowedUsePrototypeModel := new(vpcv1.VolumeAllowedUsePrototype)
		volumeAllowedUsePrototypeModel.ApiVersion = core.StringPtr("2024-06-23")
		volumeAllowedUsePrototypeModel.BareMetalServer = core.StringPtr("enable_secure_boot == true")
		volumeAllowedUsePrototypeModel.Instance = core.StringPtr("gpu.count > 0 && enable_secure_boot == true")

		encryptionKeyIdentityModel := new(vpcv1.EncryptionKeyIdentityByCRN)
		encryptionKeyIdentityModel.CRN = core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179")

		volumeProfileIdentityModel := new(vpcv1.VolumeProfileIdentityByName)
		volumeProfileIdentityModel.Name = core.StringPtr("general-purpose")

		resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityByID)
		resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		volumePrototypeInstanceByImageContextModel := new(vpcv1.VolumePrototypeInstanceByImageContext)
		volumePrototypeInstanceByImageContextModel.AllowedUse = volumeAllowedUsePrototypeModel
		volumePrototypeInstanceByImageContextModel.Bandwidth = core.Int64Ptr(int64(1000))
		volumePrototypeInstanceByImageContextModel.Capacity = core.Int64Ptr(int64(100))
		volumePrototypeInstanceByImageContextModel.EncryptionKey = encryptionKeyIdentityModel
		volumePrototypeInstanceByImageContextModel.Iops = core.Int64Ptr(int64(10000))
		volumePrototypeInstanceByImageContextModel.Name = core.StringPtr("my-volume")
		volumePrototypeInstanceByImageContextModel.Profile = volumeProfileIdentityModel
		volumePrototypeInstanceByImageContextModel.ResourceGroup = resourceGroupIdentityModel
		volumePrototypeInstanceByImageContextModel.UserTags = []string{"testString"}

		model := new(vpcv1.InstanceRescueVolumeAttachmentPrototype)
		model.DeleteVolumeOnInstanceDelete = core.BoolPtr(true)
		model.Name = core.StringPtr("my-volume-attachment")
		model.Volume = volumePrototypeInstanceByImageContextModel

		assert.Equal(t, result, model)
	}

	volumeAllowedUsePrototypeModel := make(map[string]interface{})
	volumeAllowedUsePrototypeModel["api_version"] = "2024-06-23"
	volumeAllowedUsePrototypeModel["bare_metal_server"] = "enable_secure_boot == true"
	volumeAllowedUsePrototypeModel["instance"] = "gpu.count > 0 && enable_secure_boot == true"

	encryptionKeyIdentityModel := make(map[string]interface{})
	encryptionKeyIdentityModel["crn"] = "crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"

	volumeProfileIdentityModel := make(map[string]interface{})
	volumeProfileIdentityModel["name"] = "general-purpose"

	resourceGroupIdentityModel := make(map[string]interface{})
	resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	volumePrototypeInstanceByImageContextModel := make(map[string]interface{})
	volumePrototypeInstanceByImageContextModel["allowed_use"] = []interface{}{volumeAllowedUsePrototypeModel}
	volumePrototypeInstanceByImageContextModel["bandwidth"] = int(1000)
	volumePrototypeInstanceByImageContextModel["capacity"] = int(100)
	volumePrototypeInstanceByImageContextModel["encryption_key"] = []interface{}{encryptionKeyIdentityModel}
	volumePrototypeInstanceByImageContextModel["iops"] = int(10000)
	volumePrototypeInstanceByImageContextModel["name"] = "my-volume"
	volumePrototypeInstanceByImageContextModel["profile"] = []interface{}{volumeProfileIdentityModel}
	volumePrototypeInstanceByImageContextModel["resource_group"] = []interface{}{resourceGroupIdentityModel}
	volumePrototypeInstanceByImageContextModel["user_tags"] = []interface{}{"testString"}

	model := make(map[string]interface{})
	model["delete_volume_on_instance_delete"] = true
	model["name"] = "my-volume-attachment"
	model["volume"] = []interface{}{volumePrototypeInstanceByImageContextModel}

	result, err := vpc.ResourceIBMISInstanceRescueMapToInstanceRescueVolumeAttachmentPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToVolumePrototypeInstanceByImageContext(t *testing.T) {
	checkResult := func(result *vpcv1.VolumePrototypeInstanceByImageContext) {
		volumeAllowedUsePrototypeModel := new(vpcv1.VolumeAllowedUsePrototype)
		volumeAllowedUsePrototypeModel.ApiVersion = core.StringPtr("2024-06-23")
		volumeAllowedUsePrototypeModel.BareMetalServer = core.StringPtr("enable_secure_boot == true")
		volumeAllowedUsePrototypeModel.Instance = core.StringPtr("gpu.count > 0 && enable_secure_boot == true")

		encryptionKeyIdentityModel := new(vpcv1.EncryptionKeyIdentityByCRN)
		encryptionKeyIdentityModel.CRN = core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179")

		volumeProfileIdentityModel := new(vpcv1.VolumeProfileIdentityByName)
		volumeProfileIdentityModel.Name = core.StringPtr("general-purpose")

		resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityByID)
		resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		model := new(vpcv1.VolumePrototypeInstanceByImageContext)
		model.AllowedUse = volumeAllowedUsePrototypeModel
		model.Bandwidth = core.Int64Ptr(int64(1000))
		model.Capacity = core.Int64Ptr(int64(100))
		model.EncryptionKey = encryptionKeyIdentityModel
		model.Iops = core.Int64Ptr(int64(10000))
		model.Name = core.StringPtr("my-volume")
		model.Profile = volumeProfileIdentityModel
		model.ResourceGroup = resourceGroupIdentityModel
		model.UserTags = []string{"testString"}

		assert.Equal(t, result, model)
	}

	volumeAllowedUsePrototypeModel := make(map[string]interface{})
	volumeAllowedUsePrototypeModel["api_version"] = "2024-06-23"
	volumeAllowedUsePrototypeModel["bare_metal_server"] = "enable_secure_boot == true"
	volumeAllowedUsePrototypeModel["instance"] = "gpu.count > 0 && enable_secure_boot == true"

	encryptionKeyIdentityModel := make(map[string]interface{})
	encryptionKeyIdentityModel["crn"] = "crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"

	volumeProfileIdentityModel := make(map[string]interface{})
	volumeProfileIdentityModel["name"] = "general-purpose"

	resourceGroupIdentityModel := make(map[string]interface{})
	resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	model := make(map[string]interface{})
	model["allowed_use"] = []interface{}{volumeAllowedUsePrototypeModel}
	model["bandwidth"] = int(1000)
	model["capacity"] = int(100)
	model["encryption_key"] = []interface{}{encryptionKeyIdentityModel}
	model["iops"] = int(10000)
	model["name"] = "my-volume"
	model["profile"] = []interface{}{volumeProfileIdentityModel}
	model["resource_group"] = []interface{}{resourceGroupIdentityModel}
	model["user_tags"] = []interface{}{"testString"}

	result, err := vpc.ResourceIBMISInstanceRescueMapToVolumePrototypeInstanceByImageContext(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToVolumeAllowedUsePrototype(t *testing.T) {
	checkResult := func(result *vpcv1.VolumeAllowedUsePrototype) {
		model := new(vpcv1.VolumeAllowedUsePrototype)
		model.ApiVersion = core.StringPtr("2024-06-23")
		model.BareMetalServer = core.StringPtr("enable_secure_boot == true")
		model.Instance = core.StringPtr("gpu.count > 0 && enable_secure_boot == true")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["api_version"] = "2024-06-23"
	model["bare_metal_server"] = "enable_secure_boot == true"
	model["instance"] = "gpu.count > 0 && enable_secure_boot == true"

	result, err := vpc.ResourceIBMISInstanceRescueMapToVolumeAllowedUsePrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToEncryptionKeyIdentity(t *testing.T) {
	checkResult := func(result vpcv1.EncryptionKeyIdentityIntf) {
		model := new(vpcv1.EncryptionKeyIdentity)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"

	result, err := vpc.ResourceIBMISInstanceRescueMapToEncryptionKeyIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToEncryptionKeyIdentityByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.EncryptionKeyIdentityByCRN) {
		model := new(vpcv1.EncryptionKeyIdentityByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"

	result, err := vpc.ResourceIBMISInstanceRescueMapToEncryptionKeyIdentityByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToVolumeProfileIdentity(t *testing.T) {
	checkResult := func(result vpcv1.VolumeProfileIdentityIntf) {
		model := new(vpcv1.VolumeProfileIdentity)
		model.Name = core.StringPtr("general-purpose")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "general-purpose"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"

	result, err := vpc.ResourceIBMISInstanceRescueMapToVolumeProfileIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToVolumeProfileIdentityByName(t *testing.T) {
	checkResult := func(result *vpcv1.VolumeProfileIdentityByName) {
		model := new(vpcv1.VolumeProfileIdentityByName)
		model.Name = core.StringPtr("general-purpose")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "general-purpose"

	result, err := vpc.ResourceIBMISInstanceRescueMapToVolumeProfileIdentityByName(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToVolumeProfileIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.VolumeProfileIdentityByHref) {
		model := new(vpcv1.VolumeProfileIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/volume/profiles/general-purpose"

	result, err := vpc.ResourceIBMISInstanceRescueMapToVolumeProfileIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToResourceGroupIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ResourceGroupIdentityIntf) {
		model := new(vpcv1.ResourceGroupIdentity)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMISInstanceRescueMapToResourceGroupIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToResourceGroupIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ResourceGroupIdentityByID) {
		model := new(vpcv1.ResourceGroupIdentityByID)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMISInstanceRescueMapToResourceGroupIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMISInstanceRescueMapToInstanceRescuePrototype(t *testing.T) {
	checkResult := func(result vpcv1.InstanceRescuePrototypeIntf) {
		keyIdentityModel := new(vpcv1.KeyIdentityByID)
		keyIdentityModel.ID = core.StringPtr("r006-82679077-ac3b-4c10-be16-63e9c21f0f45")

		imageIdentityModel := new(vpcv1.ImageIdentityByID)
		imageIdentityModel.ID = core.StringPtr("r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8")

		volumeAllowedUsePrototypeModel := new(vpcv1.VolumeAllowedUsePrototype)
		volumeAllowedUsePrototypeModel.ApiVersion = core.StringPtr("2024-06-23")
		volumeAllowedUsePrototypeModel.BareMetalServer = core.StringPtr("enable_secure_boot == true")
		volumeAllowedUsePrototypeModel.Instance = core.StringPtr("gpu.count > 0 && enable_secure_boot == true")

		encryptionKeyIdentityModel := new(vpcv1.EncryptionKeyIdentityByCRN)
		encryptionKeyIdentityModel.CRN = core.StringPtr("crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179")

		volumeProfileIdentityModel := new(vpcv1.VolumeProfileIdentityByName)
		volumeProfileIdentityModel.Name = core.StringPtr("general-purpose")

		resourceGroupIdentityModel := new(vpcv1.ResourceGroupIdentityByID)
		resourceGroupIdentityModel.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		volumePrototypeInstanceByImageContextModel := new(vpcv1.VolumePrototypeInstanceByImageContext)
		volumePrototypeInstanceByImageContextModel.AllowedUse = volumeAllowedUsePrototypeModel
		volumePrototypeInstanceByImageContextModel.Bandwidth = core.Int64Ptr(int64(1000))
		volumePrototypeInstanceByImageContextModel.Capacity = core.Int64Ptr(int64(100))
		volumePrototypeInstanceByImageContextModel.EncryptionKey = encryptionKeyIdentityModel
		volumePrototypeInstanceByImageContextModel.Iops = core.Int64Ptr(int64(10000))
		volumePrototypeInstanceByImageContextModel.Name = core.StringPtr("my-volume")
		volumePrototypeInstanceByImageContextModel.Profile = volumeProfileIdentityModel
		volumePrototypeInstanceByImageContextModel.ResourceGroup = resourceGroupIdentityModel
		volumePrototypeInstanceByImageContextModel.UserTags = []string{"testString"}

		instanceRescueVolumeAttachmentPrototypeModel := new(vpcv1.InstanceRescueVolumeAttachmentPrototype)
		instanceRescueVolumeAttachmentPrototypeModel.DeleteVolumeOnInstanceDelete = core.BoolPtr(true)
		instanceRescueVolumeAttachmentPrototypeModel.Name = core.StringPtr("my-volume-attachment")
		instanceRescueVolumeAttachmentPrototypeModel.Volume = volumePrototypeInstanceByImageContextModel

		model := new(vpcv1.InstanceRescuePrototype)
		model.Keys = []vpcv1.KeyIdentityIntf{keyIdentityModel}
		model.UserData = core.StringPtr("[...]")
		model.Image = imageIdentityModel
		model.RescueVolumeAttachment = instanceRescueVolumeAttachmentPrototypeModel

		assert.Equal(t, result, model)
	}

	keyIdentityModel := make(map[string]interface{})
	keyIdentityModel["id"] = "r006-82679077-ac3b-4c10-be16-63e9c21f0f45"

	imageIdentityModel := make(map[string]interface{})
	imageIdentityModel["id"] = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"

	volumeAllowedUsePrototypeModel := make(map[string]interface{})
	volumeAllowedUsePrototypeModel["api_version"] = "2024-06-23"
	volumeAllowedUsePrototypeModel["bare_metal_server"] = "enable_secure_boot == true"
	volumeAllowedUsePrototypeModel["instance"] = "gpu.count > 0 && enable_secure_boot == true"

	encryptionKeyIdentityModel := make(map[string]interface{})
	encryptionKeyIdentityModel["crn"] = "crn:v1:bluemix:public:kms:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"

	volumeProfileIdentityModel := make(map[string]interface{})
	volumeProfileIdentityModel["name"] = "general-purpose"

	resourceGroupIdentityModel := make(map[string]interface{})
	resourceGroupIdentityModel["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	volumePrototypeInstanceByImageContextModel := make(map[string]interface{})
	volumePrototypeInstanceByImageContextModel["allowed_use"] = []interface{}{volumeAllowedUsePrototypeModel}
	volumePrototypeInstanceByImageContextModel["bandwidth"] = int(1000)
	volumePrototypeInstanceByImageContextModel["capacity"] = int(100)
	volumePrototypeInstanceByImageContextModel["encryption_key"] = []interface{}{encryptionKeyIdentityModel}
	volumePrototypeInstanceByImageContextModel["iops"] = int(10000)
	volumePrototypeInstanceByImageContextModel["name"] = "my-volume"
	volumePrototypeInstanceByImageContextModel["profile"] = []interface{}{volumeProfileIdentityModel}
	volumePrototypeInstanceByImageContextModel["resource_group"] = []interface{}{resourceGroupIdentityModel}
	volumePrototypeInstanceByImageContextModel["user_tags"] = []interface{}{"testString"}

	instanceRescueVolumeAttachmentPrototypeModel := make(map[string]interface{})
	instanceRescueVolumeAttachmentPrototypeModel["delete_volume_on_instance_delete"] = true
	instanceRescueVolumeAttachmentPrototypeModel["name"] = "my-volume-attachment"
	instanceRescueVolumeAttachmentPrototypeModel["volume"] = []interface{}{volumePrototypeInstanceByImageContextModel}

	model := make(map[string]interface{})
	model["keys"] = []interface{}{keyIdentityModel}
	model["user_data"] = "[...]"
	model["image"] = []interface{}{imageIdentityModel}
	model["rescue_volume_attachment"] = []interface{}{instanceRescueVolumeAttachmentPrototypeModel}

	result, err := vpc.ResourceIBMISInstanceRescueMapToInstanceRescuePrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}
