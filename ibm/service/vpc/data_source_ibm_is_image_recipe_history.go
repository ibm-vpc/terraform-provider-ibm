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

func DataSourceIBMIsImageRecipeHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImageRecipeHistoryRead,

		Schema: map[string]*schema.Schema{
			"is_image_recipe_history_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image recipe identifier.",
			},
			"history": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of image recipes ordered by `created_at`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
							Description: "The date and time that this specific version of the image recipe was created.",
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this image recipe.",
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
							Description: "The date and time that this specific version of the image recipe was superseded.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsImageRecipeHistoryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe_history", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listImageRecipeHistoryOptions := &vpcv1.ListImageRecipeHistoryOptions{}

	listImageRecipeHistoryOptions.SetID(d.Get("is_image_recipe_history_id").(string))

	var pager *vpcv1.ImageRecipeHistoryPager
	pager, err = vpcClient.NewImageRecipeHistoryPager(listImageRecipeHistoryOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe_history", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ImageRecipeHistoryPager.GetAll() failed %s", err), "(Data) ibm_is_image_recipe_history", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsImageRecipeHistoryID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsImageRecipeHistoryImageRecipeHistoryToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_recipe_history", "read", "ImageRecipes-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("history", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting history %s", err), "(Data) ibm_is_image_recipe_history", "read", "history-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsImageRecipeHistoryID returns a reasonable ID for the list.
func dataSourceIBMIsImageRecipeHistoryID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsImageRecipeHistoryImageRecipeHistoryToMap(model *vpcv1.ImageRecipeHistory) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	baseImageMap, err := DataSourceIBMIsImageRecipeHistoryImageReferenceToMap(model.BaseImage)
	if err != nil {
		return modelMap, err
	}
	modelMap["base_image"] = []map[string]interface{}{baseImageMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["crn"] = *model.CRN
	if model.DeprecationAt != nil {
		modelMap["deprecation_at"] = model.DeprecationAt.String()
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	imageComponents := []map[string]interface{}{}
	for _, imageComponentsItem := range model.ImageComponents {
		imageComponentsItemMap, err := DataSourceIBMIsImageRecipeHistoryImageComponentReferenceToMap(&imageComponentsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		imageComponents = append(imageComponents, imageComponentsItemMap)
	}
	modelMap["image_components"] = imageComponents
	modelMap["name"] = *model.Name
	if model.ObsolescenceAt != nil {
		modelMap["obsolescence_at"] = model.ObsolescenceAt.String()
	}
	resourceGroupMap, err := DataSourceIBMIsImageRecipeHistoryImageRecipeResourceGroupToMap(model.ResourceGroup)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	modelMap["resource_type"] = *model.ResourceType
	modelMap["status"] = *model.Status
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range model.StatusReasons {
		statusReasonsItemMap, err := DataSourceIBMIsImageRecipeHistoryImageRecipeStatusReasonToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	modelMap["status_reasons"] = statusReasons
	if model.SupersededAt != nil {
		modelMap["superseded_at"] = model.SupersededAt.String()
	}
	modelMap["updated_at"] = model.UpdatedAt.String()
	modelMap["version"] = *model.Version
	modelMap["visibility"] = *model.Visibility
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeHistoryImageReferenceToMap(model *vpcv1.ImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageRecipeHistoryDeletedToMap(model.Deleted)
		if err != nil {
			return modelMap, err
		}
		modelMap["deleted"] = []map[string]interface{}{deletedMap}
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.Remote != nil {
		remoteMap, err := DataSourceIBMIsImageRecipeHistoryImageRemoteContextImageReferenceToMap(model.Remote)
		if err != nil {
			return modelMap, err
		}
		modelMap["remote"] = []map[string]interface{}{remoteMap}
	}
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeHistoryDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeHistoryImageRemoteContextImageReferenceToMap(model *vpcv1.ImageRemoteContextImageReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.Account != nil {
		accountMap, err := DataSourceIBMIsImageRecipeHistoryAccountReferenceToMap(model.Account)
		if err != nil {
			return modelMap, err
		}
		modelMap["account"] = []map[string]interface{}{accountMap}
	}
	if model.Region != nil {
		regionMap, err := DataSourceIBMIsImageRecipeHistoryRegionReferenceToMap(model.Region)
		if err != nil {
			return modelMap, err
		}
		modelMap["region"] = []map[string]interface{}{regionMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeHistoryAccountReferenceToMap(model *vpcv1.AccountReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["id"] = *model.ID
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeHistoryRegionReferenceToMap(model *vpcv1.RegionReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeHistoryImageComponentReferenceToMap(model *vpcv1.ImageComponentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageRecipeHistoryDeletedToMap(model.Deleted)
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

func DataSourceIBMIsImageRecipeHistoryImageRecipeResourceGroupToMap(model *vpcv1.ImageRecipeResourceGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageRecipeHistoryImageRecipeStatusReasonToMap(model *vpcv1.ImageRecipeStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
