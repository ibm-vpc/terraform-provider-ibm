// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmIsSharesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsSharesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "total_count"),
				),
			},
		},
	})
}

func TestAccIbmIsSharesDataSourceAllArgs(t *testing.T) {
	shareIops := 3000
	shareName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	shareSize := acctest.RandIntRange(10, 16000)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsSharesDataSourceConfig(shareIops, shareName, shareSize),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.encryption"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.iops"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "shares.0.size"),
					resource.TestCheckResourceAttrSet("data.ibm_is_shares.is_shares", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIbmIsSharesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		resource "ibm_is_share" "is_share" {
			zone = "us-south-2"
			size = 200
		}

		data "ibm_is_shares" "is_shares" {
		}
	`)
}

func testAccCheckIbmIsSharesDataSourceConfig(shareIops int, shareName string, shareSize int) string {
	return fmt.Sprintf(`
		resource "ibm_is_share" "is_share" {
			zone = "us-south-2"
			name = "%s
			iops = %d
			size = %d

		}

		data "ibm_is_shares" "is_shares" {
		}
	`, shareName, shareIops, shareSize)
}
