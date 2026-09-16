// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsPublicAddressRangesDataSourceBasic(t *testing.T) {
	ipv4AddressCount := "16"
	name := fmt.Sprintf("tf-name-par%d", acctest.RandIntRange(10, 100))
	vpcName := fmt.Sprintf("tf-name-vpc%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangesDataSourceConfigBasic(vpcName, name, ipv4AddressCount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.cidr"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.ipv4_address_count"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.network_prefix_length"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.is_public_address_ranges_instance", "public_address_ranges.0.profile.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangesDataSourceConfigBasic(vpcName, name, ipv4AddressCount string) string {
	return testAccCheckIBMPublicAddressRangeConfigBasic(vpcName, name, ipv4AddressCount) + fmt.Sprintf(`
	data "ibm_is_public_address_ranges" "is_public_address_ranges_instance" {
	  depends_on = [ibm_is_public_address_range.public_address_range_instance]
	}
`)
}

// TestAccIBMIsPublicAddressRangesDataSourceVNITarget verifies that when an IPv6 PAR targets
// a VNI, the list data source correctly populates target.virtual_network_interface on that entry.
func TestAccIBMIsPublicAddressRangesDataSourceVNITarget(t *testing.T) {
	vpcName := fmt.Sprintf("tf-vpc-pars-ds-%d", acctest.RandIntRange(10, 100))
	subnetName := fmt.Sprintf("tf-subnet-pars-ds-%d", acctest.RandIntRange(10, 100))
	vniName := fmt.Sprintf("tf-vni-pars-ds-%d", acctest.RandIntRange(10, 100))
	authCIDRName := fmt.Sprintf("tf-authcidr-pars-%d", acctest.RandIntRange(10, 100))
	parName := fmt.Sprintf("tf-par-pars-ds-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsPublicAddressRangesDataSourceVNIConfig(vpcName, subnetName, vniName, authCIDRName, parName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_ranges.pars_vni_ds", "public_address_ranges.#"),
					// Verify the PAR resource itself carries the VNI target
					resource.TestCheckResourceAttr("ibm_is_public_address_range.testacc_par", "ip_version", "ipv6"),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.testacc_par", "target.#", "1"),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.testacc_par", "target.0.virtual_network_interface.#", "1"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.testacc_par", "target.0.virtual_network_interface.0.id"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.testacc_par", "target.0.virtual_network_interface.0.crn"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.testacc_par", "target.0.virtual_network_interface.0.href"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.testacc_par", "target.0.virtual_network_interface.0.name"),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.testacc_par", "target.0.virtual_network_interface.0.resource_type", "virtual_network_interface"),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.testacc_par", "target.0.vpc.#", "0"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangesDataSourceVNIConfig(vpcName, subnetName, vniName, authCIDRName, parName string) string {
	return fmt.Sprintf(`
resource "ibm_is_vpc" "testacc_vpc" {
  name = "%s"
}

resource "ibm_is_subnet" "testacc_subnet" {
  name                     = "%s"
  vpc                      = ibm_is_vpc.testacc_vpc.id
  zone                     = "%s"
  total_ipv4_address_count = 256
}

resource "ibm_is_virtual_network_interface" "testacc_vni" {
  name   = "%s"
  subnet = ibm_is_subnet.testacc_subnet.id
}

resource "ibm_is_public_address_range_authorized_cidr" "testacc_auth_cidr" {
  name                  = "%s"
  ip_version            = "ipv6"
  availability_mode     = "zonal"
  zone                  = "%s"
  network_prefix_length = 64
}

resource "ibm_is_public_address_range" "testacc_par" {
  name                  = "%s"
  network_prefix_length = 112
  authorized_cidr {
    id = ibm_is_public_address_range_authorized_cidr.testacc_auth_cidr.id
  }
  target {
    virtual_network_interface {
      id = ibm_is_virtual_network_interface.testacc_vni.id
    }
  }
}

data "ibm_is_public_address_ranges" "pars_vni_ds" {
  depends_on = [ibm_is_public_address_range.testacc_par]
}
`, vpcName, subnetName, acc.ISZoneName, vniName, authCIDRName, acc.ISZoneName, parName)
}
