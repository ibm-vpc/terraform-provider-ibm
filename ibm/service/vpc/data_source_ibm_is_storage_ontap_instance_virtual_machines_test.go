// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIbmIsStorageOntapInstanceVirtualMachinesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceVirtualMachinesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machines.is_storage_ontap_instance_virtual_machines", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machines.is_storage_ontap_instance_virtual_machines", "storage_ontap_instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machines.is_storage_ontap_instance_virtual_machines", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machines.is_storage_ontap_instance_virtual_machines", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machines.is_storage_ontap_instance_virtual_machines", "storage_virtual_machines.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machines.is_storage_ontap_instance_virtual_machines", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIbmIsStorageOntapInstanceVirtualMachinesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_storage_ontap_instance_virtual_machines" "is_storage_ontap_instance_virtual_machines_instance" {
			storage_ontap_instance_id = "storage_ontap_instance_id"
			name = "name"
		}
	`)
}
