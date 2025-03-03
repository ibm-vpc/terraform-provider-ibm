// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsVPNGatewayServiceConnection() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNGatewayServiceConnectionRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_name", "vpn_gateway"},
				Description:  "The VPN gateway identifier.",
			},
			"vpn_gateway_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_name", "vpn_gateway"},
				Description:  "The VPN gateway name.",
			},
			"vpn_gateway_connection": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_connection", "vpn_gateway_connection_name"},
				Description:  "The VPN gateway connection identifier.",
			},
			"vpn_gateway_connection_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"vpn_gateway_connection", "vpn_gateway_connection_name"},
				Description:  "The VPN gateway connection name.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this VPN service connection was created.",
			},
			"creator": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for transit gateway resource.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for transit gateway resource.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this VPN gateway service connection",
			},
			"lifecycle_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current `lifecycle_state` (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for this lifecycle state:- `internal_error`: internal error (contact IBM support)- `resource_suspended_by_provider`: The resource has been suspended (contact IBM  support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about the reason for this lifecycle state.",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the VPN service connection.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this service connection:- `up`: operating normally- `degraded`: operating with compromised performance- `down`: not operational.",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current VPN service connection status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason. The enumerated values for this property may https://cloud.ibm.com/apidocs/vpc#property-value-expansion in the future.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the reason for this VPN service connection's status.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsVPNGatewayServiceConnectionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	vpn_gateway_id := d.Get("vpn_gateway").(string)
	vpn_gateway_name := d.Get("vpn_gateway_name").(string)
	vpn_gateway_connection := d.Get("vpn_gateway_connection").(string)
	vpn_gateway_connection_name := d.Get("vpn_gateway_connection_name").(string)

	var vpnGatewayServiceConn vpcv1.VPNServiceConnection

	if vpn_gateway_name != "" {
		listvpnGWOptions := vpcClient.NewListVPNGatewaysOptions()

		start := ""
		allrecs := []vpcv1.VPNGatewayIntf{}
		for {
			if start != "" {
				listvpnGWOptions.Start = &start
			}
			availableVPNGateways, detail, err := vpcClient.ListVPNGatewaysWithContext(context, listvpnGWOptions)
			if err != nil || availableVPNGateways == nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error reading list of VPN Gateways:%s\n%s", err, detail))
			}
			start = flex.GetNext(availableVPNGateways.Next)
			allrecs = append(allrecs, availableVPNGateways.VPNGateways...)
			if start == "" {
				break
			}
		}
		vpn_gateway_found := false
		for _, vpnGatewayIntfItem := range allrecs {
			if *vpnGatewayIntfItem.(*vpcv1.VPNGateway).Name == vpn_gateway_name {
				vpnGateway := vpnGatewayIntfItem.(*vpcv1.VPNGateway)
				vpn_gateway_id = *vpnGateway.ID
				vpn_gateway_found = true
				break
			}
		}
		if !vpn_gateway_found {
			log.Printf("[DEBUG] No vpn gateway found with given name %s", vpn_gateway_name)
			return diag.FromErr(fmt.Errorf("No vpn gateway found with given name %s", vpn_gateway_name))
		}
	}

	if vpn_gateway_connection_name != "" {
		listvpnGWConnectionOptions := vpcClient.NewListVPNGatewayConnectionsOptions(vpn_gateway_id)

		availableVPNGatewayConnections, detail, err := vpcClient.ListVPNGatewayConnections(listvpnGWConnectionOptions)
		if err != nil || availableVPNGatewayConnections == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error reading list of VPN Gateway service Connections:%s\n%s", err, detail))
		}

		vpn_gateway_conn_found := false
		for _, connectionItem := range availableVPNGatewayConnections.Connections {
			switch reflect.TypeOf(connectionItem).String() {
			case "*vpcv1.VPNGatewayConnection":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnection)
					if *connection.Name == vpn_gateway_connection_name {
						vpn_gateway_connection = *connection.ID
						vpn_gateway_conn_found = true
						break
					}
				}
			case "*vpcv1.VPNGatewayConnectionRouteMode":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnectionRouteMode)
					if *connection.Name == vpn_gateway_connection_name {
						vpn_gateway_connection = *connection.ID
						vpn_gateway_conn_found = true
						break
					}
				}
			case "*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionStaticRouteMode)
					if *connection.Name == vpn_gateway_connection_name {
						vpn_gateway_connection = *connection.ID
						vpn_gateway_conn_found = true
						break
					}
				}
			case "*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionDynamicRouteMode":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnectionRouteModeVPNGatewayConnectionDynamicRouteMode)
					if *connection.Name == vpn_gateway_connection_name {
						vpn_gateway_connection = *connection.ID
						vpn_gateway_conn_found = true
						break
					}
				}
			case "*vpcv1.VPNGatewayConnectionPolicyMode":
				{
					connection := connectionItem.(*vpcv1.VPNGatewayConnectionPolicyMode)
					if *connection.Name == vpn_gateway_connection_name {
						vpn_gateway_connection = *connection.ID
						vpn_gateway_conn_found = true
						break
					}
				}
			}
		}
		if !vpn_gateway_conn_found {
			return diag.FromErr(fmt.Errorf("VPN gateway connection %s not found", vpn_gateway_connection_name))
		}
	}

	getVPNGatewayServiceConnectionOptions := &vpcv1.GetVPNGatewayServiceConnectionOptions{}

	getVPNGatewayServiceConnectionOptions.SetVPNGatewayID(vpn_gateway_id)
	getVPNGatewayServiceConnectionOptions.SetID(vpn_gateway_connection)

	vpnGatewayServiceConnection, response, err := vpcClient.GetVPNGatewayServiceConnectionWithContext(context, getVPNGatewayServiceConnectionOptions)
	if err != nil || vpnGatewayServiceConnection == nil {
		log.Printf("[DEBUG] GetVPNGatewayServiceConnectionWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVPNGatewayServiceConnectionWithContext failed %s\n%s", err, response))
	}
	vpnGatewayServiceConn = *vpnGatewayServiceConnection

	d.SetId(*vpnGatewayServiceConn.ID)

	if err = d.Set("created_at", flex.DateTimeToString(vpnGatewayServiceConn.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err := d.Set("creator", resourceVPNGatewayServiceConnectionFlattenCreator(*vpnGatewayServiceConn.Creator)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting creator: %s", err))
	}
	if err := d.Set("lifecycle_reasons", resourceVPNGatewayServiceConnectionFlattenLifecycleReasons(vpnGatewayServiceConn.LifecycleReasons)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_reasons: %s", err))
	}
	if err = d.Set("lifecycle_state", vpnGatewayServiceConn.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("status", vpnGatewayServiceConn.Status); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status: %s", err))
	}
	if err := d.Set("status_reasons", resourceVPNGatewayServiceConnectionFlattenStateReasons(vpnGatewayServiceConn.StatusReasons)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting status_reasons: %s", err))
	}
	return nil
}
