// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
 */

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVPNGatewayRouteReport() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNGatewayRouteReportRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN gateway identifier.",
			},
			"is_vpn_gateway_route_report_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN route report identifier.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this route report was created.",
			},
			"routes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The routes of this report.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"as_path": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "AS path numbers of this route.",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"best_path": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this route is best path.",
						},
						"next_hops": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "Next hop list of this route.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unicast IP address, which must not be any of the following values:- `0.0.0.0` (the sentinel IP address)- `224.0.0.0` to `239.255.255.255` (multicast IP addresses)- `255.255.255.255` (the broadcast IP address)This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 addresses in the future.",
									},
									"used": &schema.Schema{
										Type:        schema.TypeBool,
										Computed:    true,
										Description: "The next hop is used for traffic forward.",
									},
								},
							},
						},
						"peer": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The peer of this route.",
						},
						"prefix": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination of this route.",
						},
						"valid": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether this route is valid.",
						},
						"weight": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Weight of this route.",
						},
					},
				},
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Route report status. The list of enumerated values for this property may expand in the future. Code and processes using this field must tolerate unexpected values.- `pending`: generic routing encapsulation tunnel attached- `complete`: generic routing encapsulation tunnel detached.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this route report was updated.",
			},
		},
	}
}

func dataSourceIBMIsVPNGatewayRouteReportRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_vpn_gateway_route_report", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVPNRouteReportOptions := &vpcv1.GetVPNRouteReportOptions{}

	getVPNRouteReportOptions.SetVPNGatewayID(d.Get("vpn_gateway_id").(string))
	getVPNRouteReportOptions.SetID(d.Get("is_vpn_gateway_route_report_id").(string))

	vpnRouteReport, _, err := vpcClient.GetVPNRouteReportWithContext(context, getVPNRouteReportOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVPNRouteReportWithContext failed: %s", err.Error()), "(Data) ibm_is_vpn_gateway_route_report", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *getVPNRouteReportOptions.VPNGatewayID, *getVPNRouteReportOptions.ID))

	if err = d.Set("created_at", flex.DateTimeToString(vpnRouteReport.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting created_at: %s", err), "(Data) ibm_is_vpn_gateway_route_report", "read", "set-created_at").GetDiag()
	}

	routes := []map[string]interface{}{}
	if vpnRouteReport.Routes != nil {
		for _, modelItem := range vpnRouteReport.Routes {
			modelMap, err := DataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_vpn_gateway_route_report", "read", "routes-to-map").GetDiag()
			}
			routes = append(routes, modelMap)
		}
	}
	if err = d.Set("routes", routes); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting routes: %s", err), "(Data) ibm_is_vpn_gateway_route_report", "read", "set-routes").GetDiag()
	}

	if err = d.Set("status", vpnRouteReport.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting status: %s", err), "(Data) ibm_is_vpn_gateway_route_report", "read", "set-status").GetDiag()
	}

	if err = d.Set("updated_at", flex.DateTimeToString(vpnRouteReport.UpdatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting updated_at: %s", err), "(Data) ibm_is_vpn_gateway_route_report", "read", "set-updated_at").GetDiag()
	}

	return nil
}

func DataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(model *vpcv1.VPNRouteReportRoute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["as_path"] = model.AsPath
	modelMap["best_path"] = *model.BestPath
	nextHops := []map[string]interface{}{}
	for _, nextHopsItem := range model.NextHops {
		nextHopsItemMap, err := DataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(&nextHopsItem)
		if err != nil {
			return modelMap, err
		}
		nextHops = append(nextHops, nextHopsItemMap)
	}
	modelMap["next_hops"] = nextHops
	modelMap["peer"] = *model.Peer
	modelMap["prefix"] = *model.Prefix
	modelMap["valid"] = *model.Valid
	modelMap["weight"] = flex.IntValue(model.Weight)
	return modelMap, nil
}

func DataSourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(model *vpcv1.VPNRouteReportRouteNextHop) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	modelMap["used"] = *model.Used
	return modelMap, nil
}
