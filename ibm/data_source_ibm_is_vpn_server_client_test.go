// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsVPNServerClientDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerClientDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "identifier"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "vpn_server"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "client_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "remote_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "remote_port"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client.is_vpn_server_client", "status"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServerClientDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_server_client" "is_vpn_server_client" {
			vpn_server = "vpn_server_id"
			identifier = "id"
		}
	`)
}
