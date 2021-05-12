// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmIsShareTargetsDataSource(t *testing.T) {
	vpcName := fmt.Sprintf("tf-vpc-%d", acctest.RandIntRange(10, 100))
	targetName := fmt.Sprintf("tf-share-target-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareTargetsDataSourceConfigBasic(vpcName, targetName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.0.mount_path"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share_targets.is_share_targets", "targets.0.resource_type"),
				),
			},
		},
	})
}

func testAccCheckIbmIsShareTargetsDataSourceConfigBasic(vpcName, targetName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_share" "is_share" {
			zone = "us-south-2"
			size = 200
		}

		resource "ibm_is_vpc" "testacc_vpc" {
			name = "%s"
		}

		resource "ibm_is_share_target" "is_share_target" {
			share = ibm_is_share.is_share.id
			name = "%s"
			vpc = ibm_is_vpc.testacc_vpc.id
		}

		data "ibm_is_share_targets" "is_share_targets" {
			share = ibm_is_share_target.is_share_target.share
		}
	`, vpcName, targetName)
}
