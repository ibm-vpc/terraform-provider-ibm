// Copyright IBM Corp. 2026 All Rights Reserved.
// Licensed under the Mozilla Public License v2.0

/*
 * IBM OpenAPI Terraform Generator Version: 3.111.0-1bfb72c2-20260206-185521
 */

package vpc_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	acc "github.com/IBM-Cloud/terraform-provider-ibm/ibm/acctest"
	"github.com/IBM-Cloud/terraform-provider-ibm/ibm/service/vpc"
	"github.com/IBM/go-sdk-core/v5/core"
	"github.com/IBM/vpc-go-sdk/vpcv1"
	"github.com/stretchr/testify/assert"
)

func TestAccIBMIsPublicAddressRangeProfilesDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { acc.TestAccPreCheck(t) },
		Providers: acc.TestAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckIBMIsPublicAddressRangeProfilesDataSourceConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "id"),
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "profiles.#"),
				),
			},
		},
	})
}

func testAccCheckIBMIsPublicAddressRangeProfilesDataSourceConfigBasic() string {
	return fmt.Sprintf(`
		data "ibm_is_public_address_range_profiles" "is_public_address_range_profiles_instance" {
		}
	`)
}

func TestDataSourceIBMIsPublicAddressRangeProfilesPublicAddressRangeProfileToMap(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["family"] = "user"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/public_address_range/profiles/user-ipv4"
		model["ip_version"] = "ipv4"
		model["name"] = "user-ipv4"
		model["resource_type"] = "public_address_range_profile"

		assert.Equal(t, result, model)
	}

	model := new(vpcv1.PublicAddressRangeProfile)
	model.Family = core.StringPtr("user")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/public_address_range/profiles/user-ipv4")
	model.IPVersion = core.StringPtr("ipv4")
	model.Name = core.StringPtr("user-ipv4")
	model.ResourceType = core.StringPtr("public_address_range_profile")

	result, err := vpc.DataSourceIBMIsPublicAddressRangeProfilesPublicAddressRangeProfileToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
