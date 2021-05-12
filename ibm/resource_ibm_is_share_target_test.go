// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIbmIsShareTargetBasic(t *testing.T) {
	var conf vpcv1.ShareTarget
	vpcname := fmt.Sprintf("tf-vpc-name-%d", acctest.RandIntRange(10, 100))
	targetName := fmt.Sprintf("tf-target-%d", acctest.RandIntRange(10, 100))
	targetNameUpdate := fmt.Sprintf("tf-target-%d", acctest.RandIntRange(10, 100))
	sname := fmt.Sprintf("tf-fs-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmIsShareTargetDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareTargetConfigBasic(vpcname, sname, targetName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsShareTargetExists("ibm_is_share_target.is_share_target", conf),
					resource.TestCheckResourceAttr("ibm_is_share_target.is_share_target", "name", targetName),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmIsShareTargetConfigBasic(vpcname, sname, targetNameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_share_target.is_share_target", "name", targetNameUpdate),
				),
			},
		},
	})
}

func testAccCheckIbmIsShareTargetConfigBasic(vpcName, sname, targetName string) string {
	return fmt.Sprintf(`
	data "ibm_resource_group" "group" {
		is_default = "true"
	}
	resource "ibm_is_share" "is_share" {
		zone = "us-south-2"
		size = 200
		name = "%s"
	}
	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}
	resource "ibm_is_share_target" "is_share_target" {
		share = ibm_is_share.is_share.id
		vpc = ibm_is_vpc.testacc_vpc.id
		name = "%s"
	}
	`, sname, vpcName, targetName)
}

func testAccCheckIbmIsShareTargetExists(n string, obj vpcv1.ShareTarget) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := testAccProvider.Meta().(ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getShareTargetOptions := &vpcv1.GetShareTargetOptions{}

		parts, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		getShareTargetOptions.SetShareID(parts[0])
		getShareTargetOptions.SetID(parts[1])

		shareTarget, _, err := vpcClient.GetShareTarget(getShareTargetOptions)
		if err != nil {
			return err
		}

		obj = *shareTarget
		return nil
	}
}

func testAccCheckIbmIsShareTargetDestroy(s *terraform.State) error {
	vpcClient, err := testAccProvider.Meta().(ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_share_target" {
			continue
		}

		getShareTargetOptions := &vpcv1.GetShareTargetOptions{}

		parts, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		getShareTargetOptions.SetShareID(parts[0])
		getShareTargetOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetShareTarget(getShareTargetOptions)

		if err == nil {
			return fmt.Errorf("ShareTarget still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for ShareTarget (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
