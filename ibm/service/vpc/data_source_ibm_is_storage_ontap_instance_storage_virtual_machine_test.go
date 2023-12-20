// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIbmIsStorageOntapInstanceStorageVirtualMachineDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceStorageVirtualMachineDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "storage_ontap_instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "endpoints.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_storage_virtual_machine.is_storage_ontap_instance_storage_virtual_machine", "resource_type"),
				),
			},
		},
	})
}

func testAccCheckIbmIsStorageOntapInstanceStorageVirtualMachineDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_storage_ontap_instance_storage_virtual_machine" "is_storage_ontap_instance_storage_virtual_machine_instance" {
			storage_ontap_instance_id = "storage_ontap_instance_id"
			id = "id"
		}
	`)
}
