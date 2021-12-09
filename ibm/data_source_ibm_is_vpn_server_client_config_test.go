// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsVPNServerClientConfigDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerClientConfigDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server_client_configuration.is_vpn_server_client_configuration", "vpn_server"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServerClientConfigDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_server_client_configuration" "is_vpn_server_client_configuration" {
			vpn_server = "r134-b04f5b2a-2452-41ba-a90b-14afd7a52cd6"
			file_path = "sunithatestingfile.txt"
		}
	`)
}
