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

func ResourceIBMIsDynamicRouteServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsDynamicRouteServerCreate,
		ReadContext:   resourceIBMIsDynamicRouteServerRead,
		UpdateContext: resourceIBMIsDynamicRouteServerUpdate,
		DeleteContext: resourceIBMIsDynamicRouteServerDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"asn": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The local autonomous system number (ASN) for this dynamic route server.",
			},
			"ips": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The reserved IPs bound to this dynamic route server.",
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
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_dynamic_route_server", "name"),
				Description:  "The name for this dynamic route server. The name is unique across all dynamic route servers in the region.",
			},
			"redistribute_service_routes": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether all service routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `service`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the routeAdditionally, the CIDRs `161.26.0.0/16` (IBM services) and `166.8.0.0/14` (Cloud Service Endpoints) will also be redistributed to all peers through the routing protocol.",
			},
			"redistribute_subnets": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether subnets meet the following conditions will be redistributed through the routing protocol to all peers as route destinations:- The subnet is attached to a routing table in the VPC this dynamic route server is  serving.- The routing table's `accept_routes_from` property includes the value  `dynamic_route_server`The routing protocol will redistribute routes with these subnets as route destinations.",
			},
			"redistribute_user_routes": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether all user routes are redistributed through the routing protocol.All routes will be redistributed to all peers through the routing protocol that meet the following conditions:- The route's property origin is `user`- The route is in a routing table in the VPC that this dynamic route server is serving- The dynamic route server's peer IP is in a subnet that is attached to the routing  table with the route- The dynamic route server's peer IP is in a subnet with the same zone as the route.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
							Required:    true,
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
			"security_groups": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The security groups targeting this dynamic route server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The security group's CRN.",
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
							Description: "The security group's canonical URL.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this security group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name for this security group. The name is unique across all security groups for the VPC.",
						},
					},
				},
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The VPC this dynamic route server resides in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN for this VPC.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Computed:    true,
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
							Description: "The URL for this VPC.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the dynamic route server.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMIsDynamicRouteServerValidator() *validate.ResourceValidator {
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_dynamic_route_server", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsDynamicRouteServerCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createDynamicRouteServerOptions := &vpcv1.CreateDynamicRouteServerOptions{}

	createDynamicRouteServerOptions.SetAsn(int64(d.Get("asn").(int)))
	var ips []vpcv1.ReservedIPIdentityIntf
	for _, v := range d.Get("ips").([]interface{}) {
		value := v.(map[string]interface{})
		ipsItem, err := resourceIBMIsDynamicRouteServerMapToReservedIPIdentity(value)
		if err != nil {
			return diag.FromErr(err)
		}
		ips = append(ips, ipsItem)
	}
	createDynamicRouteServerOptions.SetIps(ips)
	vpcModel, err := resourceIBMIsDynamicRouteServerMapToVPCIdentity(d.Get("vpc.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createDynamicRouteServerOptions.SetVPC(vpcModel)
	if _, ok := d.GetOk("name"); ok {
		createDynamicRouteServerOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("redistribute_service_routes"); ok {
		createDynamicRouteServerOptions.SetRedistributeServiceRoutes(d.Get("redistribute_service_routes").(bool))
	}
	if _, ok := d.GetOk("redistribute_subnets"); ok {
		createDynamicRouteServerOptions.SetRedistributeSubnets(d.Get("redistribute_subnets").(bool))
	}
	if _, ok := d.GetOk("redistribute_user_routes"); ok {
		createDynamicRouteServerOptions.SetRedistributeUserRoutes(d.Get("redistribute_user_routes").(bool))
	}
	if _, ok := d.GetOk("resource_group"); ok {
		resourceGroupModel, err := resourceIBMIsDynamicRouteServerMapToResourceGroupIdentity(d.Get("resource_group.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createDynamicRouteServerOptions.SetResourceGroup(resourceGroupModel)
	}
	if _, ok := d.GetOk("security_groups"); ok {
		var securityGroups []vpcv1.SecurityGroupIdentityIntf
		for _, v := range d.Get("security_groups").([]interface{}) {
			value := v.(map[string]interface{})
			securityGroupsItem, err := resourceIBMIsDynamicRouteServerMapToSecurityGroupIdentity(value)
			if err != nil {
				return diag.FromErr(err)
			}
			securityGroups = append(securityGroups, securityGroupsItem)
		}
		createDynamicRouteServerOptions.SetSecurityGroups(securityGroups)
	}

	dynamicRouteServer, response, err := vpcClient.CreateDynamicRouteServerWithContext(context, createDynamicRouteServerOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateDynamicRouteServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateDynamicRouteServerWithContext failed %s\n%s", err, response))
	}

	d.SetId(*dynamicRouteServer.ID)

	return resourceIBMIsDynamicRouteServerRead(context, d, meta)
}

func resourceIBMIsDynamicRouteServerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getDynamicRouteServerOptions := &vpcv1.GetDynamicRouteServerOptions{}

	getDynamicRouteServerOptions.SetID(d.Id())

	dynamicRouteServer, response, err := vpcClient.GetDynamicRouteServerWithContext(context, getDynamicRouteServerOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetDynamicRouteServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetDynamicRouteServerWithContext failed %s\n%s", err, response))
	}

	if err = d.Set("asn", flex.IntValue(dynamicRouteServer.Asn)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting asn: %s", err))
	}
	ips := []map[string]interface{}{}
	for _, ipsItem := range dynamicRouteServer.Ips {
		ipsItemMap, err := resourceIBMIsDynamicRouteServerReservedIPReferenceToMap(&ipsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		ips = append(ips, ipsItemMap)
	}
	if err = d.Set("ips", ips); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting ips: %s", err))
	}
	if !core.IsNil(dynamicRouteServer.Name) {
		if err = d.Set("name", dynamicRouteServer.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if !core.IsNil(dynamicRouteServer.RedistributeServiceRoutes) {
		if err = d.Set("redistribute_service_routes", dynamicRouteServer.RedistributeServiceRoutes); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting redistribute_service_routes: %s", err))
		}
	}
	if !core.IsNil(dynamicRouteServer.RedistributeSubnets) {
		if err = d.Set("redistribute_subnets", dynamicRouteServer.RedistributeSubnets); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting redistribute_subnets: %s", err))
		}
	}
	if !core.IsNil(dynamicRouteServer.RedistributeUserRoutes) {
		if err = d.Set("redistribute_user_routes", dynamicRouteServer.RedistributeUserRoutes); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting redistribute_user_routes: %s", err))
		}
	}
	if !core.IsNil(dynamicRouteServer.ResourceGroup) {
		resourceGroupMap, err := resourceIBMIsDynamicRouteServerResourceGroupReferenceToMap(dynamicRouteServer.ResourceGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if !core.IsNil(dynamicRouteServer.SecurityGroups) {
		securityGroups := []map[string]interface{}{}
		for _, securityGroupsItem := range dynamicRouteServer.SecurityGroups {
			securityGroupsItemMap, err := resourceIBMIsDynamicRouteServerSecurityGroupReferenceToMap(&securityGroupsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			securityGroups = append(securityGroups, securityGroupsItemMap)
		}
		if err = d.Set("security_groups", securityGroups); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting security_groups: %s", err))
		}
	}
	vpcMap, err := resourceIBMIsDynamicRouteServerVPCReferenceToMap(dynamicRouteServer.VPC)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vpc", []map[string]interface{}{vpcMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vpc: %s", err))
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
	if err = d.Set("lifecycle_state", dynamicRouteServer.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", dynamicRouteServer.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIBMIsDynamicRouteServerUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateDynamicRouteServerOptions := &vpcv1.UpdateDynamicRouteServerOptions{}

	updateDynamicRouteServerOptions.SetID(d.Id())

	hasChange := false

	patchVals := &vpcv1.DynamicRouteServerPatch{}
	if d.HasChange("asn") {
		newAsn := int64(d.Get("asn").(int))
		patchVals.Asn = &newAsn
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("redistribute_service_routes") {
		newRedistributeServiceRoutes := d.Get("redistribute_service_routes").(bool)
		patchVals.RedistributeServiceRoutes = &newRedistributeServiceRoutes
		hasChange = true
	}
	if d.HasChange("redistribute_subnets") {
		newRedistributeSubnets := d.Get("redistribute_subnets").(bool)
		patchVals.RedistributeSubnets = &newRedistributeSubnets
		hasChange = true
	}
	if d.HasChange("redistribute_user_routes") {
		newRedistributeUserRoutes := d.Get("redistribute_user_routes").(bool)
		patchVals.RedistributeUserRoutes = &newRedistributeUserRoutes
		hasChange = true
	}
	updateDynamicRouteServerOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateDynamicRouteServerOptions.DynamicRouteServerPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateDynamicRouteServerWithContext(context, updateDynamicRouteServerOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateDynamicRouteServerWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateDynamicRouteServerWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsDynamicRouteServerRead(context, d, meta)
}

func resourceIBMIsDynamicRouteServerDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteDynamicRouteServerOptions := &vpcv1.DeleteDynamicRouteServerOptions{}

	deleteDynamicRouteServerOptions.SetID(d.Id())

	deleteDynamicRouteServerOptions.SetIfMatch(d.Get("etag").(string))

	_, response, err := vpcClient.DeleteDynamicRouteServerWithContext(context, deleteDynamicRouteServerOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteDynamicRouteServerWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteDynamicRouteServerWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMIsDynamicRouteServerMapToReservedIPIdentity(modelMap map[string]interface{}) (vpcv1.ReservedIPIdentityIntf, error) {
	model := &vpcv1.ReservedIPIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

// func resourceIBMIsDynamicRouteServerMapToReservedIPIdentityReservedIPIdentityByID(modelMap map[string]interface{}) (*vpcv1.ReservedIPIdentityReservedIPIdentityByID, error) {
// 	model := &vpcv1.ReservedIPIdentityReservedIPIdentityByID{}
// 	model.ID = core.StringPtr(modelMap["id"].(string))
// 	return model, nil
// }

// func resourceIBMIsDynamicRouteServerMapToReservedIPIdentityReservedIPIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ReservedIPIdentityReservedIPIdentityByHref, error) {
// 	model := &vpcv1.ReservedIPIdentityReservedIPIdentityByHref{}
// 	model.Href = core.StringPtr(modelMap["href"].(string))
// 	return model, nil
// }

func resourceIBMIsDynamicRouteServerMapToVPCIdentity(modelMap map[string]interface{}) (vpcv1.VPCIdentityIntf, error) {
	model := &vpcv1.VPCIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToVPCIdentityByID(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByID, error) {
	model := &vpcv1.VPCIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToVPCIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByCRN, error) {
	model := &vpcv1.VPCIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToVPCIdentityByHref(modelMap map[string]interface{}) (*vpcv1.VPCIdentityByHref, error) {
	model := &vpcv1.VPCIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToResourceGroupIdentity(modelMap map[string]interface{}) (vpcv1.ResourceGroupIdentityIntf, error) {
	model := &vpcv1.ResourceGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ResourceGroupIdentityByID, error) {
	model := &vpcv1.ResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToSecurityGroupIdentity(modelMap map[string]interface{}) (vpcv1.SecurityGroupIdentityIntf, error) {
	model := &vpcv1.SecurityGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToSecurityGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByID, error) {
	model := &vpcv1.SecurityGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToSecurityGroupIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByCRN, error) {
	model := &vpcv1.SecurityGroupIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerMapToSecurityGroupIdentityByHref(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByHref, error) {
	model := &vpcv1.SecurityGroupIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsDynamicRouteServerReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsDynamicRouteServerReservedIPReferenceDeletedToMap(model.Deleted)
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

func resourceIBMIsDynamicRouteServerReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsDynamicRouteServerResourceGroupReferenceToMap(model *vpcv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIBMIsDynamicRouteServerSecurityGroupReferenceToMap(model *vpcv1.SecurityGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsDynamicRouteServerSecurityGroupReferenceDeletedToMap(model.Deleted)
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

func resourceIBMIsDynamicRouteServerSecurityGroupReferenceDeletedToMap(model *vpcv1.SecurityGroupReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsDynamicRouteServerVPCReferenceToMap(model *vpcv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsDynamicRouteServerVPCReferenceDeletedToMap(model.Deleted)
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

func resourceIBMIsDynamicRouteServerVPCReferenceDeletedToMap(model *vpcv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
