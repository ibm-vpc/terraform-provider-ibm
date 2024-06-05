// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	. "github.com/IBM-Cloud/terraform-provider-ibm/ibm/unittest"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsClusterNetworkSubnetsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnets.is_cluster_network_subnets_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnets.is_cluster_network_subnets_instance", "cluster_network_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnets.is_cluster_network_subnets_instance", "subnets.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkSubnetsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_cluster_network_subnets" "is_cluster_network_subnets_instance" {
			cluster_network_id = "cluster_network_id"
			name = "name"
			sort = "name"
		}
	`)
}

func TestDataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		clusterNetworkSubnetLifecycleReasonModel := make(map[string]interface{})
		clusterNetworkSubnetLifecycleReasonModel["code"] = "resource_suspended_by_provider"
		clusterNetworkSubnetLifecycleReasonModel["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		clusterNetworkSubnetLifecycleReasonModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		model := make(map[string]interface{})
		model["available_ipv4_address_count"] = int(15)
		model["created_at"] = "2019-01-01T12:00:00.000Z"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		model["id"] = "0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930"
		model["ip_version"] = "ipv4"
		model["ipv4_cidr_block"] = "10.0.0.0/24"
		model["lifecycle_reasons"] = []map[string]interface{}{clusterNetworkSubnetLifecycleReasonModel}
		model["lifecycle_state"] = "stable"
		model["name"] = "my-cluster-network-subnet"
		model["resource_type"] = "cluster_network_subnet"
		model["total_ipv4_address_count"] = int(256)

		assert.Equal(t, result, model)
	}

	clusterNetworkSubnetLifecycleReasonModel := new(vpcv1.ClusterNetworkSubnetLifecycleReason)
	clusterNetworkSubnetLifecycleReasonModel.Code = core.StringPtr("resource_suspended_by_provider")
	clusterNetworkSubnetLifecycleReasonModel.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	clusterNetworkSubnetLifecycleReasonModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	model := new(vpcv1.ClusterNetworkSubnet)
	model.AvailableIpv4AddressCount = core.Int64Ptr(int64(15))
	model.CreatedAt = CreateMockDateTime("2019-01-01T12:00:00.000Z")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0767-da0df18c-7598-4633-a648-fdaac28a5573/subnets/0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	model.ID = core.StringPtr("0767-7931845c-65c4-4b0a-80cd-7d9c1a6d7930")
	model.IPVersion = core.StringPtr("ipv4")
	model.Ipv4CIDRBlock = core.StringPtr("10.0.0.0/24")
	model.LifecycleReasons = []vpcv1.ClusterNetworkSubnetLifecycleReason{*clusterNetworkSubnetLifecycleReasonModel}
	model.LifecycleState = core.StringPtr("stable")
	model.Name = core.StringPtr("my-cluster-network-subnet")
	model.ResourceType = core.StringPtr("cluster_network_subnet")
	model.TotalIpv4AddressCount = core.Int64Ptr(int64(256))

	result, err := vpc.DataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetLifecycleReasonToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["code"] = "resource_suspended_by_provider"
		model["message"] = "The resource has been suspended. Contact IBM support with the CRN for next steps."
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#resource-suspension"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.ClusterNetworkSubnetLifecycleReason)
	model.Code = core.StringPtr("resource_suspended_by_provider")
	model.Message = core.StringPtr("The resource has been suspended. Contact IBM support with the CRN for next steps.")
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#resource-suspension")

	result, err := vpc.DataSourceIBMIsClusterNetworkSubnetsClusterNetworkSubnetLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
