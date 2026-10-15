// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.117.0-7f07c563-20260915-094553
*/

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsImageJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImageJobsRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to image jobs with a `name` property matching the specified value.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to image jobs with a `status` property matching the specified value.",
			},
			"jobs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of image jobs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"completed_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the image job was completed.If absent, the image job has not yet completed.",
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the image job was created.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this image job.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this image job.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this image job. The name is unique across all image jobs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"started_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the image job started running.If absent, the image job has not yet started.",
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of this image job:- `deleting`: The image job is being deleted- `failed`: The image job could not be completed successfully. The customer may  retry by creating a new image job for the same image.- `queued`: The image job has been accepted and is waiting to be scheduled- `running`: The image job is in progress- `succeeded`: The image job was completed successfullyThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"status_reasons": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The reasons for the current status (if any).",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"code": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A reason code for the image job status:The following reason codes may be set when the image job `status` is `failed`:- `failed_component_installation`: one or more image components failed to install while  building the image- `failed_component_validation`: one or more image components failed validation- `failed_image_capture`: failed to capture the image from the source volume- `failed_timeout`: the image job exceeded the configured `timeout`- `failed_worker_provisioning`: the worker used to run the image job failed to  provision- `internal_error`: internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
									},
									"message": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An explanation of the status reason.",
									},
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A link to documentation about this status reason.",
									},
								},
							},
						},
						"type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of image job:- `build`: build an image from an image recipeThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"delete_image_job_on_image_delete": &schema.Schema{
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Indicates whether deleting the image will also delete the image job that created it.",
						},
						"image": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
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
						"image_recipe": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The image recipe and version used by this image job.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The CRN for this image recipe.",
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
										Description: "The URL for this image recipe.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this image recipe.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique name for this image recipe.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "A generated content fingerprint for this version of the image recipe. The value changes whenever the image recipe is updated.",
									},
								},
							},
						},
						"image_template": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The template for creating an image from the image job.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name for this image. The name must not be used by another image in the region. Names starting with `ibm-` are reserved for system-provided images, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
									},
									"resource_group": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) will be used.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this resource group.",
												},
											},
										},
									},
								},
							},
						},
						"timeout": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The maximum time (in minutes) allowed for this image job to complete.The time begins when the job moves to a `status` of `running`. If the timeout is reached, the image and worker will be deleted.",
						},
						"worker": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The worker that executes this job.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The instance resource that is the worker for this job. If absent, the instance has notbeen created yet.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this virtual server instance.",
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
													Description: "The URL for this virtual server instance.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this virtual server instance.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The name for this virtual server instance. The name is unique across all virtual server instances in the region.",
												},
											},
										},
									},
									"profile": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Description: "The [profile](https://cloud.ibm.com/docs/vpc?topic=vpc-profiles) used to build the image.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"href": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The URL for this virtual server instance profile.",
												},
												"name": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The globally unique name for this virtual server instance profile.",
												},
												"resource_type": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The resource type.",
												},
											},
										},
									},
									"subnet": &schema.Schema{
										Type:        schema.TypeList,
										Computed:    true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"crn": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The CRN for this subnet.",
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
													Description: "The URL for this subnet.",
												},
												"id": &schema.Schema{
													Type:        schema.TypeString,
													Computed:    true,
													Description: "The unique identifier for this subnet.",
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
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsImageJobsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_jobs", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listImageJobsOptions := &vpcv1.ListImageJobsOptions{}

	if _, ok := d.GetOk("name"); ok {
		listImageJobsOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("status"); ok {
		listImageJobsOptions.SetStatus(d.Get("status").(string))
	}

	var pager *vpcv1.ImageJobsPager
	pager, err = vpcClient.NewImageJobsPager(listImageJobsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_image_jobs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ImageJobsPager.GetAll() failed %s", err), "(Data) ibm_is_image_jobs", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsImageJobsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsImageJobsImageJobToMap(modelItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_jobs", "read", "Images-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("jobs", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting jobs %s", err), "(Data) ibm_is_image_jobs", "read", "jobs-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsImageJobsID returns a reasonable ID for the list.
func dataSourceIBMIsImageJobsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsImageJobsImageJobToMap(model vpcv1.ImageJobIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ImageJobBuild); ok {
		return DataSourceIBMIsImageJobsImageJobBuildToMap(model.(*vpcv1.ImageJobBuild))
	} else if _, ok := model.(*vpcv1.ImageJob); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ImageJob)
		if model.CompletedAt != nil {
			modelMap["completed_at"] = model.CompletedAt.String()
		}
		modelMap["created_at"] = model.CreatedAt.String()
		modelMap["href"] = *model.Href
		modelMap["id"] = *model.ID
		modelMap["name"] = *model.Name
		modelMap["resource_type"] = *model.ResourceType
		if model.StartedAt != nil {
			modelMap["started_at"] = model.StartedAt.String()
		}
		modelMap["status"] = *model.Status
		statusReasons := []map[string]interface{}{}
		for _, statusReasonsItem := range model.StatusReasons {
			statusReasonsItemMap, err := DataSourceIBMIsImageJobsImageJobStatusReasonsToMap(&statusReasonsItem) // #nosec G601
			if err != nil {
				return modelMap, err
			}
			statusReasons = append(statusReasons, statusReasonsItemMap)
		}
		modelMap["status_reasons"] = statusReasons
		modelMap["type"] = *model.Type
		if model.DeleteImageJobOnImageDelete != nil {
			modelMap["delete_image_job_on_image_delete"] = *model.DeleteImageJobOnImageDelete
		}
		if model.Image != nil {
			imageMap, err := DataSourceIBMIsImageJobsImageReferenceToMap(model.Image)
			if err != nil {
				return modelMap, err
			}
			modelMap["image"] = []map[string]interface{}{imageMap}
		}
		if model.ImageRecipe != nil {
			imageRecipeMap, err := DataSourceIBMIsImageJobsImageJobBuildImageRecipeToMap(model.ImageRecipe)
			if err != nil {
				return modelMap, err
			}
			modelMap["image_recipe"] = []map[string]interface{}{imageRecipeMap}
		}
		if model.ImageTemplate != nil {
			imageTemplateMap, err := DataSourceIBMIsImageJobsImageJobBuildImageTemplateToMap(model.ImageTemplate)
			if err != nil {
				return modelMap, err
			}
			modelMap["image_template"] = []map[string]interface{}{imageTemplateMap}
		}
		if model.Timeout != nil {
			modelMap["timeout"] = flex.IntValue(model.Timeout)
		}
		if model.Worker != nil {
			workerMap, err := DataSourceIBMIsImageJobsImageJobBuildWorkerToMap(model.Worker)
			if err != nil {
				return modelMap, err
			}
			modelMap["worker"] = []map[string]interface{}{workerMap}
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("unrecognized vpcv1.ImageJobIntf subtype encountered")
	}
}

func DataSourceIBMIsImageJobsImageJobStatusReasonsToMap(model *vpcv1.ImageJobStatusReasons) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsImageJobsImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobsDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsImageJobsImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageJobsDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsImageJobsImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := DataSourceIBMIsImageJobsAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsImageJobsRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsImageJobsAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageJobsRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageJobsImageJobBuildImageRecipeToMap(model *vpcv1.ImageJobBuildImageRecipe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobsDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	modelMap["version"] = *model.Version
	return modelMap, nil
}

func DataSourceIBMIsImageJobsImageJobBuildImageTemplateToMap(model *vpcv1.ImageJobBuildImageTemplate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceGroup != nil {
		resourceGroupMap, err := DataSourceIBMIsImageJobsResourceGroupIdentityToMap(model.ResourceGroup)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsImageJobsResourceGroupIdentityToMap(model vpcv1.ResourceGroupIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ResourceGroupIdentityResourceGroupIdentityByID); ok {
		return DataSourceIBMIsImageJobsResourceGroupIdentityResourceGroupIdentityByIDToMap(model.(*vpcv1.ResourceGroupIdentityResourceGroupIdentityByID))
	} else if _, ok := model.(*vpcv1.ResourceGroupIdentity); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ResourceGroupIdentity)
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("unrecognized vpcv1.ResourceGroupIdentityIntf subtype encountered")
	}
}

func DataSourceIBMIsImageJobsResourceGroupIdentityResourceGroupIdentityByIDToMap(model *vpcv1.ResourceGroupIdentityResourceGroupIdentityByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMIsImageJobsImageJobBuildWorkerToMap(model *vpcv1.ImageJobBuildWorker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instance != nil {
		instanceMap, err := DataSourceIBMIsImageJobsInstanceReferenceToMap(model.Instance)
		if err != nil {
			return modelMap, err
		}
		modelMap["instance"] = []map[string]interface{}{instanceMap}
	}
	profileMap, err := DataSourceIBMIsImageJobsImageJobBuildWorkerProfileToMap(model.Profile)
	if err != nil {
		return modelMap, err
	}
	modelMap["profile"] = []map[string]interface{}{profileMap}
	subnetMap, err := DataSourceIBMIsImageJobsSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	return modelMap, nil
}

func DataSourceIBMIsImageJobsInstanceReferenceToMap(model *vpcv1.InstanceReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobsDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageJobsImageJobBuildWorkerProfileToMap(model *vpcv1.ImageJobBuildWorkerProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageJobsSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobsDeletedToMap(model.Deleted)
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

func DataSourceIBMIsImageJobsImageJobBuildToMap(model *vpcv1.ImageJobBuild) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.CompletedAt != nil {
		modelMap["completed_at"] = model.CompletedAt.String()
	}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	if model.StartedAt != nil {
		modelMap["started_at"] = model.StartedAt.String()
	}
	modelMap["status"] = *model.Status
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range model.StatusReasons {
		statusReasonsItemMap, err := DataSourceIBMIsImageJobsImageJobStatusReasonsToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	modelMap["status_reasons"] = statusReasons
	modelMap["delete_image_job_on_image_delete"] = *model.DeleteImageJobOnImageDelete
	if model.Image != nil {
		imageMap, err := DataSourceIBMIsImageJobsImageReferenceToMap(model.Image)
		if err != nil {
			return modelMap, err
		}
		modelMap["image"] = []map[string]interface{}{imageMap}
	}
	imageRecipeMap, err := DataSourceIBMIsImageJobsImageJobBuildImageRecipeToMap(model.ImageRecipe)
	if err != nil {
		return modelMap, err
	}
	modelMap["image_recipe"] = []map[string]interface{}{imageRecipeMap}
	if model.ImageTemplate != nil {
		imageTemplateMap, err := DataSourceIBMIsImageJobsImageJobBuildImageTemplateToMap(model.ImageTemplate)
		if err != nil {
			return modelMap, err
		}
		modelMap["image_template"] = []map[string]interface{}{imageTemplateMap}
	}
	modelMap["timeout"] = flex.IntValue(model.Timeout)
	modelMap["type"] = *model.Type
	workerMap, err := DataSourceIBMIsImageJobsImageJobBuildWorkerToMap(model.Worker)
	if err != nil {
		return modelMap, err
	}
	modelMap["worker"] = []map[string]interface{}{workerMap}
	return modelMap, nil
}
