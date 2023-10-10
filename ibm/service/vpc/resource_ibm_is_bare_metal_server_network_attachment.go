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

func ResourceIBMIsBareMetalServerNetworkAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsBareMetalServerNetworkAttachmentCreate,
		ReadContext:   resourceIBMIsBareMetalServerNetworkAttachmentRead,
		UpdateContext: resourceIBMIsBareMetalServerNetworkAttachmentUpdate,
		DeleteContext: resourceIBMIsBareMetalServerNetworkAttachmentDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"bare_metal_server": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "bare_metal_server"),
				Description:  "The bare metal server identifier.",
			},

			"network_attachment": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The network attachment's id.",
			},
			"interface_type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "interface_type"),
				Description:  "The network attachment's interface type:- `hipersocket`: a virtual network device that provides high-speed TCP/IP connectivity  within a `s390x` based system- `pci`: a physical PCI device which can only be created or deleted when the bare metal  server is stopped  - Has an `allowed_vlans` property which controls the VLANs that will be permitted    to use the PCI attachment  - Cannot directly use an IEEE 802.1q VLAN tag.- `vlan`: a virtual device, used through a `pci` device that has the `vlan` in its  array of `allowed_vlans`.  - Must use an IEEE 802.1q tag.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "name"),
				Description:  "The name for this bare metal server network attachment. The name is unique across all network attachments for the bare metal server.",
			},
			"virtual_network_interface": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "A virtual network interface for the bare metal server network attachment. This can be specified using an existing virtual network interface, or a prototype object for a new virtual network interface.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The virtual network interface id for this bare metal server network attachment.",
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
						"name": &schema.Schema{
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
							Description: "The primary IP address of the virtual network interface for the bare metal server networkattachment.",
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
						"resource_type": &schema.Schema{
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
				},
			},
			"allowed_vlans": &schema.Schema{
				Type:        schema.TypeSet,
				Optional:    true,
				Set:         schema.HashInt,
				Description: "Indicates what VLAN IDs (for VLAN type only) can use this physical (PCI type) attachment.",
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"allow_to_float": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Indicates if the bare metal server network attachment can automatically float to any other server within the same `resource_group`. The bare metal server network attachment will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to bare metal server network attachments with `vlan` interface type.",
			},
			"vlan": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_network_attachment", "vlan"),
				Description:  "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this attachment.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the bare metal server network attachment was created.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this bare metal server network attachment.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the bare metal server network attachment.",
			},
			"port_speed": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The port speed for this bare metal server network attachment in Mbps.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The bare metal server network attachment type.",
			},
		},
	}
}

func ResourceIBMIsBareMetalServerNetworkAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "bare_metal_server",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "interface_type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "hipersocket, pci, vlan",
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
		validate.ValidateSchema{
			Identifier:                 "vlan",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1",
			MaxValue:                   "4094",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_bare_metal_server_network_attachment", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsBareMetalServerNetworkAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	bodyModelMap := map[string]interface{}{}
	createBareMetalServerNetworkAttachmentOptions := &vpcv1.CreateBareMetalServerNetworkAttachmentOptions{}

	bodyModelMap["interface_type"] = d.Get("interface_type")
	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	bodyModelMap["virtual_network_interface"] = d.Get("virtual_network_interface")
	if _, ok := d.GetOk("allowed_vlans"); ok {
		bodyModelMap["allowed_vlans"] = d.Get("allowed_vlans")
	}
	if _, ok := d.GetOk("allow_to_float"); ok {
		bodyModelMap["allow_to_float"] = d.Get("allow_to_float")
	}
	if _, ok := d.GetOk("vlan"); ok {
		bodyModelMap["vlan"] = d.Get("vlan")
	}
	createBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(d.Get("bare_metal_server").(string))
	convertedModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototype(bodyModelMap)
	if err != nil {
		return diag.FromErr(err)
	}
	createBareMetalServerNetworkAttachmentOptions.BareMetalServerNetworkAttachmentPrototype = convertedModel

	bareMetalServerNetworkAttachmentIntf, response, err := vpcClient.CreateBareMetalServerNetworkAttachmentWithContext(context, createBareMetalServerNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan)
		d.SetId(fmt.Sprintf("%s/%s", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID, *bareMetalServerNetworkAttachment.ID))
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci)
		d.SetId(fmt.Sprintf("%s/%s", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID, *bareMetalServerNetworkAttachment.ID))
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment)
		d.SetId(fmt.Sprintf("%s/%s", *createBareMetalServerNetworkAttachmentOptions.BareMetalServerID, *bareMetalServerNetworkAttachment.ID))
	} else {
		return diag.FromErr(fmt.Errorf("Unrecognized vpcv1.BareMetalServerNetworkAttachmentIntf subtype encountered"))
	}

	return resourceIBMIsBareMetalServerNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsBareMetalServerNetworkAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	getBareMetalServerNetworkAttachmentOptions := &vpcv1.GetBareMetalServerNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(parts[0])
	getBareMetalServerNetworkAttachmentOptions.SetID(parts[1])

	bareMetalServerNetworkAttachmentIntf, response, err := vpcClient.GetBareMetalServerNetworkAttachmentWithContext(context, getBareMetalServerNetworkAttachmentOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByVlan)
		if err = d.Set("interface_type", bareMetalServerNetworkAttachment.InterfaceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Name) {
			if err = d.Set("name", bareMetalServerNetworkAttachment.Name); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
			}
		}
		virtualNetworkInterfaceMap, err := resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface, vpcClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("virtual_network_interface", []map[string]interface{}{virtualNetworkInterfaceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowToFloat) {
			if err = d.Set("allow_to_float", bareMetalServerNetworkAttachment.AllowToFloat); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allow_to_float: %s", err))
			}
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Vlan) {
			if err = d.Set("vlan", flex.IntValue(bareMetalServerNetworkAttachment.Vlan)); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting vlan: %s", err))
			}
		}
		if err = d.Set("created_at", flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
		if err = d.Set("href", bareMetalServerNetworkAttachment.Href); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
		}
		if err = d.Set("lifecycle_state", bareMetalServerNetworkAttachment.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
		if err = d.Set("port_speed", flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
		}
		if err = d.Set("resource_type", bareMetalServerNetworkAttachment.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
		}
		if err = d.Set("type", bareMetalServerNetworkAttachment.Type); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
		}
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachmentByPci)
		if err = d.Set("interface_type", bareMetalServerNetworkAttachment.InterfaceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Name) {
			if err = d.Set("name", bareMetalServerNetworkAttachment.Name); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
			}
		}
		virtualNetworkInterfaceMap, err := resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface, vpcClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("virtual_network_interface", []map[string]interface{}{virtualNetworkInterfaceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowedVlans) {
			allowedVlans := []interface{}{}
			for _, allowedVlansItem := range bareMetalServerNetworkAttachment.AllowedVlans {
				allowedVlans = append(allowedVlans, int64(allowedVlansItem))
			}
			if err = d.Set("allowed_vlans", allowedVlans); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allowed_vlans: %s", err))
			}
		}
		if err = d.Set("created_at", flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
		if err = d.Set("href", bareMetalServerNetworkAttachment.Href); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
		}
		if err = d.Set("lifecycle_state", bareMetalServerNetworkAttachment.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
		if err = d.Set("port_speed", flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
		}
		if err = d.Set("resource_type", bareMetalServerNetworkAttachment.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
		}

		if err = d.Set("type", bareMetalServerNetworkAttachment.Type); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
		}
	} else if _, ok := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment); ok {
		bareMetalServerNetworkAttachment := bareMetalServerNetworkAttachmentIntf.(*vpcv1.BareMetalServerNetworkAttachment)
		// parent class argument: bare_metal_server string
		if err = d.Set("bare_metal_server", getBareMetalServerNetworkAttachmentOptions.BareMetalServerID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting bare_metal_server: %s", err))
		}
		// parent class argument: interface_type string
		// parent class argument: name string
		// parent class argument: virtual_network_interface VirtualNetworkInterfaceReferenceAttachmentContext
		// parent class argument: allowed_vlans []int64
		// parent class argument: allow_to_float bool
		// parent class argument: vlan int64
		if err = d.Set("interface_type", bareMetalServerNetworkAttachment.InterfaceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting interface_type: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Name) {
			if err = d.Set("name", bareMetalServerNetworkAttachment.Name); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting name: %s", err))
			}
		}
		virtualNetworkInterfaceMap, err := resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface, vpcClient)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("virtual_network_interface", []map[string]interface{}{virtualNetworkInterfaceMap}); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting virtual_network_interface: %s", err))
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowedVlans) {
			allowedVlans := []interface{}{}
			for _, allowedVlansItem := range bareMetalServerNetworkAttachment.AllowedVlans {
				allowedVlans = append(allowedVlans, int64(allowedVlansItem))
			}
			if err = d.Set("allowed_vlans", allowedVlans); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allowed_vlans: %s", err))
			}
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.AllowToFloat) {
			if err = d.Set("allow_to_float", bareMetalServerNetworkAttachment.AllowToFloat); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting allow_to_float: %s", err))
			}
		}
		if !core.IsNil(bareMetalServerNetworkAttachment.Vlan) {
			if err = d.Set("vlan", flex.IntValue(bareMetalServerNetworkAttachment.Vlan)); err != nil {
				return diag.FromErr(fmt.Errorf("[ERROR] Error setting vlan: %s", err))
			}
		}
		if err = d.Set("created_at", flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting created_at: %s", err))
		}
		if err = d.Set("href", bareMetalServerNetworkAttachment.Href); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting href: %s", err))
		}
		if err = d.Set("lifecycle_state", bareMetalServerNetworkAttachment.LifecycleState); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting lifecycle_state: %s", err))
		}
		if err = d.Set("port_speed", flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting port_speed: %s", err))
		}
		if err = d.Set("resource_type", bareMetalServerNetworkAttachment.ResourceType); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting resource_type: %s", err))
		}
		if err = d.Set("type", bareMetalServerNetworkAttachment.Type); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting type: %s", err))
		}
		if err = d.Set("network_attachment", bareMetalServerNetworkAttachment.ID); err != nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Error setting network_attachment: %s", err))
		}
	} else {
		return diag.FromErr(fmt.Errorf("Unrecognized vpcv1.BareMetalServerNetworkAttachmentIntf subtype encountered"))
	}

	return nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	updateBareMetalServerNetworkAttachmentOptions := &vpcv1.UpdateBareMetalServerNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(parts[0])
	updateBareMetalServerNetworkAttachmentOptions.SetID(parts[1])

	hasChange := false

	patchVals := &vpcv1.BareMetalServerNetworkAttachmentPatch{}
	if d.HasChange("bare_metal_server") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "bare_metal_server"))
	}
	if d.HasChange("allowed_vlans") {
		var allowedVlans []int64
		for _, v := range d.Get("allowed_vlans").(*schema.Set).List() {
			allowedVlansItem := int64(v.(int))
			allowedVlans = append(allowedVlans, allowedVlansItem)
		}
		patchVals.AllowedVlans = allowedVlans
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}

	if hasChange {
		updateBareMetalServerNetworkAttachmentOptions.BareMetalServerNetworkAttachmentPatch, _ = patchVals.AsPatch()
		_, response, err := vpcClient.UpdateBareMetalServerNetworkAttachmentWithContext(context, updateBareMetalServerNetworkAttachmentOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIBMIsBareMetalServerNetworkAttachmentRead(context, d, meta)
}

func resourceIBMIsBareMetalServerNetworkAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBareMetalServerNetworkAttachmentOptions := &vpcv1.DeleteBareMetalServerNetworkAttachmentOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteBareMetalServerNetworkAttachmentOptions.SetBareMetalServerID(parts[0])
	deleteBareMetalServerNetworkAttachmentOptions.SetID(parts[1])

	response, err := vpcClient.DeleteBareMetalServerNetworkAttachmentWithContext(context, deleteBareMetalServerNetworkAttachmentOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteBareMetalServerNetworkAttachmentWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap map[string]interface{}) (vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceIntf, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface{}
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
			ipsItemModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(ipsItem.(map[string]interface{}))
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
		PrimaryIPModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && len(modelMap["resource_group"].([]interface{})) > 0 {
		ResourceGroupModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToResourceGroupIdentity(modelMap["resource_group"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceGroup = ResourceGroupModel
	}
	if modelMap["security_groups"] != nil {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		for _, securityGroupsItem := range modelMap["security_groups"].([]interface{}) {
			securityGroupsItemModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToSecurityGroupIdentity(securityGroupsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			securityGroups = append(securityGroups, securityGroupsItemModel)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeIntf, error) {
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

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContext(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextIntf, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref, error) {
	model := &vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPIdentityVirtualNetworkInterfaceIPsContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfaceIPPrototypeReservedIPPrototypeVirtualNetworkInterfaceIPsContext, error) {
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

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeIntf, error) {
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

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextIntf, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContext{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref, error) {
	model := &vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPIdentityVirtualNetworkInterfacePrimaryIPContextByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext(modelMap map[string]interface{}) (*vpcv1.VirtualNetworkInterfacePrimaryIPPrototypeReservedIPPrototypeVirtualNetworkInterfacePrimaryIPContext, error) {
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

func resourceIBMIsBareMetalServerNetworkAttachmentMapToResourceGroupIdentity(modelMap map[string]interface{}) (vpcv1.ResourceGroupIdentityIntf, error) {
	model := &vpcv1.ResourceGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ResourceGroupIdentityByID, error) {
	model := &vpcv1.ResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSecurityGroupIdentity(modelMap map[string]interface{}) (vpcv1.SecurityGroupIdentityIntf, error) {
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

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSecurityGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByID, error) {
	model := &vpcv1.SecurityGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSecurityGroupIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByCRN, error) {
	model := &vpcv1.SecurityGroupIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSecurityGroupIdentityByHref(modelMap map[string]interface{}) (*vpcv1.SecurityGroupIdentityByHref, error) {
	model := &vpcv1.SecurityGroupIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSubnetIdentity(modelMap map[string]interface{}) (vpcv1.SubnetIdentityIntf, error) {
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

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSubnetIdentityByID(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByID, error) {
	model := &vpcv1.SubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSubnetIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByCRN, error) {
	model := &vpcv1.SubnetIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToSubnetIdentityByHref(modelMap map[string]interface{}) (*vpcv1.SubnetIdentityByHref, error) {
	model := &vpcv1.SubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeBareMetalServerNetworkAttachmentContext(modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeBareMetalServerNetworkAttachmentContext, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfacePrototypeBareMetalServerNetworkAttachmentContext{}
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
			ipsItemModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfaceIPPrototype(ipsItem.(map[string]interface{}))
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
		PrimaryIPModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToVirtualNetworkInterfacePrimaryIPPrototype(modelMap["primary_ip"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.PrimaryIP = PrimaryIPModel
	}
	if modelMap["resource_group"] != nil && len(modelMap["resource_group"].([]interface{})) > 0 {
		ResourceGroupModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToResourceGroupIdentity(modelMap["resource_group"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceGroup = ResourceGroupModel
	}
	if modelMap["security_groups"] != nil {
		securityGroups := []vpcv1.SecurityGroupIdentityIntf{}
		for _, securityGroupsItem := range modelMap["security_groups"].([]interface{}) {
			securityGroupsItemModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToSecurityGroupIdentity(securityGroupsItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			securityGroups = append(securityGroups, securityGroupsItemModel)
		}
		model.SecurityGroups = securityGroups
	}
	if modelMap["subnet"] != nil && len(modelMap["subnet"].([]interface{})) > 0 {
		SubnetModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToSubnetIdentity(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Subnet = SubnetModel
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentity(modelMap map[string]interface{}) (vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityIntf, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentity{}
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

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByID(modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByID, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByHref(modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByHref, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterfaceVirtualNetworkInterfaceIdentityVirtualNetworkInterfaceIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototype(modelMap map[string]interface{}) (vpcv1.BareMetalServerNetworkAttachmentPrototypeIntf, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototype{}
	model.InterfaceType = core.StringPtr(modelMap["interface_type"].(string))
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	if modelMap["allowed_vlans"] != nil {
		allowedVlans := []int64{}
		for _, allowedVlansItem := range modelMap["allowed_vlans"].(*schema.Set).List() {
			allowedVlans = append(allowedVlans, int64(allowedVlansItem.(int)))
		}
		model.AllowedVlans = allowedVlans
	}
	if modelMap["allow_to_float"] != nil {
		model.AllowToFloat = core.BoolPtr(modelMap["allow_to_float"].(bool))
	}
	if modelMap["vlan"] != nil {
		model.Vlan = core.Int64Ptr(int64(modelMap["vlan"].(int)))
	}
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype(modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByPciPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	if modelMap["allowed_vlans"] != nil {
		allowedVlans := []int64{}
		for _, allowedVlansItem := range modelMap["allowed_vlans"].([]interface{}) {
			allowedVlans = append(allowedVlans, int64(allowedVlansItem.(int)))
		}
		model.AllowedVlans = allowedVlans
	}
	model.InterfaceType = core.StringPtr(modelMap["interface_type"].(string))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype(modelMap map[string]interface{}) (*vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype, error) {
	model := &vpcv1.BareMetalServerNetworkAttachmentPrototypeBareMetalServerNetworkAttachmentByVlanPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VirtualNetworkInterfaceModel, err := resourceIBMIsBareMetalServerNetworkAttachmentMapToBareMetalServerNetworkAttachmentPrototypeVirtualNetworkInterface(modelMap["virtual_network_interface"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.VirtualNetworkInterface = VirtualNetworkInterfaceModel
	if modelMap["allow_to_float"] != nil {
		model.AllowToFloat = core.BoolPtr(modelMap["allow_to_float"].(bool))
	}
	model.InterfaceType = core.StringPtr(modelMap["interface_type"].(string))
	model.Vlan = core.Int64Ptr(int64(modelMap["vlan"].(int)))
	return model, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(model *vpcv1.VirtualNetworkInterfaceReferenceAttachmentContext, sess *vpcv1.VpcV1) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	vniid := *model.ID
	getVirtualNetworkInterfaceOptions := &vpcv1.GetVirtualNetworkInterfaceOptions{
		ID: &vniid,
	}
	vniDetails, response, err := sess.GetVirtualNetworkInterface(getVirtualNetworkInterfaceOptions)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Error on GetInstanceNetworkAttachment in instance : %s\n%s", err, response)
	}
	modelMap["allow_ip_spoofing"] = vniDetails.AllowIPSpoofing
	modelMap["auto_delete"] = vniDetails.AutoDelete
	modelMap["enable_infrastructure_nat"] = vniDetails.EnableInfrastructureNat
	modelMap["resource_group"] = vniDetails.ResourceGroup.ID
	primaryipId := *vniDetails.PrimaryIP.ID
	if !core.IsNil(vniDetails.Ips) {
		ips := []map[string]interface{}{}
		for _, ipsItem := range vniDetails.Ips {
			if *ipsItem.ID != primaryipId {
				ipsItemMap, err := resourceIBMIsVirtualNetworkInterfaceReservedIPReferenceToMap(&ipsItem, true)
				if err != nil {
					return nil, err
				}
				ips = append(ips, ipsItemMap)
			}
		}
		modelMap["ips"] = ips
	}
	primaryIPMap, err := resourceIBMIsBareMetalServerReservedIPReferenceToMap(vniDetails.PrimaryIP)
	if err != nil {
		return modelMap, err
	}
	modelMap["primary_ip"] = []map[string]interface{}{primaryIPMap}

	if !core.IsNil(vniDetails.SecurityGroups) {
		securityGroups := make([]string, 0)
		for _, securityGroupsItem := range vniDetails.SecurityGroups {
			if securityGroupsItem.ID != nil {
				securityGroups = append(securityGroups, *securityGroupsItem.ID)
			}
		}
		modelMap["security_groups"] = securityGroups
	}
	if model.ResourceType != nil {
		modelMap["resource_type"] = *model.ResourceType
	}
	if vniDetails.Subnet != nil {
		modelMap["subnet"] = *vniDetails.Subnet.ID
	}
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceToMap(model *vpcv1.ReservedIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceDeletedToMap(model.Deleted)
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

func resourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceDeletedToMap(model *vpcv1.ReservedIPReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.CRN
	if model.Deleted != nil {
		deletedMap, err := resourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceDeletedToMap(model.Deleted)
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

func resourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceDeletedToMap(model *vpcv1.SubnetReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
