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

func DataSourceIBMIsFloatingIPProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsFloatingIPProfileRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The floating IP profile name.",
			},
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
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIBMIsFloatingIPProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_floating_ip_profile", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getFloatingIPProfileOptions := &vpcv1.GetFloatingIPProfileOptions{}

	getFloatingIPProfileOptions.SetName(d.Get("name").(string))

	floatingIPProfile, _, err := vpcClient.GetFloatingIPProfileWithContext(context, getFloatingIPProfileOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetFloatingIPProfileWithContext failed: %s", err.Error()), "(Data) ibm_is_floating_ip_profile", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsFloatingIPProfileID(d))

	if err = d.Set("family", floatingIPProfile.Family); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting family: %s", err), "(Data) ibm_is_floating_ip_profile", "read", "set-family").GetDiag()
	}

	if err = d.Set("href", floatingIPProfile.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_floating_ip_profile", "read", "set-href").GetDiag()
	}

	if err = d.Set("ip_version", floatingIPProfile.IPVersion); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ip_version: %s", err), "(Data) ibm_is_floating_ip_profile", "read", "set-ip_version").GetDiag()
	}

	if err = d.Set("resource_type", floatingIPProfile.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_floating_ip_profile", "read", "set-resource_type").GetDiag()
	}

	return nil
}

// dataSourceIBMIsFloatingIPProfileID returns a reasonable ID for the list.
func dataSourceIBMIsFloatingIPProfileID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
