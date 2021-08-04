// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIBMIsBackupPolicyJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPolicyJobRead,

		Schema: map[string]*schema.Schema{
			"backup_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The backup policy identifier.",
			},
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The backup policy job identifier.",
			},
			"backup_info": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The snapshot created by this backup policy job (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this snapshot.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the reference resource was deleted.",
									},
									"final_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name the referenced resource had at the time it was deleted.",
									},
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
							Description: "The URL for this snapshot.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this snapshot.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this snapshot.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"backup_resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A resource type this backup policy job applied to.",
			},
			"completed_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy job was completed.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this backup policy job.",
			},
			"plan_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier for this backup policy plan.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"source_volume": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The source volume this backup was created from (may be[deleted](https://cloud.ibm.com/apidocs/vpc#deleted-resources)).",
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
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"deleted_at": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The date and time that the reference resource was deleted.",
									},
									"final_name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name the referenced resource had at the time it was deleted.",
									},
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
							Description: "The unique user-defined name for this volume.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"started_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the backup policy job was started.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the backup policy job.The enumerated values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the volume on which the unexpected property value was encountered.",
			},
			"status_reasons": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).The enumerated reason code values for this property will expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the resource on which the unexpected reason code was encountered.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A snake case string succinctly identifying the status reason.",
						},
						"message": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "An explanation of the status reason.",
						},
						"more_info": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Link to documentation about this status reason.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsBackupPolicyJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getBackupPolicyJobOptions := &vpcv1.GetBackupPolicyJobOptions{}

	getBackupPolicyJobOptions.SetBackupPolicyID(d.Get("backup_policy_id").(string))
	getBackupPolicyJobOptions.SetID(d.Get("id").(string))

	backupPolicyJob, response, err := vpcClient.GetBackupPolicyJobWithContext(context, getBackupPolicyJobOptions)
	if err != nil {
		log.Printf("[DEBUG] GetBackupPolicyJobWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("GetBackupPolicyJobWithContext failed %s\n%s", err, response))
	}

	d.SetId(*backupPolicyJob.ID)

	if backupPolicyJob.BackupInfo != nil {
		err = d.Set("backup_info", dataSourceBackupPolicyJobFlattenBackupInfo(*backupPolicyJob.BackupInfo))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting backup_info %s", err))
		}
	}
	if err = d.Set("backup_resource_type", backupPolicyJob.BackupResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting backup_resource_type: %s", err))
	}
	if err = d.Set("completed_at", dateTimeToString(backupPolicyJob.CompletedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting completed_at: %s", err))
	}
	if err = d.Set("href", backupPolicyJob.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("plan_id", backupPolicyJob.PlanID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting plan_id: %s", err))
	}
	if err = d.Set("resource_type", backupPolicyJob.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	if backupPolicyJob.SourceVolume != nil {
		err = d.Set("source_volume", dataSourceBackupPolicyJobFlattenSourceVolume(*backupPolicyJob.SourceVolume))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting source_volume %s", err))
		}
	}
	if err = d.Set("started_at", dateTimeToString(backupPolicyJob.StartedAt)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting started_at: %s", err))
	}
	if err = d.Set("status", backupPolicyJob.Status); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting status: %s", err))
	}

	if backupPolicyJob.StatusReasons != nil {
		err = d.Set("status_reasons", dataSourceBackupPolicyJobFlattenStatusReasons(backupPolicyJob.StatusReasons))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting status_reasons %s", err))
		}
	}

	return nil
}

func dataSourceBackupPolicyJobFlattenBackupInfo(result vpcv1.SnapshotReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceBackupPolicyJobBackupInfoToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceBackupPolicyJobBackupInfoToMap(backupInfoItem vpcv1.SnapshotReference) (backupInfoMap map[string]interface{}) {
	backupInfoMap = map[string]interface{}{}

	if backupInfoItem.CRN != nil {
		backupInfoMap["crn"] = backupInfoItem.CRN
	}
	if backupInfoItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobBackupInfoDeletedToMap(*backupInfoItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		backupInfoMap["deleted"] = deletedList
	}
	if backupInfoItem.Href != nil {
		backupInfoMap["href"] = backupInfoItem.Href
	}
	if backupInfoItem.ID != nil {
		backupInfoMap["id"] = backupInfoItem.ID
	}
	if backupInfoItem.Name != nil {
		backupInfoMap["name"] = backupInfoItem.Name
	}
	if backupInfoItem.ResourceType != nil {
		backupInfoMap["resource_type"] = backupInfoItem.ResourceType
	}

	return backupInfoMap
}

func dataSourceBackupPolicyJobBackupInfoDeletedToMap(deletedItem vpcv1.SnapshotReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.DeletedAt != nil {
		deletedMap["deleted_at"] = deletedItem.DeletedAt.String()
	}
	if deletedItem.FinalName != nil {
		deletedMap["final_name"] = deletedItem.FinalName
	}
	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyJobFlattenSourceVolume(result vpcv1.VolumeReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceBackupPolicyJobSourceVolumeToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceBackupPolicyJobSourceVolumeToMap(sourceVolumeItem vpcv1.VolumeReference) (sourceVolumeMap map[string]interface{}) {
	sourceVolumeMap = map[string]interface{}{}

	if sourceVolumeItem.CRN != nil {
		sourceVolumeMap["crn"] = sourceVolumeItem.CRN
	}
	if sourceVolumeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobSourceVolumeDeletedToMap(*sourceVolumeItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		sourceVolumeMap["deleted"] = deletedList
	}
	if sourceVolumeItem.Href != nil {
		sourceVolumeMap["href"] = sourceVolumeItem.Href
	}
	if sourceVolumeItem.ID != nil {
		sourceVolumeMap["id"] = sourceVolumeItem.ID
	}
	if sourceVolumeItem.Name != nil {
		sourceVolumeMap["name"] = sourceVolumeItem.Name
	}
	if sourceVolumeItem.ResourceType != nil {
		sourceVolumeMap["resource_type"] = sourceVolumeItem.ResourceType
	}

	return sourceVolumeMap
}

func dataSourceBackupPolicyJobSourceVolumeDeletedToMap(deletedItem vpcv1.VolumeReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.DeletedAt != nil {
		deletedMap["deleted_at"] = deletedItem.DeletedAt.String()
	}
	if deletedItem.FinalName != nil {
		deletedMap["final_name"] = deletedItem.FinalName
	}
	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyJobFlattenStatusReasons(result []vpcv1.BackupPolicyJobStatusReason) (statusReasons []map[string]interface{}) {
	for _, statusReasonsItem := range result {
		statusReasons = append(statusReasons, dataSourceBackupPolicyJobStatusReasonsToMap(statusReasonsItem))
	}

	return statusReasons
}

func dataSourceBackupPolicyJobStatusReasonsToMap(statusReasonsItem vpcv1.BackupPolicyJobStatusReason) (statusReasonsMap map[string]interface{}) {
	statusReasonsMap = map[string]interface{}{}

	if statusReasonsItem.Code != nil {
		statusReasonsMap["code"] = statusReasonsItem.Code
	}
	if statusReasonsItem.Message != nil {
		statusReasonsMap["message"] = statusReasonsItem.Message
	}
	if statusReasonsItem.MoreInfo != nil {
		statusReasonsMap["more_info"] = statusReasonsItem.MoreInfo
	}

	return statusReasonsMap
}
