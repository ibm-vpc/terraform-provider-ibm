// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsVPNGatewayAdvertisedCidrBasic(t *testing.T) {
	var conf *string
	vpnGatewayID := fmt.Sprintf("tf_vpn_gateway_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsVPNGatewayAdvertisedCidrDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNGatewayAdvertisedCidrConfigBasic(vpnGatewayID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsVPNGatewayAdvertisedCidrExists("ibm_is_vpn_gateway_advertised_cidr.is_vpn_gateway_advertised_cidr_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_vpn_gateway_advertised_cidr.is_vpn_gateway_advertised_cidr_instance", "vpn_gateway", vpnGatewayID),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_vpn_gateway_advertised_cidr.is_vpn_gateway_advertised_cidr",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsVPNGatewayAdvertisedCidrConfigBasic(vpnGatewayID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_vpn_gateway_advertised_cidr" "is_vpn_gateway_advertised_cidr_instance" {
			vpn_gateway = "%s"
		}
	`, vpnGatewayID)
}

func testAccCheckIBMIsVPNGatewayAdvertisedCidrExists(n string, obj *string) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getVPNAdvertisedCidrOptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getVPNAdvertisedCidrOptions.SetVPNGatewayID(parts[0])
		getVPNAdvertisedCidrOptions.SetCIDR(parts[1])

		_, err = vpcClient.CheckVPNGatewayAdvertisedCIDR(getVPNAdvertisedCidrOptions)
		if err != nil {
			return err
		}
		return nil
	}
}

func testAccCheckIBMIsVPNGatewayAdvertisedCidrDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_vpn_gateway_advertised_cidr" {
			continue
		}

		getVPNAdvertisedCidrOptions := &vpcv1.CheckVPNGatewayAdvertisedCIDROptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getVPNAdvertisedCidrOptions.SetVPNGatewayID(parts[0])
		getVPNAdvertisedCidrOptions.SetCIDR(parts[1])

		// Try to find the key
		response, err := vpcClient.CheckVPNGatewayAdvertisedCIDR(getVPNAdvertisedCidrOptions)

		if err == nil {
			return fmt.Errorf("VPNAdvertisedCidr still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for VPNAdvertisedCidr (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
