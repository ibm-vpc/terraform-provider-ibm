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
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVPNGatewayRouteReports() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNGatewayRouteReportsRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN gateway identifier.",
			},
			"route_reports": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of VPN gateway route reports.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that this route report was created.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this route report.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsVPNGatewayRouteReportsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_vpn_gateway_route_reports", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listVPNRouteReportOptions := &vpcv1.ListVPNRouteReportOptions{}

	listVPNRouteReportOptions.SetVPNGatewayID(d.Get("vpn_gateway_id").(string))

	vpnRouteReportCollection, _, err := vpcClient.ListVPNRouteReportWithContext(context, listVPNRouteReportOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListVPNRouteReportWithContext failed: %s", err.Error()), "(Data) ibm_is_vpn_gateway_route_reports", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsVPNGatewayRouteReportsID(d))

	routeReports := []map[string]interface{}{}
	if vpnRouteReportCollection.RouteReports != nil {
		for _, modelItem := range vpnRouteReportCollection.RouteReports {
			modelMap, err := DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportToMap(&modelItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_vpn_gateway_route_reports", "read", "route_reports-to-map").GetDiag()
			}
			routeReports = append(routeReports, modelMap)
		}
	}
	if err = d.Set("route_reports", routeReports); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting route_reports: %s", err), "(Data) ibm_is_vpn_gateway_route_reports", "read", "set-route_reports").GetDiag()
	}

	return nil
}

// dataSourceIBMIsVPNGatewayRouteReportsID returns a reasonable ID for the list.
func dataSourceIBMIsVPNGatewayRouteReportsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportToMap(model *vpcv1.VPNRouteReport) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["id"] = *model.ID
	routes := []map[string]interface{}{}
	for _, routesItem := range model.Routes {
		routesItemMap, err := DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteToMap(&routesItem)
		if err != nil {
			return modelMap, err
		}
		routes = append(routes, routesItemMap)
	}
	modelMap["routes"] = routes
	modelMap["status"] = *model.Status
	modelMap["updated_at"] = model.UpdatedAt.String()
	return modelMap, nil
}

func DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteToMap(model *vpcv1.VPNRouteReportRoute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["as_path"] = model.AsPath
	modelMap["best_path"] = *model.BestPath
	nextHops := []map[string]interface{}{}
	for _, nextHopsItem := range model.NextHops {
		nextHopsItemMap, err := DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteNextHopToMap(&nextHopsItem)
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

func DataSourceIBMIsVPNGatewayRouteReportsVPNRouteReportRouteNextHopToMap(model *vpcv1.VPNRouteReportRouteNextHop) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	modelMap["used"] = *model.Used
	return modelMap, nil
}
