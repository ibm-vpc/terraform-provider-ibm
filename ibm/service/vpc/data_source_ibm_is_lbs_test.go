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

func TestAccIBMISLBSDatasource_basic(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	var lb string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{

			{
				Config: testDSCheckIBMISLBSConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISLBExists("ibm_is_lb.testacc_lb", lb),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.availability"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.instance_groups_supported"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.source_ip_session_persistence_supported"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.asymmetric_routing_supported"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.mtls_supported"),
				),
			},
			{
				Config: testDSCheckIBMISLBSDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.route_mode"),
				),
			},
		},
	})

}
func TestAccIBMISLBSDatasource_address_mode(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	floatingipname := fmt.Sprintf("tflb-fip-%d", acctest.RandIntRange(10, 100))
	floatingipname2 := fmt.Sprintf("tflb-fip2-%d", acctest.RandIntRange(10, 100))
	floatingipname3 := fmt.Sprintf("tflb-fip3-%d", acctest.RandIntRange(10, 100))
	reservedipname := fmt.Sprintf("tflb-rip-%d", acctest.RandIntRange(10, 100))
	reservedipname2 := fmt.Sprintf("tflb-rip2-%d", acctest.RandIntRange(10, 100))
	reservedipname3 := fmt.Sprintf("tflb-rip3-%d", acctest.RandIntRange(10, 100))
	var lb string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBSAddressModeConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name, floatingipname, floatingipname2, floatingipname3, reservedipname, reservedipname2, reservedipname3),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISLBExists("ibm_is_lb.testacc_LB", lb),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.address_mode"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.public_ips.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.public_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ips.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ip.#"),
				),
			},
		},
	})
}

func TestAccIBMISLBSDatasource_ReservedIp(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	var lb string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{

			{
				Config: testDSCheckIBMISLBSConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISLBExists("ibm_is_lb.testacc_lb", lb),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
				),
			},
			{
				Config: testDSCheckIBMISLBSDatasourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.route_mode"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ip.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ip.0.address"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ip.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ip.0.reserved_ip"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.private_ip.0.resource_type"),
				),
			},
		},
	})

}

func TestAccIBMISLBSDatasource_mTLS(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	var lb string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBSConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISLBExists("ibm_is_lb.testacc_lb", lb),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.test_lbs", "load_balancers.0.mtls_supported"),
				),
			},
		},
	})

}

func testDSCheckIBMISLBSAddressModeConfig(vpcname, subnetname, zone, cidr, name, fipname, fipname2, fipname3, ripname, ripname2, ripname3 string) string {
	return testAccCheckIBMISLBStaticAddressModeConfig(vpcname, subnetname, zone, cidr, name, fipname, fipname2, fipname3, ripname, ripname2, ripname3) + `
data "ibm_is_lbs" "test_lbs" {
	depends_on = [ibm_is_lb.testacc_LB]
}`
}

func testDSCheckIBMISLBSConfig(vpcname, subnetname, zone, cidr, name string) string {
	// status filter defaults to empty
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
	  resource "ibm_is_lb" "testacc_lb" {
		name    = "%s"
		subnets = [ibm_is_subnet.testacc_subnet.id]
	  }
	  data "ibm_is_lb" "ds_lb" {
		name = ibm_is_lb.testacc_lb.name
	  }
	  data "ibm_is_lbs" "test_lbs" {
	  	depends_on = [ibm_is_lb.testacc_lb]
	  }	
	  `, vpcname, subnetname, zone, cidr, name)
}
func testDSCheckIBMISLBSDatasourceConfig() string {
	// status filter defaults to empty
	return fmt.Sprintf(`
      data "ibm_is_lbs" "test_lbs" {
	  }`)
}
