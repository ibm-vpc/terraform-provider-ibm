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

func TestAccIBMIsBareMetalServerNetworkAttachmentDataSourceBasic(t *testing.T) {
	bareMetalServerNetworkAttachmentBareMetalServerID := fmt.Sprintf("tf_bare_metal_server_id_%d", acctest.RandIntRange(10, 100))
	bareMetalServerNetworkAttachmentInterfaceType := "hipersocket"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentDataSourceConfigBasic(bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_network_attachment_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "interface_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "port_speed"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "primary_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "subnet.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "virtual_network_interface.#"),
				),
			},
		},
	})
}

func TestAccIBMIsBareMetalServerNetworkAttachmentDataSourceAllArgs(t *testing.T) {
	bareMetalServerNetworkAttachmentBareMetalServerID := fmt.Sprintf("tf_bare_metal_server_id_%d", acctest.RandIntRange(10, 100))
	bareMetalServerNetworkAttachmentInterfaceType := "hipersocket"
	bareMetalServerNetworkAttachmentName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	bareMetalServerNetworkAttachmentAllowToFloat := "true"
	bareMetalServerNetworkAttachmentVlan := fmt.Sprintf("%d", acctest.RandIntRange(1, 4094))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentDataSourceConfig(bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType, bareMetalServerNetworkAttachmentName, bareMetalServerNetworkAttachmentAllowToFloat, bareMetalServerNetworkAttachmentVlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_network_attachment_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "interface_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "port_speed"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "primary_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "subnet.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "virtual_network_interface.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "allowed_vlans.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "allow_to_float"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "vlan"),
				),
			},
		},
	})
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentDataSourceConfigBasic(bareMetalServerNetworkAttachmentBareMetalServerID string, bareMetalServerNetworkAttachmentInterfaceType string) string {
	return fmt.Sprintf(`
		resource "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
			bare_metal_server_id = "%s"
			interface_type = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
				resource_type = "virtual_network_interface"
			}
		}

		data "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
			bare_metal_server_id = ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment.bare_metal_server_id
			id = "id"
		}
	`, bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType)
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentDataSourceConfig(bareMetalServerNetworkAttachmentBareMetalServerID string, bareMetalServerNetworkAttachmentInterfaceType string, bareMetalServerNetworkAttachmentName string, bareMetalServerNetworkAttachmentAllowToFloat string, bareMetalServerNetworkAttachmentVlan string) string {
	return fmt.Sprintf(`
		resource "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
			bare_metal_server_id = "%s"
			interface_type = "%s"
			name = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
				resource_type = "virtual_network_interface"
			}
			allowed_vlans = "FIXME"
			allow_to_float = %s
			vlan = %s
		}

		data "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
			bare_metal_server_id = ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment.bare_metal_server_id
			id = "id"
		}
	`, bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType, bareMetalServerNetworkAttachmentName, bareMetalServerNetworkAttachmentAllowToFloat, bareMetalServerNetworkAttachmentVlan)
}
