// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	isBareMetalServerID               = "bare_metal_server"
	isBareMetalServerDiskHref         = "href"
	isBareMetalServerDiskResourceType = "resource_type"
)

func DataSourceIBMIsBareMetalServerDisks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerDisksRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},

			//disks

			isBareMetalServerDisks: {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of bare metal server disks. Disk is a block device that is locally attached to the physical server. By default, the listed disks are sorted by their created_at property values, with the newest disk first.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						isBareMetalServerDiskHref: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this bare metal server disk",
						},
						isBareMetalServerDiskID: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this bare metal server disk",
						},
						isBareMetalServerDiskInterfaceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The disk interface used for attaching the disk. Supported values are [ nvme, sata ]",
						},
						isBareMetalServerDiskName: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this disk",
						},
						isBareMetalServerDiskResourceType: {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type",
						},
						isBareMetalServerDiskSize: {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The size of the disk in GB (gigabytes)",
						},
						"allowed_use": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The usage constraints to match against the requested instance or bare metal server properties to determine compatibility.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bare_metal_server": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "An image can only be used for bare metal instantiation if this expression resolves to true.The expression follows [Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md), but does not support built-in functions and macros. In addition, the following property is supported:- `enable_secure_boot` - (boolean) Indicates whether secure boot is enabled for this bare metal server.",
									},
									"instance": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "This image can only be used to provision a virtual server instance if the resulting instance would have property values that satisfy this expression.The expression follows [Common Expression Language](https://github.com/google/cel-spec/blob/master/doc/langdef.md), but does not support built-in functions and macros. In addition, the following variables are supported, corresponding to `Instance` properties:- `gpu.count` - (integer) The number of GPUs assigned to the instance- `gpu.manufacturer` - (string) The GPU manufacturer- `gpu.memory` - (integer) The overall amount of GPU memory in GiB (gibibytes)- `gpu.model` - (string) The GPU model- `enable_secure_boot` - (boolean) Indicates whether secure boot is enabled.",
									},
									"api_version": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The API version with which to evaluate the expressions.",
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

func dataSourceIBMISBareMetalServerDisksRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.ListBareMetalServerDisksOptions{
		BareMetalServerID: &bareMetalServerID,
	}

	diskCollection, response, err := sess.ListBareMetalServerDisksWithContext(context, options)
	disks := diskCollection.Disks
	if err != nil || disks == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) disks: %s\n%s", bareMetalServerID, err, response))
	}
	disksInfo := make([]map[string]interface{}, 0)
	for _, disk := range disks {
		l := map[string]interface{}{
			isBareMetalServerDiskHref:          disk.Href,
			isBareMetalServerDiskID:            disk.ID,
			isBareMetalServerDiskInterfaceType: disk.InterfaceType,
			isBareMetalServerDiskName:          disk.Name,
			isBareMetalServerDiskResourceType:  disk.ResourceType,
			isBareMetalServerDiskSize:          disk.Size,
		}
		if disk.AllowedUse != nil {
			usageConstraintList := []map[string]interface{}{}
			modelMap, err := ResourceceIBMIsBareMetalServerDiskAllowedUseToMap(disk.AllowedUse)
			if err != nil {
				tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_bare_metal_server_disks", "read")
				log.Println(tfErr.GetDiag())
			}
			usageConstraintList = append(usageConstraintList, modelMap)
			l["allowed_use"] = usageConstraintList
		}
		disksInfo = append(disksInfo, l)
	}

	d.SetId(dataSourceIBMISBMSDisksID(d))
	d.Set(isBareMetalServerDisks, disksInfo)
	return nil
}

// dataSourceIBMISBMSProfilesID returns a reasonable ID for a Bare Metal Server Disks list.
func dataSourceIBMISBMSDisksID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}
