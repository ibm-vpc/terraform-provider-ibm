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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsDynamicRouteServerRoutes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsDynamicRouteServerRoutesRead,

		Schema: map[string]*schema.Schema{
			"dynamic_route_server_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dynamic route server identifier.",
			},
			"sort": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "-created_at",
				Description: "Sorts the returned collection by the specified property name in ascending order. A `-` may be prepended to the name to sort in descending order. For example, the value `-destination` sorts the collection by the `destination` property in descending order.",
			},
			"peer_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to dynamic route server routes with `peer.id` matching the specified value.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to dynamic route server routes with `type` matching the specified value.",
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
			"routes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dynamic route server routes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"as_path": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The ordered sequence of autonomous systems that network packets will traverse to get to this dynamic route server, per the rules defined in [RFC 4271](https://www.rfc-editor.org/rfc/rfc4271).",
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the route was created.",
						},
						"destination": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination of the route. Each learned route must have a unique combination of`destination`, `source_ip`, and `next_hop`. Similarly, each redistributed route must have a unique combination of `destination` and `next_hop`. The learned route is not added to the VPC routing table if its `destination` is the same as the redistributed route.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this dynamic route server route.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this dynamic route server route.",
						},
						"next_hop": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The next hop packets will be routed to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
								},
							},
						},
						"peer": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
										Description: "The URL for this dynamic route server peer.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this dynamic route server peer.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this dynamic route server peer. The name is unique across all peers for the dynamic route server.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"source_ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The source IP of the dynamic route server used to establish routing protocol withthis dynamic route server peer.This property will be present only when the route `type` is `learned`.",
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
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of this route:- `learned`: route was learned from the dynamic route server peer via the routing  protocol. The learned route was evaluated based on [the best route selection  algorithm](https://cloud.ibm.com/docs/vpc?topic=drs-best-route-selection) to determine if  it was added to the VPC routing table- `redistributed_subnets`: route was redistributed to the dynamic route  server peer via the routing protocol, and route's destination is a subnet IP CIDR  block- `redistributed_user_routes`: route was redistributed to the dynamic route server  peer via the routing protocol, and it is from a VPC routing table with the `origin`  set as `user`- `redistributed_service_routes`: route was redistributed to the dynamic route server  peer via the routing protocol, and it is from a VPC routing table with the `origin`  set as `service`The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the route on which the unexpected property value was encountered.",
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

func dataSourceIBMIsDynamicRouteServerRoutesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listDynamicRouteServerRoutesOptions := &vpcv1.ListDynamicRouteServerRoutesOptions{}

	listDynamicRouteServerRoutesOptions.SetDynamicRouteServerID(d.Get("dynamic_route_server_id").(string))
	if _, ok := d.GetOk("sort"); ok {
		listDynamicRouteServerRoutesOptions.SetSort(d.Get("sort").(string))
	}
	if _, ok := d.GetOk("peer_id"); ok {
		listDynamicRouteServerRoutesOptions.SetPeerID(d.Get("peer_id").(string))
	}
	if _, ok := d.GetOk("type"); ok {
		listDynamicRouteServerRoutesOptions.SetType(d.Get("type").(string))
	}

	var pager *vpcv1.DynamicRouteServerRoutesPager
	pager, err = vpcClient.NewDynamicRouteServerRoutesPager(listDynamicRouteServerRoutesOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] DynamicRouteServerRoutesPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("DynamicRouteServerRoutesPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIBMIsDynamicRouteServerRoutesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerRouteToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("routes", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting routes %s", err))
	}

	return nil
}

// dataSourceIBMIsDynamicRouteServerRoutesID returns a reasonable ID for the list.
func dataSourceIBMIsDynamicRouteServerRoutesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerRouteCollectionFirstToMap(model *vpcv1.DynamicRouteServerRouteCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerRouteCollectionNextToMap(model *vpcv1.DynamicRouteServerRouteCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerRouteToMap(model *vpcv1.DynamicRouteServerRoute) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["as_path"] = model.AsPath
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["destination"] = model.Destination
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	nextHopMap, err := dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerRouteNextHopToMap(model.NextHop)
	if err != nil {
		return modelMap, err
	}
	modelMap["next_hop"] = []map[string]interface{}{nextHopMap}
	peerMap, err := dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerPeerReferenceToMap(model.Peer)
	if err != nil {
		return modelMap, err
	}
	modelMap["peer"] = []map[string]interface{}{peerMap}
	modelMap["resource_type"] = model.ResourceType
	if model.SourceIP != nil {
		sourceIPMap, err := dataSourceIBMIsDynamicRouteServerRoutesReservedIPReferenceToMap(model.SourceIP)
		if err != nil {
			return modelMap, err
		}
		modelMap["source_ip"] = []map[string]interface{}{sourceIPMap}
	}
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerRouteNextHopToMap(model *vpcv1.DynamicRouteServerRouteNextHop) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerPeerReferenceToMap(model *vpcv1.DynamicRouteServerPeerReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerPeerReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsDynamicRouteServerRoutesDynamicRouteServerPeerReferenceDeletedToMap(model *vpcv1.DynamicRouteServerPeerReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerRoutesReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServerRoutesReservedIPReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsDynamicRouteServerRoutesReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
