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

func dataSourceIBMIsBackupPolicies() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMIsBackupPoliciesRead,

		Schema: map[string]*schema.Schema{

			"backup_policies": &schema.Schema{
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Collection of backup policies.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"created_at": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The date and time that the backup policy was created.",
						},
						"crn": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The CRN for this backup policy.",
						},
						"href": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The URL for this backup policy.",
						},
						"id": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique identifier for this backup policy.",
						},
						"lifecycle_state": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The lifecycle state of the backup policy.",
						},
						"match_resource_types": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "A resource type this backup policy applies to. Resources that have both a matching type and a matching user tag will be subject to the backup policy.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"match_user_tags": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The user tags this backup policy applies to. Resources that have both a matching user tag and a matching type will be subject to the backup policy.",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"name": &schema.Schema{
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique user-defined name for this backup policy.",
						},
						"plans": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The plans for the backup policy.",
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
										Description: "The URL for this backup policy plan.",
									},
									"id": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique identifier for this backup policy plan.",
									},
									"name": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The unique user-defined name for this backup policy plan.",
									},
									"resource_type": &schema.Schema{
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The type of resource referenced.",
									},
								},
							},
						},
						"resource_group": &schema.Schema{
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The resource group for this backup policy.",
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
							Description: "The type of resource referenced.",
						},
					},
				},
			},
			// "first": &schema.Schema{
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "A link to the first page of resources.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"href": &schema.Schema{
			// 				Type:        schema.TypeString,
			// 				Computed:    true,
			// 				Description: "The URL for a page of resources.",
			// 			},
			// 		},
			// 	},
			// },
			// "limit": &schema.Schema{
			// 	Type:        schema.TypeInt,
			// 	Computed:    true,
			// 	Description: "The maximum number of resources that can be returned by the request.",
			// },
			// "next": &schema.Schema{
			// 	Type:        schema.TypeList,
			// 	Computed:    true,
			// 	Description: "A link to the next page of resources. This property is present for all pagesexcept the last page.",
			// 	Elem: &schema.Resource{
			// 		Schema: map[string]*schema.Schema{
			// 			"href": &schema.Schema{
			// 				Type:        schema.TypeString,
			// 				Computed:    true,
			// 				Description: "The URL for a page of resources.",
			// 			},
			// 		},
			// 	},
			// },
			"total_count": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total number of resources across all pages.",
			},
		},
	}
}

func dataSourceIBMIsBackupPoliciesRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		return diag.FromErr(err)
	}

	listBackupPoliciesOptions := &vpcv1.ListBackupPoliciesOptions{}

	backupPolicyCollection, response, err := sess.ListBackupPoliciesWithContext(context, listBackupPoliciesOptions)
	if err != nil {
		log.Printf("[DEBUG] ListBackupPoliciesWithContext failed %s\n%s", err, response)
		return diag.FromErr(fmt.Errorf("ListBackupPoliciesWithContext failed %s\n%s", err, response))
	}

	// name := v.(string)
	start := ""
	matchBackupPolicies := []vpcv1.BackupPolicy{}
	for {
		listBackupPoliciesOptions := &vpcv1.ListBackupPoliciesOptions{}
		if start != "" {
			listBackupPoliciesOptions.Start = &start
		}
		backupPolicyCollection, response, err := sess.ListBackupPoliciesWithContext(context, listBackupPoliciesOptions)
		if err != nil {
			log.Printf("[DEBUG] ListBackupPoliciesWithContext failed %s\n%s", err, response)
			return diag.FromErr(fmt.Errorf("ListBackupPoliciesWithContext failed %s\n%s", err, response))
		}
		if backupPolicyCollection != nil && *backupPolicyCollection.TotalCount == int64(0) {
			break
		}
		start = GetNext(backupPolicyCollection.Next)
		matchBackupPolicies = append(matchBackupPolicies, backupPolicyCollection.BackupPolicies...)
		if start == "" {
			break
		}

	}
	// if len(matchBackupPolicies) == 0 {
	// 	return diag.FromErr(fmt.Errorf("no BackupPolicies found"))
	// }

	// Use the provided filter argument and construct a new list with only the requested resource(s)
	// var matchBackupPolicies []vpcv1.BackupPolicy
	// var name string
	// var suppliedFilter bool

	// if v, ok := d.GetOk("name"); ok {
	// 	name = v.(string)
	// 	suppliedFilter = true
	// 	for _, data := range backupPolicyCollection.BackupPolicies {
	// 		if *data.Name == name {
	// 			matchBackupPolicies = append(matchBackupPolicies, data)
	// 		}
	// 	}
	// } else {
	// 	matchBackupPolicies = backupPolicyCollection.BackupPolicies
	// }
	// backupPolicyCollection.BackupPolicies = matchBackupPolicies

	// if suppliedFilter {
	// 	if len(backupPolicyCollection.BackupPolicies) == 0 {
	// 		return diag.FromErr(fmt.Errorf("no BackupPolicies found with name %s", name))
	// 	}
	// 	d.SetId(name)
	// } else {
	d.SetId(dataSourceIBMIsBackupPoliciesID(d))
	// }

	if matchBackupPolicies != nil {
		err = d.Set("backup_policies", dataSourceBackupPolicyCollectionFlattenBackupPolicies(matchBackupPolicies))
		if err != nil {
			return diag.FromErr(fmt.Errorf("Error setting backup_policies %s", err))
		}
	}

	// if backupPolicyCollection.First != nil {
	// 	err = d.Set("first", dataSourceBackupPolicyCollectionFlattenFirst(*backupPolicyCollection.First))
	// 	if err != nil {
	// 		return diag.FromErr(fmt.Errorf("Error setting first %s", err))
	// 	}
	// }
	// if err = d.Set("limit", intValue(backupPolicyCollection.Limit)); err != nil {
	// 	return diag.FromErr(fmt.Errorf("Error setting limit: %s", err))
	// }

	// if backupPolicyCollection.Next != nil {
	// 	err = d.Set("next", dataSourceBackupPolicyCollectionFlattenNext(*backupPolicyCollection.Next))
	// 	if err != nil {
	// 		return diag.FromErr(fmt.Errorf("Error setting next %s", err))
	// 	}
	// }
	if err = d.Set("total_count", intValue(backupPolicyCollection.TotalCount)); err != nil {
		return diag.FromErr(fmt.Errorf("Error setting total_count: %s", err))
	}

	return nil
}

// dataSourceIBMIsBackupPoliciesID returns a reasonable ID for the list.
func dataSourceIBMIsBackupPoliciesID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceBackupPolicyCollectionFlattenBackupPolicies(result []vpcv1.BackupPolicy) (backupPolicies []map[string]interface{}) {
	for _, backupPoliciesItem := range result {
		backupPolicies = append(backupPolicies, dataSourceBackupPolicyCollectionBackupPoliciesToMap(backupPoliciesItem))
	}

	return backupPolicies
}

func dataSourceBackupPolicyCollectionBackupPoliciesToMap(backupPoliciesItem vpcv1.BackupPolicy) (backupPoliciesMap map[string]interface{}) {
	backupPoliciesMap = map[string]interface{}{}

	if backupPoliciesItem.CreatedAt != nil {
		backupPoliciesMap["created_at"] = backupPoliciesItem.CreatedAt.String()
	}
	if backupPoliciesItem.CRN != nil {
		backupPoliciesMap["crn"] = backupPoliciesItem.CRN
	}
	if backupPoliciesItem.Href != nil {
		backupPoliciesMap["href"] = backupPoliciesItem.Href
	}
	if backupPoliciesItem.ID != nil {
		backupPoliciesMap["id"] = backupPoliciesItem.ID
	}
	if backupPoliciesItem.LifecycleState != nil {
		backupPoliciesMap["lifecycle_state"] = backupPoliciesItem.LifecycleState
	}
	if backupPoliciesItem.MatchResourceTypes != nil {
		backupPoliciesMap["match_resource_types"] = backupPoliciesItem.MatchResourceTypes
	}
	if backupPoliciesItem.MatchUserTags != nil {
		backupPoliciesMap["match_user_tags"] = backupPoliciesItem.MatchUserTags
	}
	if backupPoliciesItem.Name != nil {
		backupPoliciesMap["name"] = backupPoliciesItem.Name
	}
	if backupPoliciesItem.Plans != nil {
		plansList := []map[string]interface{}{}
		for _, plansItem := range backupPoliciesItem.Plans {
			plansList = append(plansList, dataSourceBackupPolicyCollectionBackupPoliciesPlansToMap(plansItem))
		}
		backupPoliciesMap["plans"] = plansList
	}
	if backupPoliciesItem.ResourceGroup != nil {
		resourceGroupList := []map[string]interface{}{}
		resourceGroupMap := dataSourceBackupPolicyCollectionBackupPoliciesResourceGroupToMap(*backupPoliciesItem.ResourceGroup)
		resourceGroupList = append(resourceGroupList, resourceGroupMap)
		backupPoliciesMap["resource_group"] = resourceGroupList
	}
	if backupPoliciesItem.ResourceType != nil {
		backupPoliciesMap["resource_type"] = backupPoliciesItem.ResourceType
	}

	return backupPoliciesMap
}

func dataSourceBackupPolicyCollectionBackupPoliciesPlansToMap(plansItem vpcv1.BackupPolicyPlanReference) (plansMap map[string]interface{}) {
	plansMap = map[string]interface{}{}

	if plansItem.Deleted != nil {
		deletedList := []map[string]interface{}{}
		deletedMap := dataSourceBackupPolicyCollectionPlansDeletedToMap(*plansItem.Deleted)
		deletedList = append(deletedList, deletedMap)
		plansMap["deleted"] = deletedList
	}
	if plansItem.Href != nil {
		plansMap["href"] = plansItem.Href
	}
	if plansItem.ID != nil {
		plansMap["id"] = plansItem.ID
	}
	if plansItem.Name != nil {
		plansMap["name"] = plansItem.Name
	}
	if plansItem.ResourceType != nil {
		plansMap["resource_type"] = plansItem.ResourceType
	}

	return plansMap
}

func dataSourceBackupPolicyCollectionPlansDeletedToMap(deletedItem vpcv1.BackupPolicyPlanReferenceDeleted) (deletedMap map[string]interface{}) {
	deletedMap = map[string]interface{}{}

	if deletedItem.MoreInfo != nil {
		deletedMap["more_info"] = deletedItem.MoreInfo
	}

	return deletedMap
}

func dataSourceBackupPolicyCollectionBackupPoliciesResourceGroupToMap(resourceGroupItem vpcv1.ResourceGroupReference) (resourceGroupMap map[string]interface{}) {
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

// func dataSourceBackupPolicyCollectionFlattenFirst(result vpcv1.BackupPolicyCollectionFirst) (finalList []map[string]interface{}) {
// 	finalList = []map[string]interface{}{}
// 	finalMap := dataSourceBackupPolicyCollectionFirstToMap(result)
// 	finalList = append(finalList, finalMap)

// 	return finalList
// }

// func dataSourceBackupPolicyCollectionFirstToMap(firstItem vpcv1.BackupPolicyCollectionFirst) (firstMap map[string]interface{}) {
// 	firstMap = map[string]interface{}{}

// 	if firstItem.Href != nil {
// 		firstMap["href"] = firstItem.Href
// 	}

// 	return firstMap
// }

// func dataSourceBackupPolicyCollectionFlattenNext(result vpcv1.BackupPolicyCollectionNext) (finalList []map[string]interface{}) {
// 	finalList = []map[string]interface{}{}
// 	finalMap := dataSourceBackupPolicyCollectionNextToMap(result)
// 	finalList = append(finalList, finalMap)

// 	return finalList
// }

// func dataSourceBackupPolicyCollectionNextToMap(nextItem vpcv1.BackupPolicyCollectionNext) (nextMap map[string]interface{}) {
// 	nextMap = map[string]interface{}{}

// 	if nextItem.Href != nil {
// 		nextMap["href"] = nextItem.Href
// 	}

// 	return nextMap
// }
