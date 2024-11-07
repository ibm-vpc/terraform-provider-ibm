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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsClusterNetworkBasic(t *testing.T) {
	var conf vpcv1.ClusterNetwork

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkExists("ibm_is_cluster_network.is_cluster_network_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network.is_cluster_network_instance", "cluster_network_id", "clusterNetworkID"),
				),
			},
		},
	})
}

func TestAccIBMIsClusterNetworkAllArgs(t *testing.T) {
	var conf vpcv1.ClusterNetwork
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsClusterNetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsClusterNetworkExists("ibm_is_cluster_network.is_cluster_network_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_cluster_network.is_cluster_network_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkConfig(nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_cluster_network.is_cluster_network_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_cluster_network.is_cluster_network",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkConfigBasic() string {
	return fmt.Sprintf(`
		resource "ibm_is_cluster_network" "is_cluster_network_instance" {
			profile = "h100"
			vpc {
				id = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
			}
			zone  = "us-south-1"
		}
	`)
}

func testAccCheckIBMIsClusterNetworkConfig(name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_cluster_network" "is_cluster_network_instance" {
			name = "%s"
			profile {
				href = "https://us-south.iaas.cloud.ibm.com/v1/cluster_network/profiles/h100"
				name = "h100"
			}
			resource_group {
				id = "fee82deba12e4c0fb69c3b09d1f12345"
			}
			subnet_prefixes {
				allocation_policy = "auto"
				cidr = "10.0.0.0/24"
			}
			vpc {
				crn = "crn:v1:bluemix:public:is:us-south:a/aa2432b1fa4d4ace891e9b80fc104e34::vpc:r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
				href = "https://us-south.iaas.cloud.ibm.com/v1/vpcs/r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
				id = "r006-4727d842-f94f-4a2d-824a-9bc9b02c523b"
			}
			zone {
				href = "https://us-south.iaas.cloud.ibm.com/v1/regions/us-south/zones/us-south-1"
				name = "us-south-1"
			}
		}
	`, name)
}

func testAccCheckIBMIsClusterNetworkExists(n string, obj vpcv1.ClusterNetwork) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{}

		getClusterNetworkOptions.SetID(rs.Primary.ID)

		clusterNetwork, _, err := vpcClient.GetClusterNetwork(getClusterNetworkOptions)
		if err != nil {
			return err
		}

		obj = *clusterNetwork
		return nil
	}
}

func testAccCheckIBMIsClusterNetworkDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_cluster_network" {
			continue
		}

		getClusterNetworkOptions := &vpcv1.GetClusterNetworkOptions{}

		getClusterNetworkOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := vpcClient.GetClusterNetwork(getClusterNetworkOptions)

		if err == nil {
			return fmt.Errorf("ClusterNetwork still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for ClusterNetwork (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
