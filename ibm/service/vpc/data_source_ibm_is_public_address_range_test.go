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

func TestAccIBMIsPublicAddressRangeDataSourceBasic(t *testing.T) {
	ipv4AddressCount := "16"
	// ipv4AddressCountUpdate := "8"
	name := fmt.Sprintf("tf-name-par%d", acctest.RandIntRange(10, 100))
	vpcName := fmt.Sprintf("tf-name-vpc%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeDataSourceConfigBasic(vpcName, name, ipv4AddressCount),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "cidr"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "ipv4_address_count"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "resource_group.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "network_prefix_length"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.is_public_address_range_instance", "profile.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeDataSourceConfigBasic(vpcName, name, ipv4AddressCount string) string {
	return testAccCheckIBMPublicAddressRangeConfigBasic(vpcName, name, ipv4AddressCount) + fmt.Sprintf(`
	data "ibm_is_public_address_range" "is_public_address_range_instance" {
		identifier = ibm_is_public_address_range.public_address_range_instance.id
	}
`)
}

// TestAccIBMIsPublicAddressRangeDataSourceVNITarget verifies that when an IPv6 PAR targets
// a VNI, the data source correctly populates target.virtual_network_interface fields.
func TestAccIBMIsPublicAddressRangeDataSourceVNITarget(t *testing.T) {
	vpcName := fmt.Sprintf("tf-vpc-par-ds-%d", acctest.RandIntRange(10, 100))
	subnetName := fmt.Sprintf("tf-subnet-par-ds-%d", acctest.RandIntRange(10, 100))
	vniName := fmt.Sprintf("tf-vni-par-ds-%d", acctest.RandIntRange(10, 100))
	authCIDRName := fmt.Sprintf("tf-authcidr-ds-%d", acctest.RandIntRange(10, 100))
	parName := fmt.Sprintf("tf-par-ds-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsPublicAddressRangeDataSourceVNIConfig(vpcName, subnetName, vniName, authCIDRName, parName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "cidr"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range.par_vni_ds", "ip_version", "ipv6"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "network_prefix_length"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "authorized_cidr.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "authorized_cidr.0.id"),
					// target.virtual_network_interface must be populated
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range.par_vni_ds", "target.#", "1"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range.par_vni_ds", "target.0.virtual_network_interface.#", "1"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "target.0.virtual_network_interface.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "target.0.virtual_network_interface.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "target.0.virtual_network_interface.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range.par_vni_ds", "target.0.virtual_network_interface.0.name"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range.par_vni_ds", "target.0.virtual_network_interface.0.resource_type", "virtual_network_interface"),
					// target.vpc must be absent for an IPv6 PAR
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range.par_vni_ds", "target.0.vpc.#", "0"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeDataSourceVNIConfig(vpcName, subnetName, vniName, authCIDRName, parName string) string {
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

data "ibm_is_public_address_range" "par_vni_ds" {
  identifier = ibm_is_public_address_range.testacc_par.id
}
`, vpcName, subnetName, acc.ISZoneName, vniName, authCIDRName, acc.ISZoneName, parName)
}
