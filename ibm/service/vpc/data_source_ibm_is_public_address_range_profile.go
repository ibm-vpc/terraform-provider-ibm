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

func DataSourceIBMIsPublicAddressRangeProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPublicAddressRangeProfileRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public address range profile name.",
			},
			"family": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The product family this public address range profile belongs to.- `provider`: The public IP addresses in the public address range with this profile are   owned by the provider.- `user`: The public IP addresses in the public address range with this profile are   owned by the user.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this public address range profile.",
			},
			"ip_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP version for public address ranges with this profile:- `ipv4`: An IPv4 public address range.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIBMIsPublicAddressRangeProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_profile", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getPublicAddressRangeProfileOptions := &vpcv1.GetPublicAddressRangeProfileOptions{}

	getPublicAddressRangeProfileOptions.SetName(d.Get("name").(string))

	publicAddressRangeProfile, _, err := vpcClient.GetPublicAddressRangeProfileWithContext(context, getPublicAddressRangeProfileOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetPublicAddressRangeProfileWithContext failed: %s", err.Error()), "(Data) ibm_is_public_address_range_profile", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsPublicAddressRangeProfileID(d))

	if err = d.Set("family", publicAddressRangeProfile.Family); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting family: %s", err), "(Data) ibm_is_public_address_range_profile", "read", "set-family").GetDiag()
	}

	if err = d.Set("href", publicAddressRangeProfile.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_public_address_range_profile", "read", "set-href").GetDiag()
	}

	if err = d.Set("ip_version", publicAddressRangeProfile.IPVersion); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ip_version: %s", err), "(Data) ibm_is_public_address_range_profile", "read", "set-ip_version").GetDiag()
	}

	if err = d.Set("resource_type", publicAddressRangeProfile.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_public_address_range_profile", "read", "set-resource_type").GetDiag()
	}

	return nil
}

// dataSourceIBMIsPublicAddressRangeProfileID returns a reasonable ID for the list.
func dataSourceIBMIsPublicAddressRangeProfileID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
