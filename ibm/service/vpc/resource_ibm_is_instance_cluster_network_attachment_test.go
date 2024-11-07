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

func TestAccIBMIsInstanceClusterNetworkAttachmentBasic(t *testing.T) {
	var conf vpcv1.InstanceClusterNetworkAttachment
	instanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsInstanceClusterNetworkAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceClusterNetworkAttachmentConfigBasic(instanceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsInstanceClusterNetworkAttachmentExists("ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "instance_id", instanceID),
				),
			},
		},
	})
}

func TestAccIBMIsInstanceClusterNetworkAttachmentAllArgs(t *testing.T) {
	var conf vpcv1.InstanceClusterNetworkAttachment
	instanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsInstanceClusterNetworkAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceClusterNetworkAttachmentConfig(instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsInstanceClusterNetworkAttachmentExists("ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "instance_id", instanceID),
					resource.TestCheckResourceAttr("ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceClusterNetworkAttachmentConfig(instanceID, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "instance_id", instanceID),
					resource.TestCheckResourceAttr("ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentConfigBasic(instanceID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_instance_cluster_network_attachment" "is_instance_cluster_network_attachment_instance" {
			instance_id = "%s"
			cluster_network_interface {
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
				id = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
				name = "my-cluster-network-interface"
				primary_ip {
					address = "10.1.0.6"
					deleted {
						more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
					}
					href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
					id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
					name = "my-cluster-network-subnet-reserved-ip"
					resource_type = "cluster_network_subnet_reserved_ip"
				}
				subnet {
					deleted {
						more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
					}
					href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
					id = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
					name = "my-cluster-network-subnet"
					resource_type = "cluster_network_subnet"
				}
			}
		}
	`, instanceID)
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentConfig(instanceID string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_instance_cluster_network_attachment" "is_instance_cluster_network_attachment_instance" {
			instance_id = "%s"
			before {
				href = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"
				id = "0717-fb880975-db45-4459-8548-64e3995ac213"
			}
			cluster_network_interface {
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
				id = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
				name = "my-cluster-network-interface"
				primary_ip {
					address = "10.1.0.6"
					deleted {
						more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
					}
					href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
					id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
					name = "my-cluster-network-subnet-reserved-ip"
					resource_type = "cluster_network_subnet_reserved_ip"
				}
				subnet {
					deleted {
						more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
					}
					href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
					id = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
					name = "my-cluster-network-subnet"
					resource_type = "cluster_network_subnet"
				}
			}
			name = "%s"
		}
	`, instanceID, name)
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentExists(n string, obj vpcv1.InstanceClusterNetworkAttachment) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getInstanceClusterNetworkAttachmentOptions := &vpcv1.GetInstanceClusterNetworkAttachmentOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getInstanceClusterNetworkAttachmentOptions.SetInstanceID(parts[0])
		getInstanceClusterNetworkAttachmentOptions.SetID(parts[1])

		instanceClusterNetworkAttachment, _, err := vpcClient.GetInstanceClusterNetworkAttachment(getInstanceClusterNetworkAttachmentOptions)
		if err != nil {
			return err
		}

		obj = *instanceClusterNetworkAttachment
		return nil
	}
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_instance_cluster_network_attachment" {
			continue
		}

		getInstanceClusterNetworkAttachmentOptions := &vpcv1.GetInstanceClusterNetworkAttachmentOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getInstanceClusterNetworkAttachmentOptions.SetInstanceID(parts[0])
		getInstanceClusterNetworkAttachmentOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetInstanceClusterNetworkAttachment(getInstanceClusterNetworkAttachmentOptions)

		if err == nil {
			return fmt.Errorf("InstanceClusterNetworkAttachment still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for InstanceClusterNetworkAttachment (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0717-fb880975-db45-4459-8548-64e3995ac213"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.InstanceClusterNetworkAttachmentBefore)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0717-fb880975-db45-4459-8548-64e3995ac213")
	model.Name = core.StringPtr("my-instance-network-attachment")
	model.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		clusterNetworkSubnetReservedIPReferenceModel := make(map[string]interface{})
		clusterNetworkSubnetReservedIPReferenceModel["address"] = "10.1.0.6"
		clusterNetworkSubnetReservedIPReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		clusterNetworkSubnetReservedIPReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/63341ffa-1037-4b50-be40-676e3e9ac0c7"
		clusterNetworkSubnetReservedIPReferenceModel["id"] = "63341ffa-1037-4b50-be40-676e3e9ac0c7"
		clusterNetworkSubnetReservedIPReferenceModel["name"] = "my-cluster-network-subnet-reserved-ip"
		clusterNetworkSubnetReservedIPReferenceModel["resource_type"] = "cluster_network_subnet_reserved_ip"

		clusterNetworkSubnetReferenceModel := make(map[string]interface{})
		clusterNetworkSubnetReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		clusterNetworkSubnetReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		clusterNetworkSubnetReferenceModel["id"] = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		clusterNetworkSubnetReferenceModel["name"] = "my-cluster-network-subnet"
		clusterNetworkSubnetReferenceModel["resource_type"] = "cluster_network_subnet"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["id"] = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["name"] = "my-cluster-network-interface"
		model["primary_ip"] = []map[string]interface{}{clusterNetworkSubnetReservedIPReferenceModel}
		model["resource_type"] = "cluster_network_interface"
		model["subnet"] = []map[string]interface{}{clusterNetworkSubnetReferenceModel}

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	clusterNetworkSubnetReservedIPReferenceModel := new(vpcv1.ClusterNetworkSubnetReservedIPReference)
	clusterNetworkSubnetReservedIPReferenceModel.Address = core.StringPtr("10.1.0.6")
	clusterNetworkSubnetReservedIPReferenceModel.Deleted = deletedModel
	clusterNetworkSubnetReservedIPReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/63341ffa-1037-4b50-be40-676e3e9ac0c7")
	clusterNetworkSubnetReservedIPReferenceModel.ID = core.StringPtr("63341ffa-1037-4b50-be40-676e3e9ac0c7")
	clusterNetworkSubnetReservedIPReferenceModel.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")
	clusterNetworkSubnetReservedIPReferenceModel.ResourceType = core.StringPtr("cluster_network_subnet_reserved_ip")

	clusterNetworkSubnetReferenceModel := new(vpcv1.ClusterNetworkSubnetReference)
	clusterNetworkSubnetReferenceModel.Deleted = deletedModel
	clusterNetworkSubnetReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	clusterNetworkSubnetReferenceModel.ID = core.StringPtr("0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	clusterNetworkSubnetReferenceModel.Name = core.StringPtr("my-cluster-network-subnet")
	clusterNetworkSubnetReferenceModel.ResourceType = core.StringPtr("cluster_network_subnet")

	model := new(vpcv1.ClusterNetworkInterfaceReference)
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.ID = core.StringPtr("0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.Name = core.StringPtr("my-cluster-network-interface")
	model.PrimaryIP = clusterNetworkSubnetReservedIPReferenceModel
	model.ResourceType = core.StringPtr("cluster_network_interface")
	model.Subnet = clusterNetworkSubnetReferenceModel

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "resource_suspended_by_provider"
		model["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.InstanceClusterNetworkAttachmentLifecycleReason)
	model.Code = core.StringPtr("resource_suspended_by_provider")
	model.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(t *testing.T) {
	checkResult := func(result vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceIntf) {
		clusterNetworkInterfacePrimaryIPPrototypeModel := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext)
		clusterNetworkInterfacePrimaryIPPrototypeModel.Address = core.StringPtr("10.0.0.5")
		clusterNetworkInterfacePrimaryIPPrototypeModel.AutoDelete = core.BoolPtr(false)
		clusterNetworkInterfacePrimaryIPPrototypeModel.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")

		clusterNetworkSubnetIdentityModel := new(vpcv1.ClusterNetworkSubnetIdentityByID)
		clusterNetworkSubnetIdentityModel.ID = core.StringPtr("7ec86020-1c6e-4889-b3f0-a15f2e50f87e")

		model := new(vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface)
		model.AutoDelete = core.BoolPtr(false)
		model.Name = core.StringPtr("my-cluster-network-interface")
		model.PrimaryIP = clusterNetworkInterfacePrimaryIPPrototypeModel
		model.Subnet = clusterNetworkSubnetIdentityModel
		model.ID = core.StringPtr("0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")

		assert.Equal(t, result, model)
	}

	clusterNetworkInterfacePrimaryIPPrototypeModel := make(map[string]interface{})
	clusterNetworkInterfacePrimaryIPPrototypeModel["address"] = "10.0.0.5"
	clusterNetworkInterfacePrimaryIPPrototypeModel["auto_delete"] = false
	clusterNetworkInterfacePrimaryIPPrototypeModel["name"] = "my-cluster-network-subnet-reserved-ip"

	clusterNetworkSubnetIdentityModel := make(map[string]interface{})
	clusterNetworkSubnetIdentityModel["id"] = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"

	model := make(map[string]interface{})
	model["auto_delete"] = false
	model["name"] = "my-cluster-network-interface"
	model["primary_ip"] = []interface{}{clusterNetworkInterfacePrimaryIPPrototypeModel}
	model["subnet"] = []interface{}{clusterNetworkSubnetIdentityModel}
	model["id"] = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterface(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototype(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext(t *testing.T) {
	checkResult := func(result vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextIntf) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext)
		model.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContext(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID)
		model.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref) {
		model := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPIdentityClusterNetworkInterfacePrimaryIPContextByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentity(t *testing.T) {
	checkResult := func(result vpcv1.ClusterNetworkSubnetIdentityIntf) {
		model := new(vpcv1.ClusterNetworkSubnetIdentity)
		model.ID = core.StringPtr("0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkSubnetIdentityByID) {
		model := new(vpcv1.ClusterNetworkSubnetIdentityByID)
		model.ID = core.StringPtr("0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.ClusterNetworkSubnetIdentityByHref) {
		model := new(vpcv1.ClusterNetworkSubnetIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToClusterNetworkSubnetIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment(t *testing.T) {
	checkResult := func(result *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment) {
		clusterNetworkInterfacePrimaryIPPrototypeModel := new(vpcv1.ClusterNetworkInterfacePrimaryIPPrototypeClusterNetworkSubnetReservedIPPrototypeClusterNetworkInterfacePrimaryIPContext)
		clusterNetworkInterfacePrimaryIPPrototypeModel.Address = core.StringPtr("10.0.0.5")
		clusterNetworkInterfacePrimaryIPPrototypeModel.AutoDelete = core.BoolPtr(false)
		clusterNetworkInterfacePrimaryIPPrototypeModel.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")

		clusterNetworkSubnetIdentityModel := new(vpcv1.ClusterNetworkSubnetIdentityByID)
		clusterNetworkSubnetIdentityModel.ID = core.StringPtr("0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")

		model := new(vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment)
		model.AutoDelete = core.BoolPtr(false)
		model.Name = core.StringPtr("my-cluster-network-interface")
		model.PrimaryIP = clusterNetworkInterfacePrimaryIPPrototypeModel
		model.Subnet = clusterNetworkSubnetIdentityModel

		assert.Equal(t, result, model)
	}

	clusterNetworkInterfacePrimaryIPPrototypeModel := make(map[string]interface{})
	clusterNetworkInterfacePrimaryIPPrototypeModel["address"] = "10.0.0.5"
	clusterNetworkInterfacePrimaryIPPrototypeModel["auto_delete"] = false
	clusterNetworkInterfacePrimaryIPPrototypeModel["name"] = "my-cluster-network-subnet-reserved-ip"

	clusterNetworkSubnetIdentityModel := make(map[string]interface{})
	clusterNetworkSubnetIdentityModel["id"] = "0717-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"

	model := make(map[string]interface{})
	model["auto_delete"] = false
	model["name"] = "my-cluster-network-interface"
	model["primary_ip"] = []interface{}{clusterNetworkInterfacePrimaryIPPrototypeModel}
	model["subnet"] = []interface{}{clusterNetworkSubnetIdentityModel}

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceInstanceClusterNetworkInterfacePrototypeInstanceClusterNetworkAttachment(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity(t *testing.T) {
	checkResult := func(result vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityIntf) {
		model := new(vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity)
		model.ID = core.StringPtr("0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentity(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID) {
		model := new(vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID)
		model.ID = core.StringPtr("0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref) {
		model := new(vpcv1.InstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentPrototypeClusterNetworkInterfaceClusterNetworkInterfaceIdentityClusterNetworkInterfaceIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototype(t *testing.T) {
	checkResult := func(result vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeIntf) {
		model := new(vpcv1.InstanceClusterNetworkAttachmentBeforePrototype)
		model.ID = core.StringPtr("0717-fb880975-db45-4459-8548-64e3995ac213")
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-fb880975-db45-4459-8548-64e3995ac213"
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototype(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByID(t *testing.T) {
	checkResult := func(result *vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByID) {
		model := new(vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByID)
		model.ID = core.StringPtr("0717-fb880975-db45-4459-8548-64e3995ac213")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["id"] = "0717-fb880975-db45-4459-8548-64e3995ac213"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByID(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref(t *testing.T) {
	checkResult := func(result *vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref) {
		model := new(vpcv1.InstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref)
		model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213")

		assert.Equal(t, result, model)
	}

	model := make(map[string]interface{})
	model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"

	result, err := vpc.ResourceIBMIsInstanceClusterNetworkAttachmentMapToInstanceClusterNetworkAttachmentBeforePrototypeInstanceClusterNetworkAttachmentIdentityByHref(model)
	assert.Nil(t, err)
	checkResult(result)
}
