// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceIBMISInstanceGroupManager() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMISInstanceGroupManagerRead,

		Schema: map[string]*schema.Schema{

			"instance_group": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "instance group ID",
			},

			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the instance group manager.",
			},

			"manager_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The type of instance group manager.",
			},

			"aggregation_window": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The time window in seconds to aggregate metrics prior to evaluation",
			},

			"cooldown": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The duration of time in seconds to pause further scale actions after scaling has taken place",
			},

			"max_membership_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The maximum number of members in a managed instance group",
			},

			"min_membership_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The minimum number of members in a managed instance group",
			},

			"manager_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of instance group manager.",
			},

			"policies": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: "list of Policies associated with instancegroup manager",
			},

			"actions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"instance_group_manager_action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_group_manager_action_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMISInstanceGroupManagerRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	sess, err := vpcClient(meta)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("vpcClient creation failed: %s", err.Error()), "ibm_cloud", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	instanceGroupID := d.Get("instance_group").(string)

	// Support for pagination
	start := ""
	allrecs := []vpcv1.InstanceGroupManagerIntf{}

	for {
		listInstanceGroupManagerOptions := vpcv1.ListInstanceGroupManagersOptions{
			InstanceGroupID: &instanceGroupID,
		}
		if start != "" {
			listInstanceGroupManagerOptions.Start = &start
		}
		instanceGroupManagerCollections, response, err := sess.ListInstanceGroupManagersWithContext(context, &listInstanceGroupManagerOptions)
		if err != nil {
			tfErr := flex.TerraformErrorf(err, fmt.Sprintf("ListInstanceGroupManagersWithContext failed: %s\n%s", err.Error(), response), "ibm_cloud", "list")
			log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
			return tfErr.GetDiag()
		}
		start = flex.GetNext(instanceGroupManagerCollections.Next)
		allrecs = append(allrecs, instanceGroupManagerCollections.Managers...)

		if start == "" {
			break
		}
	}

	instanceGroupManagerName := d.Get("name").(string)
	for _, instanceGroupManagerIntf := range allrecs {
		instanceGroupManager := instanceGroupManagerIntf.(*vpcv1.InstanceGroupManager)
		if instanceGroupManagerName == *instanceGroupManager.Name {
			d.SetId(fmt.Sprintf("%s/%s", instanceGroupID, *instanceGroupManager.ID))
			d.Set("manager_type", *instanceGroupManager.ManagerType)
			d.Set("manager_id", *instanceGroupManager.ID)

			if *instanceGroupManager.ManagerType == "scheduled" {

				actions := make([]map[string]interface{}, 0)
				if instanceGroupManager.Actions != nil {
					for _, action := range instanceGroupManager.Actions {
						actn := map[string]interface{}{
							"instance_group_manager_action":      action.ID,
							"instance_group_manager_action_name": action.Name,
							"resource_type":                      action.ResourceType,
						}
						actions = append(actions, actn)
					}
					d.Set("actions", actions)
				}

			} else {
				d.Set("aggregation_window", *instanceGroupManager.AggregationWindow)
				d.Set("cooldown", *instanceGroupManager.Cooldown)
				d.Set("max_membership_count", *instanceGroupManager.MaxMembershipCount)
				d.Set("min_membership_count", *instanceGroupManager.MinMembershipCount)
				policies := make([]string, 0)
				if instanceGroupManager.Policies != nil {
					for i := 0; i < len(instanceGroupManager.Policies); i++ {
						policies = append(policies, string(*(instanceGroupManager.Policies[i].ID)))
					}
				}

				d.Set("policies", policies)
			}

			return nil
		}
	}
	tfErr := flex.TerraformErrorf(err, fmt.Sprintf("No Instance group manager found with name : %s", name), "ibm_cloud", "datasource")
	log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
	return tfErr.GetDiag()
}
