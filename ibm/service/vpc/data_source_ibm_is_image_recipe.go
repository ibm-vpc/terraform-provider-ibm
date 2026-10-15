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

func DataSourceIBMIsImageRecipe() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImageRecipeRead,

		Schema: map[string]*schema.Schema{
			"is_image_recipe_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image recipe identifier.",
			},
			"base_image": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The base image that is extended when building an image from this image recipe.",
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
			"deprecation_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The deprecation date and time (UTC) for this image recipe.If absent, no deprecation date and time has been set.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this image recipe.",
			},
			"image_components": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The image components included in this image recipe.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this component.",
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
							Description: "The URL for this component.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this component.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique name for this image component.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique name for this image recipe.",
			},
			"obsolescence_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The obsolescence date and time (UTC) for this image recipe.If absent, no obsolescence date and time has been set.",
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
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
							Computed:    true,
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

func dataSourceIBMIsImageRecipeRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getImageRecipeOptions := &vpcv1.GetImageRecipeOptions{}

	getImageRecipeOptions.SetID(d.Get("is_image_recipe_id").(string))

	imageRecipe, _, err := vpcClient.GetImageRecipeWithContext(context, getImageRecipeOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageRecipeWithContext failed: %s", err.Error()), "(Data) ibm_is_image_recipe", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*getImageRecipeOptions.ID)

	baseImage := []map[string]interface{}{}
	baseImageMap, err := DataSourceIBMIsImageRecipeImageReferenceToMap(imageRecipe.BaseImage)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe", "read", "base_image-to-map").GetDiag()
	}
	baseImage = append(baseImage, baseImageMap)
	if err = d.Set("base_image", baseImage); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting base_image: %s", err), "(Data) ibm_is_image_recipe", "read", "set-base_image").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(imageRecipe.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting created_at: %s", err), "(Data) ibm_is_image_recipe", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("crn", imageRecipe.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting crn: %s", err), "(Data) ibm_is_image_recipe", "read", "set-crn").GetDiag()
	}

	if !core.IsNil(imageRecipe.DeprecationAt) {
		if err = d.Set("deprecation_at", flex.DateTimeToString(imageRecipe.DeprecationAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting deprecation_at: %s", err), "(Data) ibm_is_image_recipe", "read", "set-deprecation_at").GetDiag()
		}
	}

	if err = d.Set("href", imageRecipe.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting href: %s", err), "(Data) ibm_is_image_recipe", "read", "set-href").GetDiag()
	}

	imageComponents := []map[string]interface{}{}
	for _, imageComponentsItem := range imageRecipe.ImageComponents {
		imageComponentsItemMap, err := DataSourceIBMIsImageRecipeImageComponentReferenceToMap(&imageComponentsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe", "read", "image_components-to-map").GetDiag()
		}
		imageComponents = append(imageComponents, imageComponentsItemMap)
	}
	if err = d.Set("image_components", imageComponents); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting image_components: %s", err), "(Data) ibm_is_image_recipe", "read", "set-image_components").GetDiag()
	}

	if err = d.Set("name", imageRecipe.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting name: %s", err), "(Data) ibm_is_image_recipe", "read", "set-name").GetDiag()
	}

	if !core.IsNil(imageRecipe.ObsolescenceAt) {
		if err = d.Set("obsolescence_at", flex.DateTimeToString(imageRecipe.ObsolescenceAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting obsolescence_at: %s", err), "(Data) ibm_is_image_recipe", "read", "set-obsolescence_at").GetDiag()
		}
	}

	resourceGroup := []map[string]interface{}{}
	resourceGroupMap, err := DataSourceIBMIsImageRecipeImageRecipeResourceGroupToMap(imageRecipe.ResourceGroup)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe", "read", "resource_group-to-map").GetDiag()
	}
	resourceGroup = append(resourceGroup, resourceGroupMap)
	if err = d.Set("resource_group", resourceGroup); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting resource_group: %s", err), "(Data) ibm_is_image_recipe", "read", "set-resource_group").GetDiag()
	}

	if err = d.Set("resource_type", imageRecipe.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting resource_type: %s", err), "(Data) ibm_is_image_recipe", "read", "set-resource_type").GetDiag()
	}

	if err = d.Set("status", imageRecipe.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting status: %s", err), "(Data) ibm_is_image_recipe", "read", "set-status").GetDiag()
	}

	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range imageRecipe.StatusReasons {
		statusReasonsItemMap, err := DataSourceIBMIsImageRecipeImageRecipeStatusReasonToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting status_reasons: %s", err), "(Data) ibm_is_image_recipe", "read", "set-status_reasons").GetDiag()
	}

	if !core.IsNil(imageRecipe.SupersededAt) {
		if err = d.Set("superseded_at", flex.DateTimeToString(imageRecipe.SupersededAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting superseded_at: %s", err), "(Data) ibm_is_image_recipe", "read", "set-superseded_at").GetDiag()
		}
	}

	if err = d.Set("updated_at", flex.DateTimeToString(imageRecipe.UpdatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting updated_at: %s", err), "(Data) ibm_is_image_recipe", "read", "set-updated_at").GetDiag()
	}

	if err = d.Set("version", imageRecipe.Version); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting version: %s", err), "(Data) ibm_is_image_recipe", "read", "set-version").GetDiag()
	}

	if err = d.Set("visibility", imageRecipe.Visibility); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting visibility: %s", err), "(Data) ibm_is_image_recipe", "read", "set-visibility").GetDiag()
	}

	return nil
}

func DataSourceIBMIsImageRecipeImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageRecipeDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := DataSourceIBMIsImageRecipeAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsImageRecipeRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeImageComponentReferenceToMap(model *vpcv1.ImageComponentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageRecipeDeletedToMap(model.Deleted)
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

func DataSourceIBMIsImageRecipeImageRecipeResourceGroupToMap(model *vpcv1.ImageRecipeResourceGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeImageRecipeStatusReasonToMap(model *vpcv1.ImageRecipeStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
