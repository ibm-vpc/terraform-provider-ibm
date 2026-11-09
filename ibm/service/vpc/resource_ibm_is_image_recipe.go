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

func ResourceIBMIsImageRecipe() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsImageRecipeCreate,
		ReadContext:   resourceIBMIsImageRecipeRead,
		UpdateContext: resourceIBMIsImageRecipeUpdate,
		DeleteContext: resourceIBMIsImageRecipeDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"base_image": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The base image from which the image will be built. The base image must meet the usage constraints of the image components specified in the image recipe, and the base image `visibility` must be `public`.",
			},
			"image_components": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				MaxItems:    20,
				Description: "The image components to include in this image recipe.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_recipe", "name"),
				Description:  "The name for this image recipe. The name must not be used by another image recipe in the region. Names starting with `ibm-` are reserved for system-provided image recipes; If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"deprecation_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deprecation date and time (UTC) for this image recipe. If absent, no deprecation date and time has been set. Specify \"null\" to remove an existing deprecation date and time.",
			},
			"obsolescence_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The obsolescence date and time (UTC) for this image recipe. If absent, no obsolescence date and time has been set. Specify \"null\" to remove an existing obsolescence date and time.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The resource group for this image recipe. Uses account default if omitted.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this image recipe.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this image recipe.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of an image recipe.\n- `available`: The recipe can be used to create images.\n- `degraded`:  The recipe is degraded, and cannot be used.\n- `deleting`: The recipe is being deleted, and can no longer be used.\n- `deprecated`: The recipe is deprecated, but can still be used.\n- `obsolete`: The recipe can no longer be used to create images.\n- `pending`: The recipe is being created, and is not ready to be used.\n\nThe enumerated values for this property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
							Description: "A reason code for this status:\n- `base_image_unusable`: The base image for this recipe is in an `obsolete` state or has been deleted. Image creation from this recipe will fail until the recipe no longer has an unusable base image.\n- `image_component_unusable`: One or more image components referenced by this image recipe have a `status` of `degraded`, `obsolete`, or have been deleted. Image creation from this recipe will fail until the recipe no longer has unusable image components.\n- `internal_error`: internal error (contact IBM support)\n\nThe enumerated values for this property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image recipe was created.",
			},
			"superseded_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this image recipe was superseded by a newer version.\n\nIf absent, the recipe has not been superseded by a newer version.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image recipe was last updated.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A generated content fingerprint for this version of the image recipe. The value changes whenever the image recipe is updated.",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The visibility of this image recipe.\n- `private`: Visible only to this account\n- `public`: Visible to all accounts",
			},
		},
	}
}

func ResourceIBMIsImageRecipeValidator() *validate.ResourceValidator {
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
	)
	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_image_recipe", Schema: validateSchema}
	return &resourceValidator
}

func resourceIBMIsImageRecipeCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	createImageRecipeOptions := &vpcv1.CreateImageRecipeOptions{}

	baseImageID := d.Get("base_image").(string)
	createImageRecipeOptions.SetBaseImage(&vpcv1.ImageRecipeBaseImagePrototypeImageIdentityImageIdentityByID{
		ID: &baseImageID,
	})

	var imageComponents []vpcv1.ImageRecipeImageComponentPrototypeIntf
	for _, v := range d.Get("image_components").([]interface{}) {
		id := v.(string)
		imageComponents = append(imageComponents, &vpcv1.ImageRecipeImageComponentPrototypeImageComponentIdentityImageComponentIdentityImageComponentByID{
			ID: &id,
		})
	}
	createImageRecipeOptions.SetImageComponents(imageComponents)

	if _, ok := d.GetOk("name"); ok {
		createImageRecipeOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("deprecation_at"); ok {
		fmtDateTimeDeprecationAt, err := core.ParseDateTime(d.Get("deprecation_at").(string))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "parse-deprecation_at").GetDiag()
		}
		createImageRecipeOptions.SetDeprecationAt(&fmtDateTimeDeprecationAt)
	}
	if _, ok := d.GetOk("obsolescence_at"); ok {
		fmtDateTimeObsolescenceAt, err := core.ParseDateTime(d.Get("obsolescence_at").(string))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "parse-obsolescence_at").GetDiag()
		}
		createImageRecipeOptions.SetObsolescenceAt(&fmtDateTimeObsolescenceAt)
	}
	if grp, ok := d.GetOk("resource_group"); ok {
		rgID := grp.(string)
		createImageRecipeOptions.SetResourceGroup(&vpcv1.ImageRecipeResourceGroupPrototypeResourceGroupIdentityByID{
			ID: &rgID,
		})
	}

	imageRecipe, _, err := vpcClient.CreateImageRecipeWithContext(context, createImageRecipeOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("CreateImageRecipeWithContext failed: %s", err.Error()), "ibm_is_image_recipe", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*imageRecipe.ID)

	return resourceIBMIsImageRecipeRead(context, d, meta)
}

func resourceIBMIsImageRecipeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getImageRecipeOptions := &vpcv1.GetImageRecipeOptions{}
	getImageRecipeOptions.SetID(d.Id())

	imageRecipe, response, err := vpcClient.GetImageRecipeWithContext(context, getImageRecipeOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageRecipeWithContext failed: %s", err.Error()), "ibm_is_image_recipe", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if imageRecipe.BaseImage != nil {
		if err = d.Set("base_image", *imageRecipe.BaseImage.ID); err != nil {
			err = fmt.Errorf("error setting base_image: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-base_image").GetDiag()
		}
	}

	imageComponentIDs := make([]string, 0, len(imageRecipe.ImageComponents))
	for i := 0; i < len(imageRecipe.ImageComponents); i++ {
		imageComponentIDs = append(imageComponentIDs, *imageRecipe.ImageComponents[i].ID)
	}
	if err = d.Set("image_components", imageComponentIDs); err != nil {
		err = fmt.Errorf("error setting image_components: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-image_components").GetDiag()
	}
	if !core.IsNil(imageRecipe.Name) {
		if err = d.Set("name", imageRecipe.Name); err != nil {
			err = fmt.Errorf("error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.DeprecationAt) {
		if err = d.Set("deprecation_at", flex.DateTimeToString(imageRecipe.DeprecationAt)); err != nil {
			err = fmt.Errorf("error setting deprecation_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-deprecation_at").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.ObsolescenceAt) {
		if err = d.Set("obsolescence_at", flex.DateTimeToString(imageRecipe.ObsolescenceAt)); err != nil {
			err = fmt.Errorf("error setting obsolescence_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-obsolescence_at").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.ResourceGroup) {
		if err = d.Set("resource_group", *imageRecipe.ResourceGroup.ID); err != nil {
			err = fmt.Errorf("error setting resource_group: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-resource_group").GetDiag()
		}
	}
	if err = d.Set("crn", imageRecipe.CRN); err != nil {
		err = fmt.Errorf("error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-crn").GetDiag()
	}
	if err = d.Set("href", imageRecipe.Href); err != nil {
		err = fmt.Errorf("error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-href").GetDiag()
	}
	if err = d.Set("resource_type", imageRecipe.ResourceType); err != nil {
		err = fmt.Errorf("error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-resource_type").GetDiag()
	}
	if err = d.Set("status", imageRecipe.Status); err != nil {
		err = fmt.Errorf("error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-status").GetDiag()
	}

	statusReasons := make([]map[string]interface{}, 0, len(imageRecipe.StatusReasons))
	for _, sr := range imageRecipe.StatusReasons {
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
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-status_reasons").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(imageRecipe.CreatedAt)); err != nil {
		err = fmt.Errorf("error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-created_at").GetDiag()
	}
	if !core.IsNil(imageRecipe.SupersededAt) {
		if err = d.Set("superseded_at", flex.DateTimeToString(imageRecipe.SupersededAt)); err != nil {
			err = fmt.Errorf("error setting superseded_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-superseded_at").GetDiag()
		}
	}
	if err = d.Set("updated_at", flex.DateTimeToString(imageRecipe.UpdatedAt)); err != nil {
		err = fmt.Errorf("error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("version", imageRecipe.Version); err != nil {
		err = fmt.Errorf("error setting version: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-version").GetDiag()
	}
	if err = d.Set("visibility", imageRecipe.Visibility); err != nil {
		err = fmt.Errorf("error setting visibility: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-visibility").GetDiag()
	}

	return nil
}

func resourceIBMIsImageRecipeUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	updateImageRecipeOptions := &vpcv1.UpdateImageRecipeOptions{}
	updateImageRecipeOptions.SetID(d.Id())

	hasChange := false
	patchVals := &vpcv1.ImageRecipePatch{}

	if d.HasChange("base_image") {
		baseImageID := d.Get("base_image").(string)
		patchVals.BaseImage = &vpcv1.ImageRecipeBaseImagePatchImageIdentityImageIdentityByID{
			ID: &baseImageID,
		}
		hasChange = true
	}
	if d.HasChange("image_components") {
		var imageComponents []vpcv1.ImageRecipeImageComponentPrototypeIntf
		for _, v := range d.Get("image_components").([]interface{}) {
			id := v.(string)
			imageComponents = append(imageComponents, &vpcv1.ImageRecipeImageComponentPrototypeImageComponentIdentityImageComponentIdentityImageComponentByID{
				ID: &id,
			})
		}
		patchVals.ImageComponents = imageComponents
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	nullDeprecate := false
	if d.HasChange("deprecation_at") {
		deprecationAt := d.Get("deprecation_at").(string)
		if deprecationAt == "null" {
			nullDeprecate = true
		} else {
			newDeprecationAt, err := core.ParseDateTime(deprecationAt)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "update", "parse-deprecation_at").GetDiag()
			}
			patchVals.DeprecationAt = &newDeprecationAt
		}
		hasChange = true
	}
	nullObsolescence := false
	if d.HasChange("obsolescence_at") {
		obsolescenceAt := d.Get("obsolescence_at").(string)
		if obsolescenceAt == "null" {
			nullObsolescence = true
		} else {
			newObsolescenceAt, err := core.ParseDateTime(obsolescenceAt)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "update", "parse-obsolescence_at").GetDiag()
			}
			patchVals.ObsolescenceAt = &newObsolescenceAt
		}
		hasChange = true
	}

	if hasChange {
		imageRecipePatch, err := patchVals.AsPatch()
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("patchVals.AsPatch failed: %s", err.Error()), "ibm_is_image_recipe", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		if nullDeprecate {
			imageRecipePatch["deprecation_at"] = nil
		}
		if nullObsolescence {
			imageRecipePatch["obsolescence_at"] = nil
		}
		updateImageRecipeOptions.ImageRecipePatch = imageRecipePatch
		_, _, err = vpcClient.UpdateImageRecipeWithContext(context, updateImageRecipeOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateImageRecipeWithContext failed: %s", err.Error()), "ibm_is_image_recipe", "update")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}

	return resourceIBMIsImageRecipeRead(context, d, meta)
}

func resourceIBMIsImageRecipeDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	deleteImageRecipeOptions := &vpcv1.DeleteImageRecipeOptions{}
	deleteImageRecipeOptions.SetID(d.Id())

	_, err = vpcClient.DeleteImageRecipeWithContext(context, deleteImageRecipeOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("DeleteImageRecipeWithContext failed: %s", err.Error()), "ibm_is_image_recipe", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")

	return nil
}
