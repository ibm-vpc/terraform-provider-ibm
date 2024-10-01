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

func TestAccIBMIsClusterNetworkSubnetDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "cluster_network_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "is_cluster_network_subnet_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "available_ipv4_address_count"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "ipv4_cidr_block"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "lifecycle_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "total_ipv4_address_count"),
				),
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkSubnetDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_cluster_network_subnet" "is_cluster_network_subnet_instance" {
			cluster_network_id = "cluster_network_id"
			id = "id"
		}
	`)
}

func TestDataSourceIBMIsClusterNetworkSubnetClusterNetworkSubnetLifecycleReasonToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsClusterNetworkSubnetClusterNetworkSubnetLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
