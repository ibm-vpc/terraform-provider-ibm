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

func DataSourceIBMIsDynamicRouteServerPeer() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsDynamicRouteServerPeerRead,

		Schema: map[string]*schema.Schema{
			"dynamic_route_server_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dynamic route server identifier.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The dynamic route server peer identifier.",
			},
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
			"dynamic_route_server_peer_id": &schema.Schema{
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
	}
}

func dataSourceIBMIsDynamicRouteServerPeerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDynamicRouteServerPeerOptions := &vpcv1.GetDynamicRouteServerPeerOptions{}

	getDynamicRouteServerPeerOptions.SetDynamicRouteServerID(d.Get("dynamic_route_server_id").(string))
	getDynamicRouteServerPeerOptions.SetID(d.Get("id").(string))

	dynamicRouteServerPeer, response, err := vpcClient.GetDynamicRouteServerPeerWithContext(context, getDynamicRouteServerPeerOptions)
	if err != nil {
		log.Printf("[DEBUG] GetDynamicRouteServerPeerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDynamicRouteServerPeerWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *getDynamicRouteServerPeerOptions.DynamicRouteServerID, *getDynamicRouteServerPeerOptions.ID))

	if err = d.Set("asn", flex.IntValue(dynamicRouteServerPeer.Asn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting asn: %s", err))
	}

	if err = d.Set("authentication_enabled", dynamicRouteServerPeer.AuthenticationEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting authentication_enabled: %s", err))
	}

	bfd := []map[string]interface{}{}
	if dynamicRouteServerPeer.Bfd != nil {
		modelMap, err := dataSourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdToMap(dynamicRouteServerPeer.Bfd)
		if err != nil {
			return diag.FromErr(err)
		}
		bfd = append(bfd, modelMap)
	}
	if err = d.Set("bfd", bfd); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting bfd %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(dynamicRouteServerPeer.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("href", dynamicRouteServerPeer.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	if err = d.Set("dynamic_route_server_peer_id", dynamicRouteServerPeer.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dynamic_route_server_peer_id: %s", err))
	}

	ip := []map[string]interface{}{}
	if dynamicRouteServerPeer.IP != nil {
		modelMap, err := dataSourceIBMIsDynamicRouteServerPeerIPToMap(dynamicRouteServerPeer.IP)
		if err != nil {
			return diag.FromErr(err)
		}
		ip = append(ip, modelMap)
	}
	if err = d.Set("ip", ip); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ip %s", err))
	}

	if err = d.Set("lifecycle_state", dynamicRouteServerPeer.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}

	if err = d.Set("name", dynamicRouteServerPeer.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("resource_type", dynamicRouteServerPeer.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	sessions := []map[string]interface{}{}
	if dynamicRouteServerPeer.Sessions != nil {
		for _, modelItem := range dynamicRouteServerPeer.Sessions {
			modelMap, err := dataSourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBgpSessionToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			sessions = append(sessions, modelMap)
		}
	}
	if err = d.Set("sessions", sessions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting sessions %s", err))
	}

	return nil
}

func dataSourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdToMap(model *vpcv1.DynamicRouteServerPeerBfd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mode"] = model.Mode
	modelMap["role"] = model.Role
	sessions := []map[string]interface{}{}
	for _, sessionsItem := range model.Sessions {
		sessionsItemMap, err := dataSourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdSessionToMap(&sessionsItem)
		if err != nil {
			return modelMap, err
		}
		sessions = append(sessions, sessionsItemMap)
	}
	modelMap["sessions"] = sessions
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdSessionToMap(model *vpcv1.DynamicRouteServerPeerBfdSession) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	sourceIPMap, err := dataSourceIBMIsDynamicRouteServerPeerReservedIPReferenceToMap(model.SourceIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["source_ip"] = []map[string]interface{}{sourceIPMap}
	modelMap["state"] = model.State
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeerReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := dataSourceIBMIsDynamicRouteServerPeerReservedIPReferenceDeletedToMap(model.Deleted)
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

func dataSourceIBMIsDynamicRouteServerPeerReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeerIPToMap(model *vpcv1.IP) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	return modelMap, nil
}

func dataSourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBgpSessionToMap(model *vpcv1.DynamicRouteServerPeerBgpSession) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EstablishedAt != nil {
		modelMap["established_at"] = model.EstablishedAt.String()
	}
	sourceIPMap, err := dataSourceIBMIsDynamicRouteServerPeerReservedIPReferenceToMap(model.SourceIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["source_ip"] = []map[string]interface{}{sourceIPMap}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	return modelMap, nil
}
