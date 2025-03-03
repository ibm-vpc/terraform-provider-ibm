// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"

	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccIBMISVPNGatewayAdvertisedCidr_basic(t *testing.T) {
	var advertisedCidr string

	vpcname := fmt.Sprintf("tfsg-vpc-%d", acctest.RandIntRange(10, 100))
	subnetname := fmt.Sprintf("tfsg-subnet-%d", acctest.RandIntRange(10, 100))
	vpnname := fmt.Sprintf("tfvpngc-vpn-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMISVPNGatewayAdvertisedCidrDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISVPNGatewayAdvertisedCidrConfig(vpcname, subnetname, vpnname),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIBMISVPNGatewayAdvertisedCidrExists("ibm_is_vpn_gateway_advertised_cidr.VPNGatewayAdvertisedCidr", &advertisedCidr),
					// resource.TestCheckResourceAttrSet(
					// 	"ibm_is_vpn_gateway_advertised_cidr.VPNGatewayAdvertisedCidr", "advertised_cidrs"),
				),
			},
		},
	})
}

func testAccCheckIBMISVPNGatewayAdvertisedCidrDestroy(s *terraform.State) error {

	sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_vpn_gateway_advertised_cidr" {
			continue
		}
		parts, err := flex.IdParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		gID := parts[0]
		gCidr := parts[1]

		removeVPNGatewayAdvertisedCIDROptions := &vpcv1.RemoveVPNGatewayAdvertisedCIDROptions{
			VPNGatewayID: &gID,
			CIDR:         &gCidr,
		}
		response, err := sess.RemoveVPNGatewayAdvertisedCIDR(removeVPNGatewayAdvertisedCIDROptions)
		if err == nil {
			return fmt.Errorf("Advertised Cidr still exists: %v", response)
		}
	}
	return nil
}

func testAccCheckIBMISVPNGatewayAdvertisedCidrExists(n string, advertisedCidr *string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("[ERROR] No Advertised Cidr is set")
		}

		sess, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		parts, err := flex.IdParts(rs.Primary.ID)
		if err != nil {
			return err
		}
		gID := parts[0]
		gAdvertisedCidr := parts[1]

		checkVPNGatewayAdvertisedCIDROptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{
			VPNGatewayID: &gID,
			CIDR:         &gAdvertisedCidr,
		}

		response, err := sess.CheckVPNGatewayAdvertisedCIDR(checkVPNGatewayAdvertisedCIDROptions)
		if err != nil {
			if response != nil && response.StatusCode == 404 {
				*advertisedCidr = ""
				return nil
			}
			return fmt.Errorf("[ERROR] Error getting Advertised Cidr : %s\n%s", err, response)
		}
		*advertisedCidr = fmt.Sprintf("%s/%s", gID, gAdvertisedCidr)
		return nil
	}
}

func testAccCheckIBMISVPNGatewayAdvertisedCidrConfig(vpcname, subnetname, vpnname string) string {
	return fmt.Sprintf(`
		resource "ibm_is_vpc" "testacc_vpc1" {
		name = "%s"
		}

		resource "ibm_is_subnet" "testacc_subnet1" {
		name            = "%s"
		vpc             = ibm_is_vpc.testacc_vpc1.id
		zone            = "%s"
		ipv4_cidr_block = "%s"
		}
		resource "ibm_is_vpn_gateway" "testacc_VPNGateway1" {
		name   = "%s"
		subnet = ibm_is_subnet.testacc_subnet1.id
		timeouts {
			create = "18m"
			delete = "18m"
		}
		}
		resource "ibm_is_vpn_gateway_advertised_cidr" "VPNGatewayAdvertisedCidr" {
		vpn_gateway = ibm_is_vpn_gateway.testacc_VPNGateway1.id
		cidr        = "10.45.0.0/24"
	}`, vpcname, subnetname, acc.ISZoneName, acc.ISCIDR, vpnname)

}
