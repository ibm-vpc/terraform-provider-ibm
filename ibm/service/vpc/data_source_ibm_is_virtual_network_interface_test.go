// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsVirtualNetworkInterfaceDataSourceBasic(t *testing.T) {
	var conf vpcv1.VirtualNetworkInterface
	vpcname := fmt.Sprintf("tfvpngw-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfvpngw-subnet-%d", acctest.RandIntRange(10, 100))
	vniname := fmt.Sprintf("tfvpngw-createname-%d", acctest.RandIntRange(10, 100))
	tag1 := "env:test"
	tag2 := "env:dev"
	tag3 := "env:prod"
	enable_infrastructure_nat := true
	allow_ip_spoofing := true
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualNetworkInterfaceDataSourceConfigBasic(vpcname, subnetname, vniname, tag1, tag2, tag3, enable_infrastructure_nat, allow_ip_spoofing, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMIsVirtualNetworkInterfaceExists("ibm_is_virtual_network_interface.testacc_vni", conf),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "virtual_network_interface"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "resource_group.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "security_groups.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "subnet.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "vpc.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "tags.#", "3"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "zone.0.name"),
				),
			},
		},
	})
}
func TestAccIBMIsVirtualNetworkInterfaceDataSourceVniBasic(t *testing.T) {
	var conf vpcv1.VirtualNetworkInterface
	vpcname := fmt.Sprintf("tfvpngw-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfvpngw-subnet-%d", acctest.RandIntRange(10, 100))
	vniname := fmt.Sprintf("tfvpngw-createname-%d", acctest.RandIntRange(10, 100))
	tag1 := "env:test"
	tag2 := "env:dev"
	tag3 := "env:prod"
	enable_infrastructure_nat := true
	allow_ip_spoofing := true
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVirtualNetworkInterfaceDataSourceConfigBasic(vpcname, subnetname, vniname, tag1, tag2, tag3, enable_infrastructure_nat, allow_ip_spoofing, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMIsVirtualNetworkInterfaceExists("ibm_is_virtual_network_interface.testacc_vni", conf),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "virtual_network_interface"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "auto_delete"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "resource_group.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "security_groups.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "subnet.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "vpc.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "tags.#", "3"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "zone.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "enable_infrastructure_nat"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.is_virtual_network_interface", "allow_ip_spoofing"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVirtualNetworkInterfaceDataSourceConfigBasic(vpcname, subnetname, vniname, tag1, tag2, tag3 string, enablenat, allowipspoofing, isUpdate bool) string {
	return testAccCheckIBMIsVirtualNetworkInterfaceConfigBasic(vpcname, subnetname, vniname, tag1, tag2, tag3, enablenat, allowipspoofing, false) + fmt.Sprintf(`
		data "ibm_is_virtual_network_interface" "is_virtual_network_interface" {
			virtual_network_interface = ibm_is_virtual_network_interface.testacc_vni.id
		}
	`)
}

// TestAccIBMIsVirtualNetworkInterfaceDataSourcePublicAddressRanges verifies that
// public_address_ranges is populated on the single-VNI data source when an
// IPv6 PAR (via an authorized CIDR) targets the VNI.
func TestAccIBMIsVirtualNetworkInterfaceDataSourcePublicAddressRanges(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-vni-par-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tf-subnet-vni-par-%d", acctest.RandIntRange(10, 100))
	vniname := fmt.Sprintf("tf-vni-par-%d", acctest.RandIntRange(10, 100))
	authCIDRName := fmt.Sprintf("tf-authcidr-vni-%d", acctest.RandIntRange(10, 100))
	parName := fmt.Sprintf("tf-par-vni-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsVirtualNetworkInterfaceDataSourcePARConfig(vpcname, subnetname, vniname, authCIDRName, parName),
				Check: resource.ComposeTestCheckFunc(
					// VNI data source — public_address_ranges must be present
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "virtual_network_interface"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "lifecycle_state"),
					resource.TestCheckResourceAttr("data.ibm_is_virtual_network_interface.vni_ds", "public_address_ranges.#", "1"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "public_address_ranges.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "public_address_ranges.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "public_address_ranges.0.cidr"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "public_address_ranges.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_virtual_network_interface.vni_ds", "public_address_ranges.0.name"),
					resource.TestCheckResourceAttr("data.ibm_is_virtual_network_interface.vni_ds", "public_address_ranges.0.resource_type", "public_address_range"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVirtualNetworkInterfaceDataSourcePARConfig(vpcname, subnetname, vniname, authCIDRName, parName string) string {
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

data "ibm_is_virtual_network_interface" "vni_ds" {
  virtual_network_interface = ibm_is_virtual_network_interface.testacc_vni.id
  depends_on                = [ibm_is_public_address_range.testacc_par]
}
`, vpcname, subnetname, acc.ISZoneName, vniname, authCIDRName, acc.ISZoneName, parName)
}
