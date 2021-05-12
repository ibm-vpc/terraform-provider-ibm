// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIbmIsShare() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsShareRead,

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The file share identifier.",
			},
			"created_at": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time that the file share is created.",
			},
			"crn": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The CRN for this share.",
			},
			"encryption": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of encryption used for this file share.",
			},
			"encryption_key": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key used to encrypt this file share. The CRN of the [Key Protect Root Key](https://cloud.ibm.com/docs/key-protect?topic=key-protect-getting-started-tutorial) or [Hyper Protect Crypto Service Root Key](https://cloud.ibm.com/docs/hs-crypto?topic=hs-crypto-get-started) for this resource.",
			},
			"href": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The URL for this share.",
			},
			"iops": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum input/output operation performance bandwidth per second for the file share.",
			},
			"lifecycle_state": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The lifecycle state of the file share.",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique user-defined name for this file share. If unspecified, the name will be a hyphenated list of randomly-selected words.",
			},
			"profile": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The profile this file share uses.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share profile.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The globally unique name for this share profile.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
					},
				},
			},
			"resource_group": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the resource group for this file share.",
			},
			"resource_type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of resource referenced.",
			},
			"size": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the file share rounded up to the next gigabyte.",
			},
			"targets": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Mount targets for the file share.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"deleted": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "If present, this property indicates the referenced resource has been deleted and providessome supplementary information.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"more_info": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "Link to documentation about deleted resources.",
									},
								},
							},
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this share target.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this share target.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this share target.",
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of resource referenced.",
						},
					},
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The globally unique name of the zone this file share will reside in.",
			},
		},
	}
}

func dataSourceIbmIsShareRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	getShareOptions := &vpcv1.GetShareOptions{}

	getShareOptions.SetID(d.Get("id").(string))

	share, response, err := vpcClient.GetShareWithContext(context, getShareOptions)
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
		err = d.Set("profile", dataSourceShareFlattenProfile(*share.Profile))
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
		err = d.Set("targets", dataSourceShareFlattenTargets(share.Targets))
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

	return nil
}

func dataSourceShareFlattenEncryptionKey(result vpcv1.EncryptionKeyReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareEncryptionKeyToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareEncryptionKeyToMap(encryptionKeyItem vpcv1.EncryptionKeyReference) (encryptionKeyMap map[string]interface{}) {
	encryptionKeyMap = map[string]interface{}{}

	if encryptionKeyItem.CRN != nil {
		encryptionKeyMap["crn"] = encryptionKeyItem.CRN
	}

	return encryptionKeyMap
}

func dataSourceShareFlattenProfile(result vpcv1.ShareProfileReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareProfileToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareProfileToMap(profileItem vpcv1.ShareProfileReference) (profileMap map[string]interface{}) {
	profileMap = map[string]interface{}{}

	if profileItem.Href != nil {
		profileMap["href"] = profileItem.Href
	}
	if profileItem.Name != nil {
		profileMap["name"] = profileItem.Name
	}
	if profileItem.ResourceType != nil {
		profileMap["resource_type"] = profileItem.ResourceType
	}

	return profileMap
}

func dataSourceShareFlattenResourceGroup(result vpcv1.ResourceGroupReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareResourceGroupToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
	resourceGroupMap = map[string]interface{}{}

	if resourceGroupItem.Href != nil {
		resourceGroupMap["href"] = resourceGroupItem.Href
	}
	if resourceGroupItem.ID != nil {
		resourceGroupMap["id"] = resourceGroupItem.ID
	}
	if resourceGroupItem.Name != nil {
		resourceGroupMap["name"] = resourceGroupItem.Name
	}

	return resourceGroupMap
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

func dataSourceShareFlattenZone(result vpcv1.ZoneReference) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareZoneToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareZoneToMap(zoneItem vpcv1.ZoneReference) (zoneMap map[string]interface{}) {
	zoneMap = map[string]interface{}{}

	if zoneItem.Href != nil {
		zoneMap["href"] = zoneItem.Href
	}
	if zoneItem.Name != nil {
		zoneMap["name"] = zoneItem.Name
	}

	return zoneMap
}
