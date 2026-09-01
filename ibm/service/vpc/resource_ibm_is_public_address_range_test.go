// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMPublicAddressRangeBasic(t *testing.T) {
	var conf vpcv1.PublicAddressRange
	ipv4AddressCount := "16"
	name := fmt.Sprintf("tf-name-par%d", acctest.RandIntRange(10, 100))
	vpcName := fmt.Sprintf("tf-name-vpc%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMPublicAddressRangeConfigBasic(vpcName, name, ipv4AddressCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeExists("ibm_is_public_address_range.public_address_range_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.public_address_range_instance", "ipv4_address_count", ipv4AddressCount),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "network_prefix_length"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "profile.#"),
				),
			},
		},
	})
}

func TestAccIBMPublicAddressRangeNameValidation(t *testing.T) {
	vpcName := fmt.Sprintf("tf-name-vpc%d", acctest.RandIntRange(10, 100))
	invalidName := "ibm-test"
	ipv4AddressCount := "16"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckIBMPublicAddressRangeConfigBasic(vpcName, invalidName, ipv4AddressCount),
				ExpectError: regexp.MustCompile(`"name" cannot start with 'ibm-'`),
			},
		},
	})
}

func testAccCheckIBMPublicAddressRangeConfigBasic(vpcName, name, ipv4AddressCount string) string {
	return fmt.Sprintf(`
		resource ibm_is_vpc testacc_vpc {
			name = "%s"
		}
		resource "ibm_is_public_address_range" "public_address_range_instance" {			
			name = "%s"
			ipv4_address_count = "%s"
			target {
    			vpc {
      				id = ibm_is_vpc.testacc_vpc.id
    			}
    			zone {
      				name = "%s"
    			}
  			}
		}
	`, vpcName, name, ipv4AddressCount, acc.ISZoneName)
}

func testAccCheckIBMPublicAddressRangeExists(n string, obj vpcv1.PublicAddressRange) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getPublicAddressRangeOptions := &vpcv1.GetPublicAddressRangeOptions{}

		getPublicAddressRangeOptions.SetID(rs.Primary.ID)

		publicAddressRange, _, err := vpcClient.GetPublicAddressRange(getPublicAddressRangeOptions)
		if err != nil {
			return err
		}

		obj = *publicAddressRange
		return nil
	}
}

func testAccCheckIBMPublicAddressRangeDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_public_address_range" {
			continue
		}

		getPublicAddressRangeOptions := &vpcv1.GetPublicAddressRangeOptions{}

		getPublicAddressRangeOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := vpcClient.GetPublicAddressRange(getPublicAddressRangeOptions)

		if err == nil {
			return fmt.Errorf("PublicAddressRange still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for PublicAddressRange (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

// ---------------------------------------------------------------------------
// IPv6 / authorized_cidr / VNI tests
// ---------------------------------------------------------------------------

func TestAccIBMPublicAddressRangeIPv6AuthorizedCIDR(t *testing.T) {
	var conf vpcv1.PublicAddressRange
	name := fmt.Sprintf("tf-par-ipv6-%d", acctest.RandIntRange(10, 100))
	authorizedCIDRName := fmt.Sprintf("tf-authcidr-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeIPv6Config(authorizedCIDRName, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeExists("ibm_is_public_address_range.par_ipv6", conf),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.par_ipv6", "name", name),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_ipv6", "network_prefix_length"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_ipv6", "ip_version"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_ipv6", "authorized_cidr.#"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_ipv6", "authorized_cidr.0.id"),
				),
			},
		},
	})
}

func TestAccIBMPublicAddressRangeComputedFields(t *testing.T) {
	var conf vpcv1.PublicAddressRange
	name := fmt.Sprintf("tf-par-comp-%d", acctest.RandIntRange(10, 100))
	vpcName := fmt.Sprintf("tf-vpc-comp-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeConfigBasic(vpcName, name, "16"),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeExists("ibm_is_public_address_range.public_address_range_instance", conf),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "network_prefix_length"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "ip_version"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "profile.#"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "profile.0.name"),
				),
			},
		},
	})
}

func TestAccIBMPublicAddressRangeVNITarget(t *testing.T) {
	var conf vpcv1.PublicAddressRange
	name := fmt.Sprintf("tf-par-vni-%d", acctest.RandIntRange(10, 100))
	vpcName := fmt.Sprintf("tf-vpc-vni-%d", acctest.RandIntRange(10, 100))
	subnetName := fmt.Sprintf("tf-subnet-vni-%d", acctest.RandIntRange(10, 100))
	vniName := fmt.Sprintf("tf-vni-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeVNIConfig(vpcName, subnetName, vniName, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeExists("ibm_is_public_address_range.par_vni", conf),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.par_vni", "name", name),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_vni", "target.#"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_vni", "target.0.virtual_network_interface.#"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_vni", "target.0.virtual_network_interface.0.id"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_vni", "target.0.virtual_network_interface.0.crn"),
				),
			},
			{
				// Update to a second VNI — exercises the VNI update path in resourceIBMPublicAddressRangeUpdate.
				Config: testAccCheckIBMPublicAddressRangeVNIUpdateConfig(vpcName, subnetName, vniName, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeExists("ibm_is_public_address_range.par_vni", conf),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.par_vni", "target.0.virtual_network_interface.0.id"),
				),
			},
		},
	})
}

// --- new config helpers ---

func testAccCheckIBMPublicAddressRangeIPv6Config(authorizedCIDRName, name string) string {
	return fmt.Sprintf(`
resource "ibm_is_public_address_range_authorized_cidr" "test_auth_cidr" {
  name                  = "%s"
  ip_version            = "ipv6"
  availability_mode     = "zonal"
  zone                  = "%s"
  network_prefix_length = 64
}

resource "ibm_is_public_address_range" "par_ipv6" {
  name                  = "%s"
  network_prefix_length = 64
  authorized_cidr {
    id = ibm_is_public_address_range_authorized_cidr.test_auth_cidr.id
  }
  target {
    zone {
      name = "%s"
    }
  }
}
`, authorizedCIDRName, acc.ISZoneName, name, acc.ISZoneName)
}

func testAccCheckIBMPublicAddressRangeVNIConfig(vpcName, subnetName, vniName, name string) string {
	return fmt.Sprintf(`
resource "ibm_is_vpc" "testacc_vpc" {
  name = "%s"
}
resource "ibm_is_subnet" "testacc_subnet" {
  name            = "%s"
  vpc             = ibm_is_vpc.testacc_vpc.id
  zone            = "%s"
  ipv4_cidr_block = "10.240.0.0/24"
}
resource "ibm_is_virtual_network_interface" "testacc_vni" {
  name   = "%s"
  subnet = ibm_is_subnet.testacc_subnet.id
}
resource "ibm_is_public_address_range" "par_vni" {
  name               = "%s"
  ipv4_address_count = 16
  target {
    virtual_network_interface {
      id = ibm_is_virtual_network_interface.testacc_vni.id
    }
  }
}
`, vpcName, subnetName, acc.ISZoneName, vniName, name)
}

func testAccCheckIBMPublicAddressRangeVNIUpdateConfig(vpcName, subnetName, vniName, name string) string {
	return fmt.Sprintf(`
resource "ibm_is_vpc" "testacc_vpc" {
  name = "%s"
}
resource "ibm_is_subnet" "testacc_subnet" {
  name            = "%s"
  vpc             = ibm_is_vpc.testacc_vpc.id
  zone            = "%s"
  ipv4_cidr_block = "10.240.0.0/24"
}
resource "ibm_is_virtual_network_interface" "testacc_vni" {
  name   = "%s"
  subnet = ibm_is_subnet.testacc_subnet.id
}
resource "ibm_is_virtual_network_interface" "testacc_vni2" {
  name   = "%s-2"
  subnet = ibm_is_subnet.testacc_subnet.id
}
resource "ibm_is_public_address_range" "par_vni" {
  name               = "%s"
  ipv4_address_count = 16
  target {
    virtual_network_interface {
      id = ibm_is_virtual_network_interface.testacc_vni2.id
    }
  }
}
`, vpcName, subnetName, acc.ISZoneName, vniName, vniName, name)
}
