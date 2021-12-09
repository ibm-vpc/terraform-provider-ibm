// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsVPNServerRoutesDataSourceBasic(t *testing.T) {
	vpnServerRouteVPNServerID := fmt.Sprintf("tf_vpn_server_id_%d", acctest.RandIntRange(10, 100))
	vpnServerRouteDestination := fmt.Sprintf("tf_destination_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerRoutesDataSourceConfigBasic(vpnServerRouteVPNServerID, vpnServerRouteDestination),
				Check: resource.ComposeTestCheckFunc(
					// resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_routes.is_vpn_server_routes", "routes.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_routes.is_vpn_server_routes", "vpn_server"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_routes.is_vpn_server_routes", "routes.#"),
					resource.TestCheckResourceAttr("data.ibm_is_vpn_server_routes.is_vpn_server_routes", "routes.0.destination", "172.16.0.0/16"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServerRoutesDataSourceConfigBasic(vpnServerRouteVPNServerID string, vpnServerRouteDestination string) string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_server_routes" "is_vpn_server_routes" {
			vpn_server = "r134-5cd03c9c-7122-46ad-97d5-b736c4345e05"
		}
	`)
}
