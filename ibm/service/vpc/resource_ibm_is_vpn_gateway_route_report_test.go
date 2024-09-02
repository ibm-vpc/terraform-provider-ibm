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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsVPNGatewayRouteReportBasic(t *testing.T) {
	var conf vpcv1.VPNRouteReport
	vpnGatewayID := fmt.Sprintf("tf_vpn_gateway_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsVPNGatewayRouteReportDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNGatewayRouteReportConfigBasic(vpnGatewayID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsVPNGatewayRouteReportExists("ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "vpn_gateway", vpnGatewayID),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsVPNGatewayRouteReportConfigBasic(vpnGatewayID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_vpn_gateway_route_report" "is_vpn_gateway_route_report_instance" {
			vpn_gateway = "%s"
		}
	`, vpnGatewayID)
}

func testAccCheckIBMIsVPNGatewayRouteReportExists(n string, obj vpcv1.VPNRouteReport) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getVPNRouteReportOptions := &vpcv1.GetVPNRouteReportOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getVPNRouteReportOptions.SetVPNGatewayID(parts[0])
		getVPNRouteReportOptions.SetID(parts[1])

		vpnRouteReport, _, err := vpcClient.GetVPNRouteReport(getVPNRouteReportOptions)
		if err != nil {
			return err
		}

		obj = *vpnRouteReport
		return nil
	}
}

func testAccCheckIBMIsVPNGatewayRouteReportDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_vpn_gateway_route_report" {
			continue
		}

		getVPNRouteReportOptions := &vpcv1.GetVPNRouteReportOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getVPNRouteReportOptions.SetVPNGatewayID(parts[0])
		getVPNRouteReportOptions.SetID(parts[1])

		// Try to find the key
		_, response, err := vpcClient.GetVPNRouteReport(getVPNRouteReportOptions)

		if err == nil {
			return fmt.Errorf("VPNRouteReport still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for VPNRouteReport (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}

func TestResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		vpnRouteReportRouteNextHopModel := make(map[string]interface{})
		vpnRouteReportRouteNextHopModel["address"] = "192.168.1.1"
		vpnRouteReportRouteNextHopModel["used"] = true

		model := make(map[string]interface{})
		model["as_path"] = []int64{4201065558, 4203065547}
		model["best_path"] = true
		model["next_hops"] = []map[string]interface{}{vpnRouteReportRouteNextHopModel}
		model["peer"] = "192.168.3.4"
		model["prefix"] = "172.16.0.0/16"
		model["valid"] = true
		model["weight"] = int(0)

		assert.Equal(t, result, model)
	}

	vpnRouteReportRouteNextHopModel := new(vpcv1.VPNRouteReportRouteNextHop)
	vpnRouteReportRouteNextHopModel.Address = core.StringPtr("192.168.1.1")
	vpnRouteReportRouteNextHopModel.Used = core.BoolPtr(true)

	model := new(vpcv1.VPNRouteReportRoute)
	model.AsPath = []int64{4201065558, 4203065547}
	model.BestPath = core.BoolPtr(true)
	model.NextHops = []vpcv1.VPNRouteReportRouteNextHop{*vpnRouteReportRouteNextHopModel}
	model.Peer = core.StringPtr("192.168.3.4")
	model.Prefix = core.StringPtr("172.16.0.0/16")
	model.Valid = core.BoolPtr(true)
	model.Weight = core.Int64Ptr(int64(0))

	result, err := vpc.ResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["address"] = "192.168.3.4"
		model["used"] = true

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPNRouteReportRouteNextHop)
	model.Address = core.StringPtr("192.168.3.4")
	model.Used = core.BoolPtr(true)

	result, err := vpc.ResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
