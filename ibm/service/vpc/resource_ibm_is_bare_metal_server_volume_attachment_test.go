// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIBMISBMSVolumeAttachment_basic(t *testing.T) {
	var bmsVolAtt string
	vpcname := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	bmsname := fmt.Sprintf("tf-server-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfip-subnet-%d", acctest.RandIntRange(10, 100))
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	sshname := fmt.Sprintf("tf-sshname-%d", acctest.RandIntRange(10, 100))
	volname := fmt.Sprintf("tf-vol-%d", acctest.RandIntRange(10, 100))
	attname1 := fmt.Sprintf("tf-att-%d", acctest.RandIntRange(10, 100))
	attname2 := fmt.Sprintf("tf-att-upd-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISBMSVolumeAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISBMSVolumeAttachmentConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname1, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISBMSVolumeAttachmentExists("ibm_is_bare_metal_server_volume_attachment.testacc_att", bmsVolAtt),
					resource.TestCheckResourceAttr(
						"ibm_is_bare_metal_server_volume_attachment.testacc_att", "name", attname1),
					resource.TestCheckResourceAttr(
						"ibm_is_bare_metal_server_volume_attachment.testacc_att", "delete_volume_on_bare_metal_server_delete", "false"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_bare_metal_server_volume_attachment.testacc_att", "status"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_bare_metal_server_volume_attachment.testacc_att", "href"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_bare_metal_server_volume_attachment.testacc_att", "volume_attachment_id"),
				),
			},
			{
				Config: testAccCheckIBMISBMSVolumeAttachmentConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname2, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISBMSVolumeAttachmentExists("ibm_is_bare_metal_server_volume_attachment.testacc_att", bmsVolAtt),
					resource.TestCheckResourceAttr(
						"ibm_is_bare_metal_server_volume_attachment.testacc_att", "name", attname2),
					resource.TestCheckResourceAttr(
						"ibm_is_bare_metal_server_volume_attachment.testacc_att", "delete_volume_on_bare_metal_server_delete", "true"),
				),
			},
		},
	})
}

func testAccCheckIBMISBMSVolumeAttachmentDestroy(s *terraform.State) error {
	sess, _ := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_bare_metal_server_volume_attachment" {
			continue
		}
		bmsId, id, err := parseVolAttTerraformID(rs.Primary.ID)
		if err != nil {
			return err
		}
		getOpts := &vpcv1.GetBareMetalServerVolumeAttachmentOptions{
			BareMetalServerID: &bmsId,
			ID:                &id,
		}
		_, _, err = sess.GetBareMetalServerVolumeAttachment(getOpts)
		if err == nil {
			return fmt.Errorf("bare metal server volume attachment still exists: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccCheckIBMISBMSVolumeAttachmentExists(n string, bmsVolAtt string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return errors.New("No Record ID is set")
		}
		bmsId, id, err := parseVolAttTerraformID(rs.Primary.ID)
		if err != nil {
			return err
		}
		sess, _ := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		getOpts := &vpcv1.GetBareMetalServerVolumeAttachmentOptions{
			BareMetalServerID: &bmsId,
			ID:                &id,
		}
		volAttIntf, _, err := sess.GetBareMetalServerVolumeAttachment(getOpts)
		if err != nil {
			return err
		}
		volAtt, ok := volAttIntf.(*vpcv1.BareMetalServerVolumeAttachment)
		if !ok || volAtt == nil {
			return fmt.Errorf("bare metal server volume attachment not found: %s", n)
		}
		bmsVolAtt = *volAtt.ID
		return nil
	}
}

func testAccCheckIBMISBMSVolumeAttachmentConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname string, deleteOnServerDelete bool) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}

	resource "ibm_is_subnet" "testacc_subnet" {
		name                     = "%s"
		vpc                      = ibm_is_vpc.testacc_vpc.id
		zone                     = "eu-gb-1"
		total_ipv4_address_count = 16
	}

	resource "ibm_is_ssh_key" "testacc_sshkey" {
		name       = "%s"
		public_key = "%s"
	}

	resource "ibm_is_bare_metal_server" "testacc_bms" {
		profile = "%s"
		name    = "%s"
		image   = "%s"
		zone    = "eu-gb-1"
		keys    = [ibm_is_ssh_key.testacc_sshkey.id]
		primary_network_interface {
			subnet = ibm_is_subnet.testacc_subnet.id
		}
		vpc = ibm_is_vpc.testacc_vpc.id
	}

	resource "ibm_is_volume" "testacc_vol" {
		name     = "%s"
		profile  = "sdp"
		zone     = "eu-gb-1"
		capacity = 10000
	}

	resource "ibm_is_bare_metal_server_volume_attachment" "testacc_att" {
		bare_metal_server                    = ibm_is_bare_metal_server.testacc_bms.id
		volume                               = ibm_is_volume.testacc_vol.id
		name                                 = "%s"
		delete_volume_on_bare_metal_server_delete = %t
		delete_volume_on_attachment_delete   = false
	}
	`, vpcname, subnetname, sshname, publicKey,
		acc.IsBareMetalServerProfileName, bmsname, acc.IsBareMetalServerImage,
		volname, attname, deleteOnServerDelete)
}
