// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsVirtualNetworkInterfacesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualNetworkInterfacesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_group.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_group.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_group.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.security_groups.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.subnet.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.vpc.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.zone.0.name"),
				),
			},
		},
	})
}
func TestAccIBMIsVirtualNetworkInterfacesDataSourceVniBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualNetworkInterfacesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_group.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_group.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_group.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.security_groups.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.subnet.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.vpc.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.zone.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.allow_ip_spoofing"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.is_virtual_network_interfaces", "virtual_network_interfaces.0.enable_infrastructure_nat"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVirtualNetworkInterfacesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_virtual_network_interfaces" "is_virtual_network_interfaces" {
		}
	`)
}

// TestAccIBMIsVirtualNetworkInterfacesDataSourcePublicAddressRanges verifies that
// public_address_ranges is populated on the list-VNI data source when an
// IPv6 PAR (via an authorized CIDR) targets a VNI.
func TestAccIBMIsVirtualNetworkInterfacesDataSourcePublicAddressRanges(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-vnis-par-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tf-subnet-vnis-par-%d", acctest.RandIntRange(10, 100))
	vniname := fmt.Sprintf("tf-vni-vnis-par-%d", acctest.RandIntRange(10, 100))
	authCIDRName := fmt.Sprintf("tf-authcidr-vnis-%d", acctest.RandIntRange(10, 100))
	parName := fmt.Sprintf("tf-par-vnis-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsVirtualNetworkInterfacesDataSourcePARConfig(vpcname, subnetname, vniname, authCIDRName, parName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interfaces.vnis_ds", "virtual_network_interfaces.#"),
					// The VNI we created must appear in the list and carry public_address_ranges
					resource.TestCheckResourceAttrSet("ibm_is_virtual_network_interface.testacc_vni", "public_address_ranges.#"),
					resource.TestCheckResourceAttrSet("ibm_is_virtual_network_interface.testacc_vni", "public_address_ranges.0.id"),
					resource.TestCheckResourceAttrSet("ibm_is_virtual_network_interface.testacc_vni", "public_address_ranges.0.crn"),
					resource.TestCheckResourceAttrSet("ibm_is_virtual_network_interface.testacc_vni", "public_address_ranges.0.cidr"),
					resource.TestCheckResourceAttrSet("ibm_is_virtual_network_interface.testacc_vni", "public_address_ranges.0.href"),
					resource.TestCheckResourceAttrSet("ibm_is_virtual_network_interface.testacc_vni", "public_address_ranges.0.name"),
					resource.TestCheckResourceAttr("ibm_is_virtual_network_interface.testacc_vni", "public_address_ranges.0.resource_type", "public_address_range"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVirtualNetworkInterfacesDataSourcePARConfig(vpcname, subnetname, vniname, authCIDRName, parName string) string {
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

data "ibm_is_virtual_network_interfaces" "vnis_ds" {
  depends_on = [ibm_is_public_address_range.testacc_par]
}
`, vpcname, subnetname, acc.ISZoneName, vniname, authCIDRName, acc.ISZoneName, parName)
}
