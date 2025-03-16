// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isDefaultRoutingTableID             = "default_routing_table"
	isDefaultRoutingTableHref           = "href"
	isDefaultRoutingTableCrn            = "crn"
	isDefaultRoutingTableName           = "name"
	isDefaultRoutingTableResourceType   = "resource_type"
	isDefaultRoutingTableCreatedAt      = "created_at"
	isDefaultRoutingTableLifecycleState = "lifecycle_state"
	isDefaultRoutingTableRoutesList     = "routes"
	isDefaultRoutingTableSubnetsList    = "subnets"
	isDefaultRTVpcID                    = "vpc"
	isDefaultRTDirectLinkIngress        = "route_direct_link_ingress"
	isDefaultRTInternetIngress          = "route_internet_ingress"
	isDefaultRTTransitGatewayIngress    = "route_transit_gateway_ingress"
	isDefaultRTVPCZoneIngress           = "route_vpc_zone_ingress"
	isDefaultRTDefault                  = "is_default"
	isDefaultRTResourceGroup            = "resource_group"
	isDefaultRTResourceGroupHref        = "href"
	isDefaultRTResourceGroupId          = "id"
	isDefaultRTResourceGroupName        = "name"
	isDefaultRTTags                     = "tags"
	isDefaultRTAccessTags               = "access_tags"
	isDefaultRTAccessTagType            = "access"
	isDefaultRTUserTagType              = "user"
)

func DataSourceIBMISVPCDefaultRoutingTable() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISVPCDefaultRoutingTableGet,
		Schema: map[string]*schema.Schema{
			isDefaultRTVpcID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "VPC identifier",
			},
			isDefaultRoutingTableID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing Table ID",
			},
			isDefaultRoutingTableHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Href",
			},
			isDefaultRoutingTableName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Name",
			},
			isDefaultRoutingTableCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Crn",
			},
			isDefaultRoutingTableResourceType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Resource Type",
			},
			isDefaultRoutingTableCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Created At",
			},
			isDefaultRoutingTableLifecycleState: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Default Routing table Lifecycle State",
			},
			isDefaultRTDirectLinkIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Direct Link to this VPC.",
			},
			isDefaultRTInternetIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from the internet. For this to succeed, the VPC must not already have a routing table with this property set to true.",
			},
			isDefaultRTTransitGatewayIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from Transit Gateway to this VPC.",
			},
			isDefaultRTVPCZoneIngress: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "If set to true, this routing table will be used to route traffic that originates from subnets in other zones in this VPC.",
			},
			isDefaultRTDefault: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether this is the default routing table for this VPC",
			},
			isDefaultRoutingTableRoutesList: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route name",
						},

						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Route ID",
						},
					},
				},
			},
			isDefaultRoutingTableSubnetsList: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet name",
						},

						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Subnet ID",
						},
					},
				},
			},
			isDefaultRTResourceGroup: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this volume.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isDefaultRTResourceGroupHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						isDefaultRTResourceGroupId: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						isDefaultRTResourceGroupName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this resource group.",
						},
					},
				},
			},
			isDefaultRTTags: {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      flex.ResourceIBMVPCHash,
			},
			isDefaultRTAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access tags",
			},
		},
	}
}

func dataSourceIBMISVPCDefaultRoutingTableGet(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "(Data) ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	vpcID := d.Get(isDefaultRTVpcID).(string)

	getVpcDefaultRoutingTableOptions := vpcClient.NewGetVPCDefaultRoutingTableOptions(vpcID)
	result, response, err := vpcClient.GetVPCDefaultRoutingTableWithContext(context, getVpcDefaultRoutingTableOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVPCDefaultRoutingTableWithContext failed %s\n%s", err, response), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.SetId(*result.ID)
	if err = d.Set(isDefaultRoutingTableID, *result.ID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableID), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRoutingTableHref, *result.Href); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableHref), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRoutingTableName, *result.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableName), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRoutingTableCrn, *result.CRN); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableCrn), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRoutingTableResourceType, *result.ResourceType); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableResourceType), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	createdAt := *result.CreatedAt
	if err = d.Set(isDefaultRoutingTableCreatedAt, createdAt.String()); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableCreatedAt), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRoutingTableLifecycleState, *result.LifecycleState); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableLifecycleState), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRTDirectLinkIngress, *result.RouteDirectLinkIngress); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTDirectLinkIngress), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRTInternetIngress, *result.RouteInternetIngress); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTInternetIngress), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRTTransitGatewayIngress, *result.RouteTransitGatewayIngress); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTTransitGatewayIngress), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRTVPCZoneIngress, *result.RouteVPCZoneIngress); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTVPCZoneIngress), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRTDefault, *result.IsDefault); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTDefault), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	subnetsInfo := make([]map[string]interface{}, 0)
	for _, subnet := range result.Subnets {
		if subnet.Name != nil && subnet.ID != nil {
			l := map[string]interface{}{
				"name": *subnet.Name,
				"id":   *subnet.ID,
			}
			subnetsInfo = append(subnetsInfo, l)
		}
	}
	if err = d.Set(isDefaultRoutingTableSubnetsList, subnetsInfo); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableSubnetsList), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	routesInfo := make([]map[string]interface{}, 0)
	for _, route := range result.Routes {
		if route.Name != nil && route.ID != nil {
			k := map[string]interface{}{
				"name": *route.Name,
				"id":   *route.ID,
			}
			routesInfo = append(routesInfo, k)
		}
	}
	if err = d.Set(isDefaultRoutingTableRoutesList, routesInfo); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRoutingTableRoutesList), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	resourceGroupList := []map[string]interface{}{}
	if result.ResourceGroup != nil {
		resourceGroupMap := routingTableResourceGroupToMap(*result.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
	}
	if err = d.Set(isDefaultRTResourceGroup, resourceGroupList); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTResourceGroup), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isDefaultRTVpcID, vpcID); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTVpcID), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	tags, err := flex.GetGlobalTagsUsingCRN(meta, *result.CRN, "", isDefaultRTUserTagType)
	if err != nil {
		log.Printf(
			"An error occured during reading of default routing table (%s) tags : %s", d.Id(), err)
	}
	if err = d.Set(isDefaultRTTags, tags); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTTags), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *result.CRN, "", isDefaultRTAccessTagType)
	if err != nil {
		log.Printf(
			"An error occured during reading of default routing table (%s) access tags: %s", d.Id(), err)
	}
	if err = d.Set(isDefaultRTAccessTags, accesstags); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isDefaultRTAccessTags), "ibm_is_vpc_default_routing_table", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	return nil
}
