// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsVPNServerDataSourceBasic(t *testing.T) {
	identifier := "r134-5cd03c9c-7122-46ad-97d5-b736c4345e05"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNServerDataSourceConfigBasic(identifier),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "identifier"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "certificate.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "client_authentication.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "client_auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "client_auto_delete_timeout"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "client_dns_server_ips.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "client_idle_timeout"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "client_ip_pool"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "enable_split_tunneling"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "health_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "hostname"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "port"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "private_ips.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "protocol"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "resource_group.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "security_groups.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_server.is_vpn_server", "subnets.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNServerDataSourceConfigBasic(identifier string) string {
	return fmt.Sprintf(`
		data "ibm_is_vpn_server" "is_vpn_server" {
			identifier = "%s"
		}
	`, identifier)
}
