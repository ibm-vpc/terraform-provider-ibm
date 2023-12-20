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
	"github.com/IBM/vpc-beta-go-sdk/ontapv1"
)

func DataSourceIbmIsStorageOntapInstanceVirtualMachineVolumes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesRead,

		Schema: map[string]*schema.Schema{
			"storage_ontap_instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The storage ontap instance identifier.",
			},
			"storage_virtual_machine_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The storage virtual machine identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `name` property matching the exact specified name.",
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
			"volumes": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of storage volumes.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The capacity of the storage volume (in gigabytes).",
						},
						"cifs_share": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The named access point that enables CIFS clients to view, browse, and manipulatefiles on this storage volumeThis will be present when `security_style` is `mixed` or `windows`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_control_list": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The access control list for the CIFS share.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"permission": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The permission granted to users matching this access control list entry.",
												},
												"users": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The users matching this access control list entry.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
									"mount_path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The SMB/CIFS mount point for the storage volume.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The share name registered in Active Directory that SMB/CIFS clients use to mount the share. The name is unique within the Active Directory domain.",
									},
								},
							},
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the storage volume was created.",
						},
						"enable_storage_efficiency": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Deprecated:  "This argument is deprecated and may be removed in a future release",
							Description: "Indicates whether storage efficiency is enabled for the storage volume.If `true`, data-deduplication, compression and other efficiencies for space-management are enabled for this volume.",
						},
						"export_policy": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The NFS export policy for the storage volume.This will be present when `security_style` is `mixed` or `unix`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"mount_path": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The NFS mount point for the storage volume.",
									},
									"rules": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The NFS export policy rules for this storage volume.Only NFS clients included in the rules will access the volume, and only according to the specified access controls and NFS protocol versions.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"access_control": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The access control that is provided to clients that match this rule.",
												},
												"clients": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The clients that match this rule. Every client in the list of `clients` is unique.",
													Elem: &schema.Resource{
														Schema: map[string]*schema.Schema{
															"hostname": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The hostname of the NFS client.",
															},
															"address": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The IP address of the NFS client.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
															},
															"cidr_block": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The CIDR block containing IP addresses of the NFS clients. The CIDR block `0.0.0.0/0` matches all client addresses.This property may add support for IPv6 CIDR blocks in the future. When processing a value in this property, verify that the CIDR block is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected CIDR block format was encountered.",
															},
															"domain_suffix": &schema.Schema{
																Type:        schema.TypeString,
																Computed:    true,
																Description: "The domain names suffixes of the NFS clients.",
															},
														},
													},
												},
												"index": &schema.Schema{
													Type:        schema.TypeInt,
													Computed:    true,
													Description: "The unique index for this rule. Rules are applied from lowest to highest index.",
												},
												"is_superuser": &schema.Schema{
													Type:        schema.TypeBool,
													Computed:    true,
													Description: "If `true`, clients matching this rule that request super-user access are honored. Otherwise, clients are mapped to the anonymous user.",
												},
												"nfs_version": &schema.Schema{
													Type:        schema.TypeList,
													Computed:    true,
													Description: "The NFS versions that is provided to clients that match this rule.",
													Elem: &schema.Schema{
														Type: schema.TypeString,
													},
												},
											},
										},
									},
								},
							},
						},
						"health_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current storage volume health_state (if any):- `primary_node_down`: The storage volume is experiencing higher latency due to  the primary node being unavailable, and I/O being routed to the secondary node.- `volume_unavailable`: The storage volume is unavailable as both the primary and secondary nodes are down.- `internal_error`: Internal error (contact IBM support).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this health state.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this health state.",
									},
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this health state.",
									},
								},
							},
						},
						"health_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The health of this resource.- `ok`: No abnormal behavior detected- `degraded`: Experiencing compromised performance, capacity, or connectivity- `faulted`: Completely unreachable, inoperative, or otherwise entirely incapacitated- `inapplicable`: The health state does not apply because of the current lifecycle state. A resource with a lifecycle state of `failed` or `deleting` will have a health state of `inapplicable`. A `pending` resource may also have this state.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this storage volume.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this storage volume.",
						},
						"junction_path": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The path clients can use to mount or access this storage volume. The path is case insensitive and is unique within a storage virtual machine.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the storage volume.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this storage volume. The name is unique across all storage volumes in the storage virtual machine.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"security_style": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The security style for the storage volume:- `unix`: NFS clients can access the storage volume.- `windows`: SMB/CIFS clients can access the storage volume.- `mixed`: Both SMB/CIFS and NFS clients can access the storage volume.- `none`: No clients can access the volume.",
						},
						"storage_efficiency": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The storage efficiency mode used for this storage volume.- `disabled`: storage efficiency methods will not be used- `enabled`: data-deduplication, compression and other methods will be used.",
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the storage volume.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	listStorageOntapInstanceStorageVirtualMachineVolumesOptions := &ontapv1.ListStorageOntapInstanceStorageVirtualMachineVolumesOptions{}

	listStorageOntapInstanceStorageVirtualMachineVolumesOptions.SetStorageOntapInstanceID(d.Get("storage_ontap_instance_id").(string))
	listStorageOntapInstanceStorageVirtualMachineVolumesOptions.SetStorageVirtualMachineID(d.Get("storage_virtual_machine_id").(string))
	if _, ok := d.GetOk("name"); ok {
		listStorageOntapInstanceStorageVirtualMachineVolumesOptions.SetName(d.Get("name").(string))
	}

	var pager *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumesPager
	pager, err = ontapClient.NewStorageOntapInstanceStorageVirtualMachineVolumesPager(listStorageOntapInstanceStorageVirtualMachineVolumesOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] StorageOntapInstanceStorageVirtualMachineVolumesPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("StorageOntapInstanceStorageVirtualMachineVolumesPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("volumes", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting volumes %s", err))
	}

	return nil
}

// dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesID returns a reasonable ID for the list.
func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeCollectionFirstToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeCollectionNextToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolume) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["capacity"] = flex.IntValue(model.Capacity)
	if model.CifsShare != nil {
		cifsShareMap, err := dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareToMap(model.CifsShare)
		if err != nil {
			return modelMap, err
		}
		modelMap["cifs_share"] = []map[string]interface{}{cifsShareMap}
	}
	modelMap["created_at"] = model.CreatedAt.String()
	if model.EnableStorageEfficiency != nil {
		modelMap["enable_storage_efficiency"] = model.EnableStorageEfficiency
	}
	if model.ExportPolicy != nil {
		exportPolicyMap, err := dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyToMap(model.ExportPolicy)
		if err != nil {
			return modelMap, err
		}
		modelMap["export_policy"] = []map[string]interface{}{exportPolicyMap}
	}
	healthReasons := []map[string]interface{}{}
	for _, healthReasonsItem := range model.HealthReasons {
		healthReasonsItemMap, err := dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeHealthReasonToMap(&healthReasonsItem)
		if err != nil {
			return modelMap, err
		}
		healthReasons = append(healthReasons, healthReasonsItemMap)
	}
	modelMap["health_reasons"] = healthReasons
	modelMap["health_state"] = model.HealthState
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	modelMap["junction_path"] = model.JunctionPath
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	modelMap["security_style"] = model.SecurityStyle
	if model.StorageEfficiency != nil {
		modelMap["storage_efficiency"] = model.StorageEfficiency
	}
	modelMap["type"] = model.Type
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShare) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	accessControlList := []map[string]interface{}{}
	for _, accessControlListItem := range model.AccessControlList {
		accessControlListItemMap, err := dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlToMap(&accessControlListItem)
		if err != nil {
			return modelMap, err
		}
		accessControlList = append(accessControlList, accessControlListItemMap)
	}
	modelMap["access_control_list"] = accessControlList
	modelMap["mount_path"] = model.MountPath
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControl) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["permission"] = model.Permission
	modelMap["users"] = model.Users
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = model.MountPath
	if model.Rules != nil {
		rules := []map[string]interface{}{}
		for _, rulesItem := range model.Rules {
			rulesItemMap, err := dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleToMap(&rulesItem)
			if err != nil {
				return modelMap, err
			}
			rules = append(rules, rulesItemMap)
		}
		modelMap["rules"] = rules
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["access_control"] = model.AccessControl
	clients := []map[string]interface{}{}
	for _, clientsItem := range model.Clients {
		clientsItemMap, err := dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientToMap(clientsItem)
		if err != nil {
			return modelMap, err
		}
		clients = append(clients, clientsItemMap)
	}
	modelMap["clients"] = clients
	modelMap["index"] = flex.IntValue(model.Index)
	modelMap["is_superuser"] = model.IsSuperuser
	modelMap["nfs_version"] = model.NfsVersion
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientToMap(model ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname); ok {
		return dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostnameToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname))
	} else if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP); ok {
		return dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIPToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP))
	} else if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR); ok {
		return dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDRToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR))
	} else if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName); ok {
		return dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainNameToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName))
	} else if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClient); ok {
		modelMap := make(map[string]interface{})
		model := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClient)
		if model.Hostname != nil {
			modelMap["hostname"] = model.Hostname
		}
		if model.Address != nil {
			modelMap["address"] = model.Address
		}
		if model.CidrBlock != nil {
			modelMap["cidr_block"] = model.CidrBlock
		}
		if model.DomainSuffix != nil {
			modelMap["domain_suffix"] = model.DomainSuffix
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientIntf subtype encountered")
	}
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostnameToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hostname"] = model.Hostname
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIPToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDRToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["cidr_block"] = model.CidrBlock
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainNameToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["domain_suffix"] = model.DomainSuffix
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVirtualMachineVolumesStorageOntapInstanceStorageVirtualMachineVolumeHealthReasonToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeHealthReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}
