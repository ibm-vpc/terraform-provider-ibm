// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"strings"
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
	vpcname := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf-server-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfip-subnet-%d", acctest.RandIntRange(10, 100))
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	sshname := fmt.Sprintf("tf-sshname-%d", acctest.RandIntRange(10, 100))
	vniname := fmt.Sprintf("tf-vni-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsBareMetalServerNetworkAttachmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBareMetalServerNetworkAttachmentConfigBasic(vpcname, subnetname, sshname, publicKey, vniname, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsBareMetalServerNetworkAttachmentExists("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", conf),
					resource.TestCheckResourceAttr("ibm_is_bare_metal_server_network_attachment.is_bare_metal_server_network_attachment", "resource_type", "bare_metal_server_network_attachment"),
				),
			},
		},
	})
}

func testAccCheckIBMIsBareMetalServerNetworkAttachmentConfigBasic(vpcname, subnetname, sshname, publicKey, vniname, name string) string {
	return testAccCheckIBMISBareMetalServerVNIConfig(vpcname, subnetname, sshname, publicKey, vniname, name) + fmt.Sprintf(`
	resource "ibm_is_virtual_network_interface" "testacc_vni2"{
		name = "test-vni-na"
		subnet = ibm_is_subnet.testacc_subnet.id
		enable_infrastructure_nat = true
		allow_ip_spoofing = true
	}	
	resource "ibm_is_bare_metal_server_network_attachment" "is_bare_metal_server_network_attachment" {
			bare_metal_server = ibm_is_bare_metal_server.testacc_bms.id
			allowed_vlans = [200, 202, 203]
			virtual_network_interface { 
				id = ibm_is_virtual_network_interface.testacc_vni2.id
			}
	}
	`)
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
