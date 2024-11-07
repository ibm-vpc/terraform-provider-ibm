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
	"github.com/IBM/vpc-go-sdk/vpcv1"
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
