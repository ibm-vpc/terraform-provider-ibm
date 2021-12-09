// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsVPNServerClientsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerClientsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_clients.is_vpn_server_clients", "vpn_server"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_clients.is_vpn_server_clients", "clients.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServerClientsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_server_clients" "is_vpn_server_clients" {
			vpn_server = "r134-5cd03c9c-7122-46ad-97d5-b736c4345e05"
		}
	`)
}
