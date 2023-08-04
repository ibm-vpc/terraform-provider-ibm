// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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

func DataSourceIBMIsDynamicRouteServers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsDynamicRouteServersRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `name` property matching the exact specified name.",
			},
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `resource_group.id` property matching the specified identifier.",
			},
			"sort": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-created_at",
				Description: "Sorts the returned collection by the specified property name in ascending order. A `-` may be prepended to the name to sort in descending order. For example, the value `-created_at` sorts the collection by the `created_at` property in descending order, and the value `name` sorts it by the `name` property in ascending order.",
			},
			"dynamic_route_servers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dynamic route servers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this dynamic route server.",
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
				},
			},
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the first page of resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of resources that can be returned by the request.",
			},
			"next": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the next page of resources. This property is present for all pagesexcept the last page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIBMIsDynamicRouteServersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listDynamicRouteServersOptions := &vpcv1.ListDynamicRouteServersOptions{}

	if _, ok := d.GetOk("name"); ok {
		listDynamicRouteServersOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("resource_group_id"); ok {
		listDynamicRouteServersOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("sort"); ok {
		listDynamicRouteServersOptions.SetSort(d.Get("sort").(string))
	}

	var pager *vpcv1.DynamicRouteServersPager
	pager, err = vpcClient.NewDynamicRouteServersPager(listDynamicRouteServersOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] DynamicRouteServersPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("DynamicRouteServersPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIBMIsDynamicRouteServersID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsDynamicRouteServersDynamicRouteServerToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("dynamic_route_servers", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dynamic_route_servers %s", err))
	}

	return nil
}

// dataSourceIBMIsDynamicRouteServersID returns a reasonable ID for the list.
func dataSourceIBMIsDynamicRouteServersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsDynamicRouteServersDynamicRouteServerToMap(model *vpcv1.DynamicRouteServer) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["asn"] = flex.IntValue(model.Asn)
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["crn"] = model.CRN
	modelMap["health_state"] = model.HealthState
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	ips := []map[string]interface{}{}
	for _, ipsItem := range model.Ips {
		ipsItemMap, err := dataSourceIBMIsDynamicRouteServersReservedIPReferenceToMap(&ipsItem)
		if err != nil {
			return modelMap, err
		}
		ips = append(ips, ipsItemMap)
	}
	modelMap["ips"] = ips
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	modelMap["redistribute_service_routes"] = model.RedistributeServiceRoutes
	modelMap["redistribute_subnets"] = model.RedistributeSubnets
	modelMap["redistribute_user_routes"] = model.RedistributeUserRoutes
	resourceGroupMap, err := dataSourceIBMIsDynamicRouteServersResourceGroupReferenceToMap(model.ResourceGroup)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	modelMap["resource_type"] = model.ResourceType
	securityGroups := []map[string]interface{}{}
	for _, securityGroupsItem := range model.SecurityGroups {
		securityGroupsItemMap, err := dataSourceIBMIsDynamicRouteServersSecurityGroupReferenceToMap(&securityGroupsItem)
		if err != nil {
			return modelMap, err
		}
		securityGroups = append(securityGroups, securityGroupsItemMap)
	}
	modelMap["security_groups"] = securityGroups
	vpcMap, err := dataSourceIBMIsDynamicRouteServersVPCReferenceToMap(model.VPC)
	if err != nil {
		return modelMap, err
	}
	modelMap["vpc"] = []map[string]interface{}{vpcMap}
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServersReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServersReservedIPReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsDynamicRouteServersReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServersResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServersSecurityGroupReferenceToMap(model *vpcv1.SecurityGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServersSecurityGroupReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsDynamicRouteServersSecurityGroupReferenceDeletedToMap(model *vpcv1.SecurityGroupReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServersVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServersVPCReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsDynamicRouteServersVPCReferenceDeletedToMap(model *vpcv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServersDynamicRouteServerCollectionFirstToMap(model *vpcv1.DynamicRouteServerCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServersDynamicRouteServerCollectionNextToMap(model *vpcv1.DynamicRouteServerCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}
