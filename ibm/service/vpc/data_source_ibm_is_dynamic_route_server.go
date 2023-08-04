// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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

func DataSourceIBMIsDynamicRouteServer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsDynamicRouteServerRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dynamic route server identifier.",
			},
			"asn": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The local autonomous system number (ASN) for this dynamic route server.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the dynamic route server was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this dynamic route server.",
			},
			"health_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this dynamic route server.",
			},
			"ips": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reserved IPs bound to this dynamic route server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this reserved IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the dynamic route server.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this dynamic route server. The name is unique across all dynamic route servers in the region.",
			},
			"redistribute_service_routes": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether all service routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `service`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the routeAdditionally, the CIDRs `161.26.0.0/16` (IBM services) and `166.8.0.0/14` (Cloud Service Endpoints) will also be redistributed to all peers through the routing protocol.",
			},
			"redistribute_subnets": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether subnets meet the following conditions will be redistributed through the routing protocol to all peers as route destinations:- The subnet is attached to a routing table in the VPC this dynamic route server is  serving.- The routing table's `accept_routes_from` property includes the value  `dynamic_route_server`The routing protocol will redistribute routes with these subnets as route destinations.",
			},
			"redistribute_user_routes": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether all user routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `user`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the route.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The resource group for this dynamic route server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this resource group.",
						},
					},
				},
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"security_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The security groups targeting this dynamic route server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group's CRN.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security group's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this security group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this security group. The name is unique across all security groups for the VPC.",
						},
					},
				},
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this dynamic route server resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsDynamicRouteServerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDynamicRouteServerOptions := &vpcv1.GetDynamicRouteServerOptions{}

	getDynamicRouteServerOptions.SetID(d.Get("id").(string))

	dynamicRouteServer, response, err := vpcClient.GetDynamicRouteServerWithContext(context, getDynamicRouteServerOptions)
	if err != nil {
		log.Printf("[DEBUG] GetDynamicRouteServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDynamicRouteServerWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getDynamicRouteServerOptions.ID))

	if err = d.Set("asn", flex.IntValue(dynamicRouteServer.Asn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting asn: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(dynamicRouteServer.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("crn", dynamicRouteServer.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("health_state", dynamicRouteServer.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting health_state: %s", err))
	}

	if err = d.Set("href", dynamicRouteServer.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	ips := []map[string]interface{}{}
	if dynamicRouteServer.Ips != nil {
		for _, modelItem := range dynamicRouteServer.Ips {
			modelMap, err := dataSourceIBMIsDynamicRouteServerReservedIPReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			ips = append(ips, modelMap)
		}
	}
	if err = d.Set("ips", ips); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ips %s", err))
	}

	if err = d.Set("lifecycle_state", dynamicRouteServer.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}

	if err = d.Set("name", dynamicRouteServer.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("redistribute_service_routes", dynamicRouteServer.RedistributeServiceRoutes); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting redistribute_service_routes: %s", err))
	}

	if err = d.Set("redistribute_subnets", dynamicRouteServer.RedistributeSubnets); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting redistribute_subnets: %s", err))
	}

	if err = d.Set("redistribute_user_routes", dynamicRouteServer.RedistributeUserRoutes); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting redistribute_user_routes: %s", err))
	}

	resourceGroup := []map[string]interface{}{}
	if dynamicRouteServer.ResourceGroup != nil {
		modelMap, err := dataSourceIBMIsDynamicRouteServerResourceGroupReferenceToMap(dynamicRouteServer.ResourceGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		resourceGroup = append(resourceGroup, modelMap)
	}
	if err = d.Set("resource_group", resourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
	}

	if err = d.Set("resource_type", dynamicRouteServer.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	securityGroups := []map[string]interface{}{}
	if dynamicRouteServer.SecurityGroups != nil {
		for _, modelItem := range dynamicRouteServer.SecurityGroups {
			modelMap, err := dataSourceIBMIsDynamicRouteServerSecurityGroupReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			securityGroups = append(securityGroups, modelMap)
		}
	}
	if err = d.Set("security_groups", securityGroups); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting security_groups %s", err))
	}

	vpc := []map[string]interface{}{}
	if dynamicRouteServer.VPC != nil {
		modelMap, err := dataSourceIBMIsDynamicRouteServerVPCReferenceToMap(dynamicRouteServer.VPC)
		if err != nil {
			return diag.FromErr(err)
		}
		vpc = append(vpc, modelMap)
	}
	if err = d.Set("vpc", vpc); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
	}

	return nil
}

func dataSourceIBMIsDynamicRouteServerReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServerReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerSecurityGroupReferenceToMap(model *vpcv1.SecurityGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServerSecurityGroupReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerSecurityGroupReferenceDeletedToMap(model *vpcv1.SecurityGroupReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServerVPCReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerVPCReferenceDeletedToMap(model *vpcv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
