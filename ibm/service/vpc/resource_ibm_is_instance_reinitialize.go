// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isInstanceReinitializeInstance               = "instance"
	isInstanceReinitializeImage                  = "image"
	isInstanceReinitializeBootVolume             = "boot_volume"
	isInstanceReinitializeBootSnapshot           = "boot_snapshot"
	isInstanceReinitializeKeys                   = "keys"
	isInstanceReinitializeUserData               = "user_data"
	isInstanceReinitializeDefaultTrustedProfile  = "default_trusted_profile"
	isInstanceReinitializeTarget                 = "target"
	isInstanceReinitializeAutoLink               = "auto_link"
	isInstanceReinitializeBootVolumeAttachment   = "boot_volume_attachment"
	isInstanceReinitializeBootVolumeName         = "name"
	isInstanceReinitializeBootVolumeCapacity     = "capacity"
	isInstanceReinitializeBootVolumeProfile      = "profile"
	isInstanceReinitializeBootVolumeEncryption   = "encryption_key"
	isInstanceReinitializeStatus                 = "status"
	isInstanceReinitializeBootVolumeAttachmentID = "boot_volume_attachment_id"
	isInstanceReinitializeAutoStop               = "auto_stop"
	isInstanceReinitializeAutoStart              = "auto_start"
	isInstanceReinitialized                      = "triggers"
	isInstanceReinitializedAt                    = "reinitialized_at"
)

func ResourceIBMISInstanceReinitialize() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISInstanceReinitializeCreate,
		ReadContext:   resourceIBMISInstanceReinitializeRead,
		UpdateContext: resourceIBMISInstanceReinitializeUpdate,
		DeleteContext: resourceIBMISInstanceReinitializeDelete,
		Exists:        resourceIBMISInstanceReinitializeExists,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			isInstanceReinitializeInstance: {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Instance identifier",
			},
			isInstanceReinitializeImage: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{isInstanceReinitializeBootVolume, isInstanceReinitializeBootSnapshot},
				Description:   "Image ID to reinitialize with",
			},
			isInstanceReinitializeBootVolume: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{isInstanceReinitializeImage, isInstanceReinitializeBootSnapshot},
				Description:   "Existing boot volume ID",
			},
			isInstanceReinitializeBootSnapshot: {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{isInstanceReinitializeImage, isInstanceReinitializeBootVolume},
				Description:   "Snapshot ID to create boot volume from",
			},
			isInstanceReinitializeKeys: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "SSH key IDs",
			},
			isInstanceReinitializeUserData: {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "User data for cloud-init",
			},
			isInstanceReinitializeDefaultTrustedProfile: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Default trusted profile configuration",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceReinitializeTarget: {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Trusted profile ID",
						},
						isInstanceReinitializeAutoLink: {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Auto-link to instance",
						},
					},
				},
			},
			isInstanceReinitializeBootVolumeAttachment: {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Boot volume configuration when using image",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isInstanceReinitializeBootVolumeName: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Boot volume name",
						},
						isInstanceReinitializeBootVolumeCapacity: {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "Boot volume capacity in GB",
						},
						isInstanceReinitializeBootVolumeProfile: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Boot volume profile",
						},
						isInstanceReinitializeBootVolumeEncryption: {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Encryption key CRN",
						},
					},
				},
			},
			isInstanceReinitializeAutoStop: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Automatically stop instance before reinitialization if it's running",
			},
			isInstanceReinitializeAutoStart: {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Automatically start instance after reinitialization (default behavior)",
			},
			isInstanceReinitialized: {
				Type:        schema.TypeMap,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Arbitrary map of values that trigger reinitialization when changed",
			},
			// Computed fields
			isInstanceReinitializeStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Instance status after reinitialization",
			},
			isInstanceReinitializeBootVolumeAttachmentID: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "New boot volume attachment ID after reinitialization",
			},
			isInstanceReinitializedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp of last reinitialization",
			},
		},
	}
}

func ResourceIBMISInstanceReinitializeValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)

	ibmISInstanceReinitializeResourceValidator := validate.ResourceValidator{
		ResourceName: "ibm_is_instance_reinitialize",
		Schema:       validateSchema,
	}
	return &ibmISInstanceReinitializeResourceValidator
}

func resourceIBMISInstanceReinitializeCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_reinitialize", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceID := d.Get(isInstanceReinitializeInstance).(string)
	autoStop := d.Get(isInstanceReinitializeAutoStop).(bool)

	// Check instance exists
	getInstanceOptions := &vpcv1.GetInstanceOptions{
		ID: &instanceID,
	}
	instance, response, err := sess.GetInstanceWithContext(context, getInstanceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceWithContext failed: %s", err.Error()), "ibm_is_instance_reinitialize", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Handle auto_stop if instance is running
	if *instance.Status != isInstanceActionStatusStopped {
		if autoStop {
			log.Printf("[INFO] Instance is running, auto_stop enabled - stopping instance %s", instanceID)
			stopAction := "stop"
			createStopOptions := &vpcv1.CreateInstanceActionOptions{
				InstanceID: &instanceID,
				Type:       &stopAction,
			}
			_, response, err = sess.CreateInstanceActionWithContext(context, createStopOptions)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Failed to stop instance: %s\n%s", err.Error(), response), "ibm_is_instance_reinitialize", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}

			// Wait for instance to stop
			_, err = isWaitForInstanceActionStop(sess, d.Timeout(schema.TimeoutCreate), instanceID, d)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceActionStop failed: %s", err.Error()), "ibm_is_instance_reinitialize", "create")
				log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
				return tfErr.GetDiag()
			}
		} else {
			err = fmt.Errorf("[ERROR] Instance must be stopped before reinitialization. Current status: %s. Set auto_stop=true to automatically stop the instance", *instance.Status)
			tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_instance_reinitialize", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	// Build reinitialization prototype
	reinitPrototype, err := buildReinitializePrototype(d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Failed to build reinitialize prototype: %s", err.Error()), "ibm_is_instance_reinitialize", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Call reinitialize API
	reinitOptions := &vpcv1.ReinitializeInstanceOptions{
		ID:                            &instanceID,
		InstanceReinitializePrototype: reinitPrototype,
	}

	response, err = sess.ReinitializeInstanceWithContext(context, reinitOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ReinitializeInstanceWithContext failed: %s\n%s", err.Error(), response), "ibm_is_instance_reinitialize", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Wait for instance to start
	_, err = isWaitForInstanceActionStart(sess, d.Timeout(schema.TimeoutCreate), instanceID, d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceActionStart failed: %s", err.Error()), "ibm_is_instance_reinitialize", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(instanceID)

	// Set reinitialized_at timestamp
	if err = d.Set(isInstanceReinitializedAt, time.Now().UTC().Format(time.RFC3339)); err != nil {
		log.Printf("[WARN] Error setting reinitialized_at: %s", err)
	}

	return resourceIBMISInstanceReinitializeRead(context, d, meta)
}

func resourceIBMISInstanceReinitializeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_reinitialize", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	id := d.Id()
	getInstanceOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	instance, response, err := sess.GetInstanceWithContext(context, getInstanceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceWithContext failed: %s", err.Error()), "ibm_is_instance_reinitialize", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set(isInstanceReinitializeStatus, *instance.Status); err != nil {
		err = fmt.Errorf("Error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_reinitialize", "read", "set-status").GetDiag()
	}

	// Get boot volume attachment ID
	if instance.BootVolumeAttachment != nil && instance.BootVolumeAttachment.ID != nil {
		if err = d.Set(isInstanceReinitializeBootVolumeAttachmentID, *instance.BootVolumeAttachment.ID); err != nil {
			err = fmt.Errorf("Error setting boot_volume_attachment_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_reinitialize", "read", "set-boot_volume_attachment_id").GetDiag()
		}
	}

	return nil
}

func resourceIBMISInstanceReinitializeUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Reinitialization is a ForceNew operation, so this should trigger recreation
	return resourceIBMISInstanceReinitializeCreate(context, d, meta)
}

func resourceIBMISInstanceReinitializeDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// No-op: reinitialization resource doesn't need cleanup
	d.SetId("")
	return nil
}

func resourceIBMISInstanceReinitializeExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_reinitialize", "exists", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return false, tfErr
	}

	id := d.Id()
	getInstanceOptions := &vpcv1.GetInstanceOptions{
		ID: &id,
	}
	_, response, err := sess.GetInstance(getInstanceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			return false, nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstance failed: %s", err.Error()), "ibm_is_instance_reinitialize", "exists")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return false, fmt.Errorf("[ERROR] Error getting instance: %s\n%s", err, response)
	}
	return true, nil
}

// Helper function to build reinitialization prototype
func buildReinitializePrototype(d *schema.ResourceData) (vpcv1.InstanceReinitializePrototypeIntf, error) {
	// Determine which type of reinitialization based on user input

	// Handle image-based reinitialization
	if imageID, ok := d.GetOk(isInstanceReinitializeImage); ok {
		imageIdentity := &vpcv1.ImageIdentity{
			ID: core.StringPtr(imageID.(string)),
		}

		prototype := &vpcv1.InstanceReinitializePrototypeInstanceReinitializeByImage{
			Image: imageIdentity,
		}

		// Handle boot volume attachment configuration
		if bootVolConfig, ok := d.GetOk(isInstanceReinitializeBootVolumeAttachment); ok {
			bootVolList := bootVolConfig.([]interface{})
			if len(bootVolList) > 0 {
				bootVolMap := bootVolList[0].(map[string]interface{})

				// Create volume prototype with required profile
				volPrototype := &vpcv1.VolumePrototypeInstanceByImageContext{
					Profile: &vpcv1.VolumeProfileIdentity{
						Name: core.StringPtr("general-purpose"), // Default profile
					},
				}

				if name, ok := bootVolMap[isInstanceReinitializeBootVolumeName].(string); ok && name != "" {
					volPrototype.Name = core.StringPtr(name)
				}
				if capacity, ok := bootVolMap[isInstanceReinitializeBootVolumeCapacity].(int); ok && capacity > 0 {
					volPrototype.Capacity = core.Int64Ptr(int64(capacity))
				}
				if profile, ok := bootVolMap[isInstanceReinitializeBootVolumeProfile].(string); ok && profile != "" {
					volPrototype.Profile = &vpcv1.VolumeProfileIdentity{
						Name: core.StringPtr(profile),
					}
				}
				if encKey, ok := bootVolMap[isInstanceReinitializeBootVolumeEncryption].(string); ok && encKey != "" {
					volPrototype.EncryptionKey = &vpcv1.EncryptionKeyIdentity{
						CRN: core.StringPtr(encKey),
					}
				}

				prototype.BootVolumeAttachment = &vpcv1.VolumeAttachmentPrototypeInstanceByImageContext{
					Volume: volPrototype,
				}
			}
		}

		// Add common fields
		addCommonFields(d, prototype)
		return prototype, nil
	}

	// Handle existing volume-based reinitialization
	if volumeID, ok := d.GetOk(isInstanceReinitializeBootVolume); ok {
		volumeIdentity := &vpcv1.VolumeIdentity{
			ID: core.StringPtr(volumeID.(string)),
		}

		volumeAttachment := &vpcv1.VolumeAttachmentPrototypeInstanceByVolumeContext{
			Volume: volumeIdentity,
		}

		prototype := &vpcv1.InstanceReinitializePrototypeInstanceReinitializeByVolume{
			BootVolumeAttachment: volumeAttachment,
		}

		// Add common fields
		addCommonFields(d, prototype)
		return prototype, nil
	}

	// Handle snapshot-based reinitialization
	if snapshotID, ok := d.GetOk(isInstanceReinitializeBootSnapshot); ok {
		snapshotIdentity := &vpcv1.SnapshotIdentity{
			ID: core.StringPtr(snapshotID.(string)),
		}

		// Create volume prototype with required profile and snapshot
		volPrototype := &vpcv1.VolumePrototypeInstanceBySourceSnapshotContext{
			Profile: &vpcv1.VolumeProfileIdentity{
				Name: core.StringPtr("general-purpose"), // Default profile
			},
			SourceSnapshot: snapshotIdentity,
		}

		volumeAttachment := &vpcv1.VolumeAttachmentPrototypeInstanceBySourceSnapshotContext{
			Volume: volPrototype,
		}

		prototype := &vpcv1.InstanceReinitializePrototypeInstanceReinitializeBySnapshot{
			BootVolumeAttachment: volumeAttachment,
		}

		// Add common fields
		addCommonFields(d, prototype)
		return prototype, nil
	}

	return nil, fmt.Errorf("must specify one of: image, boot_volume, or boot_snapshot")
}

// Helper to add common fields to any prototype type
func addCommonFields(d *schema.ResourceData, prototype interface{}) {
	// Handle SSH keys
	if keys, ok := d.GetOk(isInstanceReinitializeKeys); ok {
		keyList := keys.([]interface{})
		keyIdentities := make([]vpcv1.KeyIdentityIntf, len(keyList))
		for i, key := range keyList {
			keyIdentities[i] = &vpcv1.KeyIdentity{
				ID: core.StringPtr(key.(string)),
			}
		}

		switch p := prototype.(type) {
		case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeByImage:
			p.Keys = keyIdentities
		case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeByVolume:
			p.Keys = keyIdentities
		case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeBySnapshot:
			p.Keys = keyIdentities
		}
	}

	// Handle user data
	if userData, ok := d.GetOk(isInstanceReinitializeUserData); ok {
		userDataStr := core.StringPtr(userData.(string))
		switch p := prototype.(type) {
		case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeByImage:
			p.UserData = userDataStr
		case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeByVolume:
			p.UserData = userDataStr
		case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeBySnapshot:
			p.UserData = userDataStr
		}
	}

	// Handle default trusted profile
	if trustedProfile, ok := d.GetOk(isInstanceReinitializeDefaultTrustedProfile); ok {
		tpList := trustedProfile.([]interface{})
		if len(tpList) > 0 {
			tpMap := tpList[0].(map[string]interface{})
			profileProto := &vpcv1.InstanceDefaultTrustedProfilePrototype{}

			if target, ok := tpMap[isInstanceReinitializeTarget].(string); ok && target != "" {
				profileProto.Target = &vpcv1.TrustedProfileIdentity{
					ID: core.StringPtr(target),
				}
			}
			if autoLink, ok := tpMap[isInstanceReinitializeAutoLink].(bool); ok {
				profileProto.AutoLink = core.BoolPtr(autoLink)
			}

			switch p := prototype.(type) {
			case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeByImage:
				p.DefaultTrustedProfile = profileProto
			case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeByVolume:
				p.DefaultTrustedProfile = profileProto
			case *vpcv1.InstanceReinitializePrototypeInstanceReinitializeBySnapshot:
				p.DefaultTrustedProfile = profileProto
			}
		}
	}
}

// Made with Bob
