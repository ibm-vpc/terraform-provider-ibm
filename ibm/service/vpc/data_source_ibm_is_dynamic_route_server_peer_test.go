// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsDynamicRouteServerPeerDataSourceBasic(t *testing.T) {
	dynamicRouteServerPeerDynamicRouteServerID := fmt.Sprintf("tf_dynamic_route_server_id_%d", acctest.RandIntRange(10, 100))
	dynamicRouteServerPeerAsn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerPeerDataSourceConfigBasic(dynamicRouteServerPeerDynamicRouteServerID, dynamicRouteServerPeerAsn),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "asn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "authentication_enabled"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "bfd.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_peer_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "sessions.#"),
				),
			},
		},
	})
}

func TestAccIBMIsDynamicRouteServerPeerDataSourceAllArgs(t *testing.T) {
	dynamicRouteServerPeerDynamicRouteServerID := fmt.Sprintf("tf_dynamic_route_server_id_%d", acctest.RandIntRange(10, 100))
	dynamicRouteServerPeerAsn := fmt.Sprintf("%d", acctest.RandIntRange(10, 100))
	dynamicRouteServerPeerName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerPeerDataSourceConfig(dynamicRouteServerPeerDynamicRouteServerID, dynamicRouteServerPeerAsn, dynamicRouteServerPeerName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "asn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "authentication_enabled"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "bfd.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "dynamic_route_server_peer_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "sessions.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "sessions.0.established_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer", "sessions.0.state"),
				),
			},
		},
	})
}

func testAccCheckIBMIsDynamicRouteServerPeerDataSourceConfigBasic(dynamicRouteServerPeerDynamicRouteServerID string, dynamicRouteServerPeerAsn string) string {
	return fmt.Sprintf(`
		resource "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer_instance" {
			dynamic_route_server_id = "%s"
			asn = %s
			ip {
				address = "192.168.3.4"
			}
		}

		data "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer_instance" {
			dynamic_route_server_id = ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer.dynamic_route_server_id
			id = "id"
		}
	`, dynamicRouteServerPeerDynamicRouteServerID, dynamicRouteServerPeerAsn)
}

func testAccCheckIBMIsDynamicRouteServerPeerDataSourceConfig(dynamicRouteServerPeerDynamicRouteServerID string, dynamicRouteServerPeerAsn string, dynamicRouteServerPeerName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer_instance" {
			dynamic_route_server_id = "%s"
			asn = %s
			bfd {
				mode = "asynchronous"
				role = "active"
				sessions {
					source_ip {
						address = "192.168.3.4"
						deleted {
							more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
						}
						href = "https://us-south.iaas.cloud.ibm.com/v1/subnets/7ec86020-1c6e-4889-b3f0-a15f2e50f87e/reserved_ips/6d353a0f-aeb1-4ae1-832e-1110d10981bb"
						id = "6d353a0f-aeb1-4ae1-832e-1110d10981bb"
						name = "my-reserved-ip"
						resource_type = "subnet_reserved_ip"
					}
					state = "admin_down"
				}
			}
			ip {
				address = "192.168.3.4"
			}
			name = "%s"
		}

		data "ibm_is_dynamic_route_server_peer" "is_dynamic_route_server_peer_instance" {
			dynamic_route_server_id = ibm_is_dynamic_route_server_peer.is_dynamic_route_server_peer.dynamic_route_server_id
			id = "id"
		}
	`, dynamicRouteServerPeerDynamicRouteServerID, dynamicRouteServerPeerAsn, dynamicRouteServerPeerName)
}
