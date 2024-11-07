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
