// Copyright IBM Corp. 2017, 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccIBMISImagesDataSource_basic(t *testing.T) {
	resName := "data.ibm_is_images.test1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.0.name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.status"),
					resource.TestCheckResourceAttrSet(resName, "images.0.architecture"),
				),
			},
		},
	})
}
func TestAccIBMISImagesDataSource_Zones(t *testing.T) {
	resName := "data.ibm_is_images.test1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceNameConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.0.name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.status"),
					resource.TestCheckResourceAttrSet(resName, "images.0.architecture"),
					resource.TestCheckResourceAttrSet(resName, "images.0.zones.#"),
					resource.TestCheckResourceAttrSet(resName, "images.0.zones.0.name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.zones.0.href"),
				),
			},
		},
	})
}
func TestAccIBMISImagesDataSource_All(t *testing.T) {
	resName := "data.ibm_is_images.test1"
	imageName := fmt.Sprintf("tfimage-name-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceAllConfig(imageName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.dedicated_host_only"),
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.display_name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.family"),
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.href"),
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.vendor"),
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.version"),
					resource.TestCheckResourceAttrSet(resName, "images.0.operating_system.0.architecture"),
					resource.TestCheckResourceAttrSet(resName, "images.0.status"),
					resource.TestCheckResourceAttrSet(resName, "images.0.resource_group.0.id"),
					resource.TestCheckResourceAttrSet(resName, "images.0.resource_group.0.name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.allowed_use.#"),
					resource.TestCheckResourceAttrSet(resName, "images.0.allowed_use.0.api_version"),
					resource.TestCheckResourceAttrSet(resName, "images.0.allowed_use.0.bare_metal_server"),
					resource.TestCheckResourceAttrSet(resName, "images.0.allowed_use.0.instance"),
				),
			},
		},
	})
}
func TestAccIBMISImagesDataSource_catalog(t *testing.T) {
	resName := "data.ibm_is_images.test1"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISCatalogImagesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.0.name"),
					resource.TestCheckResourceAttr(resName, "images.0.catalog_offering.0.managed", "true"),
					resource.TestCheckResourceAttrSet(resName, "images.0.status"),
					resource.TestCheckResourceAttrSet(resName, "images.0.architecture"),
				),
			},
		},
	})
}

func TestAccIBMISImageDataSource_With_FilterVisibilty(t *testing.T) {
	resName := "data.ibm_is_images.test1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceWithVisibilityPublic("public"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.0.name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.status"),
					resource.TestCheckResourceAttrSet(resName, "images.0.architecture"),
				),
			},
		},
	})
}

func TestAccIBMISImageDataSource_With_FilterRemoteAccountId(t *testing.T) {
	resName := "data.ibm_is_images.example"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceWithRemoteAccountId("provider"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.0.remote.#"),
					resource.TestCheckResourceAttrSet(resName, "images.0.remote.0.account.#"),
					resource.TestCheckResourceAttrSet(resName, "images.0.remote.0.account.0.id"),
					resource.TestCheckResourceAttrSet(resName, "images.0.remote.0.account.0.resource_type"),
				),
			},
		},
	})
}

func TestAccIBMISImageDataSource_With_FilterStatus(t *testing.T) {
	resName := "data.ibm_is_images.test1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceWithStatusPublic("available"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.0.name"),
					resource.TestCheckResourceAttrSet(resName, "images.0.status"),
					resource.TestCheckResourceAttrSet(resName, "images.0.architecture"),
				),
			},
		},
	})
}

func TestAccIBMISImageDataSource_With_FilterSourceImageRecipe(t *testing.T) {
	resName := "data.ibm_is_images.test1"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceWithSourceImageRecipe(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resName, "images.#"),
				),
			},
		},
	})
}

func TestAccIBMISImagesDataSource_SourceImageRecipeAndBuildJob(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMISImagesDataSourceSourceImageRecipeAndBuildJobConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.source_image_recipe.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.source_image_recipe.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.source_image_recipe.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.source_image_recipe.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.source_image_recipe.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.source_image_recipe.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.source_image_recipe.0.version"),
					resource.TestCheckResourceAttrPair("data.ibm_is_images.test1", "images.0.source_image_recipe.0.id", "ibm_is_image_recipe.is_image_job_recipe", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.build_job.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.build_job.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.build_job.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.build_job.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_images.test1", "images.0.build_job.0.resource_type"),
					resource.TestCheckResourceAttrPair("data.ibm_is_images.test1", "images.0.build_job.0.id", "ibm_is_image_job.is_image_job_instance", "id"),
				),
			},
		},
	})
}

func TestIBMISImagesDataSourceSchema_SourceImageRecipeFilters(t *testing.T) {
	s := vpc.DataSourceIBMISImages().Schema
	for _, key := range []string{"source_image_recipe_id", "source_image_recipe_version"} {
		field, ok := s[key]
		if !ok {
			t.Fatalf("expected %s in ibm_is_images schema", key)
		}
		if !field.Optional || field.Computed {
			t.Fatalf("expected %s to be an optional, non-computed filter", key)
		}
		if field.Type != schema.TypeString {
			t.Fatalf("expected %s to be a string", key)
		}
	}
}

func testAccCheckIBMISImagesDataSourceConfig() string {
	// status filter defaults to empty
	return fmt.Sprintf(`
      data "ibm_is_images" "test1" {
      }`)
}

func testAccCheckIBMISImagesDataSourceNameConfig() string {
	// status filter defaults to empty
	return fmt.Sprintf(`
      data "ibm_is_images" "test1" {
	  	name = "ibm-ubuntu-22-04-5-minimal-amd64-13"
      }`)
}
func testAccCheckIBMISImagesDataSourceAllConfig(imageName string) string {
	return fmt.Sprintf(`
	data "ibm_is_images" "test1" {
		status = "available"
	}`)
}
func testAccCheckIBMISCatalogImagesDataSourceConfig() string {
	// status filter defaults to empty
	return fmt.Sprintf(`
      data "ibm_is_images" "test1" {
		catalog_managed = true
      }`)
}

func testAccCheckIBMISImagesDataSourceWithVisibilityPublic(visibility string) string {
	return fmt.Sprintf(`
	data "ibm_is_images" "test1" {
		visibility = "%s"
	}
	`, visibility)
}

func testAccCheckIBMISImagesDataSourceWithStatusPublic(status string) string {
	return fmt.Sprintf(`
	data "ibm_is_images" "test1" {
		status = "%s"
	}
	`, status)
}

func testAccCheckIBMISImagesDataSourceWithSourceImageRecipe() string {
	return testAccCheckIBMIsImageJobConfigDependencies() + `
	data "ibm_is_images" "test1" {
		source_image_recipe_id      = ibm_is_image_recipe.is_image_job_recipe.id
		source_image_recipe_version = ibm_is_image_recipe.is_image_job_recipe.version
	}
	`
}

func testAccCheckIBMISImagesDataSourceSourceImageRecipeAndBuildJobConfig() string {
	return testAccCheckIBMIsImageJobConfigDependencies() + `
	resource "ibm_is_image_job" "is_image_job_instance" {
		type          = "build"
		image_recipe  = ibm_is_image_recipe.is_image_job_recipe.id
		worker_subnet = ibm_is_subnet.is_image_job_subnet.id
	}

	data "ibm_is_images" "test1" {
		source_image_recipe_id      = ibm_is_image_recipe.is_image_job_recipe.id
		source_image_recipe_version = ibm_is_image_recipe.is_image_job_recipe.version
		depends_on                  = [ibm_is_image_job.is_image_job_instance]
	}
	`
}

func testAccCheckIBMISImagesDataSourceWithRemoteAccountId(remoteAccountId string) string {
	return fmt.Sprintf(`
	data "ibm_is_images" "example" {
		remote_account_id = "%s"
		status = "available"
	}
	`, remoteAccountId)
}
