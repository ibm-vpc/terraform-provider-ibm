// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ()

func ResourceIBMISInstanceDiskManagement() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMisInstanceDiskManagementCreate,
		ReadContext:   resourceIBMisInstanceDiskManagementRead,
		UpdateContext: resourceIBMisInstanceDiskManagementUpdate,
		DeleteContext: resourceIBMisInstanceDiskManagementDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"instance": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the instance for which disks has to be managed",
			},
			"disks": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Disk information that has to be updated.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this instance disk.",
						},
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.InvokeValidator("ibm_is_instance_disk_management", "name"),
							Description:  "The user-defined name for this disk. The disk will be updated with this new name",
						},
					},
				},
			},
		},
	}
}

func ResourceIBMISInstanceDiskManagementValidator() *validate.ResourceValidator {

	validateSchema := make([]validate.ValidateSchema, 0)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Required:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`,
			MinValueLength:             1,
			MaxValueLength:             63})

	ibmISInstanceDiskManagementValidator := validate.ResourceValidator{ResourceName: "ibm_is_instance_disk_management", Schema: validateSchema}
	return &ibmISInstanceDiskManagementValidator
}

func resourceIBMisInstanceDiskManagementCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "ibm_cloud", "create")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	instance := d.Get("instance").(string)
	disks := d.Get("disks")
	diskUpdate := disks.([]interface{})

	for _, disk := range diskUpdate {
		diskItem := disk.(map[string]interface{})

		namestr := diskItem["name"].(string)
		diskid := diskItem["id"].(string)

		updateInstanceDiskOptions := &vpcv1.UpdateInstanceDiskOptions{}
		updateInstanceDiskOptions.SetInstanceID(instance)
		updateInstanceDiskOptions.SetID(diskid)
		instanceDiskPatchModel := &vpcv1.InstanceDiskPatch{
			Name: &namestr,
		}

		instanceDiskPatch, err := instanceDiskPatchModel.AsPatch()
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling asPatch for InstanceDiskPatch: %s", err)
		}
		updateInstanceDiskOptions.SetInstanceDiskPatch(instanceDiskPatch)

		_, response, err := sess.UpdateInstanceDisk(updateInstanceDiskOptions)
		if err != nil {
			return fmt.Errorf("[ERROR] Error calling UpdateInstanceDisk: %s %s", err, response)
		}

	}
	d.SetId(instance)
	return resourceIBMisInstanceDiskManagementRead(d, meta)
}

func resourceIBMisInstanceDiskManagementUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "ibm_cloud", "create")
		log.Printf("[DEBUG] %s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}
	if d.HasChange("disks") && !d.IsNewResource() {

		disks := d.Get("disks")
		diskUpdate := disks.([]interface{})

		for _, disk := range diskUpdate {
			diskItem := disk.(map[string]interface{})
			namestr := diskItem["name"].(string)
			diskid := diskItem["id"].(string)

			updateInstanceDiskOptions := &vpcv1.UpdateInstanceDiskOptions{}
			updateInstanceDiskOptions.SetInstanceID(d.Id())
			updateInstanceDiskOptions.SetID(diskid)
			instanceDiskPatchModel := &vpcv1.InstanceDiskPatch{
				Name: &namestr,
			}

			instanceDiskPatch, err := instanceDiskPatchModel.AsPatch()
			if err != nil {
				return fmt.Errorf("[ERROR] Error calling asPatch for InstanceDiskPatch: %s", err)
			}
			updateInstanceDiskOptions.SetInstanceDiskPatch(instanceDiskPatch)

			_, _, err = sess.UpdateInstanceDisk(updateInstanceDiskOptions)
			if err != nil {
				return fmt.Errorf("[ERROR] Error updating instance disk: %s", err)
			}

		}
	}
	return resourceIBMisInstanceDiskManagementRead(d, meta)
}

func resourceIBMisInstanceDiskManagementDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.SetId("")
	return nil
}

func resourceIBMisInstanceDiskManagementRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	d.Set("instance", d.Id())

	return nil
}
