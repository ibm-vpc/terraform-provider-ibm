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

func TestAccIBMPublicAddressRangeLBTarget(t *testing.T) {
	var conf vpcv1.PublicAddressRange
	ipv4AddressCount := "16"
	name := fmt.Sprintf("tf-name-par%d", acctest.RandIntRange(10, 100))
	vpcName := fmt.Sprintf("tf-name-vpc%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeConfigBasic(vpcName, name, ipv4AddressCount),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeExists("ibm_is_public_address_range.public_address_range_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.public_address_range_instance", "ipv4_address_count", ipv4AddressCount),
					// target.load_balancer is read-only and should be absent when bound to a vpc/zone
					resource.TestCheckResourceAttr("ibm_is_public_address_range.public_address_range_instance", "target.0.load_balancer.#", "0"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "target.0.vpc.0.id"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "target.0.zone.0.name"),
				),
			},
		},
	})
}

// testAccCheckIBMPublicAddressRangeConfigWithLBTarget creates an IPv6-enabled LB and a PAR
// whose target is that LB. The PAR has ip_version ipv6 and a /128 prefix (required by the API).
func testAccCheckIBMPublicAddressRangeConfigWithLBTarget(vpcName, subnetName, zone, cidr, lbName, parName string) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}
	resource "ibm_is_subnet" "testacc_subnet" {
		name            = "%s"
		vpc             = ibm_is_vpc.testacc_vpc.id
		zone            = "%s"
		ipv4_cidr_block = "%s"
	}
	resource "ibm_is_subnet_reserved_ip" "testacc_rip" {
		subnet = ibm_is_subnet.testacc_subnet.id
	}
	resource "ibm_is_lb" "testacc_lb" {
		name         = "%s"
		type         = "public"
		profile      = "dynamic"
		ipv6_enabled = true
		address_mode = "static"
		subnets      = [ibm_is_subnet.testacc_subnet.id]
		private_ips  = [ibm_is_subnet_reserved_ip.testacc_rip.reserved_ip]
	}
	resource "ibm_is_public_address_range" "public_address_range_instance" {
		name               = "%s"
		ipv4_address_count = 1
		target {
			load_balancer {
				id = ibm_is_lb.testacc_lb.id
			}
		}
	}`, vpcName, subnetName, zone, cidr, lbName, parName)
}

func TestAccIBMPublicAddressRangeLBTargetBound(t *testing.T) {
	var conf vpcv1.PublicAddressRange
	vpcName := fmt.Sprintf("tf-name-vpc%d", acctest.RandIntRange(10, 100))
	subnetName := fmt.Sprintf("tf-name-subnet%d", acctest.RandIntRange(10, 100))
	lbName := fmt.Sprintf("tf-name-lb%d", acctest.RandIntRange(10, 100))
	parName := fmt.Sprintf("tf-name-par%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMPublicAddressRangeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMPublicAddressRangeConfigWithLBTarget(vpcName, subnetName, acc.ISZoneName, acc.ISCIDR, lbName, parName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMPublicAddressRangeExists("ibm_is_public_address_range.public_address_range_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.public_address_range_instance", "name", parName),
					// load_balancer block is populated after binding
					resource.TestCheckResourceAttr("ibm_is_public_address_range.public_address_range_instance", "target.0.load_balancer.#", "1"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "target.0.load_balancer.0.id"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "target.0.load_balancer.0.crn"),
					resource.TestCheckResourceAttrSet("ibm_is_public_address_range.public_address_range_instance", "target.0.load_balancer.0.name"),
					resource.TestCheckResourceAttr("ibm_is_public_address_range.public_address_range_instance", "target.0.load_balancer.0.resource_type", "load_balancer"),
					// vpc/zone are absent when targeting an LB
					resource.TestCheckResourceAttr("ibm_is_public_address_range.public_address_range_instance", "target.0.vpc.#", "0"),
				),
			},
		},
	})
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
