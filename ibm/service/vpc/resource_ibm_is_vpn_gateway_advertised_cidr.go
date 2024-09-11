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

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIBMIsVPNGatewayAdvertisedCidr() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMIsVPNGatewayAdvertisedCidrCreate,
		ReadContext:   resourceIBMIsVPNGatewayAdvertisedCidrRead,
		DeleteContext: resourceIBMIsVPNGatewayAdvertisedCidrDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"vpn_gateway": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The VPN gateway identifier.",
			},
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The IP address range in CIDR block notation.",
			},

			"status": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Response status code.",
			},
		},
	}
}

func resourceIBMIsVPNGatewayAdvertisedCidrCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		// Error is coming from SDK client, so it doesn't need to be discriminated.
		tfErr := flex.TerraformErrorf(err, err.Error(), "ibm_is_vpn_gateway_route_report", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	vpngatewayid := d.Get("vpn_gateway").(string)

	addVPNGatewayAdvertisedCIDROptions := &vpcv1.AddVPNGatewayAdvertisedCIDROptions{}

	addVPNGatewayAdvertisedCIDROptions.VPNGatewayID = core.StringPtr(vpngatewayid)
	addVPNGatewayAdvertisedCIDROptions.SetCIDR(d.Get("vpn_gateway").(string))

	response, err := vpcClient.AddVPNGatewayAdvertisedCIDRWithContext(context, addVPNGatewayAdvertisedCIDROptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("AddVPNGatewayAdvertisedCIDRWithContext failed: %s", err.Error()), "ibm_is_vpn_gateway_route_report", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	d.Set("status", int64(response.StatusCode))
	d.SetId(vpngatewayid)

	return resourceIBMIsVPNGatewayAdvertisedCidrRead(context, d, meta)
}

func resourceIBMIsVPNGatewayAdvertisedCidrRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	return nil
}

func resourceIBMIsVPNGatewayAdvertisedCidrDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
