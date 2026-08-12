// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsBareMetalServerVolumeAttachment() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerVolumeAttachmentRead,

		Schema: map[string]*schema.Schema{
			isBMSVolAttBareMetalServer: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttBareMetalServer),
				Description:  "The bare metal server identifier",
			},
			isBMSVolAttName: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_bare_metal_server_volume_attachment", isBMSVolAttName),
				Description:  "The name for this volume attachment. The name is unique across all volume attachments on the instance or bare metal server.",
			},
			isBMSVolAttBandwidth: {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum bandwidth (in megabits per second) for the volume when attached to this bare metal server.",
			},
			isBMSVolAttCreatedAt: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the volume attachment was created.",
			},
			isBMSVolAttDeleteOnServerDelete: {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Indicates whether deleting the bare metal server will also delete the attached volume. This property must be false if the volume's attachment_mode is multiple.",
			},
			isBMSVolAttDevice: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A unique identifier for the device which is exposed to the bare metal server operating system. This property may be absent if the status of the volume attachment is not available.",
			},
			isBMSVolAttHref: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this bare metal server volume attachment.",
			},
			isBMSVolAttId: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this bare metal volume attachment.",
			},
			isBMSVolAttStatus: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this volume attachment: attaching, detaching, available, unusable",
			},
			isBMSVolAttStatusReason: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status reason code: volume_encryption_key_deleted (The key associated with the data volume attached to the bare metal server is deleted).",
						},
						"message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A link to documentation about this status reason.",
						},
					},
				},
			},
			isBMSVolAttType: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of volume attachment.",
			},
			isBMSVolAttProtocol: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The protocol used for this volume attachment: nvme_tcp (Non-Volatile Memory Express (NVMe) over TCP/IP, which allows bare metal servers to connect to volumes over the network using the NVMe protocol).",
			},
			isBMSVolAttNvmeQualifiedName: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The NVMe Qualified Name (NQN) of the subsystem. This unique identifier is used by the bare metal server to establish a connection to the volume over the nvme-tcp protocol. The NQN must be used when configuring the NVMe initiator on the bare metal server to access the attached volume.",
			},
			isBMSVolAttIps: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The IP addresses for connecting to the volume using nvme_tcp.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			isBMSVolAttVol: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The attached volume.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBMSVolAttVol: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this volume.",
						},
						isBMSVolAttVolName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the attached volume.",
						},
						isBMSVolAttVolCRN: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN of the attached volume.",
						},
						isBMSVolAttVolHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL of the attached volume.",
						},
					},
				},
			},
		},
	}
}

func DataSourceIBMIsBareMetalServerVolumeAttachmentValidator() *validate.ResourceValidator {
	validateSchema := []validate.ValidateSchema{
		{
			Identifier:                 isBMSVolAttBareMetalServer,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
			Required:                   true,
		},
		{
			Identifier:                 isBMSVolAttName,
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	}

	validator := validate.ResourceValidator{
		ResourceName: "ibm_is_bare_metal_server_volume_attachment",
		Schema:       validateSchema,
	}
	return &validator
}

func dataSourceIBMISBareMetalServerVolumeAttachmentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBMSVolAttBareMetalServer).(string)
	name := d.Get(isBMSVolAttName).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_bare_metal_server_volume_attachment", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &vpcv1.ListBareMetalServerVolumeAttachmentsOptions{
		BareMetalServerID: &bareMetalServerID,
	}
	volumeAttachments, _, err := sess.ListBareMetalServerVolumeAttachmentsWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListBareMetalServerVolumeAttachmentsWithContext failed: %s", err.Error()), "(Data) ibm_is_bare_metal_server_volume_attachment", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	for _, volumeAttachmentIntf := range volumeAttachments.VolumeAttachments {
		volumeAttachment, ok := volumeAttachmentIntf.(*vpcv1.BareMetalServerVolumeAttachment)
		if !ok || volumeAttachment == nil || volumeAttachment.Name == nil || *volumeAttachment.Name != name {
			continue
		}

		d.SetId(makeTerraformVolAttID(bareMetalServerID, *volumeAttachment.ID))
		d.Set(isBMSVolAttBareMetalServer, bareMetalServerID)
		d.Set(isBMSVolAttName, *volumeAttachment.Name)

		if volumeAttachment.Bandwidth != nil {
			d.Set(isBMSVolAttBandwidth, int(*volumeAttachment.Bandwidth))
		}
		if volumeAttachment.CreatedAt != nil {
			d.Set(isBMSVolAttCreatedAt, volumeAttachment.CreatedAt.String())
		}
		d.Set(isBMSVolAttDeleteOnServerDelete, *volumeAttachment.DeleteVolumeOnBareMetalServerDelete)

		if volumeAttachment.Device != nil && volumeAttachment.Device.ID != nil {
			d.Set(isBMSVolAttDevice, *volumeAttachment.Device.ID)
		}
		d.Set(isBMSVolAttHref, *volumeAttachment.Href)
		d.Set(isBMSVolAttId, *volumeAttachment.ID)
		d.Set(isBMSVolAttStatus, *volumeAttachment.Status)

		if volumeAttachment.StatusReason != nil {
			statusReasonList := make([]map[string]interface{}, 0)
			statusReasonMap := map[string]interface{}{}
			if volumeAttachment.StatusReason.Code != nil {
				statusReasonMap["code"] = *volumeAttachment.StatusReason.Code
			}
			if volumeAttachment.StatusReason.Message != nil {
				statusReasonMap["message"] = *volumeAttachment.StatusReason.Message
			}
			if volumeAttachment.StatusReason.MoreInfo != nil {
				statusReasonMap["more_info"] = *volumeAttachment.StatusReason.MoreInfo
			}
			if len(statusReasonMap) > 0 {
				statusReasonList = append(statusReasonList, statusReasonMap)
				d.Set(isBMSVolAttStatusReason, statusReasonList)
			}
		}

		d.Set(isBMSVolAttType, *volumeAttachment.Type)

		if volumeAttachment.Protocol != nil {
			d.Set(isBMSVolAttProtocol, *volumeAttachment.Protocol)
		}
		if volumeAttachment.NvmeQualifiedName != nil {
			d.Set(isBMSVolAttNvmeQualifiedName, *volumeAttachment.NvmeQualifiedName)
		}
		if len(volumeAttachment.Ips) > 0 {
			ips := make([]string, len(volumeAttachment.Ips))
			for i, ip := range volumeAttachment.Ips {
				ips[i] = *ip.Address
			}
			d.Set(isBMSVolAttIps, ips)
		}

		if volumeAttachment.Volume != nil {
			volList := make([]map[string]interface{}, 0)
			currentVol := map[string]interface{}{}
			currentVol[isBMSVolAttVol] = *volumeAttachment.Volume.ID
			currentVol[isBMSVolAttVolName] = *volumeAttachment.Volume.Name
			currentVol[isBMSVolAttVolCRN] = *volumeAttachment.Volume.CRN
			currentVol[isBMSVolAttVolHref] = *volumeAttachment.Volume.Href
			volList = append(volList, currentVol)
			if err = d.Set(isBMSVolAttVol, volList); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting volume: %s", err), "(Data) ibm_is_bare_metal_server_volume_attachment", "read", "set-volume").GetDiag()
			}
		}

		return nil
	}

	err = fmt.Errorf("No bare metal server volume attachment found with name %s on bare metal server %s", name, bareMetalServerID)
	tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListBareMetalServerVolumeAttachmentsWithContext failed: %s", err.Error()), "(Data) ibm_is_bare_metal_server_volume_attachment", "read")
	log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
	return tfErr.GetDiag()
}
