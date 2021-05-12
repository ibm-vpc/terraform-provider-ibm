// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func resourceIbmIsShareTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsShareTargetCreate,
		ReadContext:   resourceIbmIsShareTargetRead,
		UpdateContext: resourceIbmIsShareTargetUpdate,
		DeleteContext: resourceIbmIsShareTargetDelete,
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
				Optional:     true,
				ValidateFunc: InvokeValidator("ibm_is_share_target", "name"),
				Description:  "The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"subnet": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The unique identifier of the subnet associated with this file share target.Only virtual server instances in the same VPC as this subnetwill be allowed to mount the file share.In the future, this property may be required and used to assignan IP address for the file share target.",
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
			"security_groups": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of security groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"crn": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The security group's CRN.",
						},
						"deleted": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The security group's canonical URL.",
						},
						"id": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier for this security group.",
						},
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The user-defined name for this security group. Names must be unique within the VPC the security group resides in.",
						},
					},
				},
			},
		},
	}
}

func resourceIbmIsShareTargetValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_share_target", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsShareTargetCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createShareTargetOptions := &vpcv1.CreateShareTargetOptions{}

	createShareTargetOptions.SetShareID(d.Get("share").(string))
	vpcid := d.Get("vpc").(string)
	vpc := &vpcv1.ShareTargetPrototypeVPC{
		ID: &vpcid,
	}
	createShareTargetOptions.SetVPC(vpc)
	if _, ok := d.GetOk("name"); ok {
		createShareTargetOptions.SetName(d.Get("name").(string))
	}
	if subnetIntf, ok := d.GetOk("subnet"); ok {
		subnet := subnetIntf.(string)
		subnetIdentity := &vpcv1.SubnetIdentity{
			ID: &subnet,
		}
		createShareTargetOptions.Subnet = subnetIdentity
	}

	shareTarget, response, err := vpcClient.CreateShareTargetWithContext(context, createShareTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateShareTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] ************** calling wait for target available")
	_, err = isWaitForTargetAvailable(context, vpcClient, *createShareTargetOptions.ShareID, *shareTarget.ID, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", *createShareTargetOptions.ShareID, *shareTarget.ID))

	return resourceIbmIsShareTargetRead(context, d, meta)
}

func resourceIbmIsShareTargetMapToShareTargetPrototypeVpc(shareTargetPrototypeVpcMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPC {
	shareTargetPrototypeVpc := vpcv1.ShareTargetPrototypeVPC{}

	if shareTargetPrototypeVpcMap["id"] != nil {
		shareTargetPrototypeVpc.ID = core.StringPtr(shareTargetPrototypeVpcMap["id"].(string))
	}
	if shareTargetPrototypeVpcMap["crn"] != nil {
		shareTargetPrototypeVpc.CRN = core.StringPtr(shareTargetPrototypeVpcMap["crn"].(string))
	}
	if shareTargetPrototypeVpcMap["href"] != nil {
		shareTargetPrototypeVpc.Href = core.StringPtr(shareTargetPrototypeVpcMap["href"].(string))
	}

	return shareTargetPrototypeVpc
}

func resourceIbmIsShareTargetMapToShareTargetPrototypeVpcVPCIdentityByID(shareTargetPrototypeVpcVPCIdentityByIDMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPCVPCIdentityByID {
	shareTargetPrototypeVpcVPCIdentityByID := vpcv1.ShareTargetPrototypeVPCVPCIdentityByID{}

	shareTargetPrototypeVpcVPCIdentityByID.ID = core.StringPtr(shareTargetPrototypeVpcVPCIdentityByIDMap["id"].(string))

	return shareTargetPrototypeVpcVPCIdentityByID
}

func resourceIbmIsShareTargetMapToShareTargetPrototypeVpcVPCIdentityByCRN(shareTargetPrototypeVpcVPCIdentityByCRNMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPCVPCIdentityByCRN {
	shareTargetPrototypeVpcVPCIdentityByCRN := vpcv1.ShareTargetPrototypeVPCVPCIdentityByCRN{}

	shareTargetPrototypeVpcVPCIdentityByCRN.CRN = core.StringPtr(shareTargetPrototypeVpcVPCIdentityByCRNMap["crn"].(string))

	return shareTargetPrototypeVpcVPCIdentityByCRN
}

func resourceIbmIsShareTargetMapToShareTargetPrototypeVpcVPCIdentityByHref(shareTargetPrototypeVpcVPCIdentityByHrefMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPCVPCIdentityByHref {
	shareTargetPrototypeVpcVPCIdentityByHref := vpcv1.ShareTargetPrototypeVPCVPCIdentityByHref{}

	shareTargetPrototypeVpcVPCIdentityByHref.Href = core.StringPtr(shareTargetPrototypeVpcVPCIdentityByHrefMap["href"].(string))

	return shareTargetPrototypeVpcVPCIdentityByHref
}

func resourceIbmIsShareTargetMapToSubnetIdentity(subnetIdentityMap map[string]interface{}) vpcv1.SubnetIdentity {
	subnetIdentity := vpcv1.SubnetIdentity{}

	if subnetIdentityMap["id"] != nil {
		subnetIdentity.ID = core.StringPtr(subnetIdentityMap["id"].(string))
	}
	if subnetIdentityMap["crn"] != nil {
		subnetIdentity.CRN = core.StringPtr(subnetIdentityMap["crn"].(string))
	}
	if subnetIdentityMap["href"] != nil {
		subnetIdentity.Href = core.StringPtr(subnetIdentityMap["href"].(string))
	}

	return subnetIdentity
}

func resourceIbmIsShareTargetMapToSubnetIdentityByID(subnetIdentityByIDMap map[string]interface{}) vpcv1.SubnetIdentityByID {
	subnetIdentityByID := vpcv1.SubnetIdentityByID{}

	subnetIdentityByID.ID = core.StringPtr(subnetIdentityByIDMap["id"].(string))

	return subnetIdentityByID
}

func resourceIbmIsShareTargetMapToSubnetIdentityByCRN(subnetIdentityByCRNMap map[string]interface{}) vpcv1.SubnetIdentityByCRN {
	subnetIdentityByCRN := vpcv1.SubnetIdentityByCRN{}

	subnetIdentityByCRN.CRN = core.StringPtr(subnetIdentityByCRNMap["crn"].(string))

	return subnetIdentityByCRN
}

func resourceIbmIsShareTargetMapToSubnetIdentityByHref(subnetIdentityByHrefMap map[string]interface{}) vpcv1.SubnetIdentityByHref {
	subnetIdentityByHref := vpcv1.SubnetIdentityByHref{}

	subnetIdentityByHref.Href = core.StringPtr(subnetIdentityByHrefMap["href"].(string))

	return subnetIdentityByHref
}

func resourceIbmIsShareTargetRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareTargetOptions := &vpcv1.GetShareTargetOptions{}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	getShareTargetOptions.SetShareID(parts[0])
	getShareTargetOptions.SetID(parts[1])

	shareTarget, response, err := vpcClient.GetShareTargetWithContext(context, getShareTargetOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareTargetWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if err = d.Set("name", shareTarget.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if shareTarget.Subnet != nil {
		if err = d.Set("subnet", *shareTarget.Subnet.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting subnet: %s", err))
		}
	}
	if err = d.Set("created_at", shareTarget.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("href", shareTarget.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", shareTarget.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("mount_path", shareTarget.MountPath); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting mount_path: %s", err))
	}
	if err = d.Set("resource_type", shareTarget.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	securityGroups := []map[string]interface{}{}
	for _, securityGroupsItem := range shareTarget.SecurityGroups {
		securityGroupsItemMap := resourceIbmIsShareTargetSecurityGroupReferenceToMap(securityGroupsItem)
		securityGroups = append(securityGroups, securityGroupsItemMap)
	}
	if err = d.Set("security_groups", securityGroups); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting security_groups: %s", err))
	}

	return nil
}

func resourceIbmIsShareTargetShareTargetPrototypeVpcToMap(shareTargetPrototypeVpc vpcv1.ShareTargetPrototypeVPC) map[string]interface{} {
	shareTargetPrototypeVpcMap := map[string]interface{}{}

	shareTargetPrototypeVpcMap["id"] = shareTargetPrototypeVpc.ID
	shareTargetPrototypeVpcMap["crn"] = shareTargetPrototypeVpc.CRN
	shareTargetPrototypeVpcMap["href"] = shareTargetPrototypeVpc.Href

	return shareTargetPrototypeVpcMap
}

func resourceIbmIsShareTargetShareTargetPrototypeVpcVPCIdentityByIDToMap(shareTargetPrototypeVpcVPCIdentityByID vpcv1.ShareTargetPrototypeVPCVPCIdentityByID) map[string]interface{} {
	shareTargetPrototypeVpcVPCIdentityByIDMap := map[string]interface{}{}

	shareTargetPrototypeVpcVPCIdentityByIDMap["id"] = shareTargetPrototypeVpcVPCIdentityByID.ID

	return shareTargetPrototypeVpcVPCIdentityByIDMap
}

func resourceIbmIsShareTargetShareTargetPrototypeVpcVPCIdentityByCRNToMap(shareTargetPrototypeVpcVPCIdentityByCRN vpcv1.ShareTargetPrototypeVPCVPCIdentityByCRN) map[string]interface{} {
	shareTargetPrototypeVpcVPCIdentityByCRNMap := map[string]interface{}{}

	shareTargetPrototypeVpcVPCIdentityByCRNMap["crn"] = shareTargetPrototypeVpcVPCIdentityByCRN.CRN

	return shareTargetPrototypeVpcVPCIdentityByCRNMap
}

func resourceIbmIsShareTargetShareTargetPrototypeVpcVPCIdentityByHrefToMap(shareTargetPrototypeVpcVPCIdentityByHref vpcv1.ShareTargetPrototypeVPCVPCIdentityByHref) map[string]interface{} {
	shareTargetPrototypeVpcVPCIdentityByHrefMap := map[string]interface{}{}

	shareTargetPrototypeVpcVPCIdentityByHrefMap["href"] = shareTargetPrototypeVpcVPCIdentityByHref.Href

	return shareTargetPrototypeVpcVPCIdentityByHrefMap
}

func resourceIbmIsShareTargetSubnetIdentityToMap(subnetIdentity vpcv1.SubnetIdentity) map[string]interface{} {
	subnetIdentityMap := map[string]interface{}{}

	subnetIdentityMap["id"] = subnetIdentity.ID
	subnetIdentityMap["crn"] = subnetIdentity.CRN
	subnetIdentityMap["href"] = subnetIdentity.Href

	return subnetIdentityMap
}

func resourceIbmIsShareTargetSubnetIdentityByIDToMap(subnetIdentityByID vpcv1.SubnetIdentityByID) map[string]interface{} {
	subnetIdentityByIDMap := map[string]interface{}{}

	subnetIdentityByIDMap["id"] = subnetIdentityByID.ID

	return subnetIdentityByIDMap
}

func resourceIbmIsShareTargetSubnetIdentityByCRNToMap(subnetIdentityByCRN vpcv1.SubnetIdentityByCRN) map[string]interface{} {
	subnetIdentityByCRNMap := map[string]interface{}{}

	subnetIdentityByCRNMap["crn"] = subnetIdentityByCRN.CRN

	return subnetIdentityByCRNMap
}

func resourceIbmIsShareTargetSubnetIdentityByHrefToMap(subnetIdentityByHref vpcv1.SubnetIdentityByHref) map[string]interface{} {
	subnetIdentityByHrefMap := map[string]interface{}{}

	subnetIdentityByHrefMap["href"] = subnetIdentityByHref.Href

	return subnetIdentityByHrefMap
}

func resourceIbmIsShareTargetSecurityGroupReferenceToMap(securityGroupReference vpcv1.SecurityGroupReference) map[string]interface{} {
	securityGroupReferenceMap := map[string]interface{}{}

	securityGroupReferenceMap["crn"] = securityGroupReference.CRN
	if securityGroupReference.Deleted != nil {
		DeletedMap := resourceIbmIsShareTargetSecurityGroupReferenceDeletedToMap(*securityGroupReference.Deleted)
		securityGroupReferenceMap["deleted"] = []map[string]interface{}{DeletedMap}
	}
	securityGroupReferenceMap["href"] = securityGroupReference.Href
	securityGroupReferenceMap["id"] = securityGroupReference.ID
	securityGroupReferenceMap["name"] = securityGroupReference.Name

	return securityGroupReferenceMap
}

func resourceIbmIsShareTargetSecurityGroupReferenceDeletedToMap(securityGroupReferenceDeleted vpcv1.SecurityGroupReferenceDeleted) map[string]interface{} {
	securityGroupReferenceDeletedMap := map[string]interface{}{}

	securityGroupReferenceDeletedMap["more_info"] = securityGroupReferenceDeleted.MoreInfo

	return securityGroupReferenceDeletedMap
}

func resourceIbmIsShareTargetUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateShareTargetOptions := &vpcv1.UpdateShareTargetOptions{}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateShareTargetOptions.SetShareID(parts[0])
	updateShareTargetOptions.SetID(parts[1])

	hasChange := false

	shareTargetPatchModel := &vpcv1.ShareTargetPatch{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		shareTargetPatchModel.Name = &name
		hasChange = true
	}

	if hasChange {
		shareTargetPatch, err := shareTargetPatchModel.AsPatch()
		if err != nil {
			log.Printf("[DEBUG] ShareTargetPatch AsPatch failed %s", err)
			return diag.FromErr(err)
		}
		updateShareTargetOptions.SetShareTargetPatch(shareTargetPatch)
		_, response, err := vpcClient.UpdateShareTargetWithContext(context, updateShareTargetOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateShareTargetWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}

	return resourceIbmIsShareTargetRead(context, d, meta)
}

func resourceIbmIsShareTargetDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	deleteShareTargetOptions := &vpcv1.DeleteShareTargetOptions{}

	parts, err := idParts(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteShareTargetOptions.SetShareID(parts[0])
	deleteShareTargetOptions.SetID(parts[1])

	_, response, err := vpcClient.DeleteShareTargetWithContext(context, deleteShareTargetOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteShareTargetWithContext failed %s\n%s", err, response)
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
		shareTargetOptions := &vpcv1.GetShareTargetOptions{}

		shareTargetOptions.SetShareID(shareid)
		shareTargetOptions.SetID(targetid)

		target, response, err := vpcClient.GetShareTargetWithContext(context, shareTargetOptions)
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
			shareTargetOptions := &vpcv1.GetShareTargetOptions{}

			shareTargetOptions.SetShareID(shareid)
			shareTargetOptions.SetID(targetid)

			target, response, err := vpcClient.GetShareTargetWithContext(context, shareTargetOptions)
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
