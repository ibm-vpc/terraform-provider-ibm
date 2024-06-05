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
	. "github.com/IBM-Cloud/terraform-provider-ibm/ibm/unittest"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsClusterNetworkInterfacesDataSourceBasic(t *testing.T) {
	clusterNetworkInterfaceClusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkInterfacesDataSourceConfigBasic(clusterNetworkInterfaceClusterNetworkID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "cluster_network_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.#"),
				),
			},
		},
	})
}

func TestAccIBMIsClusterNetworkInterfacesDataSourceAllArgs(t *testing.T) {
	clusterNetworkInterfaceClusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))
	clusterNetworkInterfaceName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkInterfacesDataSourceConfig(clusterNetworkInterfaceClusterNetworkID, clusterNetworkInterfaceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "cluster_network_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "sort"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.allow_ip_spoofing"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.enable_infrastructure_nat"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.mac_address"),
					resource.TestCheckResourceAttr("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.name", clusterNetworkInterfaceName),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.protocol_state_filtering_mode"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_interfaces.is_cluster_network_interfaces_instance", "interfaces.0.resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkInterfacesDataSourceConfigBasic(clusterNetworkInterfaceClusterNetworkID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_cluster_network_interface" "is_cluster_network_interface_instance" {
			cluster_network_id = "%s"
		}

		data "ibm_is_cluster_network_interfaces" "is_cluster_network_interfaces_instance" {
			cluster_network_id = ibm_is_cluster_network_interface.is_cluster_network_interface_instance.cluster_network_id
			name = ibm_is_cluster_network_interface.is_cluster_network_interface_instance.name
			sort = "name"
		}
	`, clusterNetworkInterfaceClusterNetworkID)
}

func testAccCheckIBMIsClusterNetworkInterfacesDataSourceConfig(clusterNetworkInterfaceClusterNetworkID string, clusterNetworkInterfaceName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_cluster_network_interface" "is_cluster_network_interface_instance" {
			cluster_network_id = "%s"
			name = "%s"
			primary_ip {
				address = "10.1.0.6"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/3f6241cd-5abc-4afb-aaa2-058661f969d8/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
				name = "my-cluster-network-subnet-reserved-ip"
				resource_type = "cluster_network_subnet_reserved_ip"
			}
			subnet {
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
				id = "0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
				name = "my-cluster-network-subnet"
				resource_type = "cluster_network_subnet"
			}
		}

		data "ibm_is_cluster_network_interfaces" "is_cluster_network_interfaces_instance" {
			cluster_network_id = ibm_is_cluster_network_interface.is_cluster_network_interface_instance.cluster_network_id
			name = ibm_is_cluster_network_interface.is_cluster_network_interface_instance.name
			sort = "name"
		}
	`, clusterNetworkInterfaceClusterNetworkID, clusterNetworkInterfaceName)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		clusterNetworkInterfaceLifecycleReasonModel := make(map[string]interface{})
		clusterNetworkInterfaceLifecycleReasonModel["code"] = "resource_suspended_by_provider"
		clusterNetworkInterfaceLifecycleReasonModel["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		clusterNetworkInterfaceLifecycleReasonModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

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

		clusterNetworkInterfaceTargetModel := make(map[string]interface{})
		clusterNetworkInterfaceTargetModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213"
		clusterNetworkInterfaceTargetModel["id"] = "0767-fb880975-db45-4459-8548-64e3995ac213"
		clusterNetworkInterfaceTargetModel["name"] = "my-instance-network-attachment"
		clusterNetworkInterfaceTargetModel["resource_type"] = "instance_cluster_network_attachment"

		vpcReferenceDeletedModel := make(map[string]interface{})
		vpcReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		vpcReferenceModel := make(map[string]interface{})
		vpcReferenceModel["crn"] = "crn:[...]"
		vpcReferenceModel["deleted"] = []map[string]interface{}{vpcReferenceDeletedModel}
		vpcReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/0767-a0819609-0997-4f92-9409-86c95ddf59d3"
		vpcReferenceModel["id"] = "0767-a0819609-0997-4f92-9409-86c95ddf59d3"
		vpcReferenceModel["name"] = "my-vpc"
		vpcReferenceModel["resource_type"] = "vpc"

		zoneReferenceModel := make(map[string]interface{})
		zoneReferenceModel["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		zoneReferenceModel["name"] = "us-south-1"

		model := make(map[string]interface{})
		model["allow_ip_spoofing"] = true
		model["auto_delete"] = false
		model["created_at"] = "2019-01-01T12:00:00.000Z"
		model["enable_infrastructure_nat"] = false
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["id"] = "0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["lifecycle_reasons"] = []map[string]interface{}{clusterNetworkInterfaceLifecycleReasonModel}
		model["lifecycle_state"] = "stable"
		model["mac_address"] = "02:00:4D:45:45:4D"
		model["name"] = "my-cluster-network-interface"
		model["primary_ip"] = []map[string]interface{}{clusterNetworkSubnetReservedIPReferenceModel}
		model["protocol_state_filtering_mode"] = "enabled"
		model["resource_type"] = "cluster_network_interface"
		model["subnet"] = []map[string]interface{}{clusterNetworkSubnetReferenceModel}
		model["target"] = []map[string]interface{}{clusterNetworkInterfaceTargetModel}
		model["vpc"] = []map[string]interface{}{vpcReferenceModel}
		model["zone"] = []map[string]interface{}{zoneReferenceModel}

		assert.Equal(t, result, model)
	}

	clusterNetworkInterfaceLifecycleReasonModel := new(vpcv1.ClusterNetworkInterfaceLifecycleReason)
	clusterNetworkInterfaceLifecycleReasonModel.Code = core.StringPtr("resource_suspended_by_provider")
	clusterNetworkInterfaceLifecycleReasonModel.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	clusterNetworkInterfaceLifecycleReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

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

	clusterNetworkInterfaceTargetModel := new(vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReference)
	clusterNetworkInterfaceTargetModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213")
	clusterNetworkInterfaceTargetModel.ID = core.StringPtr("0767-fb880975-db45-4459-8548-64e3995ac213")
	clusterNetworkInterfaceTargetModel.Name = core.StringPtr("my-instance-network-attachment")
	clusterNetworkInterfaceTargetModel.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	vpcReferenceDeletedModel := new(vpcv1.VPCReferenceDeleted)
	vpcReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	vpcReferenceModel := new(vpcv1.VPCReference)
	vpcReferenceModel.CRN = core.StringPtr("crn:[...]")
	vpcReferenceModel.Deleted = vpcReferenceDeletedModel
	vpcReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/vpcs/0767-a0819609-0997-4f92-9409-86c95ddf59d3")
	vpcReferenceModel.ID = core.StringPtr("0767-a0819609-0997-4f92-9409-86c95ddf59d3")
	vpcReferenceModel.Name = core.StringPtr("my-vpc")
	vpcReferenceModel.ResourceType = core.StringPtr("vpc")

	zoneReferenceModel := new(vpcv1.ZoneReference)
	zoneReferenceModel.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	zoneReferenceModel.Name = core.StringPtr("us-south-1")

	model := new(vpcv1.ClusterNetworkInterface)
	model.AllowIPSpoofing = core.BoolPtr(true)
	model.AutoDelete = core.BoolPtr(false)
	model.CreatedAt = CreateMockDateTime("2019-01-01T12:00:00.000Z")
	model.EnableInfrastructureNat = core.BoolPtr(false)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.ID = core.StringPtr("0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.LifecycleReasons = []vpcv1.ClusterNetworkInterfaceLifecycleReason{*clusterNetworkInterfaceLifecycleReasonModel}
	model.LifecycleState = core.StringPtr("stable")
	model.MacAddress = core.StringPtr("02:00:4D:45:45:4D")
	model.Name = core.StringPtr("my-cluster-network-interface")
	model.PrimaryIP = clusterNetworkSubnetReservedIPReferenceModel
	model.ProtocolStateFilteringMode = core.StringPtr("enabled")
	model.ResourceType = core.StringPtr("cluster_network_interface")
	model.Subnet = clusterNetworkSubnetReferenceModel
	model.Target = clusterNetworkInterfaceTargetModel
	model.VPC = vpcReferenceModel
	model.Zone = zoneReferenceModel

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceLifecycleReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReservedIPReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReservedIPReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReservedIPReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetReservedIPReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReservedIPReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReferenceToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkSubnetReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceTargetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0767-fb880975-db45-4459-8548-64e3995ac213"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceTarget)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0767-fb880975-db45-4459-8548-64e3995ac213")
	model.Name = core.StringPtr("my-instance-network-attachment")
	model.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceTargetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213"
		model["id"] = "0767-fb880975-db45-4459-8548-64e3995ac213"
		model["name"] = "my-instance-network-attachment"
		model["resource_type"] = "instance_cluster_network_attachment"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/instances/5dd61d72-acaa-47c2-a336-3d849660d010/cluster_network_attachments/0767-fb880975-db45-4459-8548-64e3995ac213")
	model.ID = core.StringPtr("0767-fb880975-db45-4459-8548-64e3995ac213")
	model.Name = core.StringPtr("my-instance-network-attachment")
	model.ResourceType = core.StringPtr("instance_cluster_network_attachment")

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesClusterNetworkInterfaceTargetInstanceClusterNetworkAttachmentReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesVPCReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		vpcReferenceDeletedModel := make(map[string]interface{})
		vpcReferenceDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["crn"] = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		model["deleted"] = []map[string]interface{}{vpcReferenceDeletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		model["id"] = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
		model["name"] = "my-vpc"
		model["resource_type"] = "vpc"

		assert.Equal(t, result, model)
	}

	vpcReferenceDeletedModel := new(vpcv1.VPCReferenceDeleted)
	vpcReferenceDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.VPCReference)
	model.CRN = core.StringPtr("crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
	model.Deleted = vpcReferenceDeletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
	model.ID = core.StringPtr("r006-4727d842-f94f-4a2d-824a-9bc9b02c523b")
	model.Name = core.StringPtr("my-vpc")
	model.ResourceType = core.StringPtr("vpc")

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesVPCReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesVPCReferenceDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPCReferenceDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesVPCReferenceDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkInterfacesZoneReferenceToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
		model["name"] = "us-south-1"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ZoneReference)
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1")
	model.Name = core.StringPtr("us-south-1")

	result, err := vpc.DataSourceIBMIsClusterNetworkInterfacesZoneReferenceToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
