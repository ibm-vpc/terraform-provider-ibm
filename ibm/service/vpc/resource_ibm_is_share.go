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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

const (
	isFileShareAccessTags = "access_tags"
	isFileShareTags       = "tags"
)

func ResourceIbmIsShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsShareCreate,
		ReadContext:   resourceIbmIsShareRead,
		UpdateContext: resourceIbmIsShareUpdate,
		DeleteContext: resourceIbmIsShareDelete,
		Importer:      &schema.ResourceImporter{},

		CustomizeDiff: customdiff.All(
			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceTagsCustomizeDiff(diff)
				},
			),

			customdiff.Sequence(
				func(_ context.Context, diff *schema.ResourceDiff, v interface{}) error {
					return flex.ResourceSharesValidate(diff)
				}),
		),

		Schema: map[string]*schema.Schema{
			"encryption_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The CRN of the key to use for encrypting this file share.If no encryption key is provided, the share will not be encrypted.",
			},
			"initial_owner_gid": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The owner assigned to the file share at creation. The initial group identifier for the file share. Subsequent changes to the owner must be performed by a virtual server instance that has mounted the file share.",
			},
			"initial_owner_uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The owner assigned to the file share at creation. The initial user identifier for the file share. Subsequent changes to the owner must be performed by a virtual server instance that has mounted the file share.",
			},
			"iops": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share", "iops"),
				Description:  "The maximum input/output operation performance bandwidth per second for the file share.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share", "name"),
				Description:  "The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique name for this share profile.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validate.InvokeValidator("ibm_is_share", "size"),
				Description:  "The size of the file share rounded up to the next gigabyte.",
			},
			"share_target_prototype": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Share targets for the file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The user-defined name for this share target. Names must be unique within the share the share target resides in. If unspecified, the name will be a hyphenated list of randomly-selected words.",
						},
						"vpc": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The unique identifier of the VPC in which instances can mount the file share using this share target.This property will be removed in a future release.The `subnet` property should be used instead.",
						},
					},
				},
			},
			"zone": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The globally unique name of the zone this file share will reside in.",
			},
			isFileShareAccessTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share", "accesstag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
			isFileShareTags: {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString, ValidateFunc: validate.InvokeValidator("ibm_is_share", "tag")},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags",
			},
			"share_targets": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Mount targets for the file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share target.",
						},
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this share target.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this share target.",
						},
						"resource_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
					},
				},
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the file share is created.",
			},
			"crn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this share.",
			},
			"encryption": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of encryption used for this file share.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the file share.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
		},
	}
}

func ResourceIbmIsShareValidator() *validate.ResourceValidator {
	validateSchema := make([]validate.ValidateSchema, 1)
	validateSchema = append(validateSchema,
		validate.ValidateSchema{
			Identifier:                 "iops",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "100",
			MaxValue:                   "48000",
		},
		validate.ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		validate.ValidateSchema{
			Identifier:                 "size",
			ValidateFunctionIdentifier: validate.IntBetween,
			Type:                       validate.TypeInt,
			Optional:                   true,
			MinValue:                   "10",
			MaxValue:                   "32000",
		},
		validate.ValidateSchema{
			Identifier:                 "tag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^[A-Za-z0-9:_ .-]+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
		validate.ValidateSchema{
			Identifier:                 "accesstag",
			ValidateFunctionIdentifier: validate.ValidateRegexpLen,
			Type:                       validate.TypeString,
			Optional:                   true,
			Regexp:                     `^([ ]*[A-Za-z0-9:_.-]+[ ]*)+$`,
			MinValueLength:             1,
			MaxValueLength:             128,
		},
	)

	resourceValidator := validate.ResourceValidator{ResourceName: "ibm_is_share", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsShareCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createShareOptions := &vpcv1.CreateShareOptions{}
	sharePrototype := &vpcv1.SharePrototype{}
	if encryptionKeyIntf, ok := d.GetOk("encryption_key"); ok {
		encryptionKey := encryptionKeyIntf.(string)
		encryptionKeyIdentity := &vpcv1.EncryptionKeyIdentity{
			CRN: &encryptionKey,
		}
		sharePrototype.EncryptionKey = encryptionKeyIdentity
	}
	initial_owner := &vpcv1.ShareInitialOwner{}
	if o_gid, ok := d.GetOk("initial_owner_gid"); ok {
		o_gidstr := o_gid.(int64)
		initial_owner.Gid = &o_gidstr
		sharePrototype.InitialOwner = initial_owner
	}
	if u_gid, ok := d.GetOk("initial_owner_gid"); ok {
		o_uidstr := u_gid.(int64)
		initial_owner.Uid = &o_uidstr
		sharePrototype.InitialOwner = initial_owner
	}
	if iopsIntf, ok := d.GetOk("iops"); ok {
		iops := int64(iopsIntf.(int))
		sharePrototype.Iops = &iops
	}
	if nameIntf, ok := d.GetOk("name"); ok {
		name := nameIntf.(string)
		sharePrototype.Name = &name
	}
	if profileIntf, ok := d.GetOk("profile"); ok {
		profileStr := profileIntf.(string)
		profile := &vpcv1.SharePrototypeProfile{
			Name: &profileStr,
		}
		sharePrototype.Profile = profile
	}
	if resgrp, ok := d.GetOk("resource_group"); ok {
		resgrpstr := resgrp.(string)
		resourceGroup := &vpcv1.ResourceGroupIdentity{
			ID: &resgrpstr,
		}
		sharePrototype.ResourceGroup = resourceGroup
	}
	if sizeIntf, ok := d.GetOk("size"); ok {
		size := int64(sizeIntf.(int))
		sharePrototype.Size = &size
	}
	if shareTargetPrototypeIntf, ok := d.GetOk("share_target_prototype"); ok {
		var targets []vpcv1.ShareTargetPrototype
		for _, e := range shareTargetPrototypeIntf.([]interface{}) {
			value := e.(map[string]interface{})
			targetsItem := resourceIbmIsShareMapToShareTargetPrototype(value)
			targets = append(targets, targetsItem)
		}
		sharePrototype.Targets = targets
	}
	if zone, ok := d.GetOk("zone"); ok {
		zonestr := zone.(string)
		zone := &vpcv1.ZoneIdentity{
			Name: &zonestr,
		}
		sharePrototype.Zone = zone
	}
	createShareOptions.SharePrototype = sharePrototype
	share, response, err := vpcClient.CreateShareWithContext(context, createShareOptions)
	if err != nil {
		log.Printf("[DEBUG] CreateShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	_, err = isWaitForShareAvailable(context, vpcClient, *share.ID, d, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(*share.ID)
	if _, ok := d.GetOk(isFileShareTags); ok {
		oldList, newList := d.GetChange(isFileShareTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *share.CRN, "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error creating file share (%s) tags: %s", d.Id(), err)
		}
	}
	if _, ok := d.GetOk(isFileShareAccessTags); ok {
		oldList, newList := d.GetChange(isFileShareAccessTags)
		err = flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, *share.CRN, "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error creating file share (%s) access tags: %s", d.Id(), err)
		}
	}
	return resourceIbmIsShareRead(context, d, meta)
}

func resourceIbmIsShareMapToShareTargetPrototype(shareTargetPrototypeMap map[string]interface{}) vpcv1.ShareTargetPrototype {
	shareTargetPrototype := vpcv1.ShareTargetPrototype{}

	if nameIntf, ok := shareTargetPrototypeMap["name"]; ok && nameIntf != "" {
		shareTargetPrototype.Name = core.StringPtr(nameIntf.(string))
	}

	if vpcIntf, ok := shareTargetPrototypeMap["vpc"]; ok && vpcIntf != "" {
		vpc := vpcIntf.(string)
		shareTargetPrototype.VPC = &vpcv1.ShareTargetPrototypeVPC{
			ID: &vpc,
		}
	}

	return shareTargetPrototype
}

func resourceIbmIsShareRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareOptions := &vpcv1.GetShareOptions{}

	getShareOptions.SetID(d.Id())

	share, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	if share.EncryptionKey != nil {
		if err = d.Set("encryption_key", *share.EncryptionKey.CRN); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting encryption_key: %s", err))
		}
	}

	if err = d.Set("iops", flex.IntValue(share.Iops)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting iops: %s", err))
	}
	if err = d.Set("name", share.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}
	if share.Profile != nil {
		if err = d.Set("profile", *share.Profile.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting profile: %s", err))
		}
	}
	if share.ResourceGroup != nil {
		if err = d.Set("resource_group", *share.ResourceGroup.ID); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group: %s", err))
		}
	}
	if err = d.Set("size", flex.IntValue(share.Size)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting size: %s", err))
	}
	targets := []map[string]interface{}{}
	if share.Targets != nil {
		for _, targetsItem := range share.Targets {
			targetsItemMap := dataSourceShareTargetsToMap(targetsItem)
			targets = append(targets, targetsItemMap)
		}
	}
	if err = d.Set("share_targets", targets); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting targets: %s", err))
	}
	if share.Zone != nil {
		if err = d.Set("zone", *share.Zone.Name); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting zone: %s", err))
		}
	}
	if err = d.Set("created_at", share.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", share.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("encryption", share.Encryption); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption: %s", err))
	}
	if err = d.Set("href", share.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("lifecycle_state", share.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("resource_type", share.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	tags, err := flex.GetGlobalTagsUsingCRN(meta, *share.CRN, "", isUserTagType)
	if err != nil {
		log.Printf(
			"Error getting shares (%s) tags: %s", d.Id(), err)
	}

	accesstags, err := flex.GetGlobalTagsUsingCRN(meta, *share.CRN, "", isAccessTagType)
	if err != nil {
		log.Printf(
			"Error getting shares (%s) access tags: %s", d.Id(), err)
	}

	d.Set(isFileShareTags, tags)
	d.Set(isFileShareAccessTags, accesstags)
	return nil
}

func resourceIbmIsShareUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	updateShareOptions := &vpcv1.UpdateShareOptions{}

	updateShareOptions.SetID(d.Id())

	hasChange := false

	sharePatchModel := &vpcv1.SharePatch{}
	if d.HasChange("name") {
		name := d.Get("name").(string)
		sharePatchModel.Name = &name
		hasChange = true
	}

	if d.HasChange("size") {
		size := int64(d.Get("size").(int))
		sharePatchModel.Size = &size
		hasChange = true
	}

	if hasChange {

		sharePatch, err := sharePatchModel.AsPatch()

		if err != nil {
			log.Printf("[DEBUG] SharePatch AsPatch failed %s", err)
			return diag.FromErr(err)
		}
		updateShareOptions.SetSharePatch(sharePatch)
		_, response, err := vpcClient.UpdateShareWithContext(context, updateShareOptions)
		if err != nil {
			log.Printf("[DEBUG] UpdateShareWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
	}
	if d.HasChange(isFileShareTags) {
		oldList, newList := d.GetChange(isFileShareTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get("crn").(string), "", isUserTagType)
		if err != nil {
			log.Printf(
				"Error updating shares (%s) tags: %s", d.Id(), err)
		}
	}

	if d.HasChange(isFileShareAccessTags) {
		oldList, newList := d.GetChange(isFileShareAccessTags)
		err := flex.UpdateGlobalTagsUsingCRN(oldList, newList, meta, d.Get("crn").(string), "", isAccessTagType)
		if err != nil {
			log.Printf(
				"Error updating shares (%s) access tags: %s", d.Id(), err)
		}
	}
	return resourceIbmIsShareRead(context, d, meta)
}

func resourceIbmIsShareDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareOptions := &vpcv1.GetShareOptions{}

	getShareOptions.SetID(d.Id())

	share, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}
	if share.Targets != nil {
		for _, targetsItem := range share.Targets {

			deleteShareTargetOptions := &vpcv1.DeleteShareTargetOptions{}

			deleteShareTargetOptions.SetShareID(d.Id())
			deleteShareTargetOptions.SetID(*targetsItem.ID)

			_, response, err := vpcClient.DeleteShareTargetWithContext(context, deleteShareTargetOptions)
			if err != nil {
				log.Printf("[DEBUG] DeleteShareTargetWithContext failed %s\n%s", err, response)
				return diag.FromErr(err)
			}
			_, err = isWaitForTargetDelete(context, vpcClient, d, d.Id(), *targetsItem.ID)
			if err != nil {
				return diag.FromErr(err)
			}
		}

	}
	deleteShareOptions := &vpcv1.DeleteShareOptions{}

	deleteShareOptions.SetID(d.Id())

	_, response, err = vpcClient.DeleteShareWithContext(context, deleteShareOptions)
	if err != nil {
		log.Printf("[DEBUG] DeleteShareWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	_, err = isWaitForShareDelete(context, vpcClient, d, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func isWaitForShareAvailable(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, d *schema.ResourceData, timeout time.Duration) (interface{}, error) {
	log.Printf("Waiting for share (%s) to be available.", shareid)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"updating", "pending", "waiting"},
		Target:     []string{"stable", "failed"},
		Refresh:    isShareRefreshFunc(context, vpcClient, shareid, d),
		Timeout:    timeout,
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}

func isShareRefreshFunc(context context.Context, vpcClient *vpcv1.VpcV1, shareid string, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		shareOptions := &vpcv1.GetShareOptions{}

		shareOptions.SetID(shareid)

		share, response, err := vpcClient.GetShareWithContext(context, shareOptions)
		if err != nil {
			return nil, "", fmt.Errorf("Error Getting share: %s\n%s", err, response)
		}
		d.Set("lifecycle_state", *share.LifecycleState)
		if *share.LifecycleState == "stable" || *share.LifecycleState == "failed" {

			return share, *share.LifecycleState, nil

		}
		return share, "pending", nil
	}
}

func isWaitForShareDelete(context context.Context, vpcClient *vpcv1.VpcV1, d *schema.ResourceData, shareid string) (interface{}, error) {

	stateConf := &resource.StateChangeConf{
		Pending: []string{"deleting", "stable", "waiting"},
		Target:  []string{"done"},
		Refresh: func() (interface{}, string, error) {
			shareOptions := &vpcv1.GetShareOptions{}

			shareOptions.SetID(shareid)

			share, response, err := vpcClient.GetShareWithContext(context, shareOptions)
			if err != nil {
				if response != nil && response.StatusCode == 404 {
					return share, "done", nil
				}
				return nil, "", fmt.Errorf("Error Getting Target: %s\n%s", err, response)
			}
			if *share.LifecycleState == isInstanceFailed {
				return share, *share.LifecycleState, fmt.Errorf("The  target %s failed to delete: %v", shareid, err)
			}
			return share, "deleting", nil
		},
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      10 * time.Second,
		MinTimeout: 10 * time.Second,
	}

	return stateConf.WaitForState()
}
