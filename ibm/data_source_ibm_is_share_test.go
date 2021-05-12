// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmIsShareDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "encryption"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "encryption_key"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "iops"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "profile.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "resource_group"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "size"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "targets.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "zone"),
				),
			},
		},
	})
}

func TestAccIbmIsShareDataSourceAllArgs(t *testing.T) {
	shareIops := 3000
	shareName := fmt.Sprintf("tf-name-%d", acctest.RandIntRange(10, 100))
	shareSize := acctest.RandIntRange(10, 16000)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareDataSourceConfig(shareIops, shareName, shareSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "encryption"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "encryption_key"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "iops"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "profile.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "resource_group"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "size"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "targets.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "targets.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "targets.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_share.is_share", "targets.0.name", shareName),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "targets.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_share.is_share", "zone"),
				),
			},
		},
	})
}

func testAccCheckIbmIsShareDataSourceConfigBasic() string {
	return `
		resource "ibm_is_share" "is_share" {
			zone = "us-south-2"
			size = 200
		}

		data "ibm_is_share" "is_share" {
			id = ibm_is_share.is_share.id
		}
	`
}

func testAccCheckIbmIsShareDataSourceConfig(shareIops int, shareName string, shareSize int) string {
	return fmt.Sprintf(`
		resource "ibm_is_share" "is_share" {
			zone = "us-south-2"
			name = "%s"
			iops = %d
			size = %d

		}
		data "ibm_is_share" "is_share" {
			id = ibm_is_share.is_share.id
		}
	`, shareName, shareIops, shareSize)
}
