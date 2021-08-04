// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsBackupPolicyJobsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsBackupPolicyJobsDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_jobs.is_backup_policy_jobs", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_jobs.is_backup_policy_jobs", "backup_policy_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_backup_policy_jobs.is_backup_policy_jobs", "jobs.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsBackupPolicyJobsDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_backup_policy_jobs" "is_backup_policy_jobs" {
			backup_policy_id = "backup_policy_id"
		}
	`)
}
