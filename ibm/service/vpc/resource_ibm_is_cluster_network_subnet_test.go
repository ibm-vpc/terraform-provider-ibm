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

func TestAccIBMIsClusterNetworkSubnetBasic(t *testing.T) {
	var conf vpcv1.ClusterNetworkSubnet
	clusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkSubnetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetConfigBasic(clusterNetworkID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkSubnetExists("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "cluster_network_id", clusterNetworkID),
				),
			},
		},
	})
}

func TestAccIBMIsClusterNetworkSubnetAllArgs(t *testing.T) {
	var conf vpcv1.ClusterNetworkSubnet
	clusterNetworkID := fmt.Sprintf("tf_cluster_network_id_%d", acctest.RandIntRange(10, 100))
	ipVersion := "ipv4"
	ipv4CIDRBlock := fmt.Sprintf("tf_ipv4_cidr_block_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	totalIpv4AddressCount := fmt.Sprintf("%d", acctest.RandIntRange(8, 16777216))
	ipVersionUpdate := "ipv4"
	ipv4CIDRBlockUpdate := fmt.Sprintf("tf_ipv4_cidr_block_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	totalIpv4AddressCountUpdate := fmt.Sprintf("%d", acctest.RandIntRange(8, 16777216))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkSubnetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetConfig(clusterNetworkID, ipVersion, ipv4CIDRBlock, name, totalIpv4AddressCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkSubnetExists("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "cluster_network_id", clusterNetworkID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "ip_version", ipVersion),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "ipv4_cidr_block", ipv4CIDRBlock),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "total_ipv4_address_count", totalIpv4AddressCount),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkSubnetConfig(clusterNetworkID, ipVersionUpdate, ipv4CIDRBlockUpdate, nameUpdate, totalIpv4AddressCountUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "cluster_network_id", clusterNetworkID),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "ip_version", ipVersionUpdate),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "ipv4_cidr_block", ipv4CIDRBlockUpdate),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_cluster_network_subnet.is_cluster_network_subnet_instance", "total_ipv4_address_count", totalIpv4AddressCountUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_cluster_network_subnet.is_cluster_network_subnet",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkSubnetConfigBasic(clusterNetworkID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_cluster_network_subnet" "is_cluster_network_subnet_instance" {
			cluster_network_id = "%s"
		}
	`, clusterNetworkID)
}

func testAccCheckIBMIsClusterNetworkSubnetConfig(clusterNetworkID string, ipVersion string, ipv4CIDRBlock string, name string, totalIpv4AddressCount string) string {
	return fmt.Sprintf(`

		resource "ibm_is_cluster_network_subnet" "is_cluster_network_subnet_instance" {
			cluster_network_id = "%s"
			ip_version = "%s"
			ipv4_cidr_block = "%s"
			name = "%s"
			total_ipv4_address_count = %s
		}
	`, clusterNetworkID, ipVersion, ipv4CIDRBlock, name, totalIpv4AddressCount)
}

func testAccCheckIBMIsClusterNetworkSubnetExists(n string, obj vpcv1.ClusterNetworkSubnet) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getClusterNetworkSubnetOptions := &vpcv1.GetClusterNetworkSubnetOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getClusterNetworkSubnetOptions.SetClusterNetworkID(parts[0])
		getClusterNetworkSubnetOptions.SetID(parts[1])

		clusterNetworkSubnet, _, err := vpcClient.GetClusterNetworkSubnet(getClusterNetworkSubnetOptions)
		if err != nil {
			return err
		}

		obj = *clusterNetworkSubnet
		return nil
	}
}

func testAccCheckIBMIsClusterNetworkSubnetDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_cluster_network_subnet" {
			continue
		}

		getClusterNetworkSubnetOptions := &vpcv1.GetClusterNetworkSubnetOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getClusterNetworkSubnetOptions.SetClusterNetworkID(parts[0])
		getClusterNetworkSubnetOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetClusterNetworkSubnet(getClusterNetworkSubnetOptions)

		if err == nil {
			return fmt.Errorf("ClusterNetworkSubnet still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for ClusterNetworkSubnet (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
