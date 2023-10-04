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

func TestAccIBMIsInstanceNetworkAttachmentDataSourceBasic(t *testing.T) {
	instanceByNetworkAttachmentInstanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceNetworkAttachmentDataSourceConfigBasic(instanceByNetworkAttachmentInstanceID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "instance_by_network_attachment_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "port_speed"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "primary_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "subnet.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "virtual_network_interface.#"),
				),
			},
		},
	})
}

func TestAccIBMIsInstanceNetworkAttachmentDataSourceAllArgs(t *testing.T) {
	instanceByNetworkAttachmentInstanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))
	instanceByNetworkAttachmentName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceNetworkAttachmentDataSourceConfig(instanceByNetworkAttachmentInstanceID, instanceByNetworkAttachmentName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "instance_by_network_attachment_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "port_speed"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "primary_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "subnet.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_instance_network_attachment.is_instance_network_attachment", "virtual_network_interface.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsInstanceNetworkAttachmentDataSourceConfigBasic(instanceByNetworkAttachmentInstanceID string) string {
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

		data "ibm_is_instance_network_attachment" "is_instance_network_attachment_instance" {
			instance_id = ibm_is_instance_network_attachment.is_instance_network_attachment.instance_id
			id = "id"
		}
	`, instanceByNetworkAttachmentInstanceID)
}

func testAccCheckIBMIsInstanceNetworkAttachmentDataSourceConfig(instanceByNetworkAttachmentInstanceID string, instanceByNetworkAttachmentName string) string {
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

		data "ibm_is_instance_network_attachment" "is_instance_network_attachment_instance" {
			instance_id = ibm_is_instance_network_attachment.is_instance_network_attachment.instance_id
			id = "id"
		}
	`, instanceByNetworkAttachmentInstanceID, instanceByNetworkAttachmentName)
}
