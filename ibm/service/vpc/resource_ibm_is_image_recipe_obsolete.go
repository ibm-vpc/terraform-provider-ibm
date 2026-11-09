// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceIBMIsImageRecipeObsolete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsImageRecipeObsoleteCreate,
		ReadContext:   resourceIBMIsImageRecipeObsoleteRead,
		DeleteContext: resourceIBMIsImageRecipeObsoleteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"image_recipe": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The image recipe identifier.",
			},
			"base_image": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The base image from which the image will be built.",
			},
			"image_components": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The image components included in this image recipe.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this image recipe.",
			},
			"deprecation_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deprecation date and time (UTC) for this image recipe. If absent, no deprecation date and time has been set.",
			},
			"obsolescence_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The obsolescence date and time (UTC) for this image recipe. If absent, no obsolescence date and time has been set.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource group for this image recipe.",
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
				Description: "The status of an image recipe.",
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
							Description: "A reason code for this status.",
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
				Description: "The date and time that this image recipe was superseded by a newer version.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image recipe was last updated.",
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A generated content fingerprint for this version of the image recipe.",
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The visibility of this image recipe.",
			},
		},
	}
}

func resourceIBMIsImageRecipeObsoleteCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	imageRecipeID := d.Get("image_recipe").(string)

	obsoleteImageRecipeOptions := &vpcv1.ObsoleteImageRecipeOptions{}
	obsoleteImageRecipeOptions.SetID(imageRecipeID)

	_, err = vpcClient.ObsoleteImageRecipeWithContext(context, obsoleteImageRecipeOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ObsoleteImageRecipeWithContext failed: %s", err.Error()), "ibm_is_image_recipe_obsolete", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(imageRecipeID)

	return resourceIBMIsImageRecipeObsoleteRead(context, d, meta)
}

func resourceIBMIsImageRecipeObsoleteRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "initialize-client")
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
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageRecipeWithContext failed: %s", err.Error()), "ibm_is_image_recipe_obsolete", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("image_recipe", d.Id()); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-image_recipe").GetDiag()
	}

	if imageRecipe.BaseImage != nil {
		if err = d.Set("base_image", *imageRecipe.BaseImage.ID); err != nil {
			err = fmt.Errorf("error setting base_image: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-base_image").GetDiag()
		}
	}

	imageComponentIDs := make([]string, 0, len(imageRecipe.ImageComponents))
	for i := 0; i < len(imageRecipe.ImageComponents); i++ {
		imageComponentIDs = append(imageComponentIDs, *imageRecipe.ImageComponents[i].ID)
	}
	if err = d.Set("image_components", imageComponentIDs); err != nil {
		err = fmt.Errorf("error setting image_components: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-image_components").GetDiag()
	}

	if !core.IsNil(imageRecipe.Name) {
		if err = d.Set("name", imageRecipe.Name); err != nil {
			err = fmt.Errorf("error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.DeprecationAt) {
		if err = d.Set("deprecation_at", flex.DateTimeToString(imageRecipe.DeprecationAt)); err != nil {
			err = fmt.Errorf("error setting deprecation_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-deprecation_at").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.ObsolescenceAt) {
		if err = d.Set("obsolescence_at", flex.DateTimeToString(imageRecipe.ObsolescenceAt)); err != nil {
			err = fmt.Errorf("error setting obsolescence_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-obsolescence_at").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.ResourceGroup) {
		if err = d.Set("resource_group", *imageRecipe.ResourceGroup.ID); err != nil {
			err = fmt.Errorf("error setting resource_group: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-resource_group").GetDiag()
		}
	}
	if err = d.Set("crn", imageRecipe.CRN); err != nil {
		err = fmt.Errorf("error setting crn: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-crn").GetDiag()
	}
	if err = d.Set("href", imageRecipe.Href); err != nil {
		err = fmt.Errorf("error setting href: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-href").GetDiag()
	}
	if err = d.Set("resource_type", imageRecipe.ResourceType); err != nil {
		err = fmt.Errorf("error setting resource_type: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-resource_type").GetDiag()
	}
	if err = d.Set("status", imageRecipe.Status); err != nil {
		err = fmt.Errorf("error setting status: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-status").GetDiag()
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
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-status_reasons").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(imageRecipe.CreatedAt)); err != nil {
		err = fmt.Errorf("error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-created_at").GetDiag()
	}
	if !core.IsNil(imageRecipe.SupersededAt) {
		if err = d.Set("superseded_at", flex.DateTimeToString(imageRecipe.SupersededAt)); err != nil {
			err = fmt.Errorf("error setting superseded_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-superseded_at").GetDiag()
		}
	}
	if err = d.Set("updated_at", flex.DateTimeToString(imageRecipe.UpdatedAt)); err != nil {
		err = fmt.Errorf("error setting updated_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-updated_at").GetDiag()
	}
	if err = d.Set("version", imageRecipe.Version); err != nil {
		err = fmt.Errorf("error setting version: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-version").GetDiag()
	}
	if err = d.Set("visibility", imageRecipe.Visibility); err != nil {
		err = fmt.Errorf("error setting visibility: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe_obsolete", "read", "set-visibility").GetDiag()
	}

	return nil
}

func resourceIBMIsImageRecipeObsoleteDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
