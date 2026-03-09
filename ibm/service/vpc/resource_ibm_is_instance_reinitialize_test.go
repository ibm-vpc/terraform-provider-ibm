// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
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

func TestAccIBMISInstanceReinitialize_basic(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf-server-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfip-subnet-%d", acctest.RandIntRange(10, 100))
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	sshname := fmt.Sprintf("tf-sshname-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceReinitializeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceReinitializeConfig(vpcname, subnetname, sshname, publicKey, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_is_instance.testacc_instance", "name", name),
					resource.TestCheckResourceAttr(
						"ibm_is_instance_reinitialize.testacc_reinit", "auto_stop", "true"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_instance_reinitialize.testacc_reinit", "status"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_instance_reinitialize.testacc_reinit", "boot_volume_attachment_id"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_instance_reinitialize.testacc_reinit", "reinitialized_at"),
				),
			},
		},
	})
}

func TestAccIBMISInstanceReinitialize_withTriggers(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf-server-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfip-subnet-%d", acctest.RandIntRange(10, 100))
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	sshname := fmt.Sprintf("tf-sshname-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceReinitializeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceReinitializeWithTriggersConfig(vpcname, subnetname, sshname, publicKey, name, "v1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_is_instance_reinitialize.testacc_reinit", "triggers.version", "v1"),
				),
			},
			{
				Config: testAccCheckIBMISInstanceReinitializeWithTriggersConfig(vpcname, subnetname, sshname, publicKey, name, "v2"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_is_instance_reinitialize.testacc_reinit", "triggers.version", "v2"),
				),
			},
		},
	})
}

func TestAccIBMISInstanceReinitialize_withUserData(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	name := fmt.Sprintf("tf-server-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfip-subnet-%d", acctest.RandIntRange(10, 100))
	publicKey := strings.TrimSpace(`
ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCKVmnMOlHKcZK8tpt3MP1lqOLAcqcJzhsvJcjscgVERRN7/9484SOBJ3HSKxxNG5JN8owAjy5f9yYwcUg+JaUVuytn5Pv3aeYROHGGg+5G346xaq3DAwX6Y5ykr2fvjObgncQBnuU5KHWCECO/4h8uWuwh/kfniXPVjFToc+gnkqA+3RKpAecZhFXwfalQ9mMuYGFxn+fwn8cYEApsJbsEmb0iJwPiZ5hjFC8wREuiTlhPHDgkBLOiycd20op2nXzDbHfCHInquEe/gYxEitALONxm0swBOwJZwlTDOB7C6y2dzlrtxr1L59m7pCkWI4EtTRLvleehBoj3u7jB4usR
`)
	sshname := fmt.Sprintf("tf-sshname-%d", acctest.RandIntRange(10, 100))
	userData := "#!/bin/bash\necho 'Hello from reinitialized instance'"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISInstanceReinitializeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISInstanceReinitializeWithUserDataConfig(vpcname, subnetname, sshname, publicKey, name, userData),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_is_instance_reinitialize.testacc_reinit", "user_data", userData),
				),
			},
		},
	})
}

func testAccCheckIBMISInstanceReinitializeDestroy(s *terraform.State) error {
	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_instance_reinitialize" {
			continue
		}

		// Instance should still exist after reinitialize resource is destroyed
		instanceID := rs.Primary.ID
		getInstanceOptions := &vpcv1.GetInstanceOptions{
			ID: &instanceID,
		}
		_, response, err := sess.GetInstance(getInstanceOptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				// Instance was deleted, which is fine
				continue
			}
			return fmt.Errorf("Error checking for instance: %s\n%s", err, response)
		}
	}

	return nil
}

func testAccCheckIBMISInstanceReinitializeConfig(vpcname, subnetname, sshname, publicKey, name string) string {
	return fmt.Sprintf(`
	data "ibm_is_images" "im_images" {
	}

	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}

	resource "ibm_is_subnet" "testacc_subnet" {
		name            = "%s"
		vpc             = ibm_is_vpc.testacc_vpc.id
		zone            = "%s"
		ipv4_cidr_block = "%s"
	}

	resource "ibm_is_ssh_key" "testacc_sshkey" {
		name       = "%s"
		public_key = "%s"
	}

	resource "ibm_is_instance" "testacc_instance" {
		name    = "%s"
		image   = data.ibm_is_images.im_images.images.4.id
		profile = "bx2-2x8"
		primary_network_interface {
			subnet = ibm_is_subnet.testacc_subnet.id
		}
		vpc  = ibm_is_vpc.testacc_vpc.id
		zone = "%s"
		keys = [ibm_is_ssh_key.testacc_sshkey.id]

		lifecycle {
			ignore_changes = [boot_volume]
		}
	}

	resource "ibm_is_instance_reinitialize" "testacc_reinit" {
		instance   = ibm_is_instance.testacc_instance.id
		image      = data.ibm_is_images.im_images.images.5.id
		keys       = [ibm_is_ssh_key.testacc_sshkey.id]
		auto_stop  = true
		auto_start = true
	}
	`, vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, sshname, publicKey, name, acc.ISZoneName)
}

func testAccCheckIBMISInstanceReinitializeWithTriggersConfig(vpcname, subnetname, sshname, publicKey, name, version string) string {
	return fmt.Sprintf(`
	data "ibm_is_images" "im_images" {
	}

	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}

	resource "ibm_is_subnet" "testacc_subnet" {
		name            = "%s"
		vpc             = ibm_is_vpc.testacc_vpc.id
		zone            = "%s"
		ipv4_cidr_block = "%s"
	}

	resource "ibm_is_ssh_key" "testacc_sshkey" {
		name       = "%s"
		public_key = "%s"
	}

	resource "ibm_is_instance" "testacc_instance" {
		name    = "%s"
		image   = data.ibm_is_images.im_images.images.4.id
		profile = "bx2-2x8"
		primary_network_interface {
			subnet = ibm_is_subnet.testacc_subnet.id
		}
		vpc  = ibm_is_vpc.testacc_vpc.id
		zone = "%s"
		keys = [ibm_is_ssh_key.testacc_sshkey.id]

		lifecycle {
			ignore_changes = [boot_volume]
		}
	}

	resource "ibm_is_instance_reinitialize" "testacc_reinit" {
		instance   = ibm_is_instance.testacc_instance.id
		image      = data.ibm_is_images.im_images.images.5.id
		keys       = [ibm_is_ssh_key.testacc_sshkey.id]
		auto_stop  = true
		auto_start = true

		triggers = {
			version = "%s"
		}
	}
	`, vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, sshname, publicKey, name, acc.ISZoneName, version)
}

func testAccCheckIBMISInstanceReinitializeWithUserDataConfig(vpcname, subnetname, sshname, publicKey, name, userData string) string {
	return fmt.Sprintf(`
	data "ibm_is_images" "im_images" {
	}

	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}

	resource "ibm_is_subnet" "testacc_subnet" {
		name            = "%s"
		vpc             = ibm_is_vpc.testacc_vpc.id
		zone            = "%s"
		ipv4_cidr_block = "%s"
	}

	resource "ibm_is_ssh_key" "testacc_sshkey" {
		name       = "%s"
		public_key = "%s"
	}

	resource "ibm_is_instance" "testacc_instance" {
		name    = "%s"
		image   = data.ibm_is_images.im_images.images.4.id
		profile = "bx2-2x8"
		primary_network_interface {
			subnet = ibm_is_subnet.testacc_subnet.id
		}
		vpc  = ibm_is_vpc.testacc_vpc.id
		zone = "%s"
		keys = [ibm_is_ssh_key.testacc_sshkey.id]

		lifecycle {
			ignore_changes = [boot_volume]
		}
	}

	resource "ibm_is_instance_reinitialize" "testacc_reinit" {
		instance   = ibm_is_instance.testacc_instance.id
		image      = data.ibm_is_images.im_images.images.5.id
		keys       = [ibm_is_ssh_key.testacc_sshkey.id]
		user_data  = "%s"
		auto_stop  = true
		auto_start = true
	}
	`, vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, sshname, publicKey, name, acc.ISZoneName, userData)
}

// Made with Bob
