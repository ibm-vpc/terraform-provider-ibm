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
	"github.com/IBM/vpc-beta-go-sdk/ontapv1"
)

func DataSourceIbmIsStorageOntapInstance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsStorageOntapInstanceRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The storage ontap instance identifier.",
			},
			"active_subnet": &schema.Schema{
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
			"admin_password": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The password credential for the cluster administrator of the storage ontap instance.",
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
			"standby_subnet": &schema.Schema{
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
	}
}

func dataSourceIbmIsStorageOntapInstanceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getStorageOntapInstanceOptions := &ontapv1.GetStorageOntapInstanceOptions{}

	getStorageOntapInstanceOptions.SetID(d.Get("id").(string))

	storageOntapInstance, response, err := ontapClient.GetStorageOntapInstanceWithContext(context, getStorageOntapInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] GetStorageOntapInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetStorageOntapInstanceWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s", *getStorageOntapInstanceOptions.ID))

	activeSubnet := []map[string]interface{}{}
	if storageOntapInstance.ActiveSubnet != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceActiveSubnetToMap(storageOntapInstance.ActiveSubnet)
		if err != nil {
			return diag.FromErr(err)
		}
		activeSubnet = append(activeSubnet, modelMap)
	}
	if err = d.Set("active_subnet", activeSubnet); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting active_subnet %s", err))
	}

	addressPrefix := []map[string]interface{}{}
	if storageOntapInstance.AddressPrefix != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceAddressPrefixReferenceToMap(storageOntapInstance.AddressPrefix)
		if err != nil {
			return diag.FromErr(err)
		}
		addressPrefix = append(addressPrefix, modelMap)
	}
	if err = d.Set("address_prefix", addressPrefix); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting address_prefix %s", err))
	}

	adminCredentials := []map[string]interface{}{}
	if storageOntapInstance.AdminCredentials != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceAdminCredentialsToMap(storageOntapInstance.AdminCredentials)
		if err != nil {
			return diag.FromErr(err)
		}
		adminCredentials = append(adminCredentials, modelMap)
	}
	if err = d.Set("admin_credentials", adminCredentials); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting admin_credentials %s", err))
	}

	adminPassword := []map[string]interface{}{}
	if storageOntapInstance.AdminPassword != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceCredentialReferenceToMap(storageOntapInstance.AdminPassword)
		if err != nil {
			return diag.FromErr(err)
		}
		adminPassword = append(adminPassword, modelMap)
	}
	if err = d.Set("admin_password", adminPassword); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting admin_password %s", err))
	}

	if err = d.Set("capacity", flex.IntValue(storageOntapInstance.Capacity)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting capacity: %s", err))
	}

	if err = d.Set("created_at", flex.DateTimeToString(storageOntapInstance.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}

	if err = d.Set("crn", storageOntapInstance.Crn); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}

	if err = d.Set("encryption", storageOntapInstance.Encryption); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption: %s", err))
	}

	encryptionKey := []map[string]interface{}{}
	if storageOntapInstance.EncryptionKey != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceEncryptionKeyReferenceToMap(storageOntapInstance.EncryptionKey)
		if err != nil {
			return diag.FromErr(err)
		}
		encryptionKey = append(encryptionKey, modelMap)
	}
	if err = d.Set("encryption_key", encryptionKey); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption_key %s", err))
	}

	endpoints := []map[string]interface{}{}
	if storageOntapInstance.Endpoints != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointsToMap(storageOntapInstance.Endpoints)
		if err != nil {
			return diag.FromErr(err)
		}
		endpoints = append(endpoints, modelMap)
	}
	if err = d.Set("endpoints", endpoints); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoints %s", err))
	}

	healthReasons := []map[string]interface{}{}
	if storageOntapInstance.HealthReasons != nil {
		for _, modelItem := range storageOntapInstance.HealthReasons {
			modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceHealthReasonToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			healthReasons = append(healthReasons, modelMap)
		}
	}
	if err = d.Set("health_reasons", healthReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting health_reasons %s", err))
	}

	if err = d.Set("health_state", storageOntapInstance.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting health_state: %s", err))
	}

	if err = d.Set("href", storageOntapInstance.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}

	lifecycleReasons := []map[string]interface{}{}
	if storageOntapInstance.LifecycleReasons != nil {
		for _, modelItem := range storageOntapInstance.LifecycleReasons {
			modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceLifecycleReasonToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			lifecycleReasons = append(lifecycleReasons, modelMap)
		}
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_reasons %s", err))
	}

	if err = d.Set("lifecycle_state", storageOntapInstance.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}

	if err = d.Set("name", storageOntapInstance.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	primarySubnet := []map[string]interface{}{}
	if storageOntapInstance.PrimarySubnet != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceSubnetReferenceToMap(storageOntapInstance.PrimarySubnet)
		if err != nil {
			return diag.FromErr(err)
		}
		primarySubnet = append(primarySubnet, modelMap)
	}
	if err = d.Set("primary_subnet", primarySubnet); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting primary_subnet %s", err))
	}

	resourceGroup := []map[string]interface{}{}
	if storageOntapInstance.ResourceGroup != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceResourceGroupReferenceToMap(storageOntapInstance.ResourceGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		resourceGroup = append(resourceGroup, modelMap)
	}
	if err = d.Set("resource_group", resourceGroup); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
	}

	if err = d.Set("resource_type", storageOntapInstance.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	routingTables := []map[string]interface{}{}
	if storageOntapInstance.RoutingTables != nil {
		for _, modelItem := range storageOntapInstance.RoutingTables {
			modelMap, err := dataSourceIbmIsStorageOntapInstanceRoutingTableReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			routingTables = append(routingTables, modelMap)
		}
	}
	if err = d.Set("routing_tables", routingTables); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting routing_tables %s", err))
	}

	secondarySubnet := []map[string]interface{}{}
	if storageOntapInstance.SecondarySubnet != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceSubnetReferenceToMap(storageOntapInstance.SecondarySubnet)
		if err != nil {
			return diag.FromErr(err)
		}
		secondarySubnet = append(secondarySubnet, modelMap)
	}
	if err = d.Set("secondary_subnet", secondarySubnet); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting secondary_subnet %s", err))
	}

	securityGroups := []map[string]interface{}{}
	if storageOntapInstance.SecurityGroups != nil {
		for _, modelItem := range storageOntapInstance.SecurityGroups {
			modelMap, err := dataSourceIbmIsStorageOntapInstanceSecurityGroupReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			securityGroups = append(securityGroups, modelMap)
		}
	}
	if err = d.Set("security_groups", securityGroups); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting security_groups %s", err))
	}

	standbySubnet := []map[string]interface{}{}
	if storageOntapInstance.StandbySubnet != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceStandbySubnetToMap(storageOntapInstance.StandbySubnet)
		if err != nil {
			return diag.FromErr(err)
		}
		standbySubnet = append(standbySubnet, modelMap)
	}
	if err = d.Set("standby_subnet", standbySubnet); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting standby_subnet %s", err))
	}

	storageVirtualMachines := []map[string]interface{}{}
	if storageOntapInstance.StorageVirtualMachines != nil {
		for _, modelItem := range storageOntapInstance.StorageVirtualMachines {
			modelMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceToMap(&modelItem)
			if err != nil {
				return diag.FromErr(err)
			}
			storageVirtualMachines = append(storageVirtualMachines, modelMap)
		}
	}
	if err = d.Set("storage_virtual_machines", storageVirtualMachines); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting storage_virtual_machines %s", err))
	}

	vpc := []map[string]interface{}{}
	if storageOntapInstance.Vpc != nil {
		modelMap, err := dataSourceIbmIsStorageOntapInstanceVPCReferenceToMap(storageOntapInstance.Vpc)
		if err != nil {
			return diag.FromErr(err)
		}
		vpc = append(vpc, modelMap)
	}
	if err = d.Set("vpc", vpc); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vpc %s", err))
	}

	return nil
}

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceActiveSubnetToMap(model *ontapv1.StorageOntapInstanceActiveSubnet) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceSubnetReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstanceSubnetReferenceDeletedToMap(model *ontapv1.SubnetReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceAddressPrefixReferenceToMap(model *ontapv1.AddressPrefixReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceAddressPrefixReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstanceAddressPrefixReferenceDeletedToMap(model *ontapv1.AddressPrefixReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceAdminCredentialsToMap(model *ontapv1.StorageOntapInstanceAdminCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Http != nil {
		httpMap, err := dataSourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model.Http)
		if err != nil {
			return modelMap, err
		}
		modelMap["http"] = []map[string]interface{}{httpMap}
	}
	if model.Password != nil {
		passwordMap, err := dataSourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model.Password)
		if err != nil {
			return modelMap, err
		}
		modelMap["password"] = []map[string]interface{}{passwordMap}
	}
	if model.Ssh != nil {
		sshMap, err := dataSourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model.Ssh)
		if err != nil {
			return modelMap, err
		}
		modelMap["ssh"] = []map[string]interface{}{sshMap}
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model *ontapv1.CredentialReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceEncryptionKeyReferenceToMap(model *ontapv1.EncryptionKeyReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointsToMap(model *ontapv1.StorageOntapInstanceEndpoints) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	interClusterMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointToMap(model.InterCluster)
	if err != nil {
		return modelMap, err
	}
	modelMap["inter_cluster"] = []map[string]interface{}{interClusterMap}
	managementMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointToMap(model.Management)
	if err != nil {
		return modelMap, err
	}
	modelMap["management"] = []map[string]interface{}{managementMap}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointToMap(model *ontapv1.StorageOntapInstanceEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["ipv4_address"] = model.Ipv4Address
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceHealthReasonToMap(model *ontapv1.StorageOntapInstanceHealthReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceLifecycleReasonToMap(model *ontapv1.StorageOntapInstanceLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceSubnetReferenceToMap(model *ontapv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceSubnetReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstanceResourceGroupReferenceToMap(model *ontapv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceRoutingTableReferenceToMap(model *ontapv1.RoutingTableReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceRoutingTableReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstanceRoutingTableReferenceDeletedToMap(model *ontapv1.RoutingTableReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceSecurityGroupReferenceToMap(model *ontapv1.SecurityGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceSecurityGroupReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceSecurityGroupReferenceDeletedToMap(model *ontapv1.SecurityGroupReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceStandbySubnetToMap(model *ontapv1.StorageOntapInstanceStandbySubnet) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceSubnetReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceDeletedToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func dataSourceIbmIsStorageOntapInstanceVPCReferenceToMap(model *ontapv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := dataSourceIbmIsStorageOntapInstanceVPCReferenceDeletedToMap(model.Deleted)
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

func dataSourceIbmIsStorageOntapInstanceVPCReferenceDeletedToMap(model *ontapv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
