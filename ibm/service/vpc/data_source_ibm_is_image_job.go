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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsImageJob() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImageJobRead,

		Schema: map[string]*schema.Schema{
			"is_image_job_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image job identifier.",
			},
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
	}
}

func dataSourceIBMIsImageJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_job", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getImageJobOptions := &vpcv1.GetImageJobOptions{}

	getImageJobOptions.SetID(d.Get("is_image_job_id").(string))

	imageJobIntf, _, err := vpcClient.GetImageJobWithContext(context, getImageJobOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageJobWithContext failed: %s", err.Error()), "(Data) ibm_is_image_job", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	imageJob := imageJobIntf.(*vpcv1.ImageJob)

	d.SetId(*getImageJobOptions.ID)

	if !core.IsNil(imageJob.CompletedAt) {
		if err = d.Set("completed_at", flex.DateTimeToString(imageJob.CompletedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting completed_at: %s", err), "(Data) ibm_is_image_job", "read", "set-completed_at").GetDiag()
		}
	}

	if err = d.Set("created_at", flex.DateTimeToString(imageJob.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting created_at: %s", err), "(Data) ibm_is_image_job", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("href", imageJob.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting href: %s", err), "(Data) ibm_is_image_job", "read", "set-href").GetDiag()
	}

	if err = d.Set("name", imageJob.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting name: %s", err), "(Data) ibm_is_image_job", "read", "set-name").GetDiag()
	}

	if err = d.Set("resource_type", imageJob.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting resource_type: %s", err), "(Data) ibm_is_image_job", "read", "set-resource_type").GetDiag()
	}

	if !core.IsNil(imageJob.StartedAt) {
		if err = d.Set("started_at", flex.DateTimeToString(imageJob.StartedAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting started_at: %s", err), "(Data) ibm_is_image_job", "read", "set-started_at").GetDiag()
		}
	}

	if err = d.Set("status", imageJob.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting status: %s", err), "(Data) ibm_is_image_job", "read", "set-status").GetDiag()
	}

	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range imageJob.StatusReasons {
		statusReasonsItemMap, err := DataSourceIBMIsImageJobImageJobStatusReasonsToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_job", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting status_reasons: %s", err), "(Data) ibm_is_image_job", "read", "set-status_reasons").GetDiag()
	}

	if err = d.Set("type", imageJob.Type); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting type: %s", err), "(Data) ibm_is_image_job", "read", "set-type").GetDiag()
	}

	if !core.IsNil(imageJob.DeleteImageJobOnImageDelete) {
		if err = d.Set("delete_image_job_on_image_delete", imageJob.DeleteImageJobOnImageDelete); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting delete_image_job_on_image_delete: %s", err), "(Data) ibm_is_image_job", "read", "set-delete_image_job_on_image_delete").GetDiag()
		}
	}

	if !core.IsNil(imageJob.Image) {
		image := []map[string]interface{}{}
		imageMap, err := DataSourceIBMIsImageJobImageReferenceToMap(imageJob.Image)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_job", "read", "image-to-map").GetDiag()
		}
		image = append(image, imageMap)
		if err = d.Set("image", image); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting image: %s", err), "(Data) ibm_is_image_job", "read", "set-image").GetDiag()
		}
	}

	if !core.IsNil(imageJob.ImageRecipe) {
		imageRecipe := []map[string]interface{}{}
		imageRecipeMap, err := DataSourceIBMIsImageJobImageJobBuildImageRecipeToMap(imageJob.ImageRecipe)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_job", "read", "image_recipe-to-map").GetDiag()
		}
		imageRecipe = append(imageRecipe, imageRecipeMap)
		if err = d.Set("image_recipe", imageRecipe); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting image_recipe: %s", err), "(Data) ibm_is_image_job", "read", "set-image_recipe").GetDiag()
		}
	}

	if !core.IsNil(imageJob.ImageTemplate) {
		imageTemplate := []map[string]interface{}{}
		imageTemplateMap, err := DataSourceIBMIsImageJobImageJobBuildImageTemplateToMap(imageJob.ImageTemplate)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_job", "read", "image_template-to-map").GetDiag()
		}
		imageTemplate = append(imageTemplate, imageTemplateMap)
		if err = d.Set("image_template", imageTemplate); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting image_template: %s", err), "(Data) ibm_is_image_job", "read", "set-image_template").GetDiag()
		}
	}

	if !core.IsNil(imageJob.Timeout) {
		if err = d.Set("timeout", flex.IntValue(imageJob.Timeout)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting timeout: %s", err), "(Data) ibm_is_image_job", "read", "set-timeout").GetDiag()
		}
	}

	if !core.IsNil(imageJob.Worker) {
		worker := []map[string]interface{}{}
		workerMap, err := DataSourceIBMIsImageJobImageJobBuildWorkerToMap(imageJob.Worker)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_job", "read", "worker-to-map").GetDiag()
		}
		worker = append(worker, workerMap)
		if err = d.Set("worker", worker); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting worker: %s", err), "(Data) ibm_is_image_job", "read", "set-worker").GetDiag()
		}
	}

	return nil
}

func DataSourceIBMIsImageJobImageJobStatusReasonsToMap(model *vpcv1.ImageJobStatusReasons) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func DataSourceIBMIsImageJobImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsImageJobImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageJobDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsImageJobImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := DataSourceIBMIsImageJobAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsImageJobRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsImageJobAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageJobRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageJobImageJobBuildImageRecipeToMap(model *vpcv1.ImageJobBuildImageRecipe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobDeletedToMap(model.Deleted)
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

func DataSourceIBMIsImageJobImageJobBuildImageTemplateToMap(model *vpcv1.ImageJobBuildImageTemplate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceGroup != nil {
		resourceGroupMap, err := DataSourceIBMIsImageJobResourceGroupIdentityToMap(model.ResourceGroup)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsImageJobResourceGroupIdentityToMap(model vpcv1.ResourceGroupIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ResourceGroupIdentityResourceGroupIdentityByID); ok {
		return DataSourceIBMIsImageJobResourceGroupIdentityResourceGroupIdentityByIDToMap(model.(*vpcv1.ResourceGroupIdentityResourceGroupIdentityByID))
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

func DataSourceIBMIsImageJobResourceGroupIdentityResourceGroupIdentityByIDToMap(model *vpcv1.ResourceGroupIdentityResourceGroupIdentityByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func DataSourceIBMIsImageJobImageJobBuildWorkerToMap(model *vpcv1.ImageJobBuildWorker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instance != nil {
		instanceMap, err := DataSourceIBMIsImageJobInstanceReferenceToMap(model.Instance)
		if err != nil {
			return modelMap, err
		}
		modelMap["instance"] = []map[string]interface{}{instanceMap}
	}
	profileMap, err := DataSourceIBMIsImageJobImageJobBuildWorkerProfileToMap(model.Profile)
	if err != nil {
		return modelMap, err
	}
	modelMap["profile"] = []map[string]interface{}{profileMap}
	subnetMap, err := DataSourceIBMIsImageJobSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	return modelMap, nil
}

func DataSourceIBMIsImageJobInstanceReferenceToMap(model *vpcv1.InstanceReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobDeletedToMap(model.Deleted)
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

func DataSourceIBMIsImageJobImageJobBuildWorkerProfileToMap(model *vpcv1.ImageJobBuildWorkerProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageJobSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageJobDeletedToMap(model.Deleted)
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
