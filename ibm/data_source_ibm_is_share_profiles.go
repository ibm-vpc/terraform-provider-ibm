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

func dataSourceIbmIsShareProfiles() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsShareProfilesRead,

		Schema: map[string]*schema.Schema{
			"first": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the first page of resources.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"limit": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of resources that can be returned by the request.",
			},
			"next": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A link to the next page of resources. This property is present for all pagesexcept the last page.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for a page of resources.",
						},
					},
				},
			},
			"profiles": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of share profiles.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"family": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The product family this share profile belongs to.",
						},
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
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIbmIsShareProfilesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listShareProfilesOptions := &vpcv1.ListShareProfilesOptions{}

	shareProfileCollection, response, err := vpcClient.ListShareProfilesWithContext(context, listShareProfilesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListShareProfilesWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(dataSourceIbmIsShareProfilesID(d))

	if shareProfileCollection.First != nil {
		err = d.Set("first", dataSourceShareProfileCollectionFlattenFirst(*shareProfileCollection.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting first %s", err))
		}
	}
	if err = d.Set("limit", shareProfileCollection.Limit); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting limit: %s", err))
	}

	if shareProfileCollection.Next != nil {
		err = d.Set("next", dataSourceShareProfileCollectionFlattenNext(*shareProfileCollection.Next))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting next %s", err))
		}
	}

	if shareProfileCollection.Profiles != nil {
		err = d.Set("profiles", dataSourceShareProfileCollectionFlattenProfiles(shareProfileCollection.Profiles))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting profiles %s", err))
		}
	}
	if err = d.Set("total_count", shareProfileCollection.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	return nil
}

// dataSourceIbmIsShareProfilesID returns a reasonable ID for the list.
func dataSourceIbmIsShareProfilesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceShareProfileCollectionFlattenFirst(result vpcv1.ShareProfileCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareProfileCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareProfileCollectionFirstToMap(firstItem vpcv1.ShareProfileCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourceShareProfileCollectionFlattenNext(result vpcv1.ShareProfileCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourceShareProfileCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourceShareProfileCollectionNextToMap(nextItem vpcv1.ShareProfileCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourceShareProfileCollectionFlattenProfiles(result []vpcv1.ShareProfile) (profiles []map[string]interface{}) {
	for _, profilesItem := range result {
		profiles = append(profiles, dataSourceShareProfileCollectionProfilesToMap(profilesItem))
	}

	return profiles
}

func dataSourceShareProfileCollectionProfilesToMap(profilesItem vpcv1.ShareProfile) (profilesMap map[string]interface{}) {
	profilesMap = map[string]interface{}{}

	if profilesItem.Family != nil {
		profilesMap["family"] = profilesItem.Family
	}
	if profilesItem.Href != nil {
		profilesMap["href"] = profilesItem.Href
	}
	if profilesItem.Name != nil {
		profilesMap["name"] = profilesItem.Name
	}
	if profilesItem.ResourceType != nil {
		profilesMap["resource_type"] = profilesItem.ResourceType
	}

	return profilesMap
}
