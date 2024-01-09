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
	"github.com/IBM/vpc-beta-go-sdk/ontapv1"
)

func ResourceIbmIsStorageOntapInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsStorageOntapInstanceCreate,
		ReadContext:   resourceIbmIsStorageOntapInstanceRead,
		UpdateContext: resourceIbmIsStorageOntapInstanceUpdate,
		DeleteContext: resourceIbmIsStorageOntapInstanceDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"address_prefix": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
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
										Required:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"address_prefix.0.id", "address_prefix.0.href"},
							Description:  "The URL for this address prefix.",
						},
						"id": &schema.Schema{
							Type:         schema.TypeString,
							ExactlyOneOf: []string{"address_prefix.0.id", "address_prefix.0.href"},
							Optional:     true,
							Description:  "The unique identifier for this address prefix.",
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
				MaxItems:    1,
				Optional:    true,
				Description: "The credentials used (from Secrets Manager) for the cluster administrator to access thestorage ontap instance. At least one of `password`, `ssh`, or `http` will be present.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The security certificate credential for ONTAP REST API access for the clusteradministrator of the storage ontap instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "The password credential for the cluster administrator of the storage ontap instance.If present, this password credential is used by the cluster administrator for bothONTAP CLI SSH access and ONTAP REST API access. If absent, the storage ontapinstance is not accessible through either the ONTAP CLI or ONTAP REST API usingpassword-based authentication.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
							MaxItems:    1,
							Optional:    true,
							Description: "The public key credential for ONTAP CLI SSH access for the cluster administratorof the storage ontap instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance", "capacity"),
				Description:  "The capacity to use for the storage ontap instance (in terabytes). Volumes in this storage ontap instance will be allocated from this capacity.",
			},
			"encryption_key": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The root key used to wrap the data encryption key for the storage ontap instance.This property will be present for storage ontap instance with an `encryption` type of`user_managed`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Services Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance", "name"),
				Description:  "The name for this storage ontap instance. The name is unique across all storage ontap instances in the region.",
			},
			"primary_subnet": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The subnet where the primary Cloud Volumes ONTAP node is provisioned in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"primary_subnet.0.id", "primary_subnet.0.crn"},
							Description:  "The CRN for this subnet.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
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
						"id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"primary_subnet.0.id", "primary_subnet.0.crn"},
							Description:  "The unique identifier for this subnet.",
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
				MaxItems:    1,
				Optional:    true,
				Description: "The resource group for this storage ontap instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
			"routing_tables": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The VPC routing tables for this storage ontap instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The URL for this routing table.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this routing table.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name for this routing table. The name is unique across all routing tables for the VPC.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"secondary_subnet": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The subnet where the secondary Cloud Volumes ONTAP node is provisioned in.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"secondary_subnet.0.id", "secondary_subnet.0.crn"},
							Description:  "The CRN for this subnet.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
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
						"id": &schema.Schema{
							Type:         schema.TypeString,
							Optional:     true,
							ExactlyOneOf: []string{"secondary_subnet.0.id", "secondary_subnet.0.crn"},
							Description:  "The unique identifier for this subnet.",
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
				Optional:    true,
				Description: "The security groups for this storage ontap instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Optional:    true,
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
				Required:    true,
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
						"active_directory": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The Active Directory service this storage virtual machine is joined to.If absent, this storage virtual machine is not joined to an Active Directory service.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"administrators_group": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the domain group whose members have been granted administrative privileges for this storage virtual machine.",
									},
									"dns_ips": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The IP addresses of the Active Directory DNS servers or domain controllers.",
										Elem: &schema.Schema{
											Type: schema.TypeString,
										},
									},
									"domain_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The fully qualified domain name of the self-managed Active Directory.",
									},
									"domain_password_credential": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										MaxItems:    1,
										Description: "The password credential for the Active Directory domain.",
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
									"netbios_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The name of the Active Directory computer object that will be created for the storage virtual machine.",
									},
									"organizational_unit_distinguished_name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The distinguished name of the organizational unit within the self-managed Active Directory.",
									},
									"username": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The username that this storage virtual machine will use when joining the Active Directory domain. This username will be the same as the username credential used in `domain_password_credential`.",
									},
								},
							},
						},
						"admin_credentials": &schema.Schema{
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The credentials used for the administrator to access the storage virtual machine of thestorage ontap instance. At least one of `password`, `ssh`, or `http` will be present.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"http": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Description: "The security certificate credential for ONTAP REST API access for the storage virtualmachine administrator.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
										Optional:    true,
										Description: "The password credential for the storage virtual machine administrator.If present, this password credential is used by the storage virtual machineadministrator for both ONTAP CLI SSH access and ONTAP REST API access.If absent, the storage virtual machine is not accessible through either the ONTAP CLIor ONTAP REST API using password-based authentication.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
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
										Optional:    true,
										Description: "The public key credential for ONTAP CLI based ssh login for the storage virtualmachine administrator.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this credential.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The resource type.",
												},
											},
										},
									},
								},
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name for this storage virtual machine. The name is unique across all storage virtual machines in the storage ontap instance.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type.",
						},
					},
				},
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
			"endpoints": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The endpoints for this storage ontap instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"inter_cluster": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The NetApp SnapMirror management endpoint for this storage ontap instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4_address": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The unique IP address of an endpoint.",
									},
								},
							},
						},
						"management": &schema.Schema{
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The NetApp management endpoint for this storage ontap instance. Management may beperformed using the ONTAP CLI, ONTAP API, or NetApp CloudManager.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"ipv4_address": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
							Required:    true,
							Description: "A snake case string succinctly identifying the reason for this health state.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "An explanation of the reason for this health state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
							Required:    true,
							Description: "A snake case string succinctly identifying the reason for this lifecycle state.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "An explanation of the reason for this lifecycle state.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
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
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"vpc": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The VPC this storage ontap instance resides in.",
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this VPC.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The name for this VPC. The name is unique across all VPCs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmIsStorageOntapInstanceValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "capacity",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Required:                   true,
			MinValue:                   "1",
			MaxValue:                   "64",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^(ibmss)-([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             7,
			MaxValueLength:             40,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_storage_ontap_instance", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsStorageOntapInstanceCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	createStorageOntapInstanceOptions := &ontapv1.CreateStorageOntapInstanceOptions{}

	addressPrefixModel, err := resourceIbmIsStorageOntapInstanceMapToAddressPrefixIdentity(d.Get("address_prefix.0").(map[string]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}
	createStorageOntapInstanceOptions.SetAddressPrefix(addressPrefixModel)
	createStorageOntapInstanceOptions.SetCapacity(int64(d.Get("capacity").(int)))
	var storageVirtualMachines []ontapv1.StorageOntapInstanceStorageVirtualMachinePrototype
	for _, v := range d.Get("storage_virtual_machines").([]interface{}) {
		value := v.(map[string]interface{})
		storageVirtualMachinesItem, err := resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceStorageVirtualMachinePrototype(value)
		if err != nil {
			return diag.FromErr(err)
		}
		storageVirtualMachines = append(storageVirtualMachines, *storageVirtualMachinesItem)
	}
	createStorageOntapInstanceOptions.SetStorageVirtualMachines(storageVirtualMachines)
	if _, ok := d.GetOk("admin_credentials"); ok {
		adminCredentialsModel, err := resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceAdminCredentialsPrototype(d.Get("admin_credentials.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createStorageOntapInstanceOptions.SetAdminCredentials(adminCredentialsModel)
	}
	if _, ok := d.GetOk("encryption_key"); ok {
		encryptionKeyModel, err := resourceIbmIsStorageOntapInstanceMapToEncryptionKeyIdentity(d.Get("encryption_key.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createStorageOntapInstanceOptions.SetEncryptionKey(encryptionKeyModel)
	}
	if _, ok := d.GetOk("name"); ok {
		createStorageOntapInstanceOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("primary_subnet"); ok {
		primarySubnetModel, err := resourceIbmIsStorageOntapInstanceMapToSubnetIdentity(d.Get("primary_subnet.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createStorageOntapInstanceOptions.SetPrimarySubnet(primarySubnetModel)
	}
	if _, ok := d.GetOk("resource_group"); ok {
		resourceGroupModel, err := resourceIbmIsStorageOntapInstanceMapToResourceGroupIdentity(d.Get("resource_group.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createStorageOntapInstanceOptions.SetResourceGroup(resourceGroupModel)
	}
	if _, ok := d.GetOk("routing_tables"); ok {
		var routingTables []ontapv1.RoutingTableIdentityIntf
		for _, v := range d.Get("routing_tables").([]interface{}) {
			value := v.(map[string]interface{})
			routingTablesItem, err := resourceIbmIsStorageOntapInstanceMapToRoutingTableIdentity(value)
			if err != nil {
				return diag.FromErr(err)
			}
			routingTables = append(routingTables, routingTablesItem)
		}
		createStorageOntapInstanceOptions.SetRoutingTables(routingTables)
	}
	if _, ok := d.GetOk("secondary_subnet"); ok {
		secondarySubnetModel, err := resourceIbmIsStorageOntapInstanceMapToSubnetIdentity(d.Get("secondary_subnet.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		createStorageOntapInstanceOptions.SetSecondarySubnet(secondarySubnetModel)
	}
	if _, ok := d.GetOk("security_groups"); ok {
		var securityGroups []ontapv1.SecurityGroupIdentityIntf
		for _, v := range d.Get("security_groups").([]interface{}) {
			value := v.(map[string]interface{})
			securityGroupsItem, err := resourceIbmIsStorageOntapInstanceMapToSecurityGroupIdentity(value)
			if err != nil {
				return diag.FromErr(err)
			}
			securityGroups = append(securityGroups, securityGroupsItem)
		}
		createStorageOntapInstanceOptions.SetSecurityGroups(securityGroups)
	}

	storageOntapInstance, response, err := ontapClient.CreateStorageOntapInstanceWithContext(context, createStorageOntapInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateStorageOntapInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateStorageOntapInstanceWithContext failed %s\n%s", err, response))
	}

	d.SetId(*storageOntapInstance.ID)

	return resourceIbmIsStorageOntapInstanceRead(context, d, meta)
}

func resourceIbmIsStorageOntapInstanceRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getStorageOntapInstanceOptions := &ontapv1.GetStorageOntapInstanceOptions{}

	getStorageOntapInstanceOptions.SetID(d.Id())

	storageOntapInstance, response, err := ontapClient.GetStorageOntapInstanceWithContext(context, getStorageOntapInstanceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetStorageOntapInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetStorageOntapInstanceWithContext failed %s\n%s", err, response))
	}

	addressPrefixMap, err := resourceIbmIsStorageOntapInstanceAddressPrefixReferenceToMap(storageOntapInstance.AddressPrefix)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("address_prefix", []map[string]interface{}{addressPrefixMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting address_prefix: %s", err))
	}
	if !core.IsNil(storageOntapInstance.AdminCredentials) {
		adminCredentialsMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceAdminCredentialsToMap(storageOntapInstance.AdminCredentials)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("admin_credentials", []map[string]interface{}{adminCredentialsMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting admin_credentials: %s", err))
		}
	}
	if err = d.Set("capacity", flex.IntValue(storageOntapInstance.Capacity)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting capacity: %s", err))
	}
	if !core.IsNil(storageOntapInstance.EncryptionKey) {
		encryptionKeyMap, err := resourceIbmIsStorageOntapInstanceEncryptionKeyReferenceToMap(storageOntapInstance.EncryptionKey)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("encryption_key", []map[string]interface{}{encryptionKeyMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting encryption_key: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstance.Name) {
		if err = d.Set("name", storageOntapInstance.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstance.PrimarySubnet) {
		primarySubnetMap, err := resourceIbmIsStorageOntapInstanceSubnetReferenceToMap(storageOntapInstance.PrimarySubnet)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("primary_subnet", []map[string]interface{}{primarySubnetMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting primary_subnet: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstance.ResourceGroup) {
		resourceGroupMap, err := resourceIbmIsStorageOntapInstanceResourceGroupReferenceToMap(storageOntapInstance.ResourceGroup)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstance.RoutingTables) {
		routingTables := []map[string]interface{}{}
		for _, routingTablesItem := range storageOntapInstance.RoutingTables {
			routingTablesItemMap, err := resourceIbmIsStorageOntapInstanceRoutingTableReferenceToMap(&routingTablesItem)
			if err != nil {
				return diag.FromErr(err)
			}
			routingTables = append(routingTables, routingTablesItemMap)
		}
		if err = d.Set("routing_tables", routingTables); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting routing_tables: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstance.SecondarySubnet) {
		secondarySubnetMap, err := resourceIbmIsStorageOntapInstanceSubnetReferenceToMap(storageOntapInstance.SecondarySubnet)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("secondary_subnet", []map[string]interface{}{secondarySubnetMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting secondary_subnet: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstance.SecurityGroups) {
		securityGroups := []map[string]interface{}{}
		for _, securityGroupsItem := range storageOntapInstance.SecurityGroups {
			securityGroupsItemMap, err := resourceIbmIsStorageOntapInstanceSecurityGroupReferenceToMap(&securityGroupsItem)
			if err != nil {
				return diag.FromErr(err)
			}
			securityGroups = append(securityGroups, securityGroupsItemMap)
		}
		if err = d.Set("security_groups", securityGroups); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting security_groups: %s", err))
		}
	}
	storageVirtualMachines := []map[string]interface{}{}
	for _, storageVirtualMachinesItem := range storageOntapInstance.StorageVirtualMachines {
		storageVirtualMachinesItemMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceToMap(&storageVirtualMachinesItem)
		if err != nil {
			return diag.FromErr(err)
		}
		storageVirtualMachines = append(storageVirtualMachines, storageVirtualMachinesItemMap)
	}
	if err = d.Set("storage_virtual_machines", storageVirtualMachines); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting storage_virtual_machines: %s", err))
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
	endpointsMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointsToMap(storageOntapInstance.Endpoints)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("endpoints", []map[string]interface{}{endpointsMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting endpoints: %s", err))
	}
	healthReasons := []map[string]interface{}{}
	for _, healthReasonsItem := range storageOntapInstance.HealthReasons {
		healthReasonsItemMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceHealthReasonToMap(&healthReasonsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		healthReasons = append(healthReasons, healthReasonsItemMap)
	}
	if err = d.Set("health_reasons", healthReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting health_reasons: %s", err))
	}
	if err = d.Set("health_state", storageOntapInstance.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting health_state: %s", err))
	}
	if err = d.Set("href", storageOntapInstance.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	lifecycleReasons := []map[string]interface{}{}
	for _, lifecycleReasonsItem := range storageOntapInstance.LifecycleReasons {
		lifecycleReasonsItemMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceLifecycleReasonToMap(&lifecycleReasonsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		lifecycleReasons = append(lifecycleReasons, lifecycleReasonsItemMap)
	}
	if err = d.Set("lifecycle_reasons", lifecycleReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_reasons: %s", err))
	}
	if err = d.Set("lifecycle_state", storageOntapInstance.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", storageOntapInstance.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	vpcMap, err := resourceIbmIsStorageOntapInstanceVPCReferenceToMap(storageOntapInstance.Vpc)
	if err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("vpc", []map[string]interface{}{vpcMap}); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting vpc: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIbmIsStorageOntapInstanceUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	updateStorageOntapInstanceOptions := &ontapv1.UpdateStorageOntapInstanceOptions{}

	updateStorageOntapInstanceOptions.SetID(d.Id())

	hasChange := false

	patchVals := &ontapv1.StorageOntapInstancePatch{}
	if d.HasChange("admin_credentials") {
		adminCredentials, err := resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceAdminCredentialsPatch(d.Get("admin_credentials.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		patchVals.AdminCredentials = adminCredentials
		hasChange = true
	}
	if d.HasChange("capacity") {
		newCapacity := int64(d.Get("capacity").(int))
		patchVals.Capacity = &newCapacity
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("routing_tables") {
		var routingTables []ontapv1.RoutingTableIdentityIntf
		for _, v := range d.Get("routing_tables").([]interface{}) {
			value := v.(map[string]interface{})
			routingTablesItem, err := resourceIbmIsStorageOntapInstanceMapToRoutingTableIdentity(value)
			if err != nil {
				return diag.FromErr(err)
			}
			routingTables = append(routingTables, routingTablesItem)
		}
		patchVals.RoutingTables = routingTables
		hasChange = true
	}
	updateStorageOntapInstanceOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateStorageOntapInstanceOptions.StorageOntapInstancePatch, _ = patchVals.AsPatch()
		_, response, err := ontapClient.UpdateStorageOntapInstanceWithContext(context, updateStorageOntapInstanceOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateStorageOntapInstanceWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateStorageOntapInstanceWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmIsStorageOntapInstanceRead(context, d, meta)
}

func resourceIbmIsStorageOntapInstanceDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteStorageOntapInstanceOptions := &ontapv1.DeleteStorageOntapInstanceOptions{}

	deleteStorageOntapInstanceOptions.SetID(d.Id())

	_, response, err := ontapClient.DeleteStorageOntapInstanceWithContext(context, deleteStorageOntapInstanceOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteStorageOntapInstanceWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteStorageOntapInstanceWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmIsStorageOntapInstanceMapToAddressPrefixIdentity(modelMap map[string]interface{}) (ontapv1.AddressPrefixIdentityIntf, error) {
	model := &ontapv1.AddressPrefixIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToAddressPrefixIdentityByID(modelMap map[string]interface{}) (*ontapv1.AddressPrefixIdentityByID, error) {
	model := &ontapv1.AddressPrefixIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToAddressPrefixIdentityByHref(modelMap map[string]interface{}) (*ontapv1.AddressPrefixIdentityByHref, error) {
	model := &ontapv1.AddressPrefixIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceStorageVirtualMachinePrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachinePrototype, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachinePrototype{}
	if modelMap["active_directory"] != nil && len(modelMap["active_directory"].([]interface{})) > 0 {
		ActiveDirectoryModel, err := resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceStorageVirtualMachineActiveDirectoryPrototype(modelMap["active_directory"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ActiveDirectory = ActiveDirectoryModel
	}
	if modelMap["admin_credentials"] != nil && len(modelMap["admin_credentials"].([]interface{})) > 0 {
		AdminCredentialsModel, err := resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceStorageVirtualMachineAdminCredentialsPrototype(modelMap["admin_credentials"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AdminCredentials = AdminCredentialsModel
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceStorageVirtualMachineActiveDirectoryPrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineActiveDirectoryPrototype, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineActiveDirectoryPrototype{}
	model.AdministratorsGroup = core.StringPtr(modelMap["administrators_group"].(string))
	dnsIps := []string{}
	for _, dnsIpsItem := range modelMap["dns_ips"].([]interface{}) {
		dnsIps = append(dnsIps, dnsIpsItem.(string))
	}
	model.DnsIps = dnsIps
	model.DomainName = core.StringPtr(modelMap["domain_name"].(string))
	if modelMap["domain_password_credential"] != nil && len(modelMap["domain_password_credential"].([]interface{})) > 0 {
		DomainPasswordCredentialModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["domain_password_credential"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.DomainPasswordCredential = DomainPasswordCredentialModel
	}
	model.NetbiosName = core.StringPtr(modelMap["netbios_name"].(string))
	model.OrganizationalUnitDistinguishedName = core.StringPtr(modelMap["organizational_unit_distinguished_name"].(string))
	model.Username = core.StringPtr(modelMap["username"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap map[string]interface{}) (ontapv1.CredentialIdentityIntf, error) {
	model := &ontapv1.CredentialIdentity{}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.Crn = core.StringPtr(modelMap["crn"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToCredentialIdentityByCRN(modelMap map[string]interface{}) (*ontapv1.CredentialIdentityByCRN, error) {
	model := &ontapv1.CredentialIdentityByCRN{}
	model.Crn = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceStorageVirtualMachineAdminCredentialsPrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineAdminCredentialsPrototype, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineAdminCredentialsPrototype{}
	if modelMap["http"] != nil && len(modelMap["http"].([]interface{})) > 0 {
		HttpModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["http"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Http = HttpModel
	}
	if modelMap["password"] != nil && len(modelMap["password"].([]interface{})) > 0 {
		PasswordModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["password"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Password = PasswordModel
	}
	if modelMap["ssh"] != nil && len(modelMap["ssh"].([]interface{})) > 0 {
		SshModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["ssh"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ssh = SshModel
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceAdminCredentialsPrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceAdminCredentialsPrototype, error) {
	model := &ontapv1.StorageOntapInstanceAdminCredentialsPrototype{}
	if modelMap["http"] != nil && len(modelMap["http"].([]interface{})) > 0 {
		HttpModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["http"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Http = HttpModel
	}
	if modelMap["password"] != nil && len(modelMap["password"].([]interface{})) > 0 {
		PasswordModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["password"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Password = PasswordModel
	}
	if modelMap["ssh"] != nil && len(modelMap["ssh"].([]interface{})) > 0 {
		SshModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["ssh"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ssh = SshModel
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToEncryptionKeyIdentity(modelMap map[string]interface{}) (ontapv1.EncryptionKeyIdentityIntf, error) {
	model := &ontapv1.EncryptionKeyIdentity{}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.Crn = core.StringPtr(modelMap["crn"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToEncryptionKeyIdentityByCRN(modelMap map[string]interface{}) (*ontapv1.EncryptionKeyIdentityByCRN, error) {
	model := &ontapv1.EncryptionKeyIdentityByCRN{}
	model.Crn = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSubnetIdentity(modelMap map[string]interface{}) (ontapv1.SubnetIdentityIntf, error) {
	model := &ontapv1.SubnetIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.Crn = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSubnetIdentityByID(modelMap map[string]interface{}) (*ontapv1.SubnetIdentityByID, error) {
	model := &ontapv1.SubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSubnetIdentityByCRN(modelMap map[string]interface{}) (*ontapv1.SubnetIdentityByCRN, error) {
	model := &ontapv1.SubnetIdentityByCRN{}
	model.Crn = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSubnetIdentityByHref(modelMap map[string]interface{}) (*ontapv1.SubnetIdentityByHref, error) {
	model := &ontapv1.SubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToResourceGroupIdentity(modelMap map[string]interface{}) (ontapv1.ResourceGroupIdentityIntf, error) {
	model := &ontapv1.ResourceGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToResourceGroupIdentityByID(modelMap map[string]interface{}) (*ontapv1.ResourceGroupIdentityByID, error) {
	model := &ontapv1.ResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToRoutingTableIdentity(modelMap map[string]interface{}) (ontapv1.RoutingTableIdentityIntf, error) {
	model := &ontapv1.RoutingTableIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToRoutingTableIdentityByID(modelMap map[string]interface{}) (*ontapv1.RoutingTableIdentityByID, error) {
	model := &ontapv1.RoutingTableIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToRoutingTableIdentityByHref(modelMap map[string]interface{}) (*ontapv1.RoutingTableIdentityByHref, error) {
	model := &ontapv1.RoutingTableIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSecurityGroupIdentity(modelMap map[string]interface{}) (ontapv1.SecurityGroupIdentityIntf, error) {
	model := &ontapv1.SecurityGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.Crn = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSecurityGroupIdentityByID(modelMap map[string]interface{}) (*ontapv1.SecurityGroupIdentityByID, error) {
	model := &ontapv1.SecurityGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSecurityGroupIdentityByCRN(modelMap map[string]interface{}) (*ontapv1.SecurityGroupIdentityByCRN, error) {
	model := &ontapv1.SecurityGroupIdentityByCRN{}
	model.Crn = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToSecurityGroupIdentityByHref(modelMap map[string]interface{}) (*ontapv1.SecurityGroupIdentityByHref, error) {
	model := &ontapv1.SecurityGroupIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceMapToStorageOntapInstanceAdminCredentialsPatch(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceAdminCredentialsPatch, error) {
	model := &ontapv1.StorageOntapInstanceAdminCredentialsPatch{}
	if modelMap["http"] != nil && len(modelMap["http"].([]interface{})) > 0 {
		HttpModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["http"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Http = HttpModel
	}
	if modelMap["password"] != nil && len(modelMap["password"].([]interface{})) > 0 {
		PasswordModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["password"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Password = PasswordModel
	}
	if modelMap["ssh"] != nil && len(modelMap["ssh"].([]interface{})) > 0 {
		SshModel, err := resourceIbmIsStorageOntapInstanceMapToCredentialIdentity(modelMap["ssh"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Ssh = SshModel
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceAddressPrefixReferenceToMap(model *ontapv1.AddressPrefixReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIbmIsStorageOntapInstanceAddressPrefixReferenceDeletedToMap(model.Deleted)
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

func resourceIbmIsStorageOntapInstanceAddressPrefixReferenceDeletedToMap(model *ontapv1.AddressPrefixReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceStorageOntapInstanceAdminCredentialsToMap(model *ontapv1.StorageOntapInstanceAdminCredentials) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Http != nil {
		httpMap, err := resourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model.Http)
		if err != nil {
			return modelMap, err
		}
		modelMap["http"] = []map[string]interface{}{httpMap}
	}
	if model.Password != nil {
		passwordMap, err := resourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model.Password)
		if err != nil {
			return modelMap, err
		}
		modelMap["password"] = []map[string]interface{}{passwordMap}
	}
	if model.Ssh != nil {
		sshMap, err := resourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model.Ssh)
		if err != nil {
			return modelMap, err
		}
		modelMap["ssh"] = []map[string]interface{}{sshMap}
	}
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceCredentialReferenceToMap(model *ontapv1.CredentialReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	modelMap["resource_type"] = model.ResourceType
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceEncryptionKeyReferenceToMap(model *ontapv1.EncryptionKeyReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceSubnetReferenceToMap(model *ontapv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := resourceIbmIsStorageOntapInstanceSubnetReferenceDeletedToMap(model.Deleted)
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

func resourceIbmIsStorageOntapInstanceSubnetReferenceDeletedToMap(model *ontapv1.SubnetReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceResourceGroupReferenceToMap(model *ontapv1.ResourceGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceRoutingTableReferenceToMap(model *ontapv1.RoutingTableReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIbmIsStorageOntapInstanceRoutingTableReferenceDeletedToMap(model.Deleted)
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

func resourceIbmIsStorageOntapInstanceRoutingTableReferenceDeletedToMap(model *ontapv1.RoutingTableReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceSecurityGroupReferenceToMap(model *ontapv1.SecurityGroupReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := resourceIbmIsStorageOntapInstanceSecurityGroupReferenceDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["id"] = model.ID
	modelMap["name"] = model.Name
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceSecurityGroupReferenceDeletedToMap(model *ontapv1.SecurityGroupReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceDeletedToMap(model.Deleted)
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

func resourceIbmIsStorageOntapInstanceStorageOntapInstanceStorageVirtualMachineReferenceDeletedToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointsToMap(model *ontapv1.StorageOntapInstanceEndpoints) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	interClusterMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointToMap(model.InterCluster)
	if err != nil {
		return modelMap, err
	}
	modelMap["inter_cluster"] = []map[string]interface{}{interClusterMap}
	managementMap, err := resourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointToMap(model.Management)
	if err != nil {
		return modelMap, err
	}
	modelMap["management"] = []map[string]interface{}{managementMap}
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceStorageOntapInstanceEndpointToMap(model *ontapv1.StorageOntapInstanceEndpoint) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["ipv4_address"] = model.Ipv4Address
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceStorageOntapInstanceHealthReasonToMap(model *ontapv1.StorageOntapInstanceHealthReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceStorageOntapInstanceLifecycleReasonToMap(model *ontapv1.StorageOntapInstanceLifecycleReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceVPCReferenceToMap(model *ontapv1.VPCReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = model.Crn
	if model.Deleted != nil {
		deletedMap, err := resourceIbmIsStorageOntapInstanceVPCReferenceDeletedToMap(model.Deleted)
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

func resourceIbmIsStorageOntapInstanceVPCReferenceDeletedToMap(model *ontapv1.VPCReferenceDeleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = model.MoreInfo
	return modelMap, nil
}
