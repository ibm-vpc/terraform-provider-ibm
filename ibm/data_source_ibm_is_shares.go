// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func dataSourceIbmIsShares() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsSharesRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the share.",
			},
			"shares": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of file shares.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this file share.",
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
				},
			},
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsSharesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}
	shareName := ""
	if shareNameIntf, ok := d.GetOk("name"); ok {
		shareName = shareNameIntf.(string)
	}
	start := ""
	allrecs := []vpcv1.Share{}
	totalCount := 0
	for {
		listSharesOptions := &vpcv1.ListSharesOptions{}

		if start != "" {
			listSharesOptions.Start = &start
		}
		if shareName != "" {
			listSharesOptions.Name = &shareName
		}
		shareCollection, response, err := vpcClient.ListSharesWithContext(context, listSharesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListSharesWithContext failed %s\n%s", err, response)
			return diag.FromErr(err)
		}
		if totalCount == 0 {
			totalCount = int(*shareCollection.TotalCount)
		}
		start = GetNext(shareCollection.Next)
		allrecs = append(allrecs, shareCollection.Shares...)
		if start == "" {
			break
		}
	}

	d.SetId(dataSourceIbmIsSharesID(d))

	if allrecs != nil {
		err = d.Set("shares", dataSourceShareCollectionFlattenShares(allrecs))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting shares %s", err))
		}
	}
	if err = d.Set("total_count", totalCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	return nil
}

// dataSourceIbmIsSharesID returns a reasonable ID for the list.
func dataSourceIbmIsSharesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceShareCollectionFlattenFirst(result vpcv1.ShareCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareCollectionFirstToMap(firstItem vpcv1.ShareCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceShareCollectionFlattenNext(result vpcv1.ShareCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareCollectionNextToMap(nextItem vpcv1.ShareCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceShareCollectionFlattenShares(result []vpcv1.Share) (shares []map[string]interface{}) {
	for _, sharesItem := range result {
		shares = append(shares, dataSourceShareCollectionSharesToMap(sharesItem))
	}

	return shares
}

func dataSourceShareCollectionSharesToMap(sharesItem vpcv1.Share) (sharesMap map[string]interface{}) {
	sharesMap = map[string]interface{}{}

	if sharesItem.CreatedAt != nil {
		sharesMap["created_at"] = sharesItem.CreatedAt.String()
	}
	if sharesItem.CRN != nil {
		sharesMap["crn"] = sharesItem.CRN
	}
	if sharesItem.Encryption != nil {
		sharesMap["encryption"] = sharesItem.Encryption
	}
	if sharesItem.EncryptionKey != nil {
		sharesMap["encryption_key"] = *sharesItem.EncryptionKey.CRN
	}
	if sharesItem.Href != nil {
		sharesMap["href"] = sharesItem.Href
	}
	if sharesItem.ID != nil {
		sharesMap["id"] = sharesItem.ID
	}
	if sharesItem.Iops != nil {
		sharesMap["iops"] = sharesItem.Iops
	}
	if sharesItem.LifecycleState != nil {
		sharesMap["lifecycle_state"] = sharesItem.LifecycleState
	}
	if sharesItem.Name != nil {
		sharesMap["name"] = sharesItem.Name
	}
	if sharesItem.Profile != nil {
		profileList := []map[string]interface{}{}
		profileMap := dataSourceShareCollectionSharesProfileToMap(*sharesItem.Profile)
		profileList = append(profileList, profileMap)
		sharesMap["profile"] = profileList
	}
	if sharesItem.ResourceGroup != nil {
		sharesMap["resource_group"] = *sharesItem.ResourceGroup.ID
	}
	if sharesItem.ResourceType != nil {
		sharesMap["resource_type"] = sharesItem.ResourceType
	}
	if sharesItem.Size != nil {
		sharesMap["size"] = sharesItem.Size
	}
	if sharesItem.Targets != nil {
		targetsList := []map[string]interface{}{}
		for _, targetsItem := range sharesItem.Targets {
			targetsList = append(targetsList, dataSourceShareCollectionSharesTargetsToMap(targetsItem))
		}
		sharesMap["targets"] = targetsList
	}
	if sharesItem.Zone != nil {
		sharesMap["zone"] = *sharesItem.Zone.Name
	}

	return sharesMap
}

func dataSourceShareCollectionSharesEncryptionKeyToMap(encryptionKeyItem vpcv1.EncryptionKeyReference) (encryptionKeyMap map[string]interface{}) {
	encryptionKeyMap = map[string]interface{}{}

	if encryptionKeyItem.CRN != nil {
		encryptionKeyMap["crn"] = encryptionKeyItem.CRN
	}

	return encryptionKeyMap
}

func dataSourceShareCollectionSharesProfileToMap(profileItem vpcv1.ShareProfileReference) (profileMap map[string]interface{}) {
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

func dataSourceShareCollectionSharesResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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

func dataSourceShareCollectionSharesTargetsToMap(targetsItem vpcv1.ShareTargetReference) (targetsMap map[string]interface{}) {
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

func dataSourceShareCollectionSharesZoneToMap(zoneItem vpcv1.ZoneReference) (zoneMap map[string]interface{}) {
	zoneMap = map[string]interface{}{}

	if zoneItem.Href != nil {
		zoneMap["href"] = zoneItem.Href
	}
	if zoneItem.Name != nil {
		zoneMap["name"] = zoneItem.Name
	}

	return zoneMap
}
