/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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

func dataSourceIbmIsPlacementGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIbmIsPlacementGroupsRead,

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
			"placement_groups": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of placement groups.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the placement group was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this placement group.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this placement group.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this placement group.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the placement group.",
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The user-defined name for this placement group.",
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this placement group.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"href": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The URL for this resource group.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this resource group.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The user-defined name for this resource group.",
									},
								},
							},
						},
						"resource_type": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resource type.",
						},
						"strategy": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The strategy for this placement group- `host_spread`: place on different compute hosts- `power_spread`: place on compute hosts that use different power sourcesThe enumerated values for this property may expand in the future. When processing this property, check for and log unknown values. Optionally halt processing and surface the error, or bypass the placement group on which the unexpected strategy was encountered.",
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

func dataSourceIbmIsPlacementGroupsRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	vpcClient, err := meta.(ClientSession).VpcV1API()
	if err != nil {
		return diag.FromErr(err)
	}

	listPlacementGroupsOptions := &vpcv1.ListPlacementGroupsOptions{}

	placementGroupCollection, response, err := vpcClient.ListPlacementGroupsWithContext(context, listPlacementGroupsOptions)
	if err != nil {
		log.Printf("[DEBUG] ListPlacementGroupsWithContext failed %s\n%s", err, response)
		return diag.FromErr(err)
	}

	d.SetId(dataSourceIbmIsPlacementGroupsID(d))

	if placementGroupCollection.First != nil {
		err = d.Set("first", dataSourcePlacementGroupCollectionFlattenFirst(*placementGroupCollection.First))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting first %s", err))
		}
	}
	if err = d.Set("limit", placementGroupCollection.Limit); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting limit: %s", err))
	}

	if placementGroupCollection.Next != nil {
		err = d.Set("next", dataSourcePlacementGroupCollectionFlattenNext(*placementGroupCollection.Next))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting next %s", err))
		}
	}

	if placementGroupCollection.PlacementGroups != nil {
		err = d.Set("placement_groups", dataSourcePlacementGroupCollectionFlattenPlacementGroups(placementGroupCollection.PlacementGroups))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting placement_groups %s", err))
		}
	}
	if err = d.Set("total_count", placementGroupCollection.TotalCount); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	return nil
}

// dataSourceIbmIsPlacementGroupsID returns a reasonable ID for the list.
func dataSourceIbmIsPlacementGroupsID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourcePlacementGroupCollectionFlattenFirst(result vpcv1.PlacementGroupCollectionFirst) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourcePlacementGroupCollectionFirstToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourcePlacementGroupCollectionFirstToMap(firstItem vpcv1.PlacementGroupCollectionFirst) (firstMap map[string]interface{}) {
	firstMap = map[string]interface{}{}

	if firstItem.Href != nil {
		firstMap["href"] = firstItem.Href
	}

	return firstMap
}

func dataSourcePlacementGroupCollectionFlattenNext(result vpcv1.PlacementGroupCollectionNext) (finalList []map[string]interface{}) {
	finalList = []map[string]interface{}{}
	finalMap := dataSourcePlacementGroupCollectionNextToMap(result)
	finalList = append(finalList, finalMap)

	return finalList
}

func dataSourcePlacementGroupCollectionNextToMap(nextItem vpcv1.PlacementGroupCollectionNext) (nextMap map[string]interface{}) {
	nextMap = map[string]interface{}{}

	if nextItem.Href != nil {
		nextMap["href"] = nextItem.Href
	}

	return nextMap
}

func dataSourcePlacementGroupCollectionFlattenPlacementGroups(result []vpcv1.PlacementGroup) (placementGroups []map[string]interface{}) {
	for _, placementGroupsItem := range result {
		placementGroups = append(placementGroups, dataSourcePlacementGroupCollectionPlacementGroupsToMap(placementGroupsItem))
	}

	return placementGroups
}

func dataSourcePlacementGroupCollectionPlacementGroupsToMap(placementGroupsItem vpcv1.PlacementGroup) (placementGroupsMap map[string]interface{}) {
	placementGroupsMap = map[string]interface{}{}

	if placementGroupsItem.CreatedAt != nil {
		placementGroupsMap["created_at"] = placementGroupsItem.CreatedAt.String()
	}
	if placementGroupsItem.CRN != nil {
		placementGroupsMap["crn"] = placementGroupsItem.CRN
	}
	if placementGroupsItem.Href != nil {
		placementGroupsMap["href"] = placementGroupsItem.Href
	}
	if placementGroupsItem.ID != nil {
		placementGroupsMap["id"] = placementGroupsItem.ID
	}
	if placementGroupsItem.LifecycleState != nil {
		placementGroupsMap["lifecycle_state"] = placementGroupsItem.LifecycleState
	}
	if placementGroupsItem.Name != nil {
		placementGroupsMap["name"] = placementGroupsItem.Name
	}
	if placementGroupsItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourcePlacementGroupCollectionPlacementGroupsResourceGroupToMap(*placementGroupsItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		placementGroupsMap["resource_group"] = resourceGroupList
	}
	if placementGroupsItem.ResourceType != nil {
		placementGroupsMap["resource_type"] = placementGroupsItem.ResourceType
	}
	if placementGroupsItem.Strategy != nil {
		placementGroupsMap["strategy"] = placementGroupsItem.Strategy
	}

	return placementGroupsMap
}

func dataSourcePlacementGroupCollectionPlacementGroupsResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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
