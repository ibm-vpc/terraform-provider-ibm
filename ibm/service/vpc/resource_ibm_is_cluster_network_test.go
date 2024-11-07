// Copyright IBM Corp. 2024 All Rights Reserved.
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

func TestAccIBMIsClusterNetworkBasic(t *testing.T) {
	var conf vpcv1.ClusterNetwork

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkExists("ibm_is_cluster_network.is_cluster_network_instance", conf),
				),
			},
		},
	})
}

func TestAccIBMIsClusterNetworkAllArgs(t *testing.T) {
	var conf vpcv1.ClusterNetwork
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkExists("ibm_is_cluster_network.is_cluster_network_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network.is_cluster_network_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkConfig(nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_cluster_network.is_cluster_network_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_cluster_network.is_cluster_network",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkConfigBasic() string {
	return fmt.Sprintf(`
		resource "ibm_is_cluster_network" "is_cluster_network_instance" {
			profile {
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"
				name = "h100"
			}
			vpc {
				crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
				href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
				id = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
			}
			zone {
				href = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
				name = "us-south-1"
			}
		}
	`)
}

func testAccCheckIBMIsClusterNetworkConfig(name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_cluster_network" "is_cluster_network_instance" {
			name = "%s"
			profile {
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"
				name = "h100"
			}
			resource_group {
				id = "fee82deba12e4c0fb69c3b09d1f12345"
			}
			subnet_prefixes {
				allocation_policy = "auto"
				cidr = "10.0.0.0/24"
			}
			vpc {
				crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
				href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
				id = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
			}
			zone {
				href = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
				name = "us-south-1"
			}
		}
	`, name)
}

func testAccCheckIBMIsClusterNetworkExists(n string, obj vpcv1.ClusterNetwork) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{}

		getClusterNetworkOptions.SetID(rs.Primary.ID)

		clusterNetwork, _, err := vpcClient.GetClusterNetwork(getClusterNetworkOptions)
		if err != nil {
			return err
		}

		obj = *clusterNetwork
		return nil
	}
}

func testAccCheckIBMIsClusterNetworkDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_cluster_network" {
			continue
		}

		getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{}

		getClusterNetworkOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := vpcClient.GetClusterNetwork(getClusterNetworkOptions)

		if err == nil {
			return fmt.Errorf("ClusterNetwork still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for ClusterNetwork (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsClusterNetworkClusterNetworkProfileReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"
		model["name"] = "h100"
		model["resource_type"] = "cluster_network_profile"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkProfileReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100")
	model.Name = core.StringPtr("h100")
	model.ResourceType = core.StringPtr("cluster_network_profile")

	result, err := vpc.ResourceIBMIsClusterNetworkClusterNetworkProfileReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkResourceGroupReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345"
		model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"
		model["name"] = "my-resource-group"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ResourceGroupReference)
	model.Href = core.StringPtr("https://resource-controller.cloud.ibm.com/v2/resource_groups/fee82deba12e4c0fb69c3b09d1f12345")
	model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")
	model.Name = core.StringPtr("my-resource-group")

	result, err := vpc.ResourceIBMIsClusterNetworkResourceGroupReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["allocation_policy"] = "auto"
		model["cidr"] = "10.0.0.0/24"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetPrefix)
	model.AllocationPolicy = core.StringPtr("auto")
	model.CIDR = core.StringPtr("10.0.0.0/24")

	result, err := vpc.ResourceIBMIsClusterNetworkClusterNetworkSubnetPrefixToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkVPCReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		model["id"] = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		model["name"] = "my-vpc"
		model["resource_type"] = "vpc"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.VPCReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
	model.ID = core.StringPtr("r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
	model.Name = core.StringPtr("my-vpc")
	model.ResourceType = core.StringPtr("vpc")

	result, err := vpc.ResourceIBMIsClusterNetworkVPCReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.ResourceIBMIsClusterNetworkDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkZoneReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		model["name"] = "us-south-1"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ZoneReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	model.Name = core.StringPtr("us-south-1")

	result, err := vpc.ResourceIBMIsClusterNetworkZoneReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "resource_suspended_by_provider"
		model["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkLifecycleReason)
	model.Code = core.StringPtr("resource_suspended_by_provider")
	model.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	result, err := vpc.ResourceIBMIsClusterNetworkClusterNetworkLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ClusterNetworkProfileIdentityIntf) {
		model := new(vpcv1.ClusterNetworkProfileIdentity)
		model.Name = core.StringPtr("h100")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "h100"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByName(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkProfileIdentityByName) {
		model := new(vpcv1.ClusterNetworkProfileIdentityByName)
		model.Name = core.StringPtr("h100")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "h100"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByName(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkProfileIdentityByHref) {
		model := new(vpcv1.ClusterNetworkProfileIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToClusterNetworkProfileIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToVPCIdentity(t *testing.T) {
	checkResult := func(result vpcv1.VPCIdentityIntf) {
		model := new(vpcv1.VPCIdentity)
		model.ID = core.StringPtr("r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToVPCIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToVPCIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.VPCIdentityByID) {
		model := new(vpcv1.VPCIdentityByID)
		model.ID = core.StringPtr("r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToVPCIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToVPCIdentityByCRN(t *testing.T) {
	checkResult := func(result *vpcv1.VPCIdentityByCRN) {
		model := new(vpcv1.VPCIdentityByCRN)
		model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToVPCIdentityByCRN(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToVPCIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.VPCIdentityByHref) {
		model := new(vpcv1.VPCIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToVPCIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToZoneIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ZoneIdentityIntf) {
		model := new(vpcv1.ZoneIdentity)
		model.Name = core.StringPtr("us-south-1")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "us-south-1"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToZoneIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToZoneIdentityByName(t *testing.T) {
	checkResult := func(result *vpcv1.ZoneIdentityByName) {
		model := new(vpcv1.ZoneIdentityByName)
		model.Name = core.StringPtr("us-south-1")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["name"] = "us-south-1"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToZoneIdentityByName(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToZoneIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ZoneIdentityByHref) {
		model := new(vpcv1.ZoneIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToZoneIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToResourceGroupIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ResourceGroupIdentityIntf) {
		model := new(vpcv1.ResourceGroupIdentity)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToResourceGroupIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToResourceGroupIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ResourceGroupIdentityByID) {
		model := new(vpcv1.ResourceGroupIdentityByID)
		model.ID = core.StringPtr("fee82deba12e4c0fb69c3b09d1f12345")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "fee82deba12e4c0fb69c3b09d1f12345"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToResourceGroupIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkMapToClusterNetworkSubnetPrefixPrototype(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkSubnetPrefixPrototype) {
		model := new(vpcv1.ClusterNetworkSubnetPrefixPrototype)
		model.CIDR = core.StringPtr("10.0.0.0/24")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["cidr"] = "10.0.0.0/24"

	result, err := vpc.ResourceIBMIsClusterNetworkMapToClusterNetworkSubnetPrefixPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}
