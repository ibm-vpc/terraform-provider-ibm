// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.111.0-1bfb72c2-20260206-185521
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

func DataSourceIBMIsFloatingIPProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsFloatingIPProfilesRead,

		Schema: map[string]*schema.Schema{
			"profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of floating IP profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this floating IP profile belongs to.- `provider`: The floating IP with this profile is owned by the provider.- `user`: The floating IP with this profile is owned by the user.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this floating IP profile.",
						},
						"ip_version": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The IP version for floating IPs with this profile:- `ipv4`: An IPv4 floating IP.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this floating IP profile.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsFloatingIPProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_floating_ip_profiles", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listFloatingIPProfilesOptions := &vpcv1.ListFloatingIPProfilesOptions{}

	var pager *vpcv1.FloatingIPProfilesPager
	pager, err = vpcClient.NewFloatingIPProfilesPager(listFloatingIPProfilesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_floating_ip_profiles", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("FloatingIPProfilesPager.GetAll() failed %s", err), "(Data) ibm_is_floating_ip_profiles", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsFloatingIPProfilesID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsFloatingIPProfilesFloatingIPProfileToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_floating_ip_profiles", "read", "FloatingIPs-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("profiles", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting profiles %s", err), "(Data) ibm_is_floating_ip_profiles", "read", "profiles-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsFloatingIPProfilesID returns a reasonable ID for the list.
func dataSourceIBMIsFloatingIPProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsFloatingIPProfilesFloatingIPProfileToMap(model *vpcv1.FloatingIPProfile) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["family"] = *model.Family
	modelMap["href"] = *model.Href
	modelMap["ip_version"] = *model.IPVersion
	modelMap["name"] = *model.Name
	modelMap["resource_type"] = *model.ResourceType
	return modelMap, nil
}
