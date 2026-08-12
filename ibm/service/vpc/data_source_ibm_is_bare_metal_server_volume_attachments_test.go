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

func TestAccIBMISBMSVolumeAttachmentsDataSource_basic(t *testing.T) {
	resName := "data.ibm_is_bare_metal_server_volume_attachments.test1"
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
				Config: testAccCheckIBMISBMSVolumeAttachmentsDataSourceConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISBareMetalServerExists("ibm_is_bare_metal_server.testacc_bms", server),
					resource.TestCheckResourceAttrSet(resName, "bare_metal_server"),
					resource.TestCheckResourceAttrSet(resName, "volume_attachments.#"),
					resource.TestCheckResourceAttrSet(resName, "volume_attachments.0.volume_attachment_id"),
					resource.TestCheckResourceAttrSet(resName, "volume_attachments.0.name"),
					resource.TestCheckResourceAttrSet(resName, "volume_attachments.0.status"),
					resource.TestCheckResourceAttrSet(resName, "volume_attachments.0.type"),
					resource.TestCheckResourceAttrSet(resName, "volume_attachments.0.href"),
				),
			},
		},
	})
}

func testAccCheckIBMISBMSVolumeAttachmentsDataSourceConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname string) string {
	return testAccCheckIBMISBMSWithSdpVolumeAttachmentConfig(vpcname, subnetname, sshname, publicKey, bmsname, volname, attname) + `
	data "ibm_is_bare_metal_server_volume_attachments" "test1" {
		bare_metal_server = ibm_is_bare_metal_server.testacc_bms.id
		depends_on        = [ibm_is_bare_metal_server_volume_attachment.testacc_att]
	}`
}
