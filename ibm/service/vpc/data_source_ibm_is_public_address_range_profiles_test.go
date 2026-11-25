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
					// user-ipv4 profile (index 0): targetable_resource_types should contain "vpc"
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "profiles.0.targetable_resource_types.#"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "profiles.0.targetable_resource_types.0.type", "enum"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "profiles.0.targetable_resource_types.0.values.0", "vpc"),
					// provider-ipv6 profile (index 1): targetable_resource_types should contain "virtual_network_interface"
					resource.TestCheckResourceAttrSet("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "profiles.1.targetable_resource_types.#"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "profiles.1.targetable_resource_types.0.type", "enum"),
					resource.TestCheckResourceAttr("data.ibm_is_public_address_range_profiles.is_public_address_range_profiles_instance", "profiles.1.targetable_resource_types.0.values.0", "virtual_network_interface"),
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
		model["targetable_resource_types"] = []map[string]interface{}{
			{
				"type":   "enum",
				"values": []string{"vpc"},
			},
		}

		assert.Equal(t, result, model)
	}

	trtType := "enum"
	trtValues := []string{"vpc"}
	model := new(vpcv1.PublicAddressRangeProfile)
	model.Family = core.StringPtr("user")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/public_address_range/profiles/user-ipv4")
	model.IPVersion = core.StringPtr("ipv4")
	model.Name = core.StringPtr("user-ipv4")
	model.ResourceType = core.StringPtr("public_address_range_profile")
	model.TargetableResourceTypes = &vpcv1.PublicAddressRangeProfileTargetableResourceTypes{
		Type:   &trtType,
		Values: trtValues,
	}

	result, err := vpc.DataSourceIBMIsPublicAddressRangeProfilesPublicAddressRangeProfileToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeProfilesToMapIPv6(t *testing.T) {
	checkResult := func(result map[string]interface{}) {
		model := make(map[string]interface{})
		model["family"] = "provider"
		model["href"] = "https://us-south.iaas.cloud.ibm.com/v1/public_address_range/profiles/provider-ipv6"
		model["ip_version"] = "ipv6"
		model["name"] = "provider-ipv6"
		model["resource_type"] = "public_address_range_profile"
		model["targetable_resource_types"] = []map[string]interface{}{
			{
				"type":   "enum",
				"values": []string{"virtual_network_interface"},
			},
		}

		assert.Equal(t, result, model)
	}

	trtType := "enum"
	trtValues := []string{"virtual_network_interface"}
	model := new(vpcv1.PublicAddressRangeProfile)
	model.Family = core.StringPtr("provider")
	model.Href = core.StringPtr("https://us-south.iaas.cloud.ibm.com/v1/public_address_range/profiles/provider-ipv6")
	model.IPVersion = core.StringPtr("ipv6")
	model.Name = core.StringPtr("provider-ipv6")
	model.ResourceType = core.StringPtr("public_address_range_profile")
	model.TargetableResourceTypes = &vpcv1.PublicAddressRangeProfileTargetableResourceTypes{
		Type:   &trtType,
		Values: trtValues,
	}

	result, err := vpc.DataSourceIBMIsPublicAddressRangeProfilesPublicAddressRangeProfileToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}

func TestDataSourceIBMIsPublicAddressRangeProfilesToMapNoTRT(t *testing.T) {
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
	// TargetableResourceTypes intentionally nil

	result, err := vpc.DataSourceIBMIsPublicAddressRangeProfilesPublicAddressRangeProfileToMap(model)
	assert.Nil(t, err)
	checkResult(result)
}
