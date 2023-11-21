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

func DataSourceIBMIsDynamicRouteServerPeers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsDynamicRouteServerPeersRead,

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
				Description: "Sorts the returned collection by the specified property name in ascending order. A `-` may be prepended to the name to sort in descending order. For example, the value `-created_at` sorts the collection by the `created_at` property in descending order, and the value `name` sorts it by the `name` property in ascending order.",
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
			"peers": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of dynamic route server peers.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"asn": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The autonomous system number (ASN) for this dynamic route server peer.",
						},
						"authentication_enabled": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether TCP MD5 authentication key is configured and enabled in this dynamic route server peer.",
						},
						"bfd": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The bidirectional forwarding detection (BFD) configuration for this dynamic route serverpeer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mode": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The bidirectional forwarding detection operating mode on this peer.",
									},
									"role": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The bidirectional forwarding detection role in session initialization.",
									},
									"sessions": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The sessions for this bidirectional forwarding detection for this peer.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"source_ip": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The source IP of the dynamic route server used to establish bidirectional forwardingdetection session with this dynamic route server peer.",
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
												"state": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The current bidirectional forwarding detection session state as seen by this dynamic route server.",
												},
											},
										},
									},
								},
							},
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the dynamic route server peer was created.",
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
						"ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The IP address of this dynamic route server peer.The peer IP must be in a subnet in the VPC this dynamic route server is serving.",
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
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the dynamic route server peer.",
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
						"sessions": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The sessions for this dynamic route server peer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"established_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the BGP session was established.This property will be present only when the session `state` is `established`.",
									},
									"source_ip": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The source IP of the dynamic route server used to establish routing protocol with thisdynamic route server peer.",
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
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The state of the routing protocol with this dynamic route server peer.",
									},
								},
							},
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

func dataSourceIBMIsDynamicRouteServerPeersRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listDynamicRouteServerPeersOptions := &vpcv1.ListDynamicRouteServerPeersOptions{}

	listDynamicRouteServerPeersOptions.SetDynamicRouteServerID(d.Get("dynamic_route_server_id").(string))
	if _, ok := d.GetOk("sort"); ok {
		listDynamicRouteServerPeersOptions.SetSort(d.Get("sort").(string))
	}

	var pager *vpcv1.DynamicRouteServerPeersPager
	pager, err = vpcClient.NewDynamicRouteServerPeersPager(listDynamicRouteServerPeersOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] DynamicRouteServerPeersPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("DynamicRouteServerPeersPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIBMIsDynamicRouteServerPeersID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("peers", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting peers %s", err))
	}

	return nil
}

// dataSourceIBMIsDynamicRouteServerPeersID returns a reasonable ID for the list.
func dataSourceIBMIsDynamicRouteServerPeersID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerCollectionFirstToMap(model *vpcv1.DynamicRouteServerPeerCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerCollectionNextToMap(model *vpcv1.DynamicRouteServerPeerCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerToMap(model *vpcv1.DynamicRouteServerPeer) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["asn"] = flex.IntValue(model.Asn)
	modelMap["authentication_enabled"] = model.AuthenticationEnabled
	bfdMap, err := dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerBfdToMap(model.Bfd)
	if err != nil {
		return modelMap, err
	}
	modelMap["bfd"] = []map[string]interface{}{bfdMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	ipMap, err := dataSourceIBMIsDynamicRouteServerPeersIPToMap(model.IP)
	if err != nil {
		return modelMap, err
	}
	modelMap["ip"] = []map[string]interface{}{ipMap}
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	sessions := []map[string]interface{}{}
	for _, sessionsItem := range model.Sessions {
		sessionsItemMap, err := dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerBgpSessionToMap(&sessionsItem)
		if err != nil {
			return modelMap, err
		}
		sessions = append(sessions, sessionsItemMap)
	}
	modelMap["sessions"] = sessions
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerBfdToMap(model *vpcv1.DynamicRouteServerPeerBfd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mode"] = model.Mode
	modelMap["role"] = model.Role
	sessions := []map[string]interface{}{}
	for _, sessionsItem := range model.Sessions {
		sessionsItemMap, err := dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerBfdSessionToMap(&sessionsItem)
		if err != nil {
			return modelMap, err
		}
		sessions = append(sessions, sessionsItemMap)
	}
	modelMap["sessions"] = sessions
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerBfdSessionToMap(model *vpcv1.DynamicRouteServerPeerBfdSession) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	sourceIPMap, err := dataSourceIBMIsDynamicRouteServerPeersReservedIPReferenceToMap(model.SourceIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["source_ip"] = []map[string]interface{}{sourceIPMap}
	modelMap["state"] = model.State
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeersReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServerPeersReservedIPReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsDynamicRouteServerPeersReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeersIPToMap(model *vpcv1.IP) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeersDynamicRouteServerPeerBgpSessionToMap(model *vpcv1.DynamicRouteServerPeerBgpSession) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EstablishedAt != nil {
		modelMap["established_at"] = model.EstablishedAt.String()
	}
	sourceIPMap, err := dataSourceIBMIsDynamicRouteServerPeersReservedIPReferenceToMap(model.SourceIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["source_ip"] = []map[string]interface{}{sourceIPMap}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	return modelMap, nil
}
