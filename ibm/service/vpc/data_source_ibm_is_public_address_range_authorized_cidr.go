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

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDR() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPublicAddressRangeAuthorizedCIDRRead,

		Schema: map[string]*schema.Schema{
			"is_public_address_range_authorized_cidr_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public address range authorized CIDR identifier.",
			},
			"allocation": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": &schema.Schema{
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of resources allocated from this public address range authorized CIDR.",
						},
						"profile_family": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The profile `family` for resources allocated from this public address range authorized CIDR.- `provider`: The resources allocated from this authorized CIDR will have a profile  with a `family` value of `provider`.- `user`: The resources allocated from this authorized CIDR will have a profile with  a `family` value of `user`.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
						},
					},
				},
			},
			"availability_mode": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The availability mode of the public address range authorized CIDR:- `regional`: Resources allocated from the authorized CIDR can reside in any zone in the  region.- `zonal`: Resources allocated from the authorized CIDR must reside in the authorized  CIDR's `zone`.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IP address block for the public address range authorized CIDR, expressed in CIDR format.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 address blocks in the future.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this public address range authorized CIDR.",
			},
			"ip_version": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The IP version for this public address range authorized CIDR:- `ipv4`: An IPv4 public address range authorized CIDR.The enumerated values for this property may[expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) in the future.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this public address range authorized CIDR. The name is unique across all public address range authorized CIDRs in the region.",
			},
			"network_prefix_length": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The network prefix length for this public address range authorized CIDR.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
			"zone": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The zone for this public address range authorized CIDR. Resources allocated from thisauthorized CIDR must reside in this zone.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this zone.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this zone.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsPublicAddressRangeAuthorizedCIDRRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getPublicAddressRangeAuthorizedCIDROptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDROptions{}

	getPublicAddressRangeAuthorizedCIDROptions.SetID(d.Get("is_public_address_range_authorized_cidr_id").(string))

	publicAddressRangeAuthorizedCIDR, _, err := vpcClient.GetPublicAddressRangeAuthorizedCIDRWithContext(context, getPublicAddressRangeAuthorizedCIDROptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetPublicAddressRangeAuthorizedCIDRWithContext failed: %s", err.Error()), "(Data) ibm_is_public_address_range_authorized_cidr", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(*publicAddressRangeAuthorizedCIDR.ID)

	allocation := []map[string]interface{}{}
	allocationMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCIDRPublicAddressRangeAuthorizedCIDRAllocationToMap(publicAddressRangeAuthorizedCIDR.Allocation)
	if err != nil {
		return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "allocation-to-map").GetDiag()
	}
	allocation = append(allocation, allocationMap)
	if err = d.Set("allocation", allocation); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting allocation: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-allocation").GetDiag()
	}

	if err = d.Set("availability_mode", publicAddressRangeAuthorizedCIDR.AvailabilityMode); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting availability_mode: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-availability_mode").GetDiag()
	}

	if err = d.Set("cidr", publicAddressRangeAuthorizedCIDR.CIDR); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cidr: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-cidr").GetDiag()
	}

	if err = d.Set("href", publicAddressRangeAuthorizedCIDR.Href); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-href").GetDiag()
	}

	if err = d.Set("ip_version", publicAddressRangeAuthorizedCIDR.IPVersion); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting ip_version: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-ip_version").GetDiag()
	}

	if err = d.Set("name", publicAddressRangeAuthorizedCIDR.Name); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-name").GetDiag()
	}

	if err = d.Set("network_prefix_length", flex.IntValue(publicAddressRangeAuthorizedCIDR.NetworkPrefixLength)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting network_prefix_length: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-network_prefix_length").GetDiag()
	}

	if err = d.Set("resource_type", publicAddressRangeAuthorizedCIDR.ResourceType); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-resource_type").GetDiag()
	}

	if !core.IsNil(publicAddressRangeAuthorizedCIDR.Zone) {
		zone := []map[string]interface{}{}
		zoneMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCIDRZoneReferenceToMap(publicAddressRangeAuthorizedCIDR.Zone)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "zone-to-map").GetDiag()
		}
		zone = append(zone, zoneMap)
		if err = d.Set("zone", zone); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting zone: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr", "read", "set-zone").GetDiag()
		}
	}

	return nil
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRPublicAddressRangeAuthorizedCIDRAllocationToMap(model *vpcv1.PublicAddressRangeAuthorizedCIDRAllocation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["count"] = flex.IntValue(model.Count)
	modelMap["profile_family"] = *model.ProfileFamily
	return modelMap, nil
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
