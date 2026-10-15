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

func DataSourceIBMIsImageInventory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImageInventoryRead,

		Schema: map[string]*schema.Schema{
			"is_image_inventory_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image identifier.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when image inventory was created.",
			},
			"packages": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The package name.",
						},
						"version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The package version.",
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
	}
}

func dataSourceIBMIsImageInventoryRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_inventory", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getImageInventoryOptions := &vpcv1.GetImageInventoryOptions{}

	getImageInventoryOptions.SetID(d.Get("is_image_inventory_id").(string))

	imageInventory, _, err := vpcClient.GetImageInventoryWithContext(context, getImageInventoryOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetImageInventoryWithContext failed: %s", err.Error()), "(Data) ibm_is_image_inventory", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsImageInventoryID(d))

	if err = d.Set("created_at", flex.DateTimeToString(imageInventory.CreatedAt)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting created_at: %s", err), "(Data) ibm_is_image_inventory", "read", "set-created_at").GetDiag()
	}

	packages := []map[string]interface{}{}
	for _, packagesItem := range imageInventory.Packages {
		packagesItemMap, err := DataSourceIBMIsImageInventoryImageInventoryPackageToMap(&packagesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_inventory", "read", "packages-to-map").GetDiag()
		}
		packages = append(packages, packagesItemMap)
	}
	if err = d.Set("packages", packages); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting packages: %s", err), "(Data) ibm_is_image_inventory", "read", "set-packages").GetDiag()
	}

	if err = d.Set("resource_type", imageInventory.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting resource_type: %s", err), "(Data) ibm_is_image_inventory", "read", "set-resource_type").GetDiag()
	}

	return nil
}

// dataSourceIBMIsImageInventoryID returns a reasonable ID for the list.
func dataSourceIBMIsImageInventoryID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsImageInventoryImageInventoryPackageToMap(model *vpcv1.ImageInventoryPackage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["version"] = *model.Version
	return modelMap, nil
}
