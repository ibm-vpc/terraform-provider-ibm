// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMIsBareMetalServerVolumeAttachments() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerVolumeAttachmentsRead,

		Schema: map[string]*schema.Schema{
			isBMSVolAttBareMetalServer: {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeDataSourceValidator("ibm_is_bare_metal_server_volume_attachments", isBMSVolAttBareMetalServer),
				Description:  "The bare metal server identifier.",
			},
			"volume_attachments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of volume attachments for the bare metal server.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						isBMSVolAttName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this volume attachment. The name is unique across all volume attachments on the instance or bare metal server.",
						},
						isBMSVolAttStatus: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of this volume attachment: attaching (volume attachment is being initialized and not yet usable), detaching (volume attachment is being removed), available (volume attachment is usable, connection to the volume can be established from the server's operating system), unusable (volume attachment is unusable due to the underlying volume state).",
						},
						isBMSVolAttType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of volume attachment.",
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
				},
			},
		},
	}
}

func DataSourceIBMIsBareMetalServerVolumeAttachmentsValidator() *validate.ResourceValidator {
	validateSchema := []validate.ValidateSchema{
		{
			Identifier:                 isBMSVolAttBareMetalServer,
			ValidateFunctionIdentifier: validate.ValidateNoZeroValues,
			Type:                       validate.TypeString,
			Required:                   true,
		},
	}
	validator := validate.ResourceValidator{
		ResourceName: "ibm_is_bare_metal_server_volume_attachments",
		Schema:       validateSchema,
	}
	return &validator
}

func dataSourceIBMISBareMetalServerVolumeAttachmentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBMSVolAttBareMetalServer).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_bare_metal_server_volume_attachments", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	options := &vpcv1.ListBareMetalServerVolumeAttachmentsOptions{
		BareMetalServerID: &bareMetalServerID,
	}
	collection, _, err := sess.ListBareMetalServerVolumeAttachmentsWithContext(context, options)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListBareMetalServerVolumeAttachmentsWithContext failed: %s", err.Error()), "(Data) ibm_is_bare_metal_server_volume_attachments", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	volAttsList := make([]map[string]interface{}, 0)
	for _, volAttIntf := range collection.VolumeAttachments {
		volAtt, ok := volAttIntf.(*vpcv1.BareMetalServerVolumeAttachment)
		if !ok || volAtt == nil {
			continue
		}
		entry := map[string]interface{}{
			isBMSVolAttDeleteOnServerDelete: *volAtt.DeleteVolumeOnBareMetalServerDelete,
			isBMSVolAttHref:                 *volAtt.Href,
			isBMSVolAttId:                   *volAtt.ID,
			isBMSVolAttName:                 *volAtt.Name,
			isBMSVolAttStatus:               *volAtt.Status,
		}
		if volAtt.Bandwidth != nil {
			entry[isBMSVolAttBandwidth] = int(*volAtt.Bandwidth)
		}
		if volAtt.CreatedAt != nil {
			entry[isBMSVolAttCreatedAt] = volAtt.CreatedAt.String()
		}
		if volAtt.Device != nil && volAtt.Device.ID != nil {
			entry[isBMSVolAttDevice] = *volAtt.Device.ID
		}

		if volAtt.StatusReason != nil {
			statusReasonList := make([]map[string]interface{}, 0)
			statusReasonMap := map[string]interface{}{}
			if volAtt.StatusReason.Code != nil {
				statusReasonMap["code"] = *volAtt.StatusReason.Code
			}
			if volAtt.StatusReason.Message != nil {
				statusReasonMap["message"] = *volAtt.StatusReason.Message
			}
			if volAtt.StatusReason.MoreInfo != nil {
				statusReasonMap["more_info"] = *volAtt.StatusReason.MoreInfo
			}
			if len(statusReasonMap) > 0 {
				statusReasonList = append(statusReasonList, statusReasonMap)
				entry["status_reason"] = statusReasonList
			}
		}

		entry[isBMSVolAttType] = *volAtt.Type

		if volAtt.Protocol != nil {
			entry[isBMSVolAttProtocol] = *volAtt.Protocol
		}
		if volAtt.NvmeQualifiedName != nil {
			entry[isBMSVolAttNvmeQualifiedName] = *volAtt.NvmeQualifiedName
		}
		if len(volAtt.Ips) > 0 {
			ips := make([]string, len(volAtt.Ips))
			for i, ip := range volAtt.Ips {
				ips[i] = *ip.Address
			}
			entry[isBMSVolAttIps] = ips
		}
		if volAtt.Volume != nil {
			volList := make([]map[string]interface{}, 0)
			currentVol := map[string]interface{}{}
			currentVol[isBMSVolAttVol] = *volAtt.Volume.ID
			currentVol[isBMSVolAttVolName] = *volAtt.Volume.Name
			currentVol[isBMSVolAttVolCRN] = *volAtt.Volume.CRN
			currentVol[isBMSVolAttVolHref] = *volAtt.Volume.Href
			volList = append(volList, currentVol)
			entry["volume"] = volList
		}
		volAttsList = append(volAttsList, entry)
	}

	d.SetId(time.Now().UTC().String())
	d.Set(isBMSVolAttBareMetalServer, bareMetalServerID)
	if err = d.Set("volume_attachments", volAttsList); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting volume_attachments: %s", err), "(Data) ibm_is_bare_metal_server_volume_attachments", "read", "set-volume_attachments").GetDiag()
	}
	return nil
}
