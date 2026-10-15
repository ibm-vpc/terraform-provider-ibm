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

func ResourceIBMIsImageRecipe() *schema.Resource {
	return &schema.Resource{
		CreateContext:   resourceIBMIsImageRecipeCreate,
		ReadContext:     resourceIBMIsImageRecipeRead,
		UpdateContext:   resourceIBMIsImageRecipeUpdate,
		DeleteContext:   resourceIBMIsImageRecipeDelete,
		Importer: &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"base_image": &schema.Schema{
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    1,
				Required:    true,
				Description: "The base image that is extended when building an image from this image recipe.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
							Required:    true,
							Description: "The URL for this image.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
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
			"deprecation_at": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The deprecation date and time (UTC) for this image recipe.If absent, no deprecation date and time has been set.",
			},
			"image_components": &schema.Schema{
				Type:        schema.TypeList,
				Required:    true,
				Description: "The image components included in this image recipe.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The CRN for this component.",
						},
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "If present, this property indicates the referenced resource has been deleted, and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Required:    true,
										Description: "A link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The URL for this component.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this component.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique name for this image component.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc: validate.InvokeValidator("ibm_is_image_recipe", "name"),
				Description: "The unique name for this image recipe.",
			},
			"obsolescence_at": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The obsolescence date and time (UTC) for this image recipe.If absent, no obsolescence date and time has been set.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Description: "The resource group for this image recipe.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this resource group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this resource group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this resource group.",
						},
					},
				},
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image recipe was created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this image recipe.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this image recipe.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"status": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of an image recipe.- `available`: The recipe can be used to create images.- `degraded`:  The recipe is degraded, and cannot be used.- `deleting`: The recipe is being deleted, and can no longer be used.- `deprecated`: The recipe is deprecated, but can still be used.- `obsolete`: The recipe can no longer be used to create images.- `pending`: The recipe is being created, and is not ready to be used.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
							Description: "A reason code for this status:- `base_image_unusable`: The base image for this recipe is in an `obsolete` state or  has been deleted. Image creation from this recipe will fail until the recipe no  longer has an unusable base image.- `image_component_unusable`: One or more image components referenced by this image  recipe have a `status` of `degraded`, `obsolete`, or have been deleted. Image  creation from this recipe will fail until the recipe no longer has unusable image  components.- `internal_error`: internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
			"superseded_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that this image recipe was superseded by a newer version.If absent, the recipe has not been superseded by a newer version.",
			},
			"updated_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the image recipe was last updated.",
			},
			"version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "A generated content fingerprint for this version of the image recipe. The value changes whenever the image recipe is updated.",
			},
			"visibility": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The visibility of this image recipe.- `private`: Visible only to this account- `public`: Visible to all accounts.",
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

	baseImageModel, err := ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototype(d.Get("base_image.0").(map[string]interface{}))
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "parse-base_image").GetDiag()
	}
	createImageRecipeOptions.SetBaseImage(baseImageModel)
	var imageComponents []vpcv1.ImageComponentIdentityIntf
	for _, v := range d.Get("image_components").([]interface{}) {
		value := v.(map[string]interface{})
		imageComponentsItem, err := ResourceIBMIsImageRecipeMapToImageComponentIdentity(value)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "parse-image_components").GetDiag()
		}
		imageComponents = append(imageComponents, imageComponentsItem)
	}
	createImageRecipeOptions.SetImageComponents(imageComponents)
	if _, ok := d.GetOk("deprecation_at"); ok {
		fmtDateTimeDeprecationAt, err := core.ParseDateTime(d.Get("deprecation_at").(string))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "parse-deprecation_at").GetDiag()
		}
		createImageRecipeOptions.SetDeprecationAt(&fmtDateTimeDeprecationAt)
	}
	if _, ok := d.GetOk("name"); ok {
		createImageRecipeOptions.SetName(d.Get("name").(string))
	}
	if _, ok := d.GetOk("obsolescence_at"); ok {
		fmtDateTimeObsolescenceAt, err := core.ParseDateTime(d.Get("obsolescence_at").(string))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "parse-obsolescence_at").GetDiag()
		}
		createImageRecipeOptions.SetObsolescenceAt(&fmtDateTimeObsolescenceAt)
	}
	if _, ok := d.GetOk("resource_group"); ok {
		resourceGroupModel, err := ResourceIBMIsImageRecipeMapToImageRecipeResourceGroupPrototype(d.Get("resource_group.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "create", "parse-resource_group").GetDiag()
		}
		createImageRecipeOptions.SetResourceGroup(resourceGroupModel)
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

	baseImageMap, err := ResourceIBMIsImageRecipeImageReferenceToMap(imageRecipe.BaseImage)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "base_image-to-map").GetDiag()
	}
	if err = d.Set("base_image", []map[string]interface{}{baseImageMap}); err != nil {
		err = fmt.Errorf("error setting base_image: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-base_image").GetDiag()
	}
	if !core.IsNil(imageRecipe.DeprecationAt) {
		if err = d.Set("deprecation_at", flex.DateTimeToString(imageRecipe.DeprecationAt)); err != nil {
			err = fmt.Errorf("error setting deprecation_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-deprecation_at").GetDiag()
		}
	}
	imageComponents := []map[string]interface{}{}
	for _, imageComponentsItem := range imageRecipe.ImageComponents {
		imageComponentsItemMap, err := ResourceIBMIsImageRecipeImageComponentReferenceToMap(&imageComponentsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "image_components-to-map").GetDiag()
		}
		imageComponents = append(imageComponents, imageComponentsItemMap)
	}
	if err = d.Set("image_components", imageComponents); err != nil {
		err = fmt.Errorf("error setting image_components: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-image_components").GetDiag()
	}
	if !core.IsNil(imageRecipe.Name) {
		if err = d.Set("name", imageRecipe.Name); err != nil {
			err = fmt.Errorf("error setting name: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-name").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.ObsolescenceAt) {
		if err = d.Set("obsolescence_at", flex.DateTimeToString(imageRecipe.ObsolescenceAt)); err != nil {
			err = fmt.Errorf("error setting obsolescence_at: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-obsolescence_at").GetDiag()
		}
	}
	if !core.IsNil(imageRecipe.ResourceGroup) {
		resourceGroupMap, err := ResourceIBMIsImageRecipeImageRecipeResourceGroupToMap(imageRecipe.ResourceGroup)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "resource_group-to-map").GetDiag()
		}
		if err = d.Set("resource_group", []map[string]interface{}{resourceGroupMap}); err != nil {
			err = fmt.Errorf("error setting resource_group: %s", err)
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-resource_group").GetDiag()
		}
	}
	if err = d.Set("created_at", flex.DateTimeToString(imageRecipe.CreatedAt)); err != nil {
		err = fmt.Errorf("error setting created_at: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-created_at").GetDiag()
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
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range imageRecipe.StatusReasons {
		statusReasonsItemMap, err := ResourceIBMIsImageRecipeImageRecipeStatusReasonToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		err = fmt.Errorf("error setting status_reasons: %s", err)
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "read", "set-status_reasons").GetDiag()
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
		baseImage, err := ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatch(d.Get("base_image.0").(map[string]interface{}))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "update", "parse-base_image").GetDiag()
		}
		patchVals.BaseImage = baseImage
		hasChange = true
	}
	if d.HasChange("deprecation_at") {
		newDeprecationAt, err := core.ParseDateTime(d.Get("deprecation_at").(string))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "update", "parse-deprecation_at").GetDiag()
		}
		patchVals.DeprecationAt = &newDeprecationAt
		hasChange = true
	}
	if d.HasChange("image_components") {
		var imageComponents []vpcv1.ImageComponentIdentityIntf
		for _, v := range d.Get("image_components").([]interface{}) {
			value := v.(map[string]interface{})
			imageComponentsItem, err := ResourceIBMIsImageRecipeMapToImageComponentIdentity(value)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "update", "parse-image_components").GetDiag()
			}
			imageComponents = append(imageComponents, imageComponentsItem)
		}
		patchVals.ImageComponents = imageComponents
		hasChange = true
	}
	if d.HasChange("name") {
		newName := d.Get("name").(string)
		patchVals.Name = &newName
		hasChange = true
	}
	if d.HasChange("obsolescence_at") {
		newObsolescenceAt, err := core.ParseDateTime(d.Get("obsolescence_at").(string))
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_is_image_recipe", "update", "parse-obsolescence_at").GetDiag()
		}
		patchVals.ObsolescenceAt = &newObsolescenceAt
		hasChange = true
	}

	if hasChange {
		// Fields with `nil` values are omitted from the generic map,
		// so we need to re-add them to support removing arguments
		// in merge-patch operations sent to the service.
		updateImageRecipeOptions.ImageRecipePatch = ResourceIBMIsImageRecipeImageRecipePatchAsPatch(patchVals, d)

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

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototype(modelMap map[string]interface{}) (vpcv1.ImageRecipeBaseImagePrototypeIntf, error) {
	model := &vpcv1.ImageRecipeBaseImagePrototype{}
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

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByID(modelMap map[string]interface{}) (*vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByID, error) {
	model := &vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByCRN, error) {
	model := &vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePrototypeImageIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByHref, error) {
	model := &vpcv1.ImageRecipeBaseImagePrototypeImageIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageComponentIdentity(modelMap map[string]interface{}) (vpcv1.ImageComponentIdentityIntf, error) {
	model := &vpcv1.ImageComponentIdentity{}
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

func ResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByID(modelMap map[string]interface{}) (*vpcv1.ImageComponentIdentityImageComponentByID, error) {
	model := &vpcv1.ImageComponentIdentityImageComponentByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByCRN(modelMap map[string]interface{}) (*vpcv1.ImageComponentIdentityImageComponentByCRN, error) {
	model := &vpcv1.ImageComponentIdentityImageComponentByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageComponentIdentityImageComponentByHref(modelMap map[string]interface{}) (*vpcv1.ImageComponentIdentityImageComponentByHref, error) {
	model := &vpcv1.ImageComponentIdentityImageComponentByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageRecipeResourceGroupPrototype(modelMap map[string]interface{}) (vpcv1.ImageRecipeResourceGroupPrototypeIntf, error) {
	model := &vpcv1.ImageRecipeResourceGroupPrototype{}
	if modelMap["id"] != nil && modelMap["id"].(string) != "" {
		model.ID = core.StringPtr(modelMap["id"].(string))
	}
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageRecipeResourceGroupPrototypeResourceGroupIdentityByID(modelMap map[string]interface{}) (*vpcv1.ImageRecipeResourceGroupPrototypeResourceGroupIdentityByID, error) {
	model := &vpcv1.ImageRecipeResourceGroupPrototypeResourceGroupIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatch(modelMap map[string]interface{}) (vpcv1.ImageRecipeBaseImagePatchIntf, error) {
	model := &vpcv1.ImageRecipeBaseImagePatch{}
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

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByID(modelMap map[string]interface{}) (*vpcv1.ImageRecipeBaseImagePatchImageIdentityByID, error) {
	model := &vpcv1.ImageRecipeBaseImagePatchImageIdentityByID{}
	model.ID = core.StringPtr(modelMap["id"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByCRN(modelMap map[string]interface{}) (*vpcv1.ImageRecipeBaseImagePatchImageIdentityByCRN, error) {
	model := &vpcv1.ImageRecipeBaseImagePatchImageIdentityByCRN{}
	model.CRN = core.StringPtr(modelMap["crn"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeMapToImageRecipeBaseImagePatchImageIdentityByHref(modelMap map[string]interface{}) (*vpcv1.ImageRecipeBaseImagePatchImageIdentityByHref, error) {
	model := &vpcv1.ImageRecipeBaseImagePatchImageIdentityByHref{}
	model.Href = core.StringPtr(modelMap["href"].(string))
	return model, nil
}

func ResourceIBMIsImageRecipeImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsImageRecipeDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := ResourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsImageRecipeDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func ResourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := ResourceIBMIsImageRecipeAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := ResourceIBMIsImageRecipeRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func ResourceIBMIsImageRecipeAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func ResourceIBMIsImageRecipeRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsImageRecipeImageComponentReferenceToMap(model *vpcv1.ImageComponentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := ResourceIBMIsImageRecipeDeletedToMap(model.Deleted)
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

func ResourceIBMIsImageRecipeImageRecipeResourceGroupToMap(model *vpcv1.ImageRecipeResourceGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func ResourceIBMIsImageRecipeImageRecipeStatusReasonToMap(model *vpcv1.ImageRecipeStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}

func ResourceIBMIsImageRecipeImageRecipePatchAsPatch(patchVals *vpcv1.ImageRecipePatch, d *schema.ResourceData) map[string]interface{} {
	patch, _ := patchVals.AsPatch()
	var path string

	path = "base_image"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["base_image"] = nil
	} else if exists && patch["base_image"] != nil {
		ResourceIBMIsImageRecipeImageRecipeBaseImagePatchAsPatch(patch["base_image"].(map[string]interface{}), d, fmt.Sprintf("%s.0", path))
	} else if !exists {
		delete(patch, "base_image")
	}
	path = "deprecation_at"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["deprecation_at"] = nil
	} else if !exists {
		delete(patch, "deprecation_at")
	}
	path = "image_components"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["image_components"] = nil
	} else if exists && patch["image_components"] != nil {
		image_componentsList := patch["image_components"].([]map[string]interface{})
		for i, image_componentsItem := range image_componentsList {
			ResourceIBMIsImageRecipeImageComponentIdentityAsPatch(image_componentsItem, d, fmt.Sprintf("%s.%d", path, i))
		}
	} else if !exists {
		delete(patch, "image_components")
	}
	path = "name"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["name"] = nil
	} else if !exists {
		delete(patch, "name")
	}
	path = "obsolescence_at"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["obsolescence_at"] = nil
	} else if !exists {
		delete(patch, "obsolescence_at")
	}

	return patch
}

func ResourceIBMIsImageRecipeImageComponentIdentityAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["id"] = nil
	} else if !exists {
		delete(patch, "id")
	}
	path = rootPath + ".crn"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["crn"] = nil
	} else if !exists {
		delete(patch, "crn")
	}
	path = rootPath + ".href"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["href"] = nil
	} else if !exists {
		delete(patch, "href")
	}
}

func ResourceIBMIsImageRecipeImageRecipeBaseImagePatchAsPatch(patch map[string]interface{}, d *schema.ResourceData, rootPath string) {
	var path string

	path = rootPath + ".id"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["id"] = nil
	} else if !exists {
		delete(patch, "id")
	}
	path = rootPath + ".crn"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["crn"] = nil
	} else if !exists {
		delete(patch, "crn")
	}
	path = rootPath + ".href"
	if _, exists := d.GetOk(path); d.HasChange(path) && !exists {
		patch["href"] = nil
	} else if !exists {
		delete(patch, "href")
	}
}
