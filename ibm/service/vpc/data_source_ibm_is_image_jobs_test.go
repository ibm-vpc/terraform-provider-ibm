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

func TestAccIBMIsImageJobsDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageJobsDataSourceConfigBasic("build"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.#"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.type", "build"),
				),
			},
		},
	})
}

func TestAccIBMIsImageJobsDataSourceAllArgs(t *testing.T) {
	name := fmt.Sprintf("tf-jobs-ds-%d", acctest.RandIntRange(10, 100))
	timeout := fmt.Sprintf("%d", acctest.RandIntRange(30, 480))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageJobsDataSourceConfig(name, "build", "true", timeout),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "id"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "name", name),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.#", "1"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.id"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.name", name),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.status"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.type", "build"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.delete_image_job_on_image_delete", "true"),
					resource.TestCheckResourceAttr("data.ibm_is_image_jobs.is_image_jobs_instance", "jobs.0.timeout", timeout),
				),
			},
		},
	})
}

func testAccCheckIBMIsImageJobsDataSourceConfigBasic(typeVar string) string {
	return testAccCheckIBMIsImageJobConfigDependencies() + fmt.Sprintf(`
		resource "ibm_is_image_job" "is_image_job_instance" {
			type          = "%s"
			image_recipe  = ibm_is_image_recipe.is_image_job_recipe.id
			worker_subnet = ibm_is_subnet.is_image_job_subnet.id
		}

		data "ibm_is_image_jobs" "is_image_jobs_instance" {
			name = ibm_is_image_job.is_image_job_instance.name
		}
	`, typeVar)
}

func testAccCheckIBMIsImageJobsDataSourceConfig(name string, typeVar string, deleteImageJobOnImageDelete string, timeout string) string {
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

		data "ibm_is_image_jobs" "is_image_jobs_instance" {
			name = ibm_is_image_job.is_image_job_instance.name
		}
	`, name, typeVar, deleteImageJobOnImageDelete, acc.IsResourceGroupID, timeout)
}
