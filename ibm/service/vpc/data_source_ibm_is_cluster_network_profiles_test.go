// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsClusterNetworkProfilesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkProfilesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profiles.is_cluster_network_profiles_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profiles.is_cluster_network_profiles_instance", "profiles.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkProfilesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_cluster_network_profiles" "is_cluster_network_profiles_instance" {
		}
	`)
}

func TestDataSourceIBMIsClusterNetworkProfilesClusterNetworkProfileToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		instanceProfileReferenceModel := make(map[string]interface{})
		instanceProfileReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/tx3-160x1536-8h100"
		instanceProfileReferenceModel["name"] = "tx3-160x1536-8h100"
		instanceProfileReferenceModel["resource_type"] = "instance_profile"

		zoneReferenceModel := make(map[string]interface{})
		zoneReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		zoneReferenceModel["name"] = "us-south-1"

		model := make(map[string]interface{})
		model["family"] = "vela"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"
		model["name"] = "h100"
		model["resource_type"] = "cluster_network_profile"
		model["supported_instance_profiles"] = []map[string]interface{}{instanceProfileReferenceModel}
		model["zones"] = []map[string]interface{}{zoneReferenceModel}

		assert.Equal(t, result, model)
	}

	instanceProfileReferenceModel := new(vpcv1.InstanceProfileReference)
	instanceProfileReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/tx3-160x1536-8h100")
	instanceProfileReferenceModel.Name = core.StringPtr("tx3-160x1536-8h100")
	instanceProfileReferenceModel.ResourceType = core.StringPtr("instance_profile")

	zoneReferenceModel := new(vpcv1.ZoneReference)
	zoneReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	zoneReferenceModel.Name = core.StringPtr("us-south-1")

	model := new(vpcv1.ClusterNetworkProfile)
	model.Family = core.StringPtr("vela")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100")
	model.Name = core.StringPtr("h100")
	model.ResourceType = core.StringPtr("cluster_network_profile")
	model.SupportedInstanceProfiles = []vpcv1.InstanceProfileReference{*instanceProfileReferenceModel}
	model.Zones = []vpcv1.ZoneReference{*zoneReferenceModel}

	result, err := vpc.DataSourceIBMIsClusterNetworkProfilesClusterNetworkProfileToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkProfilesInstanceProfileReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16"
		model["name"] = "bx2-4x16"
		model["resource_type"] = "instance_profile"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.InstanceProfileReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instance/profiles/bx2-4x16")
	model.Name = core.StringPtr("bx2-4x16")
	model.ResourceType = core.StringPtr("instance_profile")

	result, err := vpc.DataSourceIBMIsClusterNetworkProfilesInstanceProfileReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkProfilesZoneReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		model["name"] = "us-south-1"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ZoneReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	model.Name = core.StringPtr("us-south-1")

	result, err := vpc.DataSourceIBMIsClusterNetworkProfilesZoneReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
