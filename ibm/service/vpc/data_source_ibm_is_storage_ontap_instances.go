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

func DataSourceIbmIsStorageOntapInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsStorageOntapInstancesRead,

		Schema: map[string]*schema.Schema{
			"resource_group_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with a `resource_group.id` property matching the specified identifier.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to storage ontap instances with a `lifecycle_state` property matching the specified value.",
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
			"storage_ontap_instances": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of storage ontap instances.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address_prefix": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "An address prefix in the VPC which will be used to allocate `endpoints` for thisstorage ontap instance and its storage virtual machines.",
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
										Description: "The URL for this address prefix.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this address prefix.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this address prefix. The name is unique across all address prefixes for the VPC.",
									},
								},
							},
						},
						"admin_credentials": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The credentials used (from Secrets Manager) for the cluster administrator to access thestorage ontap instance. At least one of `password`, `ssh`, or `http` will be present.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The security certificate credential for ONTAP REST API access for the clusteradministrator of the storage ontap instance.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this credential.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"password": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The password credential for the cluster administrator of the storage ontap instance.If present, this password credential is used by the cluster administrator for bothONTAP CLI SSH access and ONTAP REST API access. If absent, the storage ontapinstance is not accessible through either the ONTAP CLI or ONTAP REST API usingpassword-based authentication.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this credential.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"ssh": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The public key credential for ONTAP CLI SSH access for the cluster administratorof the storage ontap instance.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this credential.",
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
						"capacity": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The capacity to use for the storage ontap instance (in terabytes). Volumes in this storage ontap instance will be allocated from this capacity.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the storage ontap instance was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this storage ontap instance.",
						},
						"encryption": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of encryption used on the storage ontap instance.",
						},
						"encryption_key": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The root key used to wrap the data encryption key for the storage ontap instance.This property will be present for storage ontap instance with an `encryption` type of`user_managed`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Services Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
									},
								},
							},
						},
						"endpoints": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The endpoints for this storage ontap instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"inter_cluster": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The NetApp SnapMirror management endpoint for this storage ontap instance.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4_address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique IP address of an endpoint.",
												},
											},
										},
									},
									"management": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The NetApp management endpoint for this storage ontap instance. Management may beperformed using the ONTAP CLI, ONTAP API, or NetApp CloudManager.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"ipv4_address": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique IP address of an endpoint.",
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
							Description: "The reasons for the current storage ontap instance health_state (if any):- `cluster_down`: This storage ontap instance is unavailable as both the  primary and secondary nodes are unavailable.- `failback_unavailable`: The capability to failback is unavailable. The secondary node  continues to be available.- `failover_unavailable`: The capability to failover is unavailable. The primary  node continues to be available without any performance impact to clients.- `internal_error`: Internal error (contact IBM support).- `maintenance_in_progress`: A planned maintenance activity is in progress.- `primary_node_down`: The primary node is unavailable, and I/O has failed over to  the secondary node. Clients running in the same zone as the primary node may  experience higher access latency.- `secondary_node_down`: The secondary node is unavailable. Therefore, the capability  to failover is unavailable.",
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
							Description: "The URL for this storage ontap instance.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this storage ontap instance.",
						},
						"lifecycle_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current `lifecycle_state` (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the reason for this lifecycle state.",
									},
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about the reason for this lifecycle state.",
									},
								},
							},
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the storage ontap instance.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this storage ontap instance. The name is unique across all storage ontap instances in the region.",
						},
						"primary_subnet": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnet where the primary Cloud Volumes ONTAP node is provisioned in.",
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
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this storage ontap instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
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
						"routing_tables": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC routing tables for this storage ontap instance.",
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
										Description: "The URL for this routing table.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this routing table.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this routing table. The name is unique across all routing tables for the VPC.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"secondary_subnet": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The subnet where the secondary Cloud Volumes ONTAP node is provisioned in.",
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
						"security_groups": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The security groups for this storage ontap instance.",
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
						"storage_virtual_machines": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The storage virtual machines for this storage ontap instance.",
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
										Description: "The URL for this storage virtual machine.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this storage virtual machine.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this storage virtual machine. The name is unique across all storage virtual machines in the storage ontap instance.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
								},
							},
						},
						"vpc": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The VPC this storage ontap instance resides in.",
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
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsStorageOntapInstancesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	listStorageOntapInstancesOptions := &ontapv1.ListStorageOntapInstancesOptions{}

	if _, ok := d.GetOk("resource_group_id"); ok {
		listStorageOntapInstancesOptions.SetResourceGroupID(d.Get("resource_group_id").(string))
	}
	if _, ok := d.GetOk("lifecycle_state"); ok {
		listStorageOntapInstancesOptions.SetLifecycleState(d.Get("lifecycle_state").(string))
	}

	var pager *ontapv1.StorageOntapInstancesPager
	pager, err = ontapClient.NewStorageOntapInstancesPager(listStorageOntapInstancesOptions)
	if err != nil {
		return diag.FromErr(err)
	}

	allItems, err := pager.GetAll()
	if err != nil {
		log.Printf("[DEBUG] StorageOntapInstancesPager.GetAll() failed %s", err)
		return diag.FromErr(fmt.Errorf("StorageOntapInstancesPager.GetAll() failed %s", err))
	}

	d.SetId(dataSourceIbmIsStorageOntapInstancesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceToMap(&modelItem)
		if err != nil {
			return diag.FromErr(err)
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("storage_ontap_instances", mapSlice); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting storage_ontap_instances %s", err))
	}

	return nil
}

// dataSourceIbmIsStorageOntapInstancesID returns a reasonable ID for the list.
func dataSourceIbmIsStorageOntapInstancesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceCollectionFirstToMap(model *ontapv1.StorageOntapInstanceCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceCollectionNextToMap(model *ontapv1.StorageOntapInstanceCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = model.Href
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceToMap(model *ontapv1.StorageOntapInstance) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	addressPrefixMap, err := dataSourceIbmIsStorageOntapInstancesAddressPrefixReferenceToMap(model.AddressPrefix)
	if err != nil {
		return modelMap, err
	}
	modelMap["address_prefix"] = []map[string]interface{}{addressPrefixMap}
	if model.AdminCredentials != nil {
		adminCredentialsMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceAdminCredentialsToMap(model.AdminCredentials)
		if err != nil {
			return modelMap, err
		}
		modelMap["admin_credentials"] = []map[string]interface{}{adminCredentialsMap}
	}
	modelMap["capacity"] = flex.IntValue(model.Capacity)
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["crn"] = model.Crn
	modelMap["encryption"] = model.Encryption
	if model.EncryptionKey != nil {
		encryptionKeyMap, err := dataSourceIbmIsStorageOntapInstancesEncryptionKeyReferenceToMap(model.EncryptionKey)
		if err != nil {
			return modelMap, err
		}
		modelMap["encryption_key"] = []map[string]interface{}{encryptionKeyMap}
	}
	endpointsMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceEndpointsToMap(model.Endpoints)
	if err != nil {
		return modelMap, err
	}
	modelMap["endpoints"] = []map[string]interface{}{endpointsMap}
	healthReasons := []map[string]interface{}{}
	for _, healthReasonsItem := range model.HealthReasons {
		healthReasonsItemMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceHealthReasonToMap(&healthReasonsItem)
		if err != nil {
			return modelMap, err
		}
		healthReasons = append(healthReasons, healthReasonsItemMap)
	}
	modelMap["health_reasons"] = healthReasons
	modelMap["health_state"] = model.HealthState
	modelMap["href"] = model.Href
	modelMap["id"] = model.ID
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range model.LifecycleReasons {
		lifecycleReasonsItemMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceLifecycleReasonToMap(&lifecycleReasonsItem)
		if err != nil {
			return modelMap, err
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	modelMap["lifecycle_reasons"] = lifecycleReasons
	modelMap["lifecycle_state"] = model.LifecycleState
	modelMap["name"] = model.Name
	if model.PrimarySubnet != nil {
		primarySubnetMap, err := dataSourceIbmIsStorageOntapInstancesSubnetReferenceToMap(model.PrimarySubnet)
		if err != nil {
			return modelMap, err
		}
		modelMap["primary_subnet"] = []map[string]interface{}{primarySubnetMap}
	}
	resourceGroupMap, err := dataSourceIbmIsStorageOntapInstancesResourceGroupReferenceToMap(model.ResourceGroup)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	modelMap["resource_type"] = model.ResourceType
	routingTables := []map[string]interface{}{}
	for _, routingTablesItem := range model.RoutingTables {
		routingTablesItemMap, err := dataSourceIbmIsStorageOntapInstancesRoutingTableReferenceToMap(&routingTablesItem)
		if err != nil {
			return modelMap, err
		}
		routingTables = append(routingTables, routingTablesItemMap)
	}
	modelMap["routing_tables"] = routingTables
	if model.SecondarySubnet != nil {
		secondarySubnetMap, err := dataSourceIbmIsStorageOntapInstancesSubnetReferenceToMap(model.SecondarySubnet)
		if err != nil {
			return modelMap, err
		}
		modelMap["secondary_subnet"] = []map[string]interface{}{secondarySubnetMap}
	}
	securityGroups := []map[string]interface{}{}
	for _, securityGroupsItem := range model.SecurityGroups {
		securityGroupsItemMap, err := dataSourceIbmIsStorageOntapInstancesSecurityGroupReferenceToMap(&securityGroupsItem)
		if err != nil {
			return modelMap, err
		}
		securityGroups = append(securityGroups, securityGroupsItemMap)
	}
	modelMap["security_groups"] = securityGroups
	storageVirtualMachines := []map[string]interface{}{}
	for _, storageVirtualMachinesItem := range model.StorageVirtualMachines {
		storageVirtualMachinesItemMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceStorageVirtualMachineReferenceToMap(&storageVirtualMachinesItem)
		if err != nil {
			return modelMap, err
		}
		storageVirtualMachines = append(storageVirtualMachines, storageVirtualMachinesItemMap)
	}
	modelMap["storage_virtual_machines"] = storageVirtualMachines
	vpcMap, err := dataSourceIbmIsStorageOntapInstancesVPCReferenceToMap(model.Vpc)
	if err != nil {
		return modelMap, err
	}
	modelMap["vpc"] = []map[string]interface{}{vpcMap}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesAddressPrefixReferenceToMap(model *ontapv1.AddressPrefixReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstancesAddressPrefixReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstancesAddressPrefixReferenceDeletedToMap(model *ontapv1.AddressPrefixReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceAdminCredentialsToMap(model *ontapv1.StorageOntapInstanceAdminCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Http != nil {
		httpMap, err := dataSourceIbmIsStorageOntapInstancesCredentialReferenceToMap(model.Http)
		if err != nil {
			return modelMap, err
		}
		modelMap["http"] = []map[string]interface{}{httpMap}
	}
	if model.Password != nil {
		passwordMap, err := dataSourceIbmIsStorageOntapInstancesCredentialReferenceToMap(model.Password)
		if err != nil {
			return modelMap, err
		}
		modelMap["password"] = []map[string]interface{}{passwordMap}
	}
	if model.Ssh != nil {
		sshMap, err := dataSourceIbmIsStorageOntapInstancesCredentialReferenceToMap(model.Ssh)
		if err != nil {
			return modelMap, err
		}
		modelMap["ssh"] = []map[string]interface{}{sshMap}
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesCredentialReferenceToMap(model *ontapv1.CredentialReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesEncryptionKeyReferenceToMap(model *ontapv1.EncryptionKeyReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceEndpointsToMap(model *ontapv1.StorageOntapInstanceEndpoints) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	interClusterMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceEndpointToMap(model.InterCluster)
	if err != nil {
		return modelMap, err
	}
	modelMap["inter_cluster"] = []map[string]interface{}{interClusterMap}
	managementMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceEndpointToMap(model.Management)
	if err != nil {
		return modelMap, err
	}
	modelMap["management"] = []map[string]interface{}{managementMap}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceEndpointToMap(model *ontapv1.StorageOntapInstanceEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["ipv4_address"] = model.Ipv4Address
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceHealthReasonToMap(model *ontapv1.StorageOntapInstanceHealthReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceLifecycleReasonToMap(model *ontapv1.StorageOntapInstanceLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesSubnetReferenceToMap(model *ontapv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstancesSubnetReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesSubnetReferenceDeletedToMap(model *ontapv1.SubnetReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesResourceGroupReferenceToMap(model *ontapv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesRoutingTableReferenceToMap(model *ontapv1.RoutingTableReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstancesRoutingTableReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstancesRoutingTableReferenceDeletedToMap(model *ontapv1.RoutingTableReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesSecurityGroupReferenceToMap(model *ontapv1.SecurityGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstancesSecurityGroupReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesSecurityGroupReferenceDeletedToMap(model *ontapv1.SecurityGroupReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceStorageVirtualMachineReferenceToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceStorageVirtualMachineReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstancesStorageOntapInstanceStorageVirtualMachineReferenceDeletedToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesVPCReferenceToMap(model *ontapv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstancesVPCReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstancesVPCReferenceDeletedToMap(model *ontapv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
