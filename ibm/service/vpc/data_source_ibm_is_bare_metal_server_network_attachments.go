// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsBareMetalServerNetworkAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBareMetalServerNetworkAttachmentsRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},

			//network attachments properties
			"network_attachments": {
				Type:        schema.TypeList,
				Description: "A list of all network attachments on a bare metal server. A network interface is an abstract representation of a network interface card and connects a bare metal server to a subnet. While each network interface can attach to only one subnet, multiple network interfaces can be created to attach to multiple subnets. Multiple interfaces may also attach to the same subnet.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bare_metal_server_id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The bare metal server identifier.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The network attachment identifier.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the network attachment was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this network attachment.",
						},
						"bare_metal_server_network_attachment_id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this network attachment.",
						},
						"interface_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The network attachment's interface type:- `hipersocket`: a virtual network device that provides high-speed TCP/IP connectivity  within a `s390x` based system- `pci`: a physical PCI device which can only be created or deleted when the bare metal  server is stopped  - Has an `allowed_vlans` property which controls the VLANs that will be permitted    to use the PCI attachment  - Cannot directly use an IEEE 802.1q VLAN tag.- `vlan`: a virtual device, used through a `pci` device that has the `vlan` in its  array of `allowed_vlans`.  - Must use an IEEE 802.1q tag.The enumerated values for this property are expected to expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected property value was encountered.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the network attachment.",
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_speed": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The port speed for this network attachment in Mbps.",
						},
						"primary_ip": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The primary IP address of the virtual network interface for the network attachment.",
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
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"subnet": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnet of the virtual network interface for the network attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this subnet.",
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
										Description: "The URL for this subnet.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this subnet.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this subnet. The name is unique across all subnets in the VPC.",
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
							Description: "This bare metal server network attachment's interface type.",
						},
						"virtual_network_interface": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The virtual network interface for this network attachment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this virtual network interface.",
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
										Description: "The URL for this virtual network interface.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this virtual network interface.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this virtual network interface. The name is unique across all virtual network interfaces in the VPC.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"allowed_vlans": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"allow_to_float": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates if the network attachment can automatically float to any other server within the same `resource_group`. The network attachment will float automatically if the network detects a GARP or RARP on another bare metal server in the resource group. Applies only to network attachments with `vlan` interface type.",
						},
						"vlan": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Indicates the 802.1Q VLAN ID tag that must be used for all traffic on this attachment.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsBareMetalServerNetworkAttachmentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listBareMetalServerNetworkAttachmentsOptions := &vpcv1.ListBareMetalServerNetworkAttachmentsOptions{}

	listBareMetalServerNetworkAttachmentsOptions.SetBareMetalServerID(d.Get("bare_metal_server_id").(string))
	nattchs := []vpcv1.BareMetalServerNetworkAttachmentIntf{}

	bareMetalServerNetworkAttachmentsIntf, response, err := sess.ListBareMetalServerNetworkAttachmentsWithContext(context, listBareMetalServerNetworkAttachmentsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListBareMetalServerNetworkAttachmentsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListBareMetalServerNetworkAttachmentsWithContext failed %s\n%s", err, response))
	}
	nattchs = append(nattchs, bareMetalServerNetworkAttachmentsIntf.NetworkAttachments...)
	nattchInfo := make([]map[string]interface{}, 0)
	for _, nattch := range nattchs {
		l := map[string]interface{}{}
		switch reflect.TypeOf(nattch).String() {
		case "*vpcv1.BareMetalServerNetworkAttachmentByHiperSocket":
			{
				bareMetalServerNetworkAttachment := nattch.(*vpcv1.BareMetalServerNetworkAttachmentByHiperSocket)
				l["created_at"] = flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)

				l["href"] = bareMetalServerNetworkAttachment.Href

				l["bare_metal_server_network_attachment_id"] = bareMetalServerNetworkAttachment.ID

				l["interface_type"] = bareMetalServerNetworkAttachment.InterfaceType

				l["lifecycle_state"] = bareMetalServerNetworkAttachment.LifecycleState

				l["name"] = bareMetalServerNetworkAttachment.Name

				l["port_speed"] = flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)

				primaryIP := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.PrimaryIP != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceToMap(bareMetalServerNetworkAttachment.PrimaryIP)
					if err != nil {
						return diag.FromErr(err)
					}
					primaryIP = append(primaryIP, modelMap)
				}
				l["primary_ip"] = primaryIP

				l["resource_type"] = bareMetalServerNetworkAttachment.ResourceType

				subnet := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.Subnet != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceToMap(bareMetalServerNetworkAttachment.Subnet)
					if err != nil {
						return diag.FromErr(err)
					}
					subnet = append(subnet, modelMap)
				}
				l["subnet"] = subnet

				l["type"] = bareMetalServerNetworkAttachment.Type
				virtualNetworkInterface := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.VirtualNetworkInterface != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface)
					if err != nil {
						return diag.FromErr(err)
					}
					virtualNetworkInterface = append(virtualNetworkInterface, modelMap)
				}
				l["virtual_network_interface"] = virtualNetworkInterface
				nattchInfo = append(nattchInfo, l)
			}
		case "*vpcv1.BareMetalServerNetworkAttachmentByPci":
			{
				bareMetalServerNetworkAttachment := nattch.(*vpcv1.BareMetalServerNetworkAttachmentByPci)
				l["created_at"] = flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)

				l["href"] = bareMetalServerNetworkAttachment.Href

				l["bare_metal_server_network_attachment_id"] = bareMetalServerNetworkAttachment.ID

				l["interface_type"] = bareMetalServerNetworkAttachment.InterfaceType

				l["lifecycle_state"] = bareMetalServerNetworkAttachment.LifecycleState

				l["name"] = bareMetalServerNetworkAttachment.Name

				l["port_speed"] = flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)

				primaryIP := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.PrimaryIP != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceToMap(bareMetalServerNetworkAttachment.PrimaryIP)
					if err != nil {
						return diag.FromErr(err)
					}
					primaryIP = append(primaryIP, modelMap)
				}
				l["primary_ip"] = primaryIP

				l["resource_type"] = bareMetalServerNetworkAttachment.ResourceType

				subnet := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.Subnet != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceToMap(bareMetalServerNetworkAttachment.Subnet)
					if err != nil {
						return diag.FromErr(err)
					}
					subnet = append(subnet, modelMap)
				}
				l["subnet"] = subnet

				l["type"] = bareMetalServerNetworkAttachment.Type
				virtualNetworkInterface := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.VirtualNetworkInterface != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface)
					if err != nil {
						return diag.FromErr(err)
					}
					virtualNetworkInterface = append(virtualNetworkInterface, modelMap)
				}
				l["virtual_network_interface"] = virtualNetworkInterface
				if bareMetalServerNetworkAttachment.AllowedVlans != nil {
					var out = make([]interface{}, len(bareMetalServerNetworkAttachment.AllowedVlans), len(bareMetalServerNetworkAttachment.AllowedVlans))
					for i, v := range bareMetalServerNetworkAttachment.AllowedVlans {
						out[i] = int(v)
					}
					l["allowed_vlans"] = schema.NewSet(schema.HashInt, out)
				}
				nattchInfo = append(nattchInfo, l)
			}
		case "*vpcv1.BareMetalServerNetworkAttachmentByVlan":
			{
				bareMetalServerNetworkAttachment := nattch.(*vpcv1.BareMetalServerNetworkAttachmentByVlan)
				l["created_at"] = flex.DateTimeToString(bareMetalServerNetworkAttachment.CreatedAt)

				l["href"] = bareMetalServerNetworkAttachment.Href

				l["bare_metal_server_network_attachment_id"] = bareMetalServerNetworkAttachment.ID

				l["interface_type"] = bareMetalServerNetworkAttachment.InterfaceType

				l["lifecycle_state"] = bareMetalServerNetworkAttachment.LifecycleState

				l["name"] = bareMetalServerNetworkAttachment.Name

				l["port_speed"] = flex.IntValue(bareMetalServerNetworkAttachment.PortSpeed)

				primaryIP := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.PrimaryIP != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentReservedIPReferenceToMap(bareMetalServerNetworkAttachment.PrimaryIP)
					if err != nil {
						return diag.FromErr(err)
					}
					primaryIP = append(primaryIP, modelMap)
				}
				l["primary_ip"] = primaryIP

				l["resource_type"] = bareMetalServerNetworkAttachment.ResourceType

				subnet := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.Subnet != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentSubnetReferenceToMap(bareMetalServerNetworkAttachment.Subnet)
					if err != nil {
						return diag.FromErr(err)
					}
					subnet = append(subnet, modelMap)
				}
				l["subnet"] = subnet

				l["type"] = bareMetalServerNetworkAttachment.Type
				virtualNetworkInterface := []map[string]interface{}{}
				if bareMetalServerNetworkAttachment.VirtualNetworkInterface != nil {
					modelMap, err := dataSourceIBMIsBareMetalServerNetworkAttachmentVirtualNetworkInterfaceReferenceAttachmentContextToMap(bareMetalServerNetworkAttachment.VirtualNetworkInterface)
					if err != nil {
						return diag.FromErr(err)
					}
					virtualNetworkInterface = append(virtualNetworkInterface, modelMap)
				}
				l["virtual_network_interface"] = virtualNetworkInterface

				l["allow_to_float"] = bareMetalServerNetworkAttachment.AllowToFloat

				l["vlan"] = flex.IntValue(bareMetalServerNetworkAttachment.Vlan)
				nattchInfo = append(nattchInfo, l)
			}
		}
	}
	d.SetId(dataSourceIBMIsBareMetalServerNetworkAttachmentsId(d))
	d.Set("network_attachments", nattchInfo)
	return nil
}

// dataSourceIBMIsBareMetalServerNetworkAttachmentsId returns a reasonable ID for a BMS network attachments.
func dataSourceIBMIsBareMetalServerNetworkAttachmentsId(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
