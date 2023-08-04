// Copyright IBM Corp. 2023 All Rights Reserved.
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

func TestAccIBMIsDynamicRouteServerPeerBasic(t *testing.T) {
	var conf vpcv1.DynamicRouteServerPeer
	dynamicRouteServerID := fmt.Sprintf("tf_dynamic_route_server_id_%d", acctest.RandIntRange(10, 100))
	asn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	asnUpdate := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsDynamicRouteServerPeerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerPeerConfigBasic(dynamicRouteServerID, asn),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsDynamicRouteServerPeerExists("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", conf),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_id", dynamicRouteServerID),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "asn", asn),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerPeerConfigBasic(dynamicRouteServerID, asnUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_id", dynamicRouteServerID),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "asn", asnUpdate),
				),
			},
		},
	})
}

func TestAccIBMIsDynamicRouteServerPeerAllArgs(t *testing.T) {
	var conf vpcv1.DynamicRouteServerPeer
	dynamicRouteServerID := fmt.Sprintf("tf_dynamic_route_server_id_%d", acctest.RandIntRange(10, 100))
	asn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	asnUpdate := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsDynamicRouteServerPeerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerPeerConfig(dynamicRouteServerID, asn, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsDynamicRouteServerPeerExists("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", conf),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_id", dynamicRouteServerID),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "asn", asn),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerPeerConfig(dynamicRouteServerID, asnUpdate, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_id", dynamicRouteServerID),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "asn", asnUpdate),
					resource.TestCheckResourceAttr("ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsDynamicRouteServerPeerConfigBasic(dynamicRouteServerID string, asn string) string {
	return fmt.Sprintf(`
		resource "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer_instance" {
			dynamic_route_server_id = "%s"
			asn = %s
			ip {
				address = "192.168.3.4"
			}
		}
	`, dynamicRouteServerID, asn)
}

func testAccCheckIBMIsDynamicRouteServerPeerConfig(dynamicRouteServerID string, asn string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer_instance" {
			dynamic_route_server_id = "%s"
			asn = %s
			bfd {
				role = "active"
			}
			ip {
				address = "192.168.3.4"
			}
			name = "%s"
		}
	`, dynamicRouteServerID, asn, name)
}

func testAccCheckIBMIsDynamicRouteServerPeerExists(n string, obj vpcv1.DynamicRouteServerPeer) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getDynamicRouteServerPeerOptions := &vpcv1.GetDynamicRouteServerPeerOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getDynamicRouteServerPeerOptions.SetDynamicRouteServerID(parts[0])
		getDynamicRouteServerPeerOptions.SetID(parts[1])

		dynamicRouteServerPeer, _, err := vpcClient.GetDynamicRouteServerPeer(getDynamicRouteServerPeerOptions)
		if err != nil {
			return err
		}

		obj = *dynamicRouteServerPeer
		return nil
	}
}

func testAccCheckIBMIsDynamicRouteServerPeerDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_dynamic_route_server_peer" {
			continue
		}

		getDynamicRouteServerPeerOptions := &vpcv1.GetDynamicRouteServerPeerOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getDynamicRouteServerPeerOptions.SetDynamicRouteServerID(parts[0])
		getDynamicRouteServerPeerOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetDynamicRouteServerPeer(getDynamicRouteServerPeerOptions)

		if err == nil {
			return fmt.Errorf("DynamicRouteServerPeer still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for DynamicRouteServerPeer (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
