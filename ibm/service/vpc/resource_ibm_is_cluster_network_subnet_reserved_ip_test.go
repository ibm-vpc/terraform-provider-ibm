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

func TestAccIBMIsClusterNetworkSubnetReservedIPBasic(t *testing.T) {
	var conf vpcv1.ClusterNetworkSubnetReservedIP
	clusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))
	clusterNetworkSubnetID := fmt.Sprintf("tf_cluster_network_subnet_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkSubnetReservedIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetReservedIPConfigBasic(clusterNetworkID, clusterNetworkSubnetID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkSubnetReservedIPExists("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_id", clusterNetworkID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_subnet_id", clusterNetworkSubnetID),
				),
			},
		},
	})
}

func TestAccIBMIsClusterNetworkSubnetReservedIPAllArgs(t *testing.T) {
	var conf vpcv1.ClusterNetworkSubnetReservedIP
	clusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))
	clusterNetworkSubnetID := fmt.Sprintf("tf_cluster_network_subnet_id_%d", acctest.RandIntRange(10, 100))
	address := fmt.Sprintf("tf_address_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	addressUpdate := fmt.Sprintf("tf_address_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkSubnetReservedIPDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetReservedIPConfig(clusterNetworkID, clusterNetworkSubnetID, address, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkSubnetReservedIPExists("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_id", clusterNetworkID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_subnet_id", clusterNetworkSubnetID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "address", address),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetReservedIPConfig(clusterNetworkID, clusterNetworkSubnetID, addressUpdate, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_id", clusterNetworkID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "cluster_network_subnet_id", clusterNetworkSubnetID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "address", addressUpdate),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_cluster_network_subnet_reserved_ip.is_cluster_network_subnet_reserved_ip",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkSubnetReservedIPConfigBasic(clusterNetworkID string, clusterNetworkSubnetID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_cluster_network_subnet_reserved_ip" "is_cluster_network_subnet_reserved_ip_instance" {
			cluster_network_id = "%s"
			cluster_network_subnet_id = "%s"
		}
	`, clusterNetworkID, clusterNetworkSubnetID)
}

func testAccCheckIBMIsClusterNetworkSubnetReservedIPConfig(clusterNetworkID string, clusterNetworkSubnetID string, address string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_cluster_network_subnet_reserved_ip" "is_cluster_network_subnet_reserved_ip_instance" {
			cluster_network_id = "%s"
			cluster_network_subnet_id = "%s"
			address = "%s"
			name = "%s"
		}
	`, clusterNetworkID, clusterNetworkSubnetID, address, name)
}

func testAccCheckIBMIsClusterNetworkSubnetReservedIPExists(n string, obj vpcv1.ClusterNetworkSubnetReservedIP) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getClusterNetworkSubnetReservedIPOptions := &vpcv1.GetClusterNetworkSubnetReservedIPOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkID(parts[0])
		getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkSubnetID(parts[1])
		getClusterNetworkSubnetReservedIPOptions.SetID(parts[2])

		clusterNetworkSubnetReservedIP, _, err := vpcClient.GetClusterNetworkSubnetReservedIP(getClusterNetworkSubnetReservedIPOptions)
		if err != nil {
			return err
		}

		obj = *clusterNetworkSubnetReservedIP
		return nil
	}
}

func testAccCheckIBMIsClusterNetworkSubnetReservedIPDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_cluster_network_subnet_reserved_ip" {
			continue
		}

		getClusterNetworkSubnetReservedIPOptions := &vpcv1.GetClusterNetworkSubnetReservedIPOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkID(parts[0])
		getClusterNetworkSubnetReservedIPOptions.SetClusterNetworkSubnetID(parts[1])
		getClusterNetworkSubnetReservedIPOptions.SetID(parts[2])

		// Try to find the key
		_, response, err := vpcClient.GetClusterNetworkSubnetReservedIP(getClusterNetworkSubnetReservedIPOptions)

		if err == nil {
			return fmt.Errorf("ClusterNetworkSubnetReservedIP still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for ClusterNetworkSubnetReservedIP (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(t *testing.T) {
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

	result, err := vpc.ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPLifecycleReasonToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["id"] = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["name"] = "my-cluster-network-interface"
		model["resource_type"] = "cluster_network_interface"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReservedIPTarget)
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.ID = core.StringPtr("0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.Name = core.StringPtr("my-cluster-network-interface")
	model.ResourceType = core.StringPtr("cluster_network_interface")

	result, err := vpc.ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.Deleted)
	model.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	result, err := vpc.ResourceIBMIsClusterNetworkSubnetReservedIPDeletedToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		deletedModel := make(map[string]interface{})
		deletedModel["more_info"] = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"

		model := make(map[string]interface{})
		model["deleted"] = []map[string]interface{}{deletedModel}
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["id"] = "0717-ffc092f7-5d02-4b93-ab69-26860529b9fb"
		model["name"] = "my-cluster-network-interface"
		model["resource_type"] = "cluster_network_interface"

		assert.Equal(t, result, model)
	}

	deletedModel := new(vpcv1.Deleted)
	deletedModel.MoreInfo = core.StringPtr("https://cloud.ibm.com/apidocs/vpc#deleted-resources")

	model := new(vpcv1.ClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContext)
	model.Deleted = deletedModel
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/cluster_networks/0717-da0df18c-7598-4633-a648-fdaac28a5573/interfaces/0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.ID = core.StringPtr("0717-ffc092f7-5d02-4b93-ab69-26860529b9fb")
	model.Name = core.StringPtr("my-cluster-network-interface")
	model.ResourceType = core.StringPtr("cluster_network_interface")

	result, err := vpc.ResourceIBMIsClusterNetworkSubnetReservedIPClusterNetworkSubnetReservedIPTargetClusterNetworkInterfaceReferenceClusterNetworkSubnetReservedIPTargetContextToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
