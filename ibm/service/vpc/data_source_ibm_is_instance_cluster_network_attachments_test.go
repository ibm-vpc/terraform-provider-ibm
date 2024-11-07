// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsInstanceClusterNetworkAttachmentsDataSourceBasic(t *testing.T) {
	instanceClusterNetworkAttachmentInstanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceClusterNetworkAttachmentsDataSourceConfigBasic(instanceClusterNetworkAttachmentInstanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.#"),
				),
			},
		},
	})
}

func TestAccIBMIsInstanceClusterNetworkAttachmentsDataSourceAllArgs(t *testing.T) {
	instanceClusterNetworkAttachmentInstanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))
	instanceClusterNetworkAttachmentName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceClusterNetworkAttachmentsDataSourceConfig(instanceClusterNetworkAttachmentInstanceID, instanceClusterNetworkAttachmentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.0.lifecycle_state"),
					resource.TestCheckResourceAttr("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.0.name", instanceClusterNetworkAttachmentName),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.0.resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentsDataSourceConfigBasic(instanceClusterNetworkAttachmentInstanceID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_instance_cluster_network_attachment" "is_instance_cluster_network_attachment_instance" {
			instance_id = "%s"
			cluster_network_interface {
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
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
				resource_type = "cluster_network_interface"
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

		data "ibm_is_instance_cluster_network_attachments" "is_instance_cluster_network_attachments_instance" {
			instance_id = ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance.instance_id
		}
	`, instanceClusterNetworkAttachmentInstanceID)
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentsDataSourceConfig(instanceClusterNetworkAttachmentInstanceID string, instanceClusterNetworkAttachmentName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_instance_cluster_network_attachment" "is_instance_cluster_network_attachment_instance" {
			instance_id = "%s"
			before {
				href = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"
				id = "0717-fb880975-db45-4459-8548-64e3995ac213"
				name = "my-instance-network-attachment"
				resource_type = "instance_cluster_network_attachment"
			}
			cluster_network_interface {
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
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
				resource_type = "cluster_network_interface"
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

		data "ibm_is_instance_cluster_network_attachments" "is_instance_cluster_network_attachments_instance" {
			instance_id = ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance.instance_id
		}
	`, instanceClusterNetworkAttachmentInstanceID, instanceClusterNetworkAttachmentName)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		instanceClusterNetworkAttachmentBeforeModel := make(map[string]interface{})
		instanceClusterNetworkAttachmentBeforeModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-a69563fa-0415-4d6e-aeb3-a3f14654bf90"
		instanceClusterNetworkAttachmentBeforeModel["id"] = "a69563fa-0415-4d6e-aeb3-a3f14654bf90"
		instanceClusterNetworkAttachmentBeforeModel["name"] = "other-instance-cluster-network-attachment"
		instanceClusterNetworkAttachmentBeforeModel["resource_type"] = "instance_cluster_network_attachment"

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

		clusterNetworkInterfaceReferenceModel := make(map[string]interface{})
		clusterNetworkInterfaceReferenceModel["deleted"] = []map[string]interface{}{deletedModel}
		clusterNetworkInterfaceReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		clusterNetworkInterfaceReferenceModel["id"] = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		clusterNetworkInterfaceReferenceModel["name"] = "my-cluster-network-interface"
		clusterNetworkInterfaceReferenceModel["primary_ip"] = []map[string]interface{}{clusterNetworkSubnetReservedIPReferenceModel}
		clusterNetworkInterfaceReferenceModel["resource_type"] = "cluster_network_interface"
		clusterNetworkInterfaceReferenceModel["subnet"] = []map[string]interface{}{clusterNetworkSubnetReferenceModel}

		instanceClusterNetworkAttachmentLifecycleReasonModel := make(map[string]interface{})
		instanceClusterNetworkAttachmentLifecycleReasonModel["code"] = "resource_suspended_by_provider"
		instanceClusterNetworkAttachmentLifecycleReasonModel["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		instanceClusterNetworkAttachmentLifecycleReasonModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		model := make(map[string]interface{})
		model["before"] = []map[string]interface{}{instanceClusterNetworkAttachmentBeforeModel}
		model["cluster_network_interface"] = []map[string]interface{}{clusterNetworkInterfaceReferenceModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0717-fb880975-db45-4459-8548-64e3995ac213"
		model["lifecycle_reasons"] = []map[string]interface{}{instanceClusterNetworkAttachmentLifecycleReasonModel}
		model["lifecycle_state"] = "stable"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	instanceClusterNetworkAttachmentBeforeModel := new(vpcv1.InstanceClusterNetworkAttachmentBefore)
	instanceClusterNetworkAttachmentBeforeModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-a69563fa-0415-4d6e-aeb3-a3f14654bf90")
	instanceClusterNetworkAttachmentBeforeModel.ID = core.StringPtr("a69563fa-0415-4d6e-aeb3-a3f14654bf90")
	instanceClusterNetworkAttachmentBeforeModel.Name = core.StringPtr("other-instance-cluster-network-attachment")
	instanceClusterNetworkAttachmentBeforeModel.ResourceType = core.StringPtr("instance_cluster_network_attachment")

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

	clusterNetworkInterfaceReferenceModel := new(vpcv1.ClusterNetworkInterfaceReference)
	clusterNetworkInterfaceReferenceModel.Deleted = deletedModel
	clusterNetworkInterfaceReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	clusterNetworkInterfaceReferenceModel.ID = core.StringPtr("0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	clusterNetworkInterfaceReferenceModel.Name = core.StringPtr("my-cluster-network-interface")
	clusterNetworkInterfaceReferenceModel.PrimaryIP = clusterNetworkSubnetReservedIPReferenceModel
	clusterNetworkInterfaceReferenceModel.ResourceType = core.StringPtr("cluster_network_interface")
	clusterNetworkInterfaceReferenceModel.Subnet = clusterNetworkSubnetReferenceModel

	instanceClusterNetworkAttachmentLifecycleReasonModel := new(vpcv1.InstanceClusterNetworkAttachmentLifecycleReason)
	instanceClusterNetworkAttachmentLifecycleReasonModel.Code = core.StringPtr("resource_suspended_by_provider")
	instanceClusterNetworkAttachmentLifecycleReasonModel.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	instanceClusterNetworkAttachmentLifecycleReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	model := new(vpcv1.InstanceClusterNetworkAttachment)
	model.Before = instanceClusterNetworkAttachmentBeforeModel
	model.ClusterNetworkInterface = clusterNetworkInterfaceReferenceModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/0717_e21b7391-2ca2-4ab5-84a8-b92157a633b0/cluster_network_attachments/0717-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0717-fb880975-db45-4459-8548-64e3995ac213")
	model.LifecycleReasons = []vpcv1.InstanceClusterNetworkAttachmentLifecycleReason{*instanceClusterNetworkAttachmentLifecycleReasonModel}
	model.LifecycleState = core.StringPtr("stable")
	model.Name = core.StringPtr("my-instance-network-attachment")
	model.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentBeforeToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentBeforeToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentLifecycleReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
