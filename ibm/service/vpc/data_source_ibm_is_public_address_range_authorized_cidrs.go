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

func DataSourceIBMIsPublicAddressRangeAuthorizedCidrs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPublicAddressRangeAuthorizedCidrsRead,

		Schema: map[string]*schema.Schema{
			"allocation_profile_family": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with an `allocation.profile_family` property matching the exact specified value.",
			},
			"availability_mode": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with an `availability_mode` property matching the exact specified value.",
			},
			"authorized_cidrs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A page of public address range authorized CIDRs.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this public address range authorized CIDR.",
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
				},
			},
		},
	}
}

func dataSourceIBMIsPublicAddressRangeAuthorizedCidrsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidrs", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listPublicAddressRangeAuthorizedCidrsOptions := &vpcv1.ListPublicAddressRangeAuthorizedCIDRsOptions{}

	if _, ok := d.GetOk("allocation_profile_family"); ok {
		listPublicAddressRangeAuthorizedCidrsOptions.SetAllocationProfileFamily(d.Get("allocation_profile_family").(string))
	}
	if _, ok := d.GetOk("availability_mode"); ok {
		listPublicAddressRangeAuthorizedCidrsOptions.SetAvailabilityMode(d.Get("availability_mode").(string))
	}

	var pager *vpcv1.PublicAddressRangeAuthorizedCIDRsPager
	pager, err = vpcClient.NewPublicAddressRangeAuthorizedCIDRsPager(listPublicAddressRangeAuthorizedCidrsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidrs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("PublicAddressRangeAuthorizedCidrsPager.GetAll() failed %s", err), "(Data) ibm_is_public_address_range_authorized_cidrs", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsPublicAddressRangeAuthorizedCidrsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRToMap(&modelItem) // #nosec G601
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidrs", "read", "PublicAddressRanges-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("authorized_cidrs", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting authorized_cidrs %s", err), "(Data) ibm_is_public_address_range_authorized_cidrs", "read", "authorized_cidrs-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsPublicAddressRangeAuthorizedCidrsID returns a reasonable ID for the list.
func dataSourceIBMIsPublicAddressRangeAuthorizedCidrsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRToMap(model *vpcv1.PublicAddressRangeAuthorizedCIDR) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	allocationMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRAllocationToMap(model.Allocation)
	if err != nil {
		return modelMap, err
	}
	modelMap["allocation"] = []map[string]interface{}{allocationMap}
	modelMap["availability_mode"] = *model.AvailabilityMode
	modelMap["cidr"] = *model.CIDR
	modelMap["href"] = *model.Href
	modelMap["id"] = *model.ID
	modelMap["ip_version"] = *model.IPVersion
	modelMap["name"] = *model.Name
	modelMap["network_prefix_length"] = flex.IntValue(model.NetworkPrefixLength)
	modelMap["resource_type"] = *model.ResourceType
	if model.Zone != nil {
		zoneMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCidrsZoneReferenceToMap(model.Zone)
		if err != nil {
			return modelMap, err
		}
		modelMap["zone"] = []map[string]interface{}{zoneMap}
	}
	return modelMap, nil
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCidrsPublicAddressRangeAuthorizedCIDRAllocationToMap(model *vpcv1.PublicAddressRangeAuthorizedCIDRAllocation) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["count"] = flex.IntValue(model.Count)
	modelMap["profile_family"] = *model.ProfileFamily
	return modelMap, nil
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCidrsZoneReferenceToMap(model *vpcv1.ZoneReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	modelMap["name"] = *model.Name
	return modelMap, nil
}
