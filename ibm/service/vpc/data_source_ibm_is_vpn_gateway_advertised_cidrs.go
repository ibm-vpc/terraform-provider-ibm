// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
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

func DataSourceIBMIsVPNGatewayAdvertisedCidrs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsVPNGatewayAdvertisedCidrsRead,

		Schema: map[string]*schema.Schema{
			"vpn_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The VPN gateway identifier.",
			},
			"advertised_cidrs": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The additional CIDRs advertised through any enabled routing protocol (for example, BGP). The routing protocol will advertise routes with these CIDRs and VPC prefixes as route destinations.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the first page of resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of resources that can be returned by the request.",
			},
			"next": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the next page of resources. This property is present for all pagesexcept the last page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIBMIsVPNGatewayAdvertisedCidrsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	listVPNGatewayAdvertisedCidrsOptions := &vpcv1.ListVPNGatewayAdvertisedCIDRsOptions{}

	listVPNGatewayAdvertisedCidrsOptions.SetVPNGatewayID(d.Get("vpn_gateway").(string))

	vpnGatewayAdvertisedCIDRCollection, _, err := vpcClient.ListVPNGatewayAdvertisedCIDRs(listVPNGatewayAdvertisedCidrsOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListVPNGatewayAdvertisedCidrsWithContext failed: %s", err.Error()), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId(dataSourceIBMIsVPNGatewayAdvertisedCidrsID(d))

	first := []map[string]interface{}{}
	if vpnGatewayAdvertisedCIDRCollection.First != nil {
		modelMap, err := DataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionFirstToMap(vpnGatewayAdvertisedCIDRCollection.First)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read", "first-to-map").GetDiag()
		}
		first = append(first, modelMap)
	}
	if err = d.Set("first", first); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting first: %s", err), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read", "set-first").GetDiag()
	}

	if err = d.Set("limit", flex.IntValue(vpnGatewayAdvertisedCIDRCollection.Limit)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting limit: %s", err), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read", "set-limit").GetDiag()
	}

	next := []map[string]interface{}{}
	if vpnGatewayAdvertisedCIDRCollection.Next != nil {
		modelMap, err := DataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionNextToMap(vpnGatewayAdvertisedCIDRCollection.Next)
		if err != nil {
			return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read", "next-to-map").GetDiag()
		}
		next = append(next, modelMap)
	}
	if err = d.Set("next", next); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting next: %s", err), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read", "set-next").GetDiag()
	}

	if err = d.Set("total_count", flex.IntValue(vpnGatewayAdvertisedCIDRCollection.TotalCount)); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_is_vpn_gateway_advertised_cidrs", "read", "set-total_count").GetDiag()
	}

	return nil
}

// dataSourceIBMIsVPNGatewayAdvertisedCidrsID returns a reasonable ID for the list.
func dataSourceIBMIsVPNGatewayAdvertisedCidrsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func DataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionFirstToMap(model *vpcv1.VPNGatewayAdvertisedCIDRCollectionFirst) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}

func DataSourceIBMIsVPNGatewayAdvertisedCidrsVPNGatewayAdvertisedCIDRCollectionNextToMap(model *vpcv1.VPNGatewayAdvertisedCIDRCollectionNext) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	modelMap["href"] = *model.Href
	return modelMap, nil
}
