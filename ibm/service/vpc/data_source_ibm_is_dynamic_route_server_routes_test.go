// Copyright IBM Corp. 2023 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsDynamicRouteServerRoutesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsDynamicRouteServerRoutesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_routes.is_dynamic_route_server_routes", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_routes.is_dynamic_route_server_routes", "dynamic_route_server_id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_routes.is_dynamic_route_server_routes", "first.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_routes.is_dynamic_route_server_routes", "limit"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_routes.is_dynamic_route_server_routes", "routes.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_dynamic_route_server_routes.is_dynamic_route_server_routes", "total_count"),
				),
			},
		},
	})
}

func testAccCheckIBMIsDynamicRouteServerRoutesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_dynamic_route_server_routes" "is_dynamic_route_server_routes_instance" {
			dynamic_route_server_id = "dynamic_route_server_id"
			sort = "created_at"
			peer.id = "peer.id"
			type = "learned"
		}
	`)
}
