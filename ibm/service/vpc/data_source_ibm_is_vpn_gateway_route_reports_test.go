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
	. "github.com/IBM-Cloud/terraform-provider-ibm/ibm/unittest"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsVPNGatewayRouteReportsDataSourceBasic(t *testing.T) {
	vpnRouteReportVPNGatewayID := fmt.Sprintf("tf_vpn_gateway_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsVPNGatewayRouteReportsDataSourceConfigBasic(vpnRouteReportVPNGatewayID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_reports.is_vpn_gateway_route_reports_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_reports.is_vpn_gateway_route_reports_instance", "vpn_gateway_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_vpn_gateway_route_reports.is_vpn_gateway_route_reports_instance", "route_reports.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsVPNGatewayRouteReportsDataSourceConfigBasic(vpnRouteReportVPNGatewayID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_vpn_gateway_route_report" "is_vpn_gateway_route_report_instance" {
			vpn_gateway_id = "%s"
		}

		data "ibm_is_vpn_gateway_route_reports" "is_vpn_gateway_route_reports_instance" {
			vpn_gateway_id = ibm_is_vpn_gateway_route_report.is_vpn_gateway_route_report_instance.vpn_gateway_id
		}
	`, vpnRouteReportVPNGatewayID)
}

func TestDataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		vpnRouteReportRouteNextHopModel := make(map[string]interface{})
		vpnRouteReportRouteNextHopModel["address"] = "192.168.1.1"
		vpnRouteReportRouteNextHopModel["used"] = true

		vpnRouteReportRouteModel := make(map[string]interface{})
		vpnRouteReportRouteModel["as_path"] = []int64{4201065558, 4203065547}
		vpnRouteReportRouteModel["best_path"] = true
		vpnRouteReportRouteModel["next_hops"] = []map[string]interface{}{vpnRouteReportRouteNextHopModel}
		vpnRouteReportRouteModel["peer"] = "10.10.10.1"
		vpnRouteReportRouteModel["prefix"] = "192.168.20.0/24"
		vpnRouteReportRouteModel["valid"] = true
		vpnRouteReportRouteModel["weight"] = int(0)

		model := make(map[string]interface{})
		model["created_at"] = "2019-01-01T12:00:00.000Z"
		model["id"] = "ddf51bec-3424-11e8-b467-0ed5f89f718b"
		model["routes"] = []map[string]interface{}{vpnRouteReportRouteModel}
		model["status"] = "pending"
		model["updated_at"] = "2019-01-01T12:00:00.000Z"

		assert.Equal(t, result, model)
	}

	vpnRouteReportRouteNextHopModel := new(vpcv1.VPNRouteReportRouteNextHop)
	vpnRouteReportRouteNextHopModel.Address = core.StringPtr("192.168.1.1")
	vpnRouteReportRouteNextHopModel.Used = core.BoolPtr(true)

	vpnRouteReportRouteModel := new(vpcv1.VPNRouteReportRoute)
	vpnRouteReportRouteModel.AsPath = []int64{4201065558, 4203065547}
	vpnRouteReportRouteModel.BestPath = core.BoolPtr(true)
	vpnRouteReportRouteModel.NextHops = []vpcv1.VPNRouteReportRouteNextHop{*vpnRouteReportRouteNextHopModel}
	vpnRouteReportRouteModel.Peer = core.StringPtr("10.10.10.1")
	vpnRouteReportRouteModel.Prefix = core.StringPtr("192.168.20.0/24")
	vpnRouteReportRouteModel.Valid = core.BoolPtr(true)
	vpnRouteReportRouteModel.Weight = core.Int64Ptr(int64(0))

	model := new(vpcv1.VPNRouteReport)
	model.CreatedAt = CreateMockDateTime("2019-01-01T12:00:00.000Z")
	model.ID = core.StringPtr("ddf51bec-3424-11e8-b467-0ed5f89f718b")
	model.Routes = []vpcv1.VPNRouteReportRoute{*vpnRouteReportRouteModel}
	model.Status = core.StringPtr("pending")
	model.UpdatedAt = CreateMockDateTime("2019-01-01T12:00:00.000Z")

	result, err := vpc.DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteToMap(t *testing.T) {
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

	result, err := vpc.DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteNextHopToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["address"] = "192.168.3.4"
		model["used"] = true

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.VPNRouteReportRouteNextHop)
	model.Address = core.StringPtr("192.168.3.4")
	model.Used = core.BoolPtr(true)

	result, err := vpc.DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteNextHopToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
