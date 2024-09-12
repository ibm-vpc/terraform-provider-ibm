// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVolumeProfile       = "name"
	isVolumeProfileFamily = "family"
)

func DataSourceIBMISVolumeProfile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISVolumeProfileRead,

		Schema: map[string]*schema.Schema{

			isVolumeProfile: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Volume profile name",
			},

			isVolumeProfileFamily: {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Volume profile family",
			},
		},
	}
}

func dataSourceIBMISVolumeProfileRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	name := d.Get(isVolumeProfile).(string)

	err := volumeProfileGet(d, meta, name)
	if err != nil {
		return err
	}
	return nil
}

func volumeProfileGet(d *schema.ResourceData, meta interface{}, name string) error {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "ibm_cloud", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	getVolumeProfileOptions := &vpcv1.GetVolumeProfileOptions{
		Name: &name,
	}
	profile, _, err := sess.GetVolumeProfile(getVolumeProfileOptions)
	if err != nil {
		return err
	}
	// For lack of anything better, compose our id from profile name.
	d.SetId(*profile.Name)
	d.Set(isVolumeProfile, *profile.Name)
	d.Set(isVolumeProfileFamily, *profile.Family)
	return nil
}
