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

func DataSourceIBMIsImageComponent() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImageComponentRead,

		Schema: map[string]*schema.Schema{
			"is_image_component_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image component identifier.",
			},
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
	}
}

func dataSourceIBMIsImageComponentRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_component", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getImageComponentsOptions := &vpcv1.GetImageComponentsOptions{}

	getImageComponentsOptions.SetID(d.Get("is_image_component_id").(string))

	imageComponent, _, err := vpcClient.GetImageComponentsWithContext(context, getImageComponentsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageComponentsWithContext failed: %s", err.Error()), "(Data) ibm_is_image_component", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*imageComponent.ID)

	allowedUse := []map[string]interface{}{}
	allowedUseMap, err := DataSourceIBMIsImageComponentImageComponentAllowedUseToMap(imageComponent.AllowedUse)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_component", "read", "allowed_use-to-map").GetDiag()
	}
	allowedUse = append(allowedUse, allowedUseMap)
	if err = d.Set("allowed_use", allowedUse); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting allowed_use: %s", err), "(Data) ibm_is_image_component", "read", "set-allowed_use").GetDiag()
	}

	if err = d.Set("created_at", flex.DateTimeToString(imageComponent.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting created_at: %s", err), "(Data) ibm_is_image_component", "read", "set-created_at").GetDiag()
	}

	if err = d.Set("crn", imageComponent.CRN); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting crn: %s", err), "(Data) ibm_is_image_component", "read", "set-crn").GetDiag()
	}

	dependencies := []map[string]interface{}{}
	for _, dependenciesItem := range imageComponent.Dependencies {
		dependenciesItemMap, err := DataSourceIBMIsImageComponentImageComponentReferenceToMap(&dependenciesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_component", "read", "dependencies-to-map").GetDiag()
		}
		dependencies = append(dependencies, dependenciesItemMap)
	}
	if err = d.Set("dependencies", dependencies); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting dependencies: %s", err), "(Data) ibm_is_image_component", "read", "set-dependencies").GetDiag()
	}

	if !core.IsNil(imageComponent.DeprecationAt) {
		if err = d.Set("deprecation_at", flex.DateTimeToString(imageComponent.DeprecationAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting deprecation_at: %s", err), "(Data) ibm_is_image_component", "read", "set-deprecation_at").GetDiag()
		}
	}

	if err = d.Set("href", imageComponent.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting href: %s", err), "(Data) ibm_is_image_component", "read", "set-href").GetDiag()
	}

	if err = d.Set("name", imageComponent.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting name: %s", err), "(Data) ibm_is_image_component", "read", "set-name").GetDiag()
	}

	if !core.IsNil(imageComponent.ObsolescenceAt) {
		if err = d.Set("obsolescence_at", flex.DateTimeToString(imageComponent.ObsolescenceAt)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting obsolescence_at: %s", err), "(Data) ibm_is_image_component", "read", "set-obsolescence_at").GetDiag()
		}
	}

	resourceGroup := []map[string]interface{}{}
	resourceGroupMap, err := DataSourceIBMIsImageComponentImageComponentResourceGroupToMap(imageComponent.ResourceGroup)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_component", "read", "resource_group-to-map").GetDiag()
	}
	resourceGroup = append(resourceGroup, resourceGroupMap)
	if err = d.Set("resource_group", resourceGroup); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting resource_group: %s", err), "(Data) ibm_is_image_component", "read", "set-resource_group").GetDiag()
	}

	if err = d.Set("resource_type", imageComponent.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting resource_type: %s", err), "(Data) ibm_is_image_component", "read", "set-resource_type").GetDiag()
	}

	software := []map[string]interface{}{}
	softwareMap, err := DataSourceIBMIsImageComponentImageComponentSoftwareToMap(imageComponent.Software)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_component", "read", "software-to-map").GetDiag()
	}
	software = append(software, softwareMap)
	if err = d.Set("software", software); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting software: %s", err), "(Data) ibm_is_image_component", "read", "set-software").GetDiag()
	}

	if err = d.Set("status", imageComponent.Status); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting status: %s", err), "(Data) ibm_is_image_component", "read", "set-status").GetDiag()
	}

	statusReasons := []map[string]interface{}{}
	for _, statusReasonsItem := range imageComponent.StatusReasons {
		statusReasonsItemMap, err := DataSourceIBMIsImageComponentImageComponentStatusReasonToMap(&statusReasonsItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_component", "read", "status_reasons-to-map").GetDiag()
		}
		statusReasons = append(statusReasons, statusReasonsItemMap)
	}
	if err = d.Set("status_reasons", statusReasons); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting status_reasons: %s", err), "(Data) ibm_is_image_component", "read", "set-status_reasons").GetDiag()
	}

	if err = d.Set("visibility", imageComponent.Visibility); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting visibility: %s", err), "(Data) ibm_is_image_component", "read", "set-visibility").GetDiag()
	}

	return nil
}

func DataSourceIBMIsImageComponentImageComponentAllowedUseToMap(model *vpcv1.ImageComponentAllowedUse) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["api_version"] = *model.ApiVersion
	modelMap["image"] = *model.Image
	return modelMap, nil
}

func DataSourceIBMIsImageComponentImageComponentReferenceToMap(model *vpcv1.ImageComponentReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsImageComponentDeletedToMap(model.Deleted)
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

func DataSourceIBMIsImageComponentDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsImageComponentImageComponentResourceGroupToMap(model *vpcv1.ImageComponentResourceGroup) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["name"] = *model.Name
	return modelMap, nil
}

func DataSourceIBMIsImageComponentImageComponentSoftwareToMap(model vpcv1.ImageComponentSoftwareIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.ImageComponentSoftwarePackage); ok {
		return DataSourceIBMIsImageComponentImageComponentSoftwarePackageToMap(model.(*vpcv1.ImageComponentSoftwarePackage))
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

func DataSourceIBMIsImageComponentImageComponentSoftwarePackageToMap(model *vpcv1.ImageComponentSoftwarePackage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	modelMap["version"] = *model.Version
	return modelMap, nil
}

func DataSourceIBMIsImageComponentImageComponentStatusReasonToMap(model *vpcv1.ImageComponentStatusReason) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["code"] = *model.Code
	modelMap["message"] = *model.Message
	if model.MoreInfo != nil {
		modelMap["more_info"] = *model.MoreInfo
	}
	return modelMap, nil
}
