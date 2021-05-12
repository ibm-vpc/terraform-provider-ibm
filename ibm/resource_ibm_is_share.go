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

func resourceIbmIsShare() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIbmIsShareCreate,
		ReadContext:   resourceIbmIsShareRead,
		UpdateContext: resourceIbmIsShareUpdate,
		DeleteContext: resourceIbmIsShareDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"encryption_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The CRN of the key to use for encrypting this file share.If no encryption key is provided, the share will not be encrypted.",
			},
			"initial_owner_gid": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The owner assigned to the file share at creation. The initial group identifier for the file share. Subsequent changes to the owner must be performed by a virtual server instance that has mounted the file share.",
			},
			"initial_owner_uid": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The owner assigned to the file share at creation. The initial user identifier for the file share. Subsequent changes to the owner must be performed by a virtual server instance that has mounted the file share.",
			},
			"iops": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator("ibm_is_share", "iops"),
				Description:  "The maximum input/output operation performance bandwidth per second for the file share.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: InvokeValidator("ibm_is_share", "name"),
				Description:  "The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"profile": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "tier-3iops",
				ValidateFunc: InvokeValidator("ibm_is_share", "name"),
				Description:  "The globally unique name for this share profile.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the resource group to use. If unspecified, the account's [default resourcegroup](https://cloud.ibm.com/apidocs/resource-manager#introduction) is used.",
			},
			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: InvokeValidator("ibm_is_share", "size"),
				Description:  "The size of the file share rounded up to the next gigabyte.",
			},
			"targets": {
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

func resourceIbmIsShareValidator() *ResourceValidator {
	validateSchema := make([]ValidateSchema, 1)
	validateSchema = append(validateSchema,
		ValidateSchema{
			Identifier:                 "iops",
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Optional:                   true,
			MinValue:                   "100",
			MaxValue:                   "48000",
		},
		ValidateSchema{
			Identifier:                 "name",
			ValidateFunctionIdentifier: ValidateRegexpLen,
			Type:                       TypeString,
			Optional:                   true,
			Regexp:                     `^([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$`,
			MinValueLength:             1,
			MaxValueLength:             63,
		},
		ValidateSchema{
			Identifier:                 "size",
			ValidateFunctionIdentifier: IntBetween,
			Type:                       TypeInt,
			Optional:                   true,
			MinValue:                   "10",
			MaxValue:                   "16000",
		},
	)

	resourceValidator := ResourceValidator{ResourceName: "ibm_is_share", Schema: validateSchema}
	return &resourceValidator
}

func resourceIbmIsShareCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	createShareOptions := &vpcv1.CreateShareOptions{}

	if _, ok := d.GetOk("encryption_key"); ok {
		encryptionKey := resourceIbmIsShareMapToEncryptionKeyIdentity(d.Get("encryption_key.0").(map[string]interface{}))
		createShareOptions.SetEncryptionKey(&encryptionKey)
	}
	initial_owner := &vpcv1.ShareInitialOwner{}
	if o_gid, ok := d.GetOk("initial_owner_gid"); ok {
		o_gidstr := o_gid.(int64)
		initial_owner.Gid = &o_gidstr
		createShareOptions.SetInitialOwner(initial_owner)
	}
	if u_gid, ok := d.GetOk("initial_owner_gid"); ok {
		o_uidstr := u_gid.(int64)
		initial_owner.Uid = &o_uidstr
		createShareOptions.SetInitialOwner(initial_owner)
	}
	if _, ok := d.GetOk("iops"); ok {
		createShareOptions.SetIops(int64(d.Get("iops").(int)))
	}
	if _, ok := d.GetOk("name"); ok {
		createShareOptions.SetName(d.Get("name").(string))
	}
	if profileIntf, ok := d.GetOk("profile"); ok {
		profileStr := profileIntf.(string)
		profile := &vpcv1.SharePrototypeProfile{
			Name: &profileStr,
		}
		createShareOptions.SetProfile(profile)
	}
	if resgrp, ok := d.GetOk("resource_group"); ok {
		resgrpstr := resgrp.(string)
		resourceGroup := &vpcv1.ResourceGroupIdentity{
			ID: &resgrpstr,
		}
		createShareOptions.SetResourceGroup(resourceGroup)
	}
	if _, ok := d.GetOk("size"); ok {
		createShareOptions.SetSize(int64(d.Get("size").(int)))
	}
	if _, ok := d.GetOk("targets"); ok {
		var targets []vpcv1.ShareTargetPrototype
		for _, e := range d.Get("targets").([]interface{}) {
			value := e.(map[string]interface{})
			targetsItem := resourceIbmIsShareMapToShareTargetPrototype(value)
			targets = append(targets, targetsItem)
		}
		createShareOptions.SetTargets(targets)
	}
	if zone, ok := d.GetOk("zone"); ok {
		zonestr := zone.(string)
		zone := &vpcv1.ZoneIdentity{
			Name: &zonestr,
		}
		createShareOptions.SetZone(zone)
	}

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

	return resourceIbmIsShareRead(context, d, meta)
}

func resourceIbmIsShareMapToEncryptionKeyIdentity(encryptionKeyIdentityMap map[string]interface{}) vpcv1.EncryptionKeyIdentity {
	encryptionKeyIdentity := vpcv1.EncryptionKeyIdentity{}

	if encryptionKeyIdentityMap["crn"] != nil {
		encryptionKeyIdentity.CRN = core.StringPtr(encryptionKeyIdentityMap["crn"].(string))
	}

	return encryptionKeyIdentity
}

func resourceIbmIsShareMapToEncryptionKeyIdentityByCRN(encryptionKeyIdentityByCRNMap map[string]interface{}) vpcv1.EncryptionKeyIdentityByCRN {
	encryptionKeyIdentityByCRN := vpcv1.EncryptionKeyIdentityByCRN{}

	encryptionKeyIdentityByCRN.CRN = core.StringPtr(encryptionKeyIdentityByCRNMap["crn"].(string))

	return encryptionKeyIdentityByCRN
}

func resourceIbmIsShareMapToShareInitialOwner(shareInitialOwnerMap map[string]interface{}) vpcv1.ShareInitialOwner {
	shareInitialOwner := vpcv1.ShareInitialOwner{}

	if shareInitialOwnerMap["gid"] != nil {
		shareInitialOwner.Gid = core.Int64Ptr(int64(shareInitialOwnerMap["gid"].(int)))
	}
	if shareInitialOwnerMap["uid"] != nil {
		shareInitialOwner.Uid = core.Int64Ptr(int64(shareInitialOwnerMap["uid"].(int)))
	}

	return shareInitialOwner
}

func resourceIbmIsShareMapToSharePrototypeProfile(sharePrototypeProfileMap map[string]interface{}) vpcv1.SharePrototypeProfile {
	sharePrototypeProfile := vpcv1.SharePrototypeProfile{}

	if sharePrototypeProfileMap["name"] != nil {
		sharePrototypeProfile.Name = core.StringPtr(sharePrototypeProfileMap["name"].(string))
	}
	if sharePrototypeProfileMap["href"] != nil {
		sharePrototypeProfile.Href = core.StringPtr(sharePrototypeProfileMap["href"].(string))
	}

	return sharePrototypeProfile
}

func resourceIbmIsShareMapToSharePrototypeProfileShareProfileIdentityByName(sharePrototypeProfileShareProfileIdentityByNameMap map[string]interface{}) vpcv1.SharePrototypeProfileShareProfileIdentityByName {
	sharePrototypeProfileShareProfileIdentityByName := vpcv1.SharePrototypeProfileShareProfileIdentityByName{}

	sharePrototypeProfileShareProfileIdentityByName.Name = core.StringPtr(sharePrototypeProfileShareProfileIdentityByNameMap["name"].(string))

	return sharePrototypeProfileShareProfileIdentityByName
}

func resourceIbmIsShareMapToSharePrototypeProfileShareProfileIdentityByHref(sharePrototypeProfileShareProfileIdentityByHrefMap map[string]interface{}) vpcv1.SharePrototypeProfileShareProfileIdentityByHref {
	sharePrototypeProfileShareProfileIdentityByHref := vpcv1.SharePrototypeProfileShareProfileIdentityByHref{}

	sharePrototypeProfileShareProfileIdentityByHref.Href = core.StringPtr(sharePrototypeProfileShareProfileIdentityByHrefMap["href"].(string))

	return sharePrototypeProfileShareProfileIdentityByHref
}

func resourceIbmIsShareMapToResourceGroupIdentity(resourceGroupIdentityMap map[string]interface{}) vpcv1.ResourceGroupIdentity {
	resourceGroupIdentity := vpcv1.ResourceGroupIdentity{}

	if resourceGroupIdentityMap["id"] != nil {
		resourceGroupIdentity.ID = core.StringPtr(resourceGroupIdentityMap["id"].(string))
	}

	return resourceGroupIdentity
}

func resourceIbmIsShareMapToResourceGroupIdentityByID(resourceGroupIdentityByIDMap map[string]interface{}) vpcv1.ResourceGroupIdentityByID {
	resourceGroupIdentityByID := vpcv1.ResourceGroupIdentityByID{}

	resourceGroupIdentityByID.ID = core.StringPtr(resourceGroupIdentityByIDMap["id"].(string))

	return resourceGroupIdentityByID
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

func resourceIbmIsShareMapToSubnetIdentity(subnetIdentityMap map[string]interface{}) vpcv1.SubnetIdentity {
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

func resourceIbmIsShareMapToSubnetIdentityByID(subnetIdentityByIDMap map[string]interface{}) vpcv1.SubnetIdentityByID {
	subnetIdentityByID := vpcv1.SubnetIdentityByID{}

	subnetIdentityByID.ID = core.StringPtr(subnetIdentityByIDMap["id"].(string))

	return subnetIdentityByID
}

func resourceIbmIsShareMapToSubnetIdentityByCRN(subnetIdentityByCRNMap map[string]interface{}) vpcv1.SubnetIdentityByCRN {
	subnetIdentityByCRN := vpcv1.SubnetIdentityByCRN{}

	subnetIdentityByCRN.CRN = core.StringPtr(subnetIdentityByCRNMap["crn"].(string))

	return subnetIdentityByCRN
}

func resourceIbmIsShareMapToSubnetIdentityByHref(subnetIdentityByHrefMap map[string]interface{}) vpcv1.SubnetIdentityByHref {
	subnetIdentityByHref := vpcv1.SubnetIdentityByHref{}

	subnetIdentityByHref.Href = core.StringPtr(subnetIdentityByHrefMap["href"].(string))

	return subnetIdentityByHref
}

func resourceIbmIsShareMapToShareTargetPrototypeVpc(shareTargetPrototypeVpcMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPC {
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

func resourceIbmIsShareMapToShareTargetPrototypeVpcVPCIdentityByID(shareTargetPrototypeVpcVPCIdentityByIDMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPCVPCIdentityByID {
	shareTargetPrototypeVpcVPCIdentityByID := vpcv1.ShareTargetPrototypeVPCVPCIdentityByID{}

	shareTargetPrototypeVpcVPCIdentityByID.ID = core.StringPtr(shareTargetPrototypeVpcVPCIdentityByIDMap["id"].(string))

	return shareTargetPrototypeVpcVPCIdentityByID
}

func resourceIbmIsShareMapToShareTargetPrototypeVpcVPCIdentityByCRN(shareTargetPrototypeVpcVPCIdentityByCRNMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPCVPCIdentityByCRN {
	shareTargetPrototypeVpcVPCIdentityByCRN := vpcv1.ShareTargetPrototypeVPCVPCIdentityByCRN{}

	shareTargetPrototypeVpcVPCIdentityByCRN.CRN = core.StringPtr(shareTargetPrototypeVpcVPCIdentityByCRNMap["crn"].(string))

	return shareTargetPrototypeVpcVPCIdentityByCRN
}

func resourceIbmIsShareMapToShareTargetPrototypeVpcVPCIdentityByHref(shareTargetPrototypeVpcVPCIdentityByHrefMap map[string]interface{}) vpcv1.ShareTargetPrototypeVPCVPCIdentityByHref {
	shareTargetPrototypeVpcVPCIdentityByHref := vpcv1.ShareTargetPrototypeVPCVPCIdentityByHref{}

	shareTargetPrototypeVpcVPCIdentityByHref.Href = core.StringPtr(shareTargetPrototypeVpcVPCIdentityByHrefMap["href"].(string))

	return shareTargetPrototypeVpcVPCIdentityByHref
}

func resourceIbmIsShareMapToZoneIdentity(zoneIdentityMap map[string]interface{}) vpcv1.ZoneIdentity {
	zoneIdentity := vpcv1.ZoneIdentity{}

	if zoneIdentityMap["name"] != nil {
		zoneIdentity.Name = core.StringPtr(zoneIdentityMap["name"].(string))
	}
	if zoneIdentityMap["href"] != nil {
		zoneIdentity.Href = core.StringPtr(zoneIdentityMap["href"].(string))
	}

	return zoneIdentity
}

func resourceIbmIsShareMapToZoneIdentityByName(zoneIdentityByNameMap map[string]interface{}) vpcv1.ZoneIdentityByName {
	zoneIdentityByName := vpcv1.ZoneIdentityByName{}

	zoneIdentityByName.Name = core.StringPtr(zoneIdentityByNameMap["name"].(string))

	return zoneIdentityByName
}

func resourceIbmIsShareMapToZoneIdentityByHref(zoneIdentityByHrefMap map[string]interface{}) vpcv1.ZoneIdentityByHref {
	zoneIdentityByHref := vpcv1.ZoneIdentityByHref{}

	zoneIdentityByHref.Href = core.StringPtr(zoneIdentityByHrefMap["href"].(string))

	return zoneIdentityByHref
}

func resourceIbmIsShareRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
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

	if err = d.Set("iops", intValue(share.Iops)); err != nil {
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
	if err = d.Set("size", intValue(share.Size)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting size: %s", err))
	}
	if share.Targets != nil {
		targets := []map[string]interface{}{}
		for _, targetsItem := range share.Targets {
			targetsItemMap := dataSourceShareTargetsToMap(targetsItem)
			targets = append(targets, targetsItemMap)
		}
		if err = d.Set("share_targets", targets); err != nil {
			return diag.FromErr(fmt.Errorf("Error setting targets: %s", err))
		}
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

	return nil
}

func resourceIbmIsShareEncryptionKeyIdentityToMap(encryptionKeyIdentity vpcv1.EncryptionKeyIdentity) map[string]interface{} {
	encryptionKeyIdentityMap := map[string]interface{}{}

	encryptionKeyIdentityMap["crn"] = encryptionKeyIdentity.CRN

	return encryptionKeyIdentityMap
}

func resourceIbmIsShareEncryptionKeyIdentityByCRNToMap(encryptionKeyIdentityByCRN vpcv1.EncryptionKeyIdentityByCRN) map[string]interface{} {
	encryptionKeyIdentityByCRNMap := map[string]interface{}{}

	encryptionKeyIdentityByCRNMap["crn"] = encryptionKeyIdentityByCRN.CRN

	return encryptionKeyIdentityByCRNMap
}

func resourceIbmIsShareShareInitialOwnerToMap(shareInitialOwner vpcv1.ShareInitialOwner) map[string]interface{} {
	shareInitialOwnerMap := map[string]interface{}{}

	shareInitialOwnerMap["gid"] = intValue(shareInitialOwner.Gid)
	shareInitialOwnerMap["uid"] = intValue(shareInitialOwner.Uid)

	return shareInitialOwnerMap
}

func resourceIbmIsShareSharePrototypeProfileToMap(sharePrototypeProfile vpcv1.SharePrototypeProfile) map[string]interface{} {
	sharePrototypeProfileMap := map[string]interface{}{}

	sharePrototypeProfileMap["name"] = sharePrototypeProfile.Name
	sharePrototypeProfileMap["href"] = sharePrototypeProfile.Href

	return sharePrototypeProfileMap
}

func resourceIbmIsShareSharePrototypeProfileShareProfileIdentityByNameToMap(sharePrototypeProfileShareProfileIdentityByName vpcv1.SharePrototypeProfileShareProfileIdentityByName) map[string]interface{} {
	sharePrototypeProfileShareProfileIdentityByNameMap := map[string]interface{}{}

	sharePrototypeProfileShareProfileIdentityByNameMap["name"] = sharePrototypeProfileShareProfileIdentityByName.Name

	return sharePrototypeProfileShareProfileIdentityByNameMap
}

func resourceIbmIsShareSharePrototypeProfileShareProfileIdentityByHrefToMap(sharePrototypeProfileShareProfileIdentityByHref vpcv1.SharePrototypeProfileShareProfileIdentityByHref) map[string]interface{} {
	sharePrototypeProfileShareProfileIdentityByHrefMap := map[string]interface{}{}

	sharePrototypeProfileShareProfileIdentityByHrefMap["href"] = sharePrototypeProfileShareProfileIdentityByHref.Href

	return sharePrototypeProfileShareProfileIdentityByHrefMap
}

func resourceIbmIsShareResourceGroupIdentityToMap(resourceGroupIdentity vpcv1.ResourceGroupIdentity) map[string]interface{} {
	resourceGroupIdentityMap := map[string]interface{}{}

	resourceGroupIdentityMap["id"] = resourceGroupIdentity.ID

	return resourceGroupIdentityMap
}

func resourceIbmIsShareResourceGroupIdentityByIDToMap(resourceGroupIdentityByID vpcv1.ResourceGroupIdentityByID) map[string]interface{} {
	resourceGroupIdentityByIDMap := map[string]interface{}{}

	resourceGroupIdentityByIDMap["id"] = resourceGroupIdentityByID.ID

	return resourceGroupIdentityByIDMap
}

func resourceIbmIsShareShareTargetPrototypeToMap(shareTargetPrototype vpcv1.ShareTargetPrototype) map[string]interface{} {
	shareTargetPrototypeMap := map[string]interface{}{}

	shareTargetPrototypeMap["name"] = shareTargetPrototype.Name
	if shareTargetPrototype.Subnet != nil {
		SubnetMap := resourceIbmIsShareSubnetIdentityToMap(*shareTargetPrototype.Subnet.(*vpcv1.SubnetIdentity))
		shareTargetPrototypeMap["subnet"] = []map[string]interface{}{SubnetMap}
	}
	VpcMap := resourceIbmIsShareShareTargetPrototypeVpcToMap(*shareTargetPrototype.VPC.(*vpcv1.ShareTargetPrototypeVPC))
	shareTargetPrototypeMap["vpc"] = []map[string]interface{}{VpcMap}

	return shareTargetPrototypeMap
}

func resourceIbmIsShareSubnetIdentityToMap(subnetIdentity vpcv1.SubnetIdentity) map[string]interface{} {
	subnetIdentityMap := map[string]interface{}{}

	subnetIdentityMap["id"] = subnetIdentity.ID
	subnetIdentityMap["crn"] = subnetIdentity.CRN
	subnetIdentityMap["href"] = subnetIdentity.Href

	return subnetIdentityMap
}

func resourceIbmIsShareSubnetIdentityByIDToMap(subnetIdentityByID vpcv1.SubnetIdentityByID) map[string]interface{} {
	subnetIdentityByIDMap := map[string]interface{}{}

	subnetIdentityByIDMap["id"] = subnetIdentityByID.ID

	return subnetIdentityByIDMap
}

func resourceIbmIsShareSubnetIdentityByCRNToMap(subnetIdentityByCRN vpcv1.SubnetIdentityByCRN) map[string]interface{} {
	subnetIdentityByCRNMap := map[string]interface{}{}

	subnetIdentityByCRNMap["crn"] = subnetIdentityByCRN.CRN

	return subnetIdentityByCRNMap
}

func resourceIbmIsShareSubnetIdentityByHrefToMap(subnetIdentityByHref vpcv1.SubnetIdentityByHref) map[string]interface{} {
	subnetIdentityByHrefMap := map[string]interface{}{}

	subnetIdentityByHrefMap["href"] = subnetIdentityByHref.Href

	return subnetIdentityByHrefMap
}

func resourceIbmIsShareShareTargetPrototypeVpcToMap(shareTargetPrototypeVpc vpcv1.ShareTargetPrototypeVPC) map[string]interface{} {
	shareTargetPrototypeVpcMap := map[string]interface{}{}

	shareTargetPrototypeVpcMap["id"] = shareTargetPrototypeVpc.ID
	shareTargetPrototypeVpcMap["crn"] = shareTargetPrototypeVpc.CRN
	shareTargetPrototypeVpcMap["href"] = shareTargetPrototypeVpc.Href

	return shareTargetPrototypeVpcMap
}

func resourceIbmIsShareShareTargetPrototypeVpcVPCIdentityByIDToMap(shareTargetPrototypeVpcVPCIdentityByID vpcv1.ShareTargetPrototypeVPCVPCIdentityByID) map[string]interface{} {
	shareTargetPrototypeVpcVPCIdentityByIDMap := map[string]interface{}{}

	shareTargetPrototypeVpcVPCIdentityByIDMap["id"] = shareTargetPrototypeVpcVPCIdentityByID.ID

	return shareTargetPrototypeVpcVPCIdentityByIDMap
}

func resourceIbmIsShareShareTargetPrototypeVpcVPCIdentityByCRNToMap(shareTargetPrototypeVpcVPCIdentityByCRN vpcv1.ShareTargetPrototypeVPCVPCIdentityByCRN) map[string]interface{} {
	shareTargetPrototypeVpcVPCIdentityByCRNMap := map[string]interface{}{}

	shareTargetPrototypeVpcVPCIdentityByCRNMap["crn"] = shareTargetPrototypeVpcVPCIdentityByCRN.CRN

	return shareTargetPrototypeVpcVPCIdentityByCRNMap
}

func resourceIbmIsShareShareTargetPrototypeVpcVPCIdentityByHrefToMap(shareTargetPrototypeVpcVPCIdentityByHref vpcv1.ShareTargetPrototypeVPCVPCIdentityByHref) map[string]interface{} {
	shareTargetPrototypeVpcVPCIdentityByHrefMap := map[string]interface{}{}

	shareTargetPrototypeVpcVPCIdentityByHrefMap["href"] = shareTargetPrototypeVpcVPCIdentityByHref.Href

	return shareTargetPrototypeVpcVPCIdentityByHrefMap
}

func resourceIbmIsShareZoneIdentityToMap(zoneIdentity vpcv1.ZoneIdentity) map[string]interface{} {
	zoneIdentityMap := map[string]interface{}{}

	zoneIdentityMap["name"] = zoneIdentity.Name
	zoneIdentityMap["href"] = zoneIdentity.Href

	return zoneIdentityMap
}

func resourceIbmIsShareZoneIdentityByNameToMap(zoneIdentityByName vpcv1.ZoneIdentityByName) map[string]interface{} {
	zoneIdentityByNameMap := map[string]interface{}{}

	zoneIdentityByNameMap["name"] = zoneIdentityByName.Name

	return zoneIdentityByNameMap
}

func resourceIbmIsShareZoneIdentityByHrefToMap(zoneIdentityByHref vpcv1.ZoneIdentityByHref) map[string]interface{} {
	zoneIdentityByHrefMap := map[string]interface{}{}

	zoneIdentityByHrefMap["href"] = zoneIdentityByHref.Href

	return zoneIdentityByHrefMap
}

func resourceIbmIsShareUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
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
	/*
		if d.HasChange("size") {
			size := int64(d.Get("size").(int))
			log.Println("******* update size", size)
			sharePatchModel.Size = &size
			hasChange = true
		}
	*/
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

	return resourceIbmIsShareRead(context, d, meta)
}

func resourceIbmIsShareDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
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
