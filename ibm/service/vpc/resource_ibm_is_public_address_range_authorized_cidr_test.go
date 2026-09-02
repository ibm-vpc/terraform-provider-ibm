// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMPublicAddressRangeAuthorizedCIDRBasic(t *testing.T) {
	var conf vpcv1.PublicAddressRangeAuthorizedCIDR
	name := fmt.Sprintf("tf-authcidr-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeAuthorizedCIDRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeAuthorizedCIDRConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeAuthorizedCIDRExists(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", conf),
					resource.TestCheckResourceAttr(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "name", name),
					resource.TestCheckResourceAttr(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "ip_version", "ipv6"),
					resource.TestCheckResourceAttr(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "availability_mode", "zonal"),
					resource.TestCheckResourceAttr(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "network_prefix_length", "64"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "cidr"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "crn"),
					resource.TestCheckResourceAttrSet(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "lifecycle_state"),
				),
			},
		},
	})
}

func TestAccIBMPublicAddressRangeAuthorizedCIDRUpdate(t *testing.T) {
	var conf vpcv1.PublicAddressRangeAuthorizedCIDR
	name := fmt.Sprintf("tf-authcidr-%d", acctest.RandIntRange(10, 100))
	updatedName := fmt.Sprintf("%s-updated", name)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeAuthorizedCIDRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeAuthorizedCIDRConfig(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeAuthorizedCIDRExists(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", conf),
					resource.TestCheckResourceAttr(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "name", name),
				),
			},
			{
				Config: testAccCheckIBMPublicAddressRangeAuthorizedCIDRConfig(updatedName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeAuthorizedCIDRExists(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", conf),
					resource.TestCheckResourceAttr(
						"ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr", "name", updatedName),
				),
			},
		},
	})
}

func TestAccIBMPublicAddressRangeAuthorizedCIDRImport(t *testing.T) {
	name := fmt.Sprintf("tf-authcidr-import-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeAuthorizedCIDRDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeAuthorizedCIDRConfig(name),
			},
			{
				ResourceName:      "ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMPublicAddressRangeAuthorizedCIDRConfig(name string) string {
	return fmt.Sprintf(`
resource "ibm_is_public_address_range_authorized_cidr" "testacc_auth_cidr" {
  name                  = "%s"
  ip_version            = "ipv6"
  availability_mode     = "zonal"
  zone                  = "%s"
  network_prefix_length = 64
}
`, name, acc.ISZoneName)
}

func testAccCheckIBMPublicAddressRangeAuthorizedCIDRExists(n string, obj vpcv1.PublicAddressRangeAuthorizedCIDR) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set for %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getOptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDROptions{}
		getOptions.SetID(rs.Primary.ID)

		authorizedCIDR, _, err := vpcClient.GetPublicAddressRangeAuthorizedCIDR(getOptions)
		if err != nil {
			return err
		}

		obj = *authorizedCIDR
		return nil
	}
}

func testAccCheckIBMPublicAddressRangeAuthorizedCIDRDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_public_address_range_authorized_cidr" {
			continue
		}

		getOptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDROptions{}
		getOptions.SetID(rs.Primary.ID)

		_, response, err := vpcClient.GetPublicAddressRangeAuthorizedCIDR(getOptions)
		if err == nil {
			return fmt.Errorf("PublicAddressRangeAuthorizedCIDR still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for PublicAddressRangeAuthorizedCIDR (%s) destruction: %s", rs.Primary.ID, err)
		}
	}
	return nil
}
