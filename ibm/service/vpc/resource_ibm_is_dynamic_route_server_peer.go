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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsDynamicRouteServerPeer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsDynamicRouteServerPeerCreate,
		ReadContext:   resourceIBMIsDynamicRouteServerPeerRead,
		UpdateContext: resourceIBMIsDynamicRouteServerPeerUpdate,
		DeleteContext: resourceIBMIsDynamicRouteServerPeerDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"dynamic_route_server_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The dynamic route server identifier.",
			},
			"asn": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The autonomous system number (ASN) for this dynamic route server peer.",
			},
			"bfd": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
							Required:    true,
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
										MinItems:    1,
										MaxItems:    1,
										Required:    true,
										Description: "The source IP of the dynamic route server used to establish bidirectional forwardingdetection session with this dynamic route server peer.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												"deleted": &schema.Schema{
													Type:        schema.TypeList,
													MaxItems:    1,
													Optional:    true,
													Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"more_info": &schema.Schema{
																Type:        schema.TypeString,
																Required:    true,
																Description: "Link to documentation about deleted resources.",
															},
														},
													},
												},
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The URL for this reserved IP.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The unique identifier for this reserved IP.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"state": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The current bidirectional forwarding detection session state as seen by this dynamic route server.",
									},
								},
							},
						},
					},
				},
			},
			"ip": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The IP address of this dynamic route server peer.The peer IP must be in a subnet in the VPC this dynamic route server is serving.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The IP address.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_dynamic_route_server_peer", "name"),
				Description:  "The name for this dynamic route server peer. The name is unique across all peers for the dynamic route server.",
			},
			"authentication_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether TCP MD5 authentication key is configured and enabled in this dynamic route server peer.",
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
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the dynamic route server peer.",
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
							Optional:    true,
							Description: "The date and time that the BGP session was established.This property will be present only when the session `state` is `established`.",
						},
						"source_ip": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The source IP of the dynamic route server used to establish routing protocol with thisdynamic route server peer.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"address": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The IP address.If the address has not yet been selected, the value will be `0.0.0.0`.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										MaxItems:    1,
										Optional:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Required:    true,
													Description: "Link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The URL for this reserved IP.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The unique identifier for this reserved IP.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name for this reserved IP. The name is unique across all reserved IPs in a subnet.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
			"dynamic_route_server_peer_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this dynamic route server peer.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIsDynamicRouteServerPeerValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_dynamic_route_server_peer", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsDynamicRouteServerPeerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createDynamicRouteServerPeerOptions := &vpcv1.CreateDynamicRouteServerPeerOptions{}

	createDynamicRouteServerPeerOptions.SetDynamicRouteServerID(d.Get("dynamic_route_server_id").(string))
	createDynamicRouteServerPeerOptions.SetAsn(int64(d.Get("asn").(int)))
	ipModel, err := resourceIBMIsDynamicRouteServerPeerMapToIP(d.Get("ip.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createDynamicRouteServerPeerOptions.SetIP(ipModel)
	if _, ok := d.GetOk("bfd"); ok {
		bfdModel, err := resourceIBMIsDynamicRouteServerPeerMapToDynamicRouteServerPeerBfdPrototype(d.Get("bfd.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createDynamicRouteServerPeerOptions.SetBfd(bfdModel)
	}
	if _, ok := d.GetOk("md5_authentication_key"); ok {
		createDynamicRouteServerPeerOptions.SetMd5AuthenticationKey(d.Get("md5_authentication_key").(string))
	}
	if _, ok := d.GetOk("name"); ok {
		createDynamicRouteServerPeerOptions.SetName(d.Get("name").(string))
	}

	dynamicRouteServerPeer, response, err := vpcClient.CreateDynamicRouteServerPeerWithContext(context, createDynamicRouteServerPeerOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateDynamicRouteServerPeerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateDynamicRouteServerPeerWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createDynamicRouteServerPeerOptions.DynamicRouteServerID, *dynamicRouteServerPeer.ID))

	return resourceIBMIsDynamicRouteServerPeerRead(context, d, meta)
}

func resourceIBMIsDynamicRouteServerPeerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDynamicRouteServerPeerOptions := &vpcv1.GetDynamicRouteServerPeerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getDynamicRouteServerPeerOptions.SetDynamicRouteServerID(parts[0])
	getDynamicRouteServerPeerOptions.SetID(parts[1])

	dynamicRouteServerPeer, response, err := vpcClient.GetDynamicRouteServerPeerWithContext(context, getDynamicRouteServerPeerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDynamicRouteServerPeerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDynamicRouteServerPeerWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("asn", flex.IntValue(dynamicRouteServerPeer.Asn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting asn: %s", err))
	}
	if !core.IsNil(dynamicRouteServerPeer.Bfd) {
		bfdMap, err := resourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdToMap(dynamicRouteServerPeer.Bfd)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("bfd", []map[string]interface{}{bfdMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting bfd: %s", err))
		}
	}
	ipMap, err := resourceIBMIsDynamicRouteServerPeerIPToMap(dynamicRouteServerPeer.IP)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("ip", []map[string]interface{}{ipMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ip: %s", err))
	}
	if !core.IsNil(dynamicRouteServerPeer.Name) {
		if err = d.Set("name", dynamicRouteServerPeer.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if err = d.Set("authentication_enabled", dynamicRouteServerPeer.AuthenticationEnabled); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting authentication_enabled: %s", err))
	}
	if err = d.Set("created_at", flex.DateTimeToString(dynamicRouteServerPeer.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", dynamicRouteServerPeer.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", dynamicRouteServerPeer.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", dynamicRouteServerPeer.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	sessions := []map[string]interface{}{}
	for _, sessionsItem := range dynamicRouteServerPeer.Sessions {
		sessionsItemMap, err := resourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBgpSessionToMap(&sessionsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		sessions = append(sessions, sessionsItemMap)
	}
	if err = d.Set("sessions", sessions); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting sessions: %s", err))
	}
	if err = d.Set("dynamic_route_server_peer_id", dynamicRouteServerPeer.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting dynamic_route_server_peer_id: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIBMIsDynamicRouteServerPeerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateDynamicRouteServerPeerOptions := &vpcv1.UpdateDynamicRouteServerPeerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateDynamicRouteServerPeerOptions.SetDynamicRouteServerID(parts[0])
	updateDynamicRouteServerPeerOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.DynamicRouteServerPeerPatch{}
	if d.HasChange("dynamic_route_server_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "dynamic_route_server_id"))
	}
	if d.HasChange("asn") {
		newAsn := int64(d.Get("asn").(int))
		patchVals.Asn = &newAsn
		hasChange = true
	}
	if d.HasChange("bfd") {
		bfd, err := resourceIBMIsDynamicRouteServerPeerMapToDynamicRouteServerPeerBfdPatch(d.Get("bfd.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		patchVals.Bfd = bfd
		hasChange = true
	}
	if d.HasChange("md5_authentication_key") {
		newMd5AuthenticationKey := d.Get("md5_authentication_key").(string)
		patchVals.Md5AuthenticationKey = &newMd5AuthenticationKey
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	updateDynamicRouteServerPeerOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateDynamicRouteServerPeerOptions.DynamicRouteServerPeerPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateDynamicRouteServerPeerWithContext(context, updateDynamicRouteServerPeerOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateDynamicRouteServerPeerWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateDynamicRouteServerPeerWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsDynamicRouteServerPeerRead(context, d, meta)
}

func resourceIBMIsDynamicRouteServerPeerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDynamicRouteServerPeerOptions := &vpcv1.DeleteDynamicRouteServerPeerOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDynamicRouteServerPeerOptions.SetDynamicRouteServerID(parts[0])
	deleteDynamicRouteServerPeerOptions.SetID(parts[1])

	deleteDynamicRouteServerPeerOptions.SetIfMatch(d.Get("etag").(string))

	_, response, err := vpcClient.DeleteDynamicRouteServerPeerWithContext(context, deleteDynamicRouteServerPeerOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteDynamicRouteServerPeerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteDynamicRouteServerPeerWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMIsDynamicRouteServerPeerMapToIP(modelMap map[string]interface{}) (*vpcv1.IP, error) {
	model := &vpcv1.IP{}
	model.Address = core.StringPtr(modelMap["address"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerPeerMapToDynamicRouteServerPeerBfdPrototype(modelMap map[string]interface{}) (*vpcv1.DynamicRouteServerPeerBfdPrototype, error) {
	model := &vpcv1.DynamicRouteServerPeerBfdPrototype{}
	if modelMap["role"] != nil && modelMap["role"].(string) != "" {
		model.Role = core.StringPtr(modelMap["role"].(string))
	}
	return model, nil
}

func resourceIBMIsDynamicRouteServerPeerMapToDynamicRouteServerPeerBfdPatch(modelMap map[string]interface{}) (*vpcv1.DynamicRouteServerPeerBfdPatch, error) {
	model := &vpcv1.DynamicRouteServerPeerBfdPatch{}
	if modelMap["role"] != nil && modelMap["role"].(string) != "" {
		model.Role = core.StringPtr(modelMap["role"].(string))
	}
	return model, nil
}

func resourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdToMap(model *vpcv1.DynamicRouteServerPeerBfd) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mode"] = model.Mode
	modelMap["role"] = model.Role
	sessions := []map[string]interface{}{}
	for _, sessionsItem := range model.Sessions {
		sessionsItemMap, err := resourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdSessionToMap(&sessionsItem)
		if err != nil {
			return modelMap, err
		}
		sessions = append(sessions, sessionsItemMap)
	}
	modelMap["sessions"] = sessions
	return modelMap, nil
}

func resourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBfdSessionToMap(model *vpcv1.DynamicRouteServerPeerBfdSession) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	sourceIPMap, err := resourceIBMIsDynamicRouteServerPeerReservedIPReferenceToMap(model.SourceIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["source_ip"] = []map[string]interface{}{sourceIPMap}
	modelMap["state"] = model.State
	return modelMap, nil
}

func resourceIBMIsDynamicRouteServerPeerReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsDynamicRouteServerPeerReservedIPReferenceDeletedToMap(model.Deleted)
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

func resourceIBMIsDynamicRouteServerPeerReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsDynamicRouteServerPeerIPToMap(model *vpcv1.IP) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	return modelMap, nil
}

func resourceIBMIsDynamicRouteServerPeerDynamicRouteServerPeerBgpSessionToMap(model *vpcv1.DynamicRouteServerPeerBgpSession) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.EstablishedAt != nil {
		modelMap["established_at"] = model.EstablishedAt.String()
	}
	sourceIPMap, err := resourceIBMIsDynamicRouteServerPeerReservedIPReferenceToMap(model.SourceIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["source_ip"] = []map[string]interface{}{sourceIPMap}
	if model.State != nil {
		modelMap["state"] = model.State
	}
	return modelMap, nil
}
