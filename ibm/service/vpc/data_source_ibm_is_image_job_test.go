// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsImageJobDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageJobDataSourceConfigBasic("build"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "image_job"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "status"),
					resource.TestCheckResourceAttr("data.ibm_is_image_job.is_image_job_instance", "type", "build"),
				),
			},
		},
	})
}

func TestAccIBMIsImageJobDataSourceAllArgs(t *testing.T) {
	name := fmt.Sprintf("tf-job-ds-%d", acctest.RandIntRange(10, 100))
	timeout := fmt.Sprintf("%d", acctest.RandIntRange(30, 480))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageJobDataSourceConfig(name, "build", "true", timeout),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "image_job"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "href"),
					resource.TestCheckResourceAttr("data.ibm_is_image_job.is_image_job_instance", "name", name),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "status"),
					resource.TestCheckResourceAttr("data.ibm_is_image_job.is_image_job_instance", "type", "build"),
					resource.TestCheckResourceAttr("data.ibm_is_image_job.is_image_job_instance", "delete_image_job_on_image_delete", "true"),
					resource.TestCheckResourceAttr("data.ibm_is_image_job.is_image_job_instance", "timeout", timeout),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "image_recipe.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "image_template.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_job.is_image_job_instance", "worker.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageJobDataSourceConfigBasic(typeVar string) string {
	return testAccCheckIBMIsImageJobConfigDependencies() + fmt.Sprintf(`
		resource "ibm_is_image_job" "is_image_job_instance" {
			type          = "%s"
			image_recipe  = ibm_is_image_recipe.is_image_job_recipe.id
			worker_subnet = ibm_is_subnet.is_image_job_subnet.id
		}

		data "ibm_is_image_job" "is_image_job_instance" {
			image_job = ibm_is_image_job.is_image_job_instance.id
		}
	`, typeVar)
}

func testAccCheckIBMIsImageJobDataSourceConfig(name string, typeVar string, deleteImageJobOnImageDelete string, timeout string) string {
	return testAccCheckIBMIsImageJobConfigDependencies() + fmt.Sprintf(`
		resource "ibm_is_image_job" "is_image_job_instance" {
			name                             = "%s"
			type                             = "%s"
			delete_image_job_on_image_delete = %s
			image_recipe                     = ibm_is_image_recipe.is_image_job_recipe.id
			image_recipe_version             = ibm_is_image_recipe.is_image_job_recipe.version
			image_template_name              = "tf-image-job-image"
			image_template_resource_group    = "%s"
			timeout                          = %s
			worker_profile                   = "bx2-4x16"
			worker_subnet                    = ibm_is_subnet.is_image_job_subnet.id
		}

		data "ibm_is_image_job" "is_image_job_instance" {
			image_job = ibm_is_image_job.is_image_job_instance.id
		}
	`, name, typeVar, deleteImageJobOnImageDelete, acc.IsResourceGroupID, timeout)
}
