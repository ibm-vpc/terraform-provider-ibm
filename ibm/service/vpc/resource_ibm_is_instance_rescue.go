// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 99-SNAPSHOT-c913c18c-20260623-001551
 */

package vpc

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMISInstanceRescue() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMISInstanceRescueCreate,
		ReadContext:   resourceIBMISInstanceRescueRead,
		DeleteContext: resourceIBMISInstanceRescueDelete,
		Importer:      &schema.ResourceImporter{},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_instance_rescue", "instance_id"),
				Description:  "The virtual server instance identifier.",
			},
			"image": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "The image to use for rescuing the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CRN for this image.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL for this image.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unique identifier for this image.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this image. The name is unique across all images in the region.",
						},
						"remote": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates that the resource associated with this referenceis remote and therefore may not be directly retrievable.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"account": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisaccount, and identifies the owning account.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this account.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"region": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "If present, this property indicates that the referenced resource is remote to thisregion, and identifies the native region.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this region.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this region.",
												},
											},
										},
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"keys": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Description: "The public SSH keys used at initialization for the rescue volume.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The unique identifier for this key.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The CRN for this key.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The URL for this key.",
						},
						"fingerprint": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The fingerprint for this key.",
						},
					},
				},
			},
			"user_data": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "User data to be made available when setting up the rescue volume.",
			},
			"lifecycle_status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle status of the instance (e.g., 'rescue', 'stable').",
			},
			"rescue_volume_attachment": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				ForceNew:    true,
				Description: "The rescue volume attachment for this instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delete_volume_on_instance_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     true,
							ForceNew:    true,
							Description: "Indicates whether to delete the rescue volume when the instance is deleted.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The name for this volume attachment. The name is unique across all volume attachments on the instance.",
						},
						"volume": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The volume prototype or reference for the rescue volume.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The unique identifier for this volume.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: "The name for this volume.",
									},
									"profile": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The profile name for this volume.",
									},
									"capacity": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The capacity of the volume in gigabytes.",
									},
									"iops": &schema.Schema{
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum I/O operations per second (IOPS) for the volume.",
									},
									"encryption_key": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The CRN of the encryption key to use for this volume.",
									},
									"resource_group": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The resource group ID for this volume.",
									},
									"user_tags": &schema.Schema{
										Type:        schema.TypeSet,
										Optional:    true,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Description: "User tags for this volume.",
									},
								},
							},
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A link to documentation about deleted resources.",
									},
								},
							},
						},
						"device": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The configuration for the volume as a device in the instance operating system.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique identifier for the device.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this volume attachment.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this volume attachment.",
						},
					},
				},
			},
			"password": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"encrypted_password": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The administrator password at rescue, encrypted using `encryption_key`, and returned base64-encoded.",
						},
						"encryption_key": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The public SSH key used to encrypt the administrator password.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this key.",
									},
									"deleted": &schema.Schema{
										Type:        schema.TypeList,
										Optional:    true,
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"more_info": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "A link to documentation about deleted resources.",
												},
											},
										},
									},
									"fingerprint": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The fingerprint for this key.  The value is returned base64-encoded and prefixed with the hash algorithm (such as `SHA256`).The length of this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this key.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this key.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this key. The name is unique across all keys in the region.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISInstanceRescueValidator() *validate.ResourceValidator {
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
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_rescue", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMISInstanceRescueCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceID := d.Get("instance_id").(string)

	// Check if instance exists
	getInstanceOptions := vpcClient.NewGetInstanceOptions(instanceID)
	instance, _, err := vpcClient.GetInstanceWithContext(context, getInstanceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceWithContext failed: %s", err.Error()), "ibm_is_instance_rescue", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Stop instance if it's not already stopped
	if *instance.Status != "stopped" {
		log.Printf("[INFO] Instance %s is in status '%s'. Stopping instance before rescue operation.", instanceID, *instance.Status)

		actionType := "stop"
		createInstanceActionOptions := &vpcv1.CreateInstanceActionOptions{
			InstanceID: &instanceID,
			Type:       &actionType,
		}

		_, response, err := vpcClient.CreateInstanceActionWithContext(context, createInstanceActionOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateInstanceActionWithContext (stop) failed: %s\n%s", err.Error(), response), "ibm_is_instance_rescue", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		// Wait for instance to stop
		log.Printf("[INFO] Waiting for instance %s to stop...", instanceID)
		_, err = isWaitForInstanceActionStop(vpcClient, d.Timeout(schema.TimeoutCreate), instanceID, d)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceActionStop failed: %s", err.Error()), "ibm_is_instance_rescue", "create")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		log.Printf("[INFO] Instance %s successfully stopped", instanceID)
	}

	bodyModelMap := map[string]interface{}{}
	if _, ok := d.GetOk("keys"); ok {
		bodyModelMap["keys"] = d.Get("keys")
	}
	if _, ok := d.GetOk("user_data"); ok {
		bodyModelMap["user_data"] = d.Get("user_data")
	}
	if _, ok := d.GetOk("image"); ok {
		bodyModelMap["image"] = d.Get("image")
	}
	if _, ok := d.GetOk("rescue_volume_attachment"); ok {
		bodyModelMap["rescue_volume_attachment"] = d.Get("rescue_volume_attachment")
	}

	convertedModel, err := ResourceIBMISInstanceRescueMapToInstanceRescuePrototype(bodyModelMap)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "create", "parse-request-body").GetDiag()
	}

	createInstanceRescueOptions := &vpcv1.CreateInstanceRescueOptions{
		InstanceID:              &instanceID,
		InstanceRescuePrototype: convertedModel,
	}

	_, _, err = vpcClient.CreateInstanceRescueWithContext(context, createInstanceRescueOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateInstanceRescueWithContext failed: %s", err.Error()), "ibm_is_instance_rescue", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Wait for instance to start (API automatically starts the instance after rescue is created)
	log.Printf("[INFO] Waiting for instance %s to start after rescue creation...", instanceID)
	_, err = isWaitForInstanceActionStart(vpcClient, d.Timeout(schema.TimeoutCreate), instanceID, d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceActionStart failed: %s", err.Error()), "ibm_is_instance_rescue", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Wait for instance to enter rescue mode
	log.Printf("[INFO] Waiting for instance %s to enter rescue mode...", instanceID)
	_, err = isWaitForInstanceRescueMode(vpcClient, d.Timeout(schema.TimeoutCreate), instanceID, d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceRescueMode failed: %s", err.Error()), "ibm_is_instance_rescue", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	log.Printf("[INFO] Instance %s successfully entered rescue mode and is running", instanceID)

	d.SetId(instanceID)

	return resourceIBMISInstanceRescueRead(context, d, meta)
}

func resourceIBMISInstanceRescueRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceID := d.Id()

	getInstanceRescueOptions := &vpcv1.GetInstanceRescueOptions{
		InstanceID: &instanceID,
	}

	instanceRescue, response, err := vpcClient.GetInstanceRescueWithContext(context, getInstanceRescueOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceRescueWithContext failed: %s", err.Error()), "ibm_is_instance_rescue", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("instance_id", d.Id()); err != nil {
		err = fmt.Errorf("Error setting instance_id: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "set-instance_id").GetDiag()
	}
	if !core.IsNil(instanceRescue.Image) {
		imageMap, err := ResourceIBMISInstanceRescueImageReferenceToMap(instanceRescue.Image)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "image-to-map").GetDiag()
		}
		if err = d.Set("image", []map[string]interface{}{imageMap}); err != nil {
			err = fmt.Errorf("Error setting image: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "set-image").GetDiag()
		}
	}
	if !core.IsNil(instanceRescue.Keys) {
		keys := []map[string]interface{}{}
		for _, keysItem := range instanceRescue.Keys {
			keysItemMap, err := ResourceIBMISInstanceRescueKeyReferenceToMap(&keysItem) // #nosec G601
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "keys-to-map").GetDiag()
			}
			keys = append(keys, keysItemMap)
		}
		if err = d.Set("keys", keys); err != nil {
			err = fmt.Errorf("Error setting keys: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "set-keys").GetDiag()
		}
	}
	if !core.IsNil(instanceRescue.RescueVolumeAttachment) {
		rescueVolumeAttachmentMap, err := ResourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(instanceRescue.RescueVolumeAttachment)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "rescue_volume_attachment-to-map").GetDiag()
		}
		if err = d.Set("rescue_volume_attachment", []map[string]interface{}{rescueVolumeAttachmentMap}); err != nil {
			err = fmt.Errorf("Error setting rescue_volume_attachment: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "set-rescue_volume_attachment").GetDiag()
		}
	}
	if !core.IsNil(instanceRescue.Password) {
		passwordMap, err := ResourceIBMISInstanceRescueInstanceRescuePasswordToMap(instanceRescue.Password)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "password-to-map").GetDiag()
		}
		if err = d.Set("password", []map[string]interface{}{passwordMap}); err != nil {
			err = fmt.Errorf("Error setting password: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "set-password").GetDiag()
		}
	}

	// Get instance to retrieve lifecycle_status and user_data
	getInstanceOptions := vpcClient.NewGetInstanceOptions(instanceID)
	instance, _, err := vpcClient.GetInstanceWithContext(context, getInstanceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceWithContext failed: %s", err.Error()), "ibm_is_instance_rescue", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if !core.IsNil(instance.LifecycleState) {
		if err = d.Set("lifecycle_status", *instance.LifecycleState); err != nil {
			err = fmt.Errorf("Error setting lifecycle_status: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "read", "set-lifecycle_status").GetDiag()
		}
	}

	return nil
}

func resourceIBMISInstanceRescueDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceID := d.Id()

	// Verify instance exists and get current status
	getInstanceOptions := vpcClient.NewGetInstanceOptions(instanceID)
	instance, _, err := vpcClient.GetInstanceWithContext(context, getInstanceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceWithContext failed: %s", err.Error()), "ibm_is_instance_rescue", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Check if instance is in rescue mode
	if instance.LifecycleState != nil && *instance.LifecycleState != "rescue" {
		log.Printf("[WARN] Instance %s is not in rescue mode (current state: %s), skipping rescue delete", instanceID, *instance.LifecycleState)
		d.SetId("")
		return nil
	}

	// Stop instance if it's not already stopped (required by DELETE API)
	if *instance.Status != "stopped" {
		log.Printf("[INFO] Instance %s is in status '%s'. Stopping instance before exiting rescue mode.", instanceID, *instance.Status)

		actionType := "stop"
		createInstanceActionOptions := &vpcv1.CreateInstanceActionOptions{
			InstanceID: &instanceID,
			Type:       &actionType,
		}

		_, response, err := vpcClient.CreateInstanceActionWithContext(context, createInstanceActionOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateInstanceActionWithContext (stop) failed: %s\n%s", err.Error(), response), "ibm_is_instance_rescue", "delete")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}

		// Wait for instance to stop
		log.Printf("[INFO] Waiting for instance %s to stop...", instanceID)
		_, err = isWaitForInstanceActionStop(vpcClient, d.Timeout(schema.TimeoutDelete), instanceID, d)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceActionStop failed: %s", err.Error()), "ibm_is_instance_rescue", "delete")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		log.Printf("[INFO] Instance %s successfully stopped", instanceID)
	}

	deleteInstanceRescueOptions := &vpcv1.DeleteInstanceRescueOptions{}
	deleteInstanceRescueOptions.SetInstanceID(instanceID)

	_, err = vpcClient.DeleteInstanceRescueWithContext(context, deleteInstanceRescueOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteInstanceRescueWithContext failed: %s", err.Error()), "ibm_is_instance_rescue", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Wait for instance to exit rescue mode and return to stopped state
	_, err = isWaitForInstanceExitRescue(vpcClient, d.Timeout(schema.TimeoutDelete), instanceID, d)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("isWaitForInstanceExitRescue failed: %s", err.Error()), "ibm_is_instance_rescue", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Verify instance returned to stopped state
	instance, _, err = vpcClient.GetInstanceWithContext(context, getInstanceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceWithContext failed after rescue exit: %s", err.Error()), "ibm_is_instance_rescue", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if *instance.Status != "stopped" {
		err = fmt.Errorf("[ERROR] Instance did not return to stopped state after exiting rescue mode. Current status: %s", *instance.Status)
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_instance_rescue", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMISInstanceRescueMapToKeyIdentity(modelMap map[string]interface{}) (vpcv1.KeyIdentityIntf, error) {
	model := &vpcv1.KeyIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	if modelMap["fingerprint"] != nil && modelMap["fingerprint"].(string) != "" {
		model.Fingerprint = core.StringPtr(modelMap["fingerprint"].(string))
	}
	return model, nil
}

func ResourceIBMISInstanceRescueMapToKeyIdentityByID(modelMap map[string]interface{}) (*vpcv1.KeyIdentityByID, error) {
	model := &vpcv1.KeyIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToKeyIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.KeyIdentityByCRN, error) {
	model := &vpcv1.KeyIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToKeyIdentityByHref(modelMap map[string]interface{}) (*vpcv1.KeyIdentityByHref, error) {
	model := &vpcv1.KeyIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToKeyIdentityByFingerprint(modelMap map[string]interface{}) (*vpcv1.KeyIdentityByFingerprint, error) {
	model := &vpcv1.KeyIdentityByFingerprint{}
	model.Fingerprint = core.StringPtr(modelMap["fingerprint"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToImageIdentity(modelMap map[string]interface{}) (vpcv1.ImageIdentityIntf, error) {
	model := &vpcv1.ImageIdentity{}
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

func ResourceIBMISInstanceRescueMapToImageIdentityByID(modelMap map[string]interface{}) (*vpcv1.ImageIdentityByID, error) {
	model := &vpcv1.ImageIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToImageIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.ImageIdentityByCRN, error) {
	model := &vpcv1.ImageIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToImageIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ImageIdentityByHref, error) {
	model := &vpcv1.ImageIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToInstanceRescueVolumeAttachmentPrototype(modelMap map[string]interface{}) (*vpcv1.InstanceRescueVolumeAttachmentPrototype, error) {
	model := &vpcv1.InstanceRescueVolumeAttachmentPrototype{}
	if modelMap["delete_volume_on_instance_delete"] != nil {
		model.DeleteVolumeOnInstanceDelete = core.BoolPtr(modelMap["delete_volume_on_instance_delete"].(bool))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	VolumeModel, err := ResourceIBMISInstanceRescueMapToVolumePrototypeInstanceByImageContext(modelMap["volume"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Volume = VolumeModel
	return model, nil
}

func ResourceIBMISInstanceRescueMapToVolumePrototypeInstanceByImageContext(modelMap map[string]interface{}) (*vpcv1.VolumePrototypeInstanceByImageContext, error) {
	model := &vpcv1.VolumePrototypeInstanceByImageContext{}
	if modelMap["allowed_use"] != nil && len(modelMap["allowed_use"].([]interface{})) > 0 {
		AllowedUseModel, err := ResourceIBMISInstanceRescueMapToVolumeAllowedUsePrototype(modelMap["allowed_use"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.AllowedUse = AllowedUseModel
	}
	if modelMap["bandwidth"] != nil {
		model.Bandwidth = core.Int64Ptr(int64(modelMap["bandwidth"].(int)))
	}
	if modelMap["capacity"] != nil {
		model.Capacity = core.Int64Ptr(int64(modelMap["capacity"].(int)))
	}
	if modelMap["encryption_key"] != nil && len(modelMap["encryption_key"].([]interface{})) > 0 {
		EncryptionKeyModel, err := ResourceIBMISInstanceRescueMapToEncryptionKeyIdentity(modelMap["encryption_key"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.EncryptionKey = EncryptionKeyModel
	}
	if modelMap["iops"] != nil {
		model.Iops = core.Int64Ptr(int64(modelMap["iops"].(int)))
	}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	ProfileModel, err := ResourceIBMISInstanceRescueMapToVolumeProfileIdentity(modelMap["profile"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Profile = ProfileModel
	if modelMap["resource_group"] != nil && len(modelMap["resource_group"].([]interface{})) > 0 {
		ResourceGroupModel, err := ResourceIBMISInstanceRescueMapToResourceGroupIdentity(modelMap["resource_group"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceGroup = ResourceGroupModel
	}
	if modelMap["user_tags"] != nil {
		userTags := []string{}
		for _, userTagsItem := range modelMap["user_tags"].([]interface{}) {
			userTags = append(userTags, userTagsItem.(string))
		}
		model.UserTags = userTags
	}
	return model, nil
}

func ResourceIBMISInstanceRescueMapToVolumeAllowedUsePrototype(modelMap map[string]interface{}) (*vpcv1.VolumeAllowedUsePrototype, error) {
	model := &vpcv1.VolumeAllowedUsePrototype{}
	if modelMap["api_version"] != nil && modelMap["api_version"].(string) != "" {
		model.ApiVersion = core.StringPtr(modelMap["api_version"].(string))
	}
	if modelMap["bare_metal_server"] != nil && modelMap["bare_metal_server"].(string) != "" {
		model.BareMetalServer = core.StringPtr(modelMap["bare_metal_server"].(string))
	}
	if modelMap["instance"] != nil && modelMap["instance"].(string) != "" {
		model.Instance = core.StringPtr(modelMap["instance"].(string))
	}
	return model, nil
}

func ResourceIBMISInstanceRescueMapToEncryptionKeyIdentity(modelMap map[string]interface{}) (vpcv1.EncryptionKeyIdentityIntf, error) {
	model := &vpcv1.EncryptionKeyIdentity{}
	if modelMap["crn"] != nil && modelMap["crn"].(string) != "" {
		model.CRN = core.StringPtr(modelMap["crn"].(string))
	}
	return model, nil
}

func ResourceIBMISInstanceRescueMapToEncryptionKeyIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.EncryptionKeyIdentityByCRN, error) {
	model := &vpcv1.EncryptionKeyIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToVolumeProfileIdentity(modelMap map[string]interface{}) (vpcv1.VolumeProfileIdentityIntf, error) {
	model := &vpcv1.VolumeProfileIdentity{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMISInstanceRescueMapToVolumeProfileIdentityByName(modelMap map[string]interface{}) (*vpcv1.VolumeProfileIdentityByName, error) {
	model := &vpcv1.VolumeProfileIdentityByName{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToVolumeProfileIdentityByHref(modelMap map[string]interface{}) (*vpcv1.VolumeProfileIdentityByHref, error) {
	model := &vpcv1.VolumeProfileIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToResourceGroupIdentity(modelMap map[string]interface{}) (vpcv1.ResourceGroupIdentityIntf, error) {
	model := &vpcv1.ResourceGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func ResourceIBMISInstanceRescueMapToResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ResourceGroupIdentityByID, error) {
	model := &vpcv1.ResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMISInstanceRescueMapToInstanceRescuePrototype(modelMap map[string]interface{}) (vpcv1.InstanceRescuePrototypeIntf, error) {
	model := &vpcv1.InstanceRescuePrototype{}
	if modelMap["keys"] != nil {
		keys := []vpcv1.KeyIdentityIntf{}
		for _, keysItem := range modelMap["keys"].([]interface{}) {
			keysItemModel, err := ResourceIBMISInstanceRescueMapToKeyIdentity(keysItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			keys = append(keys, keysItemModel)
		}
		model.Keys = keys
	}
	if modelMap["user_data"] != nil && modelMap["user_data"].(string) != "" {
		model.UserData = core.StringPtr(modelMap["user_data"].(string))
	}
	if modelMap["image"] != nil && len(modelMap["image"].([]interface{})) > 0 {
		ImageModel, err := ResourceIBMISInstanceRescueMapToImageIdentity(modelMap["image"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Image = ImageModel
	}
	if modelMap["rescue_volume_attachment"] != nil && len(modelMap["rescue_volume_attachment"].([]interface{})) > 0 {
		RescueVolumeAttachmentModel, err := ResourceIBMISInstanceRescueMapToInstanceRescueVolumeAttachmentPrototype(modelMap["rescue_volume_attachment"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.RescueVolumeAttachment = RescueVolumeAttachmentModel
	}
	return model, nil
}

func ResourceIBMISInstanceRescueMapToInstanceRescuePrototypeInstanceRescueByImage(modelMap map[string]interface{}) (*vpcv1.InstanceRescuePrototypeInstanceRescueByImage, error) {
	model := &vpcv1.InstanceRescuePrototypeInstanceRescueByImage{}
	if modelMap["keys"] != nil {
		keys := []vpcv1.KeyIdentityIntf{}
		for _, keysItem := range modelMap["keys"].([]interface{}) {
			keysItemModel, err := ResourceIBMISInstanceRescueMapToKeyIdentity(keysItem.(map[string]interface{}))
			if err != nil {
				return model, err
			}
			keys = append(keys, keysItemModel)
		}
		model.Keys = keys
	}
	if modelMap["user_data"] != nil && modelMap["user_data"].(string) != "" {
		model.UserData = core.StringPtr(modelMap["user_data"].(string))
	}
	ImageModel, err := ResourceIBMISInstanceRescueMapToImageIdentity(modelMap["image"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Image = ImageModel
	RescueVolumeAttachmentModel, err := ResourceIBMISInstanceRescueMapToInstanceRescueVolumeAttachmentPrototype(modelMap["rescue_volume_attachment"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.RescueVolumeAttachment = RescueVolumeAttachmentModel
	return model, nil
}

func ResourceIBMISInstanceRescueImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMISInstanceRescueDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := ResourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMISInstanceRescueDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

// Helper function to map deleted field to avoid code duplication
func mapDeletedField(deleted *vpcv1.Deleted, modelMap map[string]interface{}) error {
	if deleted != nil {
		deletedMap, err := ResourceIBMISInstanceRescueDeletedToMap(deleted)
		if err != nil {
			return err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	return nil
}

func ResourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := ResourceIBMISInstanceRescueAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := ResourceIBMISInstanceRescueRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func ResourceIBMISInstanceRescueAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMISInstanceRescueRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMISInstanceRescueKeyReferenceToMap(model *vpcv1.KeyReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMISInstanceRescueDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["fingerprint"] = *model.Fingerprint
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(model *vpcv1.VolumeAttachmentReferenceInstanceContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if err := mapDeletedField(model.Deleted, modelMap); err != nil {
		return modelMap, err
	}
	if model.Device != nil {
		deviceMap, err := ResourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(model.Device)
		if err != nil {
			return modelMap, err
		}
		modelMap["device"] = []map[string]interface{}{deviceMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Volume != nil {
		volumeMap, err := ResourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(model.Volume)
		if err != nil {
			return modelMap, err
		}
		modelMap["volume"] = []map[string]interface{}{volumeMap}
	}
	return modelMap, nil
}

func ResourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(model *vpcv1.VolumeAttachmentDevice) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(model *vpcv1.VolumeReferenceVolumeAttachmentContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if err := mapDeletedField(model.Deleted, modelMap); err != nil {
		return modelMap, err
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMISInstanceRescueInstanceRescuePasswordToMap(model *vpcv1.InstanceRescuePassword) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["encrypted_password"] = base64.StdEncoding.EncodeToString(*model.EncryptedPassword)
	encryptionKeyMap, err := ResourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(model.EncryptionKey)
	if err != nil {
		return modelMap, err
	}
	modelMap["encryption_key"] = []map[string]interface{}{encryptionKeyMap}
	return modelMap, nil
}

func ResourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(model *vpcv1.InstanceRescueEncryptionKey) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if err := mapDeletedField(model.Deleted, modelMap); err != nil {
		return modelMap, err
	}
	modelMap["fingerprint"] = *model.Fingerprint
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

const (
	isInstanceRescuePollingDelay      = 10 * time.Second
	isInstanceRescuePollingMinTimeout = 10 * time.Second
)

// Wait for instance to enter rescue mode
func isWaitForInstanceRescueMode(client *vpcv1.VpcV1, timeout time.Duration, instanceID string, d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Waiting for instance (%s) to enter rescue mode", instanceID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"stopping", "stopped", "starting", "stable"},
		Target:  []string{"rescue"},
		Refresh: func() (interface{}, string, error) {
			getOptions := &vpcv1.GetInstanceOptions{
				ID: &instanceID,
			}
			instance, _, err := client.GetInstance(getOptions)
			if err != nil {
				return nil, "", err
			}
			if instance.LifecycleState == nil {
				return instance, "", fmt.Errorf("instance lifecycle_state is nil")
			}
			log.Printf("[DEBUG] Instance (%s) lifecycle_state: %s", instanceID, *instance.LifecycleState)
			return instance, *instance.LifecycleState, nil
		},
		Timeout:    timeout,
		Delay:      isInstanceRescuePollingDelay,
		MinTimeout: isInstanceRescuePollingMinTimeout,
	}
	return stateConf.WaitForState()
}

// Wait for instance to exit rescue mode and return to stable state
func isWaitForInstanceExitRescue(client *vpcv1.VpcV1, timeout time.Duration, instanceID string, d *schema.ResourceData) (interface{}, error) {
	log.Printf("[INFO] Waiting for instance (%s) to exit rescue mode", instanceID)

	stateConf := &resource.StateChangeConf{
		Pending: []string{"rescue", "updating"},
		Target:  []string{"stable"},
		Refresh: func() (interface{}, string, error) {
			getOptions := &vpcv1.GetInstanceOptions{
				ID: &instanceID,
			}
			instance, _, err := client.GetInstance(getOptions)
			if err != nil {
				return nil, "", err
			}
			if instance.LifecycleState == nil {
				return instance, "", fmt.Errorf("instance lifecycle_state is nil")
			}
			log.Printf("[DEBUG] Instance (%s) lifecycle_state: %s", instanceID, *instance.LifecycleState)
			return instance, *instance.LifecycleState, nil
		},
		Timeout:    timeout,
		Delay:      isInstanceRescuePollingDelay,
		MinTimeout: isInstanceRescuePollingMinTimeout,
	}
	return stateConf.WaitForState()
}
