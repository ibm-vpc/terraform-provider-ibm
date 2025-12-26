// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
)

func ResourceIBMCmShareApproval() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIBMCmShareApprovalCreate,
		ReadContext:   resourceIBMCmShareApprovalRead,
		UpdateContext: resourceIBMCmShareApprovalUpdate,
		DeleteContext: resourceIBMCmShareApprovalDelete,
		Importer:      &schema.ResourceImporter{},

		Schema: map[string]*schema.Schema{
			"object_kind": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"offering", "vpe", "proxy_source", "preset_configuration"}, false),
				Description:  "The object kind for share approval. Options are offering, vpe, proxy_source, or preset_configuration.",
			},
			"approval_state": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"approved", "pending", "rejected"}, false),
				Description:  "The approval state. Options are approved, pending, or rejected.",
			},
			"account_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of account IDs to set approval for. Prefix with -acct- for regular accounts, -ent- for enterprise accounts, and -entgrp- for enterprise account groups.",
			},
			"enterprise_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Enterprise or enterprise account group ID to view or manage requests for the enterprise. Prefix with -ent- for an enterprise and -entgrp- for an account group.",
			},
		},
	}
}

func resourceIBMCmShareApprovalCreate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_share_approval", "create", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	objectKind := d.Get("object_kind").(string)
	approvalState := d.Get("approval_state").(string)
	accountIDs := flex.ExpandStringList(d.Get("account_ids").([]interface{}))

	updateShareApprovalListAsSourceOptions := catalogManagementClient.NewUpdateShareApprovalListAsSourceOptions(
		objectKind,
		approvalState,
		accountIDs,
	)

	if enterpriseID, ok := d.GetOk("enterprise_id"); ok {
		updateShareApprovalListAsSourceOptions.SetEnterpriseID(enterpriseID.(string))
	}

	accessListBulkResponse, response, err := catalogManagementClient.UpdateShareApprovalListAsSourceWithContext(context, updateShareApprovalListAsSourceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateShareApprovalListAsSourceWithContext failed: %s\n%s", err.Error(), response), "ibm_cm_share_approval", "create")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Create a composite ID from object_kind and approval_state
	d.SetId(fmt.Sprintf("%s/%s", objectKind, approvalState))

	// Log the response for debugging
	if accessListBulkResponse != nil {
		log.Printf("[DEBUG] Share approval update response: %+v", accessListBulkResponse)
	}

	return resourceIBMCmShareApprovalRead(context, d, meta)
}

func resourceIBMCmShareApprovalRead(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_share_approval", "read", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Parse the ID to get object_kind and approval_state
	idParts := strings.Split(d.Id(), "/")
	if len(idParts) != 2 {
		return diag.Errorf("Invalid ID format. Expected format: object_kind/approval_state")
	}

	objectKind := idParts[0]
	approvalState := idParts[1]

	getShareApprovalListAsSourceOptions := catalogManagementClient.NewGetShareApprovalListAsSourceOptions(
		objectKind,
		approvalState,
	)

	if enterpriseID, ok := d.GetOk("enterprise_id"); ok {
		getShareApprovalListAsSourceOptions.SetEnterpriseID(enterpriseID.(string))
	}

	shareApprovalListAccessResult, response, err := catalogManagementClient.GetShareApprovalListAsSourceWithContext(context, getShareApprovalListAsSourceOptions)
	if err != nil {
		if response != nil && response.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("GetShareApprovalListAsSourceWithContext failed: %s", err.Error()), "ibm_cm_share_approval", "read")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	if err = d.Set("object_kind", objectKind); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting object_kind: %s", err), "ibm_cm_share_approval", "read", "set-object_kind").GetDiag()
	}
	if err = d.Set("approval_state", approvalState); err != nil {
		return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting approval_state: %s", err), "ibm_cm_share_approval", "read", "set-approval_state").GetDiag()
	}

	// Extract account IDs from the response
	if shareApprovalListAccessResult != nil && shareApprovalListAccessResult.Resources != nil {
		accountIDs := make([]string, 0)
		for _, resource := range shareApprovalListAccessResult.Resources {
			if resource.Account != nil {
				accountIDs = append(accountIDs, *resource.Account)
			}
		}
		if err = d.Set("account_ids", accountIDs); err != nil {
			return flex.DiscriminatedTerraformErrorf(err, fmt.Sprintf("Error setting account_ids: %s", err), "ibm_cm_share_approval", "read", "set-account_ids").GetDiag()
		}
	}

	return nil
}

func resourceIBMCmShareApprovalUpdate(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_share_approval", "update", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	objectKind := d.Get("object_kind").(string)
	approvalState := d.Get("approval_state").(string)
	accountIDs := flex.ExpandStringList(d.Get("account_ids").([]interface{}))

	updateShareApprovalListAsSourceOptions := catalogManagementClient.NewUpdateShareApprovalListAsSourceOptions(
		objectKind,
		approvalState,
		accountIDs,
	)

	if enterpriseID, ok := d.GetOk("enterprise_id"); ok {
		updateShareApprovalListAsSourceOptions.SetEnterpriseID(enterpriseID.(string))
	}

	_, response, err := catalogManagementClient.UpdateShareApprovalListAsSourceWithContext(context, updateShareApprovalListAsSourceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateShareApprovalListAsSourceWithContext failed: %s\n%s", err.Error(), response), "ibm_cm_share_approval", "update")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// Update the ID if approval_state changed
	if d.HasChange("approval_state") {
		d.SetId(fmt.Sprintf("%s/%s", objectKind, approvalState))
	}

	return resourceIBMCmShareApprovalRead(context, d, meta)
}

func resourceIBMCmShareApprovalDelete(context context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	catalogManagementClient, err := meta.(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		tfErr := flex.DiscriminatedTerraformErrorf(err, err.Error(), "ibm_cm_share_approval", "delete", "initialize-client")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	// On delete, set the approval state to "rejected" to revoke access
	objectKind := d.Get("object_kind").(string)
	accountIDs := flex.ExpandStringList(d.Get("account_ids").([]interface{}))

	updateShareApprovalListAsSourceOptions := catalogManagementClient.NewUpdateShareApprovalListAsSourceOptions(
		objectKind,
		"rejected",
		accountIDs,
	)

	if enterpriseID, ok := d.GetOk("enterprise_id"); ok {
		updateShareApprovalListAsSourceOptions.SetEnterpriseID(enterpriseID.(string))
	}

	_, response, err := catalogManagementClient.UpdateShareApprovalListAsSourceWithContext(context, updateShareApprovalListAsSourceOptions)
	if err != nil {
		tfErr := flex.TerraformErrorf(err, fmt.Sprintf("UpdateShareApprovalListAsSourceWithContext failed during delete: %s\n%s", err.Error(), response), "ibm_cm_share_approval", "delete")
		log.Printf("[DEBUG]\n%s", tfErr.GetDebugMessage())
		return tfErr.GetDiag()
	}

	d.SetId("")
	return nil
}
