// Copyright IBM Corp. 2026 All Rights Reserved.
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
	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIBMIsImageJobBasic(t *testing.T) {
	var conf vpcv1.ImageJob
	typeVar := "build"
	typeVarUpdate := "build"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageJobConfigBasic(typeVar),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageJobExists("ibm_is_image_job.is_image_job_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVar),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "image_recipe"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "worker_subnet"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "href"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "status"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "created_at"),
				),
			},
			{
				Config: testAccCheckIBMIsImageJobConfigBasic(typeVarUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVarUpdate),
				),
			},
		},
	})
}

func TestAccIBMIsImageJobAllArgs(t *testing.T) {
	var conf vpcv1.ImageJob
	name := fmt.Sprintf("tf-job-%d", acctest.RandIntRange(10, 100))
	typeVar := "build"
	deleteImageJobOnImageDelete := "true"
	timeout := fmt.Sprintf("%d", acctest.RandIntRange(30, 480))
	nameUpdate := fmt.Sprintf("tf-job-up-%d", acctest.RandIntRange(10, 100))
	typeVarUpdate := "build"
	deleteImageJobOnImageDeleteUpdate := "false"
	timeoutUpdate := fmt.Sprintf("%d", acctest.RandIntRange(30, 480))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acc.TestAccPreCheck(t) },
		Providers:    acc.TestAccProviders,
		CheckDestroy: testAccCheckIBMIsImageJobDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsImageJobConfig(name, typeVar, deleteImageJobOnImageDelete, timeout),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIBMIsImageJobExists("ibm_is_image_job.is_image_job_instance", conf),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "name", name),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVar),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "delete_image_job_on_image_delete", deleteImageJobOnImageDelete),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "timeout", timeout),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "worker_profile", "bx2-4x16"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "image_recipe"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "worker_subnet"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "href"),
					resource.TestCheckResourceAttrSet("ibm_is_image_job.is_image_job_instance", "status"),
				),
			},
			{
				Config: testAccCheckIBMIsImageJobConfig(nameUpdate, typeVarUpdate, deleteImageJobOnImageDeleteUpdate, timeoutUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "name", nameUpdate),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "type", typeVarUpdate),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "delete_image_job_on_image_delete", deleteImageJobOnImageDeleteUpdate),
					resource.TestCheckResourceAttr("ibm_is_image_job.is_image_job_instance", "timeout", timeoutUpdate),
				),
			},
			{
				ResourceName:      "ibm_is_image_job.is_image_job_instance",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMIsImageJobConfigDependencies() string {
	return fmt.Sprintf(`
		resource "ibm_is_vpc" "is_image_job_vpc" {
			name = "tf-image-job-vpc"
		}

		resource "ibm_is_subnet" "is_image_job_subnet" {
			name            = "tf-image-job-subnet"
			vpc             = ibm_is_vpc.is_image_job_vpc.id
			zone            = "%s"
			ipv4_cidr_block = "%s"
		}

		resource "ibm_is_image_recipe" "is_image_job_recipe" {
			base_image       = "r006-72b27b5c-f4b0-48bb-b954-5becc7c1dcb8"
			image_components = ["r006-72b27b5c-f4b0-48bb-b954-5becc7f1bcb8"]
		}
	`, acc.ISZoneName, acc.ISCIDR)
}

func testAccCheckIBMIsImageJobConfigBasic(typeVar string) string {
	return testAccCheckIBMIsImageJobConfigDependencies() + fmt.Sprintf(`
		resource "ibm_is_image_job" "is_image_job_instance" {
			type          = "%s"
			image_recipe  = ibm_is_image_recipe.is_image_job_recipe.id
			worker_subnet = ibm_is_subnet.is_image_job_subnet.id
		}
	`, typeVar)
}

func testAccCheckIBMIsImageJobConfig(name string, typeVar string, deleteImageJobOnImageDelete string, timeout string) string {
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
	`, name, typeVar, deleteImageJobOnImageDelete, acc.IsResourceGroupID, timeout)
}

func testAccCheckIBMIsImageJobExists(n string, obj vpcv1.ImageJob) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getImageJobOptions := &vpcv1.GetImageJobOptions{}
		getImageJobOptions.SetID(rs.Primary.ID)

		imageJobIntf, _, err := vpcClient.GetImageJob(getImageJobOptions)
		if err != nil {
			return err
		}

		imageJob := imageJobIntf.(*vpcv1.ImageJob)
		obj = *imageJob
		return nil
	}
}

func testAccCheckIBMIsImageJobDestroy(s *terraform.State) error {
	vpcClient, err := acc.TestAccProvider.Meta().(conns.ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_image_job" {
			continue
		}

		getImageJobOptions := &vpcv1.GetImageJobOptions{}
		getImageJobOptions.SetID(rs.Primary.ID)

		_, response, err := vpcClient.GetImageJob(getImageJobOptions)
		if err == nil {
			return fmt.Errorf("ImageJob still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("error checking for ImageJob (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}
	return nil
}
