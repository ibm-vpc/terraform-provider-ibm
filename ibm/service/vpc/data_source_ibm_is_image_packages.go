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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsImagePackages() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsImagePackagesRead,

		Schema: map[string]*schema.Schema{
			"image_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The image identifier.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to packages with a `name` property matching the exact specified name.",
			},
			"packages": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of packages installed in the image.",
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
		},
	}
}

func dataSourceIBMIsImagePackagesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_packages", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	imageID := d.Get("image_id").(string)

	listImagePackagesOptions := &vpcv1.ListImagePackagesOptions{}
	listImagePackagesOptions.SetID(imageID)
	if _, ok := d.GetOk("name"); ok {
		listImagePackagesOptions.SetName(d.Get("name").(string))
	}

	pager, err := vpcClient.NewImagePackagesPager(listImagePackagesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_image_packages", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAllWithContext(context)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ImagePackagesPager.GetAll() failed: %s", err), "(Data) ibm_is_image_packages", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(imageID)

	packages := []map[string]interface{}{}
	for _, packagesItem := range allItems {
		packagesItemMap, err := DataSourceIBMIsImagePackagesImagePackageToMap(&packagesItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_image_packages", "read", "packages-to-map").GetDiag()
		}
		packages = append(packages, packagesItemMap)
	}
	if err = d.Set("packages", packages); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("error setting packages: %s", err), "(Data) ibm_is_image_packages", "read", "set-packages").GetDiag()
	}

	return nil
}

func DataSourceIBMIsImagePackagesImagePackageToMap(model *vpcv1.ImagePackage) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["name"] = *model.Name
	modelMap["version"] = *model.Version
	return modelMap, nil
}
