// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsInstanceNetworkAttachmentBasic(t *testing.T) {
	var conf vpcv1.InstanceNetworkAttachment
	instanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsInstanceNetworkAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceNetworkAttachmentConfigBasic(instanceID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsInstanceNetworkAttachmentExists("ibm_is_instance_network_attachment.is_instance_network_attachment", conf),
					resource.TestCheckResourceAttr("ibm_is_instance_network_attachment.is_instance_network_attachment", "instance_id", instanceID),
				),
			},
		},
	})
}

func TestAccIBMIsInstanceNetworkAttachmentAllArgs(t *testing.T) {
	var conf vpcv1.InstanceNetworkAttachment
	instanceID := fmt.Sprintf("tf_instance_id_%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsInstanceNetworkAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceNetworkAttachmentConfig(instanceID, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsInstanceNetworkAttachmentExists("ibm_is_instance_network_attachment.is_instance_network_attachment", conf),
					resource.TestCheckResourceAttr("ibm_is_instance_network_attachment.is_instance_network_attachment", "instance_id", instanceID),
					resource.TestCheckResourceAttr("ibm_is_instance_network_attachment.is_instance_network_attachment", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsInstanceNetworkAttachmentConfig(instanceID, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_instance_network_attachment.is_instance_network_attachment", "instance_id", instanceID),
					resource.TestCheckResourceAttr("ibm_is_instance_network_attachment.is_instance_network_attachment", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_instance_network_attachment.is_instance_network_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsInstanceNetworkAttachmentConfigBasic(instanceID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_instance_network_attachment" "is_instance_network_attachment_instance" {
			instance = "%s"
			virtual_network_interface {
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
			}
		}
	`, instanceID)
}

func testAccCheckIBMIsInstanceNetworkAttachmentConfig(instanceID string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_instance_network_attachment" "is_instance_network_attachment_instance" {
			instance_id = "%s"
			name = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
			}
		}
	`, instanceID, name)
}

func testAccCheckIBMIsInstanceNetworkAttachmentExists(n string, obj vpcv1.InstanceNetworkAttachment) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := vpcClient(acc.TestAccProvider.Meta())
		if err != nil {
			return err
		}

		getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
		getInstanceNetworkAttachmentOptions.SetID(parts[1])

		instanceByNetworkAttachment, _, err := vpcClient.GetInstanceNetworkAttachment(getInstanceNetworkAttachmentOptions)
		if err != nil {
			return err
		}

		obj = *instanceByNetworkAttachment
		return nil
	}
}

func testAccCheckIBMIsInstanceNetworkAttachmentDestroy(s *terraform.State) error {
	vpcClient, err := vpcClient(acc.TestAccProvider.Meta())
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_instance_network_attachment" {
			continue
		}

		getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
		getInstanceNetworkAttachmentOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetInstanceNetworkAttachment(getInstanceNetworkAttachmentOptions)

		if err == nil {
			return fmt.Errorf("InstanceByNetworkAttachment still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("[ERROR] Error checking for InstanceByNetworkAttachment (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
