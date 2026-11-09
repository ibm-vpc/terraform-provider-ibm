// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

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
		CreateContext: resourceIBMIsImageJobCreate,
		ReadContext:   resourceIBMIsImageJobRead,
		UpdateContext: resourceIBMIsImageJobUpdate,
		DeleteContext: resourceIBMIsImageJobDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"delete_image_job_on_image_delete": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Indicates whether deleting the image will also delete the image job that created it.",
			},
			"image_recipe": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the image recipe to use for this image job.",
			},
			"image_recipe_version": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The version of the image recipe to use. If unspecified, the latest version of the specified recipe will be used.",
			},
			"image_template_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name for the image created by this job. The name must not be used by another image in the region. Names starting with `ibm-` are reserved for system-provided images, and are not allowed. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"image_template_resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The unique identifier of the resource group to use for the image created by this job. If unspecified, the account's default resource group will be used.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_job", "name"),
				Description:  "The name for this image job. The name must not be used by another image job in the region. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      120,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_job", "timeout"),
				Description:  "The maximum time (in minutes) allowed for this image job to complete. The time begins when the job moves to a `status` of `running`. If the timeout is reached, the image and worker will be deleted.",
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_job", "type"),
				Description:  "The type of image job:\n- `build`: build an image from an image recipe\n\nThe enumerated values for this property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"worker_profile": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The name of the profile to use for the worker executing the image job. If unspecified, the `cxf-4x8` instance profile will be used.",
			},
			"worker_subnet": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the subnet the worker will be part of while executing the image job.",
			},
			// computed response fields
			"completed_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image job was completed. If absent, the image job has not yet completed.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image job was created.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this image job.",
			},
			"image_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the image created by this job, if the job has succeeded.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"started_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image job started running. If absent, the image job has not yet started.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of this image job:\n- `deleting`: The image job is being deleted\n- `failed`: The image job could not be completed successfully\n- `queued`: The image job has been accepted and is waiting to be scheduled\n- `running`: The image job is in progress\n- `succeeded`: The image job was completed successfully\n\nThe enumerated values for this property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"status_reasons": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The reasons for the current status (if any).",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A reason code for the image job status.",
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
			"worker_instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the worker instance, if it has been created.",
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

	recipeID := d.Get("image_recipe").(string)
	imageRecipePrototype := &vpcv1.ImageJobBuildImageRecipePrototypeImageRecipeByID{
		ID: &recipeID,
	}
	if v, ok := d.GetOk("image_recipe_version"); ok {
		version := v.(string)
		imageRecipePrototype.Version = &version
	}

	subnetID := d.Get("worker_subnet").(string)
	workerPrototype := &vpcv1.ImageJobBuildWorkerPrototype{
		Subnet: &vpcv1.ImageJobBuildWorkerSubnetPrototypeSubnetIdentityByID{
			ID: &subnetID,
		},
	}
	if v, ok := d.GetOk("worker_profile"); ok {
		profileName := v.(string)
		workerPrototype.Profile = &vpcv1.ImageJobBuildWorkerProfilePrototypeInstanceProfileIdentityByName{
			Name: &profileName,
		}
	}

	jobPrototype := &vpcv1.ImageJobPrototypeImageJobBuildPrototype{
		Type:        core.StringPtr(d.Get("type").(string)),
		ImageRecipe: imageRecipePrototype,
		Worker:      workerPrototype,
	}

	if v, ok := d.GetOk("name"); ok {
		name := v.(string)
		jobPrototype.Name = &name
	}
	if v, ok := d.GetOk("timeout"); ok {
		timeout := int64(v.(int))
		jobPrototype.Timeout = &timeout
	}
	if v, ok := d.GetOk("delete_image_job_on_image_delete"); ok {
		delete := v.(bool)
		jobPrototype.DeleteImageJobOnImageDelete = &delete
	}

	if v, ok := d.GetOk("image_template_name"); ok {
		tmplName := v.(string)
		imageTemplate := &vpcv1.ImageJobBuildImageTemplate{
			Name: &tmplName,
		}
		if rg, ok := d.GetOk("image_template_resource_group"); ok {
			rgID := rg.(string)
			imageTemplate.ResourceGroup = &vpcv1.ResourceGroupIdentity{
				ID: &rgID,
			}
		}
		jobPrototype.ImageTemplate = imageTemplate
	}

	createImageJobOptions := &vpcv1.CreateImageJobOptions{
		ImageJobPrototype: jobPrototype,
	}

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

	if !core.IsNil(imageJob.DeleteImageJobOnImageDelete) {
		if err = d.Set("delete_image_job_on_image_delete", imageJob.DeleteImageJobOnImageDelete); err != nil {
			err = fmt.Errorf("error setting delete_image_job_on_image_delete: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-delete_image_job_on_image_delete").GetDiag()
		}
	}
	if imageJob.ImageRecipe != nil {
		if err = d.Set("image_recipe", *imageJob.ImageRecipe.ID); err != nil {
			err = fmt.Errorf("error setting image_recipe: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-image_recipe").GetDiag()
		}
	}
	if !core.IsNil(imageJob.Name) {
		if err = d.Set("name", imageJob.Name); err != nil {
			err = fmt.Errorf("error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(imageJob.Timeout) {
		if err = d.Set("timeout", flex.IntValue(imageJob.Timeout)); err != nil {
			err = fmt.Errorf("error setting timeout: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-timeout").GetDiag()
		}
	}
	if err = d.Set("type", imageJob.Type); err != nil {
		err = fmt.Errorf("error setting type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-type").GetDiag()
	}
	if imageJob.Worker != nil && imageJob.Worker.Profile != nil {
		if err = d.Set("worker_profile", *imageJob.Worker.Profile.Name); err != nil {
			err = fmt.Errorf("error setting worker_profile: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-worker_profile").GetDiag()
		}
	}
	if imageJob.Worker != nil && imageJob.Worker.Subnet != nil {
		if err = d.Set("worker_subnet", *imageJob.Worker.Subnet.ID); err != nil {
			err = fmt.Errorf("error setting worker_subnet: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-worker_subnet").GetDiag()
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
	if !core.IsNil(imageJob.Image) {
		if err = d.Set("image_id", *imageJob.Image.ID); err != nil {
			err = fmt.Errorf("error setting image_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-image_id").GetDiag()
		}
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

	statusReasons := make([]map[string]interface{}, 0, len(imageJob.StatusReasons))
	for _, sr := range imageJob.StatusReasons {
		srMap := map[string]interface{}{
			"code":    *sr.Code,
			"message": *sr.Message,
		}
		if sr.MoreInfo != nil {
			srMap["more_info"] = *sr.MoreInfo
		}
		statusReasons = append(statusReasons, srMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		err = fmt.Errorf("error setting status_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-status_reasons").GetDiag()
	}

	if imageJob.Worker != nil && imageJob.Worker.Instance != nil {
		if err = d.Set("worker_instance_id", *imageJob.Worker.Instance.ID); err != nil {
			err = fmt.Errorf("error setting worker_instance_id: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_job", "read", "set-worker_instance_id").GetDiag()
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
		newDel := d.Get("delete_image_job_on_image_delete").(bool)
		patchVals.DeleteImageJobOnImageDelete = &newDel
		hasChange = true
	}

	if hasChange {
		updateImageJobOptions.ImageJobPatch, _ = patchVals.AsPatch()
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
