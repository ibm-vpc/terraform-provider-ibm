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

func ResourceIBMIsVPNGatewayRouteReport() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsVPNGatewayRouteReportCreate,
		ReadContext:   resourceIBMIsVPNGatewayRouteReportRead,
		DeleteContext: resourceIBMIsVPNGatewayRouteReportDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"vpn_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPN gateway identifier.",
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
							Elem:        &schema.Schema{Type: schema.TypeInt},
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
			"route_report_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this route report.",
			},
		},
	}
}

func resourceIBMIsVPNGatewayRouteReportCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createVPNRouteReportOptions := &vpcv1.CreateVPNRouteReportOptions{}

	createVPNRouteReportOptions.SetVPNGatewayID(d.Get("vpn_gateway").(string))

	vpnRouteReport, _, err := vpcClient.CreateVPNRouteReportWithContext(context, createVPNRouteReportOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateVPNRouteReportWithContext failed: %s", err.Error()), "ibm_is_vpn_gateway_route_report", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(fmt.Sprintf("%s/%s", *createVPNRouteReportOptions.VPNGatewayID, *vpnRouteReport.ID))

	return resourceIBMIsVPNGatewayRouteReportRead(context, d, meta)
}

func resourceIBMIsVPNGatewayRouteReportRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getVPNRouteReportOptions := &vpcv1.GetVPNRouteReportOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read", "sep-id-parts").GetDiag()
	}

	getVPNRouteReportOptions.SetVPNGatewayID(parts[0])
	getVPNRouteReportOptions.SetID(parts[1])

	vpnRouteReport, response, err := vpcClient.GetVPNRouteReportWithContext(context, getVPNRouteReportOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVPNRouteReportWithContext failed: %s", err.Error()), "ibm_is_vpn_gateway_route_report", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(vpnRouteReport.CreatedAt)); err != nil {
		err = fmt.Errorf("Error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read", "set-created_at").GetDiag()
	}
	routes := []map[string]interface{}{}
	for _, routesItem := range vpnRouteReport.Routes {
		routesItemMap, err := ResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(&routesItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read", "routes-to-map").GetDiag()
		}
		routes = append(routes, routesItemMap)
	}
	if err = d.Set("routes", routes); err != nil {
		err = fmt.Errorf("Error setting routes: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read", "set-routes").GetDiag()
	}
	if err = d.Set("status", vpnRouteReport.Status); err != nil {
		err = fmt.Errorf("Error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read", "set-status").GetDiag()
	}
	if err = d.Set("updated_at", flex.DateTimeToString(vpnRouteReport.UpdatedAt)); err != nil {
		err = fmt.Errorf("Error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("route_report_id", vpnRouteReport.ID); err != nil {
		err = fmt.Errorf("Error setting route_report_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "read", "set-route_report_id").GetDiag()
	}

	return nil
}

func resourceIBMIsVPNGatewayRouteReportDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteVPNRouteReportOptions := &vpcv1.DeleteVPNRouteReportOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "delete", "sep-id-parts").GetDiag()
	}

	deleteVPNRouteReportOptions.SetVPNGatewayID(parts[0])
	deleteVPNRouteReportOptions.SetID(parts[1])

	_, err = vpcClient.DeleteVPNRouteReportWithContext(context, deleteVPNRouteReportOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteVPNRouteReportWithContext failed: %s", err.Error()), "ibm_is_vpn_gateway_route_report", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteToMap(model *vpcv1.VPNRouteReportRoute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["as_path"] = model.AsPath
	modelMap["best_path"] = *model.BestPath
	nextHops := []map[string]interface{}{}
	for _, nextHopsItem := range model.NextHops {
		nextHopsItemMap, err := ResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(&nextHopsItem)
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

func ResourceIBMIsVPNGatewayRouteReportVPNRouteReportRouteNextHopToMap(model *vpcv1.VPNRouteReportRouteNextHop) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	modelMap["used"] = *model.Used
	return modelMap, nil
}
