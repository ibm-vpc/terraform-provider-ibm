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

func DataSourceIBMIsImageComponents() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImageComponentsRead,

		Schema: map[string]*schema.Schema{
			"is_image_components_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image identifier.",
			},
			"image_components": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of image components.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_use": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"api_version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API version with which to evaluate the expressions.",
									},
									"image": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The expression that must be satisfied by the properties of an image that uses this image component.The expression follows [Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md), but does not support built-in functions and macros. In addition, the following variables are supported, corresponding to `Image` properties:- `operating_system.architecture` (string): The operating system supported system  architecture- `operating_system.family` (string): The operating system family- `operating_system.version` (string): The operating system version.",
									},
								},
							},
						},
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the image component was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this component.",
						},
						"dependencies": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The dependencies for this image component. An image recipe that uses this component must contain its required dependencies.",
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
						"deprecation_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The deprecation date and time for this image component.If absent, no deprecation date and time has been set.",
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
						"obsolescence_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The obsolescence date and time for this component.If absent, no obsolescence date and time has been set.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this component.",
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
						"software": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The name of the software package.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The resource type.",
									},
									"version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The version of the software package.",
									},
								},
							},
						},
						"status": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of this image component- `available`: image component is ready to be used.- `degraded`: image component is degraded, and cannot be used.- `deleting`: image component is being deleted.- `deprecated`: image component is deprecated, but can still be used.- `obsolete`: image component is obsolete, and can no longer be used.- `pending`: image component is being created and is not ready to be used.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
										Description: "A reason code for this status:- `dependent_image_component_unusable`: One or more dependent image components  referenced by this image component have a status of `degraded`, `obsolete`, or  `deleting` or have been deleted. Image recipe creation with this image component  will fail until the image component no longer has unusable dependent image components.- `internal_error`: internal error (contact IBM support)The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
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
						"visibility": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The visibility of this image component.- `private`: Visible only to this account- `public`: Visible to all accounts.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsImageComponentsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_components", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listImageImageComponentsOptions := &vpcv1.ListImageImageComponentsOptions{}

	listImageImageComponentsOptions.SetID(d.Get("is_image_components_id").(string))

	var pager *vpcv1.ImageImageComponentsPager
	pager, err = vpcClient.NewImageImageComponentsPager(listImageImageComponentsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_image_components", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ImageImageComponentsPager.GetAll() failed %s", err), "(Data) ibm_is_image_components", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsImageComponentsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsImageComponentsImageComponentToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_components", "read", "Images-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("image_components", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting image_components %s", err), "(Data) ibm_is_image_components", "read", "image_components-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsImageComponentsID returns a reasonable ID for the list.
func dataSourceIBMIsImageComponentsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsImageComponentsImageComponentToMap(model *vpcv1.ImageComponent) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	allowedUseMap, err := DataSourceIBMIsImageComponentsImageComponentAllowedUseToMap(model.AllowedUse)
	if err != nil {
		return modelMap, err
	}
	modelMap["allowed_use"] = []map[string]interface{}{allowedUseMap}
	modelMap["created_at"] = model.CreatedAt.String()
	modelMap["crn"] = *model.CRN
	dependencies := []map[string]interface{}{}
	for _, dependenciesItem := range model.Dependencies {
		dependenciesItemMap, err := DataSourceIBMIsImageComponentsImageComponentReferenceToMap(&dependenciesItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		dependencies = append(dependencies, dependenciesItemMap)
	}
	modelMap["dependencies"] = dependencies
	if model.DeprecationAt != nil {
		modelMap["deprecation_at"] = model.DeprecationAt.String()
	}
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	if model.ObsolescenceAt != nil {
		modelMap["obsolescence_at"] = model.ObsolescenceAt.String()
	}
	resourceGroupMap, err := DataSourceIBMIsImageComponentsImageComponentResourceGroupToMap(model.ResourceGroup)
	if err != nil {
		return modelMap, err
	}
	modelMap["resource_group"] = []map[string]interface{}{resourceGroupMap}
	modelMap["resource_type"] = *model.ResourceType
	softwareMap, err := DataSourceIBMIsImageComponentsImageComponentSoftwareToMap(model.Software)
	if err != nil {
		return modelMap, err
	}
	modelMap["software"] = []map[string]interface{}{softwareMap}
	modelMap["status"] = *model.Status
	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range model.StatusReasons {
		statusReasonsItemMap, err := DataSourceIBMIsImageComponentsImageComponentStatusReasonToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return modelMap, err
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	modelMap["status_reasons"] = statusReasons
	modelMap["visibility"] = *model.Visibility
	return modelMap, nil
}

func DataSourceIBMIsImageComponentsImageComponentAllowedUseToMap(model *vpcv1.ImageComponentAllowedUse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["api_version"] = *model.ApiVersion
	modelMap["image"] = *model.Image
	return modelMap, nil
}

func DataSourceIBMIsImageComponentsImageComponentReferenceToMap(model *vpcv1.ImageComponentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageComponentsDeletedToMap(model.Deleted)
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

func DataSourceIBMIsImageComponentsDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsImageComponentsImageComponentResourceGroupToMap(model *vpcv1.ImageComponentResourceGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageComponentsImageComponentSoftwareToMap(model vpcv1.ImageComponentSoftwareIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ImageComponentSoftwarePackage); ok {
		return DataSourceIBMIsImageComponentsImageComponentSoftwarePackageToMap(model.(*vpcv1.ImageComponentSoftwarePackage))
	} else if _, ok := model.(*vpcv1.ImageComponentSoftware); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.ImageComponentSoftware)
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		if model.Version != nil {
			modelMap["version"] = *model.Version
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("unrecognized vpcv1.ImageComponentSoftwareIntf subtype encountered")
	}
}

func DataSourceIBMIsImageComponentsImageComponentSoftwarePackageToMap(model *vpcv1.ImageComponentSoftwarePackage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	modelMap["version"] = *model.Version
	return modelMap, nil
}

func DataSourceIBMIsImageComponentsImageComponentStatusReasonToMap(model *vpcv1.ImageComponentStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
