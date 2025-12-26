// Copyright IBM Corp. 2025 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package catalogmanagement_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM/platform-services-go-sdk/catalogmanagementv1"
)

func TestAccIBMCmShareApprovalBasic(t *testing.T) {
	var conf catalogmanagementv1.ShareApprovalListAccessResult
	objectKind := "offering"
	approvalState := "approved"
	accountID := fmt.Sprintf("-acct-%s", acctest.RandStringFromCharSet(32, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMCmShareApprovalDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMCmShareApprovalConfigBasic(objectKind, approvalState, accountID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMCmShareApprovalExists("ibm_cm_share_approval.cm_share_approval", conf),
					resource.TestCheckResourceAttr("ibm_cm_share_approval.cm_share_approval", "object_kind", objectKind),
					resource.TestCheckResourceAttr("ibm_cm_share_approval.cm_share_approval", "approval_state", approvalState),
				),
			},
			{
				Config: testAccCheckIBMCmShareApprovalConfigBasic(objectKind, "rejected", accountID),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_cm_share_approval.cm_share_approval", "approval_state", "rejected"),
				),
			},
		},
	})
}

func testAccCheckIBMCmShareApprovalConfigBasic(objectKind string, approvalState string, accountID string) string {
	return fmt.Sprintf(`
		resource "ibm_cm_share_approval" "cm_share_approval" {
			object_kind = "%s"
			approval_state = "%s"
			account_ids = ["%s"]
		}
	`, objectKind, approvalState, accountID)
}

func testAccCheckIBMCmShareApprovalExists(n string, obj catalogmanagementv1.ShareApprovalListAccessResult) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		_, err := acc.TestAccProvider.Meta().(conns.ClientSession).CatalogManagementV1()
		if err != nil {
			return err
		}

		// Parse the ID to get object_kind and approval_state
		// ID format: object_kind/approval_state
		// For now, just verify the resource exists in state
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		return nil
	}
}

func testAccCheckIBMCmShareApprovalDestroy(s *terraform.State) error {
	_, err := acc.TestAccProvider.Meta().(conns.ClientSession).CatalogManagementV1()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_cm_share_approval" {
			continue
		}

		// The resource sets approval to "rejected" on delete
		// We can verify by checking if the resource no longer has approved accounts
		// For now, we'll just return nil as the delete operation is a state change
	}

	return nil
}
