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

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISInstanceRescue() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISInstanceRescueRead,

		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The virtual server instance identifier.",
			},
			"image": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The image to use for rescuing the instance.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this image.",
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
										Description: "A link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this image.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
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
				Computed:    true,
				Description: "The public SSH keys used at initialization for the rescue volume.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this key.",
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
										Computed:    true,
										Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
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
			"rescue_volume_attachment": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The rescue volume attachment for this instance.",
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
										Description: "A link to documentation about deleted resources.",
									},
								},
							},
						},
						"device": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The configuration for the volume as a device in the instance operating system.This property may be absent if the volume attachment's `status` is not `attached`.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A unique identifier for the device which is exposed to the instance operating system.",
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
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this volume attachment. The name is unique across all volume attachments on the instance.",
						},
						"volume": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The attached volume.This property will be absent if the volume has not yet been provisioned.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this volume.",
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
													Description: "A link to documentation about deleted resources.",
												},
											},
										},
									},
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this volume.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this volume.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this volume. The name is unique across all volumes in the region.",
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
		},
	}
}

func dataSourceIBMISInstanceRescueRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_rescue", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getInstanceRescueOptions := vpcClient.NewGetInstanceRescueOptions(d.Get("instance_id").(string))

	getInstanceRescueOptions.SetInstanceID(d.Get("instance_id").(string))

	instanceRescue, _, err := vpcClient.GetInstanceRescueWithContext(context, getInstanceRescueOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetInstanceRescueWithContext failed: %s", err.Error()), "(Data) ibm_is_instance_rescue", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getInstanceRescueOptions.InstanceID)

	image := []map[string]interface{}{}
	imageMap, err := DataSourceIBMISInstanceRescueImageReferenceToMap(instanceRescue.Image)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_rescue", "read", "image-to-map").GetDiag()
	}
	image = append(image, imageMap)
	if err = d.Set("image", image); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting image: %s", err), "(Data) ibm_is_instance_rescue", "read", "set-image").GetDiag()
	}

	keys := []map[string]interface{}{}
	for _, keysItem := range instanceRescue.Keys {
		keysItemMap, err := DataSourceIBMISInstanceRescueKeyReferenceToMap(&keysItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_rescue", "read", "keys-to-map").GetDiag()
		}
		keys = append(keys, keysItemMap)
	}
	if err = d.Set("keys", keys); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting keys: %s", err), "(Data) ibm_is_instance_rescue", "read", "set-keys").GetDiag()
	}

	if !core.IsNil(instanceRescue.Password) {
		password := []map[string]interface{}{}
		passwordMap, err := DataSourceIBMISInstanceRescueInstanceRescuePasswordToMap(instanceRescue.Password)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_rescue", "read", "password-to-map").GetDiag()
		}
		password = append(password, passwordMap)
		if err = d.Set("password", password); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting password: %s", err), "(Data) ibm_is_instance_rescue", "read", "set-password").GetDiag()
		}
	}

	rescueVolumeAttachment := []map[string]interface{}{}
	rescueVolumeAttachmentMap, err := DataSourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(instanceRescue.RescueVolumeAttachment)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_instance_rescue", "read", "rescue_volume_attachment-to-map").GetDiag()
	}
	rescueVolumeAttachment = append(rescueVolumeAttachment, rescueVolumeAttachmentMap)
	if err = d.Set("rescue_volume_attachment", rescueVolumeAttachment); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting rescue_volume_attachment: %s", err), "(Data) ibm_is_instance_rescue", "read", "set-rescue_volume_attachment").GetDiag()
	}

	return nil
}

func DataSourceIBMISInstanceRescueImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMISInstanceRescueDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := DataSourceIBMISInstanceRescueAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := DataSourceIBMISInstanceRescueRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueKeyReferenceToMap(model *vpcv1.KeyReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMISInstanceRescueDeletedToMap(model.Deleted)
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

func DataSourceIBMISInstanceRescueInstanceRescuePasswordToMap(model *vpcv1.InstanceRescuePassword) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["encrypted_password"] = base64.StdEncoding.EncodeToString(*model.EncryptedPassword)
	encryptionKeyMap, err := DataSourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(model.EncryptionKey)
	if err != nil {
		return modelMap, err
	}
	modelMap["encryption_key"] = []map[string]interface{}{encryptionKeyMap}
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueInstanceRescueEncryptionKeyToMap(model *vpcv1.InstanceRescueEncryptionKey) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMISInstanceRescueDeletedToMap(model.Deleted)
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

func DataSourceIBMISInstanceRescueVolumeAttachmentReferenceInstanceContextToMap(model *vpcv1.VolumeAttachmentReferenceInstanceContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMISInstanceRescueDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	if model.Device != nil {
		deviceMap, err := DataSourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(model.Device)
		if err != nil {
			return modelMap, err
		}
		modelMap["device"] = []map[string]interface{}{deviceMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Volume != nil {
		volumeMap, err := DataSourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(model.Volume)
		if err != nil {
			return modelMap, err
		}
		modelMap["volume"] = []map[string]interface{}{volumeMap}
	}
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueVolumeAttachmentDeviceToMap(model *vpcv1.VolumeAttachmentDevice) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMISInstanceRescueVolumeReferenceVolumeAttachmentContextToMap(model *vpcv1.VolumeReferenceVolumeAttachmentContext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMISInstanceRescueDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}
