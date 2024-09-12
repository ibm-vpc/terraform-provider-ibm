// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isVolumeProfiles = "profiles"
)

func DataSourceIBMISVolumeProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISVolumeProfilesRead,

		Schema: map[string]*schema.Schema{

			isVolumeProfiles: {
				Type:        schema.TypeList,
				Description: "List of Volume profile maps",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"family": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISVolumeProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	err := volumeProfilesList(d, meta)
	if err != nil {
		return err
	}
	return nil
}

func volumeProfilesList(d *schema.ResourceData, meta interface{}) error {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "ibm_cloud", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	start := ""
	allrecs := []vpcv1.VolumeProfile{}
	for {
		listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{}
		if start != "" {
			listVolumeProfilesOptions.Start = &start
		}
		availableProfiles, response, err := sess.ListVolumeProfiles(listVolumeProfilesOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error Fetching Volume Profiles %s\n%s", err, response)
		}
		start = flex.GetNext(availableProfiles.Next)
		allrecs = append(allrecs, availableProfiles.Profiles...)
		if start == "" {
			break
		}
	}

	// listVolumeProfilesOptions := &vpcv1.ListVolumeProfilesOptions{}
	// availableProfiles, response, err := sess.ListVolumeProfiles(listVolumeProfilesOptions)
	// if err != nil {
	// 	return fmt.Errorf("[ERROR] Error Fetching Volume Profiles %s\n%s", err, response)
	// }
	profilesInfo := make([]map[string]interface{}, 0)
	for _, profile := range allrecs {

		l := map[string]interface{}{
			"name":   *profile.Name,
			"family": *profile.Family,
		}
		profilesInfo = append(profilesInfo, l)
	}
	d.SetId(dataSourceIBMISVolumeProfilesID(d))
	d.Set(isVolumeProfiles, profilesInfo)
	return nil
}

// dataSourceIBMISVolumeProfilesID returns a reasonable ID for a Volume Profile list.
func dataSourceIBMISVolumeProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
