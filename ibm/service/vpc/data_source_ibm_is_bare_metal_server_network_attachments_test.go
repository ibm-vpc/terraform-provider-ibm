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

func TestAccIBMIsBareMetalServerNetworkAttachmentsDataSourceBasic(t *testing.T) {
	bareMetalServerNetworkAttachmentBareMetalServerID := fmt.Sprintf("tf_bare_metal_server_id_%d", acctest.RandIntRange(10, 100))
	bareMetalServerNetworkAttachmentInterfaceType := "hipersocket"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentsDataSourceConfigBasic(bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "bare_metal_server_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.#"),
					resource.TestCheckResourceAttr("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.interface_type", bareMetalServerNetworkAttachmentInterfaceType),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "total_count"),
				),
			},
		},
	})
}

func TestAccIBMIsBareMetalServerNetworkAttachmentsDataSourceAllArgs(t *testing.T) {
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
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentsDataSourceConfig(bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType, bareMetalServerNetworkAttachmentName, bareMetalServerNetworkAttachmentAllowToFloat, bareMetalServerNetworkAttachmentVlan),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "bare_metal_server_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.interface_type", bareMetalServerNetworkAttachmentInterfaceType),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.lifecycle_state"),
					resource.TestCheckResourceAttr("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.name", bareMetalServerNetworkAttachmentName),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.port_speed"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.type"),
					resource.TestCheckResourceAttr("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.allow_to_float", bareMetalServerNetworkAttachmentAllowToFloat),
					resource.TestCheckResourceAttr("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "network_attachments.0.vlan", bareMetalServerNetworkAttachmentVlan),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "next.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_bare_metal_server_network_attachments.is_bare_metal_server_network_attachments", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentsDataSourceConfigBasic(bareMetalServerNetworkAttachmentBareMetalServerID string, bareMetalServerNetworkAttachmentInterfaceType string) string {
	return fmt.Sprintf(`
		resource "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
			bare_metal_server_id = "%s"
			interface_type = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
				resource_type = "virtual_network_interface"
			}
		}

		data "ibm_is_bare_metal_server_network_attachments" "is_bare_metal_server_network_attachments_instance" {
			bare_metal_server_id = ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment.bare_metal_server_id
		}
	`, bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType)
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentsDataSourceConfig(bareMetalServerNetworkAttachmentBareMetalServerID string, bareMetalServerNetworkAttachmentInterfaceType string, bareMetalServerNetworkAttachmentName string, bareMetalServerNetworkAttachmentAllowToFloat string, bareMetalServerNetworkAttachmentVlan string) string {
	return fmt.Sprintf(`
		resource "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
			bare_metal_server_id = "%s"
			interface_type = "%s"
			name = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
				resource_type = "virtual_network_interface"
			}
			allowed_vlans = "FIXME"
			allow_to_float = %s
			vlan = %s
		}

		data "ibm_is_bare_metal_server_network_attachments" "is_bare_metal_server_network_attachments_instance" {
			bare_metal_server_id = ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment.bare_metal_server_id
		}
	`, bareMetalServerNetworkAttachmentBareMetalServerID, bareMetalServerNetworkAttachmentInterfaceType, bareMetalServerNetworkAttachmentName, bareMetalServerNetworkAttachmentAllowToFloat, bareMetalServerNetworkAttachmentVlan)
}
