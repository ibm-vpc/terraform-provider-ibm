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

func TestAccIBMIsInstanceNetworkAttachmentsDataSourceBasic(t *testing.T) {
	instanceByNetworkAttachmentInstanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceNetworkAttachmentsDataSourceConfigBasic(instanceByNetworkAttachmentInstanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.#"),
				),
			},
		},
	})
}

func TestAccIBMIsInstanceNetworkAttachmentsDataSourceAllArgs(t *testing.T) {
	instanceByNetworkAttachmentInstanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))
	instanceByNetworkAttachmentName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceNetworkAttachmentsDataSourceConfig(instanceByNetworkAttachmentInstanceID, instanceByNetworkAttachmentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.lifecycle_state"),
					resource.TestCheckResourceAttr("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.name", instanceByNetworkAttachmentName),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.port_speed"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachments.is_instance_network_attachments", "network_attachments.0.type"),
				),
			},
		},
	})
}

func testAccCheckIBMIsInstanceNetworkAttachmentsDataSourceConfigBasic(instanceByNetworkAttachmentInstanceID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_instance_network_attachment" "is_instance_network_attachment_instance" {
			instance_id = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
				resource_type = "virtual_network_interface"
			}
		}

		data "ibm_is_instance_network_attachments" "is_instance_network_attachments_instance" {
			instance_id = ibm_is_instance_network_attachment.is_instance_network_attachment.instance_id
		}
	`, instanceByNetworkAttachmentInstanceID)
}

func testAccCheckIBMIsInstanceNetworkAttachmentsDataSourceConfig(instanceByNetworkAttachmentInstanceID string, instanceByNetworkAttachmentName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_instance_network_attachment" "is_instance_network_attachment_instance" {
			instance_id = "%s"
			name = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
				resource_type = "virtual_network_interface"
			}
		}

		data "ibm_is_instance_network_attachments" "is_instance_network_attachments_instance" {
			instance_id = ibm_is_instance_network_attachment.is_instance_network_attachment.instance_id
		}
	`, instanceByNetworkAttachmentInstanceID, instanceByNetworkAttachmentName)
}
