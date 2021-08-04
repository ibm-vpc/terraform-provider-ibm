// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsBackupPolicyJobDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBackupPolicyJobDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "backup_policy_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "backup_info.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "backup_resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "completed_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "plan_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "source_volume.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "started_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "status"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_job.is_backup_policy_job", "status_reasons.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsBackupPolicyJobDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_backup_policy_job" "is_backup_policy_job" {
			backup_policy_id = "backup_policy_id"
			id = "id"
		}
	`)
}
