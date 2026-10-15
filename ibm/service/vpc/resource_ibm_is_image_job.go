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
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsImageJob() *schema.Resource {
	return &schema.Resource{
		CreateContext:   resourceIBMIsImageJobCreate,
		ReadContext:     resourceIBMIsImageJobRead,
		UpdateContext:   resourceIBMIsImageJobUpdate,
		DeleteContext:   resourceIBMIsImageJobDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_job", "name"),
				Description: "The name for this image job. The name is unique across all image jobs in the region.",
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_job", "type"),
				Description: "The type of image job:- `build`: build an image from an image recipeThe enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"delete_image_job_on_image_delete": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether deleting the image will also delete the image job that created it.",
			},
			"image_recipe": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The image recipe and version used by this image job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
							Required:    true,
							Description: "The URL for this image recipe.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
							Required:    true,
							Description: "A generated content fingerprint for this version of the image recipe. The value changes whenever the image recipe is updated.",
						},
					},
				},
			},
			"image_template": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The template for creating an image from the image job.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name for this image. The name must not be used by another image in the region. Names starting with `ibm-` are reserved for system-provided images, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "The resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) will be used.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Optional:    true,
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
				Optional:    true,
				Default:     120,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_job", "timeout"),
				Description: "The maximum time (in minutes) allowed for this image job to complete.The time begins when the job moves to a `status` of `running`. If the timeout is reached, the image and worker will be deleted.",
			},
			"worker": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
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
										Optional:    true,
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
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Description: "The [profile](https://cloud.ibm.com/docs/vpc?topic=vpc-profiles) used to build the image.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "The URL for this virtual server instance profile.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
							MinItems:    1,
							MaxItems:    1,
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"crn": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
										Required:    true,
										Description: "The URL for this subnet.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
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
							Optional:    true,
							Computed:    true,
							Description: "A link to documentation about this status reason.",
						},
					},
				},
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
							Optional:    true,
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
							Optional:    true,
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
		},
	}
}

func ResourceIBMIsImageJobValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "type",
			ValidateFunctionIdentifier: validate.ValidateAllowedStringValue,
			Type:                       validate.TypeString,
			Required:                   true,
			AllowedValues:              "build",
			Regexp:                     `^[a-z][a-z0-9]*(_[a-z0-9]+)*$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "timeout",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "30",
			MaxValue:                   "480",
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_image_job", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsImageJobCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	bodyModelMap := map[string]interface{}{}
	createImageJobOptions := &vpcv1.CreateImageJobOptions{}

	if _, ok := d.GetOk("name"); ok {
		bodyModelMap["name"] = d.Get("name")
	}
	bodyModelMap["type"] = d.Get("type")
	if _, ok := d.GetOk("delete_image_job_on_image_delete"); ok {
		bodyModelMap["delete_image_job_on_image_delete"] = d.Get("delete_image_job_on_image_delete")
	}
	if _, ok := d.GetOk("image_recipe"); ok {
		bodyModelMap["image_recipe"] = d.Get("image_recipe")
	}
	if _, ok := d.GetOk("image_template"); ok {
		bodyModelMap["image_template"] = d.Get("image_template")
	}
	if _, ok := d.GetOk("timeout"); ok {
		bodyModelMap["timeout"] = d.Get("timeout")
	}
	if _, ok := d.GetOk("worker"); ok {
		bodyModelMap["worker"] = d.Get("worker")
	}
	convertedModel, err := ResourceIBMIsImageJobMapToImageJobPrototype(bodyModelMap)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "create", "parse-request-body").GetDiag()
	}
	createImageJobOptions.ImageJobPrototype = convertedModel

	imageJobIntf, _, err := vpcClient.CreateImageJobWithContext(context, createImageJobOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateImageJobWithContext failed: %s", err.Error()), "ibm_is_image_job", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	imageJob := imageJobIntf.(*vpcv1.ImageJob)
	d.SetId(*imageJob.ID)

	return resourceIBMIsImageJobRead(context, d, meta)
}

func resourceIBMIsImageJobRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getImageJobOptions := &vpcv1.GetImageJobOptions{}

	getImageJobOptions.SetID(d.Id())

	imageJobIntf, response, err := vpcClient.GetImageJobWithContext(context, getImageJobOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageJobWithContext failed: %s", err.Error()), "ibm_is_image_job", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	imageJob := imageJobIntf.(*vpcv1.ImageJob)
	if !core.IsNil(imageJob.Name) {
		if err = d.Set("name", imageJob.Name); err != nil {
			err = fmt.Errorf("error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-name").GetDiag()
		}
	}
	if err = d.Set("type", imageJob.Type); err != nil {
		err = fmt.Errorf("error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-type").GetDiag()
	}
	if !core.IsNil(imageJob.DeleteImageJobOnImageDelete) {
		if err = d.Set("delete_image_job_on_image_delete", imageJob.DeleteImageJobOnImageDelete); err != nil {
			err = fmt.Errorf("error setting delete_image_job_on_image_delete: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-delete_image_job_on_image_delete").GetDiag()
		}
	}
	if !core.IsNil(imageJob.ImageRecipe) {
		imageRecipeMap, err := ResourceIBMIsImageJobImageJobBuildImageRecipeToMap(imageJob.ImageRecipe)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "image_recipe-to-map").GetDiag()
		}
		if err = d.Set("image_recipe", []map[string]interface{}{imageRecipeMap}); err != nil {
			err = fmt.Errorf("error setting image_recipe: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-image_recipe").GetDiag()
		}
	}
	if !core.IsNil(imageJob.ImageTemplate) {
		imageTemplateMap, err := ResourceIBMIsImageJobImageJobBuildImageTemplateToMap(imageJob.ImageTemplate)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "image_template-to-map").GetDiag()
		}
		if err = d.Set("image_template", []map[string]interface{}{imageTemplateMap}); err != nil {
			err = fmt.Errorf("error setting image_template: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-image_template").GetDiag()
		}
	}
	if !core.IsNil(imageJob.Timeout) {
		if err = d.Set("timeout", flex.IntValue(imageJob.Timeout)); err != nil {
			err = fmt.Errorf("error setting timeout: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-timeout").GetDiag()
		}
	}
	if !core.IsNil(imageJob.Worker) {
		workerMap, err := ResourceIBMIsImageJobImageJobBuildWorkerToMap(imageJob.Worker)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "worker-to-map").GetDiag()
		}
		if err = d.Set("worker", []map[string]interface{}{workerMap}); err != nil {
			err = fmt.Errorf("error setting worker: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-worker").GetDiag()
		}
	}
	if !core.IsNil(imageJob.CompletedAt) {
		if err = d.Set("completed_at", flex.DateTimeToString(imageJob.CompletedAt)); err != nil {
			err = fmt.Errorf("error setting completed_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-completed_at").GetDiag()
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(imageJob.CreatedAt)); err != nil {
		err = fmt.Errorf("error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-created_at").GetDiag()
	}
	if err = d.Set("href", imageJob.Href); err != nil {
		err = fmt.Errorf("error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-href").GetDiag()
	}
	if err = d.Set("resource_type", imageJob.ResourceType); err != nil {
		err = fmt.Errorf("error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-resource_type").GetDiag()
	}
	if !core.IsNil(imageJob.StartedAt) {
		if err = d.Set("started_at", flex.DateTimeToString(imageJob.StartedAt)); err != nil {
			err = fmt.Errorf("error setting started_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-started_at").GetDiag()
		}
	}
	if err = d.Set("status", imageJob.Status); err != nil {
		err = fmt.Errorf("error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-status").GetDiag()
	}
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range imageJob.StatusReasons {
		statusReasonsItemMap, err := ResourceIBMIsImageJobImageJobStatusReasonsToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		err = fmt.Errorf("error setting status_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-status_reasons").GetDiag()
	}
	if !core.IsNil(imageJob.Image) {
		imageMap, err := ResourceIBMIsImageJobImageReferenceToMap(imageJob.Image)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "image-to-map").GetDiag()
		}
		if err = d.Set("image", []map[string]interface{}{imageMap}); err != nil {
			err = fmt.Errorf("error setting image: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-image").GetDiag()
		}
	}

	return nil
}

func resourceIBMIsImageJobUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateImageJobOptions := &vpcv1.UpdateImageJobOptions{}

	updateImageJobOptions.SetID(d.Id())

	hasChange := false

	patchVals := &vpcv1.ImageJobPatch{}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("delete_image_job_on_image_delete") {
		newDeleteImageJobOnImageDelete := d.Get("delete_image_job_on_image_delete").(bool)
		patchVals.DeleteImageJobOnImageDelete = &newDeleteImageJobOnImageDelete
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateImageJobOptions.ImageJobPatch = ResourceIBMIsImageJobImageJobPatchAsPatch(patchVals, d)

		_, _, err = vpcClient.UpdateImageJobWithContext(context, updateImageJobOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateImageJobWithContext failed: %s", err.Error()), "ibm_is_image_job", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsImageJobRead(context, d, meta)
}

func resourceIBMIsImageJobDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteImageJobOptions := &vpcv1.DeleteImageJobOptions{}

	deleteImageJobOptions.SetID(d.Id())

	_, err = vpcClient.DeleteImageJobWithContext(context, deleteImageJobOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteImageJobWithContext failed: %s", err.Error()), "ibm_is_image_job", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}

func ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototype(modelMap map[string]interface{}) (vpcv1.ImageJobBuildImageRecipePrototypeIntf, error) {
	model := &vpcv1.ImageJobBuildImageRecipePrototype{}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
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

func ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByID(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByID, error) {
	model := &vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByID{}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByCRN(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByCRN, error) {
	model := &vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByCRN{}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototypeImageRecipeByHref(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByHref, error) {
	model := &vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByHref{}
	if modelMap["version"] != nil && modelMap["version"].(string) != "" {
		model.Version = core.StringPtr(modelMap["version"].(string))
	}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildImageTemplate(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildImageTemplate, error) {
	model := &vpcv1.ImageJobBuildImageTemplate{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["resource_group"] != nil && len(modelMap["resource_group"].([]interface{})) > 0 {
		ResourceGroupModel, err := ResourceIBMIsImageJobMapToResourceGroupIdentity(modelMap["resource_group"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ResourceGroup = ResourceGroupModel
	}
	return model, nil
}

func ResourceIBMIsImageJobMapToResourceGroupIdentity(modelMap map[string]interface{}) (vpcv1.ResourceGroupIdentityIntf, error) {
	model := &vpcv1.ResourceGroupIdentity{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func ResourceIBMIsImageJobMapToResourceGroupIdentityResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ResourceGroupIdentityResourceGroupIdentityByID, error) {
	model := &vpcv1.ResourceGroupIdentityResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildWorkerPrototype(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildWorkerPrototype, error) {
	model := &vpcv1.ImageJobBuildWorkerPrototype{}
	if modelMap["profile"] != nil && len(modelMap["profile"].([]interface{})) > 0 {
		ProfileModel, err := ResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototype(modelMap["profile"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Profile = ProfileModel
	}
	SubnetModel, err := ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototype(modelMap["subnet"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Subnet = SubnetModel
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototype(modelMap map[string]interface{}) (vpcv1.ImageJobBuildWorkerProfilePrototypeIntf, error) {
	model := &vpcv1.ImageJobBuildWorkerProfilePrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["href"] != nil && modelMap["href"].(string) != "" {
		model.Href = core.StringPtr(modelMap["href"].(string))
	}
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName, error) {
	model := &vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName{}
	model.Name = core.StringPtr(modelMap["name"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByHref, error) {
	model := &vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototype(modelMap map[string]interface{}) (vpcv1.ImageJobBuildWorkerSubnetPrototypeIntf, error) {
	model := &vpcv1.ImageJobBuildWorkerSubnetPrototype{}
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

func ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID, error) {
	model := &vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByCRN, error) {
	model := &vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobBuildWorkerSubnetPrototypeSubnetIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByHref, error) {
	model := &vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobPrototype(modelMap map[string]interface{}) (vpcv1.ImageJobPrototypeIntf, error) {
	model := &vpcv1.ImageJobPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	model.Type = core.StringPtr(modelMap["type"].(string))
	if modelMap["delete_image_job_on_image_delete"] != nil {
		model.DeleteImageJobOnImageDelete = core.BoolPtr(modelMap["delete_image_job_on_image_delete"].(bool))
	}
	if modelMap["image_recipe"] != nil && len(modelMap["image_recipe"].([]interface{})) > 0 {
		ImageRecipeModel, err := ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototype(modelMap["image_recipe"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ImageRecipe = ImageRecipeModel
	}
	if modelMap["image_template"] != nil && len(modelMap["image_template"].([]interface{})) > 0 {
		ImageTemplateModel, err := ResourceIBMIsImageJobMapToImageJobBuildImageTemplate(modelMap["image_template"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ImageTemplate = ImageTemplateModel
	}
	if modelMap["timeout"] != nil {
		model.Timeout = core.Int64Ptr(int64(modelMap["timeout"].(int)))
	}
	if modelMap["worker"] != nil && len(modelMap["worker"].([]interface{})) > 0 {
		WorkerModel, err := ResourceIBMIsImageJobMapToImageJobBuildWorkerPrototype(modelMap["worker"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.Worker = WorkerModel
	}
	return model, nil
}

func ResourceIBMIsImageJobMapToImageJobPrototypeImageJobBuildPrototype(modelMap map[string]interface{}) (*vpcv1.ImageJobPrototypeImageJobBuildPrototype, error) {
	model := &vpcv1.ImageJobPrototypeImageJobBuildPrototype{}
	if modelMap["name"] != nil && modelMap["name"].(string) != "" {
		model.Name = core.StringPtr(modelMap["name"].(string))
	}
	if modelMap["delete_image_job_on_image_delete"] != nil {
		model.DeleteImageJobOnImageDelete = core.BoolPtr(modelMap["delete_image_job_on_image_delete"].(bool))
	}
	ImageRecipeModel, err := ResourceIBMIsImageJobMapToImageJobBuildImageRecipePrototype(modelMap["image_recipe"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.ImageRecipe = ImageRecipeModel
	if modelMap["image_template"] != nil && len(modelMap["image_template"].([]interface{})) > 0 {
		ImageTemplateModel, err := ResourceIBMIsImageJobMapToImageJobBuildImageTemplate(modelMap["image_template"].([]interface{})[0].(map[string]interface{}))
		if err != nil {
			return model, err
		}
		model.ImageTemplate = ImageTemplateModel
	}
	if modelMap["timeout"] != nil {
		model.Timeout = core.Int64Ptr(int64(modelMap["timeout"].(int)))
	}
	model.Type = core.StringPtr(modelMap["type"].(string))
	WorkerModel, err := ResourceIBMIsImageJobMapToImageJobBuildWorkerPrototype(modelMap["worker"].([]interface{})[0].(map[string]interface{}))
	if err != nil {
		return model, err
	}
	model.Worker = WorkerModel
	return model, nil
}

func ResourceIBMIsImageJobImageJobBuildImageRecipeToMap(model *vpcv1.ImageJobBuildImageRecipe) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsImageJobDeletedToMap(model.Deleted)
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

func ResourceIBMIsImageJobDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsImageJobImageJobBuildImageTemplateToMap(model *vpcv1.ImageJobBuildImageTemplate) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Name != nil {
		modelMap["name"] = *model.Name
	}
	if model.ResourceGroup != nil {
		resourceGroupMap, err := ResourceIBMIsImageJobResourceGroupIdentityToMap(model.ResourceGroup)
		if err != nil {
			return modelMap, err
		}
		modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	}
	return modelMap, nil
}

func ResourceIBMIsImageJobResourceGroupIdentityToMap(model vpcv1.ResourceGroupIdentityIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ResourceGroupIdentityResourceGroupIdentityByID); ok {
		return ResourceIBMIsImageJobResourceGroupIdentityResourceGroupIdentityByIDToMap(model.(*vpcv1.ResourceGroupIdentityResourceGroupIdentityByID))
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

func ResourceIBMIsImageJobResourceGroupIdentityResourceGroupIdentityByIDToMap(model *vpcv1.ResourceGroupIdentityResourceGroupIdentityByID) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	return modelMap, nil
}

func ResourceIBMIsImageJobImageJobBuildWorkerToMap(model *vpcv1.ImageJobBuildWorker) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Instance != nil {
		instanceMap, err := ResourceIBMIsImageJobInstanceReferenceToMap(model.Instance)
		if err != nil {
			return modelMap, err
		}
		modelMap["instance"] = []map[string]interface{}{instanceMap}
	}
	profileMap, err := ResourceIBMIsImageJobImageJobBuildWorkerProfileToMap(model.Profile)
	if err != nil {
		return modelMap, err
	}
	modelMap["profile"] = []map[string]interface{}{profileMap}
	subnetMap, err := ResourceIBMIsImageJobSubnetReferenceToMap(model.Subnet)
	if err != nil {
		return modelMap, err
	}
	modelMap["subnet"] = []map[string]interface{}{subnetMap}
	return modelMap, nil
}

func ResourceIBMIsImageJobInstanceReferenceToMap(model *vpcv1.InstanceReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsImageJobDeletedToMap(model.Deleted)
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

func ResourceIBMIsImageJobImageJobBuildWorkerProfileToMap(model *vpcv1.ImageJobBuildWorkerProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsImageJobSubnetReferenceToMap(model *vpcv1.SubnetReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsImageJobDeletedToMap(model.Deleted)
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

func ResourceIBMIsImageJobImageJobStatusReasonsToMap(model *vpcv1.ImageJobStatusReasons) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsImageJobImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsImageJobDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := ResourceIBMIsImageJobImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsImageJobImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := ResourceIBMIsImageJobAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := ResourceIBMIsImageJobRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func ResourceIBMIsImageJobAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsImageJobRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsImageJobImageJobPatchAsPatch(patchVals *vpcv1.ImageJobPatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	} else if !exists {
		delete(patch, "name")
	}
	path = "delete_image_job_on_image_delete"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["delete_image_job_on_image_delete"] = nil
	} else if !exists {
		delete(patch, "delete_image_job_on_image_delete")
	}

	return patch
}
