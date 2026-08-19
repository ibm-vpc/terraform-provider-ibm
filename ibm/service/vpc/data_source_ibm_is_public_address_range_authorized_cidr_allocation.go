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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocation() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationRead,

		Schema: map[string]*schema.Schema{
			"authorized_cidr_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public address range authorized CIDR identifier.",
			},
			"authorized_cidr_allocation_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The public address range authorized CIDR allocation identifier.",
			},
			"address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique IP address.",
			},
			"cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public IP address block for this public address range, expressed in CIDR format.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this allocation.",
			},
			"deleted": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "If present, this property indicates the referenced resource has been deleted, and provides some supplementary information.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"more_info": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A link to documentation about deleted resources.",
						},
					},
				},
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this allocation.",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name for this allocation.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The resource type.",
			},
		},
	}
}

func dataSourceIBMIsPublicAddressRangeAuthorizedCIDRAllocationRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	getOptions := &vpcv1.GetPublicAddressRangeAuthorizedCIDRAllocationOptions{}
	getOptions.SetAuthorizedCIDRID(d.Get("authorized_cidr_id").(string))
	getOptions.SetID(d.Get("authorized_cidr_allocation_id").(string))

	allocationIntf, _, err := vpcClient.GetPublicAddressRangeAuthorizedCIDRAllocationWithContext(context, getOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetPublicAddressRangeAuthorizedCIDRAllocationWithContext failed: %s", err.Error()), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	switch allocation := allocationIntf.(type) {
	case *vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemFloatingIPReference:
		d.SetId(*allocation.ID)
		if err = d.Set("address", allocation.Address); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting address: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-address").GetDiag()
		}
		if err = d.Set("crn", allocation.CRN); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-crn").GetDiag()
		}
		if allocation.Deleted != nil {
			deleted := []map[string]interface{}{{"more_info": *allocation.Deleted.MoreInfo}}
			if err = d.Set("deleted", deleted); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting deleted: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-deleted").GetDiag()
			}
		}
		if err = d.Set("href", allocation.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-href").GetDiag()
		}
		if err = d.Set("name", allocation.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-name").GetDiag()
		}
		if err = d.Set("resource_type", allocation.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-resource_type").GetDiag()
		}
	case *vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItemPublicAddressRangeReference:
		d.SetId(*allocation.ID)
		if err = d.Set("cidr", allocation.CIDR); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cidr: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-cidr").GetDiag()
		}
		if err = d.Set("crn", allocation.CRN); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-crn").GetDiag()
		}
		if allocation.Deleted != nil {
			deleted := []map[string]interface{}{{"more_info": *allocation.Deleted.MoreInfo}}
			if err = d.Set("deleted", deleted); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting deleted: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-deleted").GetDiag()
			}
		}
		if err = d.Set("href", allocation.Href); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-href").GetDiag()
		}
		if err = d.Set("name", allocation.Name); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-name").GetDiag()
		}
		if err = d.Set("resource_type", allocation.ResourceType); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-resource_type").GetDiag()
		}
	case *vpcv1.PublicAddressRangeAuthorizedCIDRAllocationItem:
		// Base type returned when SDK unmarshaler doesn't dispatch to a concrete subtype.
		// Dispatch manually based on resource_type.
		if allocation.ID != nil {
			d.SetId(*allocation.ID)
		}
		if allocation.ResourceType != nil {
			if err = d.Set("resource_type", allocation.ResourceType); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_type: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-resource_type").GetDiag()
			}
		}
		if allocation.CIDR != nil {
			if err = d.Set("cidr", allocation.CIDR); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting cidr: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-cidr").GetDiag()
			}
		}
		if allocation.Address != nil {
			if err = d.Set("address", allocation.Address); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting address: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-address").GetDiag()
			}
		}
		if allocation.CRN != nil {
			if err = d.Set("crn", allocation.CRN); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting crn: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-crn").GetDiag()
			}
		}
		if allocation.Deleted != nil {
			deleted := []map[string]interface{}{{"more_info": *allocation.Deleted.MoreInfo}}
			if err = d.Set("deleted", deleted); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting deleted: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-deleted").GetDiag()
			}
		}
		if allocation.Href != nil {
			if err = d.Set("href", allocation.Href); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting href: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-href").GetDiag()
			}
		}
		if allocation.Name != nil {
			if err = d.Set("name", allocation.Name); err != nil {
				return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting name: %s", err), "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "set-name").GetDiag()
			}
		}
	default:
		return flex.DiscriminatedTerraformErrorf(fmt.Errorf("unrecognized allocation type"), "unrecognized allocation type", "(Data) ibm_is_public_address_range_authorized_cidr_allocation", "read", "type-switch").GetDiag()
	}

	return nil
}
