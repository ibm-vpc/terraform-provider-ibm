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

func TestAccIBMIsInstanceClusterNetworkAttachmentDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceClusterNetworkAttachmentDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "is_instance_cluster_network_attachment_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "cluster_network_interface.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "lifecycle_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachment.is_instance_cluster_network_attachment_instance", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_instance_cluster_network_attachment" "is_instance_cluster_network_attachment_instance" {
			instance_id = "instance_id"
			id = "id"
		}
	`)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0767-fb880975-db45-4459-8548-64e3995ac213"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.InstanceClusterNetworkAttachmentBefore)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0767-fb880975-db45-4459-8548-64e3995ac213")
	model.Name = core.StringPtr("my-instance-network-attachment")
	model.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentBeforeToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		clusterNetworkInterfaceReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkInterfaceReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		clusterNetworkSubnetReservedIPReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkSubnetReservedIPReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		clusterNetworkSubnetReservedIPReferenceModel := make(map[string]interface{})
		clusterNetworkSubnetReservedIPReferenceModel["address"] = "10.1.0.6"
		clusterNetworkSubnetReservedIPReferenceModel["deleted"] = []map[string]interface{}{clusterNetworkSubnetReservedIPReferenceDeletedModel}
		clusterNetworkSubnetReservedIPReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/3f6241cd-5abc-4afb-aaa2-058661f969d8/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		clusterNetworkSubnetReservedIPReferenceModel["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		clusterNetworkSubnetReservedIPReferenceModel["name"] = "my-cluster-network-subnet-reserved-ip"
		clusterNetworkSubnetReservedIPReferenceModel["resource_type"] = "cluster_network_subnet_reserved_ip"

		clusterNetworkSubnetReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkSubnetReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		clusterNetworkSubnetReferenceModel := make(map[string]interface{})
		clusterNetworkSubnetReferenceModel["deleted"] = []map[string]interface{}{clusterNetworkSubnetReferenceDeletedModel}
		clusterNetworkSubnetReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		clusterNetworkSubnetReferenceModel["id"] = "0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		clusterNetworkSubnetReferenceModel["name"] = "my-cluster-network-subnet"
		clusterNetworkSubnetReferenceModel["resource_type"] = "cluster_network_subnet"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{clusterNetworkInterfaceReferenceDeletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["id"] = "0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["name"] = "my-cluster-network-interface"
		model["primary_ip"] = []map[string]interface{}{clusterNetworkSubnetReservedIPReferenceModel}
		model["resource_type"] = "cluster_network_interface"
		model["subnet"] = []map[string]interface{}{clusterNetworkSubnetReferenceModel}

		assert.Equal(t, result, model)
	}

	clusterNetworkInterfaceReferenceDeletedModel := new(vpcv1.ClusterNetworkInterfaceReferenceDeleted)
	clusterNetworkInterfaceReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	clusterNetworkSubnetReservedIPReferenceDeletedModel := new(vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted)
	clusterNetworkSubnetReservedIPReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	clusterNetworkSubnetReservedIPReferenceModel := new(vpcv1.ClusterNetworkSubnetReservedIPReference)
	clusterNetworkSubnetReservedIPReferenceModel.Address = core.StringPtr("10.1.0.6")
	clusterNetworkSubnetReservedIPReferenceModel.Deleted = clusterNetworkSubnetReservedIPReferenceDeletedModel
	clusterNetworkSubnetReservedIPReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/3f6241cd-5abc-4afb-aaa2-058661f969d8/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	clusterNetworkSubnetReservedIPReferenceModel.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	clusterNetworkSubnetReservedIPReferenceModel.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")
	clusterNetworkSubnetReservedIPReferenceModel.ResourceType = core.StringPtr("cluster_network_subnet_reserved_ip")

	clusterNetworkSubnetReferenceDeletedModel := new(vpcv1.ClusterNetworkSubnetReferenceDeleted)
	clusterNetworkSubnetReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	clusterNetworkSubnetReferenceModel := new(vpcv1.ClusterNetworkSubnetReference)
	clusterNetworkSubnetReferenceModel.Deleted = clusterNetworkSubnetReferenceDeletedModel
	clusterNetworkSubnetReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	clusterNetworkSubnetReferenceModel.ID = core.StringPtr("0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	clusterNetworkSubnetReferenceModel.Name = core.StringPtr("my-cluster-network-subnet")
	clusterNetworkSubnetReferenceModel.ResourceType = core.StringPtr("cluster_network_subnet")

	model := new(vpcv1.ClusterNetworkInterfaceReference)
	model.Deleted = clusterNetworkInterfaceReferenceDeletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.ID = core.StringPtr("0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.Name = core.StringPtr("my-cluster-network-interface")
	model.PrimaryIP = clusterNetworkSubnetReservedIPReferenceModel
	model.ResourceType = core.StringPtr("cluster_network_interface")
	model.Subnet = clusterNetworkSubnetReferenceModel

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkInterfaceReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		clusterNetworkSubnetReservedIPReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkSubnetReservedIPReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["address"] = "10.1.0.6"
		model["deleted"] = []map[string]interface{}{clusterNetworkSubnetReservedIPReferenceDeletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/3f6241cd-5abc-4afb-aaa2-058661f969d8/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		model["id"] = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
		model["name"] = "my-cluster-network-subnet-reserved-ip"
		model["resource_type"] = "cluster_network_subnet_reserved_ip"

		assert.Equal(t, result, model)
	}

	clusterNetworkSubnetReservedIPReferenceDeletedModel := new(vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted)
	clusterNetworkSubnetReservedIPReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReservedIPReference)
	model.Address = core.StringPtr("10.1.0.6")
	model.Deleted = clusterNetworkSubnetReservedIPReferenceDeletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/3f6241cd-5abc-4afb-aaa2-058661f969d8/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	model.ID = core.StringPtr("6d353a0f-aeb1-4ae1-832e-1110d10981bb")
	model.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip")
	model.ResourceType = core.StringPtr("cluster_network_subnet_reserved_ip")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReservedIPReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		clusterNetworkSubnetReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkSubnetReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{clusterNetworkSubnetReferenceDeletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		model["id"] = "0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		model["name"] = "my-cluster-network-subnet"
		model["resource_type"] = "cluster_network_subnet"

		assert.Equal(t, result, model)
	}

	clusterNetworkSubnetReferenceDeletedModel := new(vpcv1.ClusterNetworkSubnetReferenceDeleted)
	clusterNetworkSubnetReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReference)
	model.Deleted = clusterNetworkSubnetReferenceDeletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	model.ID = core.StringPtr("0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	model.Name = core.StringPtr("my-cluster-network-subnet")
	model.ResourceType = core.StringPtr("cluster_network_subnet")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentClusterNetworkSubnetReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentInstanceClusterNetworkAttachmentLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
