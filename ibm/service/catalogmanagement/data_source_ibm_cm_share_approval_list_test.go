// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMCmShareApprovalListDataSourceBasic(t *testing.T) {
	objectKind := "offering"
	approvalState := "approved"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCmShareApprovalListDataSourceConfigBasic(objectKind, approvalState),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_cm_share_approval_list.cm_share_approval_list", "id"),
					resource.TestCheckResourceAttr("data.ibm_cm_share_approval_list.cm_share_approval_list", "object_kind", objectKind),
					resource.TestCheckResourceAttr("data.ibm_cm_share_approval_list.cm_share_approval_list", "approval_state", approvalState),
					resource.TestCheckResourceAttrSet("data.ibm_cm_share_approval_list.cm_share_approval_list", "total_count"),
					resource.TestCheckResourceAttrSet("data.ibm_cm_share_approval_list.cm_share_approval_list", "resource_count"),
				),
			},
		},
	})
}

func TestAccIBMCmShareApprovalListDataSourceAllArgs(t *testing.T) {
	objectKind := "offering"
	approvalState := "pending"
	limit := "50"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCmShareApprovalListDataSourceConfig(objectKind, approvalState, limit),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_cm_share_approval_list.cm_share_approval_list", "id"),
					resource.TestCheckResourceAttr("data.ibm_cm_share_approval_list.cm_share_approval_list", "object_kind", objectKind),
					resource.TestCheckResourceAttr("data.ibm_cm_share_approval_list.cm_share_approval_list", "approval_state", approvalState),
					resource.TestCheckResourceAttr("data.ibm_cm_share_approval_list.cm_share_approval_list", "limit", limit),
				),
			},
		},
	})
}

func testAccCheckIBMCmShareApprovalListDataSourceConfigBasic(objectKind string, approvalState string) string {
	return fmt.Sprintf(`
		data "ibm_cm_share_approval_list" "cm_share_approval_list" {
			object_kind = "%s"
			approval_state = "%s"
		}
	`, objectKind, approvalState)
}

func testAccCheckIBMCmShareApprovalListDataSourceConfig(objectKind string, approvalState string, limit string) string {
	return fmt.Sprintf(`
		data "ibm_cm_share_approval_list" "cm_share_approval_list" {
			object_kind = "%s"
			approval_state = "%s"
			limit = %s
		}
	`, objectKind, approvalState, limit)
}
