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

func ResourceIbmIsStorageOntapInstanceVirtualMachineVolume() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsStorageOntapInstanceVirtualMachineVolumeCreate,
		ReadContext:   resourceIbmIsStorageOntapInstanceVirtualMachineVolumeRead,
		UpdateContext: resourceIbmIsStorageOntapInstanceVirtualMachineVolumeUpdate,
		DeleteContext: resourceIbmIsStorageOntapInstanceVirtualMachineVolumeDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"storage_ontap_instance_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance_virtual_machine_volume", "storage_ontap_instance_id"),
				Description:  "The storage ontap instance identifier.",
			},
			"storage_virtual_machine_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance_virtual_machine_volume", "storage_virtual_machine_id"),
				Description:  "The storage virtual machine identifier.",
			},
			"capacity": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance_virtual_machine_volume", "capacity"),
				Description:  "The capacity of the storage volume (in gigabytes).",
			},
			"cifs_share": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The named access point that enables CIFS clients to view, browse, and manipulatefiles on this storage volumeThis will be present when `security_style` is `mixed` or `windows`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_control_list": &schema.Schema{
							Type:        schema.TypeList,
							Required:    true,
							Description: "The access control list for the CIFS share.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"permission": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The permission granted to users matching this access control list entry.",
									},
									"users": &schema.Schema{
										Type:        schema.TypeList,
										Required:    true,
										Description: "The users matching this access control list entry.",
										Elem:        &schema.Schema{Type: schema.TypeString},
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
							Required:    true,
							Description: "The share name registered in Active Directory that SMB/CIFS clients use to mount the share. The name is unique within the Active Directory domain.",
						},
					},
				},
			},
			"enable_storage_efficiency": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Deprecated:  "This argument is deprecated and may be removed in a future release",
				Default:     true,
				Description: "Indicates whether storage efficiency is enabled for the storage volume.If `true`, data-deduplication, compression and other efficiencies for space-management are enabled for this volume.",
			},
			"export_policy": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
							Optional:    true,
							Description: "The NFS export policy rules for this storage volume.Only NFS clients included in the rules will access the volume, and only according to the specified access controls and NFS protocol versions.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"access_control": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The access control that is provided to clients that match this rule.",
									},
									"clients": &schema.Schema{
										Type:        schema.TypeList,
										Required:    true,
										Description: "The clients that match this rule. Every client in the list of `clients` is unique.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"hostname": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The hostname of the NFS client.",
												},
												"address": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The IP address of the NFS client.This property may add support for IPv6 addresses in the future. When processing a value in this property, verify that the address is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected IP address format was encountered.",
												},
												"cidr_block": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The CIDR block containing IP addresses of the NFS clients. The CIDR block `0.0.0.0/0` matches all client addresses.This property may add support for IPv6 CIDR blocks in the future. When processing a value in this property, verify that the CIDR block is in an expected format. If it is not, log an error. Optionally halt processing and surface the error, or bypass the resource on which the unexpected CIDR block format was encountered.",
												},
												"domain_suffix": &schema.Schema{
													Type:        schema.TypeString,
													Optional:    true,
													Description: "The domain names suffixes of the NFS clients.",
												},
											},
										},
									},
									"index": &schema.Schema{
										Type:        schema.TypeInt,
										Required:    true,
										Description: "The unique index for this rule. Rules are applied from lowest to highest index.",
									},
									"is_superuser": &schema.Schema{
										Type:        schema.TypeBool,
										Required:    true,
										Description: "If `true`, clients matching this rule that request super-user access are honored. Otherwise, clients are mapped to the anonymous user.",
									},
									"nfs_version": &schema.Schema{
										Type:        schema.TypeList,
										Required:    true,
										Description: "The NFS versions that is provided to clients that match this rule.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance_virtual_machine_volume", "name"),
				Description:  "The name for this storage volume. The name is unique across all storage volumes in the storage virtual machine.",
			},
			"security_style": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "mixed",
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance_virtual_machine_volume", "security_style"),
				Description:  "The security style for the storage volume:- `unix`: NFS clients can access the storage volume.- `windows`: SMB/CIFS clients can access the storage volume.- `mixed`: Both SMB/CIFS and NFS clients can access the storage volume.- `none`: No clients can access the volume.",
			},
			"storage_efficiency": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "enabled",
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance_virtual_machine_volume", "storage_efficiency"),
				Description:  "The storage efficiency mode used for this storage volume.- `disabled`: storage efficiency methods will not be used- `enabled`: data-deduplication, compression and other methods will be used.",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "read_write",
				ValidateFunc: validate.InvokeValidator("ibm_is_storage_ontap_instance_virtual_machine_volume", "type"),
				Description:  "The type of the storage volume.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the storage volume was created.",
			},
			"health_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current storage volume health_state (if any):- `primary_node_down`: The storage volume is experiencing higher latency due to  the primary node being unavailable, and I/O being routed to the secondary node.- `volume_unavailable`: The storage volume is unavailable as both the primary and secondary nodes are down.- `internal_error`: Internal error (contact IBM support).",
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
				Description: "The URL for this storage volume.",
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
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"storage_ontap_instance_storage_virtual_machine_volume_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this storage volume.",
			},
			"etag": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIbmIsStorageOntapInstanceVirtualMachineVolumeValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "storage_ontap_instance_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "storage_virtual_machine_id",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^[-0-9a-z_]+$`,
			MinValueLength:             1,
			MaxValueLength:             64,
		},
		validate.ValidateSchema{
			Identifier:                 "capacity",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "10",
			MaxValue:                   "16000",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[a-zA-Z][0-9a-zA-Z_]{0,202}$`,
			MinValueLength:             1,
			MaxValueLength:             203,
		},
		validate.ValidateSchema{
			Identifier:                 "security_style",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "mixed, none, unix, windows",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "storage_efficiency",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "disabled, enabled",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "data_protection, read_write",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_storage_ontap_instance_virtual_machine_volume", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	bodyModelMap := map[string]interface{}{}
	createStorageOntapInstanceStorageVirtualMachineVolumeOptions := &ontapv1.CreateStorageOntapInstanceStorageVirtualMachineVolumeOptions{}

	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	if _, ok := d.GetOk("capacity"); ok {
		bodyModelMap["capacity"] = d.Get("capacity")
	}
	if _, ok := d.GetOk("cifs_share"); ok {
		bodyModelMap["cifs_share"] = d.Get("cifs_share")
	}
	if _, ok := d.GetOk("enable_storage_efficiency"); ok {
		bodyModelMap["enable_storage_efficiency"] = d.Get("enable_storage_efficiency")
	}
	if _, ok := d.GetOk("export_policy"); ok {
		bodyModelMap["export_policy"] = d.Get("export_policy")
	}
	if _, ok := d.GetOk("security_style"); ok {
		bodyModelMap["security_style"] = d.Get("security_style")
	}
	if _, ok := d.GetOk("storage_efficiency"); ok {
		bodyModelMap["storage_efficiency"] = d.Get("storage_efficiency")
	}
	if _, ok := d.GetOk("type"); ok {
		bodyModelMap["type"] = d.Get("type")
	}
	createStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageOntapInstanceID(d.Get("storage_ontap_instance_id").(string))
	createStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageVirtualMachineID(d.Get("storage_virtual_machine_id").(string))
	convertedModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumePrototype(bodyModelMap)
	if err != nil {
		return diag.FromErr(err)
	}
	createStorageOntapInstanceStorageVirtualMachineVolumeOptions.StorageOntapInstanceStorageVirtualMachineVolumePrototype = convertedModel

	storageOntapInstanceStorageVirtualMachineVolume, response, err := ontapClient.CreateStorageOntapInstanceStorageVirtualMachineVolumeWithContext(context, createStorageOntapInstanceStorageVirtualMachineVolumeOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("CreateStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response))
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", *createStorageOntapInstanceStorageVirtualMachineVolumeOptions.StorageOntapInstanceID, *createStorageOntapInstanceStorageVirtualMachineVolumeOptions.StorageVirtualMachineID, *storageOntapInstanceStorageVirtualMachineVolume.ID))

	return resourceIbmIsStorageOntapInstanceVirtualMachineVolumeRead(context, d, meta)
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	getStorageOntapInstanceStorageVirtualMachineVolumeOptions := &ontapv1.GetStorageOntapInstanceStorageVirtualMachineVolumeOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageOntapInstanceID(parts[0])
	getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageVirtualMachineID(parts[1])
	getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetID(parts[2])

	storageOntapInstanceStorageVirtualMachineVolume, response, err := ontapClient.GetStorageOntapInstanceStorageVirtualMachineVolumeWithContext(context, getStorageOntapInstanceStorageVirtualMachineVolumeOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response))
	}

	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.Capacity) {
		if err = d.Set("capacity", flex.IntValue(storageOntapInstanceStorageVirtualMachineVolume.Capacity)); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting capacity: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.CifsShare) {
		cifsShareMap, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareToMap(storageOntapInstanceStorageVirtualMachineVolume.CifsShare)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("cifs_share", []map[string]interface{}{cifsShareMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting cifs_share: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.EnableStorageEfficiency) {
		if err = d.Set("enable_storage_efficiency", storageOntapInstanceStorageVirtualMachineVolume.EnableStorageEfficiency); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting enable_storage_efficiency: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.ExportPolicy) {
		exportPolicyMap, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyToMap(storageOntapInstanceStorageVirtualMachineVolume.ExportPolicy)
		if err != nil {
			return diag.FromErr(err)
		}
		if err = d.Set("export_policy", []map[string]interface{}{exportPolicyMap}); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting export_policy: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.Name) {
		if err = d.Set("name", storageOntapInstanceStorageVirtualMachineVolume.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.SecurityStyle) {
		if err = d.Set("security_style", storageOntapInstanceStorageVirtualMachineVolume.SecurityStyle); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting security_style: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.StorageEfficiency) {
		if err = d.Set("storage_efficiency", storageOntapInstanceStorageVirtualMachineVolume.StorageEfficiency); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting storage_efficiency: %s", err))
		}
	}
	if !core.IsNil(storageOntapInstanceStorageVirtualMachineVolume.Type) {
		if err = d.Set("type", storageOntapInstanceStorageVirtualMachineVolume.Type); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting type: %s", err))
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(storageOntapInstanceStorageVirtualMachineVolume.CreatedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	healthReasons := []map[string]interface{}{}
	for _, healthReasonsItem := range storageOntapInstanceStorageVirtualMachineVolume.HealthReasons {
		healthReasonsItemMap, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeHealthReasonToMap(&healthReasonsItem)
		if err != nil {
			return diag.FromErr(err)
		}
		healthReasons = append(healthReasons, healthReasonsItemMap)
	}
	if err = d.Set("health_reasons", healthReasons); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting health_reasons: %s", err))
	}
	if err = d.Set("health_state", storageOntapInstanceStorageVirtualMachineVolume.HealthState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting health_state: %s", err))
	}
	if err = d.Set("href", storageOntapInstanceStorageVirtualMachineVolume.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("junction_path", storageOntapInstanceStorageVirtualMachineVolume.JunctionPath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting junction_path: %s", err))
	}
	if err = d.Set("lifecycle_state", storageOntapInstanceStorageVirtualMachineVolume.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", storageOntapInstanceStorageVirtualMachineVolume.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("storage_ontap_instance_storage_virtual_machine_volume_id", storageOntapInstanceStorageVirtualMachineVolume.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting storage_ontap_instance_storage_virtual_machine_volume_id: %s", err))
	}
	if err = d.Set("etag", response.Headers.Get("Etag")); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting etag: %s", err))
	}

	return nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	updateStorageOntapInstanceStorageVirtualMachineVolumeOptions := &ontapv1.UpdateStorageOntapInstanceStorageVirtualMachineVolumeOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	updateStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageOntapInstanceID(parts[0])
	updateStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageVirtualMachineID(parts[1])
	updateStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetID(parts[2])

	hasChange := false

	patchVals := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumePatch{}
	if d.HasChange("storage_ontap_instance_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "storage_ontap_instance_id"))
	}
	if d.HasChange("storage_virtual_machine_id") {
		return diag.FromErr(fmt.Errorf("Cannot update resource property \"%s\" with the ForceNew annotation."+
			" The resource must be re-created to update this property.", "storage_virtual_machine_id"))
	}
	if d.HasChange("capacity") {
		newCapacity := int64(d.Get("capacity").(int))
		patchVals.Capacity = &newCapacity
		hasChange = true
	}
	if d.HasChange("cifs_share") {
		cifsShare, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePatch(d.Get("cifs_share.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		patchVals.CifsShare = cifsShare
		hasChange = true
	}
	if d.HasChange("export_policy") {
		exportPolicy, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPatch(d.Get("export_policy.0").(map[string]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}
		patchVals.ExportPolicy = exportPolicy
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("security_style") {
		newSecurityStyle := d.Get("security_style").(string)
		patchVals.SecurityStyle = &newSecurityStyle
		hasChange = true
	}
	updateStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetIfMatch(d.Get("etag").(string))

	if hasChange {
		updateStorageOntapInstanceStorageVirtualMachineVolumeOptions.StorageOntapInstanceStorageVirtualMachineVolumePatch, _ = patchVals.AsPatch()
		_, response, err := ontapClient.UpdateStorageOntapInstanceStorageVirtualMachineVolumeWithContext(context, updateStorageOntapInstanceStorageVirtualMachineVolumeOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("UpdateStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response))
		}
	}

	return resourceIbmIsStorageOntapInstanceVirtualMachineVolumeRead(context, d, meta)
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ontapClient, err := meta.(conns.ClientSession).OntapAPI()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteStorageOntapInstanceStorageVirtualMachineVolumeOptions := &ontapv1.DeleteStorageOntapInstanceStorageVirtualMachineVolumeOptions{}

	parts, err := flex.SepIdParts(d.Id(), "/")
	if err != nil {
		return diag.FromErr(err)
	}

	deleteStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageOntapInstanceID(parts[0])
	deleteStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageVirtualMachineID(parts[1])
	deleteStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetID(parts[2])

	_, response, err := ontapClient.DeleteStorageOntapInstanceStorageVirtualMachineVolumeWithContext(context, deleteStorageOntapInstanceStorageVirtualMachineVolumeOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("DeleteStorageOntapInstanceStorageVirtualMachineVolumeWithContext failed %s\n%s", err, response))
	}

	d.SetId("")

	return nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePrototype, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePrototype{}
	accessControlList := []ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlPrototype{}
	for _, accessControlListItem := range modelMap["access_control_list"].([]interface{}) {
		accessControlListItemModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlPrototype(accessControlListItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		accessControlList = append(accessControlList, *accessControlListItemModel)
	}
	model.AccessControlList = accessControlList
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlPrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlPrototype, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlPrototype{}
	model.Permission = core.StringPtr(modelMap["permission"].(string))
	users := []string{}
	for _, usersItem := range modelMap["users"].([]interface{}) {
		users = append(users, usersItem.(string))
	}
	model.Users = users
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPrototype, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPrototype{}
	if modelMap["rules"] != nil {
		rules := []ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRulePrototype{}
		for _, rulesItem := range modelMap["rules"].([]interface{}) {
			rulesItemModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRulePrototype(rulesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			rules = append(rules, *rulesItemModel)
		}
		model.Rules = rules
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRulePrototype(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRulePrototype, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRulePrototype{}
	if modelMap["access_control"] != nil && modelMap["access_control"].(string) != "" {
		model.AccessControl = core.StringPtr(modelMap["access_control"].(string))
	}
	clients := []ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientIntf{}
	for _, clientsItem := range modelMap["clients"].([]interface{}) {
		clientsItemModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClient(clientsItem.(map[string]interface{}))
		if err != nil {
			return model, err
		}
		clients = append(clients, clientsItemModel)
	}
	model.Clients = clients
	if modelMap["is_superuser"] != nil {
		model.IsSuperuser = core.BoolPtr(modelMap["is_superuser"].(bool))
	}
	if modelMap["nfs_version"] != nil {
		nfsVersion := []string{}
		for _, nfsVersionItem := range modelMap["nfs_version"].([]interface{}) {
			nfsVersion = append(nfsVersion, nfsVersionItem.(string))
		}
		model.NfsVersion = nfsVersion
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClient(modelMap map[string]interface{}) (ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientIntf, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClient{}
	if modelMap["hostname"] != nil && modelMap["hostname"].(string) != "" {
		model.Hostname = core.StringPtr(modelMap["hostname"].(string))
	}
	if modelMap["address"] != nil && modelMap["address"].(string) != "" {
		model.Address = core.StringPtr(modelMap["address"].(string))
	}
	if modelMap["cidr_block"] != nil && modelMap["cidr_block"].(string) != "" {
		model.CidrBlock = core.StringPtr(modelMap["cidr_block"].(string))
	}
	if modelMap["domain_suffix"] != nil && modelMap["domain_suffix"].(string) != "" {
		model.DomainSuffix = core.StringPtr(modelMap["domain_suffix"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname{}
	model.Hostname = core.StringPtr(modelMap["hostname"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP{}
	model.Address = core.StringPtr(modelMap["address"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR{}
	model.CidrBlock = core.StringPtr(modelMap["cidr_block"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName{}
	model.DomainSuffix = core.StringPtr(modelMap["domain_suffix"].(string))
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePatch(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePatch, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePatch{}
	if modelMap["access_control_list"] != nil {
		accessControlList := []ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlPrototype{}
		for _, accessControlListItem := range modelMap["access_control_list"].([]interface{}) {
			accessControlListItemModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlPrototype(accessControlListItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			accessControlList = append(accessControlList, *accessControlListItemModel)
		}
		model.AccessControlList = accessControlList
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPatch(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPatch, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPatch{}
	if modelMap["rules"] != nil {
		rules := []ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRulePrototype{}
		for _, rulesItem := range modelMap["rules"].([]interface{}) {
			rulesItemModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRulePrototype(rulesItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			rules = append(rules, *rulesItemModel)
		}
		model.Rules = rules
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumePrototype(modelMap map[string]interface{}) (ontapv1.StorageOntapInstanceStorageVirtualMachineVolumePrototypeIntf, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumePrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["capacity"] != nil {
		model.Capacity = core.Int64Ptr(int64(modelMap["capacity"].(int)))
	}
	if modelMap["cifs_share"] != nil && len(modelMap["cifs_share"].([]interface{})) > 0 {
		CifsShareModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePrototype(modelMap["cifs_share"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.CifsShare = CifsShareModel
	}
	if modelMap["enable_storage_efficiency"] != nil {
		model.EnableStorageEfficiency = core.BoolPtr(modelMap["enable_storage_efficiency"].(bool))
	}
	if modelMap["export_policy"] != nil && len(modelMap["export_policy"].([]interface{})) > 0 {
		ExportPolicyModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPrototype(modelMap["export_policy"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ExportPolicy = ExportPolicyModel
	}
	if modelMap["security_style"] != nil && modelMap["security_style"].(string) != "" {
		model.SecurityStyle = core.StringPtr(modelMap["security_style"].(string))
	}
	if modelMap["storage_efficiency"] != nil && modelMap["storage_efficiency"].(string) != "" {
		model.StorageEfficiency = core.StringPtr(modelMap["storage_efficiency"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumePrototypeStorageOntapInstanceStorageVirtualMachineVolumeByCapacity(modelMap map[string]interface{}) (*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumePrototypeStorageOntapInstanceStorageVirtualMachineVolumeByCapacity, error) {
	model := &ontapv1.StorageOntapInstanceStorageVirtualMachineVolumePrototypeStorageOntapInstanceStorageVirtualMachineVolumeByCapacity{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Capacity = core.Int64Ptr(int64(modelMap["capacity"].(int)))
	if modelMap["cifs_share"] != nil && len(modelMap["cifs_share"].([]interface{})) > 0 {
		CifsShareModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeCIFSSharePrototype(modelMap["cifs_share"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.CifsShare = CifsShareModel
	}
	if modelMap["enable_storage_efficiency"] != nil {
		model.EnableStorageEfficiency = core.BoolPtr(modelMap["enable_storage_efficiency"].(bool))
	}
	if modelMap["export_policy"] != nil && len(modelMap["export_policy"].([]interface{})) > 0 {
		ExportPolicyModel, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeMapToStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyPrototype(modelMap["export_policy"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ExportPolicy = ExportPolicyModel
	}
	if modelMap["security_style"] != nil && modelMap["security_style"].(string) != "" {
		model.SecurityStyle = core.StringPtr(modelMap["security_style"].(string))
	}
	if modelMap["storage_efficiency"] != nil && modelMap["storage_efficiency"].(string) != "" {
		model.StorageEfficiency = core.StringPtr(modelMap["storage_efficiency"].(string))
	}
	if modelMap["type"] != nil && modelMap["type"].(string) != "" {
		model.Type = core.StringPtr(modelMap["type"].(string))
	}
	return model, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShare) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	accessControlList := []map[string]interface{}{}
	for _, accessControlListItem := range model.AccessControlList {
		accessControlListItemMap, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlToMap(&accessControlListItem)
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

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControlToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeCIFSShareAccessControl) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["permission"] = model.Permission
	modelMap["users"] = model.Users
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicy) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["mount_path"] = model.MountPath
	if model.Rules != nil {
		rules := []map[string]interface{}{}
		for _, rulesItem := range model.Rules {
			rulesItemMap, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleToMap(&rulesItem)
			if err != nil {
				return modelMap, err
			}
			rules = append(rules, rulesItemMap)
		}
		modelMap["rules"] = rules
	}
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRule) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["access_control"] = model.AccessControl
	clients := []map[string]interface{}{}
	for _, clientsItem := range model.Clients {
		clientsItemMap, err := resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientToMap(clientsItem)
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

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientToMap(model ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientIntf) (map[string]interface{}, error) {
	if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname); ok {
		return resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostnameToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname))
	} else if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP); ok {
		return resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIPToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP))
	} else if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR); ok {
		return resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDRToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR))
	} else if _, ok := model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName); ok {
		return resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainNameToMap(model.(*ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName))
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

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostnameToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByHostname) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["hostname"] = model.Hostname
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIPToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByIP) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = model.Address
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDRToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByCIDR) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["cidr_block"] = model.CidrBlock
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainNameToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeExportPolicyRuleClientByDomainName) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["domain_suffix"] = model.DomainSuffix
	return modelMap, nil
}

func resourceIbmIsStorageOntapInstanceVirtualMachineVolumeStorageOntapInstanceStorageVirtualMachineVolumeHealthReasonToMap(model *ontapv1.StorageOntapInstanceStorageVirtualMachineVolumeHealthReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = model.Code
	modelMap["message"] = model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = model.MoreInfo
	}
	return modelMap, nil
}
