// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
 */

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsVPNGatewayRouteReportDataSourceBasic(t *testing.T) {
	vpnRouteReportVPNGatewayID := fmt.Sprintf("tf_vpn_gateway_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNGatewayRouteReportDataSourceConfigBasic(vpnRouteReportVPNGatewayID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "is_vpn_gateway_route_report_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "routes.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance", "updated_at"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNGatewayRouteReportDataSourceConfigBasic(vpnRouteReportVPNGatewayID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_vpn_gateway_route_report" "is_vpn_gateway_route_report_instance" {
			vpn_gateway_id = "%s"
		}

		data "ibm_is_vpn_gateway_route_report" "is_vpn_gateway_route_report_instance" {
			vpn_gateway_id = ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance.vpn_gateway_id
			is_vpn_gateway_route_report_id = ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance.is_vpn_gateway_route_report_id
		}
	`, vpnRouteReportVPNGatewayID)
}

func TestDataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["address"] = "192.168.3.4"
		model["used"] = true

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPNRouteReportRouteNextHop)
	model.Address = core.StringPtr("192.168.3.4")
	model.Used = core.BoolPtr(true)

	result, err := vpc.DataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
