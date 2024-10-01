// Copyright IBM Corp. 2024 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.91.0-d9755c53-20240605-153412
 */

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
)

func TestAccIBMIsClusterNetworkProfileDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsClusterNetworkProfileDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "family"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "supported_instance_profiles.#"),
					// resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "supported_instance_profiles.0.href"),
					// resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "supported_instance_profiles.0.name"),
					// resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "supported_instance_profiles.0.resource_type"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "zones.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "zones.0.href"),
					resource.TestCheckResourceAttrSet("data.ibm_is_cluster_network_profile.is_cluster_network_profile_instance", "zones.0.name"),
				),
			},
		},
	})
}

func testAccCheckIBMIsClusterNetworkProfileDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_cluster_network_profile" "is_cluster_network_profile_instance" {
			name = "h100"
		}
	`)
}
