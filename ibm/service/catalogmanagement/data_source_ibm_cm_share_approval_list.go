// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func DataSourceIBMCmShareApprovalList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIBMCmShareApprovalListRead,

		Schema: map[string]*schema.Schema{
			"object_kind": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"offering", "vpe", "proxy_source", "preset_configuration"}, false),
				Description:  "The object kind for share approval. Options are offering, vpe, proxy_source, or preset_configuration.",
			},
			"approval_state": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"approved", "pending", "rejected"}, false),
				Description:  "The approval state to filter by. Options are approved, pending, or rejected.",
			},
			"enterprise_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enterprise or enterprise account group ID to view or manage requests for the enterprise. Prefix with -ent- for an enterprise and -entgrp- for an account group.",
			},
			"limit": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     100,
				Description: "Number of results to return in the query.",
			},
			"total_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The total count of resources that match the request.",
			},
			"resource_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The number of resources returned in this response.",
			},
			"resources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A list of share approval access records.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Unique ID of the access record.",
						},
						"account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Account ID.",
						},
						"account_type": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Account type (normal account or enterprise).",
						},
						"target_account": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Object's owner's account.",
						},
						"target_kind": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Entity type.",
						},
						"created": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Date and time the access record was created.",
						},
						"approval_state": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Approval state of the access record.",
						},
					},
				},
			},
		},
	}
}

func dataSourceIBMCmShareApprovalListRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cm_share_approval_list", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	objectKind := d.Get("object_kind").(string)
	approvalState := d.Get("approval_state").(string)

	getShareApprovalListAsSourceOptions := catalogManagementClient.NewGetShareApprovalListAsSourceOptions(
		objectKind,
		approvalState,
	)

	if enterpriseID, ok := d.GetOk("enterprise_id"); ok {
		getShareApprovalListAsSourceOptions.SetEnterpriseID(enterpriseID.(string))
	}

	if limit, ok := d.GetOk("limit"); ok {
		getShareApprovalListAsSourceOptions.SetLimit(int64(limit.(int)))
	}

	shareApprovalListAccessResult, response, err := catalogManagementClient.GetShareApprovalListAsSourceWithContext(context, getShareApprovalListAsSourceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetShareApprovalListAsSourceWithContext failed: %s\n%s", err.Error(), response), "(Data) ibm_cm_share_approval_list", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Generate a unique ID based on the query parameters and timestamp
	d.SetId(dataSourceIBMCmShareApprovalListID(d))

	if !core.IsNil(shareApprovalListAccessResult.TotalCount) {
		if err = d.Set("total_count", flex.IntValue(shareApprovalListAccessResult.TotalCount)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting total_count: %s", err), "(Data) ibm_cm_share_approval_list", "read", "set-total_count").GetDiag()
		}
	}

	if !core.IsNil(shareApprovalListAccessResult.ResourceCount) {
		if err = d.Set("resource_count", flex.IntValue(shareApprovalListAccessResult.ResourceCount)); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resource_count: %s", err), "(Data) ibm_cm_share_approval_list", "read", "set-resource_count").GetDiag()
		}
	}

	if shareApprovalListAccessResult.Resources != nil {
		resources := []map[string]interface{}{}
		for _, resourceItem := range shareApprovalListAccessResult.Resources {
			resourceMap, err := dataSourceIBMCmShareApprovalListShareApprovalAccessToMap(&resourceItem)
			if err != nil {
				return flex.DiscriminatedTerraformErrorf(err, err.Error(), "(Data) ibm_cm_share_approval_list", "read", "resources-to-map").GetDiag()
			}
			resources = append(resources, resourceMap)
		}
		if err = d.Set("resources", resources); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting resources: %s", err), "(Data) ibm_cm_share_approval_list", "read", "set-resources").GetDiag()
		}
	}

	return nil
}

// dataSourceIBMCmShareApprovalListID returns a reasonable ID for the data source.
func dataSourceIBMCmShareApprovalListID(d *schema.ResourceData) string {
	return time.Now().UTC().String()
}

func dataSourceIBMCmShareApprovalListShareApprovalAccessToMap(model *catalogmanagementv1.ShareApprovalAccess) (map[string]interface{}, error) {
	modelMap := make(map[string]interface{})
	if model.ID != nil {
		modelMap["id"] = *model.ID
	}
	if model.Account != nil {
		modelMap["account"] = *model.Account
	}
	if model.AccountType != nil {
		modelMap["account_type"] = flex.IntValue(model.AccountType)
	}
	if model.TargetAccount != nil {
		modelMap["target_account"] = *model.TargetAccount
	}
	if model.TargetKind != nil {
		modelMap["target_kind"] = *model.TargetKind
	}
	if model.Created != nil {
		modelMap["created"] = model.Created.String()
	}
	if model.ApprovalState != nil {
		modelMap["approval_state"] = *model.ApprovalState
	}
	return modelMap, nil
}
