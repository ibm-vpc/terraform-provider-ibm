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

func TestAccIBMISLBDatasource_basic(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	routeMode := "false"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "route_mode", routeMode),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "availability"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "instance_groups_supported"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "source_ip_session_persistence_supported"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "asymmetric_routing_supported"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "mtls_supported"),
				),
			},
		},
	})
}
func TestAccIBMISLBDatasource_failsafe_policy(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	routeMode := "false"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "route_mode", routeMode),
				),
			},
		},
	})
}
func TestAccIBMISLBDatasource_address_mode(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	floatingipname := fmt.Sprintf("tflb-fip-%d", acctest.RandIntRange(10, 100))
	floatingipname2 := fmt.Sprintf("tflb-fip2-%d", acctest.RandIntRange(10, 100))
	floatingipname3 := fmt.Sprintf("tflb-fip3-%d", acctest.RandIntRange(10, 100))
	reservedipname := fmt.Sprintf("tflb-rip-%d", acctest.RandIntRange(10, 100))
	reservedipname2 := fmt.Sprintf("tflb-rip2-%d", acctest.RandIntRange(10, 100))
	reservedipname3 := fmt.Sprintf("tflb-rip3-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBAddressModeConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name, floatingipname, floatingipname2, floatingipname3, reservedipname, reservedipname2, reservedipname3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "address_mode"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "public_ips.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "public_ip.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "public_ip.0.resource_type"),
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "public_range.#", "0"),
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "public_ranges.#", "0"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "private_ips.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "private_ip.#"),
				),
			},
		},
	})
}

func TestAccIBMISLBDatasource_ReservedIp(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	routeMode := "false"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "route_mode", routeMode),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_lb.ds_lb", "private_ip.#"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_lb.ds_lb", "private_ip.0.address"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_lb.ds_lb", "private_ip.0.href"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_lb.ds_lb", "private_ip.0.name"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_lb.ds_lb", "private_ip.0.reserved_ip"),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_lb.ds_lb", "private_ip.0.resource_type"),
				),
			},
		},
	})
}

func TestAccIBMISLBDatasource_pDNS(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	var lb string

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{

			{
				Config: testAccCheckIBMISDSLBDNS(vpcname, subnetname, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISLBExists("ibm_is_lb.testacc_LB", lb),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.test_lb", "name", name),
					resource.TestCheckResourceAttrSet(
						"data.ibm_is_lb.test_lb", "dns.#"),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.test_lb", "dns.0.instance_crn", acc.DNSInstanceCRN),
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.test_lb", "dns.0.zone_id", acc.DNSZoneID),
				),
			},
		},
	})

}

func TestAccIBMISLBDatasource_mTLS(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "mtls_supported"),
				),
			},
		},
	})
}

func testDSCheckIBMISLBAddressModeConfig(vpcname, subnetname, zone, cidr, name, fipname, fipname2, fipname3, ripname, ripname2, ripname3 string) string {
	return testAccCheckIBMISLBStaticAddressModeConfig(vpcname, subnetname, zone, cidr, name, fipname, fipname2, fipname3, ripname, ripname2, ripname3) + `
data "ibm_is_lb" "ds_lb" {
	name       = ibm_is_lb.testacc_LB.name
	depends_on = [ibm_is_lb.testacc_LB]
}`
}

func testDSCheckIBMISLBConfig(vpcname, subnetname, zone, cidr, name string) string {
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
}`, vpcname, subnetname, zone, cidr, name)
}

func testDSCheckIBMISLBIPv6Config(vpcname, subnetname, zone, cidr, name string) string {
	return testAccCheckIBMISLBIPv6Config(vpcname, subnetname, zone, cidr, name, true) + fmt.Sprintf(`
data "ibm_is_lb" "ds_lb" {
  name = ibm_is_lb.testacc_LB.name
}`)
}

func TestAccIBMISLBDatasource_ipv6(t *testing.T) {
	name := fmt.Sprintf("tflb-ipv6-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBIPv6Config(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "ipv6_enabled", "true"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "public_ranges.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "public_range.#"),
				),
			},
		},
	})
}

func TestAccIBMISLBDatasource_ipv6_with_par(t *testing.T) {
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	parName := fmt.Sprintf("tflb-par-%d", acctest.RandIntRange(10, 100))
	lbName := fmt.Sprintf("tflb-ipv6-par-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISLBIPv6WithPARConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, parName, lbName) + `
	data "ibm_is_lb" "ds_lb" {
		name = ibm_is_lb.testacc_LB.name
	}
	data "ibm_is_lbs" "ds_lbs" {
		depends_on = [ibm_is_lb.testacc_LB]
	}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "ipv6_enabled", "true"),
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "public_ranges.#", "1"),
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "public_range.#", "1"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "public_range.0.cidr"),
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "public_range.0.resource_type", "public_address_range"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "public_ip.0.cidr"),
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "public_ip.0.resource_type", "public_address_range"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lbs.ds_lbs", "load_balancers.#"),
				),
			},
		},
	})
}

func testAccCheckIBMISDSLBDNS(vpcname, subnetname, name string) string {
	// status filter defaults to empty
	return testAccCheckIBMISLBDNS(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name, acc.DNSInstanceCRN, acc.DNSZoneID) + fmt.Sprintf(`
      data "ibm_is_lb" "test_lb" {
		name = ibm_is_lb.testacc_LB.name
	  }`)
}

func TestAccIBMISLBDatasource_http_bundle(t *testing.T) {
	name := fmt.Sprintf("tflb-name-%d", acctest.RandIntRange(10, 100))
	vpcname := fmt.Sprintf("tflb-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tflb-subnet-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testDSCheckIBMISLBConfig(vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.ibm_is_lb.ds_lb", "name", name),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "advanced_health_checks_supported"),
					resource.TestCheckResourceAttrSet("data.ibm_is_lb.ds_lb", "fqdn_pool_members_supported"),
				),
			},
		},
	})
}
