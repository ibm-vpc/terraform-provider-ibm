// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsVPNServerRouteBasic(t *testing.T) {
	var conf vpcv1.VPNServerRoute
	vpnServerID := "r134-5cd03c9c-7122-46ad-97d5-b736c4345e05"
	destination := "172.16.0.0/16"
	name := fmt.Sprintf("tfname%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tfname%d", acctest.RandIntRange(10, 100))
	action := "translate"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMIsVPNServerRouteDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerRouteConfigBasic(vpnServerID, destination, action, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsVPNServerRouteExists("ibm_is_vpn_server_route.is_vpn_server_route", conf),
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "vpn_server", vpnServerID),
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "destination", destination),
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "action", action),
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerRouteConfigBasic(vpnServerID, destination, action, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "vpn_server", vpnServerID),
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "destination", destination),
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "action", action),
					resource.TestCheckResourceAttr("ibm_is_vpn_server_route.is_vpn_server_route", "name", nameUpdate),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServerRouteConfigBasic(vpnServerID string, destination string, action string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_vpn_server_route" "is_vpn_server_route" {
			vpn_server = "%s"
			destination = "%s"
			action = "%s"
			name = "%s"
		}
	`, vpnServerID, destination, action, name)
}

func testAccCheckIBMIsVPNServerRouteExists(n string, obj vpcv1.VPNServerRoute) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		sess, _ := testAccProvider.Meta().(ClientSession).VpcV1API()

		getVPNServerRouteOptions := &vpcv1.GetVPNServerRouteOptions{}

		parts, err := sepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getVPNServerRouteOptions.SetVPNServerID(parts[0])
		getVPNServerRouteOptions.SetID(parts[1])

		vpnServerRoute, _, err := sess.GetVPNServerRoute(getVPNServerRouteOptions)
		if err != nil {
			return err
		}

		obj = *vpnServerRoute
		return nil
	}
}

func testAccCheckIBMIsVPNServerRouteDestroy(s *terraform.State) error {
	sess, _ := testAccProvider.Meta().(ClientSession).VpcV1API()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_vpn_server_route" {
			continue
		}

		getVPNServerRouteOptions := &vpcv1.GetVPNServerRouteOptions{}

		parts, err := sepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getVPNServerRouteOptions.SetVPNServerID(parts[0])
		getVPNServerRouteOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := sess.GetVPNServerRoute(getVPNServerRouteOptions)

		if err == nil {
			return fmt.Errorf("VPNServerRoute still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for VPNServerRoute (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
