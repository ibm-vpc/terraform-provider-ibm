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

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsRead,

		Schema: map[string]*schema.Schema{
			"authorized_cidr_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public address range authorized CIDR identifier.",
			},
			"allocations_resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filters the collection to resources with an item in the `allocations` property with a`resource_type` property matching the specified value.",
			},
			"allocations": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The floating IPs and public address ranges allocated from this public address range authorized CIDR.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"address": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique IP address.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this floating IP.",
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
							Description: "The URL for this floating IP.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this floating IP.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name for this floating IP. The name is unique across all floating IPs in the region.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"cidr": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The public IP address block for this public address range, expressed in CIDR format.This property may [expand](https://cloud.ibm.com/apidocs/vpc#property-value-expansion) to support IPv6 address blocks in the future.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidr_allocations", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listPublicAddressRangeAuthorizedCIDRAllocationsOptions := &vpcv1.ListPublicAddressRangeAuthorizedCIDRAllocationsOptions{}

	listPublicAddressRangeAuthorizedCIDRAllocationsOptions.SetAuthorizedCIDRID(d.Get("authorized_cidr_id").(string))
	if _, ok := d.GetOk("allocations_resource_type"); ok {
		listPublicAddressRangeAuthorizedCIDRAllocationsOptions.SetAllocationsResourceType(d.Get("allocations_resource_type").(string))
	}

	var pager *vpcv1.PublicAddressRangeAuthorizedCIDRAllocationsPager
	pager, err = vpcClient.NewPublicAddressRangeAuthorizedCIDRAllocationsPager(listPublicAddressRangeAuthorizedCIDRAllocationsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidr_allocations", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	allItems, err := pager.GetAll()
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("PublicAddressRangeAuthorizedCIDRAllocationsPager.GetAll() failed %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocations", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsID(d))

	mapSlice := []map[string]interface{}{}
	for _, modelItem := range allItems {
		modelMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemToMap(modelItem)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidr_allocations", "read", "PublicAddressRanges-to-map").GetDiag()
		}
		mapSlice = append(mapSlice, modelMap)
	}

	if err = d.Set("allocations", mapSlice); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting allocations %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocations", "read", "allocations-set").GetDiag()
	}

	return nil
}

// dataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsID returns a reasonable ID for the list.
func dataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemToMap(model vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemIntf) (map[string]interface{}, error) {
	if _, ok := model.(*vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReference); ok {
		return DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReferenceToMap(model.(*vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReference))
	} else if _, ok := model.(*vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReference); ok {
		return DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReferenceToMap(model.(*vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReference))
	} else if _, ok := model.(*vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItem); ok {
		modelMap := make(map[string]interface{})
		model := model.(*vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItem)
		if model.Address != nil {
			modelMap["address"] = *model.Address
		}
		if model.CRN != nil {
			modelMap["crn"] = *model.CRN
		}
		if model.Deleted != nil {
			deletedMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDeletedToMap(model.Deleted)
			if err != nil {
				return modelMap, err
			}
			modelMap["deleted"] = []map[string]interface{}{deletedMap}
		}
		if model.Href != nil {
			modelMap["href"] = *model.Href
		}
		if model.ID != nil {
			modelMap["id"] = *model.ID
		}
		if model.Name != nil {
			modelMap["name"] = *model.Name
		}
		if model.ResourceType != nil {
			modelMap["resource_type"] = *model.ResourceType
		}
		if model.CIDR != nil {
			modelMap["cidr"] = *model.CIDR
		}
		return modelMap, nil
	} else {
		return nil, fmt.Errorf("Unrecognized vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemIntf subtype encountered")
	}
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDeletedToMap(model *vpcv1.Deleted) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["more_info"] = *model.MoreInfo
	return modelMap, nil
}

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReferenceToMap(model *vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["address"] = *model.Address
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDeletedToMap(model.Deleted)
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

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsPublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReferenceToMap(model *vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReference) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["cidr"] = *model.CIDR
	modelMap["crn"] = *model.CRN
	if model.Deleted != nil {
		deletedMap, err := DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationsDeletedToMap(model.Deleted)
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
