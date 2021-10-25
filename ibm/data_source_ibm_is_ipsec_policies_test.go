// Copyright IBM Corp. 2021 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

package ibm

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIBMIsIpsecPoliciesDataSourceBasic(t *testing.T) {
	name := fmt.Sprintf("tfipsecc-name-%d", acctest.RandIntRange(10, 100))
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMIsIpsecPoliciesDataSourceConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.#"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.authentication_algorithm"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.encapsulation_mode"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.encryption_algorithm"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.key_lifetime"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.name"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.pfs"),
					resource.TestCheckResourceAttrSet("data.ibm_is_ipsec_policies.is_ipsec_policies", "ipsec_policies.0.transform_protocol"),
				),
			},
		},
	})
}

func testAccCheckIBMIsIpsecPoliciesDataSourceConfigBasic(name string) string {
	return fmt.Sprintf(`
		resource "ibm_is_ipsec_policy" "example" {
			name = "%s"
			authentication_algorithm = "sha1"
			encryption_algorithm = "aes128"
			pfs = "group_2"
		}
		data "ibm_is_ipsec_policies" "is_ipsec_policies" {
		}
	`, name)
}
