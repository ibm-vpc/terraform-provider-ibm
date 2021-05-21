/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIbmIsPlacementGroupsDataSourceBasic(t *testing.T) {
	placementGroupStrategy := "host_spread"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsPlacementGroupsDataSourceConfigBasic(placementGroupStrategy),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.strategy"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "total_count"),
				),
			},
		},
	})
}

func TestAccIbmIsPlacementGroupsDataSourceAllArgs(t *testing.T) {
	placementGroupStrategy := "host_spread"
	placementGroupName := fmt.Sprintf("tf-pg-name%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsPlacementGroupsDataSourceConfig(placementGroupStrategy, placementGroupName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.created_at"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.crn"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.lifecycle_state"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "placement_groups.0.strategy"),
					resource.TestCheckResourceAttrSet("data.ibm_is_placement_groups.is_placement_groups", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIbmIsPlacementGroupsDataSourceConfigBasic(placementGroupStrategy string) string {
	return fmt.Sprintf(`
		resource "ibm_is_placement_group" "is_placement_group" {
			strategy = "%s"
		}

		data "ibm_is_placement_groups" "is_placement_groups" {
		}
	`, placementGroupStrategy)
}

func testAccCheckIbmIsPlacementGroupsDataSourceConfig(placementGroupStrategy string, placementGroupName string) string {
	return fmt.Sprintf(`
		data "ibm_resource_group" "default" {
			is_default=true
		}
		resource "ibm_is_placement_group" "is_placement_group" {
			strategy = "%s"
			name = "%s"
			resource_group = data.ibm_resource_group.default.id
		}

		data "ibm_is_placement_groups" "is_placement_groups" {
		}
	`, placementGroupStrategy, placementGroupName)
}
