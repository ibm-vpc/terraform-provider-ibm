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
	isZoneName   = "name"
	isZoneRegion = "region"
	isZoneStatus = "status"

	isZoneDataCenter    = "data_center"
	isZoneUniversalName = "universal_name"
)

func DataSourceIBMISZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISZoneRead,

		Schema: map[string]*schema.Schema{

			isZoneName: {
				Type:     schema.TypeString,
				Required: true,
			},

			isZoneRegion: {
				Type:     schema.TypeString,
				Required: true,
			},

			isZoneStatus: {
				Type:     schema.TypeString,
				Computed: true,
			},
			isZoneDataCenter: {
				Type:     schema.TypeString,
				Computed: true,
			},
			isZoneUniversalName: {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceIBMISZoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	regionName := d.Get(isZoneRegion).(string)
	zoneName := d.Get(isZoneName).(string)
	return zoneGet(ctx, d, meta, regionName, zoneName)
}

func zoneGet(ctx context.Context, d *schema.ResourceData, meta interface{}, regionName, zoneName string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "(Data) ibm_is_zone", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	getRegionZoneOptions := &vpcv1.GetRegionZoneOptions{
		RegionName: &regionName,
		Name:       &zoneName,
	}
	zone, response, err := sess.GetRegionZoneWithContext(ctx, getRegionZoneOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetRegionZoneWithContext failed %s\n%s", err, response), "(Data) ibm_is_zone", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	// For lack of anything better, compose our id from region name + zone name.
	id := fmt.Sprintf("%s.%s", *zone.Region.Name, *zone.Name)
	d.SetId(id)
	if err = d.Set(isZoneName, *zone.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isZoneName), "(Data) ibm_is_zone", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isZoneRegion, *zone.Region.Name); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isZoneRegion), "(Data) ibm_is_zone", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if err = d.Set(isZoneStatus, *zone.Status); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isZoneStatus), "(Data) ibm_is_zone", "read")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if zone.DataCenter != nil {
		if err = d.Set(isZoneDataCenter, *zone.DataCenter); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isZoneDataCenter), "(Data) ibm_is_zone", "read")
			log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	if zone.UniversalName != nil {
		if err = d.Set(isZoneUniversalName, *zone.UniversalName); err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting %s", isZoneUniversalName), "(Data) ibm_is_zone", "read")
			log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
	}
	return nil
}
