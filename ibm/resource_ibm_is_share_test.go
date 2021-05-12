// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIbmIsShareBasic(t *testing.T) {
	var conf vpcv1.Share

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmIsShareDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareConfigBasic(),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsShareExists("ibm_is_share.is_share", conf),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmIsShareConfigBasic(),
				Check:  resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}

func TestAccIbmIsShareAllArgs(t *testing.T) {
	var conf vpcv1.Share
	iops := 200
	name := fmt.Sprintf("tf-fs-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tf-vpc-name-%d", acctest.RandIntRange(10, 100))
	size := acctest.RandIntRange(10, 3000)

	nameUpdate := fmt.Sprintf("tf-fs-name-%d", acctest.RandIntRange(10, 100))
	//sizeUpdate := acctest.RandIntRange(10, 16000)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmIsShareDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsShareConfig(vpcname, iops, name, shareProfileName, size),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsShareExists("ibm_is_share.is_share", conf),
					//resource.TestCheckResourceAttr("ibm_is_share.is_share", "iops", strconv.Itoa(iops)),
					resource.TestCheckResourceAttr("ibm_is_share.is_share", "name", name),
					resource.TestCheckResourceAttr("ibm_is_share.is_share", "size", strconv.Itoa(size)),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmIsShareConfig(vpcname, iops, nameUpdate, shareProfileName, 3000),
				Check: resource.ComposeAggregateTestCheckFunc(
					//resource.TestCheckResourceAttr("ibm_is_share.is_share", "iops", strconv.Itoa(iops)),
					resource.TestCheckResourceAttr("ibm_is_share.is_share", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_share.is_share", "size", strconv.Itoa(size)),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_share.is_share",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIbmIsShareConfigBasic() string {
	return `
		resource "ibm_is_share" "is_share" {
			zone = "us-south-2"
			size = 200
		}
	`
}

func testAccCheckIbmIsShareConfig(vpcName string, iops int, name, profile string, size int) string {
	return fmt.Sprintf(`

	data "ibm_resource_group" "group" {
		is_default = "true"
	}
	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}
	resource "ibm_is_share" "is_share" {
		
		#iops = %d
		name = "%s"
		profile = "%s"
		resource_group = data.ibm_resource_group.group.id
		size = %d
		targets {
			vpc = ibm_is_vpc.testacc_vpc.id
		}
		zone = "us-south-2"
	}
	`, vpcName, iops, name, profile, size)
}

func testAccCheckIbmIsShareExists(n string, obj vpcv1.Share) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := testAccProvider.Meta().(ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getShareOptions := &vpcv1.GetShareOptions{}

		getShareOptions.SetID(rs.Primary.ID)

		share, _, err := vpcClient.GetShare(getShareOptions)
		if err != nil {
			return err
		}

		obj = *share
		return nil
	}
}

func testAccCheckIbmIsShareDestroy(s *terraform.State) error {
	vpcClient, err := testAccProvider.Meta().(ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_share" {
			continue
		}

		getShareOptions := &vpcv1.GetShareOptions{}

		getShareOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := vpcClient.GetShare(getShareOptions)

		if err == nil {
			return fmt.Errorf("Share still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for Share (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
