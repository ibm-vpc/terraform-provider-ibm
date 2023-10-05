// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsInstanceNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsInstanceNetworkAttachmentCreate,
		ReadContext:   resourceIBMIsInstanceNetworkAttachmentRead,
		UpdateContext: resourceIBMIsInstanceNetworkAttachmentUpdate,
		DeleteContext: resourceIBMIsInstanceNetworkAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				parts, err := flex.SepIdParts(d.Id(), "/")
				if err != nil {
					return nil, err
				}
				d.Set("instance", parts[0])
				d.Set("virtual_network_interface", parts[1])
				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"instance": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_network_attachment", "instance_id"),
				Description:  "The virtual server instance identifier.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_network_attachment", "name"),
				Description:  "The name for this instance network attachment. The name is unique across all network attachments for the instance.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the instance network attachment was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this instance network attachment.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the instance network attachment.",
			},
			"port_speed": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port speed for this instance network attachment in Mbps.",
			},

			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The instance network attachment type.",
			},
			"instance_network_attachment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this instance network attachment.",
			},

			// vni properties
			"virtual_network_interface": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"primary_ip"},
				Description:   "The virtual network interface id for this instance network attachment.",
			},
			"allow_ip_spoofing": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether source IP spoofing is allowed on this interface. If `false`, source IP spoofing is prevented on this interface. If `true`, source IP spoofing is allowed on this interface.",
			},
			"auto_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Indicates whether this virtual network interface will be automatically deleted when`target` is deleted.",
			},
			"enable_infrastructure_nat": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If `true`:- The VPC infrastructure performs any needed NAT operations.- `floating_ips` must not have more than one floating IP.If `false`:- Packets are passed unchanged to/from the network interface,  allowing the workload to perform any needed NAT operations.- `allow_ip_spoofing` must be `false`.- If the virtual network interface is attached:  - The target `resource_type` must be `bare_metal_server_network_attachment`.  - The target `interface_type` must not be `hipersocket`.",
			},
			"ips": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The reserved IPs bound to this virtual network interface.May be empty when `lifecycle_state` is `pending`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
						"auto_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: "Indicates whether this reserved IP member will be automatically deleted when either target is deleted, or the reserved IP is unbound.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this reserved IP.",
						},
						"reserved_ip": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unique identifier for this reserved IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
			"vni_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_virtual_network_interface", "vni_name"),
				Description:  "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
			},
			"primary_ip": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The primary IP address of the virtual network interface for the instance networkattachment.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
						"reserved_ip": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique identifier for this reserved IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The resource group id for this virtual network interface.",
			},
			"vni_resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"security_groups": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: "The security groups for this virtual network interface.",
			},
			"subnet": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The associated subnet id.",
			},
		},
	}
}

func ResourceIBMIsInstanceNetworkAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
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

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_network_attachment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsInstanceNetworkAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	createInstanceNetworkAttachmentOptions := &vpcv1.CreateInstanceNetworkAttachmentOptions{}

	createInstanceNetworkAttachmentOptions.SetInstanceID(d.Get("instance_id").(string))
	virtualNetworkInterfaceModel, err := resourceIBMIsInstanceNetworkAttachmentMapToInstanceNetworkAttachmentPrototypeVirtualNetworkInterface(d.Get("virtual_network_interface.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createInstanceNetworkAttachmentOptions.SetVirtualNetworkInterface(virtualNetworkInterfaceModel)
	if _, ok := d.GetOk("name"); ok {
		createInstanceNetworkAttachmentOptions.SetName(d.Get("name").(string))
	}

	instanceNetworkAttachment, response, err := vpcClient.CreateInstanceNetworkAttachmentWithContext(context, createInstanceNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s", *createInstanceNetworkAttachmentOptions.InstanceID, *instanceNetworkAttachment.ID))

	return resourceIBMIsInstanceNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsInstanceNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceNetworkAttachmentOptions := &vpcv1.GetInstanceNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
	getInstanceNetworkAttachmentOptions.SetID(parts[1])

	instanceNetworkAttachment, response, err := vpcClient.GetInstanceNetworkAttachmentWithContext(context, getInstanceNetworkAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}
	// attachment details
	if !core.IsNil(instanceNetworkAttachment.Name) {
		if err = d.Set("name", instanceNetworkAttachment.Name); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(instanceNetworkAttachment.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
	}

	if err = d.Set("href", instanceNetworkAttachment.Href); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", instanceNetworkAttachment.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("port_speed", flex.IntValue(instanceNetworkAttachment.PortSpeed)); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
	}
	if err = d.Set("resource_type", instanceNetworkAttachment.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
	}
	if err = d.Set("type", instanceNetworkAttachment.Type); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
	}
	if err = d.Set("instance_network_attachment_id", instanceNetworkAttachment.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance_network_attachment_id: %s", err))
	}
	// vni details
	vniId := *instanceNetworkAttachment.VirtualNetworkInterface.ID
	if err = d.Set("virtual_network_interface", vniId); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
	}
	getVniOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
		ID: &vniId,
	}
	vniDetails, response, err := vpcClient.GetVirtualNetworkInterface(getVniOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetVirtualNetworkInterface failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetVirtualNetworkInterface failed %s\n%s", err, response))
	}
	if err = d.Set("allow_ip_spoofing", vniDetails.AllowIPSpoofing); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting allow_ip_spoofing: %s", err))
	}
	if err = d.Set("auto_delete", vniDetails.AutoDelete); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting auto_delete: %s", err))
	}
	if err = d.Set("enable_infrastructure_nat", vniDetails.EnableInfrastructureNat); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting enable_infrastructure_nat: %s", err))
	}
	if err = d.Set("vni_name", vniDetails.Name); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting vni_name: %s", err))
	}
	if err = d.Set("resource_group", vniDetails.ResourceGroup.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_group: %s", err))
	}
	if err = d.Set("vni_resource_type", vniDetails.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting vni_resource_type: %s", err))
	}
	primaryipId := *instanceNetworkAttachment.PrimaryIP.ID
	if !core.IsNil(vniDetails.Ips) {
		ips := []map[string]interface{}{}
		for _, ipsItem := range vniDetails.Ips {
			if *ipsItem.ID != primaryipId {
				ipsItemMap, err := resourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(&ipsItem)
				if err != nil {
					return diag.FromErr(err)
				}
				ips = append(ips, ipsItemMap)
			}
		}
		if err = d.Set("ips", ips); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting ips: %s", err))
		}
	}

	if !core.IsNil(vniDetails.SecurityGroups) {
		securityGroups := make([]string, 0)
		for _, securityGroupsItem := range vniDetails.SecurityGroups {
			if securityGroupsItem.ID != nil {
				securityGroups = append(securityGroups, *securityGroupsItem.ID)
			}
		}
		if err = d.Set("security_groups", securityGroups); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting security_groups for vni: %s", err))
		}
	}

	primaryIPMap, err := resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceToMap(instanceNetworkAttachment.PrimaryIP)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("primary_ip", []map[string]interface{}{primaryIPMap}); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting primary_ip: %s", err))
	}

	if err = d.Set("subnet", *instanceNetworkAttachment.Subnet.ID); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting subnet: %s", err))
	}
	if err = d.Set("instance", parts[0]); err != nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error setting instance: %s", err))
	}

	return nil
}

func resourceIBMIsInstanceNetworkAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateInstanceNetworkAttachmentOptions := &vpcv1.UpdateInstanceNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
	updateInstanceNetworkAttachmentOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.InstanceNetworkAttachmentPatch{}
	if d.HasChange("instance_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "instance_id"))
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	if hasChange {
		updateInstanceNetworkAttachmentOptions.InstanceNetworkAttachmentPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateInstanceNetworkAttachmentWithContext(context, updateInstanceNetworkAttachmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsInstanceNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsInstanceNetworkAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteInstanceNetworkAttachmentOptions := &vpcv1.DeleteInstanceNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteInstanceNetworkAttachmentOptions.SetInstanceID(parts[0])
	deleteInstanceNetworkAttachmentOptions.SetID(parts[1])

	response, err := vpcClient.DeleteInstanceNetworkAttachmentWithContext(context, deleteInstanceNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteInstanceNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteInstanceNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToInstanceNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap map[string]interface{}) (vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceIntf, error) {
	model := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterface{}
	if modelMap["allow_ip_spoofing"] != nil {
		model.AllowIPSpoofing = core.BoolPtr(modelMap["allow_ip_spoofing"].(bool))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["enable_infrastructure_nat"] != nil {
		model.EnableInfrastructureNat = core.BoolPtr(modelMap["enable_infrastructure_nat"].(bool))
	}
	if modelMap["ips"] != nil {
		ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
		for _, ipsItem := range modelMap["ips"].([]interface{}) {
			ipsItemModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(ipsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			ips = append(ips, ipsItemModel)
		}
		model.Ips = ips
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && len(modelMap["resource_group"].([]interface{})) > 0 {
		ResourceGroupModel, err := resourceIBMIsInstanceNetworkAttachmentMapToResourceGroupIdentity(modelMap["resource_group"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceGroup = ResourceGroupModel
	}
	if modelMap["security_groups"] != nil {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		for _, securityGroupsItem := range modelMap["security_groups"].([]interface{}) {
			securityGroupsItemModel, err := resourceIBMIsInstanceNetworkAttachmentMapToSecurityGroupIdentity(securityGroupsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			securityGroups = append(securityGroups, securityGroupsItemModel)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContext(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext{}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextIntf, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext{}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToResourceGroupIdentity(modelMap map[string]interface{}) (vpcv1.ResourceGroupIdentityIntf, error) {
	model := &vpcv1.ResourceGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ResourceGroupIdentityByID, error) {
	model := &vpcv1.ResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSecurityGroupIdentity(modelMap map[string]interface{}) (vpcv1.SecurityGroupIdentityIntf, error) {
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

func resourceIBMIsInstanceNetworkAttachmentMapToSecurityGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByID, error) {
	model := &vpcv1.SecurityGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSecurityGroupIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByCRN, error) {
	model := &vpcv1.SecurityGroupIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSecurityGroupIdentityByHref(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByHref, error) {
	model := &vpcv1.SecurityGroupIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentity(modelMap map[string]interface{}) (vpcv1.SubnetIdentityIntf, error) {
	model := &vpcv1.SubnetIdentity{}
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

func resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentityByID(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByID, error) {
	model := &vpcv1.SubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByCRN, error) {
	model := &vpcv1.SubnetIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentityByHref(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByHref, error) {
	model := &vpcv1.SubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentMapToInstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext(modelMap map[string]interface{}) (*vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext, error) {
	model := &vpcv1.InstanceNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeInstanceNetworkAttachmentContext{}
	if modelMap["allow_ip_spoofing"] != nil {
		model.AllowIPSpoofing = core.BoolPtr(modelMap["allow_ip_spoofing"].(bool))
	}
	if modelMap["auto_delete"] != nil {
		model.AutoDelete = core.BoolPtr(modelMap["auto_delete"].(bool))
	}
	if modelMap["enable_infrastructure_nat"] != nil {
		model.EnableInfrastructureNat = core.BoolPtr(modelMap["enable_infrastructure_nat"].(bool))
	}
	if modelMap["ips"] != nil {
		ips := []vpcv1.VirtualNetworkInterfaceIPPrototypeIntf{}
		for _, ipsItem := range modelMap["ips"].([]interface{}) {
			ipsItemModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(ipsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			ips = append(ips, ipsItemModel)
		}
		model.Ips = ips
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["primary_ip"] != nil && len(modelMap["primary_ip"].([]interface{})) > 0 {
		PrimaryIPModel, err := resourceIBMIsInstanceNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && len(modelMap["resource_group"].([]interface{})) > 0 {
		ResourceGroupModel, err := resourceIBMIsInstanceNetworkAttachmentMapToResourceGroupIdentity(modelMap["resource_group"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceGroup = ResourceGroupModel
	}
	if modelMap["security_groups"] != nil {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		for _, securityGroupsItem := range modelMap["security_groups"].([]interface{}) {
			securityGroupsItemModel, err := resourceIBMIsInstanceNetworkAttachmentMapToSecurityGroupIdentity(securityGroupsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			securityGroups = append(securityGroups, securityGroupsItemModel)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := resourceIBMIsInstanceNetworkAttachmentMapToSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	return model, nil
}

func resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = model.Href
	modelMap["reserved_ip"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIBMIsInstanceNetworkAttachmentReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsInstanceNetworkAttachmentSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsInstanceNetworkAttachmentSubnetReferenceDeletedToMap(model.Deleted)
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

func resourceIBMIsInstanceNetworkAttachmentSubnetReferenceDeletedToMap(model *vpcv1.SubnetReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
