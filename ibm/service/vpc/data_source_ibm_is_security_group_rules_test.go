// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsSecurityGroupRulesDataSourceBasic(t *testing.T) {
	vpcname := fmt.Sprintf("tfsubnet-vpc-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tfsubnet-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsSecurityGroupRulesDataSourceConfigBasic(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "security_group"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "rules.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "rules.0.direction"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "rules.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "rules.0.ip_version"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "rules.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "rules.0.protocol"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.example", "rules.0.remote.#"),
				),
			},
		},
	})
}

func TestAccIBMIsSecurityGroupRulesDataSource_IPv6Icmp(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-ipv6icmp-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tf-sg-ipv6icmp-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsSecurityGroupRulesDataSourceConfigIPv6Icmp(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.ipv6icmp", "security_group"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.ipv6icmp", "rules.#"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rules.ipv6icmp", "rules.0.protocol", "ipv6_icmp"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rules.ipv6icmp", "rules.0.ip_version", "ipv6"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rules.ipv6icmp", "rules.0.direction", "inbound"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.ipv6icmp", "rules.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.ipv6icmp", "rules.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.ipv6icmp", "rules.0.code"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.ipv6icmp", "rules.0.type"),
				),
			},
		},
	})
}

func TestAccIBMIsSecurityGroupRulesDataSource_IndividualIPv6(t *testing.T) {
	vpcname := fmt.Sprintf("tf-vpc-indv6-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tf-sg-indv6-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsSecurityGroupRulesDataSourceConfigIndividualIPv6(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.indv6", "security_group"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.indv6", "rules.#"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rules.indv6", "rules.0.protocol", "ipv6_hop_opt"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rules.indv6", "rules.0.ip_version", "ipv6"),
					resource.TestCheckResourceAttr("data.ibm_is_security_group_rules.indv6", "rules.0.direction", "outbound"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.indv6", "rules.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_security_group_rules.indv6", "rules.0.id"),
				),
			},
		},
	})
}

func testAccCheckIBMIsSecurityGroupRulesDataSourceConfigBasic(vpcname, sgname string) string {
	return fmt.Sprintf(`

	resource "ibm_is_vpc" "example" {
		name = "%s"
	  }
	  
	  resource "ibm_is_security_group" "example" {
		name = "%s"
		vpc  = ibm_is_vpc.example.id
		depends_on = [
			ibm_is_vpc.example,
		]
	  }
	  
	  resource "ibm_is_security_group_rule" "example" {
		group     = ibm_is_security_group.example.id
		direction = "outbound"
		remote    = "127.0.0.1"
		tcp {
		  port_min = 8080
		  port_max = 8080
		}
		depends_on = [
			ibm_is_security_group.example,
		]
	  }
		data "ibm_is_security_group_rules" "example" {
			depends_on = [
				ibm_is_security_group_rule.example,
		]
			security_group = ibm_is_security_group.example.id
		}
	`, vpcname, sgname)
}

func testAccCheckIBMIsSecurityGroupRulesDataSourceConfigIPv6Icmp(vpcname, sgname string) string {
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
		depends_on = [ibm_is_security_group.ipv6icmp]
	}

	data "ibm_is_security_group_rules" "ipv6icmp" {
		security_group = ibm_is_security_group.ipv6icmp.id
		depends_on     = [ibm_is_security_group_rule.ipv6icmp]
	}
	`, vpcname, sgname)
}

func testAccCheckIBMIsSecurityGroupRulesDataSourceConfigIndividualIPv6(vpcname, sgname string) string {
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
		ip_version = "ipv6"
		protocol   = "ipv6_hop_opt"
		depends_on = [ibm_is_security_group.indv6]
	}

	data "ibm_is_security_group_rules" "indv6" {
		security_group = ibm_is_security_group.indv6.id
		depends_on     = [ibm_is_security_group_rule.indv6]
	}
	`, vpcname, sgname)
}
