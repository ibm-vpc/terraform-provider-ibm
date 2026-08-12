// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"strings"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMISBMSVolumeAttachmentDataSource_basic(t *testing.T) {
	resName := "data.ibm_is_bare_metal_server_volume_attachment.test1"
	var server string
	vpcname := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	bmsname := fmt.Sprintf("tf-server-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfip-subnet-%d", acctest.RandIntRange(10, 100))
	volname := fmt.Sprintf("tf-vol-%d", acctest.RandIntRange(10, 100))
	attname := fmt.Sprintf("tf-att-%d", acctest.RandIntRange(10, 100))
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	sshname := fmt.Sprintf("tf-sshname-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISBMSVolumeAttachmentDataSourceConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISBareMetalServerExists("ibm_is_bare_metal_server.testacc_bms", server),
					resource.TestCheckResourceAttr(resName, "name", attname),
					resource.TestCheckResourceAttrSet(resName, "bare_metal_server"),
					resource.TestCheckResourceAttrSet(resName, "volume_attachment_id"),
					resource.TestCheckResourceAttrSet(resName, "status"),
					resource.TestCheckResourceAttrSet(resName, "type"),
					resource.TestCheckResourceAttrSet(resName, "href"),
					resource.TestCheckResourceAttrSet(resName, "protocol"),
					resource.TestCheckResourceAttrSet(resName, "volume.#"),
				),
			},
		},
	})
}

func testAccCheckIBMISBMSVolumeAttachmentDataSourceConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname string) string {
	return testAccCheckIBMISBMSWithSdpVolumeAttachmentConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname) + fmt.Sprintf(`
	data "ibm_is_bare_metal_server_volume_attachment" "test1" {
		bare_metal_server = ibm_is_bare_metal_server.testacc_bms.id
		name              = "%s"
		depends_on        = [ibm_is_bare_metal_server_volume_attachment.testacc_att]
	}`, attname)
}

func testAccCheckIBMISBMSWithSdpVolumeAttachmentConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname string) string {
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
		bare_metal_server = ibm_is_bare_metal_server.testacc_bms.id
		volume            = ibm_is_volume.testacc_vol.id
		name              = "%s"
		delete_volume_on_attachment_delete = false
	}
`, vpcname, subnetname, sshname, publicKey,
		acc.IsBareMetalServerProfileName, bmsname, acc.IsBareMetalServerImage,
		volname, attname)
}
