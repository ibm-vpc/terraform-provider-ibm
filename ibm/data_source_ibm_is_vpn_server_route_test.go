// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsVPNServerRouteDataSourceBasic(t *testing.T) {
	vpnServerRouteVPNServerID := fmt.Sprintf("tf_vpn_server_%d", acctest.RandIntRange(10, 100))
	vpnServerRouteDestination := fmt.Sprintf("tf_destination_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerRouteDataSourceConfigBasic(vpnServerRouteVPNServerID, vpnServerRouteDestination),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "identifier"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "vpn_server"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "action"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "destination"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_route.is_vpn_server_route", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServerRouteDataSourceConfigBasic(vpnServerRouteVPNServerID string, vpnServerRouteDestination string) string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_server_route" "is_vpn_server_route" {
			vpn_server = "r134-5cd03c9c-7122-46ad-97d5-b736c4345e05"
			identifier = "r134-1f0e1d21-c9e9-410b-8d6c-7ac4d8b2c835"
		}
	`)
}
