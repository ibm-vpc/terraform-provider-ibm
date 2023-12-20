// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIbmIsStorageOntapInstanceVirtualMachineVolumesDataSourceBasic(t *testing.T) {
	storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID := fmt.Sprintf("tf_storage_ontap_instance_id_%d", acctest.RandIntRange(10, 100))
	storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID := fmt.Sprintf("tf_storage_virtual_machine_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumesDataSourceConfigBasic(storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID, storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "storage_ontap_instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "storage_virtual_machine_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "total_count"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.#"),
				),
			},
		},
	})
}

func TestAccIbmIsStorageOntapInstanceVirtualMachineVolumesDataSourceAllArgs(t *testing.T) {
	storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID := fmt.Sprintf("tf_storage_ontap_instance_id_%d", acctest.RandIntRange(10, 100))
	storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID := fmt.Sprintf("tf_storage_virtual_machine_id_%d", acctest.RandIntRange(10, 100))
	storageOntapInstanceStorageVirtualMachineVolumeCapacity := fmt.Sprintf("%d", acctest.RandIntRange(10, 16000))
	storageOntapInstanceStorageVirtualMachineVolumeEnableStorageEfficiency := "false"
	storageOntapInstanceStorageVirtualMachineVolumeName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	storageOntapInstanceStorageVirtualMachineVolumeSecurityStyle := "mixed"
	storageOntapInstanceStorageVirtualMachineVolumeStorageEfficiency := "disabled"
	storageOntapInstanceStorageVirtualMachineVolumeType := "data_protection"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumesDataSourceConfig(storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID, storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID, storageOntapInstanceStorageVirtualMachineVolumeCapacity, storageOntapInstanceStorageVirtualMachineVolumeEnableStorageEfficiency, storageOntapInstanceStorageVirtualMachineVolumeName, storageOntapInstanceStorageVirtualMachineVolumeSecurityStyle, storageOntapInstanceStorageVirtualMachineVolumeStorageEfficiency, storageOntapInstanceStorageVirtualMachineVolumeType),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "storage_ontap_instance_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "storage_virtual_machine_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "next.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "total_count"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.#"),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.capacity", storageOntapInstanceStorageVirtualMachineVolumeCapacity),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.created_at"),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.enable_storage_efficiency", storageOntapInstanceStorageVirtualMachineVolumeEnableStorageEfficiency),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.health_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.junction_path"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.lifecycle_state"),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.name", storageOntapInstanceStorageVirtualMachineVolumeName),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.resource_type"),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.security_style", storageOntapInstanceStorageVirtualMachineVolumeSecurityStyle),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.storage_efficiency", storageOntapInstanceStorageVirtualMachineVolumeStorageEfficiency),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance_virtual_machine_volumes.is_storage_ontap_instance_virtual_machine_volumes", "volumes.0.type", storageOntapInstanceStorageVirtualMachineVolumeType),
				),
			},
		},
	})
}

func testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumesDataSourceConfigBasic(storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID string, storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_storage_ontap_instance_virtual_machine_volume" "is_storage_ontap_instance_virtual_machine_volume_instance" {
			storage_ontap_instance_id = "%s"
			storage_virtual_machine_id = "%s"
		}

		data "ibm_is_storage_ontap_instance_virtual_machine_volumes" "is_storage_ontap_instance_virtual_machine_volumes_instance" {
			storage_ontap_instance_id = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.storage_ontap_instance_id
			storage_virtual_machine_id = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.storage_virtual_machine_id
			name = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.name
		}
	`, storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID, storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID)
}

func testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumesDataSourceConfig(storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID string, storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID string, storageOntapInstanceStorageVirtualMachineVolumeCapacity string, storageOntapInstanceStorageVirtualMachineVolumeEnableStorageEfficiency string, storageOntapInstanceStorageVirtualMachineVolumeName string, storageOntapInstanceStorageVirtualMachineVolumeSecurityStyle string, storageOntapInstanceStorageVirtualMachineVolumeStorageEfficiency string, storageOntapInstanceStorageVirtualMachineVolumeType string) string {
	return fmt.Sprintf(`
		resource "ibm_is_storage_ontap_instance_virtual_machine_volume" "is_storage_ontap_instance_virtual_machine_volume_instance" {
			storage_ontap_instance_id = "%s"
			storage_virtual_machine_id = "%s"
			capacity = %s
			cifs_share {
				access_control_list {
					permission = "full_control"
					users = ["user1"]
				}
				mount_path = "//192.168.3.4/my-share"
				name = "my-share"
			}
			enable_storage_efficiency = %s
			export_policy {
				mount_path = "192.168.3.4:/vol/path"
				rules {
					access_control = "read_only"
					clients {
						hostname = "host1"
					}
					index = 1
					is_superuser = false
					nfs_version = ["nfs4"]
				}
			}
			name = "%s"
			security_style = "%s"
			storage_efficiency = "%s"
			type = "%s"
		}

		data "ibm_is_storage_ontap_instance_virtual_machine_volumes" "is_storage_ontap_instance_virtual_machine_volumes_instance" {
			storage_ontap_instance_id = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.storage_ontap_instance_id
			storage_virtual_machine_id = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.storage_virtual_machine_id
			name = ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume.name
		}
	`, storageOntapInstanceStorageVirtualMachineVolumeStorageOntapInstanceID, storageOntapInstanceStorageVirtualMachineVolumeStorageVirtualMachineID, storageOntapInstanceStorageVirtualMachineVolumeCapacity, storageOntapInstanceStorageVirtualMachineVolumeEnableStorageEfficiency, storageOntapInstanceStorageVirtualMachineVolumeName, storageOntapInstanceStorageVirtualMachineVolumeSecurityStyle, storageOntapInstanceStorageVirtualMachineVolumeStorageEfficiency, storageOntapInstanceStorageVirtualMachineVolumeType)
}
