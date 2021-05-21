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
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/IBM/vpc-go-sdk/vpcv1"
)

func TestAccIbmIsPlacementGroupBasic(t *testing.T) {
	var conf vpcv1.PlacementGroup
	strategy := "host_spread"
	strategyUpdate := "power_spread"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmIsPlacementGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIbmIsPlacementGroupConfigBasic(strategy),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsPlacementGroupExists("ibm_is_placement_group.is_placement_group", conf),
					resource.TestCheckResourceAttr("ibm_is_placement_group.is_placement_group", "strategy", strategy),
				),
			},
			resource.TestStep{
				Config: testAccCheckIbmIsPlacementGroupConfigBasic(strategyUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_placement_group.is_placement_group", "strategy", strategyUpdate),
				),
			},
		},
	})
}

func TestAccIbmIsPlacementGroupAllArgs(t *testing.T) {
	var conf vpcv1.PlacementGroup
	strategy := "host_spread"
	name := fmt.Sprintf("tf-pg-name%d", acctest.RandIntRange(10, 100))
	strategyUpdate := "power_spread"
	nameUpdate := fmt.Sprintf("tf-pg-name-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIbmIsPlacementGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIbmIsPlacementGroupConfig(strategy, name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckIbmIsPlacementGroupExists("ibm_is_placement_group.is_placement_group", conf),
					resource.TestCheckResourceAttr("ibm_is_placement_group.is_placement_group", "strategy", strategy),
					resource.TestCheckResourceAttr("ibm_is_placement_group.is_placement_group", "name", name),
				),
			},
			{
				Config: testAccCheckIbmIsPlacementGroupConfig(strategyUpdate, nameUpdate),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ibm_is_placement_group.is_placement_group", "strategy", strategyUpdate),
					resource.TestCheckResourceAttr("ibm_is_placement_group.is_placement_group", "name", nameUpdate),
				),
			},
			{
				ResourceName:      "ibm_is_placement_group.is_placement_group",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIbmIsPlacementGroupConfigBasic(strategy string) string {
	return fmt.Sprintf(`

		resource "ibm_is_placement_group" "is_placement_group" {
			strategy = "%s"
		}
	`, strategy)
}

func testAccCheckIbmIsPlacementGroupConfig(strategy string, name string) string {
	return fmt.Sprintf(`
		data "ibm_resource_group" "default" {
			is_default=true
		}
		resource "ibm_is_placement_group" "is_placement_group" {
			strategy = "%s"
			name = "%s"
			resource_group = data.ibm_resource_group.default.id
		}
	`, strategy, name)
}

func testAccCheckIbmIsPlacementGroupExists(n string, obj vpcv1.PlacementGroup) resource.TestCheckFunc {

	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		vpcClient, err := testAccProvider.Meta().(ClientSession).VpcV1API()
		if err != nil {
			return err
		}

		getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{}

		getPlacementGroupOptions.SetID(rs.Primary.ID)

		placementGroup, _, err := vpcClient.GetPlacementGroup(getPlacementGroupOptions)
		if err != nil {
			return err
		}

		obj = *placementGroup
		return nil
	}
}

func testAccCheckIbmIsPlacementGroupDestroy(s *terraform.State) error {
	vpcClient, err := testAccProvider.Meta().(ClientSession).VpcV1API()
	if err != nil {
		return err
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_is_placement_group" {
			continue
		}

		getPlacementGroupOptions := &vpcv1.GetPlacementGroupOptions{}

		getPlacementGroupOptions.SetID(rs.Primary.ID)

		// Try to find the key
		_, response, err := vpcClient.GetPlacementGroup(getPlacementGroupOptions)

		if err == nil {
			return fmt.Errorf("PlacementGroup still exists: %s", rs.Primary.ID)
		} else if response.StatusCode != 404 {
			return fmt.Errorf("Error checking for PlacementGroup (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}

	return nil
}
