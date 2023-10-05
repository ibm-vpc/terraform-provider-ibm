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

func TestAccIBMIsBareMetalServerNetworkAttachmentBasic(t *testing.T) {
	var conf vpcv1.BareMetalServerNetworkAttachment
	bareMetalServerID := fmt.Sprintf("tf_bare_metal_server_id_%d", acctest.RandIntRange(10, 100))
	interfaceType := "hipersocket"
	interfaceTypeUpdate := "vlan"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsBareMetalServerNetworkAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentConfigBasic(bareMetalServerID, interfaceType),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsBareMetalServerNetworkAttachmentExists("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", conf),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_id", bareMetalServerID),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "interface_type", interfaceType),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentConfigBasic(bareMetalServerID, interfaceTypeUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_id", bareMetalServerID),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "interface_type", interfaceTypeUpdate),
				),
			},
		},
	})
}

func TestAccIBMIsBareMetalServerNetworkAttachmentAllArgs(t *testing.T) {
	var conf vpcv1.BareMetalServerNetworkAttachment
	bareMetalServerID := fmt.Sprintf("tf_bare_metal_server_id_%d", acctest.RandIntRange(10, 100))
	interfaceType := "hipersocket"
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	allowToFloat := "true"
	vlan := fmt.Sprintf("%d", acctest.RandIntRange(1, 4094))
	interfaceTypeUpdate := "vlan"
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	allowToFloatUpdate := "false"
	vlanUpdate := fmt.Sprintf("%d", acctest.RandIntRange(1, 4094))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsBareMetalServerNetworkAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentConfig(bareMetalServerID, interfaceType, name, allowToFloat, vlan),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsBareMetalServerNetworkAttachmentExists("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", conf),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_id", bareMetalServerID),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "interface_type", interfaceType),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "name", name),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "allow_to_float", allowToFloat),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "vlan", vlan),
				),
			},
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentConfig(bareMetalServerID, interfaceTypeUpdate, nameUpdate, allowToFloatUpdate, vlanUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "bare_metal_server_id", bareMetalServerID),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "interface_type", interfaceTypeUpdate),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "allow_to_float", allowToFloatUpdate),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "vlan", vlanUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentConfigBasic(bareMetalServerID string, interfaceType string) string {
	return fmt.Sprintf(`
		resource "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment_instance" {
			bare_metal_server_id = "%s"
			interface_type = "%s"
			virtual_network_interface {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::virtual-network-interface:0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				href = "https://us-south.iaas.cloud.ibm.com/v1/virtual_network_interfaces/0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				id = "0767-fa41aecb-4f21-423d-8082-630bfba1e1d9"
				name = "my-virtual-network-interface"
			}
		}
	`, bareMetalServerID, interfaceType)
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentConfig(bareMetalServerID string, interfaceType string, name string, allowToFloat string, vlan string) string {
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
			}
			allowed_vlans = "FIXME"
			allow_to_float = %s
			vlan = %s
		}
	`, bareMetalServerID, interfaceType, name, allowToFloat, vlan)
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentExists(n string, obj vpcv1.BareMetalServerNetworkAttachment) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := vpcClient(acc.TestAccProvider.Meta())
		if err != nil {
			return err
		}

		getBareMetalServerNetworkAttachmentOptions := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(parts[0])
		getBareMetalServerNetworkAttachmentOptions.SetID(parts[1])

		bareMetalServerNetworkAttachmentIntf, _, err := vpcClient.GetBareMetalServerNetworkAttachment(getBareMetalServerNetworkAttachmentOptions)
		if err != nil {
			return err
		}

		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment)
		obj = *bareMetalServerNetworkAttachment
		return nil
	}
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentDestroy(s *terraform.State) error {
	vpcClient, err := vpcClient(acc.TestAccProvider.Meta())
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_bare_metal_server_network_attachment" {
			continue
		}

		getBareMetalServerNetworkAttachmentOptions := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(parts[0])
		getBareMetalServerNetworkAttachmentOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetBareMetalServerNetworkAttachment(getBareMetalServerNetworkAttachmentOptions)

		if err == nil {
			return fmt.Errorf("is_bare_metal_server_network_attachment still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for is_bare_metal_server_network_attachment (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
