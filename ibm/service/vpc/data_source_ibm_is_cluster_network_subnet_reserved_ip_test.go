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

func TestAccIBMIsClusterNetworkSubnetReservedIPDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetReservedIPDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_subnet_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "is_cluster_network_subnet_reserved_ip_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "address"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "lifecycle_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "owner"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkSubnetReservedIPDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_cluster_network_subnet_reserved_ip" "is_cluster_network_subnet_reserved_ip_instance" {
			cluster_network_id = "cluster_network_id"
			cluster_network_subnet_id = "cluster_network_subnet_id"
			id = "id"
		}
	`)
}

func TestDataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "resource_suspended_by_provider"
		model["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetReservedIPLifecycleReason)
	model.Code = core.StringPtr("resource_suspended_by_provider")
	model.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	result, err := vpc.DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel := make(map[string]interface{})
		clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["id"] = "0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["name"] = "my-cluster-network-interface"
		model["resource_type"] = "cluster_network_interface"

		assert.Equal(t, result, model)
	}

	clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel := new(vpcv1.ClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeleted)
	clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReservedIPTarget)
	model.Deleted = clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.ID = core.StringPtr("0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.Name = core.StringPtr("my-cluster-network-interface")
	model.ResourceType = core.StringPtr("cluster_network_interface")

	result, err := vpc.DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel := make(map[string]interface{})
		clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["id"] = "0767-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["name"] = "my-cluster-network-interface"
		model["resource_type"] = "cluster_network_interface"

		assert.Equal(t, result, model)
	}

	clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel := new(vpcv1.ClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeleted)
	clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext)
	model.Deleted = clusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextDeletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.ID = core.StringPtr("0767-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.Name = core.StringPtr("my-cluster-network-interface")
	model.ResourceType = core.StringPtr("cluster_network_interface")

	result, err := vpc.DataSourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
