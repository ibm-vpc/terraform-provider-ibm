// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMISSecurityGroupDatasource_basic(t *testing.T) {
	vpcname := fmt.Sprintf("tfsubnet-vpc-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tfsubnet-name-%d", acctest.RandIntRange(10, 100))
	dataSourceName := "data.ibm_is_security_group.sg1_rule"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISSgRuleConfig(vpcname, sgname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
				),
			},
		},
	})
}
func TestAccIBMISSecurityGroupDatasource_filters(t *testing.T) {
	vpcname := fmt.Sprintf("tfsubnet-vpc-%d", acctest.RandIntRange(10, 100))
	sgname := fmt.Sprintf("tfsubnet-name-%d", acctest.RandIntRange(10, 100))
	vpcname1 := fmt.Sprintf("tfsubnet-vpc1-%d", acctest.RandIntRange(10, 100))
	sgname1 := fmt.Sprintf("tfsubnet-name1-%d", acctest.RandIntRange(10, 100))
	dataSourceName := "data.ibm_is_security_group.sg1_rule"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISSecurityGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISSgRuleFiltersConfig(vpcname, sgname, vpcname1, sgname1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(dataSourceName, "vpc"),
					resource.TestCheckResourceAttrSet(dataSourceName, "rules.#"),
					resource.TestCheckResourceAttrSet(dataSourceName, "tags.#"),
				),
			},
		},
	})
}

func testAccCheckIBMISSgRuleConfig(vpcname, sgname string) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	  }
	  
	  resource "ibm_is_security_group" "testacc_security_group" {
		name = "%s"
		tags = ["sgtag1" , "sgTag2"]
		vpc  = ibm_is_vpc.testacc_vpc.id
	  }
	  
	  resource "ibm_is_security_group_rule" "testacc_security_group_rule_all" {
		group     = ibm_is_security_group.testacc_security_group.id
		direction = "inbound"
		remote    = "127.0.0.1"
	  }

	  resource "ibm_is_security_group_rule" "testacc_security_group_rule_icmp" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_all]
        group      = ibm_is_security_group.testacc_security_group.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        icmp {
          code = 20
          type = 30
        }
      }
      
      resource "ibm_is_security_group_rule" "testacc_security_group_rule_udp" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_icmp]
        group      = ibm_is_security_group.testacc_security_group.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        udp {
          port_min = 805
          port_max = 807
        }
      }
      
      resource "ibm_is_security_group_rule" "testacc_security_group_rule_tcp" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_udp]
        group      = ibm_is_security_group.testacc_security_group.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        tcp {
          port_min = 8080
          port_max = 8080
        }
	  }
	  
	  data "ibm_is_security_group" "sg1_rule" {
		name = ibm_is_security_group.testacc_security_group.name
	}

	  
    `, vpcname, sgname)

}
func testAccCheckIBMISSgRuleFiltersConfig(vpcname, sgname, vpcname1, sgname1 string) string {
	return fmt.Sprintf(`
	resource "ibm_is_vpc" "testacc_vpc" {
		name = "%s"
	}
	resource "ibm_is_vpc" "testacc_vpc1" {
		name = "%s"
	}
	  
	resource "ibm_is_security_group" "testacc_security_group" {
		name = "%s"
		tags = ["sgtag1" , "sgTag2"]
		vpc  = ibm_is_vpc.testacc_vpc.id
	}
	resource "ibm_is_security_group" "testacc_security_group1" {
		name = "%s"
		tags = ["sgtag1" , "sgTag2"]
		vpc  = ibm_is_vpc.testacc_vpc1.id
	}
	  
	resource "ibm_is_security_group_rule" "testacc_security_group_rule_all" {
		group     = ibm_is_security_group.testacc_security_group.id
		direction = "inbound"
		remote    = "127.0.0.1"
	}

	resource "ibm_is_security_group_rule" "testacc_security_group_rule_icmp" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_all]
        group      = ibm_is_security_group.testacc_security_group.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        icmp {
          code = 20
          type = 30
        }
    }
      
    resource "ibm_is_security_group_rule" "testacc_security_group_rule_udp" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_icmp]
        group      = ibm_is_security_group.testacc_security_group.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        udp {
          port_min = 805
          port_max = 807
        }
    }
      
    resource "ibm_is_security_group_rule" "testacc_security_group_rule_tcp" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_udp]
        group      = ibm_is_security_group.testacc_security_group.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        tcp {
          port_min = 8080
          port_max = 8080
        }
	}
	resource "ibm_is_security_group_rule" "testacc_security_group_rule_all1" {
		group     = ibm_is_security_group.testacc_security_group1.id
		direction = "inbound"
		remote    = "127.0.0.1"
	}

	resource "ibm_is_security_group_rule" "testacc_security_group_rule_icmp1" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_all1]
        group      = ibm_is_security_group.testacc_security_group1.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        icmp {
          code = 20
          type = 30
        }
    }
      
    resource "ibm_is_security_group_rule" "testacc_security_group_rule_udp1" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_icmp1]
        group      = ibm_is_security_group.testacc_security_group1.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        udp {
          port_min = 805
          port_max = 807
        }
    }
      
    resource "ibm_is_security_group_rule" "testacc_security_group_rule_tcp1" {
        depends_on = [ibm_is_security_group_rule.testacc_security_group_rule_udp1]
        group      = ibm_is_security_group.testacc_security_group1.id
        direction  = "inbound"
        remote     = "127.0.0.1"
        tcp {
          port_min = 8080
          port_max = 8080
        }
	}
	  
	data "ibm_is_security_group" "sg1_rule" {
		name = ibm_is_security_group.testacc_security_group.name
		vpc_id = ibm_is_vpc.testacc_vpc.id
	}
	data "ibm_is_security_group" "sg1_rule1" {
		name = ibm_is_security_group.testacc_security_group1.name
		vpc_crn = ibm_is_vpc.testacc_vpc1.crn
	}

	  
    `, vpcname, vpcname1, sgname, sgname1)

}
