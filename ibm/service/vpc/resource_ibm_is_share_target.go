// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/validate"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func ResourceIbmIsShareMountTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsShareMountTargetCreate,
		ReadContext:   resourceIbmIsShareMountTargetRead,
		UpdateContext: resourceIbmIsShareMountTargetUpdate,
		DeleteContext: resourceIbmIsShareMountTargetDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"share": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The file share identifier.",
			},
			"vpc": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share_mount_target", "name"),
				Description:  "The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			// "subnet": {
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// 	Description: "The unique identifier of the subnet associated with this file share target.Only virtual server instances in the same VPC as this subnet will be allowed to mount the file share. In the future, this property may be required and used to assign an IP address for the file share target.",
			// },
			"share_target": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of this target",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the share target was created.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share target.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the mount target.",
			},
			"mount_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The mount path for the share.The IP addresses used in the mount path are currently within the IBM services IP range, but are expected to change to be within one of the VPC's subnets in the future.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
		},
	}
}

func ResourceIbmIsShareMountTargetValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_share_mount_target", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsShareMountTargetCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createShareMountTargetOptions := &vpcv1.CreateShareTargetOptions{}

	createShareMountTargetOptions.SetShareID(d.Get("share").(string))
	vpcid := d.Get("vpc").(string)
	vpc := &vpcv1.VPCIdentity{
		ID: &vpcid,
	}
	createShareMountTargetOptions.SetVPC(vpc)
	if _, ok := d.GetOk("name"); ok {
		createShareMountTargetOptions.SetName(d.Get("name").(string))
	}
	// if subnetIntf, ok := d.GetOk("subnet"); ok {
	// 	subnet := subnetIntf.(string)
	// 	subnetIdentity := &vpcv1.SubnetIdentity{
	// 		ID: &subnet,
	// 	}
	// 	createShareMountTargetOptions.Subnet = subnetIdentity
	// }

	shareMountTarget, response, err := vpcClient.CreateShareTargetWithContext(context, createShareMountTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateShareMountTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	_, err = isWaitForTargetAvailable(context, vpcClient, *createShareMountTargetOptions.ShareID, *shareMountTarget.ID, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", *createShareMountTargetOptions.ShareID, *shareMountTarget.ID))
	d.Set("share_target", *shareMountTarget.ID)
	return resourceIbmIsShareMountTargetRead(context, d, meta)
}

func resourceIbmIsShareMountTargetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareMountTargetOptions := &vpcv1.GetShareTargetOptions{}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getShareMountTargetOptions.SetShareID(parts[0])
	getShareMountTargetOptions.SetID(parts[1])

	shareMountTarget, response, err := vpcClient.GetShareTargetWithContext(context, getShareMountTargetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareMountTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.Set("share_target", *shareMountTarget.ID)

	if err = d.Set("vpc", *shareMountTarget.VPC.ID); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if err = d.Set("name", *shareMountTarget.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	// if shareMountTarget.Subnet != nil {
	// 	if err = d.Set("subnet", *shareMountTarget.Subnet.ID); err != nil {
	// 		return diag.FromErr(fmt.Errorf("Error setting subnet: %s", err))
	// 	}
	// }
	if err = d.Set("created_at", shareMountTarget.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", shareMountTarget.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", shareMountTarget.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("mount_path", shareMountTarget.MountPath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mount_path: %s", err))
	}
	if err = d.Set("resource_type", shareMountTarget.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}

	return nil
}

func resourceIbmIsShareMountTargetUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateShareMountTargetOptions := &vpcv1.UpdateShareTargetOptions{}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateShareMountTargetOptions.SetShareID(parts[0])
	updateShareMountTargetOptions.SetID(parts[1])

	hasChange := false

	shareMountTargetPatchModel := &vpcv1.ShareMountTargetPatch{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		shareMountTargetPatchModel.Name = &name
		hasChange = true
	}

	if hasChange {
		shareMountTargetPatch, err := shareMountTargetPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] ShareMountTargetPatch AsPatch failed %s", err)
			return diag.FromErr(err)
		}
		updateShareMountTargetOptions.SetShareMountTargetPatch(shareMountTargetPatch)
		_, response, err := vpcClient.UpdateShareTargetWithContext(context, updateShareMountTargetOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateShareMountTargetWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIsShareMountTargetRead(context, d, meta)
}

func resourceIbmIsShareMountTargetDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteShareMountTargetOptions := &vpcv1.DeleteShareTargetOptions{}

	parts, err := flex.IdParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteShareMountTargetOptions.SetShareID(parts[0])
	deleteShareMountTargetOptions.SetID(parts[1])

	_, response, err := vpcClient.DeleteShareTargetWithContext(context, deleteShareMountTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteShareMountTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	_, err = isWaitForTargetDelete(context, vpcClient, d, parts[0], parts[1])
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func isWaitForTargetAvailable(context context.Context, vpcClient *vpcv1.VpcV1, shareid, targetid string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for target (%s) to be available.", targetid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating", "pending", "waiting"},
		Target:     []string{"stable", "failed"},
		Refresh:    isTargetRefreshFunc(context, vpcClient, shareid, targetid, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isTargetRefreshFunc(context context.Context, vpcClient *vpcv1.VpcV1, shareid, targetid string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		shareMountTargetOptions := &vpcv1.GetShareTargetOptions{}

		shareMountTargetOptions.SetShareID(shareid)
		shareMountTargetOptions.SetID(targetid)

		target, response, err := vpcClient.GetShareTargetWithContext(context, shareMountTargetOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting target: %s\n%s", err, response)
		}
		d.Set("lifecycle_state", *target.LifecycleState)
		if *target.LifecycleState == "stable" || *target.LifecycleState == "failed" {

			return target, *target.LifecycleState, nil

		}
		return target, "pending", nil
	}
}

func isWaitForTargetDelete(context context.Context, vpcClient *vpcv1.VpcV1, d *schema.ResourceData, shareid, targetid string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting", "stable"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {
			shareMountTargetOptions := &vpcv1.GetShareTargetOptions{}

			shareMountTargetOptions.SetShareID(shareid)
			shareMountTargetOptions.SetID(targetid)

			target, response, err := vpcClient.GetShareTargetWithContext(context, shareMountTargetOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return target, "done", nil
				}
				return nil, "", fmt.Errorf("Error Getting Target: %s\n%s", err, response)
			}
			if *target.LifecycleState == isInstanceFailed {
				return target, *target.LifecycleState, fmt.Errorf("The  target %s failed to delete: %v", targetid, err)
			}
			return target, "deleting", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
