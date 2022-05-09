// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func DataSourceIbmIsShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsShareRead,

		Schema: map[string]*schema.Schema{
			"share": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"name", "share"},
				Description:  "The file share identifier.",
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"share", "name"},
				Description:  "Name of the share.",
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
			"encryption_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key used to encrypt this file share. The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
			},
			"href": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share.",
			},
			"iops": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum input/output operation performance bandwidth per second for the file share.",
			},
			"lifecycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the file share.",
			},
			"profile": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the profile this file share uses.",
			},
			"resource_group": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the resource group for this file share.",
			},
			"resource_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the file share rounded up to the next gigabyte.",
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
			"zone": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the zone this file share will reside in.",
			},
			isFileShareAccessTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of access management tags",
			},
			isFileShareTags: {
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         flex.ResourceIBMVPCHash,
				Description: "List of tags",
			},
		},
	}
}

func dataSourceIbmIsShareRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(conns.ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	shareName := d.Get("name").(string)
	shareId := d.Get("share").(string)
	var share *vpcv1.Share = nil
	if shareId != "" {
		getShareOptions := &vpcv1.GetShareOptions{}

		getShareOptions.SetID(d.Get("share").(string))

		shareItem, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
		if err != nil {
			if response != nil {
				if response.StatusCode == 404 {
					d.SetId("")
				}
				log.Printf("[DEBUG] GetShareWithContext failed %s\n%s", err, response)
				return nil
			}
			log.Printf("[DEBUG] GetShareWithContext failed %s\n", err)
			return diag.FromErr(err)
		}
		share = shareItem
	} else if shareName != "" {
		listSharesOptions := &vpcv1.ListSharesOptions{}

		if shareName != "" {
			listSharesOptions.Name = &shareName
		}
		shareCollection, response, err := vpcClient.ListSharesWithContext(context, listSharesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListSharesWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		for _, sharesItem := range shareCollection.Shares {
			if *sharesItem.Name == shareName {
				share = &sharesItem
				break
			}
		}
		if share == nil {
			return diag.FromErr(fmt.Errorf("[ERROR] Share with provided name %s not found", shareName))
		}
	}

	d.SetId(*share.ID)
	if err = d.Set("created_at", share.CreatedAt.String()); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting created_at: %s", err))
	}
	if err = d.Set("crn", share.CRN); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting crn: %s", err))
	}
	if err = d.Set("encryption", share.Encryption); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting encryption: %s", err))
	}

	if share.EncryptionKey != nil {
		err = d.Set("encryption_key", *share.EncryptionKey.CRN)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting encryption_key %s", err))
		}
	}
	if err = d.Set("href", share.Href); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting href: %s", err))
	}
	if err = d.Set("iops", share.Iops); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting iops: %s", err))
	}
	if err = d.Set("lifecycle_state", share.LifecycleState); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting lifecycle_state: %s", err))
	}
	if err = d.Set("name", share.Name); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting name: %s", err))
	}

	if share.Profile != nil {
		err = d.Set("profile", *share.Profile.Name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting profile %s", err))
		}
	}

	if share.ResourceGroup != nil {
		err = d.Set("resource_group", *share.ResourceGroup.ID)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting resource_group %s", err))
		}
	}
	if err = d.Set("resource_type", share.ResourceType); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting resource_type: %s", err))
	}
	if err = d.Set("size", share.Size); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting size: %s", err))
	}

	if share.Targets != nil {
		err = d.Set("share_targets", dataSourceShareFlattenTargets(share.Targets))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting targets %s", err))
		}
	}

	if share.Zone != nil {
		err = d.Set("zone", *share.Zone.Name)
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting zone %s", err))
		}
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

func dataSourceShareFlattenTargets(result []vpcv1.ShareTargetReference) (targets []map[string]interface{}) {
	for _, targetsItem := range result {
		targets = append(targets, dataSourceShareTargetsToMap(targetsItem))
	}

	return targets
}

func dataSourceShareTargetsToMap(targetsItem vpcv1.ShareTargetReference) (targetsMap map[string]interface{}) {
	targetsMap = map[string]interface{}{}

	if targetsItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceShareTargetsDeletedToMap(*targetsItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		targetsMap["deleted"] = deletedList
	}
	if targetsItem.Href != nil {
		targetsMap["href"] = targetsItem.Href
	}
	if targetsItem.ID != nil {
		targetsMap["id"] = targetsItem.ID
	}
	if targetsItem.Name != nil {
		targetsMap["name"] = targetsItem.Name
	}
	if targetsItem.ResourceType != nil {
		targetsMap["resource_type"] = targetsItem.ResourceType
	}

	return targetsMap
}

func dataSourceShareTargetsDeletedToMap(deletedItem vpcv1.ShareTargetReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}
