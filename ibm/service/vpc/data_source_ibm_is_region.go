// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isRegionEndpoint = "endpoint"
	isRegionName     = "name"
	isRegionStatus   = "status"
)

func DataSourceIBMISRegion() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISRegionRead,

		Schema: map[string]*schema.Schema{

			isRegionEndpoint: {
				Type:     schema.TypeString,
				Computed: true,
			},

			isRegionName: {
				Type:     schema.TypeString,
				Optional: true,
			},

			isRegionStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceIBMISRegionRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	name := d.Get("name").(string)

	if name == "" {
		bmxSess, err := meta.(conns.ClientSession).BluemixSession()
		if err != nil {
			return err
		}
		name = bmxSess.Config.Region
	}
	return regionGet(d, meta, name)
}

func regionGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "ibm_cloud", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	getRegionOptions := &vpcv1.GetRegionOptions{
		Name: &name,
	}
	region, _, err := sess.GetRegion(getRegionOptions)
	if err != nil {
		return err
	}
	// For lack of anything better, compose our id from region name.
	d.SetId(*region.Name)
	d.Set(isRegionEndpoint, *region.Endpoint)
	d.Set(isRegionName, *region.Name)
	d.Set(isRegionStatus, *region.Status)
	return nil
}
