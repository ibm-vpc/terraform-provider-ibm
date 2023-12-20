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
	"github.com/IBM/vpc-beta-go-sdk/ontapv1"
)

func TestAccIbmIsStorageOntapInstanceBasic(t *testing.T) {
	var conf ontapv1.StorageOntapInstance
	capacity := fmt.Sprintf("%d", acctest.RandIntRange(1, 64))
	capacityUpdate := fmt.Sprintf("%d", acctest.RandIntRange(1, 64))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIbmIsStorageOntapInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceConfigBasic(capacity),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsStorageOntapInstanceExists("ibm_is_storage_ontap_instance.is_storage_ontap_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance.is_storage_ontap_instance", "capacity", capacity),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceConfigBasic(capacityUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance.is_storage_ontap_instance", "capacity", capacityUpdate),
				),
			},
		},
	})
}

func TestAccIbmIsStorageOntapInstanceAllArgs(t *testing.T) {
	var conf ontapv1.StorageOntapInstance
	capacity := fmt.Sprintf("%d", acctest.RandIntRange(1, 64))
	name := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))
	capacityUpdate := fmt.Sprintf("%d", acctest.RandIntRange(1, 64))
	nameUpdate := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIbmIsStorageOntapInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceConfig(capacity, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsStorageOntapInstanceExists("ibm_is_storage_ontap_instance.is_storage_ontap_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance.is_storage_ontap_instance", "capacity", capacity),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance.is_storage_ontap_instance", "name", name),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceConfig(capacityUpdate, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance.is_storage_ontap_instance", "capacity", capacityUpdate),
					resource.TestCheckResourceAttr("ibm_is_storage_ontap_instance.is_storage_ontap_instance", "name", nameUpdate),
				),
			},
			resource.TestStep{
				ResourceName:      "ibm_is_storage_ontap_instance.is_storage_ontap_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIbmIsStorageOntapInstanceConfigBasic(capacity string) string {
	return fmt.Sprintf(`
		resource "ibm_is_storage_ontap_instance" "is_storage_ontap_instance_instance" {
			address_prefix {
				href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/vpcs/8e454ead-0db7-48ac-9a8b-2698d8c470a7/address_prefixes/1a15dca5-7e33-45e1-b7c5-bc690e569531"
				id = "1a15dca5-7e33-45e1-b7c5-bc690e569531"
			}
			capacity = %s
			storage_virtual_machines {
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/storage_ontap_instances/r134-d7cc5196-9864-48c4-82d8-3f30da41ffc5/storage_virtual_machines/r134-efee5196-9864-48c4-82d8-3f30da41ffc5"
				id = "r134-efee5196-9864-48c4-82d8-3f30da41ffc5"
				name = "my-storage-virtual-machine"
				resource_type = "storage_ontap_instance_storage_virtual_machine"
			}
		}
	`, capacity)
}

func testAccCheckIbmIsStorageOntapInstanceConfig(capacity string, name string) string {
	return fmt.Sprintf(`

		resource "ibm_is_storage_ontap_instance" "is_storage_ontap_instance_instance" {
			active_subnet {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::subnet:7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
			}
			address_prefix {
				href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/vpcs/8e454ead-0db7-48ac-9a8b-2698d8c470a7/address_prefixes/1a15dca5-7e33-45e1-b7c5-bc690e569531"
				id = "1a15dca5-7e33-45e1-b7c5-bc690e569531"
			}
			admin_credentials {
				http {
					crn = "crn:v1:bluemix:public:secrets-manager:eu-gb:a/123456:43af9a51-2dca-4947-b36b-8c41363537b7:secret:0736e7b6-7fa7-1524-a370-44f09894866e"
					resource_type = "credential"
				}
				password {
					crn = "crn:v1:bluemix:public:secrets-manager:eu-gb:a/123456:43af9a51-2dca-4947-b36b-8c41363537b7:secret:0736e7b6-7fa7-1524-a370-44f09894866e"
					resource_type = "credential"
				}
				ssh {
					crn = "crn:v1:bluemix:public:secrets-manager:eu-gb:a/123456:43af9a51-2dca-4947-b36b-8c41363537b7:secret:0736e7b6-7fa7-1524-a370-44f09894866e"
					resource_type = "credential"
				}
			}
			admin_password {
				crn = "crn:v1:bluemix:public:secrets-manager:eu-gb:a/123456:43af9a51-2dca-4947-b36b-8c41363537b7:secret:0736e7b6-7fa7-1524-a370-44f09894866e"
			}
			capacity = %s
			encryption_key {
				crn = "crn:v1:bluemix:public:kms:us-south:a/123456:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"
			}
			name = "%s"
			primary_subnet {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::subnet:7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
			}
			resource_group {
				id = "fee82deba12e4c0fb69c3b09d1f12345"
			}
			routing_tables {
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/vpcs/982d72b7-db1b-4606-afb2-ed6bd4b0bed1/routing_tables/6885e83f-03b2-4603-8a86-db2a0f55c840"
				id = "1a15dca5-7e33-45e1-b7c5-bc690e569531"
				name = "my-routing-table-1"
				resource_type = "routing_table"
			}
			secondary_subnet {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::subnet:7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
			}
			security_groups {
				crn = "crn:v1:bluemix:public:is:us-south:a/123456::security-group:be5df5ca-12a0-494b-907e-aa6ec2bfa271"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				id = "be5df5ca-12a0-494b-907e-aa6ec2bfa271"
				name = "my-security-group"
			}
			standby_subnet {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::subnet:7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
			}
			storage_virtual_machines {
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/storage_ontap_instances/r134-d7cc5196-9864-48c4-82d8-3f30da41ffc5/storage_virtual_machines/r134-efee5196-9864-48c4-82d8-3f30da41ffc5"
				id = "r134-efee5196-9864-48c4-82d8-3f30da41ffc5"
				name = "my-storage-virtual-machine"
				resource_type = "storage_ontap_instance_storage_virtual_machine"
			}
		}
	`, capacity, name)
}

func testAccCheckIbmIsStorageOntapInstanceExists(n string, obj ontapv1.StorageOntapInstance) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		ontapClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).OntapAPI()
		if err != nil {
			return err
		}

		getStorageOntapInstanceOptions := &ontapv1.GetStorageOntapInstanceOptions{}

		getStorageOntapInstanceOptions.SetID(rs.Primary.ID)

		storageOntapInstance, _, err := ontapClient.GetStorageOntapInstance(getStorageOntapInstanceOptions)
		if err != nil {
			return err
		}

		obj = *storageOntapInstance
		return nil
	}
}

func testAccCheckIbmIsStorageOntapInstanceDestroy(s *terraform.State) error {
	ontapClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).OntapAPI()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_storage_ontap_instance" {
			continue
		}

		getStorageOntapInstanceOptions := &ontapv1.GetStorageOntapInstanceOptions{}

		getStorageOntapInstanceOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := ontapClient.GetStorageOntapInstance(getStorageOntapInstanceOptions)

		if err == nil {
			return fmt.Errorf("StorageOntapInstance still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for StorageOntapInstance (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
