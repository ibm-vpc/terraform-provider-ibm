// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIBMIsBackupPolicyJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPolicyJobsRead,

		Schema: map[string]*schema.Schema{
			"backup_policy_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The backup policy identifier.",
			},
			"jobs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of backup policy jobs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this backup policy job.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsBackupPolicyJobsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listBackupPolicyJobsOptions := &vpcv1.ListBackupPolicyJobsOptions{}

	listBackupPolicyJobsOptions.SetBackupPolicyID(d.Get("backup_policy_id").(string))

	backupPolicyJobCollection, response, err := vpcClient.ListBackupPolicyJobsWithContext(context, listBackupPolicyJobsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListBackupPolicyJobsWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListBackupPolicyJobsWithContext failed %s\n%s", err, response))
	}

	d.SetId(dataSourceIBMIsBackupPolicyJobsID(d))

	if backupPolicyJobCollection.Jobs != nil {
		err = d.Set("jobs", dataSourceBackupPolicyJobCollectionFlattenJobs(backupPolicyJobCollection.Jobs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting jobs %s", err))
		}
	}

	return nil
}

// dataSourceIBMIsBackupPolicyJobsID returns a reasonable ID for the list.
func dataSourceIBMIsBackupPolicyJobsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceBackupPolicyJobCollectionFlattenJobs(result []vpcv1.BackupPolicyJob) (jobs []map[string]interface{}) {
	for _, jobsItem := range result {
		jobs = append(jobs, dataSourceBackupPolicyJobCollectionJobsToMap(jobsItem))
	}

	return jobs
}

func dataSourceBackupPolicyJobCollectionJobsToMap(jobsItem vpcv1.BackupPolicyJob) (jobsMap map[string]interface{}) {
	jobsMap = map[string]interface{}{}

	if jobsItem.BackupInfo != nil {
		backupInfoList := []map[string]interface{}{}
		backupInfoMap := dataSourceBackupPolicyJobCollectionJobsBackupInfoToMap(*jobsItem.BackupInfo)
		backupInfoList = append(backupInfoList, backupInfoMap)
		jobsMap["backup_info"] = backupInfoList
	}
	if jobsItem.BackupResourceType != nil {
		jobsMap["backup_resource_type"] = jobsItem.BackupResourceType
	}
	if jobsItem.CompletedAt != nil {
		jobsMap["completed_at"] = jobsItem.CompletedAt.String()
	}
	if jobsItem.Href != nil {
		jobsMap["href"] = jobsItem.Href
	}
	if jobsItem.ID != nil {
		jobsMap["id"] = jobsItem.ID
	}
	if jobsItem.PlanID != nil {
		jobsMap["plan_id"] = jobsItem.PlanID
	}
	if jobsItem.ResourceType != nil {
		jobsMap["resource_type"] = jobsItem.ResourceType
	}
	if jobsItem.SourceVolume != nil {
		sourceVolumeList := []map[string]interface{}{}
		sourceVolumeMap := dataSourceBackupPolicyJobCollectionJobsSourceVolumeToMap(*jobsItem.SourceVolume)
		sourceVolumeList = append(sourceVolumeList, sourceVolumeMap)
		jobsMap["source_volume"] = sourceVolumeList
	}
	if jobsItem.StartedAt != nil {
		jobsMap["started_at"] = jobsItem.StartedAt.String()
	}
	if jobsItem.Status != nil {
		jobsMap["status"] = jobsItem.Status
	}
	if jobsItem.StatusReasons != nil {
		statusReasonsList := []map[string]interface{}{}
		for _, statusReasonsItem := range jobsItem.StatusReasons {
			statusReasonsList = append(statusReasonsList, dataSourceBackupPolicyJobCollectionJobsStatusReasonsToMap(statusReasonsItem))
		}
		jobsMap["status_reasons"] = statusReasonsList
	}

	return jobsMap
}

func dataSourceBackupPolicyJobCollectionJobsBackupInfoToMap(backupInfoItem vpcv1.SnapshotReference) (backupInfoMap map[string]interface{}) {
	backupInfoMap = map[string]interface{}{}

	if backupInfoItem.CRN != nil {
		backupInfoMap["crn"] = backupInfoItem.CRN
	}
	if backupInfoItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobCollectionBackupInfoDeletedToMap(*backupInfoItem.Deleted)
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

func dataSourceBackupPolicyJobCollectionBackupInfoDeletedToMap(deletedItem vpcv1.SnapshotReferenceDeleted) (deletedMap map[string]interface{}) {
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

func dataSourceBackupPolicyJobCollectionJobsSourceVolumeToMap(sourceVolumeItem vpcv1.VolumeReference) (sourceVolumeMap map[string]interface{}) {
	sourceVolumeMap = map[string]interface{}{}

	if sourceVolumeItem.CRN != nil {
		sourceVolumeMap["crn"] = sourceVolumeItem.CRN
	}
	if sourceVolumeItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyJobCollectionSourceVolumeDeletedToMap(*sourceVolumeItem.Deleted)
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

func dataSourceBackupPolicyJobCollectionSourceVolumeDeletedToMap(deletedItem vpcv1.VolumeReferenceDeleted) (deletedMap map[string]interface{}) {
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

func dataSourceBackupPolicyJobCollectionJobsStatusReasonsToMap(statusReasonsItem vpcv1.BackupPolicyJobStatusReason) (statusReasonsMap map[string]interface{}) {
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
