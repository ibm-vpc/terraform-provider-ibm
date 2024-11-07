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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsClusterNetworkInterfaceBasic(t *testing.T) {
	var conf vpcv1.ClusterNetworkInterface
	clusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkInterfaceConfigBasic(clusterNetworkID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkInterfaceExists("ibm_is_cluster_network_interface.is_cluster_network_interface_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_interface.is_cluster_network_interface_instance", "cluster_network_id", clusterNetworkID),
				),
			},
		},
	})
}

func TestAccIBMIsClusterNetworkInterfaceAllArgs(t *testing.T) {
	var conf vpcv1.ClusterNetworkInterface
	clusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkInterfaceConfig(clusterNetworkID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkInterfaceExists("ibm_is_cluster_network_interface.is_cluster_network_interface_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_interface.is_cluster_network_interface_instance", "cluster_network_id", clusterNetworkID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_interface.is_cluster_network_interface_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkInterfaceConfig(clusterNetworkID, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_cluster_network_interface.is_cluster_network_interface_instance", "cluster_network_id", clusterNetworkID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_interface.is_cluster_network_interface_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_cluster_network_interface.is_cluster_network_interface",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkInterfaceConfigBasic(clusterNetworkID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_cluster_network_interface" "is_cluster_network_interface_instance" {
			cluster_network_id = "%s"
		}
	`, clusterNetworkID)
}

func testAccCheckIBMIsClusterNetworkInterfaceConfig(clusterNetworkID string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_cluster_network_interface" "is_cluster_network_interface_instance" {
			cluster_network_id = "%s"
			name = "%s"
			primary_ip {
				address = "10.1.0.6"
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				name = "my-cluster-network-subnet-reserved-ip"
			}
			subnet {
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
				id = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
			}
		}
	`, clusterNetworkID, name)
}

func testAccCheckIBMIsClusterNetworkInterfaceExists(n string, obj vpcv1.ClusterNetworkInterface) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getClusterNetworkInterfaceOptions := &vpcv1.GetClusterNetworkInterfaceOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getClusterNetworkInterfaceOptions.SetClusterNetworkID(parts[0])
		getClusterNetworkInterfaceOptions.SetID(parts[1])

		clusterNetworkInterface, _, err := vpcClient.GetClusterNetworkInterface(getClusterNetworkInterfaceOptions)
		if err != nil {
			return err
		}

		obj = *clusterNetworkInterface
		return nil
	}
}

func testAccCheckIBMIsClusterNetworkInterfaceDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_cluster_network_interface" {
			continue
		}

		getClusterNetworkInterfaceOptions := &vpcv1.GetClusterNetworkInterfaceOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getClusterNetworkInterfaceOptions.SetClusterNetworkID(parts[0])
		getClusterNetworkInterfaceOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetClusterNetworkInterface(getClusterNetworkInterfaceOptions)

		if err == nil {
			return fmt.Errorf("ClusterNetworkInterface still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for ClusterNetworkInterface (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["address"] = "10.1.0.6"
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		model["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		model["name"] = "my-cluster-network-subnet-reserved-ip"
		model["resource_type"] = "cluster_network_subnet_reserved_ip"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReservedIPReference)
	model.Address = core.StringPtr("10.1.0.6")
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	model.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	model.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")
	model.ResourceType = core.StringPtr("cluster_network_subnet_reserved_ip")

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReservedIPReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		model["id"] = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		model["name"] = "my-cluster-network-subnet"
		model["resource_type"] = "cluster_network_subnet"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReference)
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	model.ID = core.StringPtr("0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	model.Name = core.StringPtr("my-cluster-network-subnet")
	model.ResourceType = core.StringPtr("cluster_network_subnet")

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceClusterNetworkSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceLifecycleReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "resource_suspended_by_provider"
		model["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceLifecycleReason)
	model.Code = core.StringPtr("resource_suspended_by_provider")
	model.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0717-fb880975-db45-4459-8548-64e3995ac213"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceTarget)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0717-fb880975-db45-4459-8548-64e3995ac213")
	model.Name = core.StringPtr("my-instance-network-attachment")
	model.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContextToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0717-fb880975-db45-4459-8548-64e3995ac213"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContext)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0717-fb880975-db45-4459-8548-64e3995ac213")
	model.Name = core.StringPtr("my-instance-network-attachment")
	model.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceClusterNetworkInterfaceContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceVPCReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceVPCReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceZoneReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		model["name"] = "us-south-1"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ZoneReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	model.Name = core.StringPtr("us-south-1")

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceZoneReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototype(t *testing.T) {
	checkResult := func(result vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeIntf) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototype)
		model.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")
		model.Address = core.StringPtr("192.168.3.4")
		model.AutoDelete = core.BoolPtr(false)
		model.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
	model["address"] = "192.168.3.4"
	model["auto_delete"] = false
	model["name"] = "my-cluster-network-subnet-reserved-ip"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext(t *testing.T) {
	checkResult := func(result vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext)
		model.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID)
		model.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext)
		model.Address = core.StringPtr("192.168.3.4")
		model.AutoDelete = core.BoolPtr(false)
		model.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["address"] = "192.168.3.4"
	model["auto_delete"] = false
	model["name"] = "my-cluster-network-subnet-reserved-ip"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ClusterNetworkSubnetIdentityIntf) {
		model := new(vpcv1.ClusterNetworkSubnetIdentity)
		model.ID = core.StringPtr("0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkSubnetIdentityByID) {
		model := new(vpcv1.ClusterNetworkSubnetIdentityByID)
		model.ID = core.StringPtr("0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkSubnetIdentityByHref) {
		model := new(vpcv1.ClusterNetworkSubnetIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"

	result, err := vpc.ResourceIBMIsClusterNetworkInterfaceMapToClusterNetworkSubnetIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}
