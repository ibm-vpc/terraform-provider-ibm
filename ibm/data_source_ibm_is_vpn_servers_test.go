// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsVPNServersDataSourceBasic(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServersDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.client_auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.client_auto_delete_timeout"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.health_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.hostname"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_servers.is_vpn_servers", "vpn_servers.0.resource_type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServersDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_servers" "is_vpn_servers" {
		}
	`)
}
