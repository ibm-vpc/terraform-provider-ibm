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
	isZoneNames = "zones"
)

func DataSourceIBMISZones() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISZonesRead,

		Schema: map[string]*schema.Schema{

			isZoneRegion: {
				Type:     schema.TypeString,
				Required: true,
			},

			isZoneStatus: {
				Type:     schema.TypeString,
				Optional: true,
			},

			isZoneNames: {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceIBMISZonesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	regionName := d.Get(isZoneRegion).(string)
	return zonesList(d, context, meta, regionName)
}

func zonesList(d *schema.ResourceData, context context.Context, meta interface{}, regionName string) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "(Data) ibm_is_zones", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	listRegionZonesOptions := &vpcv1.ListRegionZonesOptions{
		RegionName: &regionName,
	}
	availableZones, _, err := sess.ListRegionZonesWithContext(context, listRegionZonesOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListRegionZonesWithContext failed: %s", err.Error()), "(Data) ibm_is_zones", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	names := make([]string, 0)
	status := d.Get(isZoneStatus).(string)
	for _, zone := range availableZones.Zones {
		if status == "" || *zone.Status == status {
			names = append(names, *zone.Name)
		}
	}
	d.SetId(dataSourceIBMISZonesId(d))
	d.Set(isZoneNames, names)
	return nil
}

// dataSourceIBMISZonesId returns a reasonable ID for a zone list.
func dataSourceIBMISZonesId(d *schema.ResourceData) string {
	// Our zone list is not guaranteed to be stable because the content
	// of the list can vary between two calls if any of the following
	// events occur between calls:
	// - a zone is added to our region
	// - a zone is dropped from our region
	// - we are using the status filter and the status of one or more
	//   zones changes between calls.
	//
	// For simplicity we are using a timestamp for the required terraform id.
	// If we find through usage that this choice is too ephemeral for our users
	// then we can change this function to use a more stable id, perhaps
	// composed from a hash of the list contents. But, for now, a timestamp
	// is good enough.
	return time.Now().UTC().String()
}
