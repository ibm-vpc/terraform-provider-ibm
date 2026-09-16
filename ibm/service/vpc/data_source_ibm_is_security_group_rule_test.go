// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsSecurityGroupRuleDataSourceBasic(t *testing.T) {
	vpcname := fmt.Sprintf("tfsecgrprl-vpc-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tfsecgrprl-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsSecurityGroupRuleDataSourceConfigBasic(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.example", "security_group"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.example", "direction"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.example", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.example", "ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.example", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.example", "protocol"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.example", "remote.#"),
				),
			},
		},
	})
}

// TestAccIBMIsSecurityGroupRuleDataSource_IPv6Icmp verifies that the single-rule
// data source correctly reads back an ipv6_icmp rule with type and code.
func TestAccIBMIsSecurityGroupRuleDataSource_IPv6Icmp(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-sgr-ipv6icmp-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tf-sg-sgr-ipv6icmp-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsSecurityGroupRuleDataSourceIPv6Icmp(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.ipv6icmp", "id"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6icmp", "ip_version", "ipv6"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6icmp", "protocol", "ipv6_icmp"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6icmp", "direction", "inbound"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.ipv6icmp", "type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.ipv6icmp", "code"),
				),
			},
		},
	})
}

// TestAccIBMIsSecurityGroupRuleDataSource_IPv6Individual verifies the single-rule
// data source for an IPv6-specific individual protocol (ipv6_hop_opt).
func TestAccIBMIsSecurityGroupRuleDataSource_IPv6Individual(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-sgr-indv6-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tf-sg-sgr-indv6-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsSecurityGroupRuleDataSourceIPv6Individual(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.indv6", "id"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.indv6", "ip_version", "ipv6"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.indv6", "protocol", "ipv6_hop_opt"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.indv6", "direction", "outbound"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.indv6", "href"),
				),
			},
		},
	})
}

// TestAccIBMIsSecurityGroupRuleDataSource_IPv6TCP verifies the single-rule
// data source for a tcp rule with ip_version = "ipv6".
func TestAccIBMIsSecurityGroupRuleDataSource_IPv6TCP(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-sgr-ipv6tcp-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tf-sg-sgr-ipv6tcp-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsSecurityGroupRuleDataSourceIPv6TCP(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rule.ipv6tcp", "id"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6tcp", "ip_version", "ipv6"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6tcp", "protocol", "tcp"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6tcp", "direction", "inbound"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6tcp", "port_min", "443"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rule.ipv6tcp", "port_max", "443"),
				),
			},
		},
	})
}

func testAccCheckIBMIsSecurityGroupRuleDataSourceConfigBasic(vpcname, sgname string) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "example" {
		name = "%s"
	}
	resource "ibm_is_security_group" "example" {
		name = "%s"
		vpc  = ibm_is_vpc.example.id
	}
	resource "ibm_is_security_group_rule" "example" {
		group     = ibm_is_security_group.example.id
		direction = "inbound"
		remote    = "0.0.0.0/0"
		protocol  = "tcp"
		port_min  = 80
		port_max  = 80
	}
	data "ibm_is_security_group_rule" "example" {
		security_group      = ibm_is_security_group.example.id
		security_group_rule = ibm_is_security_group_rule.example.rule_id
	}
	`, vpcname, sgname)
}

func testAccCheckIBMIsSecurityGroupRuleDataSourceIPv6Icmp(vpcname, sgname string) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "ipv6icmp" {
		name = "%s"
	}
	resource "ibm_is_security_group" "ipv6icmp" {
		name = "%s"
		vpc  = ibm_is_vpc.ipv6icmp.id
	}
	resource "ibm_is_security_group_rule" "ipv6icmp" {
		group      = ibm_is_security_group.ipv6icmp.id
		direction  = "inbound"
		remote     = "::/0"
		ip_version = "ipv6"
		protocol   = "ipv6_icmp"
		type       = 128
		code       = 0
	}
	data "ibm_is_security_group_rule" "ipv6icmp" {
		security_group      = ibm_is_security_group.ipv6icmp.id
		security_group_rule = ibm_is_security_group_rule.ipv6icmp.rule_id
	}
	`, vpcname, sgname)
}

func testAccCheckIBMIsSecurityGroupRuleDataSourceIPv6Individual(vpcname, sgname string) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "indv6" {
		name = "%s"
	}
	resource "ibm_is_security_group" "indv6" {
		name = "%s"
		vpc  = ibm_is_vpc.indv6.id
	}
	resource "ibm_is_security_group_rule" "indv6" {
		group      = ibm_is_security_group.indv6.id
		direction  = "outbound"
		remote     = "::/0"
		ip_version = "ipv6"
		protocol   = "ipv6_hop_opt"
	}
	data "ibm_is_security_group_rule" "indv6" {
		security_group      = ibm_is_security_group.indv6.id
		security_group_rule = ibm_is_security_group_rule.indv6.rule_id
	}
	`, vpcname, sgname)
}

func testAccCheckIBMIsSecurityGroupRuleDataSourceIPv6TCP(vpcname, sgname string) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "ipv6tcp" {
		name = "%s"
	}
	resource "ibm_is_security_group" "ipv6tcp" {
		name = "%s"
		vpc  = ibm_is_vpc.ipv6tcp.id
	}
	resource "ibm_is_security_group_rule" "ipv6tcp" {
		group      = ibm_is_security_group.ipv6tcp.id
		direction  = "inbound"
		remote     = "::/0"
		ip_version = "ipv6"
		protocol   = "tcp"
		port_min   = 443
		port_max   = 443
	}
	data "ibm_is_security_group_rule" "ipv6tcp" {
		security_group      = ibm_is_security_group.ipv6tcp.id
		security_group_rule = ibm_is_security_group_rule.ipv6tcp.rule_id
	}
	`, vpcname, sgname)
}
