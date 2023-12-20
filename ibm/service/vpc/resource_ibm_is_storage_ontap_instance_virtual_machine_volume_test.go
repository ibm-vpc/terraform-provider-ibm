// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/conns"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/flex"
	"github.com/IBM/vpc-beta-go-sdk/ontapv1"
)

func TestAccIbmIsStorageOntapInstanceVirtualMachineVolumeBasic(t *testing.T) {
	var conf ontapv1.StorageOntapInstanceStorageVirtualMachineVolume
	storageOntapInstanceID := fmt.Sprintf("tf_storage_ontap_instance_id_%d", acctest.RandIntRange(10, 100))
	storageVirtualMachineID := fmt.Sprintf("tf_storage_virtual_machine_id_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeConfigBasic(storageOntapInstanceID, storageVirtualMachineID),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeExists("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", conf),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_ontap_instance_id", storageOntapInstanceID),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_virtual_machine_id", storageVirtualMachineID),
				),
			},
		},
	})
}

func TestAccIbmIsStorageOntapInstanceVirtualMachineVolumeAllArgs(t *testing.T) {
	var conf ontapv1.StorageOntapInstanceStorageVirtualMachineVolume
	storageOntapInstanceID := fmt.Sprintf("tf_storage_ontap_instance_id_%d", acctest.RandIntRange(10, 100))
	storageVirtualMachineID := fmt.Sprintf("tf_storage_virtual_machine_id_%d", acctest.RandIntRange(10, 100))
	capacity := fmt.Sprintf("%d", acctest.RandIntRange(10, 16000))
	enableStorageEfficiency := "false"
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	securityStyle := "mixed"
	storageEfficiency := "disabled"
	typeVar := "data_protection"
	capacityUpdate := fmt.Sprintf("%d", acctest.RandIntRange(10, 16000))
	enableStorageEfficiencyUpdate := "true"
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	securityStyleUpdate := "windows"
	storageEfficiencyUpdate := "enabled"
	typeVarUpdate := "read_write"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeConfig(storageOntapInstanceID, storageVirtualMachineID, capacity, enableStorageEfficiency, name, securityStyle, storageEfficiency, typeVar),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeExists("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", conf),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_ontap_instance_id", storageOntapInstanceID),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_virtual_machine_id", storageVirtualMachineID),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "capacity", capacity),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "enable_storage_efficiency", enableStorageEfficiency),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "name", name),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "security_style", securityStyle),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_efficiency", storageEfficiency),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "type", typeVar),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeConfig(storageOntapInstanceID, storageVirtualMachineID, capacityUpdate, enableStorageEfficiencyUpdate, nameUpdate, securityStyleUpdate, storageEfficiencyUpdate, typeVarUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_ontap_instance_id", storageOntapInstanceID),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_virtual_machine_id", storageVirtualMachineID),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "capacity", capacityUpdate),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "enable_storage_efficiency", enableStorageEfficiencyUpdate),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "security_style", securityStyleUpdate),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "storage_efficiency", storageEfficiencyUpdate),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume", "type", typeVarUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_storage_ontap_instance_virtual_machine_volume.is_storage_ontap_instance_virtual_machine_volume",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeConfigBasic(storageOntapInstanceID string, storageVirtualMachineID string) string {
	return fmt.Sprintf(`
		resource "ibm_is_storage_ontap_instance_virtual_machine_volume" "is_storage_ontap_instance_virtual_machine_volume_instance" {
			storage_ontap_instance_id = "%s"
			storage_virtual_machine_id = "%s"
		}
	`, storageOntapInstanceID, storageVirtualMachineID)
}

func testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeConfig(storageOntapInstanceID string, storageVirtualMachineID string, capacity string, enableStorageEfficiency string, name string, securityStyle string, storageEfficiency string, typeVar string) string {
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
				name = "my-share"
			}
			enable_storage_efficiency = %s
			export_policy {
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
	`, storageOntapInstanceID, storageVirtualMachineID, capacity, enableStorageEfficiency, name, securityStyle, storageEfficiency, typeVar)
}

func testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeExists(n string, obj ontapv1.StorageOntapInstanceStorageVirtualMachineVolume) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		ontapClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).OntapAPI()
		if err != nil {
			return err
		}

		getStorageOntapInstanceStorageVirtualMachineVolumeOptions := &ontapv1.GetStorageOntapInstanceStorageVirtualMachineVolumeOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageOntapInstanceID(parts[0])
		getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageVirtualMachineID(parts[1])
		getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetID(parts[2])

		storageOntapInstanceStorageVirtualMachineVolume, _, err := ontapClient.GetStorageOntapInstanceStorageVirtualMachineVolume(getStorageOntapInstanceStorageVirtualMachineVolumeOptions)
		if err != nil {
			return err
		}

		obj = *storageOntapInstanceStorageVirtualMachineVolume
		return nil
	}
}

func testAccCheckIbmIsStorageOntapInstanceVirtualMachineVolumeDestroy(s *terraform.State) error {
	ontapClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).OntapAPI()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_storage_ontap_instance_virtual_machine_volume" {
			continue
		}

		getStorageOntapInstanceStorageVirtualMachineVolumeOptions := &ontapv1.GetStorageOntapInstanceStorageVirtualMachineVolumeOptions{}

		parts, err := flex.SepIdParts(rs.Primary.ID, "/")
		if err != nil {
			return err
		}

		getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageOntapInstanceID(parts[0])
		getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetStorageVirtualMachineID(parts[1])
		getStorageOntapInstanceStorageVirtualMachineVolumeOptions.SetID(parts[2])

		// Try to find the key
		_, response, err := ontapClient.GetStorageOntapInstanceStorageVirtualMachineVolume(getStorageOntapInstanceStorageVirtualMachineVolumeOptions)

		if err == nil {
			return fmt.Errorf("StorageOntapInstanceStorageVirtualMachineVolume still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for StorageOntapInstanceStorageVirtualMachineVolume (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
