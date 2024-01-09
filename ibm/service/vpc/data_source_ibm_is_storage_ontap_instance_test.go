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

func TestAccIbmIsStorageOntapInstanceDataSourceBasic(t *testing.T) {
	storageOntapInstanceCapacity := fmt.Sprintf("%d", acctest.RandIntRange(1, 64))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceDataSourceConfigBasic(storageOntapInstanceCapacity),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "identifier"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "address_prefix.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "capacity"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "encryption"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "endpoints.#"),
					// resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "health_reasons.#"),
					// resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "health_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "href"),
					// resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "lifecycle_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "resource_group.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "routing_tables.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "security_groups.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "storage_virtual_machines.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "vpc.#"),
				),
			},
		},
	})
}

func TestAccIbmIsStorageOntapInstanceDataSourceAllArgs(t *testing.T) {
	storageOntapInstanceCapacity := fmt.Sprintf("%d", acctest.RandIntRange(1, 64))
	storageOntapInstanceName := fmt.Sprintf("tf_name_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsStorageOntapInstanceDataSourceConfig(storageOntapInstanceCapacity, storageOntapInstanceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "address_prefix.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "admin_credentials.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "capacity"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "encryption"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "encryption_key.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "endpoints.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "health_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "health_reasons.0.code"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "health_reasons.0.message"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "health_reasons.0.more_info"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "health_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "lifecycle_reasons.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "lifecycle_reasons.0.code"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "lifecycle_reasons.0.message"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "lifecycle_reasons.0.more_info"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "primary_subnet.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "resource_group.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "routing_tables.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "routing_tables.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "routing_tables.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "routing_tables.0.name", storageOntapInstanceName),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "routing_tables.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "secondary_subnet.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "security_groups.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "security_groups.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "security_groups.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "security_groups.0.name", storageOntapInstanceName),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "storage_virtual_machines.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "storage_virtual_machines.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "storage_virtual_machines.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "storage_virtual_machines.0.name", storageOntapInstanceName),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "storage_virtual_machines.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_storage_ontap_instance.is_storage_ontap_instance", "vpc.#"),
				),
			},
		},
	})
}

func testAccCheckIbmIsStorageOntapInstanceDataSourceConfigBasic(storageOntapInstanceCapacity string) string {
	return fmt.Sprint(`
		data "ibm_is_storage_ontap_instance" "is_storage_ontap_instance" {
			identifier = "r134-826d7442-3312-4be3-bc02-c527d4975dbd"
		}
	`)
}

func testAccCheckIbmIsStorageOntapInstanceDataSourceConfig(storageOntapInstanceCapacity string, storageOntapInstanceName string) string {
	return fmt.Sprintf(`
		resource "ibm_is_storage_ontap_instance" "is_storage_ontap_instance_instance" {
			address_prefix {
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				href = "https://us-south.iaas.cloud.ibm.com/ontap/v1/vpcs/8e454ead-0db7-48ac-9a8b-2698d8c470a7/address_prefixes/1a15dca5-7e33-45e1-b7c5-bc690e569531"
				id = "1a15dca5-7e33-45e1-b7c5-bc690e569531"
				name = "my-address-prefix-1"
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
			capacity = %s
			encryption_key {
				crn = "crn:v1:bluemix:public:kms:us-south:a/123456:e4a29d1a-2ef0-42a6-8fd2-350deb1c647e:key:5437653b-c4b1-447f-9646-b2a2a4cd6179"
			}
			name = "%s"
			primary_subnet {
				crn = "crn:v1:bluemix:public:is:us-south-1:a/123456::subnet:7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				name = "my-subnet"
				resource_type = "subnet"
			}
			resource_group {
				id = "fee82deba12e4c0fb69c3b09d1f12345"
				name = "my-resource-group"
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
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				id = "7ec86020-1c6e-4889-b3f0-a15f2e50f87e"
				name = "my-subnet"
				resource_type = "subnet"
			}
			security_groups {
				crn = "crn:v1:bluemix:public:is:us-south:a/123456::security-group:be5df5ca-12a0-494b-907e-aa6ec2bfa271"
				deleted {
					more_info = "https://cloud.ibm.com/apidocs/vpc#deleted-resources"
				}
				id = "be5df5ca-12a0-494b-907e-aa6ec2bfa271"
				name = "my-security-group"
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

		data "ibm_is_storage_ontap_instance" "is_storage_ontap_instance_instance" {
			id = "id"
		}
	`, storageOntapInstanceCapacity, storageOntapInstanceName)
}
