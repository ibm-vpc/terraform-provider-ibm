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
	isBareMetalServerDisk = "disk"
)

func DataSourceIBMIsBareMetalServerDisk() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISBareMetalServerDiskRead,

		Schema: map[string]*schema.Schema{
			isBareMetalServerID: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server identifier",
			},

			isBareMetalServerDisk: {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The bare metal server disk identifier",
			},
			//disks

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
				Description: "The usage constraints to be matched against the requested bare metal server properties to determine compatibility.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bare_metal_server": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expression that must be satisfied by the properties of a bare metal server provisioned using the image data in this disk.",
						},
						"instance": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The expression that must be satisfied by the properties of a bare metal server provisioned using the image data in this disk.",
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
	}
}

func dataSourceIBMISBareMetalServerDiskRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	bareMetalServerID := d.Get(isBareMetalServerID).(string)
	bareMetalServerDiskID := d.Get(isBareMetalServerDisk).(string)
	sess, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	options := &vpcv1.GetBareMetalServerDiskOptions{
		BareMetalServerID: &bareMetalServerID,
		ID:                &bareMetalServerDiskID,
	}

	disk, response, err := sess.GetBareMetalServerDiskWithContext(context, options)
	if err != nil || disk == nil {
		return diag.FromErr(fmt.Errorf("[ERROR] Error getting Bare Metal Server (%s) disk (%s): %s\n%s", bareMetalServerID, bareMetalServerDiskID, err, response))
	}
	d.SetId(*disk.ID)
	d.Set(isBareMetalServerDiskHref, *disk.Href)
	d.Set(isBareMetalServerDiskInterfaceType, *disk.InterfaceType)
	d.Set(isBareMetalServerDiskName, *disk.Name)
	d.Set(isBareMetalServerDiskResourceType, *disk.ResourceType)
	d.Set(isBareMetalServerDiskSize, *disk.Size)
	allowedUses := []map[string]interface{}{}
	if disk.AllowedUse != nil {
		modelMap, err := ResourceceIBMIsBareMetalServerDiskAllowedUseToMap(disk.AllowedUse)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, err.Error(), "(Data) ibm_is_bare_metal_server_disk", "read")
			log.Println(tfErr.GetDiag())
		}
		allowedUses = append(allowedUses, modelMap)
	}
	if err = d.Set("allowed_use", allowedUses); err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("Error setting allowed_use: %s", err), "(Data) ibm_is_bare_metal_server_disk", "read")
		log.Println(tfErr.GetDiag())
	}
	return nil
}
