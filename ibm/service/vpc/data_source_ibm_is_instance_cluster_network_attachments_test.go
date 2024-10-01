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

func TestAccIBMIsInstanceClusterNetworkAttachmentsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceClusterNetworkAttachmentsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_cluster_network_attachments.is_instance_cluster_network_attachments_instance", "cluster_network_attachments.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsInstanceClusterNetworkAttachmentsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_instance_cluster_network_attachments" "is_instance_cluster_network_attachments_instance" {
			instance_id = "instance_id"
		}
	`)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		instanceClusterNetworkAttachmentBeforeModel := make(map[string]interface{})
		instanceClusterNetworkAttachmentBeforeModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/123a490a-9e64-4254-a93b-9a3af3ede270/cluster_network_attachments/a69563fa-0415-4d6e-aeb3-a3f14654bf90"
		instanceClusterNetworkAttachmentBeforeModel["id"] = "a69563fa-0415-4d6e-aeb3-a3f14654bf90"
		instanceClusterNetworkAttachmentBeforeModel["name"] = "other-instance-cluster-network-attachment"
		instanceClusterNetworkAttachmentBeforeModel["resource_type"] = "instance_cluster_network_attachment"

		clusterNetworkInterfaceReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkInterfaceReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		clusterNetworkSubnetReservedIPReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkSubnetReservedIPReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		clusterNetworkSubnetReservedIPReferenceModel := make(map[string]interface{})
		clusterNetworkSubnetReservedIPReferenceModel["address"] = "10.0.0.32"
		clusterNetworkSubnetReservedIPReferenceModel["deleted"] = []map[string]interface{}{clusterNetworkSubnetReservedIPReferenceDeletedModel}
		clusterNetworkSubnetReservedIPReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-1e8f34ac-f4c6-437c-adeb-8d698ce7eaa4/subnets/0767-b28a7e6d-a66b-4de7-8713-15dcffdce401/reserved_ips/0767-7768a27e-cd6c-4a13-a9e6-d67a964e54a5"
		clusterNetworkSubnetReservedIPReferenceModel["id"] = "0767-7768a27e-cd6c-4a13-a9e6-d67a964e54a5"
		clusterNetworkSubnetReservedIPReferenceModel["name"] = "my-cluster-network-subnet-reserved-ip-1"
		clusterNetworkSubnetReservedIPReferenceModel["resource_type"] = "cluster_network_subnet_reserved_ip"

		clusterNetworkSubnetReferenceDeletedModel := make(map[string]interface{})
		clusterNetworkSubnetReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		clusterNetworkSubnetReferenceModel := make(map[string]interface{})
		clusterNetworkSubnetReferenceModel["deleted"] = []map[string]interface{}{clusterNetworkSubnetReferenceDeletedModel}
		clusterNetworkSubnetReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-1e8f34ac-f4c6-437c-adeb-8d698ce7eaa4/subnets/9270d819-c05e-4352-99e4-80c4680cdb7c"
		clusterNetworkSubnetReferenceModel["id"] = "0767-9270d819-c05e-4352-99e4-80c4680cdb7c"
		clusterNetworkSubnetReferenceModel["name"] = "my-subnet"
		clusterNetworkSubnetReferenceModel["resource_type"] = "cluster_network_subnet"

		clusterNetworkInterfaceReferenceModel := make(map[string]interface{})
		clusterNetworkInterfaceReferenceModel["deleted"] = []map[string]interface{}{clusterNetworkInterfaceReferenceDeletedModel}
		clusterNetworkInterfaceReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/1e8f34ac-f4c6-437c-adeb-8d698ce7eaa4/interfaces/0767-54eb57ee-86f2-4796-90bb-d7874e0831ef"
		clusterNetworkInterfaceReferenceModel["id"] = "0767-54eb57ee-86f2-4796-90bb-d7874e0831ef"
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
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0767-fb880975-db45-4459-8548-64e3995ac213"
		model["lifecycle_reasons"] = []map[string]interface{}{instanceClusterNetworkAttachmentLifecycleReasonModel}
		model["lifecycle_state"] = "stable"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	instanceClusterNetworkAttachmentBeforeModel := new(vpcv1.InstanceClusterNetworkAttachmentBefore)
	instanceClusterNetworkAttachmentBeforeModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/123a490a-9e64-4254-a93b-9a3af3ede270/cluster_network_attachments/a69563fa-0415-4d6e-aeb3-a3f14654bf90")
	instanceClusterNetworkAttachmentBeforeModel.ID = core.StringPtr("a69563fa-0415-4d6e-aeb3-a3f14654bf90")
	instanceClusterNetworkAttachmentBeforeModel.Name = core.StringPtr("other-instance-cluster-network-attachment")
	instanceClusterNetworkAttachmentBeforeModel.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	clusterNetworkInterfaceReferenceDeletedModel := new(vpcv1.ClusterNetworkInterfaceReferenceDeleted)
	clusterNetworkInterfaceReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	clusterNetworkSubnetReservedIPReferenceDeletedModel := new(vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted)
	clusterNetworkSubnetReservedIPReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	clusterNetworkSubnetReservedIPReferenceModel := new(vpcv1.ClusterNetworkSubnetReservedIPReference)
	clusterNetworkSubnetReservedIPReferenceModel.Address = core.StringPtr("10.0.0.32")
	clusterNetworkSubnetReservedIPReferenceModel.Deleted = clusterNetworkSubnetReservedIPReferenceDeletedModel
	clusterNetworkSubnetReservedIPReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-1e8f34ac-f4c6-437c-adeb-8d698ce7eaa4/subnets/0767-b28a7e6d-a66b-4de7-8713-15dcffdce401/reserved_ips/0767-7768a27e-cd6c-4a13-a9e6-d67a964e54a5")
	clusterNetworkSubnetReservedIPReferenceModel.ID = core.StringPtr("0767-7768a27e-cd6c-4a13-a9e6-d67a964e54a5")
	clusterNetworkSubnetReservedIPReferenceModel.Name = core.StringPtr("my-cluster-network-subnet-reserved-ip-1")
	clusterNetworkSubnetReservedIPReferenceModel.ResourceType = core.StringPtr("cluster_network_subnet_reserved_ip")

	clusterNetworkSubnetReferenceDeletedModel := new(vpcv1.ClusterNetworkSubnetReferenceDeleted)
	clusterNetworkSubnetReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	clusterNetworkSubnetReferenceModel := new(vpcv1.ClusterNetworkSubnetReference)
	clusterNetworkSubnetReferenceModel.Deleted = clusterNetworkSubnetReferenceDeletedModel
	clusterNetworkSubnetReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-1e8f34ac-f4c6-437c-adeb-8d698ce7eaa4/subnets/9270d819-c05e-4352-99e4-80c4680cdb7c")
	clusterNetworkSubnetReferenceModel.ID = core.StringPtr("0767-9270d819-c05e-4352-99e4-80c4680cdb7c")
	clusterNetworkSubnetReferenceModel.Name = core.StringPtr("my-subnet")
	clusterNetworkSubnetReferenceModel.ResourceType = core.StringPtr("cluster_network_subnet")

	clusterNetworkInterfaceReferenceModel := new(vpcv1.ClusterNetworkInterfaceReference)
	clusterNetworkInterfaceReferenceModel.Deleted = clusterNetworkInterfaceReferenceDeletedModel
	clusterNetworkInterfaceReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/1e8f34ac-f4c6-437c-adeb-8d698ce7eaa4/interfaces/0767-54eb57ee-86f2-4796-90bb-d7874e0831ef")
	clusterNetworkInterfaceReferenceModel.ID = core.StringPtr("0767-54eb57ee-86f2-4796-90bb-d7874e0831ef")
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
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0767-fb880975-db45-4459-8548-64e3995ac213")
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsInstanceClusterNetworkAttachmentBeforeToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkInterfaceReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReservedIPReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsInstanceClusterNetworkAttachmentsClusterNetworkSubnetReferenceDeletedToMap(model)
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
