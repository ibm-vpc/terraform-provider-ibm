// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBMSVolAttBareMetalServer          = "bare_metal_server"
	isBMSVolAttId                       = "volume_attachment_id"
	isBMSVolAttBandwidth                = "bandwidth"
	isBMSVolAttCreatedAt                = "created_at"
	isBMSVolAttDeleteOnServerDelete     = "delete_volume_on_bare_metal_server_delete"
	isBMSVolAttDeleteOnAttachmentDelete = "delete_volume_on_attachment_delete"
	isBMSVolAttDevice                   = "device"
	isBMSVolAttHref                     = "href"
	isBMSVolAttName                     = "name"
	isBMSVolAttStatus                   = "status"
	isBMSVolAttStatusReason             = "status_reason"
	isBMSVolAttType                     = "type"

	isBMSVolAttVol               = "volume"
	isBMSVolAttVolName           = "volume_name"
	isBMSVolAttVolCRN            = "volume_crn"
	isBMSVolAttVolHref           = "volume_href"
	isBMSVolAttVolDeleted        = "volume_deleted"
	isBMSVolAttResourceGroup     = "resource_group"
	isBMSVolAttUserTags          = "user_tags"
	isBMSVolAttAllowedUse        = "allowed_use"
	isBMSVolAttSourceSnapshot    = "source_snapshot"
	isBMSVolAttSourceSnapshotCrn = "source_snapshot_crn"

	isBMSVolAttProtocol          = "protocol"
	isBMSVolAttNvmeQualifiedName = "nvme_qualified_name"
	isBMSVolAttIps               = "ips"
	isBMSVolAttCapacity          = "capacity"
	isBMSVolAttIops              = "iops"
	isBMSVolAttProfile           = "profile"

	isBMSVolAttEncryptionKey  = "encryption_key"
	isBMSVolAttAttachmentMode = "attachment_mode"
)

func ResourceIBMISBareMetalServerVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMisBMSVolumeAttachmentCreate,
		ReadContext:   resourceIBMisBMSVolumeAttachmentRead,
		UpdateContext: resourceIBMisBMSVolumeAttachmentUpdate,
		DeleteContext: resourceIBMisBMSVolumeAttachmentDelete,
		Exists:        resourceIBMisBMSVolumeAttachmentExists,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceVolumeValidate(diff)
				}),
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				}),
		),

		Schema: map[string]*schema.Schema{
			isBMSVolAttBareMetalServer: {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttBareMetalServer),
				Description:  "The bare metal server ID",
			},
			isBMSVolAttId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this volume attachment",
			},
			isBMSVolAttName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttName),
				Description:  "The name for this volume attachment",
			},
			isBMSVolAttVol: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{isBMSVolAttCapacity, isBMSVolAttProfile, isBMSVolAttIops, isBMSVolAttEncryptionKey, isBMSVolAttAttachmentMode, isBMSVolAttResourceGroup, isBMSVolAttUserTags, isBMSVolAttAllowedUse, isBMSVolAttSourceSnapshot, isBMSVolAttBandwidth, isBMSVolAttVolName},
				ValidateFunc:  validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttVol),
				Description:   "The ID of an existing volume to attach",
			},
			isBMSVolAttCapacity: {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				AtLeastOneOf:  []string{isBMSVolAttVol, isBMSVolAttCapacity, isBMSVolAttSourceSnapshot},
				ValidateFunc:  validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttCapacity),
				Description:   "The capacity of the volume in gigabytes",
			},
			isBMSVolAttBandwidth: {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				ValidateFunc:  validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttBandwidth),
				Description:   "The maximum bandwidth (in megabits per second) for the volume",
			},
			isBMSVolAttProfile: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				ValidateFunc:  validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttProfile),
				Description:   "The profile name for the volume. Bare metal server volume attachments support only the sdp profile.",
			},
			isBMSVolAttIops: {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				Description:   "The maximum I/O operations per second for the volume",
			},
			isBMSVolAttEncryptionKey: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				Description:   "The CRN of the encryption key for the volume",
			},
			isBMSVolAttAttachmentMode: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				ValidateFunc:  validate.ValidateAllowedStringValues([]string{"single", "multiple"}),
				Description:   "The attachment mode of the volume: single or multiple",
			},
			isBMSVolAttResourceGroup: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				Description:   "The resource group ID for the new volume",
			},
			isBMSVolAttUserTags: {
				Type:          schema.TypeSet,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				Elem:          &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttUserTags)},
				Set:           flex.ResourceIBMVPCHash,
				Description:   "The user tags associated with the volume",
			},
			isBMSVolAttAllowedUse: {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{isBMSVolAttVol},
				Description:   "The usage constraints to be matched against requested instance or bare metal server properties to determine compatibility.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_version": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", "allowed_use.api_version"),
							Description:  "The API version with which to evaluate the expressions.",
						},
						"bare_metal_server": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", "allowed_use.bare_metal_server"),
							Description:  "The expression that must be satisfied by the properties of a bare metal server provisioned using this volume.",
						},
						"instance": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ForceNew:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", "allowed_use.instance"),
							Description:  "The expression that must be satisfied by the properties of a virtual server instance provisioned using this volume.",
						},
					},
				},
			},
			isBMSVolAttSourceSnapshot: {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				AtLeastOneOf:  []string{isBMSVolAttVol, isBMSVolAttCapacity, isBMSVolAttSourceSnapshot},
				ConflictsWith: []string{isBMSVolAttVol},
				Description:   "The snapshot ID to use as the source for the new volume",
			},
			isBMSVolAttSourceSnapshotCrn: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the source snapshot for the attached volume",
			},
			isBMSVolAttDeleteOnServerDelete: {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "If true, deleting the bare metal server also deletes the attached volume",
			},
			isBMSVolAttDeleteOnAttachmentDelete: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "If true, deleting the attachment also deletes the volume. Default is true.",
			},
			isBMSVolAttProtocol: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.ValidateAllowedStringValues([]string{"nvme_tcp"}),
				Description:  "The protocol used for this volume attachment (nvme_tcp)",
			},
			isBMSVolAttNvmeQualifiedName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The NVMe Qualified Name (NQN) of the subsystem",
			},
			isBMSVolAttIps: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The IP addresses for connecting to the volume using nvme_tcp",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			isBMSVolAttVolName: {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttVolName),
				Description:  "The name of the attached volume",
			},
			isBMSVolAttVolCRN: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN of the attached volume",
			},
			isBMSVolAttVolHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL of the attached volume",
			},
			isBMSVolAttVolDeleted: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Link to documentation about deleted resources",
			},
			isBMSVolAttDevice: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier for the device exposed to the bare metal server OS",
			},
			isBMSVolAttHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this volume attachment",
			},
			isBMSVolAttCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the volume attachment was created",
			},
			isBMSVolAttStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this volume attachment: attaching, available, detaching, unusable",
			},
			isBMSVolAttStatusReason: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any)",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status reason code",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason",
						},
						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A link to documentation about this status reason",
						},
					},
				},
			},
			isBMSVolAttType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of volume attachment (data)",
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceIBMISBMSVolumeAttachmentValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttBareMetalServer,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
		},
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttVol,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
		},
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttCapacity,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			MinValue:                   "10",
			MaxValue:                   "64000",
		},
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttProfile,
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Optional:                   true,
			AllowedValues:              "sdp",
		},
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttBandwidth,
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "1000",
			MaxValue:                   "8192",
		},
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttUserTags,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "allowed_use.api_version",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12]\d|3[01])$`,
		},
		validate.ValidateSchema{
			Identifier:                 "allowed_use.bare_metal_server",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-zA-Z_][a-zA-Z0-9_]*|[-+*/%]|&&|\|\||!|==|!=|<|<=|>|>=|~|\bin\b|\(|\)|\[|\]|,|\.|"|'|"|'|\s+|\d+)+$`,
		},
		validate.ValidateSchema{
			Identifier:                 "allowed_use.instance",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-zA-Z_][a-zA-Z0-9_]*|[-+*/%]|&&|\|\||!|==|!=|<|<=|>|>=|~|\bin\b|\(|\)|\[|\]|,|\.|"|'|"|'|\s+|\d+)+$`,
		},
		validate.ValidateSchema{
			Identifier:                 isBMSVolAttVolName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)
	validator := validate.ResourceValidator{ResourceName: "ibm_is_bare_metal_server_volume_attachment", Schema: validateSchema}
	return &validator
}

func resourceIBMisBMSVolumeAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bmsId := d.Get(isBMSVolAttBareMetalServer).(string)
	err := bmsVolAttachmentCreate(context, d, meta, bmsId)
	if err != nil {
		return err
	}
	return resourceIBMisBMSVolumeAttachmentRead(context, d, meta)
}

func bmsVolAttachmentCreate(context context.Context, d *schema.ResourceData, meta interface{}, bmsId string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	proto := &vpcv1.BareMetalServerVolumeAttachmentPrototype{}

	volumeIdStr := ""
	if volumeId, ok := d.GetOk(isBMSVolAttVol); ok {
		volumeIdStr = volumeId.(string)
	}

	if volumeIdStr != "" {
		// attach existing volume by ID
		proto.Volume = &vpcv1.BareMetalServerVolumeAttachmentPrototypeVolumeVolumeIdentity{
			ID: &volumeIdStr,
		}
	} else {
		// create a new volume inline
		newVol := &vpcv1.BareMetalServerVolumeAttachmentPrototypeVolumeVolumePrototypeBareMetalServerContext{}

		if volName, ok := d.GetOk(isBMSVolAttVolName); ok {
			volNameStr := volName.(string)
			newVol.Name = &volNameStr
		}

		// resolve snapshot minimum capacity, mirroring instance behaviour
		volSnapshotStr := ""
		if volSnapshot, ok := d.GetOk(isBMSVolAttSourceSnapshot); ok {
			volSnapshotStr = volSnapshot.(string)
			newVol.SourceSnapshot = &vpcv1.SnapshotIdentity{ID: &volSnapshotStr}
		}

		var snapCapacity int64
		if volSnapshotStr != "" {
			snapshotGet, _, err := sess.GetSnapshotWithContext(context, &vpcv1.GetSnapshotOptions{
				ID: &volSnapshotStr,
			})
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetSnapshotWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			snapCapacity = int64(int(*snapshotGet.MinimumCapacity))
		}

		var volCapacityInt int64
		if volCapacity, ok := d.GetOk(isBMSVolAttCapacity); ok {
			volCapacityInt = int64(volCapacity.(int))
			if volCapacityInt != 0 && volCapacityInt > snapCapacity {
				newVol.Capacity = &volCapacityInt
			}
		}

		if volBandwidth, ok := d.GetOk(isBMSVolAttBandwidth); ok {
			bw := int64(volBandwidth.(int))
			if bw != 0 {
				newVol.Bandwidth = &bw
			}
		}
		if volIops, ok := d.GetOk(isBMSVolAttIops); ok {
			iops := int64(volIops.(int))
			if iops != 0 {
				newVol.Iops = &iops
			}
			volProfileStr := d.Get(isBMSVolAttProfile).(string)
			if volProfileStr == "" {
				volProfileStr = "sdp"
			}
			newVol.Profile = &vpcv1.VolumeProfileIdentity{Name: &volProfileStr}
		} else {
			volProfileStr := "sdp"
			if volProfile, ok := d.GetOk(isBMSVolAttProfile); ok {
				volProfileStr = volProfile.(string)
			}
			newVol.Profile = &vpcv1.VolumeProfileIdentity{Name: &volProfileStr}
		}

		if encKey, ok := d.GetOk(isBMSVolAttEncryptionKey); ok {
			encKeyStr := encKey.(string)
			newVol.EncryptionKey = &vpcv1.EncryptionKeyIdentity{CRN: &encKeyStr}
		}
		if attMode, ok := d.GetOk(isBMSVolAttAttachmentMode); ok {
			attModeStr := attMode.(string)
			newVol.AttachmentMode = &attModeStr
		}
		if resourceGroupID, ok := d.GetOk(isBMSVolAttResourceGroup); ok {
			rgIDStr := resourceGroupID.(string)
			newVol.ResourceGroup = &vpcv1.ResourceGroupIdentity{ID: &rgIDStr}
		}
		if userTags, ok := d.GetOk(isBMSVolAttUserTags); ok {
			userTagsSet := userTags.(*schema.Set)
			if userTagsSet != nil && userTagsSet.Len() != 0 {
				userTagsArray := make([]string, userTagsSet.Len())
				for i, userTag := range userTagsSet.List() {
					userTagsArray[i] = userTag.(string)
				}
				schematicTags := os.Getenv("IC_ENV_TAGS")
				var envTags []string
				if schematicTags != "" {
					envTags = strings.Split(schematicTags, ",")
					userTagsArray = append(userTagsArray, envTags...)
				}
				newVol.UserTags = userTagsArray
			}
		}
		if allowedUse, ok := d.GetOk(isBMSVolAttAllowedUse); ok && len(allowedUse.([]interface{})) > 0 {
			allowedUseModel, _ := ResourceIBMIsVolumeAllowedUseMapToVolumeAllowedUsePrototype(allowedUse.([]interface{})[0].(map[string]interface{}))
			newVol.AllowedUse = allowedUseModel
		}
		proto.Volume = newVol
	}

	if autoDelete, ok := d.GetOk(isBMSVolAttDeleteOnServerDelete); ok {
		b := autoDelete.(bool)
		proto.DeleteVolumeOnBareMetalServerDelete = &b
	}
	if name, ok := d.GetOk(isBMSVolAttName); ok {
		nameStr := name.(string)
		proto.Name = &nameStr
	}
	if protocol, ok := d.GetOk(isBMSVolAttProtocol); ok {
		protocolStr := protocol.(string)
		proto.Protocol = &protocolStr
	}

	isBMSKey := "bms_key_" + bmsId
	conns.IbmMutexKV.Lock(isBMSKey)
	defer conns.IbmMutexKV.Unlock(isBMSKey)

	createOpts := &vpcv1.CreateBareMetalServerVolumeAttachmentOptions{
		BareMetalServerID:                        &bmsId,
		BareMetalServerVolumeAttachmentPrototype: proto,
	}
	volAttIntf, _, err := sess.CreateBareMetalServerVolumeAttachmentWithContext(context, createOpts)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateBareMetalServerVolumeAttachmentWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	volAtt, ok := volAttIntf.(*vpcv1.BareMetalServerVolumeAttachment)
	if !ok || volAtt == nil {
		tfErr := flex.TerraformErrorf(fmt.Errorf("unexpected type from CreateBareMetalServerVolumeAttachment"), "unexpected type", "ibm_is_bare_metal_server_volume_attachment", "create")
		return tfErr.GetDiag()
	}

	d.SetId(makeTerraformVolAttID(bmsId, *volAtt.ID))
	volAttResult, err := isWaitForBMSVolumeAttached(sess, d, bmsId, *volAtt.ID)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForBMSVolumeAttached failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	v := os.Getenv("IC_ENV_TAGS")
	if _, ok := d.GetOk(isBMSVolAttUserTags); ok || v != "" {
		volAttRef := volAttResult.(*vpcv1.BareMetalServerVolumeAttachment)
		if volAttRef != nil && volAttRef.Volume != nil {
			oldList, newList := d.GetChange(isBMSVolAttUserTags)
			err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *volAttRef.Volume.CRN, "", isInstanceUserTagType)
			if err != nil {
				log.Printf("Error on create of resource bare metal server volume attachment (%s) tags: %s", d.Id(), err)
			}
		}
	}

	log.Printf("[INFO] Bare Metal Server (%s) volume attachment: %s", bmsId, *volAtt.ID)
	return nil
}

func resourceIBMisBMSVolumeAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bmsId, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "sep-id-parts").GetDiag()
	}
	return bmsVolumeAttachmentGet(context, d, meta, bmsId, id)
}

func bmsVolumeAttachmentGet(context context.Context, d *schema.ResourceData, meta interface{}, bmsId, id string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getOpts := &vpcv1.GetBareMetalServerVolumeAttachmentOptions{
		BareMetalServerID: &bmsId,
		ID:                &id,
	}
	volAttIntf, response, err := sess.GetBareMetalServerVolumeAttachmentWithContext(context, getOpts)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBareMetalServerVolumeAttachmentWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	volAtt, ok := volAttIntf.(*vpcv1.BareMetalServerVolumeAttachment)
	if !ok || volAtt == nil {
		d.SetId("")
		return nil
	}

	if err = d.Set(isBMSVolAttBareMetalServer, bmsId); err != nil {
		err = fmt.Errorf("Error setting bare_metal_server: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-bare_metal_server").GetDiag()
	}
	if err = d.Set(isBMSVolAttId, *volAtt.ID); err != nil {
		err = fmt.Errorf("Error setting volume_attachment_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-volume_attachment_id").GetDiag()
	}
	if err = d.Set(isBMSVolAttName, *volAtt.Name); err != nil {
		err = fmt.Errorf("Error setting name: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-name").GetDiag()
	}
	if err = d.Set(isBMSVolAttHref, *volAtt.Href); err != nil {
		err = fmt.Errorf("Error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-href").GetDiag()
	}
	if err = d.Set(isBMSVolAttStatus, *volAtt.Status); err != nil {
		err = fmt.Errorf("Error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-status").GetDiag()
	}
	if err = d.Set(isBMSVolAttType, *volAtt.Type); err != nil {
		err = fmt.Errorf("Error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-type").GetDiag()
	}
	if err = d.Set(isBMSVolAttDeleteOnServerDelete, *volAtt.DeleteVolumeOnBareMetalServerDelete); err != nil {
		err = fmt.Errorf("Error setting delete_volume_on_bare_metal_server_delete: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-delete_volume_on_bare_metal_server_delete").GetDiag()
	}
	if volAtt.CreatedAt != nil {
		if err = d.Set(isBMSVolAttCreatedAt, volAtt.CreatedAt.String()); err != nil {
			err = fmt.Errorf("Error setting created_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-created_at").GetDiag()
		}
	}
	if volAtt.Bandwidth != nil {
		if err = d.Set(isBMSVolAttBandwidth, int(*volAtt.Bandwidth)); err != nil {
			err = fmt.Errorf("Error setting bandwidth: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-bandwidth").GetDiag()
		}
	}
	if volAtt.Device != nil {
		if err = d.Set(isBMSVolAttDevice, *volAtt.Device.ID); err != nil {
			err = fmt.Errorf("Error setting device: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-device").GetDiag()
		}
	}
	if volAtt.Protocol != nil {
		if err = d.Set(isBMSVolAttProtocol, *volAtt.Protocol); err != nil {
			err = fmt.Errorf("Error setting protocol: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-protocol").GetDiag()
		}
	}
	if volAtt.NvmeQualifiedName != nil {
		if err = d.Set(isBMSVolAttNvmeQualifiedName, *volAtt.NvmeQualifiedName); err != nil {
			err = fmt.Errorf("Error setting nvme_qualified_name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-nvme_qualified_name").GetDiag()
		}
	}
	if len(volAtt.Ips) > 0 {
		ips := make([]string, len(volAtt.Ips))
		for i, ip := range volAtt.Ips {
			ips[i] = *ip.Address
		}
		if err = d.Set(isBMSVolAttIps, ips); err != nil {
			err = fmt.Errorf("Error setting ips: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-ips").GetDiag()
		}
	}
	if volAtt.StatusReason != nil {
		statusReasonMap := map[string]interface{}{
			"code":    *volAtt.StatusReason.Code,
			"message": *volAtt.StatusReason.Message,
		}
		if volAtt.StatusReason.MoreInfo != nil {
			statusReasonMap["more_info"] = *volAtt.StatusReason.MoreInfo
		}
		if err = d.Set(isBMSVolAttStatusReason, []map[string]interface{}{statusReasonMap}); err != nil {
			err = fmt.Errorf("Error setting status_reason: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-status_reason").GetDiag()
		}
	}

	if volAtt.Volume != nil {
		if err = d.Set(isBMSVolAttVol, *volAtt.Volume.ID); err != nil {
			err = fmt.Errorf("Error setting volume: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-volume").GetDiag()
		}
		if err = d.Set(isBMSVolAttVolName, *volAtt.Volume.Name); err != nil {
			err = fmt.Errorf("Error setting volume_name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-volume_name").GetDiag()
		}
		if err = d.Set(isBMSVolAttVolCRN, *volAtt.Volume.CRN); err != nil {
			err = fmt.Errorf("Error setting volume_crn: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-volume_crn").GetDiag()
		}
		if err = d.Set(isBMSVolAttVolHref, *volAtt.Volume.Href); err != nil {
			err = fmt.Errorf("Error setting volume_href: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-volume_href").GetDiag()
		}
		if volAtt.Volume.Deleted != nil {
			if err = d.Set(isBMSVolAttVolDeleted, *volAtt.Volume.Deleted.MoreInfo); err != nil {
				err = fmt.Errorf("Error setting volume_deleted: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-volume_deleted").GetDiag()
			}
		}

		// second API call: GetVolume for volume-level fields (mirrors instance behaviour)
		volId := *volAtt.Volume.ID
		volumeDetail, response, err := sess.GetVolumeWithContext(context, &vpcv1.GetVolumeOptions{ID: &volId})
		if err != nil || volumeDetail == nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "read")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if err = d.Set(isBMSVolAttVol, *volumeDetail.ID); err != nil {
			err = fmt.Errorf("Error setting volume: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-volume").GetDiag()
		}
		if err = d.Set(isBMSVolAttIops, *volumeDetail.Iops); err != nil {
			err = fmt.Errorf("Error setting iops: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-iops").GetDiag()
		}
		if err = d.Set(isBMSVolAttProfile, *volumeDetail.Profile.Name); err != nil {
			err = fmt.Errorf("Error setting profile: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-profile").GetDiag()
		}
		if err = d.Set(isBMSVolAttCapacity, *volumeDetail.Capacity); err != nil {
			err = fmt.Errorf("Error setting capacity: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-capacity").GetDiag()
		}
		if volumeDetail.AttachmentMode != nil {
			if err = d.Set(isBMSVolAttAttachmentMode, *volumeDetail.AttachmentMode); err != nil {
				err = fmt.Errorf("Error setting attachment_mode: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-attachment_mode").GetDiag()
			}
		}
		if err = d.Set(isBMSVolAttBandwidth, volumeDetail.Bandwidth); err != nil {
			err = fmt.Errorf("Error setting bandwidth: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-bandwidth").GetDiag()
		}
		if volumeDetail.EncryptionKey != nil {
			if err = d.Set(isBMSVolAttEncryptionKey, *volumeDetail.EncryptionKey.CRN); err != nil {
				err = fmt.Errorf("Error setting encryption_key: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-encryption_key").GetDiag()
			}
		}
		if volumeDetail.SourceSnapshot != nil {
			if err = d.Set(isBMSVolAttSourceSnapshot, *volumeDetail.SourceSnapshot.ID); err != nil {
				err = fmt.Errorf("Error setting source_snapshot: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-source_snapshot").GetDiag()
			}
			if err = d.Set(isBMSVolAttSourceSnapshotCrn, *volumeDetail.SourceSnapshot.CRN); err != nil {
				err = fmt.Errorf("Error setting source_snapshot_crn: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-source_snapshot_crn").GetDiag()
			}
		}
		if volumeDetail.ResourceGroup != nil {
			if err = d.Set(isBMSVolAttResourceGroup, *volumeDetail.ResourceGroup.ID); err != nil {
				err = fmt.Errorf("Error setting resource_group: %s", err)
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-resource_group").GetDiag()
			}
		}
		allowedUses := []map[string]interface{}{}
		if volumeDetail.AllowedUse != nil {
			modelMap, err := ResourceceIBMIsVolumeAllowedUseToMap(volumeDetail.AllowedUse)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read")
				log.Println(tfErr.GetDiag())
			} else {
				allowedUses = append(allowedUses, modelMap)
			}
		}
		if err = d.Set(isBMSVolAttAllowedUse, allowedUses); err != nil {
			err = fmt.Errorf("Error setting allowed_use: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-allowed_use").GetDiag()
		}
		if err = d.Set("version", response.Headers.Get("Etag")); err != nil {
			err = fmt.Errorf("Error setting version: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-version").GetDiag()
		}
		userTags, err := flex.GetGlobalTagsUsingCRN(meta, *volumeDetail.CRN, "", isInstanceUserTagType)
		if err != nil {
			log.Printf("Error on get of resource bare metal server volume attachment (%s) tags: %s", d.Id(), err)
		}
		if err = d.Set(isBMSVolAttUserTags, userTags); err != nil {
			err = fmt.Errorf("Error setting user_tags: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "read", "set-user_tags").GetDiag()
		}
	}

	return nil
}

func bmsVolAttUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	bmsId, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "update", "sep-id-parts").GetDiag()
	}
	if volumeCRNOk, ok := d.GetOk(isBMSVolAttVolCRN); ok {
		volumeCRN := volumeCRNOk.(string)
		if d.HasChange(isBMSVolAttUserTags) {
			oldList, newList := d.GetChange(isBMSVolAttUserTags)
			err = flex.UpdateTagsUsingCRN(oldList, newList, meta, volumeCRN)
			if err != nil {
				log.Printf("Error on update of resource bare metal server volume attachment (%s) tags: %s", d.Id(), err)
			}
		}
	}
	flag := false
	patchModel := &vpcv1.BareMetalServerVolumeAttachmentPatch{}
	if d.HasChange(isBMSVolAttDeleteOnServerDelete) {
		b := d.Get(isBMSVolAttDeleteOnServerDelete).(bool)
		patchModel.DeleteVolumeOnBareMetalServerDelete = &b
		flag = true
	}
	if d.HasChange(isBMSVolAttName) {
		name := d.Get(isBMSVolAttName).(string)
		patchModel.Name = &name
		flag = true
	}
	if flag {
		patch, err := patchModel.AsPatch()
		if err != nil || patch == nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("BareMetalServerVolumeAttachmentPatch.AsPatch() failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		updateOpts := &vpcv1.UpdateBareMetalServerVolumeAttachmentOptions{
			BareMetalServerID:                    &bmsId,
			ID:                                   &id,
			BareMetalServerVolumeAttachmentPatch: patch,
		}
		_, _, err = sess.UpdateBareMetalServerVolumeAttachmentWithContext(context, updateOpts)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateBareMetalServerVolumeAttachmentWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	// update volume_name and bandwidth via UpdateVolume
	hasVolNameChanged := d.HasChange(isBMSVolAttVolName)
	hasBandwidthChanged := d.HasChange(isBMSVolAttBandwidth)
	volId := ""
	if volIdOk, ok := d.GetOk(isBMSVolAttVol); ok {
		volId = volIdOk.(string)
	}
	if volId != "" && (hasVolNameChanged || hasBandwidthChanged) {
		voloptions := &vpcv1.UpdateVolumeOptions{ID: &volId}
		if hasVolNameChanged {
			volumeNamePatchModel := &vpcv1.VolumePatch{}
			newname := d.Get(isBMSVolAttVolName).(string)
			volumeNamePatchModel.Name = &newname
			volumePatch, err := volumeNamePatchModel.AsPatch()
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("volumeNamePatchModel.AsPatch() failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			voloptions.VolumePatch = volumePatch
			_, _, err = sess.UpdateVolumeWithContext(context, voloptions)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
		if hasBandwidthChanged {
			volumeBandwidthPatchModel := &vpcv1.VolumePatch{}
			newBandwidth := int64(d.Get(isBMSVolAttBandwidth).(int))
			volumeBandwidthPatchModel.Bandwidth = &newBandwidth
			volumePatch, err := volumeBandwidthPatchModel.AsPatch()
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("volumeBandwidthPatchModel.AsPatch() failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
			voloptions.VolumePatch = volumePatch
			_, _, err = sess.UpdateVolumeWithContext(context, voloptions)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
	}

	// update iops/profile via UpdateVolume (sdp profile can update iops in-place)
	volProfile := ""
	if volProfileOk, ok := d.GetOk(isBMSVolAttProfile); ok {
		volProfile = volProfileOk.(string)
	}
	if volId != "" && d.HasChange(isBMSVolAttIops) && !d.HasChange(isBMSVolAttProfile) && volProfile == "sdp" {
		updateVolumeProfileOptions := &vpcv1.UpdateVolumeOptions{
			ID: &volId,
		}
		volumeProfilePatchModel := &vpcv1.VolumePatch{}
		if d.HasChange(isBMSVolAttIops) {
			iops := int64(d.Get(isBMSVolAttIops).(int))
			volumeProfilePatchModel.Iops = &iops
		}
		volumeProfilePatch, err := volumeProfilePatchModel.AsPatch()
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("volumeProfilePatchModel.AsPatch() failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		optionsget := &vpcv1.GetVolumeOptions{
			ID: &volId,
		}
		_, response, err := sess.GetVolumeWithContext(context, optionsget)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		eTag := response.Headers.Get("ETag")
		updateVolumeProfileOptions.IfMatch = &eTag
		updateVolumeProfileOptions.VolumePatch = volumeProfilePatch
		_, _, err = sess.UpdateVolumeWithContext(context, updateVolumeProfileOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		_, err = isWaitForVolumeAvailable(sess, volId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForVolumeAvailable failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	} else if volId != "" && (d.HasChange(isBMSVolAttIops) || d.HasChange(isBMSVolAttProfile)) {
		updateVolumeProfileOptions := &vpcv1.UpdateVolumeOptions{
			ID: &volId,
		}
		volumeProfilePatchModel := &vpcv1.VolumePatch{}
		if d.HasChange(isBMSVolAttProfile) {
			profile := d.Get(isBMSVolAttProfile).(string)
			volumeProfilePatchModel.Profile = &vpcv1.VolumeProfileIdentity{
				Name: &profile,
			}
		} else if d.HasChange(isBMSVolAttIops) {
			profile := d.Get(isBMSVolAttProfile).(string)
			volumeProfilePatchModel.Profile = &vpcv1.VolumeProfileIdentity{
				Name: &profile,
			}
			iops := int64(d.Get(isBMSVolAttIops).(int))
			volumeProfilePatchModel.Iops = &iops
		}
		volumeProfilePatch, err := volumeProfilePatchModel.AsPatch()
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("volumeProfilePatchModel.AsPatch() failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		optionsget := &vpcv1.GetVolumeOptions{
			ID: &volId,
		}
		_, response, err := sess.GetVolumeWithContext(context, optionsget)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		eTag := response.Headers.Get("ETag")
		updateVolumeProfileOptions.IfMatch = &eTag
		updateVolumeProfileOptions.VolumePatch = volumeProfilePatch
		_, _, err = sess.UpdateVolumeWithContext(context, updateVolumeProfileOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		_, err = isWaitForVolumeAvailable(sess, volId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForVolumeAvailable failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	// update capacity via UpdateVolume
	if volId != "" && d.HasChange(isBMSVolAttCapacity) {
		getvolumeoptions := &vpcv1.GetVolumeOptions{ID: &volId}
		vol, _, err := sess.GetVolumeWithContext(context, getvolumeoptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		// for non-sdp profiles the instance must be running to resize; sdp can resize in place
		if *vol.Profile.Name != "sdp" {
			var firstAttachment *vpcv1.VolumeAttachmentReferenceVolumeContext
			if len(vol.VolumeAttachments) > 0 {
				firstAttachment, _ = vol.VolumeAttachments[0].(*vpcv1.VolumeAttachmentReferenceVolumeContext)
			}
			if firstAttachment == nil || firstAttachment.Name == nil || *firstAttachment.Name == "" {
				err = fmt.Errorf("Error volume capacity can't be updated since volume %s is not attached to any server for VolumePatch", id)
				tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "update")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		}
		capacity := int64(d.Get(isBMSVolAttCapacity).(int))
		updateVolumeOptions := &vpcv1.UpdateVolumeOptions{ID: &volId}
		volumeCapacityPatchModel := &vpcv1.VolumePatch{}
		volumeCapacityPatchModel.Capacity = &capacity
		volumeCapacityPatch, err := volumeCapacityPatchModel.AsPatch()
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("volumeCapacityPatchModel.AsPatch() failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		updateVolumeOptions.VolumePatch = volumeCapacityPatch
		_, _, err = sess.UpdateVolumeWithContext(context, updateVolumeOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		_, err = isWaitForVolumeAvailable(sess, volId, d.Timeout(schema.TimeoutCreate))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForVolumeAvailable failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	return nil
}

func resourceIBMisBMSVolumeAttachmentUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	err := bmsVolAttUpdate(context, d, meta)
	if err != nil {
		return err
	}
	return resourceIBMisBMSVolumeAttachmentRead(context, d, meta)
}

func bmsVolAttDelete(context context.Context, d *schema.ResourceData, meta interface{}, bmsId, id, volId string, volDelete bool) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteBmsVolAttOptions := &vpcv1.DeleteBareMetalServerVolumeAttachmentOptions{
		BareMetalServerID: &bmsId,
		ID:                &id,
	}

	isBMSKey := "bms_key_" + bmsId
	conns.IbmMutexKV.Lock(isBMSKey)
	defer conns.IbmMutexKV.Unlock(isBMSKey)

	_, err = sess.DeleteBareMetalServerVolumeAttachmentWithContext(context, deleteBmsVolAttOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteBareMetalServerVolumeAttachmentWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	_, err = isWaitForBMSVolumeDetached(sess, d, bmsId, id)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForBMSVolumeDetached failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if volDelete {
		deleteVolumeOptions := &vpcv1.DeleteVolumeOptions{
			ID: &volId,
		}
		_, err := sess.DeleteVolumeWithContext(context, deleteVolumeOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteVolumeWithContext failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "delete")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		_, err = isWaitForVolumeDeleted(sess, volId, d.Timeout(schema.TimeoutDelete))
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForVolumeDeleted failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "delete")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	return nil
}
func resourceIBMisBMSVolumeAttachmentDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bmsId, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "delete", "sep-id-parts").GetDiag()
	}

	volDelete := false
	if volDeleteOk, ok := d.GetOk(isBMSVolAttDeleteOnAttachmentDelete); ok {
		volDelete = volDeleteOk.(bool)
	}
	volId := ""
	if volIdOk, ok := d.GetOk(isBMSVolAttVol); ok {
		volId = volIdOk.(string)
	}

	diagErr := bmsVolAttDelete(context, d, meta, bmsId, id, volId, volDelete)
	if diagErr != nil {
		return diagErr
	}
	d.SetId("")
	return nil
}

func resourceIBMisBMSVolumeAttachmentExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	bmsId, id, err := parseVolAttTerraformID(d.Id())
	if err != nil {
		return false, flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "exists", "sep-id-parts")
	}
	exists, err := bmsVolAttExists(d, meta, bmsId, id)
	return exists, err
}

func bmsVolAttExists(d *schema.ResourceData, meta interface{}, bmsId, id string) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_bare_metal_server_volume_attachment", "exists", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return false, tfErr
	}

	getOpts := &vpcv1.GetBareMetalServerVolumeAttachmentOptions{
		BareMetalServerID: &bmsId,
		ID:                &id,
	}
	_, response, err := sess.GetBareMetalServerVolumeAttachment(getOpts)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetBareMetalServerVolumeAttachment failed: %s", err.Error()), "ibm_is_bare_metal_server_volume_attachment", "exists")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return false, tfErr
	}
	return true, nil
}

func isWaitForBMSVolumeAttached(sess *vpcv1.VpcV1, d *schema.ResourceData, bmsId, volAttId string) (interface{}, error) {
	log.Printf("Waiting for bare metal server (%s) volume attachment (%s) to be attached.", bmsId, volAttId)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"attaching"},
		Target:     []string{"available"},
		Refresh:    isBMSVolumeRefreshFunc(sess, bmsId, volAttId),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}

func isBMSVolumeRefreshFunc(sess *vpcv1.VpcV1, bmsId, volAttId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		getOpts := &vpcv1.GetBareMetalServerVolumeAttachmentOptions{
			BareMetalServerID: &bmsId,
			ID:                &volAttId,
		}
		volAttIntf, response, err := sess.GetBareMetalServerVolumeAttachment(getOpts)
		if err != nil {
			return nil, "", fmt.Errorf("[ERROR] Error getting BMS volume attachment: %s\n%s", err, response)
		}
		volAtt, ok := volAttIntf.(*vpcv1.BareMetalServerVolumeAttachment)
		if !ok || volAtt == nil {
			return nil, "", fmt.Errorf("[ERROR] unexpected type from GetBareMetalServerVolumeAttachment")
		}
		switch *volAtt.Status {
		case "available":
			return volAtt, "available", nil
		case "unusable":
			return volAtt, "", fmt.Errorf("[ERROR] BMS volume attachment (%s) is in unusable state", volAttId)
		}
		return volAtt, "attaching", nil
	}
}

func isWaitForBMSVolumeDetached(sess *vpcv1.VpcV1, d *schema.ResourceData, bmsId, volAttId string) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"available", "detaching"},
		Target:  []string{isInstanceDeleteDone, ""},
		Refresh: func() (interface{}, string, error) {
			getOpts := &vpcv1.GetBareMetalServerVolumeAttachmentOptions{
				BareMetalServerID: &bmsId,
				ID:                &volAttId,
			}
			volAttIntf, response, err := sess.GetBareMetalServerVolumeAttachment(getOpts)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return volAttIntf, isInstanceDeleteDone, nil
				}
				return nil, "", fmt.Errorf("[ERROR] Error detaching BMS volume: %s\n%s", err, response)
			}
			return volAttIntf, "detaching", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForState()
}
